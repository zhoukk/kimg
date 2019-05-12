# kimg
> 一款基于go开发的图片服务器

[![Build Status](https://img.shields.io/travis/zhoukk/kimg.svg?style=flat)](https://travis-ci.org/zhoukk/kimg)

[English](./README.md) | 简体中文


## 快速开始

> 下载并启动

- 从linux二进制运行文件启动
```console
$ wget -O- https://github.com/zhoukk/kimg/releases/download/release-v0.4.1/kimg_v0.4.1_linux.tar.gz | tar xf -
$ cd kimg_v0.4.1_linux
$ ./kimg
```

- 从docker镜像启动
```console
$ docker pull zhoukk/kimg:v0.4.1
$ docker run --rm -p 80:80 zhoukk/kimg:v0.4.1
```

> 打开浏览器体验
```console
$ open http://localhost
```

## docker镜像中包含的文件结构

<img src="http://kimg.zhoukk.com/image/2fb0757f132497b06f0cdceda9a8d8a1?origin=1" width=480 />

## 用法

> 上传图片到kimg服务

- 使用data-binary方式post上传

<img src="http://kimg.zhoukk.com/image/c99cbbebd327c6f3b3cdb190d1a8a95a?origin=1" width=480 />

- 使用multipart-form方式post上传

<img src="http://kimg.zhoukk.com/image/c55e9ad1cc4618a5bb0e47097a2b9eb3?origin=1" width=480 />

> 获取一个指定样式的图片

```console
$ open http://localhost/image/323551c4a7e2071a28a41331b98ca821?s=1&sm=fit&sw=300&sh=300&c=1&cw=200&ch=200
```    

> 获取图片的信息

<img src="http://kimg.zhoukk.com/image/5769d4865b750885710d987d3131f16d?origin=1" width=480 />

## License

[MIT](https://github.com/zhoukk/kimg/blob/master/LICENSE)