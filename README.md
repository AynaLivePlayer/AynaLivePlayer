# AynaLivePlayer

Bilibili Audio Bot. Written by Golang.

Provider By Aynakeya

QQ group: 621035845

## Disclaimer

All APIs used in this project are  **publicly available** on the internet and not obtained through illegal means such as
reverse engineering.

The use of this project may involve access to copyrighted content. This project does **not** own or claim any rights to
such content. **To avoid potential infringement**, all users are **required to delete any copyrighted data obtained
through this project within 24 hours.**

Any direct, indirect, special, incidental, or consequential damages (including but not limited to loss of goodwill, work
stoppage, computer failure or malfunction, or any and all other commercial damages or losses) that arise from the use or
inability to use this project are **solely the responsibility of the user**.

This project is completely free and open-source, published on GitHub for global users for **technical learning and
research purposes only**. This project does **not** guarantee compliance with local laws or regulations in all
jurisdictions.

**Using this project in violation of local laws is strictly prohibited.** Any legal consequences arising from
intentional or unintentional violations are the user's responsibility. The project maintainers accept **no liability**
for such outcomes.


## build


> outdated, please refer to workflow file

```
go build -o AynaLivePlayer.exe -ldflags -H=windowsgui app/gui/main.go
```

## packaging

> outdated, please refer to workflow file

```
fyne package --src path_to_gui --exe AynaLivePlayer.exe --appVersion 0.8.4 --icon path_to_icon
```

## Windows build guide

> outdated, please refer to workflow file

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