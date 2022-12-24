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
else
	EXECUTABLE=$(NAME)
endif

${EXECUTABLE}:
	go build -o $(EXECUTABLE) -ldflags -H=windowsgui ./app/gui/main.go

run:
	go run ./app/gui/main.go

clear:
	$(RM) config.ini log.txt playlists.txt liverooms.json

bundle:
	fyne bundle --name resImageEmpty --package resource ./assets/empty.png >  ./resource/bundle.go
	fyne bundle --append --name resImageIcon --package resource ./assets/icon.jpg >> ./resource/bundle.go
	fyne bundle --append --name resFontMSYaHei --package resource ./assets/msyh.ttc >> ./resource/bundle.go
	fyne bundle --append --name resFontMSYaHeiBold --package resource ./assets/msyhbd.ttc >> ./resource/bundle.go

release: ${EXECUTABLE} bundle
	-mkdir release
ifeq ($(OS), Windows_NT)
	COPY .\$(EXECUTABLE) .\release\$(EXECUTABLE)
	COPY .\webtemplates.json .\release\webtemplates.json
	COPY LICENSE.md .\release\LICENSE.md
	XCOPY  .\assets .\release\assets /s /e /i /y /q
	XCOPY  .\music .\release\music /s /e /i /y /q
	XCOPY  .\template .\release\template /s /e /i /y /q
else
	cp ./$(EXECUTABLE) ./release/$(EXECUTABLE)
	cp ./webtemplates.json ./release/webtemplates.json
	cp LICENSE.md ./release/LICENSE.md
	cp -r ./assets ./release/assest
	cp -r ./music ./release/music
	cp -r ./template ./release/template
endif

clean:
	$(RM) $(EXECUTABLE) config.ini log.txt playlists.txt liverooms.json
	$(RRM) release

.PHONY: ${EXECUTABLE}