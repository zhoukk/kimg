package kimg

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"

	"gopkg.in/gographics/imagick.v3/imagick"
)

var orientationNames = map[imagick.OrientationType]string{
	imagick.ORIENTATION_UNDEFINED:    "UNDEFINED",
	imagick.ORIENTATION_TOP_LEFT:     "TOP_LEFT",
	imagick.ORIENTATION_TOP_RIGHT:    "TOP_RIGHT",
	imagick.ORIENTATION_BOTTOM_RIGHT: "BOTTOM_RIGHT",
	imagick.ORIENTATION_BOTTOM_LEFT:  "BOTTOM_LEFT",
	imagick.ORIENTATION_LEFT_TOP:     "LEFT_TOP",
	imagick.ORIENTATION_RIGHT_TOP:    "RIGHT_TOP",
	imagick.ORIENTATION_RIGHT_BOTTOM: "RIGHT_BOTTOM",
	imagick.ORIENTATION_LEFT_BOTTOM:  "LEFT_BOTTOM",
}

var gravityMaps = map[string]imagick.GravityType{
	"nw": imagick.GRAVITY_NORTH_WEST,
	"n":  imagick.GRAVITY_NORTH,
	"ne": imagick.GRAVITY_NORTH_EAST,
	"w":  imagick.GRAVITY_WEST,
	"c":  imagick.GRAVITY_CENTER,
	"e":  imagick.GRAVITY_EAST,
	"sw": imagick.GRAVITY_SOUTH_WEST,
	"s":  imagick.GRAVITY_SOUTH,
	"se": imagick.GRAVITY_SOUTH_EAST,
}

// KimgImagick image processor struct hold kimg context.
type KimgImagick struct {
	ctx *KimgContext
}

// NewKimgImagick create a image processor instance and initialize imagick.
func NewKimgImagick(ctx *KimgContext) *KimgImagick {
	imagick.Initialize()

	return &KimgImagick{
		ctx: ctx,
	}
}

// Release terminate imagick.
func (image *KimgImagick) Release() {
	imagick.Terminate()
}

// Info get a image information and return a kimg response.
func (image *KimgImagick) Info(req *KimgRequest, data []byte) (*KimgResponse, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.PingImageBlob(data)
	if err != nil {
		return nil, err
	}

	width := uint(0)
	height := uint(0)
	format := mw.GetImageFormat()
	size, _ := mw.GetImageLength()
	orientationType := mw.GetImageOrientation()

	if "GIF" == format {
		width, height, _, _, _ = mw.GetImagePage()
	} else {
		width = mw.GetImageWidth()
		height = mw.GetImageHeight()
	}

	exif := make(map[string]string)
	names := mw.GetImageProperties("exif:*")
	for _, name := range names {
		exif[name] = mw.GetImageProperty(name)
	}

	u, _ := url.Parse(image.ctx.Config.Httpd.URL)
	u.Path = fmt.Sprintf("image/%s", req.Md5)
	return &KimgResponse{
		Md5:         req.Md5,
		URL:         u.String(),
		Style:       req.Key(),
		Size:        int(size),
		Width:       int(width),
		Height:      int(height),
		Format:      format,
		Orientation: orientationNames[orientationType],
		Exif:        exif,
	}, nil
}

