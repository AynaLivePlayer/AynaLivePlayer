package gutil

import (
	"bytes"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/go-resty/resty/v2"
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
		get, err := resty.New().R().Get(picture.Url)
		if err != nil {
			return nil, err
		}
		img = canvas.NewImageFromReader(bytes.NewReader(get.Body()), "cover")
		// NewImageFromURI will return an image with empty resource and file
		if img == nil {
			return nil, errors.New("fail to download image")
		}
	}
	if img.Resource == nil {
		return nil, errors.New("fail to read image")
	}
	// compress image, so it won't be too large
	img.Resource = ResizeImage(img.Resource, 128, 128)
	return img, nil
}
