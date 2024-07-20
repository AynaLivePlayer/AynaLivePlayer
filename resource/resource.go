package resource

import "fyne.io/fyne/v2"

var ImageEmpty = &fyne.StaticResource{
	StaticName:    "flat-color-icons--audio-file.svg",
	StaticContent: []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 48 48\"><path fill=\"none\" d=\"M204 0h48v48h-48z\"/><path fill=\"#90caf9\" d=\"M244 45h-32V3h22l10 10z\"/><path fill=\"#e1f5fe\" d=\"M242.5 14H233V4.5z\"/><g fill=\"#1976d2\"><circle cx=\"227\" cy=\"30\" r=\"4\"/><path d=\"m234 21l-5-2v11h2v-7.1l3 1.1z\"/></g><path fill=\"#90caf9\" d=\"M40 45H8V3h22l10 10z\"/><path fill=\"#e1f5fe\" d=\"M38.5 14H29V4.5z\"/><g fill=\"#1976d2\"><circle cx=\"23\" cy=\"30\" r=\"4\"/><path d=\"m30 21l-5-2v11h2v-7.1l3 1.1z\"/></g></svg>"),
}

var ImageEmptyQrCode = &fyne.StaticResource{
	StaticName:    "flat-color-icons--qr-code.svg",
	StaticContent: []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\"><path fill=\"black\" d=\"M1 1h10v10H1zm2 2v6h6V3z\"/><path fill=\"black\" fill-rule=\"evenodd\" d=\"M5 5h2v2H5z\"/><path fill=\"black\" d=\"M13 1h10v10H13zm2 2v6h6V3z\"/><path fill=\"black\" fill-rule=\"evenodd\" d=\"M17 5h2v2h-2z\"/><path fill=\"black\" d=\"M1 13h10v10H1zm2 2v6h6v-6z\"/><path fill=\"black\" fill-rule=\"evenodd\" d=\"M5 17h2v2H5z\"/><path fill=\"black\" d=\"M23 19h-4v4h-6V13h1h-1v6h2v2h2v-6h-2v-2h-1h3v2h2v2h2v-4h2zm0 2v2h-2v-2z\"/></svg>"),
}

var ImageIcon = resImageIcon
