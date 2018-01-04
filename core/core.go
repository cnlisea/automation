package core

import (
	"github.com/cnlisea/automation/config"
)

type Instance struct {
	config   []interface{}
	Token    string
	AuthType int
	Runs     map[string]map[string]map[string]interface{}
}

func NewInstance(cfg []interface{}) *Instance {
	i := new(Instance)
	if nil == cfg {
		i.config = config.GlobalConfig
	} else {
		i.config = cfg
	}

	i.AuthType = config.AuthType
	return i
}
