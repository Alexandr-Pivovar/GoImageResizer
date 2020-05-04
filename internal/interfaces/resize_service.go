package interfaces

import (
	"GoImageZip/internal/app"
	"GoImageZip/internal/domain"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image/png"
)

type ImageResize struct{}

func (ir ImageResize) Do(image domain.Image) (i domain.Image, err error) {
	im, err := png.Decode(bytes.NewReader(image.Data))
	if err != nil {
		return i, fmt.Errorf("%s: %s", app.ErrDecode, err)
	}

	dstImage128 := imaging.Resize(im, int(image.Width), int(image.Height), imaging.Lanczos)

	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, dstImage128, imaging.PNG)
	if err != nil {
		return i, fmt.Errorf("%s: %s", app.ErrDecode, err)
	}

	i.Data = make([]byte, buf.Len())
	_, err = buf.Read(i.Data)
	if err != nil {
		return i, fmt.Errorf("%s: %s", app.ErrDecode, err)
	}

	return
}
