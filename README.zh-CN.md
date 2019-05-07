# kimg
> 一款基于go开发的图片服务器

[![Build Status](https://img.shields.io/travis/zhoukk/kimg.svg?style=flat)](https://travis-ci.org/zhoukk/kimg)

[English](./README.md) | 简体中文


## 快速开始

> 下载并启动

- 从linux二进制运行文件启动
```console
$ wget -O- https://github.com/zhoukk/kimg/releases/download/release-v0.3.1/kimg_v0.3.1_linux.tar.gz | tar xf -
$ cd kimg_v0.3.1_linux
$ ./kimg
```

- 从docker镜像启动
```console
$ docker pull zhoukk/kimg:v0.3.1
$ docker run --rm -p 80:80 zhoukk/kimg
```

> 打开浏览器体验
```console
$ open http://localhost
```

## docker镜像中包含的文件结构

<a href="https://asciinema.org/a/243736?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243736.svg" width=480 /></a>

## 用法

> 上传图片到kimg服务

- 使用raw-data方式post上传

<a href="https://asciinema.org/a/243841?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243841.svg" width=480 /></a>

- 使用multipart-form方式post上传

<a href="https://asciinema.org/a/243754?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243754.svg" width=480 /></a>

> 获取一个指定样式的图片

```console
$ open http://localhost/image/323551c4a7e2071a28a41331b98ca821?s=1&sm=fit&sw=300&sh=300&c=1&cw=200&ch=200
```    

> 获取图片的信息

<a href="https://asciinema.org/a/243758?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243758.svg" width=480 /></a>

## License

[MIT](https://github.com/zhoukk/kimg/blob/master/LICENSE)