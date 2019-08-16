package idgo

import (
	"github.com/gomodule/redigo/redis"
)

// RedisStore stores allocated id to redis.
type RedisStore struct {
	maxSize     int
	conn        redis.Conn
	keyBitArray string
	keyCount    string
}

// NewRedisStore is RedisStore constructed.
func NewRedisStore(host, key string, maxSize int) (*RedisStore, error) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	store := &RedisStore{
		maxSize:     maxSize,
		conn:        conn,
		keyBitArray: key,
		keyCount:    key + "-count",
	}

	if err := store.initializeAllocatedIDCount(); err != nil {
		return nil, err
	}

	return store, nil
}

func (r *RedisStore) initializeAllocatedIDCount() error {
	if _, err := r.conn.Do("watch", r.keyCount); err != nil {
		return err
	}

	_, err := redis.Int(r.conn.Do("get", r.keyCount))
	if err != nil {
		if _, err := r.conn.Do("set", r.keyCount, 0); err != nil {
			return err
		}
		if _, err = r.conn.Do("unwatch", r.keyCount); err != nil {
			return err
		}
	}

	if _, err = r.conn.Do("unwatch", r.keyCount); err != nil {
		return err
	}
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

	if err := r.conn.Send("incr", r.keyCount); err != nil {
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

	if err := r.conn.Send("decr", r.keyCount); err != nil {
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

	if err := r.conn.Send("set", r.keyCount, "0"); err != nil {
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
	count, err := redis.Int(r.conn.Do("get", r.keyCount))
	if err != nil {
		return 0, err
	}

	return count, nil
}
