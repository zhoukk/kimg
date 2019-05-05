package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// KimgConfig is configuration of kimg.
type KimgConfig struct {
	Httpd struct {
		Bind      string            `yaml:"bind,omitempty"`
		Headers   map[string]string `yaml:"headers,omitempty"`
		Etag      bool              `yaml:"etag,omitempty"`
		MaxAge    int               `yaml:"maxAge,omitempty"`
		FormName  string            `yaml:"formName,omitempty"`
		MaxSize   int64             `yaml:"maxSize,omitempty"`
		EnableWeb bool              `yaml:"enableWeb,omitempty"`
	} `yaml:"httpd,omitempty"`

	Image struct {
		Format       string   `yaml:"format,omitempty"`
		Quality      int      `yaml:"quality,omitempty"`
		AllowedTypes []string `yaml:"allowedTypes,omitempty"`
	} `yaml:"image,omitempty"`

	Logger struct {
		Mode  string `yaml:"mode,omitempty"`
		Level string `yaml:"level,omitempty"`
		File  string `yaml:"file,omitempty"`
	} `yaml:"logger,omitempty"`

	Cache struct {
		Mode     string `yaml:"mode,omitempty"`
		MaxSize  int    `yaml:"maxSize,omitempty"`
		Memcache struct {
			URL string `yaml:"url,omitempty"`
		} `yaml:"memcache,omitempty"`
		Redis struct {
			URL string `yaml:"url,omitempty"`
		} `yaml:"redis,omitempty"`
		Memory struct {
			Capacity int64 `yaml:"capacity,omitempty"`
		} `yaml:"memory,omitempty"`
	} `yaml:"cache,omitempty"`

	Storage struct {
		Mode    string `yaml:"mode,omitempty"`
		SaveNew bool   `yaml:"saveNew,omitempty"`
		Root    string `yaml:"root,omitempty"`
	} `yaml:"storage,omitempty"`
}

// NewKimgConfig create a config instance from config file.
func NewKimgConfig(configFile string) (*KimgConfig, error) {
	var cfg KimgConfig

	cfg.Httpd.Bind = "0.0.0.0:80"
	cfg.Httpd.Headers = map[string]string{"Server": "kimg"}
	cfg.Httpd.Etag = true
	cfg.Httpd.MaxAge = 90 * 24 * 3600
	cfg.Httpd.FormName = "file"
	cfg.Httpd.MaxSize = 100 * 1024 * 1024
	cfg.Httpd.EnableWeb = true

	cfg.Image.Format = "jpeg"
	cfg.Image.Quality = 75
	cfg.Image.AllowedTypes = []string{"jpeg", "jpg", "png", "gif", "webp"}

	cfg.Logger.Mode = "console"
	cfg.Logger.Level = "debug"
	cfg.Logger.File = "kimg.log"

	cfg.Cache.Mode = "memory"
	cfg.Cache.MaxSize = 1 * 1024 * 1024
	cfg.Cache.Memcache.URL = "127.0.0.1:11211"
	cfg.Cache.Redis.URL = "127.0.0.1:6379"
	cfg.Cache.Memory.Capacity = 100 * 1024 * 1024

	cfg.Storage.Mode = "file"
	cfg.Storage.SaveNew = true
	cfg.Storage.Root = "kimgs"

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("[WARN] %s\n", err)
		log.Println("[INFO] configuration [default] used")
	} else {
		log.Printf("[INFO] configuration [%s] used\n", configFile)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Printf("[ERROR] %s\n", err)
			return nil, err
		}
	}

	// httpd env
	if env, ok := os.LookupEnv("KIMG_HTTPD_BIND"); ok {
		cfg.Httpd.Bind = env
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_HEADERS"); ok {
		arr := strings.Split(env, ",")
		for _, v := range arr {
			s := strings.Split(v, ":")
			if len(s) == 2 {
				cfg.Httpd.Headers[s[0]] = s[1]
			}
		}
	}
	if env, ok := os.LookupEnv("KIMG_HTTPD_ETAG"); ok {
		cfg.Httpd.Etag, _ = strconv.ParseBool(env)
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
		cfg.Httpd.EnableWeb, _ = strconv.ParseBool(env)
	}

	// image env
	if env, ok := os.LookupEnv("KIMG_IMAGE_FORMAT"); ok {
		cfg.Image.Format = env
	}
	if env, ok := os.LookupEnv("KIMG_IMAGE_QUALITY"); ok {
		cfg.Image.Quality, _ = strconv.Atoi(env)
	}
	if env, ok := os.LookupEnv("KIMG_IMAGE_ALLOWED_TYPES"); ok {
		cfg.Image.AllowedTypes = strings.Split(env, ",")
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
	if env, ok := os.LookupEnv("KIMG_CACHE_MEMCACHE_URL"); ok {
		cfg.Cache.Memcache.URL = env
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_REDIS_URL"); ok {
		cfg.Cache.Redis.URL = env
	}
	if env, ok := os.LookupEnv("KIMG_CACHE_MEMORY_CAPACITY"); ok {
		cfg.Cache.Memory.Capacity, _ = strconv.ParseInt(env, 0, 64)
	}

	// storage env
	if env, ok := os.LookupEnv("KIMG_STORAGE_MODE"); ok {
		cfg.Storage.Mode = env
	}
	if env, ok := os.LookupEnv("KIMG_STORAGE_SAVE_NEW"); ok {
		cfg.Storage.SaveNew, _ = strconv.ParseBool(env)
	}
	if env, ok := os.LookupEnv("KIMG_STORAGE_ROOT"); ok {
		cfg.Storage.Root = env
	}

	return &cfg, nil
}