// Convert convert a image according kimg request and return new image data.
func (image *KimgImagick) Convert(data []byte, req KimgRequest) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageBlob(data)
	if err != nil {
		image.ctx.Logger.Warn("ReadImageBlob err: %s", err)
		return nil, err
	}
	mw.ResetIterator()

	if "none" != req.Format {
		err = mw.SetImageFormat(strings.ToUpper(req.Format))
		if err != nil {
			image.ctx.Logger.Warn("SetImageFormat %s, err: %s", req.Format, err)
			return nil, err
		}
		image.ctx.Logger.Debug("SetImageFormat %s", req.Format)
	}

	if req.AutoOrient {
		if err := mw.AutoOrientImage(); err != nil {
			image.ctx.Logger.Warn("AutoOrientImage err: %s", err)
			return nil, err
		}
		image.ctx.Logger.Debug("AutoOrientImage")
	}

	if req.Strip {
		if err := mw.StripImage(); err != nil {
			image.ctx.Logger.Warn("StripImage err: %s", err)
			return nil, err
		}
		image.ctx.Logger.Debug("StripImage")
	}

	var newData []byte
	format := mw.GetImageFormat()
	if "GIF" == format {
		delay := mw.GetImageDelay()
		aw := mw.CoalesceImages()
		defer aw.Destroy()
		mw.Destroy()

		mw = imagick.NewMagickWand()
		mw.SetImageDelay(delay)
		for i := 0; i < int(aw.GetNumberImages()); i++ {
			aw.SetIteratorIndex(i)
			img := aw.GetImage()
			defer img.Destroy()
			if err = image.convertImage(img, &req); err == nil {
				mw.AddImage(img)
			}
		}
		mw.OptimizeImageLayers()
		mw.ResetIterator()
		newData = mw.GetImagesBlob()
	} else {
		if err = image.convertImage(mw, &req); err != nil {
			return nil, err
		}
		newData = mw.GetImageBlob()
	}

	if newData == nil || len(newData) == 0 {
		image.ctx.Logger.Warn("GetImageBlob failed")
		return nil, errors.New("GetImageBlob failed")
	}

	return newData, nil
}

func (image *KimgImagick) convertImage(mw *imagick.MagickWand, req *KimgRequest) error {
	if req.Scale {
		if err := image.scale(mw, req); err != nil {
			return err
		}
	}

	if req.Crop {
		if err := image.crop(mw, req); err != nil {
			return err
		}
	}

	if req.Rotate != 0 {
		background := imagick.NewPixelWand()
		defer background.Destroy()
		if len(req.BGColor) == 0 {
			req.BGColor = "transparent"
		}
		if !background.SetColor(req.BGColor) {
			image.ctx.Logger.Warn("background.SetColor %s failed", req.BGColor)
		} else {
			image.ctx.Logger.Debug("background.SetColor %s", req.BGColor)
		}
		if err := mw.RotateImage(background, float64(req.Rotate)); err != nil {
			image.ctx.Logger.Warn("RotateImage %f, err: %s", req.Rotate, err)
			return err
		}
		image.ctx.Logger.Debug("RotateImage %d %s", req.Rotate, req.BGColor)
	}

	if image.ctx.Config.Watermark.Enable {
		if err := image.waterMark(mw); err != nil {
			return err
		}
	}

	if req.Gray {
		if err := mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE); err != nil {
			image.ctx.Logger.Warn("SetImageType gray, err: %s", err)
			return err
		}
		image.ctx.Logger.Debug("SetImageType gray")
	}

	if req.Quality > 0 {
		if err := mw.SetImageCompressionQuality(uint(req.Quality)); err != nil {
			image.ctx.Logger.Warn("SetImageCompressionQuality %d, err: %s", req.Quality, err)
			return err
		}
		image.ctx.Logger.Debug("SetImageCompressionQuality %d", req.Quality)
	}
	return nil
}

func (image *KimgImagick) scale(mw *imagick.MagickWand, req *KimgRequest) error {
	w := mw.GetImageWidth()
	h := mw.GetImageHeight()

	if req.ScaleP > 0 {
		req.ScaleW = round(float64(w) * float64(req.ScaleP) / 100.0)
		req.ScaleH = round(float64(h) * float64(req.ScaleP) / 100.0)
	}
	if req.ScaleWP > 0 {
		req.ScaleW = round(float64(w) * float64(req.ScaleWP) / 100.0)
		req.ScaleH = int(h)
	}
	if req.ScaleHP > 0 {
		req.ScaleW = int(w)
		req.ScaleH = round(float64(h) * float64(req.ScaleHP) / 100.0)
	}

	if req.ScaleW > 0 && req.ScaleH == 0 {
		req.ScaleH = round(float64(req.ScaleW) * float64(h) / float64(w))
	} else if req.ScaleH > 0 && req.ScaleW == 0 {
		req.ScaleW = round(float64(req.ScaleH) * float64(w) / float64(h))
	} else if req.ScaleW > 0 && req.ScaleH > 0 {
		ratioW := float64(req.ScaleW) / float64(w)
		ratioH := float64(req.ScaleH) / float64(h)
		switch req.ScaleM {
		case "fit":
			{
				ratio := math.Min(ratioW, ratioH)
				req.ScaleW = round(float64(w) * ratio)
				req.ScaleH = round(float64(h) * ratio)
			}
		case "fill":
			{
				ratio := math.Max(ratioW, ratioH)
				req.ScaleW = round(float64(w) * ratio)
				req.ScaleH = round(float64(h) * ratio)
			}
		}
	}

	if req.ScaleW <= 0 {
		req.ScaleW = 1
	}
	if req.ScaleH <= 0 {
		req.ScaleH = 1
	}

	if err := mw.ResizeImage(uint(req.ScaleW), uint(req.ScaleH), imagick.FILTER_LANCZOS); err != nil {
		image.ctx.Logger.Warn("ResizeImage %d %d, err: %s", req.ScaleW, req.ScaleH, err)
		return err
	}
	image.ctx.Logger.Debug("ResizeImage %d %d", req.ScaleW, req.ScaleH)
	return nil
}

