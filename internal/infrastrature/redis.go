package infrastrature

import (
	"GoImageZip/internal/interfaces"
	"encoding/json"
	"github.com/mediocregopher/radix/v3"
	"sync"
)

const LINKS = "links"

type RedisHandler struct {
	conn radix.Conn
	mu   sync.Mutex
}

// todo comment
func NewRedisConnector(addr, password string, db int) (*RedisHandler, error) {
	pool, err := radix.Dial("tcp", addr, radix.DialSelectDB(db))
	if err != nil {
		return nil, err
	}

	return &RedisHandler{conn: pool}, nil
}


// todo comment
func (rh *RedisHandler) Create(image interfaces.Image) error {
	rh.mu.Lock()

	b, _ := json.Marshal(image)

	err := rh.conn.Do(radix.FlatCmd(nil, "HSET", LINKS, image.Id, string(b)))
	rh.mu.Unlock()
	return err
}

// todo comment
func (rh *RedisHandler) GetAllItems() ([]interfaces.Image, error) {
	m := make(map[string]string)

	rh.mu.Lock()
	err := rh.conn.Do(radix.FlatCmd(&m, "HGETALL", LINKS))
	rh.mu.Unlock()

	res := make([]interfaces.Image, 0, len(m))

	for _, value := range m {
		var r interfaces.Image

		err := json.Unmarshal([]byte(value), &r)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, err
}

// todo comment
func (rh *RedisHandler) GetById(id string) (interfaces.Image, error) {
	rh.mu.Lock()

	var s string
	err := rh.conn.Do(radix.FlatCmd(&s, "HGET", LINKS, id))
	rh.mu.Unlock()

	var image interfaces.Image
	err = json.Unmarshal([]byte(s), &image)

	return interfaces.Image{
		Id:         image.Id,
		OriginUrl:  image.OriginUrl,
		ResizedUrl: image.ResizedUrl,
		Width:      image.Width,
		Height:     image.Height,
	}, err
}
