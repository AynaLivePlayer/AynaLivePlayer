NAME = AynaLivePlayer

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
ifeq ($(OS), Windows_NT)
	-DEL $(EXECUTABLE) config.ini log.txt playlists.json liverooms.json /s /q
else
	rm config.ini log.txt playlists.txt liverooms.json
endif


release: ${EXECUTABLE}
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
ifeq ($(OS), Windows_NT)
	-DEL $(EXECUTABLE) /s /q
	-rmdir .\release /s /q
else
	rm -r $(EXECUTABLE) ./release
endif

.PHONY: ${EXECUTABLE}