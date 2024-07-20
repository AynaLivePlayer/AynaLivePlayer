# AynaLivePlayer

Bilibili Audio Bot. Written by Golang.

Provider By Aynakeya

QQ group: 621035845

## build

```
go build -o AynaLivePlayer.exe -ldflags -H=windowsgui app/gui/main.go
```

## packaging
```
fyne package --src path_to_gui --exe AynaLivePlayer.exe --appVersion 0.8.4 --icon path_to_icon
```

## Windows build guide

1. install golang [link](https://go.dev/doc/install)
2. install chocolatey [link](https://chocolatey.org/install)
3. install required packages
```
choco install git
choco install mingw
```
4. install fyne
```
go install fyne.io/fyne/v2/cmd/fyne@latest
```
5. clone this repo
```bash
git clone --recurse-submodules git@github.com:AynaLivePlayer/AynaLivePlayer.git
```
if you are using https links
```
git clone https://github.com/AynaLivePlayer/AynaLivePlayer.git
git submodule set-url pkg/miaosic https://github.com/AynaLivePlayer/miaosic.git
git submodule set-url pkg/liveroom-sdk https://github.com/AynaLivePlayer/liveroom-sdk.git
git submodule update
```
6. now you can build (please check makefile for more details)
```powershell
$env:CGO_LDFLAGS="-LC:\Users\Admin\Desktop\AynaLivePlayer\libmpv\lib";$env:CGO_CFLAGS="-IC:\Users\Admin\Desktop\AynaLivePlayer\libmpv\include"
# ... more setup, see makefile
go build -o AynaLivePlayer.exe -ldflags -H=windowsgui app/main.go
```