package core

import (
	"github.com/cnlisea/automation/config"
)

func (i *Instance) Parse() (err error) {
	i.Runs, err = config.ParseConfig(i.config)
	return
}
