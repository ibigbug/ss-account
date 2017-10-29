package user

import (
	"github.com/ibigbug/ss-account/config"
)

// AddOneUser register a user with specific backend
// a random port will be used upon the port range if
// no port given
func AddOneUser(backend, username, port string) (p string, err error) {
	if port == "" {
		p = config.GetRandomPort()
	} else {
		p = port
	}

	m := Manager{
		Username: username,
		Backend:  backend,
		Port:     p,
	}

	if err = m.Start(); err != nil {
		return
	}

	DefaultManaged.Add(&m)

	return
}

func GetAllUserUsage() ([]*Manager, error) {
	usage, err := DefaultStorage.GetAllUserUsage()
	if err != nil {
		return nil, err
	}
	var rv []*Manager
	for _, u := range usage {
		rv = append(rv, &Manager{
			Username:      u.Username,
			Port:          u.Port,
			Backend:       u.Backend,
			Active:        u.Active,
			BytesDownload: u.Total[1],
			BytesUpload:   u.Total[0],
		})
	}
	return rv, nil
}
