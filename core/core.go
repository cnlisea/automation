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
	return &Instance{
		config:   cfg,
		AuthType: config.AuthType,
	}
}
