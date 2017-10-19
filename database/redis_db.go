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

func (db *RedisDatabase) GetUser(uid string) (*Binding, error) {
	u := db.r.HGetAll(uid)
	if u.Err() != nil {
		return nil, u.Err()
	}
	uVal := u.Val()
	active, _ := strconv.ParseBool(uVal["active"])

	return &Binding{
		Username: strings.Split(uid, ":")[1],
		Port:     uVal["port"],
		Active:   active,
		Backend:  uVal["backend"],
	}, nil
}

func (db *RedisDatabase) GetUsage(uid string) (*Usage, error) {
	username := strings.SplitN(uid, ":", 2)[1]

	usageKey := fmt.Sprintf("usage:%s", username)
	usage := db.r.HGetAll(usageKey)
	if usage.Err() != nil {
		return nil, usage.Err()
	}

	up, _ := strconv.ParseInt(usage.Val()["0"], 10, 0)
	down, _ := strconv.ParseInt(usage.Val()["1"], 10, 0)

	return &Usage{
		Total: u{
			0: up,
			1: down,
		},
	}, nil
}

func (db *RedisDatabase) GetAllActiveBinding() ([]*Binding, error) {
	users := db.r.Keys("user:*")
	if users.Err() != nil {
		return nil, users.Err()
	}
	var rv []*Binding

	for _, uid := range users.Val() {
		u, err := db.GetUser(uid)
		if err != nil {
			return nil, err
		}
		rv = append(rv, u)
	}
	return rv, nil
}

func (db *RedisDatabase) GetAllUserUsage() ([]*UserUsage, error) {
	users := db.r.Keys("user:*")
	if users.Err() != nil {
		return nil, users.Err()
	}
	var rv []*UserUsage
	for _, uid := range users.Val() {
		if user, err := db.GetUser(uid); err == nil {
			if usage, err := db.GetUsage(uid); err == nil {
				rv = append(rv, &UserUsage{
					Binding: user,
					Usage:   usage,
				})
			}
		}
	}
	return rv, nil
}
