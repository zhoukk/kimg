package main

import (
	"flag"
	"fmt"
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
	flag.StringVar(&configFile, "c", "kimg.ini", "config file")
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

	addr := fmt.Sprintf("%s:%d", ctx.Config.Httpd.Host, ctx.Config.Httpd.Port)
	log.Printf("[INFO] kimg#%s start at %s\n", KimgVersion, addr)
	log.Fatalln(http.ListenAndServe(addr, ctx))
}
