package config

import (
	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/http"
	"github.com/cnlisea/automation/utils"
	"encoding/json"
	"strconv"
	"strings"
)

var (
	BaseUrl = "http://test.jubao56.com:5000"
	/*
		* 0 No Auth
			* 1 Bearer Token
	*/
	AuthType     = http.BearerTokenAuth
	TokenContain = "data"
	TokenKey     = "access_token"
	Token        = ""
)

var (
	GlobalConfig = []interface{}{}
)

func ParseConfigFile(path string) ([]interface{}, error) {
	// read all
	b, err := utils.ReadFileAll(path)
	if nil != err {
		return nil, err
	}

	var config []interface{}

	data := make(map[string]interface{})
	if err = json.Unmarshal(b, &data); nil != err {
		return nil, err
	}

	if 0 == len(data) {
		return []interface{}{}, nil
	}

	for k, v := range data {
		switch strings.ToLower(k) {
		case "baseurl":
			fallthrough
		case "base_url":
			fallthrough
		case "url":
			BaseUrl = utils.ToString(v)
		case "authtype":
			fallthrough
		case "auth_type":
			fallthrough
		case "auth":
			authType, err := strconv.Atoi(utils.ToString(v))
			if nil != err {
				return nil, err
			}
			AuthType = authType
		case "token_contain":
			fallthrough
		case "tokencontain":
			TokenContain = utils.ToString(v)
		case "tokenkey":
			fallthrough
		case "token_key":
			TokenKey = utils.ToString(v)
		case "token":
			Token = utils.ToString(v)
		default:
			cData, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			if !strings.Contains(k, "_") {
				for kk, vv := range cData {
					ccData, ok := vv.(map[string]interface{})
					if !ok {
						continue
					}

					var req, res interface{}
					if val, ok := ccData["req"]; ok {
						req = val
					}

					if val, ok := ccData["request"]; ok {
						req = val
					}

					if val, ok := ccData["res"]; ok {
						res = val
					}

					if val, ok := ccData["response"]; ok {
						res = val
					}

					if nil == req || nil == res {
						//TODO logs req or res is nil
						continue
					}

					req, res = RecoverType(req), RecoverType(res)
					/*fmt.Println("req:")
					for k, v := range req.(map[string]interface{}){
						t := reflect.TypeOf(v)
						fmt.Println(k, " ", t.String())
					}
					fmt.Println("res:")
					for k, v := range res.(map[string]interface{}){
						t := reflect.TypeOf(v)
						fmt.Println(k, " ", t.String())
					}*/
					config = append(config, k+"_"+kk, req, res)
				}
			} else {
				var req, res interface{}
				if val, ok := cData["req"]; ok {
					req = val
				}

				if val, ok := cData["request"]; ok {
					req = val
				}

				if val, ok := cData["res"]; ok {
					res = val
				}

				if val, ok := cData["response"]; ok {
					res = val
				}

				if nil == req || nil == res {
					//TODO logs req or res is nil
					continue
				}

				if v, ok := req.(map[string]interface{}); ok {
					_ = v
				}

				if v, ok := res.(map[string]interface{}); ok {
					_ = v
				}

				req, res = RecoverType(req), RecoverType(res)
				/*fmt.Println("req:")
				for k, v := range req.(map[string]interface{}){
					t := reflect.TypeOf(v)
					fmt.Println(k, " ", t.String())
				}
				fmt.Println("res:")
				for k, v := range res.(map[string]interface{}){
					t := reflect.TypeOf(v)
					fmt.Println(k, " ", t.String())
				}*/
				config = append(config, k, req, res)
			}
		}
	}

	// save global config
	GlobalConfig = config

	return config, nil
}

func RecoverType(data interface{}) interface{} {
	switch data.(type) {
	case map[string]interface{}:
		data = MapToType(data.(map[string]interface{}))
	case []interface{}:
		data = SliceToType(data.([]interface{}))
	}
	return data
}

func MapToType(data map[string]interface{}) interface{} {
	for k, v := range data {
		if !strings.Contains(k, "_") {
			continue
		}

		ks := strings.Split(k, "_")
		switch strings.ToLower(ks[len(ks)-1]) {
		case "m":
		case "n":
		case "mn":
		default:
			continue
		}

		keyPre := ks[0]
		for i := 1; i < len(ks)-1; i++ {
			keyPre = keyPre + "_" + ks[i]
		}

		switch v.(type) {
		case float64:
			t, ok := data[keyPre+constant.TypeSuffix]
			_, ok2 := t.(string)
			if !ok || !ok2 {
				t = utils.TypeString(v)
			}

			data[k] = utils.StringToType(t.(string), v)
		case map[string]interface{}:
			data[k] = MapToType(v.(map[string]interface{}))
		case []interface{}:
			data[k] = SliceToType(v.([]interface{}))
		}
	}
	return data
}

func SliceToType(data []interface{}) interface{} {
	for k, v := range data {
		switch v.(type) {
		case string:
			data[k] = utils.ToString(v)
		case float64:
			data[k] = utils.ToString(v)
		case map[string]interface{}:
			data[k] = MapToType(v.(map[string]interface{}))
		case []interface{}:
			data[k] = SliceToType(v.([]interface{}))
		}
	}
	return data
}
