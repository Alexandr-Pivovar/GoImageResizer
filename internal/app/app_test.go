package app

import (
	"GoImageZip/internal/app/mocks"
	"GoImageZip/internal/domain"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestImageService_Resize(t *testing.T) {
	// arrange
	cases := []struct {
		testName   string
		repo       *mocks.ImageRepository
		resizer    *mocks.Resizer
		cloud      *mocks.Clouder
		uuid       func() string
		decode     func([]byte, []byte) (int, error)
		imageParam domain.Image
		wantImage  domain.ImageInfo
		wantErr    error
	}{
		{
			testName: "Should return error if Decode func return error",
			uuid:     func() string { return "1" },
			decode:   func(dst []byte, src []byte) (int, error) { return 0, errors.New("decode error") },
			wantErr:  errors.New("decode error"),
		},

		{
			testName: "Should return error if resizer resizes img",
			uuid:     func() string { return "1" },
			decode:   func(dst []byte, src []byte) (int, error) { return 0, nil },
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", domain.Image{}).
					Return(domain.Image{}, errors.New("resizer error")).
					Once()
				return m
			}(),
			wantErr:  errors.New("resizer error"),
		},

		{
			testName: "Should return error if cloud.Save try to save origin img",
			uuid:     func() string { return "1" },
			decode:   func(dst []byte, src []byte) (int, error) { return 0, nil },
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", domain.Image{}).
					Return(domain.Image{}, nil).
					Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("Save", "1", domain.Image{}).
					Return("", errors.New("cloud error")).
					Once()
				return m
			}(),
			wantErr:  errors.New("cloud error"),
		},

		{
			testName: "Should return error if second cloud.Save return error",
			uuid:     func() string { return "1" },
			decode:   func(dst []byte, src []byte) (int, error) { return 0, nil },
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", domain.Image{}).
					Return(domain.Image{Data:[]byte{100}}, nil).
					Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("Save", "1", domain.Image{}).
					Return("", nil).
					Once()
				m.On("Save", "1", domain.Image{Data: []byte{100}}).
					Return("", errors.New("cloud error")).
					Once()
				return m
			}(),
			wantErr:  errors.New("cloud error"),
		},

		{
			testName: "Should return error if repo.Create return error",
			uuid:     func() string { return "1" },
			decode:   func(dst []byte, src []byte) (int, error) { return 0, nil },
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", mock.Anything).
					Return(domain.Image{}, nil).
					Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("Save", "1", mock.Anything).
					Return("x.com", nil).
					Once()
				m.On("Save", "1", mock.Anything).
					Return("z.com", nil).
					Once()
				return m
			}(),
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("Create", domain.ImageInfo{
					Id: "1",
					OriginUrl: "x.com",
					ResizedUrl: "z.com",
					Param: domain.Param{
						Format: "png",
						Width: 10,
						Height: 15,
					},
				}).
					Return(errors.New("repo error")).
					Once()
				return m
			}(),
			imageParam: domain.Image{
				Data: nil,
				Param: domain.Param{
					Width: 10,
					Height: 15,
					Format: "png",
				},
			},
			wantImage: domain.ImageInfo{
				Id: "1",
				Param: domain.Param{
					Width: 10,
					Height: 15,
					Format: "png",
				},
				OriginUrl: "x.com",
				ResizedUrl: "z.com",
			},
			wantErr:  errors.New("repo error"),
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			i := &ImageService{
				repo:    c.repo,
				resizer: c.resizer,
				cloud:   c.cloud,
				uuid:    c.uuid,
				decode:  c.decode,
			}

			// act
			gotImage, gotErr := i.Resize(c.imageParam)

			// assert
			assert.Equal(t, c.wantImage, gotImage)
			assert.Equal(t, c.wantErr, gotErr)

			if c.resizer != nil {
				c.resizer.AssertExpectations(t)
			}
			if c.cloud != nil {
				c.cloud.AssertExpectations(t)
			}
			if c.repo != nil {
				c.repo.AssertExpectations(t)
			}
		})
	}
}

func TestImageService_GetHistory(t *testing.T) {
	// arrange
	cases := []struct {
		testName   string
		repo       *mocks.ImageRepository
		wantImages []domain.ImageInfo
		wantErr    error
	}{
		{
			testName: "Should return error if repo.GetHistory func return error",
			repo: func() *mocks.ImageRepository{
				m := &mocks.ImageRepository{}
				m.On("GetHistory").
					Return(nil, errors.New("repo error"))
				return m
			}(),
			wantErr:  errors.New("repo error"),
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			i := &ImageService{
				repo: c.repo,
			}

			// act
			gotImages, gotErr := i.GetHistory()

			// assert
			assert.Equal(t, c.wantImages, gotImages)
			assert.Equal(t, c.wantErr, gotErr)

			if c.repo != nil {
				c.repo.AssertExpectations(t)
			}
		})
	}
}

