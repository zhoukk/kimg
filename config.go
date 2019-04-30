package main

import (
	"log"
	"os"
	"strconv"

	gcfg "gopkg.in/gcfg.v1"
)

// KimgConfig is configuration of kimg.
type KimgConfig struct {
	Httpd struct {
		Host      string
		Port      int
		Headers   string
		Etag      int
		MaxAge    int
		FormName  string
		MaxSize   int64
		EnableWeb int
	}

	Image struct {
		Format       string
		Quality      int
		AllowedTypes string
	}

	Logger struct {
		Mode  string
		Level string
		File  string
	}

	Cache struct {
		Mode         string
		MaxSize      int
		MemcacheHost string
		MemcachePort int
		RedisHost    string
		RedisPort    int
	}

	Storage struct {
		Mode    string
		SaveNew int
		Root    string
	}
}

// NewKimgConfig create a config instance from config file.
func NewKimgConfig(configFile string) (*KimgConfig, error) {
	var cfg KimgConfig

	cfg.Httpd.Host = "0.0.0.0"
	cfg.Httpd.Port = 80
	cfg.Httpd.Headers = "Server:kimg"
	cfg.Httpd.Etag = 1
	cfg.Httpd.MaxAge = 7776000
	cfg.Httpd.FormName = "file"
	cfg.Httpd.MaxSize = 104857600
	cfg.Httpd.EnableWeb = 1

	cfg.Image.Format = "jpeg"
	cfg.Image.Quality = 75
	cfg.Image.AllowedTypes = "jpeg,jpg,png,gif,webp"

	cfg.Logger.Mode = "console"
	cfg.Logger.Level = "debug"
	cfg.Logger.File = "kimg.log"

	cfg.Cache.Mode = "none"
	cfg.Cache.MaxSize = 1048576
	cfg.Cache.MemcacheHost = "127.0.0.1"
	cfg.Cache.MemcachePort = 11211
	cfg.Cache.RedisHost = "127.0.0.1"
	cfg.Cache.RedisPort = 6379

	cfg.Storage.Mode = "file"
	cfg.Storage.SaveNew = 1
	cfg.Storage.Root = "kimgs"

	err := gcfg.ReadFileInto(&cfg, configFile)
	if err != nil {
		log.Printf("[ERROR] %s\n", err)
		log.Println("[INFO] configuration [default] used")
	} else {
		log.Printf("[INFO] configuration [%s] used\n", configFile)
	}

	// httpd env
	if env, ok := os.LookupEnv("KIMG_HTTPD_HOST"); ok {
		cfg.Httpd.Host = env
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_PORT"); ok {
		cfg.Httpd.Port, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_HEADERS"); ok {
		cfg.Httpd.Headers = env
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_ETAG"); ok {
		cfg.Httpd.Etag, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_MAX_AGE"); ok {
		cfg.Httpd.MaxAge, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_FORM_NAME"); ok {
		cfg.Httpd.FormName = env
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_MAX_SIZE"); ok {
		cfg.Httpd.MaxSize, _ = strconv.ParseInt(env, 0, 64)
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_ENABLE_WEB"); ok {
		cfg.Httpd.EnableWeb, _ = strconv.Atoi(env)
	}

	// image env
	if env, ok := os.LookupEnv("KIMG_IMAGE_FORMAT"); ok {
		cfg.Image.Format = env
	}
	if env, ok := os.LookupEnv("KIMG_IMAGE_QUALITY"); ok {
		cfg.Image.Quality, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_IMAGE_ALLOWED_TYPES"); ok {
		cfg.Image.AllowedTypes = env
	}

	// logger env
	if env, ok := os.LookupEnv("KIMG_LOGGER_MODE"); ok {
		cfg.Logger.Mode = env
	}
	if env, ok := os.LookupEnv("KIMG_LOGGER_LEVEL"); ok {
		cfg.Logger.Level = env
	}
	if env, ok := os.LookupEnv("KIMG_LOGGER_FILE"); ok {
		cfg.Logger.File = env
	}

	// cache env
	if env, ok := os.LookupEnv("KIMG_CACHE_MODE"); ok {
		cfg.Cache.Mode = env
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_MAX_SIZE"); ok {
		cfg.Cache.MaxSize, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_MEMCACHE_HOST"); ok {
		cfg.Cache.MemcacheHost = env
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_MEMCACHE_PORT"); ok {
		cfg.Cache.MemcachePort, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_REDIS_HOST"); ok {
		cfg.Cache.RedisHost = env
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_REDIS_PORT"); ok {
		cfg.Cache.RedisPort, _ = strconv.Atoi(env)
	}

	// storage env
	if env, ok := os.LookupEnv("KIMG_STORAGE_MODE"); ok {
		cfg.Storage.Mode = env
	}
	if env, ok := os.LookupEnv("KIMG_STORAGE_SAVE_NEW"); ok {
		cfg.Storage.SaveNew, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_STORAGE_ROOT"); ok {
		cfg.Storage.Root = env
	}

	return &cfg, nil
}
