package idgo

import "github.com/gomodule/redigo/redis"

// RedisStore stores allocated id to redis.
type RedisStore struct {
	maxSize int
	conn    redis.Conn
	key     string
}

// NewRedisStore is RedisStore constructed.
func NewRedisStore(host, key string, maxSize int) (*RedisStore, error) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		maxSize: maxSize,
		conn:    conn,
		key:     key,
	}, nil
}

func (r *RedisStore) isAllocated(id int) (bool, error) {
	result, err := r.conn.Do("bitfield", r.key, "get", "i1", id)
	return redis.Bool(result.([]interface{})[0], err)
}

func (r *RedisStore) allocate(id int) error {
	_, err := r.conn.Do("bitfield", r.key, "set", "i1", id, "1")
	return err
}

func (r *RedisStore) free(id int) error {
	_, err := r.conn.Do("bitfield", r.key, "set", "i1", id, "0")
	return err
}

func (r *RedisStore) freeAll() error {
	_, err := r.conn.Do("del", r.key)
	return err
}

func (r *RedisStore) getMaxSize() int {
	return r.maxSize
}
