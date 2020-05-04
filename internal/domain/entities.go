package domain

type Image struct {
	Data  []byte
	Param
}

type Param struct {
	Width  uint
	Height uint
	Format string
}

type ImageInfo struct {
	Id         string
	OriginUrl  string
	ResizedUrl string
	Param
}
