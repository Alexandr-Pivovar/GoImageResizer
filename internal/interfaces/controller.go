package interfaces

import (
	"GoImageZip/internal/app"
	"GoImageZip/internal/domain"
	"fmt"
	"github.com/labstack/echo/v4"
)

type ResizeParam struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

type ResizeImgResp struct {
	OriginUrl  string `json:"origin_url"`
	ResizedUrl string `json:"resized_url"`
	ResizeParam
}

type ResizeImg struct {
	Data   string `json:"data"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Format string `json:"format"`
}

type UpdateImg struct {
	Id string `json:"id"`
	ResizeParam
}

type Controller struct {
	app app.ResizeServicer
}

func NewController(service app.ResizeServicer) Controller {
	return Controller{
		app: service,
	}
}

func (c Controller) Run(addr string) {
	e := echo.New()

	g := e.Group("api/v1")

	g.POST("/resize", c.Resize)
	g.GET("/history", c.GetHistory)
	g.GET("/history/:id", c.GetById)
	g.POST("/history/:id", c.Update)

	e.Logger.Fatal(e.Start(addr))
}

// todo comment
func (c Controller) Resize(ctx echo.Context) error {
	r := &ResizeImg{}
	err := ctx.Bind(r)
	if err != nil {
		fmt.Println(err)
		return echo.ErrBadRequest
	}

	imageInfo, err := c.app.Resize(domain.Image{
		Data: []byte(r.Data),
		Param: domain.Param{
			Height: r.Height,
			Width:  r.Width,
			Format: r.Format,
		},
	})
	if err != nil {
		fmt.Println(err)
		return echo.ErrBadRequest
	}

	res := ResizeImgResp{
		OriginUrl:  imageInfo.OriginUrl,
		ResizedUrl: imageInfo.ResizedUrl,
		ResizeParam: ResizeParam{
			Height: imageInfo.Height,
			Width:  imageInfo.Width,
		},
	}

	return ctx.JSON(200, res)
}

// todo comment
func (c Controller) GetHistory(ctx echo.Context) error {
	images, err := c.app.GetHistory()
	if err != nil {
		return echo.ErrInternalServerError
	}

	r := make(map[string]ResizeImgResp, len(images))
	for _, image := range images {
		r[image.Id] = ResizeImgResp{
			OriginUrl:  image.OriginUrl,
			ResizedUrl: image.ResizedUrl,
			ResizeParam: ResizeParam{
				Height: image.Height,
				Width:  image.Width,
			},
		}
	}

	return ctx.JSON(200, r)
}

// todo comment
func (c Controller) GetById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	image, err := c.app.GetById(id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	r := ResizeImgResp{
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
		ResizeParam: ResizeParam{
			Width:  image.Width,
			Height: image.Height,
		},
	}

	return ctx.JSON(200, r)
}

// todo comment
func (c Controller) Update(ctx echo.Context) error {
	u := &UpdateImg{}
	err := ctx.Bind(u)
	if err != nil {
		fmt.Println(err)
		return echo.ErrBadRequest
	}

	image, err := c.app.Update(domain.ImageInfo{
		Id: u.Id,
		Param: domain.Param{
			Width:  u.Width,
			Height: u.Height,
		},
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	r := ResizeImgResp{
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
		ResizeParam: ResizeParam{
			Width:  image.Width,
			Height: image.Height,
		},
	}

	return ctx.JSON(200, r)
}
