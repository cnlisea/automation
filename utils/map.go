package utils

func MapMerge(ms ...map[string]interface{}) map[string]interface{} {
	if len(ms) == 0 {
		return map[string]interface{}{}
	}

	// 计算所有map的总key数
	mLen := 0
	for _, v := range ms {
		mLen = mLen + len(v)
	}

	result := make(map[string]interface{}, mLen)
	for _, m := range ms {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
