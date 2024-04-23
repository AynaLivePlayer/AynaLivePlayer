EXECUTABLE=AynaLivePlayer
WINDOWS=$(EXECUTABLE).exe
LINUX=$(EXECUTABLE)_linux
DARWIN=$(EXECUTABLE)_darwin

ifeq ($(OS),Windows_NT)
    RM = del /Q /F
    RRM = rmdir /Q /S
    MKDIR = mkdir
    COPY = XCOPY /Y
    MOVE = move
else
    RM = rm -f
    RRM = rm -rf
    MKDIR = mkdir
    COPY = cp -r
    MOVE = mv
endif

bundle:
	fyne bundle --name resImageIcon --package resource ./assets/icon.png > ./resource/bundle.go
#	fyne bundle --append --name resFontMSYaHei --package resource ./assets/msyh.ttc >> ./resource/bundle.go
#	fyne bundle --append --name resFontMSYaHeiBold --package resource ./assets/msyhbd.ttc >> ./resource/bundle.go
	fyne bundle --append --name resFontMSYaHei --package resource ./assets/msyh0.ttf >> ./resource/bundle.go
	fyne bundle --append --name resFontMSYaHeiBold --package resource ./assets/msyhbd0.ttf >> ./resource/bundle.go

prebuild: bundle
	$(RRM) ./release
	$(MKDIR) ./release
	$(MKDIR) ./release/assets
	$(COPY) LICENSE.md ./release/LICENSE.md
	$(COPY) ./assets/translation.json ./release/assets/translation.json
	$(COPY) ./assets/config ./release/config
	$(COPY) ./music ./release/music
	go mod tidy


$(LINUX): prebuild
	env GOOS=linux GOARCH=amd64 go build -o ./release/$(LINUX) app/main.go
	$(MOVE) ./release/$(LINUX) ./release/$(EXECUTABLE)

$(WINDOWS): prebuild
	env GOOS=windows GOARCH=amd64 go build -o ./release/$(WINDOWS) -ldflags -H=windowsgui app/main.go

$(DARWIN): prebuild
	env GOOS=darwin GOARCH=amd64 go build -o ./release/$(DARWIN) app/main.go
	$(MOVE) ./release/$(LINUX) ./release/$(EXECUTABLE)


windows: $(WINDOWS) ## Build for Windows
	$(COPY) ./assets/windows/mpv-2.dll ./release/mpv-2.dll

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

clean:
	$(RRM) ./release