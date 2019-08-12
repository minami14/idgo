package idgo

import (
	"github.com/gomodule/redigo/redis"
)

// RedisStore stores allocated id to redis.
type RedisStore struct {
	maxSize     int
	conn        redis.Conn
	keyBitArray string
	keyMax      string
}

// NewRedisStore is RedisStore constructed.
func NewRedisStore(host, key string, maxSize int) (*RedisStore, error) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		maxSize:     maxSize,
		conn:        conn,
		keyBitArray: key,
		keyMax:      key + "_max",
	}, nil
}

func (r *RedisStore) isAllocated(id int) (bool, error) {
	result, err := r.conn.Do("bitfield", r.keyBitArray, "get", "i1", id)
	return redis.Bool(result.([]interface{})[0], err)
}

func (r *RedisStore) allocate(id int) error {
	if err := r.conn.Send("multi"); err != nil {
		return err
	}

	if err := r.conn.Send("bitfield", r.keyBitArray, "set", "i1", id, "1"); err != nil {
		return err
	}

	if err := r.conn.Send("incr", r.keyMax); err != nil {
		return err
	}

	if _, err := r.conn.Do("exec"); err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) free(id int) error {
	if err := r.conn.Send("multi"); err != nil {
		return err
	}

	if err := r.conn.Send("bitfield", r.keyBitArray, "set", "i1", id, "0"); err != nil {
		return err
	}

	if err := r.conn.Send("decr", r.keyMax); err != nil {
		return err
	}

	if _, err := r.conn.Do("exec"); err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) freeAll() error {
	if err := r.conn.Send("multi"); err != nil {
		return err
	}

	if err := r.conn.Send("del", r.keyBitArray); err != nil {
		return err
	}

	if err := r.conn.Send("set", r.keyMax, "0"); err != nil {
		return err
	}

	if _, err := r.conn.Do("exec"); err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) getMaxSize() int {
	return r.maxSize
}

func (r *RedisStore) getAllocatedIDCount() (int, error) {
	if _, err := r.conn.Do("watch", r.keyMax); err != nil {
		return 0, err
	}

	count, err := redis.Int(r.conn.Do("get", r.keyMax))
	if err != nil {
		r.conn.Do("set", r.keyMax, 0)
		r.conn.Do("unwatch", r.keyMax)

		count, err := redis.Int(r.conn.Do("get", r.keyMax))
		if err != nil {
			return 0, err
		}
		return count, nil
	}

	r.conn.Do("unwatch", r.keyMax)
	return count, nil
}
