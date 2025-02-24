package cloudinary

import (
	"context"
	"io"
	"log"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Client interface {
	Upload(file io.Reader, path string) error
	Delete(path string) error
	GenerateUrl(path string) string
}

type Cloudinary struct {
	cld *cloudinary.Cloudinary
}

func New(url string) Client {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		log.Println("fail to create cloudinary client")
	}
	return &Cloudinary{
		cld: cld,
	}
}

func (c *Cloudinary) Upload(file io.Reader, path string) (err error) {
	_, err = c.cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicID: path,
	})
	if err != nil {
		return
	}

	return
}

func (c Cloudinary) Delete(path string) (err error) {
	_, err = c.cld.Admin.DeleteAssets(
		context.Background(),
		admin.DeleteAssetsParams{PublicIDs: []string{path}},
	)
	if err != nil {
		return
	}

	return
}

func (c Cloudinary) GenerateUrl(path string) string {
	resp, _ := c.cld.Admin.Asset(context.Background(), admin.AssetParams{PublicID: path})

	return resp.SecureURL
}
