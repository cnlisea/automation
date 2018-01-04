package core

import (
	"automation/config"
)

func (i *Instance) Parse() (err error) {
	i.Runs, err = config.ParseConfig(i.config)
	return
}
