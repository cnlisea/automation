package automation

import (
	"testing"
)

func TestInstance(t *testing.T) {
	requests := []interface{}{
		"user_login",
		map[string]interface{}{
			"url":             "/v1/cmn/user/login",
			"method":          "post",
			"auth":            0, // 是否认证，默认值为不认证
			"json_mobile_m":   "15971485531",
			"json_password_m": "jubao56",
		},
		map[string]interface{}{
			"json_err_code_mn": int(0),
			"json_err_msg_n":   string(""),
			"json_data_m": map[string]interface{}{
				"json_access_token_m": "",
				"json_role_id_m":      0,
			},
		},
		"banner_info",
		map[string]interface{}{
			"url":        "/v1/cmn/biz/banner/info", // 请求url
			"method":     "GET",                     // 请求方法
			"auth":       0,                         // 不认证
			"author":     "lisea",                   // 作者
			"title":      "获取banner详细信息",            // 标题
			"explain":    "任何用户都可以访问该接口",            // 说明，多条时定义为[]string类型
			"query_id_m": uint32(33),
			"query_id_d": "banner id",
		},
		map[string]interface{}{
			"json_err_code_mn": int(0),
			"json_err_code_d":  "错误代码，详见错误码说明",
			"json_data_m":      "*",
			"json_data_t":      "pricelevel",
			"json_data_d":      "banner图片，详细说明见图片(banner)结构字段说明",
		},
	}

	instance := New(requests)
	if err := instance.Parse(); nil != err {
		t.Fatal("automation parse err:", err)
	}

	if err := instance.Run(); nil != err {
		t.Fatal("automation run err:", err)
	}
	t.Log("instance successfully!!!")
}

func BenchmarkInstance(b *testing.B) {
	requests := []interface{}{
		"user_login",
		map[string]interface{}{
			"url":             "/v1/cmn/user/login",
			"method":          "post",
			"json_mobile_m":   "15971485531",
			"json_password_m": "jubao56",
		},
		map[string]interface{}{
			"json_err_code_mn": int(0),
			"json_err_msg_n":   string(""),
			"json_data_m": map[string]interface{}{
				"json_access_token_m": "",
				"json_role_id_m":      0,
			},
		},
		"banner_info",
		map[string]interface{}{
			"url":        "/v1/cmn/biz/banner/info", // 请求url
			"method":     "GET",                     // 请求方法
			"author":     "lisea",                   // 作者
			"title":      "获取banner详细信息",            // 标题
			"explain":    "任何用户都可以访问该接口",            // 说明，多条时定义为[]string类型
			"query_id_m": uint32(33),
			"query_id_d": "banner id",
		},
		map[string]interface{}{
			"json_err_code_mn": int(0),
			"json_err_code_d":  "错误代码，详见错误码说明",
			"json_data_m":      "*",
			"json_data_t":      "pricelevel",
			"json_data_d":      "banner图片，详细说明见图片(banner)结构字段说明",
		},
	}

	for i := 0; i < b.N; i++ {
		instance := New(requests)
		if err := instance.Parse(); nil != err {
			b.Fatal("automation parse err:", err)
		}

		if err := instance.Run(); nil != err {
			b.Fatal("automation run err:", err)
		}
		b.Log("instance successfully!!!")
	}
}
