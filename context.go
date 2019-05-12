package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// KimgRequest define a image request.
type KimgRequest struct {
	Md5    string `json:"-"`
	Origin bool   `json:"-"`
	Style  string `json:"-"`
	Save   bool   `json:"-"`

	// scale params
	Scale   bool   `json:"scale,omitempty"`
	ScaleM  string `json:"scale_m,omitempty"`
	ScaleW  int    `json:"scale_w,omitempty"`
	ScaleH  int    `json:"scale_h,omitempty"`
	ScaleWP int    `json:"scale_wp,omitempty"`
	ScaleHP int    `json:"scale_hp,omitempty"`
	ScaleP  int    `json:"scale_p,omitempty"`

	// crop params
	Crop    bool   `json:"crop,omitempty"`
	Gravity string `json:"gravity,omitempty"`
	CropW   int    `json:"crop_w,omitempty"`
	CropH   int    `json:"crop_h,omitempty"`
	Offset  string `json:"offset,omitempty"`
	OffsetX int    `json:"offset_x,omitempty"`
	OffsetY int    `json:"offset_y,omitempty"`

	// watermark
	Text        string `json:"text,omitempty"`
	FontSize    int    `json:"font_size,omitempty"`
	FontWeight  int    `json:"font_weight,omitempty"`
	FontColor   string `json:"font_color,omitempty"`
	StrokeColor string `json:"stroke_color,omitempty"`
	StrokeWidth int    `json:"stroke_width,omitempty"`
	TextGravity string `json:"text_gravity,omitempty"`
	TextX       int    `json:"text_x,omitempty"`
	TextY       int    `json:"text_y,omitempty"`
	TextRotate  int    `json:"text_rotate,omitempty"`
	TextOpacity int    `json:"text_opacity,omitempty"`

	Format     string `json:"format,omitempty"`
	Quality    int    `json:"quality,omitempty"`
	Rotate     int    `json:"rotate,omitempty"`
	BGColor    string `json:"bg_color,omitempty"`
	Gray       bool   `json:"gray,omitempty"`
	AutoOrient bool   `json:"auto_orient,omitempty"`
	Strip      bool   `json:"strip,omitempty"`
}

// KimgResponse define a image response.
type KimgResponse struct {
	Md5         string            `json:"md5"`
	URL         string            `json:"url"`
	Style       string            `json:"style,omitempty"`
	Size        int               `json:"size"`
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	Format      string            `json:"format"`
	Orientation string            `json:"orientation"`
	Exif        map[string]string `json:"exif"`
}

// KimgContext context of kimg.
type KimgContext struct {
	Config  *KimgConfig
	Cache   KimgCache
	Logger  KimgLogger
	Storage KimgStorage
	Image   *KimgImagick
}

// Key generate a key according to image style request params.
func (req *KimgRequest) Key() string {
	if req.Origin {
		return req.Md5
	}
	if len(req.Style) > 0 {
		return req.Style
	}
	b, _ := json.Marshal(req)
	m := md5.New()
	m.Write(b)
	return hex.EncodeToString(m.Sum(nil))
}

// NewKimgContext create a instance of kimg context.
func NewKimgContext(configFile string) (*KimgContext, error) {
	var ctx KimgContext

	config, err := NewKimgConfig(configFile)
	if err != nil {
		return nil, err
	}
	ctx.Config = config

	logger, err := NewKimgLogger(config)
	if err != nil {
		return nil, err
	}
	ctx.Logger = logger

	cache, err := NewKimgCache(config)
	if err != nil {
		return nil, err
	}
	ctx.Cache = cache

	storage, err := NewKimgStorage(&ctx)
	if err != nil {
		return nil, err
	}
	ctx.Storage = storage

	ctx.Image = NewKimgImagick(&ctx)

	return &ctx, nil
}

// Release release resource in kimg context.
func (ctx *KimgContext) Release() {
	ctx.Image.Release()
}

// SaveImage save a image to kimg and make a kimg response.
func (ctx *KimgContext) SaveImage(data []byte) (*KimgResponse, error) {
	m := md5.New()
	m.Write(data)
	md5Sum := hex.EncodeToString(m.Sum(nil))

	ctx.Logger.Debug("SaveImage md5Sum: %s", md5Sum)

	req := ctx.originRequest(md5Sum)

	err := ctx.Storage.Set(req, data)
	if err != nil {
		return nil, err
	}

	if ctx.isCacheEnable(data) {
		cacheKey := ctx.cacheKey(req)
		if err = ctx.Cache.Set(cacheKey, data); err != nil {
			ctx.Logger.Warn("SaveImage md5Sum: %s, SetCache %s err: %s", md5Sum, cacheKey, err)
		} else {
			ctx.Logger.Debug("SaveImage md5Sum: %s, SetCache %s", md5Sum, cacheKey)
		}
	}

	resp, err := ctx.Image.Info(req, data)
	if err != nil {
		ctx.Logger.Warn("SaveImage md5Sum: %s, Image.Info err: %s", md5Sum, err)
		return nil, err
	}

	return resp, nil
}

