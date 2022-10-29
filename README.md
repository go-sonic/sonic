<p align="center">
   <img src="https://raw.githubusercontent.com/go-sonic/resources/master/logo/logo.png" />
</p>

<p align="center"><b>Sonic </b> [ˈsɒnɪk] ,Sonic is a Go Blogging Platform. Simple and Powerful.</p>

<p align="center">
<a href="https://github.com/go-sonic/sonic/releases"><img alt="GitHub release" src="https://img.shields.io/github/release/go-sonic/sonic.svg?style=flat-square&include_prereleases" /></a>
<a href="https://github.com/go-sonic/sonic/releases"><img alt="GitHub All Releases" src="https://img.shields.io/github/downloads/go-sonic/sonic/total.svg?style=flat-square" /></a>
<a href="https://hub.docker.com/r/go-sonic/sonic"><img alt="Docker pulls" src="https://img.shields.io/docker/pulls/go-sonic/sonic?style=flat-square" /></a>
<a href="https://github.com/go-sonic/sonic/commits"><img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/go-sonic/sonic.svg?style=flat-square" /></a>
<a href="https://github.com/go-sonic/sonic/actions"><img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/go-sonic/sonic/Sonic%20CI?style=flat-square" /></a>
<br />
<a href="https://go-sonic.org">Website</a>
<a href="https://t.me/go_sonic">Telegram Channel</a>
</p>


English | [中文](doc/README_ZH.md)

## 📖 Introduction

Sonic means as fast as sound speed. Like its name, sonic is a high-performance blog system developed using golang

Thanks [Halo](https://github.com/halo-dev) project team,this project is inspired by Halo. Front end project fork from Halo

## 🚀 Features:
- Support multiple types of databases: SQLite、MySQL(TODO: PostgreSQL)
- Small: The installation file is only 10mb size
- High-performance: Post details page can withstand 900qps(Enviroment:   Intel Xeon Platinum 8260 4C 8G ,SQLite3)
- Support changing theme
- Support Linux、Windows、Mac OS. And Support x86、x64、Arm、Arm64、MIPS
- Object storage(MINIO、Google Cloud、AWS、AliYun)


## 🧰 How to install

### Download the latest installation package
> Please pay attention to the operating system and instruction set
```bash
wget https://github.com/go-sonic/sonic/releases/download/v1.0.0/sonic-linux-64.zip -O sonic.zip
```
### Decompression
```bash
unzip sonic.zip
```
### Launch
```bash
cd sonic
./sonic -config conf/config.yaml
```

### Initialization
**The default port is 8080**

Open http://ip:port/admin#install

Next, you can access sonic through the browser.

The URL of the admin console is http://ip:port/admin


## TODO
- [ ] i18n
- [ ] PostgreSQL
- [ ] Better error handling
- [ ] Plugin(base on Wasm)
- [ ] Use new web framework([Hertz](https://github.com/cloudwego/hertz))

## 📄 License

Source code in `sonic` is available under the [MIT License](/LICENSE.md).