func (image *KimgImagick) crop(mw *imagick.MagickWand, req *KimgRequest) error {
	var w, h uint
	var x, y int

	w = mw.GetImageWidth()
	h = mw.GetImageHeight()

	if req.CropW <= 0 {
		req.CropW = int(w)
	}
	if req.CropH <= 0 {
		req.CropH = int(h)
	}

	switch req.Gravity {
	case "nw":
		{
			x = 0
			y = 0
		}
	case "n":
		{
			x = round(float64(w) / 2.0)
			y = 0
			x -= round(float64(req.CropW) / 2.0)
		}
	case "ne":
		{
			x = int(w)
			y = 0
			x -= req.CropW
		}
	case "w":
		{
			x = 0
			y = round(float64(h) / 2.0)
			y -= round(float64(req.CropH) / 2.0)
		}
	case "c":
		{
			x = round(float64(w) / 2.0)
			y = round(float64(h) / 2.0)
			x -= round(float64(req.CropW) / 2.0)
			y -= round(float64(req.CropH) / 2.0)
		}
	case "e":
		{
			x = int(w)
			y = round(float64(h / 2.0))
			x -= req.CropW
			y -= round(float64(req.CropH) / 2.0)
		}
	case "sw":
		{
			x = 0
			y = int(h)
			y -= req.CropH
		}
	case "s":
		{
			x = round(float64(w) / 2.0)
			y = int(h)
			x -= round(float64(req.CropW) / 2.0)
			y -= req.CropH
		}
	case "se":
		{
			x = int(w)
			y = int(h)
			x -= req.CropW
			y -= req.CropH
		}
	}

	switch req.Offset {
	case "lt":
		{
			x -= req.OffsetX
			y -= req.OffsetY
		}
	case "lb":
		{
			x -= req.OffsetX
			y += req.OffsetY
		}
	case "rt":
		{
			x += req.OffsetX
			y -= req.OffsetY
		}
	case "rb":
		{
			x += req.OffsetX
			y += req.OffsetY
		}
	}

	if err := mw.CropImage(uint(req.CropW), uint(req.CropH), x, y); err != nil {
		image.ctx.Logger.Warn("CropImage %d %d %d %d, err: %s", req.CropW, req.CropH, x, y, err)
		return err
	}
	if err := mw.SetImagePage(uint(req.CropW), uint(req.CropH), 0, 0); err != nil {
		image.ctx.Logger.Warn("SetImagePage %d %d %d %d, err: %s", req.CropW, req.CropH, 0, 0, err)
		return err
	}
	image.ctx.Logger.Debug("CropImage %d %d %d %d", req.CropW, req.CropH, x, y)
	return nil
}

