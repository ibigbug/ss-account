package user

import (
	"fmt"
	"log"
	"os"

	"github.com/ibigbug/ss-account/config"
	"github.com/ibigbug/ss-account/database"
	"github.com/ibigbug/ss-account/utils"
)

var logger *log.Logger

// DefaultManaged is default user manager
var DefaultManaged = &Managed{}

// DefaultStorage is default storage provider
var DefaultStorage database.Storage

// Init does storage initialization
func Init() {
	if utils.IsaTty(os.Stdout.Fd()) {
		logger = log.New(os.Stderr, "\x1b[95m[User]: \x1b[0m", log.LstdFlags)
	} else {
		logger = log.New(os.Stderr, "[User]: ", log.LstdFlags)
	}
	redisHost, port, pass, db := config.GetRedisOptions()
	DefaultStorage = database.NewAsyncStorage(database.NewRedisDb(&database.RedisOptions{
		Addr:     fmt.Sprintf("%s:%s", redisHost, port),
		Password: pass,
		DB:       db,
	}))

	restoreBindings()
}

func restoreBindings() {
	if activeUsers, err := DefaultStorage.GetAllActiveBinding(); err == nil {
		for _, b := range activeUsers {
			m := &Manager{
				Username: b.Username,
				Port:     b.Port,
				Backend:  b.Backend,
			}
			if err := m.Start(); err != nil {
				logger.Printf("failed to start binding %s", err)
			} else {
				DefaultManaged.Add(m)
			}
		}
	} else {
		logger.Printf("failed to restore bindings: %v", err)
	}
}
