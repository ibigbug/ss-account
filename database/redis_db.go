package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisDb(opts *RedisOptions) *RedisDatabase {
	r := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
	return &RedisDatabase{
		r: r,
	}
}

// RedisDatabase storage
// Scheme
// accounting:
//   usage:@username -> {up, down}
//   usage:@username:@day -> {up, down}
//   usage:@username:@month -> {up, down}
//   usage:@username:@year -> {up, down}
// user:
//   user:@id -> {username, port, active}
type RedisDatabase struct {
	r *redis.Client
}

func (db *RedisDatabase) Write(r *Record) {
	usageKey := fmt.Sprintf("usage:%s", r.Username)
	dir := strconv.Itoa(r.Dir)
	used := int64(r.BytesUsed)

	t := time.Unix(r.Time, 0)
	daily := fmt.Sprintf("usage:%s:%s", r.Username, t.Format("2006-01-02"))
	mouthly := fmt.Sprintf("usage:%s:%s", r.Username, t.Format("2006-01"))
	yearly := fmt.Sprintf("usage:%s:%s", r.Username, t.Format("2006"))

	db.r.HIncrBy(usageKey, dir, used)
	db.r.HIncrBy(daily, dir, used)
	db.r.HIncrBy(mouthly, dir, used)
	db.r.HIncrBy(yearly, dir, used)
}

func (db *RedisDatabase) BindPort(b *Binding) {
	uid = db.r.Incr("user:id")
}
