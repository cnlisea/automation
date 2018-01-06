package config

import (
	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/http"
	"github.com/cnlisea/automation/utils"
	"strings"
)

var (
	BaseUrl = ""
	/*
		* 0 No Auth
			* 1 Bearer Token
	*/
	AuthType     = http.NoAuth
	TokenContain = "data"
	TokenKey     = "access_token"
	Token        = ""
)

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