func (image *KimgImagick) waterMark(mw *imagick.MagickWand) error {

	wm := image.ctx.Config.Watermark
	if len(wm.Logo.File) > 0 {
		logoMW := imagick.NewMagickWand()
		defer logoMW.Destroy()
		if err := logoMW.ReadImage(wm.Logo.File); err != nil {
			image.ctx.Logger.Warn("ReadImage %s, err: %s", wm.Logo.File, err)
			return err
		}

		dw := imagick.NewDrawingWand()
		pw := imagick.NewPixelWand()
		defer dw.Destroy()
		defer pw.Destroy()

		if gravity, ok := gravityMaps[wm.Gravity]; ok {
			dw.SetGravity(gravity)
			image.ctx.Logger.Debug("SetGravity %s", wm.Gravity)
		}
		if wm.Opacity > 0 {
			logoMW.SetImageAlpha(float64(wm.Opacity) / 100.0)
			image.ctx.Logger.Debug("SetImageAlpha %d", wm.Opacity)
		}
		if wm.Rotate > 0 {
			dw.Rotate(float64(wm.Rotate))
			image.ctx.Logger.Debug("Rotate %d", wm.Rotate)
		}

		if err := dw.Composite(imagick.COMPOSITE_OP_OVER, float64(wm.X), float64(wm.Y), float64(wm.Logo.W), float64(wm.Logo.H), logoMW); err != nil {
			image.ctx.Logger.Warn("Composite %d %d %d %d, err: %s", wm.X, wm.Y, wm.Logo.W, wm.Logo.H, err)
			return err
		}
		image.ctx.Logger.Debug("Composite %d %d %d %d", wm.X, wm.Y, wm.Logo.W, wm.Logo.H)

		if err := mw.DrawImage(dw); err != nil {
			image.ctx.Logger.Warn("DrawImage err: %s", err)
			return err
		}
		image.ctx.Logger.Debug("DrawImage")
	}

	if len(wm.Text.Content) > 0 {
		dw := imagick.NewDrawingWand()
		pw := imagick.NewPixelWand()
		defer dw.Destroy()
		defer pw.Destroy()

		if gravity, ok := gravityMaps[wm.Gravity]; ok {
			dw.SetGravity(gravity)
			image.ctx.Logger.Debug("SetGravity %s", wm.Gravity)
		}
		if len(wm.Text.FontName) > 0 {
			if err := dw.SetFont(wm.Text.FontName); err != nil {
				image.ctx.Logger.Warn("SetFont %s, err: %s", wm.Text.FontName, err)
			} else {
				image.ctx.Logger.Debug("SetFont %s", wm.Text.FontName)
			}
		}
		if wm.Text.FontSize > 0 {
			dw.SetFontSize(float64(wm.Text.FontSize))
			image.ctx.Logger.Debug("SetFontSize %d", wm.Text.FontSize)
		}
		if len(wm.Text.FontColor) > 0 {
			if pw.SetColor(wm.Text.FontColor) {
				image.ctx.Logger.Debug("SetAlpha %d", wm.Opacity)
				dw.SetFillColor(pw)
				image.ctx.Logger.Debug("SetFillColor %s", wm.Text.FontColor)
			} else {
				image.ctx.Logger.Warn("SetFillColor %s, err", wm.Text.FontColor)
			}
		}
		dw.SetFillOpacity(float64(wm.Opacity) / 100.0)
		image.ctx.Logger.Debug("SetFillOpacity %d", wm.Opacity)
		if wm.Text.StrokeWidth > 0 {
			dw.SetStrokeWidth(float64(wm.Text.StrokeWidth))
			image.ctx.Logger.Debug("SetStrokeWidth %d", wm.Text.StrokeWidth)
			if len(wm.Text.StrokeColor) > 0 {
				if pw.SetColor(wm.Text.StrokeColor) {
					dw.SetStrokeColor(pw)
					image.ctx.Logger.Debug("SetStrokeColor %s", wm.Text.StrokeColor)
				} else {
					image.ctx.Logger.Warn("SetStrokeColor %s, err", wm.Text.StrokeColor)
				}
			}
			dw.SetStrokeOpacity(float64(wm.Opacity) / 100.0)
			image.ctx.Logger.Debug("SetStrokeOpacity %d", wm.Opacity)
		}
		if err := mw.AnnotateImage(dw, float64(wm.X), float64(wm.Y), float64(wm.Rotate), wm.Text.Content); err != nil {
			image.ctx.Logger.Warn("AnnotateImage %d %d %d %s, err: %s", wm.X, wm.Y, wm.Rotate, wm.Text.Content, err)
			return err
		}
		image.ctx.Logger.Debug("AnnotateImage %d %d %d %s", wm.X, wm.Y, wm.Rotate, wm.Text.Content)
	}
	return nil
}

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}
