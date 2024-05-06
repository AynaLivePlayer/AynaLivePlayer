package gutil

import (
	"bytes"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/nfnt/resize"
	"image"
	"image/png"
)

func ResizeImage(resource fyne.Resource, width int, height int) fyne.Resource {
	data := resource.Content()
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return resource
	}
	img = resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)
	buf := bytes.NewBuffer([]byte{})
	err = png.Encode(buf, img)
	if err != nil {
		return resource
	}
	return fyne.NewStaticResource(resource.Name(), buf.Bytes())
}

func NewImageFromPlayerPicture(picture miaosic.Picture) (*canvas.Image, error) {
	var img *canvas.Image
	if picture.Data != nil {
		img = canvas.NewImageFromReader(bytes.NewReader(picture.Data), "cover")
		// return an error when img is nil
		if img == nil {
			return nil, errors.New("fail to read image")
		}

	} else {
		uri, err := storage.ParseURI(picture.Url)
		if err != nil {
			return nil, err
		}
		if uri == nil {
			return nil, errors.New("fail to fail url")
		}
		img = canvas.NewImageFromURI(uri)
		if img == nil || (img.File == "" && img.Resource == nil) {
			// bug fix, return a new error to indicate fail to read an image
			return nil, errors.New("fail to read image")
		}
	}
	if img.Resource == nil {
		return nil, errors.New("fail to read image")
	}
	// compress image, so it won't be too large
	img.Resource = ResizeImage(img.Resource, 128, 128)
	return img, nil
}
