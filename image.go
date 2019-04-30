package main

import (
	"errors"
	"math"

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

	size, _ := mw.GetImageLength()
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	format := mw.GetImageFormat()
	orientationType := mw.GetImageOrientation()

	exif := make(map[string]string)
	names := mw.GetImageProperties("*")
	for _, name := range names {
		exif[name] = mw.GetImageProperty(name)
	}

	return &KimgResponse{
		Md5:         req.Md5,
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

	if req.Scale {
		if err = image.scale(mw, &req); err != nil {
			return nil, err
		}
	}

	if req.Crop {
		if err = image.crop(mw, &req); err != nil {
			return nil, err
		}
	}

	if req.Rotate != 0 {
		background := imagick.NewPixelWand()
		defer background.Destroy()
		if len(req.BGColor) > 0 && !background.SetColor(req.BGColor) {
			image.ctx.Logger.Warn("background.SetColor %s failed", req.BGColor)
			return nil, errors.New("background.SetColor failed")
		}
		err = mw.RotateImage(background, float64(req.Rotate))
		if err != nil {
			image.ctx.Logger.Warn("RotateImage %f, err: %s", req.Rotate, err)
			return nil, err
		}
		image.ctx.Logger.Debug("RotateImage %d #%s", req.Rotate, req.BGColor)
	}

	if len(req.Text) > 0 {
		if err = image.waterMark(mw, &req); err != nil {
			return nil, err
		}
	}

	if req.Gray {
		err = mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE)
		if err != nil {
			image.ctx.Logger.Warn("SetImageType gray, err: %s", err)
			return nil, err
		}
		image.ctx.Logger.Debug("SetImageType gray")
	}

	err = mw.SetImageCompressionQuality(uint(req.Quality))
	if err != nil {
		image.ctx.Logger.Warn("SetImageCompressionQuality %d, err: %s", req.Quality, err)
		return nil, err
	}
	image.ctx.Logger.Debug("SetImageCompressionQuality %d", req.Quality)

	err = mw.SetImageFormat(req.Format)
	if err != nil {
		image.ctx.Logger.Warn("SetImageFormat %s, err: %s", req.Format, err)
		return nil, err
	}
	image.ctx.Logger.Debug("SetImageFormat %s", req.Format)

	newData := mw.GetImageBlob()
	if newData == nil || len(newData) == 0 {
		image.ctx.Logger.Warn("GetImageBlob failed")
		return nil, errors.New("GetImageBlob failed")
	}

	return newData, nil
}

func (image *KimgImagick) scale(mw *imagick.MagickWand, req *KimgRequest) error {
	var w, h uint

	w = mw.GetImageWidth()
	h = mw.GetImageHeight()

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
		x += req.OffsetX
		y += req.OffsetY
	case "lb":
		x -= req.OffsetX
		y += req.OffsetY
	case "rt":
		x += req.OffsetX
		y -= req.OffsetY
	case "rb":
		x -= req.OffsetX
		y -= req.OffsetY
	}

	if err := mw.CropImage(uint(req.CropW), uint(req.CropH), x, y); err != nil {
		image.ctx.Logger.Warn("CropImage %d %d %d %d, err: %s", req.CropW, req.CropH, x, y, err)
		return err
	}
	image.ctx.Logger.Debug("CropImage %d %d %d %d", req.CropW, req.CropH, x, y)
	return nil
}

func (image *KimgImagick) waterMark(mw *imagick.MagickWand, req *KimgRequest) error {
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()
	defer dw.Destroy()
	defer pw.Destroy()
	if req.FontSize > 0 {
		dw.SetFontSize(float64(req.FontSize))
	}
	if req.FontWeight >= 100 && req.FontWeight <= 900 {
		dw.SetFontWeight(uint(req.FontWeight))
	}
	if len(req.FontColor) > 0 {
		pw.SetColor(req.FontColor)
		if req.TextOpacity > 0 {
			pw.SetOpacity(float64(req.TextOpacity) / 100.0)
		}
		dw.SetFillColor(pw)
	}
	if len(req.StrokeColor) > 0 {
		pw.SetColor(req.StrokeColor)
		dw.SetStrokeColor(pw)
	}
	dw.SetStrokeWidth(float64(req.StrokeWidth))
	if gravity, ok := gravityMaps[req.TextGravity]; ok {
		dw.SetGravity(gravity)
	}
	if err := mw.AnnotateImage(dw, float64(req.TextX), float64(req.TextY), float64(req.TextRotate), req.Text); err != nil {
		image.ctx.Logger.Warn("AnnotateImage %d %d %d, err: %s", req.TextX, req.TextY, req.TextRotate, err)
		return err
	}
	image.ctx.Logger.Debug("AnnotateImage %s %d %d %d", req.Text, req.TextX, req.TextY, req.TextRotate)
	return nil
}

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}
