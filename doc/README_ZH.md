<p align="center">
   <img width="170" src="https://raw.githubusercontent.com/go-sonic/resources/master/logo/logo.svg" />
</p>

<p align="center"><b>Sonic </b> [ˈsɒnɪk] ,Sonic 是一个用Golang开发的博客平台，高效快速.</p>

<p align="center">
<a href="https://github.com/go-sonic/sonic/releases"><img alt="GitHub release" src="https://img.shields.io/github/release/go-sonic/sonic.svg?style=flat-square&include_prereleases" /></a>
<a href="https://github.com/go-sonic/sonic/releases"><img alt="GitHub All Releases" src="https://img.shields.io/github/downloads/go-sonic/sonic/total.svg?style=flat-square" /></a>
<a href="https://hub.docker.com/r/gosonic/sonic"><img alt="Docker pulls" src="https://img.shields.io/docker/pulls/gosonic/sonic?style=flat-square" /></a>
<a href="https://github.com/go-sonic/sonic/commits"><img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/go-sonic/sonic.svg?style=flat-square" /></a>

<br />
<a href="https://t.me/go_sonic">Telegram 频道</a>
</p>


## 📖 介绍

Sonic 意为声速的、声音的，正如它的名字一样, sonic 致力于成为最快速的开源博客平台。

感谢 [Halo](https://github.com/halo-dev/) 项目组，本项目的灵感来自Halo，前端项目Fork自[Console](https://github.com/halo-dev)

## 🚀 Features:
- 支持多种类型的数据库：SQLite、MySQL(TODO: PostgreSQL)
- 体积小: 安装包仅仅只有10Mb
- 高性能: 文章详情页可以达到2500 QPS(压测环境是: Intel Xeon Platinum 8260 4C 8G ,SQLite3)
- 支持更换主题
- 支持 Linux、Windows、Mac OS等主流操作系统，支持x86、x64、Arm、Arm64、MIPS等指令集架构
- 支持对象存储(MINIO、Google Cloud、AWS、AliYun)

## 🎊 Preview

![默认主题](https://github.com/go-sonic/default-theme-anatole/raw/master/screenshot.png)

![控制台](https://github.com/go-sonic/resources/raw/master/console-screenshot.png)

## 🧰 安装

**下载对应平台的安装包**
> 根据你的操作系统和指令集下载对应的安装包,注意要下载最新的版本
```bash
wget https://github.com/go-sonic/sonic/releases/download/v1.0.3/sonic-linux-amd64.zip -O sonic.zip
```
**解压**
```bash
unzip -d sonic sonic.zip
```
**运行**
> 可以通过 -config选项来指定配置文件的位置
```bash
cd sonic
./sonic -config conf/config.yaml
```

**然后你就可以通过浏览器访问sonic了，默认的端口是8080**

后台管理路径是 http://ip:port/admin

## 🔨️  构建
**1. 拉取项目**
```bash
git clone --recursive --depth 1 https://github.com/go-sonic/sonic
```
**2. 运行**
```bash
cd sonic
go run main.go
```
> 如果你在windows上构建该项目，请确保你已经正确安装了gcc编译器，比如，你可以在[这里](https://jmeubank.github.io/tdm-gcc/)找到一个TDM版本的gcc编译器。

🚀 完成! 你的项目现在已经运行起来了。

## Docker
See: https://hub.docker.com/r/gosonic/sonic

## 主题生态

| Theme   | 
|---------|
| [Anatole](https://github.com/go-sonic/default-theme-anatole) |
| [Journal](https://github.com/hooxuu/sonic-theme-Journal) |
| [Clark](https://github.com/ClarkQAQ/sonic_theme_clark)   |
| [Earth](https://github.com/Meepoljdx/sonic-theme-earth) |
| [PaperMod](https://github.com/jakezhu9/sonic-theme-papermod) |
| [Tink](https://github.com/raisons/sonic-theme-tink) |

## TODO
- [ ] i18n
- [ ] PostgreSQL
- [ ] 更好的错误处理
- [ ] 插件系统(基于 Wasm)
- [ ] 使用新的web框架([Hertz](https://github.com/cloudwego/hertz))


## 如何贡献

非常欢迎你的加入！[提一个 Issue](https://github.com/go-sonic/sonic/issues) 或者提交一个 Pull Request。


Sonic 遵循 [Contributor Covenant](http://contributor-covenant.org/version/1/3/0/) 行为规范。

### 贡献者

感谢以下参与项目的人：
<a href="https://github.com/go-sonic/sonic/graphs/contributors"><img src="https://opencollective.com/go-sonic/contributors.svg?width=890&button=false" /></a>

特别感谢 Evan(evanzhao@88.com) ，他帮助设计了Logo。


## 📄 License

Source code in `sonic` is available under the [MIT License](/LICENSE.md).

