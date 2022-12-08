<p align="center">
   <img width="200" src="https://raw.githubusercontent.com/go-sonic/resources/master/logo/logo.svg" />
</p>

<p align="center"><b>Sonic </b> [ˈsɒnɪk] ,Sonic is a Go Blogging Platform. Simple and Powerful.</p>

<p align="center">
<a href="https://github.com/go-sonic/sonic/releases"><img alt="GitHub release" src="https://img.shields.io/github/release/go-sonic/sonic.svg?style=flat-square&include_prereleases" /></a>
<a href="https://github.com/go-sonic/sonic/releases"><img alt="GitHub All Releases" src="https://img.shields.io/github/downloads/go-sonic/sonic/total.svg?style=flat-square" /></a>
<a href="https://hub.docker.com/r/gosonic/sonic"><img alt="Docker pulls" src="https://img.shields.io/docker/pulls/gosonic/sonic?style=flat-square" /></a>
<a href="https://github.com/go-sonic/sonic/commits"><img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/go-sonic/sonic.svg?style=flat-square" /></a>
<br />
<a href="https://t.me/go_sonic">Telegram Channel</a>
</p>


English | [中文](doc/README_ZH.md)

## 📖 Introduction

Sonic means as fast as sound speed. Like its name, sonic is a high-performance blog system developed using golang

Thanks [Halo](https://github.com/halo-dev) project team,this project is inspired by Halo. Front end project fork from Halo

## 🚀 Features:
- Support multiple types of databases: SQLite、MySQL(TODO: PostgreSQL)
- Small: The installation file is only 10mb size
- High-performance: Post details page can withstand 2500 QPS(Enviroment:   Intel Xeon Platinum 8260 4C 8G ,SQLite3)
- Support changing theme
- Support Linux、Windows、Mac OS. And Support x86、x64、Arm、Arm64、MIPS
- Object storage(MINIO、Google Cloud、AWS、AliYun)


## 🎊 Preview

![Default Theme](https://github.com/go-sonic/default-theme-anatole/raw/master/screenshot.png)

![Console](https://github.com/go-sonic/resources/raw/master/console-screenshot.png)

## 🧰 How to install

**Download the latest installation package**
> Please pay attention to the operating os and instruction set  and the version
```bash
wget https://github.com/go-sonic/sonic/releases/download/v1.0.3/sonic-linux-amd64.zip -O sonic.zip
```
**Decompression**
```bash
unzip -d sonic sonic.zip
```
**Launch**
```bash
cd sonic
./sonic -config conf/config.yaml
```

**Initialization**
**The default port is 8080**

Open http://ip:port/admin#install

Next, you can access sonic through the browser.

The URL of the admin console is http://ip:port/admin

## Docker
See: https://hub.docker.com/r/gosonic/sonic

## Theme ecology

| Theme   | URL                                               |
|---------|---------------------------------------------------|
| Anatole | https://github.com/go-sonic/default-theme-anatole |
| Journal | https://github.com/hooxuu/sonic-theme-Journal     |

## TODO
- [ ] i18n
- [ ] PostgreSQL
- [ ] Better error handling
- [ ] Plugin(base on Wasm)
- [ ] Use new web framework([Hertz](https://github.com/cloudwego/hertz))

## 📄 License

Source code in `sonic` is available under the [MIT License](/LICENSE.md).