func TestImageService_GetById(t *testing.T) {
	// arrange
	cases := []struct {
		testName  string
		repo      *mocks.ImageRepository
		id        string
		wantImage domain.ImageInfo
		wantErr   error
	}{
		{
			testName: "Should return error if repo.GetById func return error",
			id:       "1",
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("GetById", "1").
					Return(domain.ImageInfo{}, errors.New("repo error"))
				return m
			}(),
			wantErr: errors.New("repo error"),
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			i := &ImageService{
				repo: c.repo,
			}

			// act
			gotImage, gotErr := i.GetById(c.id)

			// assert
			assert.Equal(t, c.wantImage, gotImage)
			assert.Equal(t, c.wantErr, gotErr)

			if c.repo != nil {
				c.repo.AssertExpectations(t)
			}
		})
	}
}

func TestImageService_Update(t *testing.T) {
	// arrange
	cases := []struct {
		testName  string
		repo      *mocks.ImageRepository
		cloud     *mocks.Clouder
		resizer   *mocks.Resizer
		imgParam  domain.ImageInfo
		wantImage domain.ImageInfo
		wantErr   error
	}{

		{
			testName: "Should return error if repo.GetById func return error",
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("GetById", "1").
					Return(domain.ImageInfo{}, errors.New("repo error")).
					Once()
				return m
			}(),
			imgParam: domain.ImageInfo{Id: "1"},
			wantErr:  errors.New("repo error"),
		},

		{
			testName: "Should return error if repo.GetImage func return error",
			imgParam: domain.ImageInfo{Id: "1"},
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("GetById", "1").
					Return(domain.ImageInfo{
						OriginUrl: "x.com",
					}, nil).
					Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("GetImage", "x.com").
					Return(nil, errors.New("repo error")).
					Once()
				return m
			}(),
			wantErr: errors.New("repo error"),
		},

		{
			testName: "Should return error if resizer.Do func return error",
			imgParam: domain.ImageInfo{
				Id: "1",
				Param: domain.Param{
					Height: 10,
					Width:  15,
				},
			},
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("GetById", "1").
					Return(domain.ImageInfo{
						OriginUrl: "x.com",
						Param: domain.Param{
							Format: "png",
						},
					}, nil).
					Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("GetImage", "x.com").
					Return([]byte{100}, nil).
					Once()
				return m
			}(),
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", domain.Image{
					Data: []byte{100},
					Param: domain.Param{
						Width:  15,
						Height: 10,
						Format: "png",
					},
				}).
					Return(domain.Image{}, errors.New("resizer error")).
					Once()
				return m
			}(),
			wantErr: errors.New("resizer error"),
		},

		{
			testName: "Should return error if cloud.Save func return error",
			imgParam: domain.ImageInfo{Id: "1"},
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("GetById", mock.Anything).
					Return(domain.ImageInfo{}, nil).
					Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("GetImage", mock.Anything).
					Return(nil, nil).
					Once()
				m.On("Save", "1", domain.Image{}).
					Return("", errors.New("cloud error")).
					Once()
				return m
			}(),
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", mock.Anything).
					Return(domain.Image{}, nil).
					Once()
				return m
			}(),
			wantErr: errors.New("cloud error"),
		},

		{
			testName: "Should return error if repo.Update func return error",
			imgParam: domain.ImageInfo{Id: "1"},
			repo: func() *mocks.ImageRepository {
				m := &mocks.ImageRepository{}
				m.On("GetById", mock.Anything).
					Return(domain.ImageInfo{
						Id:         "1",
						OriginUrl:  "origin",
						ResizedUrl: "resized",
						Param: domain.Param{
							Width:  10,
							Height: 12,
							Format: "png",
						},
					}, nil).Once()
				m.On("Update", domain.ImageInfo{
					Id:         "1",
					OriginUrl:  "origin",
					ResizedUrl: "new url",
					Param: domain.Param{
						Width:  0,
						Height: 0,
						Format: "png",
					},
				}).
					Return(errors.New("cloud error")).Once()
				return m
			}(),
			cloud: func() *mocks.Clouder {
				m := &mocks.Clouder{}
				m.On("GetImage", mock.Anything).
					Return(nil, nil).
					Once()
				m.On("Save", "1", domain.Image{}).
					Return("new url", nil).
					Once()
				return m
			}(),
			resizer: func() *mocks.Resizer {
				m := &mocks.Resizer{}
				m.On("Do", mock.Anything).
					Return(domain.Image{}, nil).
					Once()
				return m
			}(),
			wantImage: domain.ImageInfo{
				Id:         "1",
				OriginUrl:  "origin",
				ResizedUrl: "new url",
				Param: domain.Param{
					Width:  0,
					Height: 0,
					Format: "png",
				},
			},
			wantErr: errors.New("cloud error"),
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			i := &ImageService{
				repo:    c.repo,
				cloud:   c.cloud,
				resizer: c.resizer,
			}

			// act
			gotImage, gotErr := i.Update(c.imgParam)

			// assert
			assert.Equal(t, c.wantImage, gotImage)
			assert.Equal(t, c.wantErr, gotErr)

			if c.repo != nil {
				c.repo.AssertExpectations(t)
			}
			if c.cloud != nil {
				c.cloud.AssertExpectations(t)
			}
			if c.resizer != nil {
				c.resizer.AssertExpectations(t)
			}
		})
	}
}
