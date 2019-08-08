package idgo

import "github.com/gomodule/redigo/redis"

// RedisStore stores allocated id to redis.
type RedisStore struct {
	maxSize int
	conn    redis.Conn
}

// NewRedisStore is RedisStore constructed.
func NewRedisStore(host string, maxSize int) (*RedisStore, error) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		maxSize: maxSize,
		conn:    conn,
	}, nil
}

const key = "AllocatedID"

func (r *RedisStore) isAllocated(id int) (bool, error) {
	return redis.Bool(r.conn.Do("bitfield", key, "get", "i1", id))
}

func (r *RedisStore) allocate(id int) error {
	_, err := r.conn.Do("bitfield", key, "set", "i1", id, "1")
	return err
}

func (r *RedisStore) free(id int) error {
	_, err := r.conn.Do("bitfield", key, "set", "i1", id, "0")
	return err
}

func (r *RedisStore) freeAll() error {
	_, err := r.conn.Do("del", key)
	return err
}

func (r *RedisStore) getMaxSize() int {
	return 0
}
