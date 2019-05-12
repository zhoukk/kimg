# kimg
> Yet another image server build in go

[![Build Status](https://img.shields.io/travis/zhoukk/kimg.svg?style=flat)](https://travis-ci.org/zhoukk/kimg)

English | [简体中文](./README.zh-CN.md)


## Quick Start

> run kimg

- run with linux binary
```console
$ wget -O- https://github.com/zhoukk/kimg/releases/download/release-v0.4.1/kimg_v0.4.1_linux.tar.gz | tar xf -
$ cd kimg_v0.4.1_linux
$ ./kimg
```

- run with docker
```console
$ docker pull zhoukk/kimg:v0.4.1
$ docker run --rm -p 80:80 zhoukk/kimg:v0.4.1
```

> open a browser and have fun
```console
$ open http://localhost
```

## What's in docker image

<img src="http://kimg.zhoukk.com/image/2fb0757f132497b06f0cdceda9a8d8a1?origin=1" width=480 />

## Usage

> Upload a image to kimg

- Upload use data-binary

<img src="http://kimg.zhoukk.com/image/c99cbbebd327c6f3b3cdb190d1a8a95a?origin=1" width=480 />

- Upload use multipart-form

<img src="http://kimg.zhoukk.com/image/c55e9ad1cc4618a5bb0e47097a2b9eb3?origin=1" width=480 />

> Fetch a image with style from kimg

```console
$ open http://localhost/image/323551c4a7e2071a28a41331b98ca821?s=1&sm=fit&sw=300&sh=300&c=1&cw=200&ch=200
```    

> Get a image information

<img src="http://kimg.zhoukk.com/image/5769d4865b750885710d987d3131f16d?origin=1" width=480 />

## License

[MIT](https://github.com/zhoukk/kimg/blob/master/LICENSE)