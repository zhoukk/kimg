package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var contentTypes = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
}

func (ctx *KimgContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mux := http.NewServeMux()

	mux.HandleFunc("/image", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				ctx.post(w, r)
			}
		}
	}))

	mux.HandleFunc("/image/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md5Sum := r.URL.Path[7:len(r.URL.Path)]
		if !ctx.isValidMd5(md5Sum) {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case "GET":
			{
				ctx.get(w, r, md5Sum)
			}
		case "DELETE":
			{
				ctx.delete(w, r, md5Sum)
			}
		}
	}))

	mux.HandleFunc("/info/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md5Sum := r.URL.Path[6:len(r.URL.Path)]
		if !ctx.isValidMd5(md5Sum) {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case "GET":
			{
				ctx.info(w, r, md5Sum)
			}
		}
	}))

	if ctx.Config.Httpd.EnableWeb {
		mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				{
					ctx.www(w, r)
				}
			}
		}))
	}

	mux.ServeHTTP(w, r)
}

func (ctx *KimgContext) www(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./www")).ServeHTTP(w, r)
}

func (ctx *KimgContext) post(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > ctx.Config.Httpd.MaxSize {
		http.Error(w, "Payload Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	var rd io.Reader
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		rd = io.LimitReader(r.Body, ctx.Config.Httpd.MaxSize)
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(ctx.Config.Httpd.MaxSize); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file, _, err := r.FormFile(ctx.Config.Httpd.FormName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rd = file
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(rd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileType := http.DetectContentType(data)
	if !ctx.isAllowedType(fileType) {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	resp, err := ctx.SaveImage(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)

	ctx.Logger.Info("POST md5: %s, size: %d", resp.Md5, resp.Size)
}

func (ctx *KimgContext) info(w http.ResponseWriter, r *http.Request, md5Sum string) {
	if err := r.ParseForm(); err != nil {
		ctx.Logger.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := ctx.genRequest(r, md5Sum)

	resp, err := ctx.InfoImage(req)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)

	ctx.Logger.Info("INFO md5: %s", md5Sum)
}

func (ctx *KimgContext) get(w http.ResponseWriter, r *http.Request, md5Sum string) {
	if err := r.ParseForm(); err != nil {
		ctx.Logger.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := ctx.genRequest(r, md5Sum)

	data, err := ctx.GetImage(req)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if ctx.Config.Httpd.MaxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ctx.Config.Httpd.MaxAge))
	}

	headers := ctx.Config.Httpd.Headers
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.Header().Set("X-Kimg-Style", req.Key())

	if ctx.Config.Httpd.Etag {
		m := md5.New()
		m.Write(data)
		newMd5 := hex.EncodeToString(m.Sum(nil))

		if ifNoneMatch, ok := r.Header["If-None-Match"]; ok {
			if ifNoneMatch[0] == newMd5 || ifNoneMatch[0] == "*" {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Header().Set("Etag", newMd5)
	}

	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(data))

	ctx.Logger.Info("GET %s, size: %d", r.RequestURI, len(data))
}

func (ctx *KimgContext) delete(w http.ResponseWriter, r *http.Request, md5Sum string) {
	err := ctx.DeleteImage(md5Sum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Info("DELETE md5: %s", md5Sum)
}

func (ctx *KimgContext) isAllowedType(fileType string) bool {
	types := ctx.Config.Image.AllowedTypes
	for _, t := range types {
		if strings.Contains(fileType, t) {
			return true
		}
	}
	return false
}

func (ctx *KimgContext) isValidMd5(md5 string) bool {
	return regexp.MustCompile(`^([0-9a-zA-Z]){32}$`).MatchString(md5)
}

func (ctx *KimgContext) genRequest(r *http.Request, md5Sum string) *KimgRequest {
	var req KimgRequest

	req.Md5 = md5Sum

	if v, ok := r.Form["origin"]; ok {
		req.Origin = v[0] != "0"
		return &req
	}

	if v, ok := r.Form["style"]; ok {
		req.Style = v[0]
		return &req
	}

	if v, ok := r.Form["save"]; ok {
		req.Save = v[0] != "0"
	} else {
		req.Save = ctx.Config.Storage.SaveNew
	}

	if v, ok := r.Form["s"]; ok {
		req.Scale = v[0] != "0"
	}

	if req.Scale {
		if v, ok := r.Form["sm"]; ok {
			req.ScaleM = v[0]
		}
		if v, ok := r.Form["sw"]; ok {
			req.ScaleW, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["sh"]; ok {
			req.ScaleH, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["sp"]; ok {
			req.ScaleP, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["swp"]; ok {
			req.ScaleWP, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["shp"]; ok {
			req.ScaleHP, _ = strconv.Atoi(v[0])
		}
	}

	if v, ok := r.Form["c"]; ok {
		req.Crop = v[0] != "0"
	}

	if req.Crop {
		if v, ok := r.Form["cg"]; ok {
			req.Gravity = v[0]
		}

		if v, ok := r.Form["cw"]; ok {
			req.CropW, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["ch"]; ok {
			req.CropH, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["co"]; ok {
			req.Offset = v[0]
		}
		if v, ok := r.Form["cx"]; ok {
			req.OffsetX, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["cy"]; ok {
			req.OffsetY, _ = strconv.Atoi(v[0])
		}
	}

	if v, ok := r.Form["t"]; ok {
		req.Text = v[0]
	}
	if len(req.Text) > 0 {
		if v, ok := r.Form["ts"]; ok {
			req.FontSize, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["tw"]; ok {
			req.FontWeight, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["tc"]; ok && len(v[0]) == 6 {
			req.FontColor = "#" + v[0]
		}
		if v, ok := r.Form["tsc"]; ok && len(v[0]) == 6 {
			req.StrokeColor = "#" + v[0]
		}
		if v, ok := r.Form["tsw"]; ok {
			req.StrokeWidth, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["tg"]; ok {
			req.TextGravity = v[0]
		}
		if v, ok := r.Form["tx"]; ok {
			req.TextX, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["ty"]; ok {
			req.TextY, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["tr"]; ok {
			req.TextRotate, _ = strconv.Atoi(v[0])
		}
		if v, ok := r.Form["to"]; ok {
			req.TextOpacity, _ = strconv.Atoi(v[0])
		}
	}

	if v, ok := r.Form["f"]; ok {
		req.Format = strings.ToLower(v[0])
		if !ctx.isAllowedType(req.Format) {
			req.Format = ctx.Config.Image.Format
		}
	} else {
		req.Format = ctx.Config.Image.Format
	}

	if v, ok := r.Form["q"]; ok {
		req.Quality, _ = strconv.Atoi(v[0])
	}
	if req.Quality <= 0 {
		req.Quality = ctx.Config.Image.Quality
	} else if req.Quality > 100 {
		req.Quality = 100
	}

	if v, ok := r.Form["r"]; ok {
		req.Rotate, _ = strconv.Atoi(v[0])
	}
	if v, ok := r.Form["bc"]; ok && len(v[0]) == 6 {
		req.BGColor = "#" + v[0]
	}
	if v, ok := r.Form["g"]; ok {
		req.Gray = v[0] != "0"
	}
	if v, ok := r.Form["ao"]; ok {
		req.AutoOrient = v[0] != "0"
	} else {
		req.AutoOrient = true
	}
	if v, ok := r.Form["st"]; ok {
		req.Strip = v[0] != "0"
	} else {
		req.Strip = true
	}
	return &req
}
