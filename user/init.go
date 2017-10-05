package user

import (
	"fmt"

	"github.com/ibigbug/ss-account/config"
	"github.com/ibigbug/ss-account/database"
)

// DefaultManaged is default user manager
var DefaultManaged = &Managed{}

// DefaultStorage is default storage provider
var DefaultStorage database.Storage

// Init does storage initialization
func Init() {
	redisHost, port, pass, db := config.GetRedisOptions()
	DefaultStorage = database.NewAsyncStorage(database.NewRedisDb(&database.RedisOptions{
		Addr:     fmt.Sprintf("%s:%s", redisHost, port),
		Password: pass,
		DB:       db,
	}))
}
