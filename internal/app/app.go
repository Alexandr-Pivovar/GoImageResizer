package app

import (
	"GoImageZip/internal/domain"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"reflect"
)

var (
	ErrCould  = errors.New("cloud error")
	ErrStore  = errors.New("store error")
	ErrDecode = errors.New("decode error")
)

//go:generate mockery -name ImageRepository -case underscore
// ImageRepository interface for working with store of image info
type ImageRepository interface {
	Create(domain.ImageInfo) error
	Update(domain.ImageInfo) error
	GetById(string) (domain.ImageInfo, error)
	GetHistory() ([]domain.ImageInfo, error)
}

//go:generate mockery -name Clouder -case underscore
// ImageRepository interface for working with store of images
type Clouder interface {
	Save(string, domain.Image) (string, error)
	GetImage(string) ([]byte, error)
}

//go:generate mockery -name ResizeServicer -case underscore
// ImageRepository interface present interface for income requests
type ResizeServicer interface {
	Resize(domain.Image) (domain.ImageInfo, error)
	Update(domain.ImageInfo) (domain.ImageInfo, error)
	GetById(string) (domain.ImageInfo, error)
	GetHistory() ([]domain.ImageInfo, error)
}

//go:generate mockery -name Resizer -case underscore
// ImageRepository interface for resizing images
type Resizer interface {
	Do(domain.Image) (domain.Image, error)
}

type ImageService struct {
	repo    ImageRepository
	resizer Resizer
	cloud   Clouder
	uuid    func() string
	decode  func([]byte, []byte) (int, error)
}

// NewImageService create new value of ImageService
func NewImageService(repo ImageRepository, resizer Resizer, cloud Clouder) *ImageService {
	if repo == nil || reflect.ValueOf(repo).IsNil() {
		panic("repo param is nil")
	}
	if resizer == nil || reflect.ValueOf(resizer).IsNil() {
		panic("resizer param is nil")
	}
	if cloud == nil || reflect.ValueOf(cloud).IsNil() {
		panic("cloud param is nil")
	}

	return &ImageService{
		repo:    repo,
		resizer: resizer,
		cloud:   cloud,
		uuid:    func() string { return uuid.New().String() },
		decode: func(dst []byte, src []byte) (int, error) {
			return base64.StdEncoding.Decode(dst, src)
		},
	}
}

// Resize resizes image with height and width parameters
func (is ImageService) Resize(image domain.Image) (ii domain.ImageInfo, err error) {
	id := is.uuid()

	_, err = is.decode(image.Data, image.Data)
	if err != nil {
		return ii, err
	}

	resizedImg, err := is.resizer.Do(image)
	if err != nil {
		return ii, err
	}

	originUrl, err := is.cloud.Save(id, image)
	if err != nil {
		return ii, err
	}

	resizedUrl, err := is.cloud.Save(id, resizedImg)
	if err != nil {
		return ii, err
	}

	ii.Id = id
	ii.Width = image.Width
	ii.Height = image.Height
	ii.OriginUrl = originUrl
	ii.ResizedUrl = resizedUrl

	return ii, is.repo.Create(ii)
}

// GetHistory gets all items
func (is ImageService) GetHistory() ([]domain.ImageInfo, error) {
	return is.repo.GetHistory()
}

// GetById gets item by id
func (is ImageService) GetById(id string) (domain.ImageInfo, error) {
	return is.repo.GetById(id)
}

// Update updates item by id
func (is ImageService) Update(imageInfo domain.ImageInfo) (domain.ImageInfo, error) {
	im, err := is.repo.GetById(imageInfo.Id)
	if err != nil {
		return domain.ImageInfo{}, err
	}

	data, err := is.cloud.GetImage(im.OriginUrl)
	if err != nil {
		return domain.ImageInfo{}, err
	}

	resizedImg, err := is.resizer.Do(domain.Image{
		Data: data,
		Param: domain.Param{
			Height: imageInfo.Height,
			Width:  imageInfo.Width,
		},
	})
	if err != nil {
		return domain.ImageInfo{}, err
	}

	newUrl, err := is.cloud.Save(imageInfo.Id, resizedImg)
	if err != nil {
		return domain.ImageInfo{}, err
	}

	im.ResizedUrl = newUrl
	im.Height = imageInfo.Height
	im.Width = imageInfo.Width

	return im, is.repo.Update(im)
}
