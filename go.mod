module AynaLivePlayer

go 1.23.0

toolchain go1.24.4

replace (
	github.com/AynaLivePlayer/liveroom-sdk v0.1.0 => ./pkg/liveroom-sdk // submodule
	github.com/AynaLivePlayer/miaosic v0.2.1 => ./pkg/miaosic // submodule

	github.com/saltosystems/winrt-go => github.com/go-musicfox/winrt-go v0.1.4 // winrt with media foundation
)

require (
	fyne.io/fyne/v2 v2.6.1
	github.com/AynaLivePlayer/liveroom-sdk v0.1.0
	github.com/AynaLivePlayer/miaosic v0.2.1
	github.com/adrg/libvlc-go/v3 v3.1.6
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/aynakeya/go-mpv v0.0.8
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20250301202403-da16c1255728
	github.com/go-ole/go-ole v1.3.0
	github.com/go-resty/resty/v2 v2.16.5
	github.com/gorilla/websocket v1.5.3
	github.com/mattn/go-colorable v0.1.14
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/saltosystems/winrt-go v0.0.0-20241223121953-98e32661f6ff
	github.com/sirupsen/logrus v1.9.3
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.10.0
	github.com/tidwall/gjson v1.18.0
	github.com/virtuald/go-paniclog v0.0.0-20190812204905-43a7fa316459
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b
	golang.org/x/sys v0.33.0
	gopkg.in/ini.v1 v1.67.0
)

require (
	fyne.io/systray v1.11.0 // indirect
	github.com/AynaLivePlayer/blivedm-go v0.0.0-20250527143915-74cc4b2603bc // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/XiaoMengXinX/Music163Api-Go v0.1.30 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/aynakeya/deepcolor v1.0.3 // indirect
	github.com/aynakeya/open-bilibili-live v0.0.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dhowden/tag v0.0.0-20240417053706-3d75831295e8 // indirect
	github.com/fredbi/uri v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fyne-io/gl-js v0.1.0 // indirect
	github.com/fyne-io/glfw-js v0.2.0 // indirect
	github.com/fyne-io/image v0.1.1 // indirect
	github.com/fyne-io/oksvg v0.1.0 // indirect
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/go-text/render v0.2.0 // indirect
	github.com/go-text/typesetting v0.2.1 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hack-pad/go-indexeddb v0.3.2 // indirect
	github.com/hack-pad/safejs v0.1.0 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20241217141322-fcc2cadd6f08 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jsummers/gobmp v0.0.0-20230614200233-a9de23ed2e25 // indirect
	github.com/makiuchi-d/gozxing v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rymdport/portal v0.4.1 // indirect
	github.com/sahilm/fuzzy v0.1.1 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/yuin/goldmark v1.7.8 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/image v0.24.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace (
//	github.com/aynakeya/blivedm => D:\Repository\blivedm
//	github.com/aynakeya/go-mpv => D:\Repository\go-mpv
//)
