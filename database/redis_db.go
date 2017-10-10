package database

import (
	"fmt"
	"strconv"
	"strings"
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
//   user:@username -> {port, active}
type RedisDatabase struct {
	r *redis.Client
}

func (db *RedisDatabase) Write(r *Record) error {
	usageKey := fmt.Sprintf("usage:%s", r.Username)
	dir := strconv.Itoa(r.Dir)
	used := int64(r.BytesUsed)

	t := time.Unix(r.Time, 0)
	daily := fmt.Sprintf("usage:%s:%s", r.Username, t.Format("2006-01-02"))
	mouthly := fmt.Sprintf("usage:%s:%s", r.Username, t.Format("2006-01"))
	yearly := fmt.Sprintf("usage:%s:%s", r.Username, t.Format("2006"))

	pipe := db.r.TxPipeline()
	pipe.HIncrBy(usageKey, dir, used)
	pipe.HIncrBy(daily, dir, used)
	pipe.HIncrBy(mouthly, dir, used)
	pipe.HIncrBy(yearly, dir, used)

	_, err := pipe.Exec()
	return err
}

func (db *RedisDatabase) BindPort(b *Binding) error {
	userKey := fmt.Sprintf("user:%s", b.Username)
	cmd := db.r.HMSet(userKey, map[string]interface{}{
		"port":    b.Port,
		"active":  b.Active,
		"backend": b.Backend,
	})
	return cmd.Err()
}

func (db *RedisDatabase) GetAllActiveBinding() ([]*Binding, error) {
	users := db.r.Keys("user:*")
	if users.Err() != nil {
		return nil, users.Err()
	}
	var rv []*Binding

	for _, uid := range users.Val() {
		u := db.r.HGetAll(uid)
		if u.Err() != nil {
			continue
		}
		uVal := u.Val()
		if active, _ := strconv.ParseBool(uVal["active"]); !active {
			continue
		}
		rv = append(rv, &Binding{
			Username: strings.Split(uid, ":")[1],
			Port:     uVal["port"],
			Active:   true,
			Backend:  uVal["backend"],
		})
	}
	return rv, nil
}
