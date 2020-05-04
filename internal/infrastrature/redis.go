package infrastrature

import (
	"GoImageZip/internal/app"
	"GoImageZip/internal/interfaces"
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix/v3"
	"sync"
)

const LINKS = "links"

type RedisHandler struct {
	conn radix.Conn
	mu   sync.Mutex
}

// NewRedisConnector creates new instance of db conn
func NewRedisConnector(addr, password string, db int) (*RedisHandler, error) {
	pool, err := radix.Dial("tcp", addr, radix.DialSelectDB(db))
	if err != nil {
		return nil, err
	}

	return &RedisHandler{conn: pool}, nil
}

// Create sreates a data in store
func (rh *RedisHandler) Create(image interfaces.Image) error {
	b, err := json.Marshal(image)
	if err != nil {
		return fmt.Errorf("%s: %s", app.ErrStore, err)
	}

	rh.mu.Lock()
	err = rh.conn.Do(radix.FlatCmd(nil, "HSET", LINKS, image.Id, string(b)))
	rh.mu.Unlock()
	if err != nil {
		return fmt.Errorf("%s: %s", app.ErrStore, err)
	}

	return nil
}

// GetAllItems requests all data to store
func (rh *RedisHandler) GetAllItems() ([]interfaces.Image, error) {
	m := make(map[string]string)

	rh.mu.Lock()
	err := rh.conn.Do(radix.FlatCmd(&m, "HGETALL", LINKS))
	rh.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", app.ErrStore, err)
	}

	res := make([]interfaces.Image, 0, len(m))

	for _, value := range m {
		var r interfaces.Image

		err := json.Unmarshal([]byte(value), &r)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", app.ErrStore, err)
		}

		res = append(res, r)
	}

	return res, nil
}

// GetById requests data by id to store
func (rh *RedisHandler) GetById(id string) (interfaces.Image, error) {
	rh.mu.Lock()

	var s string
	err := rh.conn.Do(radix.FlatCmd(&s, "HGET", LINKS, id))
	rh.mu.Unlock()
	if err != nil {
		return interfaces.Image{}, fmt.Errorf("%s: %s", app.ErrStore, err)
	}

	var image interfaces.Image
	err = json.Unmarshal([]byte(s), &image)
	if err != nil {
		return interfaces.Image{}, fmt.Errorf("%s: %s", app.ErrStore, err)
	}

	return interfaces.Image{
		Id:         image.Id,
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
		Width:      image.Width,
		Height:     image.Height,
	}, nil
}
