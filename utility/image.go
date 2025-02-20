package utility

import (
	"bytes"
	"mime/multipart"

	"github.com/PhantomX7/go-core/utility/errors"
	"github.com/disintegration/imaging"
)

func FormatImage(file multipart.File) (b bytes.Buffer, err error) {
	// resize the image while keeping the ratio
	dImage, err := imaging.Decode(file)
	if err != nil {
		err = errors.ErrUnprocessableEntity
		return
	}

	// resize output image
	dImage = imaging.Resize(dImage, dImage.Bounds().Max.X, dImage.Bounds().Max.Y, imaging.Lanczos)

	_ = imaging.Encode(&b, dImage, imaging.JPEG, imaging.JPEGQuality(50))

	return
}

func cropImage(file multipart.File) (b bytes.Buffer, err error) {
	// resize the image while keeping the ratio
	dImage, err := imaging.Decode(file)
	if err != nil {
		err = errors.ErrUnprocessableEntity
		return
	}

	// crop the image
	dImage = imaging.CropCenter(dImage, 100, 100)

	_ = imaging.Encode(&b, dImage, imaging.JPEG, imaging.JPEGQuality(50))

	return
}
