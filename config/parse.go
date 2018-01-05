package config

import (
	"errors"
	"strings"

	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/http"
)

func ParseConfig(cfg []interface{}) (map[string]map[string]map[string]interface{}, error) {
	if len(cfg) < 0 {
		return nil, errors.New("must use a config RunProcess")
	}

	runs := make(map[string]map[string]map[string]interface{}, 0)
	// zone config RunProcess
	for i := 0; i < len(cfg); i += 3 {
		api, ok := cfg[i].(string)
		if !ok {
			continue
		}
		apis := strings.Split(api, "_")
		if len(apis) < 2 {
			continue
		}

		/*
			// 接口类型验证
			switch apis[0] {
			case "user":
			case "banner":
			case "bizarticle":
			case "product":
			case "productprice":
			case "pricelevel":
			default:
				return nil, errors.New(fmt.Sprintln("config [", apis[0], "] bad parameter type"))
		}*/

		// get request url
		uri, ok := (cfg[i+1].(map[string]interface{}))["url"].(string)
		if !ok {
			continue
		}
		if !strings.Contains(uri, "http") {
			uri = BaseUrl + uri
		}

		// get request method
		method, ok := (cfg[i+1].(map[string]interface{}))["method"].(string)
		if !ok || !http.ValidMethod(method) {
			continue
		}

		// 类型
		if _, ok := runs[apis[0]]; !ok {
			runs[apis[0]] = make(map[string]map[string]interface{})
		}

		if _, ok := runs[apis[0]][apis[1]]; !ok {
			runs[apis[0]][apis[1]] = make(map[string]interface{})
		}

		runs[apis[0]][apis[1]][constant.ConfigRequestUrl] = uri
		runs[apis[0]][apis[1]][constant.ConfigRequestMethod] = strings.ToUpper(method)
		runs[apis[0]][apis[1]][constant.ConfigRequest] = cfg[i+1]
		runs[apis[0]][apis[1]][constant.ConfigResponse] = cfg[i+2]
	}

	//if len(config) == 0 || (AuthType != NoAuth && (len(config) < 2)) {
	if len(runs) == 0 {
		return nil, errors.New("parameter missing for config")
	}

	return runs, nil
}