// GetImage get a image data from kimg according to a image request.
func (ctx *KimgContext) GetImage(req *KimgRequest) ([]byte, error) {

	ctx.Logger.Debug("GetImage md5Sum: %s, req: %#v", req.Md5, req)

	cacheKey := ctx.cacheKey(req)

	if ctx.isCacheEnable(nil) {
		data, err := ctx.Cache.Get(cacheKey)
		if err == nil {
			ctx.Logger.Debug("GetImage md5Sum: %s, GetCache %s", req.Md5, cacheKey)
			return data, nil
		}
		ctx.Logger.Debug("GetImage md5Sum: %s, GetCache %s err: %s", req.Md5, cacheKey, err)
	}

	data, err := ctx.Storage.Get(req)
	if err == nil {
		if ctx.isCacheEnable(data) {
			if err = ctx.Cache.Set(cacheKey, data); err != nil {
				ctx.Logger.Warn("GetImage md5Sum: %s, SetCache %s err: %s", req.Md5, cacheKey, err)
			} else {
				ctx.Logger.Debug("GetImage md5Sum: %s, SetCache %s", req.Md5, cacheKey)
			}
		}
		return data, nil
	}

	var originData []byte
	saveToCache := false
	originReq := ctx.originRequest(req.Md5)
	if ctx.isCacheEnable(nil) {
		originCacheKey := ctx.cacheKey(originReq)
		originData, err = ctx.Cache.Get(originCacheKey)
		if err != nil {
			saveToCache = true
			ctx.Logger.Debug("GetImage md5Sum: %s, GetOriginCache %s err: %s", req.Md5, originCacheKey, err)
		} else {
			ctx.Logger.Debug("GetImage md5Sum: %s, GetOriginCache %s", req.Md5, originCacheKey)
		}
	}

	if originData == nil {
		originData, err = ctx.Storage.Get(originReq)
		if err != nil {
			ctx.Logger.Warn("GetImage md5Sum: %s, GetStorage err: %s", req.Md5, err)
			return nil, err
		} else if ctx.isCacheEnable(originData) && saveToCache {
			originCacheKey := ctx.cacheKey(originReq)
			if err = ctx.Cache.Set(originCacheKey, originData); err != nil {
				ctx.Logger.Warn("GetImage md5Sum: %s, SetOriginCache %s err: %s", req.Md5, originCacheKey, err)
			} else {
				ctx.Logger.Debug("GetImage md5Sum: %s, SetOriginCache %s", req.Md5, originCacheKey)
			}
		}
	}

	data, err = ctx.Image.Convert(originData, *req)
	if err != nil {
		ctx.Logger.Warn("GetImage md5Sum: %s, Image.Convert err: %s", req.Md5, err)
		return nil, err
	}

	if ctx.isCacheEnable(nil) {
		if err = ctx.Cache.Set(cacheKey, data); err != nil {
			ctx.Logger.Warn("GetImage md5Sum: %s, SetCache %s err: %s", req.Md5, cacheKey, err)
		} else {
			ctx.Logger.Debug("GetImage md5Sum: %s, SetCache %s", req.Md5, cacheKey)
		}
	}

	if req.Save {
		if err := ctx.Storage.Set(req, data); err != nil {
			ctx.Logger.Warn("GetImage md5Sum: %s, save new image err :%s", req.Md5, err)
		}
	}

	return data, nil
}

// InfoImage get a image information according the md5 key and make a image response.
func (ctx *KimgContext) InfoImage(req *KimgRequest) (*KimgResponse, error) {

	ctx.Logger.Debug("InfoImage md5Sum: %s", req.Md5)

	var data []byte
	var err error

	cacheKey := ctx.cacheKey(req)
	saveToCache := false

	if ctx.isCacheEnable(nil) {
		data, err = ctx.Cache.Get(cacheKey)
		if err != nil {
			saveToCache = true
			ctx.Logger.Debug("InfoImage md5Sum: %s, GetCache %s err: %s", req.Md5, cacheKey, err)
		}
	}

	if data == nil {
		data, err = ctx.Storage.Get(req)
		if err != nil {
			ctx.Logger.Warn("InfoImage md5Sum: %s, GetStorage err: %s", req.Md5, err)
			return nil, err
		} else if ctx.isCacheEnable(data) && saveToCache {
			if err = ctx.Cache.Set(cacheKey, data); err != nil {
				ctx.Logger.Warn("InfoImage md5Sum: %s, SetCache %s err: %s", req.Md5, cacheKey, err)
			} else {
				ctx.Logger.Debug("InfoImage md5Sum: %s, SetCache %s", req.Md5, cacheKey)
			}
		}
	}

	resp, err := ctx.Image.Info(req, data)
	if err != nil {
		ctx.Logger.Warn("InfoImage md5Sum: %s, Image.Info err: %s", req.Md5, err)
		return nil, err
	}

	return resp, err
}

// DeleteImage delete a image from kimg according the md5 key.
func (ctx *KimgContext) DeleteImage(md5Sum string) error {

	ctx.Logger.Debug("DeleteImage md5Sum: %s", md5Sum)

	req := ctx.originRequest(md5Sum)

	if ctx.isCacheEnable(nil) {
		cacheKey := ctx.cacheKey(req)
		err := ctx.Cache.Del(cacheKey)
		if err != nil {
			ctx.Logger.Warn("DeleteImage md5Sum: %s, DelCache %s err: %s", req.Md5, cacheKey, err)
		} else {
			ctx.Logger.Debug("DeleteImage md5Sum: %s, DelCache %s", req.Md5, cacheKey)
		}
	}

	err := ctx.Storage.Del(req)
	if err != nil {
		ctx.Logger.Warn("DeleteImage md5Sum: %s, DelStorage err: %s", req.Md5, err)
		return err
	}

	return nil
}

func (ctx *KimgContext) isCacheEnable(data []byte) bool {
	return ctx.Cache != nil && (data != nil || ctx.Config.Cache.MaxSize >= len(data))
}

func (ctx *KimgContext) cacheKey(req *KimgRequest) string {
	if req.Origin {
		return req.Md5
	}
	return fmt.Sprintf("%s:%s", req.Md5, req.Key())
}

func (ctx *KimgContext) originRequest(md5Sum string) *KimgRequest {
	return &KimgRequest{Md5: md5Sum, Origin: true}
}
