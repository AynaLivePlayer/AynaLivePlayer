module AynaLivePlayer

go 1.19

replace (
	github.com/AynaLivePlayer/liveroom-sdk v0.1.0 => ./pkg/liveroom-sdk // submodule
	github.com/AynaLivePlayer/miaosic v0.1.5 => ./pkg/miaosic // submodule
)

require (
	fyne.io/fyne/v2 v2.4.5
	fyne.io/x/fyne v0.0.0-20240326131024-3ba9170cc3be
	github.com/AynaLivePlayer/liveroom-sdk v0.1.0
	github.com/AynaLivePlayer/miaosic v0.1.5
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/aynakeya/go-mpv v0.0.6
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240306074159-ea2d69986ecb
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gorilla/websocket v1.5.0
	github.com/mattn/go-colorable v0.1.12
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/sirupsen/logrus v1.9.3
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/gjson v1.16.0
	github.com/virtuald/go-paniclog v0.0.0-20190812204905-43a7fa316459
	go.uber.org/zap v1.26.0
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	golang.org/x/sys v0.13.0
	gopkg.in/ini.v1 v1.67.0
)

require (
	fyne.io/systray v1.10.1-0.20231115130155-104f5ef7839e // indirect
	github.com/AynaLivePlayer/blivedm-go v0.0.0-20240408074929-6565ab41764b // indirect
	github.com/PuerkitoBio/goquery v1.7.1 // indirect
	github.com/XiaoMengXinX/Music163Api-Go v0.1.30 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/aynakeya/deepcolor v1.0.2 // indirect
	github.com/aynakeya/open-bilibili-live v0.0.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dhowden/tag v0.0.0-20230630033851-978a0926ee25 // indirect
	github.com/fredbi/uri v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fyne-io/gl-js v0.0.0-20220119005834-d2da28d9ccfe // indirect
	github.com/fyne-io/glfw-js v0.0.0-20220120001248-ee7290d23504 // indirect
	github.com/fyne-io/image v0.0.0-20220602074514-4956b0afb3d2 // indirect
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-text/render v0.1.0 // indirect
	github.com/go-text/typesetting v0.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jsummers/gobmp v0.0.0-20151104160322-e2ba15ffa76e // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sahilm/fuzzy v0.1.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/tevino/abool v1.2.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/yuin/goldmark v1.5.5 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/image v0.11.0 // indirect
	golang.org/x/mobile v0.0.0-20230531173138-3c911d8e3eda // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/js/dom v0.0.0-20210725211120-f030747120f2 // indirect
)

//replace (
//	github.com/aynakeya/blivedm => D:\Repository\blivedm
//	github.com/aynakeya/go-mpv => D:\Repository\go-mpv
//)
