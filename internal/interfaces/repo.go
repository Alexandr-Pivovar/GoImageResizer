package interfaces

import (
	"GoImageZip/internal/domain"
)

type DbHandler interface {
	Create(Image) error
	GetAllItems() ([]Image, error)
	GetById(string) (Image, error)
}

type Image struct {
	Id         string `json:"id"`
	OriginUrl  string `json:"origin_url"`
	ResizedUrl string `json:"resized_url"`
	Width      uint   `json:"width"`
	Height     uint   `json:"height"`
}

type RedisRepo struct {
	dbHandler DbHandler
}

func NewRedisRepo(db DbHandler) *RedisRepo {
	return &RedisRepo{dbHandler: db}
}

// GetHistory requests history and prepare data for app layer
func (us RedisRepo) GetHistory() ([]domain.ImageInfo, error) {
	items, err := us.dbHandler.GetAllItems()
	if err != nil {
		return nil, err
	}

	imagesInfo := make([]domain.ImageInfo, 0, len(items))

	for _, item := range items {
		imagesInfo = append(imagesInfo, domain.ImageInfo{
			Id:         item.Id,
			OriginUrl:  item.OriginUrl,
			ResizedUrl: item.ResizedUrl,
			Param: domain.Param{
				Width:  item.Width,
				Height: item.Height,
			},
		})
	}
	return imagesInfo, nil
}

// Create prepares a data for dbHandler
func (us RedisRepo) Create(image domain.ImageInfo) error {
	return us.dbHandler.Create(Image{
		Id:         image.Id,
		Height:     image.Height,
		Width:      image.Width,
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
	})
}

// Update prepares a data for dbHandler
func (us RedisRepo) Update(image domain.ImageInfo) error {
	return us.dbHandler.Create(Image{
		Id:         image.Id,
		Height:     image.Height,
		Width:      image.Width,
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
	})
}

// GetById requests image by id and prepare data for app layer
func (us RedisRepo) GetById(id string) (domain.ImageInfo, error) {
	image, err := us.dbHandler.GetById(id)
	if err != nil {
		return domain.ImageInfo{}, err
	}

	return domain.ImageInfo{
		Id:         image.Id,
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
		Param: domain.Param{
			Height: image.Height,
			Width:  image.Width,
		},
	}, nil
}
