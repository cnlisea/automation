package automation

import (
	"github.com/cnlisea/automation/config"
	"github.com/cnlisea/automation/core"
)

type Instance interface {
	// parse config
	Parse() error
	// test and generate docs
	Run() error
}

func New(cfg []interface{}) Instance {
	return core.NewInstance(cfg)
}

func ParseConfigFile(path string) ([]interface{}, error) {
	return config.ParseConfigFile(path)
}
