# kimg
> Yet another image server build in go

[![Build Status](https://img.shields.io/travis/zhoukk/kimg.svg?style=flat)](https://travis-ci.org/zhoukk/kimg)

English | [简体中文](./README.zh-CN.md)


## Quick Start

> run kimg

- run with linux binary
```console
$ wget -O- https://github.com/zhoukk/kimg/releases/download/release-v0.3.1/kimg_v0.3.1_linux.tar.gz | tar xf -
$ cd kimg_v0.3.1_linux
$ ./kimg
```

- run with docker
```console
$ docker pull zhoukk/kimg:v0.3.1
$ docker run --rm -p 80:80 zhoukk/kimg
```

> open a browser and have fun
```console
$ open http://localhost
```

## What's in docker image

<a href="https://asciinema.org/a/243736?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243736.svg" width=480 /></a>

## Usage

> Upload a image to kimg

- Upload use raw-data

<a href="https://asciinema.org/a/243841?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243841.svg" width=480 /></a>

- Upload use multipart-form

<a href="https://asciinema.org/a/243754?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243754.svg" width=480 /></a>

> Fetch a image with style from kimg

```console
$ open http://localhost/image/323551c4a7e2071a28a41331b98ca821?s=1&sm=fit&sw=300&sh=300&c=1&cw=200&ch=200
```    

> Get a image information

<a href="https://asciinema.org/a/243758?autoplay=1" target="_blank"><img src="https://asciinema.org/a/243758.svg" width=480 /></a>


## License

[MIT](https://github.com/zhoukk/kimg/blob/master/LICENSE)