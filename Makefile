NAME = AynaLivePlayer

ifeq ($(OS),Windows_NT)
    RM = del /Q /F
    RRM = rmdir /Q /S
else
    RM = rm -f
    RRM = rm -rf
endif

ifeq ($(OS), Windows_NT)
	EXECUTABLE=$(NAME).exe
	SCRIPTPATH = .\assets\scripts\windows
else
	EXECUTABLE=$(NAME)
	SCRIPTPATH = ./assets/scripts/linux
endif

gui: bundle
	go build -o $(EXECUTABLE) -ldflags -H=windowsgui main.go

run: bundle
	go run main.go

clear:
	$(RM) config.ini log.txt playlists.txt liverooms.json

bundle:
	fyne bundle --name resImageEmpty --package resource ./assets/empty.png >  ./resource/bundle.go
	fyne bundle --append --name resImageIcon --package resource ./assets/icon.jpg >> ./resource/bundle.go
#	fyne bundle --append --name resFontMSYaHei --package resource ./assets/msyh.ttc >> ./resource/bundle.go
#	fyne bundle --append --name resFontMSYaHeiBold --package resource ./assets/msyhbd.ttc >> ./resource/bundle.go
	fyne bundle --append --name resFontMSYaHei --package resource ./assets/msyh0.ttf >> ./resource/bundle.go
	fyne bundle --append --name resFontMSYaHeiBold --package resource ./assets/msyhbd0.ttf >> ./resource/bundle.go

release: gui
	-mkdir release
ifeq ($(OS), Windows_NT)
	COPY .\$(EXECUTABLE) .\release\$(EXECUTABLE)
	COPY .\webtemplates.json .\release\webtemplates.json
	mkdir .\release\assets
	COPY .\assets\mpv-2.dll .\release\mpv-2.dll
	COPY .\assets\translation.json .\release\assets\translation.json
	COPY LICENSE.md .\release\LICENSE.md
	XCOPY  .\assets\scripts\windows\* .\release\ /k /i /y /q
	XCOPY  .\assets\webinfo .\release\assets\webinfo /s /e /i /y /q
	XCOPY  .\music .\release\music /s /e /i /y /q
	XCOPY  .\template .\release\template /s /e /i /y /q
else
	cp ./$(EXECUTABLE) ./release/$(EXECUTABLE)
	cp ./webtemplates.json ./release/webtemplates.json
	cp ./assets/translation.json ./release/assets/translation.json
	mkdir ./release/assets
	cp LICENSE.md ./release/LICENSE.md
	cp ./assets/scripts/linux/* ./release/
	cp -r ./assets/webinfo ./release/assest/webinfo
	cp -r ./music ./release/music
	cp -r ./template ./release/template
endif

clean:
	$(RM) $(EXECUTABLE) config.ini log.txt playlists.txt liverooms.json
	$(RRM) release

.PHONY: ${EXECUTABLE}