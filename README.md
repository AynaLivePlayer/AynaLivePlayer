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