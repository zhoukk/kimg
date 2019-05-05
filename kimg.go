package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
)

var (
	// KimgVersion version of Kimg
	KimgVersion = "latest"

	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "kimg.yaml", "config file")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, err := NewKimgContext(configFile)
	if err != nil {
		log.Println(err)
		flag.Usage()
		return
	}
	defer ctx.Release()

	log.Printf("[INFO] kimg#%s start at %s\n", KimgVersion, ctx.Config.Httpd.Bind)
	log.Fatalln(http.ListenAndServe(ctx.Config.Httpd.Bind, ctx))
}
