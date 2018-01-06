package utils

import (
	"math"
	"strconv"
	"strings"
	"time"
)

func ToType(srcType interface{}, data interface{}) interface{} {
	switch srcType.(type) {
	case string:
		if srcType == "*" {
			return data
		}
		if data == nil {
			data = ""
		}
		return data.(string)
	case float32:
		if data == nil {
			data = float64(0)
		}
		return float32(data.(float64))
	case float64:
		if data == nil {
			data = float64(0)
		}
		return data.(float64)
	case int:
		if data == nil {
			data = float64(0)
		}
		return int(math.Trunc(data.(float64)))
	case int8:
		if data == nil {
			data = float64(0)
		}
		return int8(math.Trunc(data.(float64)))
	case int16:
		if data == nil {
			data = float64(0)
		}
		return int16(math.Trunc(data.(float64)))
	case int32:
		if data == nil {
			data = float64(0)
		}
		return int32(math.Trunc(data.(float64)))
	case int64:
		if data == nil {
			data = float64(0)
		}
		return int64(math.Trunc(data.(float64)))
	case uint8:
		if data == nil {
			data = float64(0)
		}
		return uint8(math.Trunc(data.(float64)))
	case uint16:
		if data == nil {
			data = float64(0)
		}
		return uint16(math.Trunc(data.(float64)))
	case uint32:
		if data == nil {
			data = float64(0)
		}
		return uint32(math.Trunc(data.(float64)))
	case uint64:
		if data == nil {
			data = float64(0)
		}
		return uint64(math.Trunc(data.(float64)))
	case bool:
		if data == nil {
			data = false
		}
		return data.(bool)
	case []interface{}:
		if data == nil {
			data = []interface{}{}
		}
		result := make([]interface{}, len(srcType.([]interface{})))
		for i, v := range srcType.([]interface{}) {
			result[i] = ToType(v, data.([]interface{})[i])
		}
		return result
	case map[string]interface{}:
		if data == nil {
			data = make(map[string]interface{})
		}
		result := make(map[string]interface{})
		for k, v := range srcType.(map[string]interface{}) {
			ks := strings.Split(k, "_")
			if len(ks) < 3 {
				continue
			}

			key := ks[1]
			if len(ks) > 3 {
				for i := 2; i < len(ks)-1; i++ {
					key = key + "_" + ks[i]
				}
			}
			result[key] = ToType(v, data.(map[string]interface{})[key])
		}
		return result
	}

	return data
}

func StringToType(t string, data interface{}) (r interface{}) {
	// recover
	defer func() {
		if err := recover(); nil != err {
			r = data
		}
	}()

	switch strings.ToLower(t) {
	case "string":
		if data == nil {
			data = ""
		}
		r = data.(string)
	case "float32":
		if data == nil {
			data = float64(0)
		}
		r = float32(data.(float64))
	case "float64":
		if data == nil {
			data = float64(0)
		}
		r = data.(float64)
	case "int":
		if data == nil {
			data = float64(0)
		}
		r = int(math.Trunc(data.(float64)))
	case "int8":
		if data == nil {
			data = float64(0)
		}
		r = int8(math.Trunc(data.(float64)))
	case "int16":
		if data == nil {
			data = float64(0)
		}
		r = int16(math.Trunc(data.(float64)))
	case "int32":
		if data == nil {
			data = float64(0)
		}
		r = int32(math.Trunc(data.(float64)))
	case "int64":
		if data == nil {
			data = float64(0)
		}
		r = int64(math.Trunc(data.(float64)))
	case "uint8":
		if data == nil {
			data = float64(0)
		}
		r = uint8(math.Trunc(data.(float64)))
	case "uint16":
		if data == nil {
			data = float64(0)
		}
		r = uint16(math.Trunc(data.(float64)))
	case "uint32":
		if data == nil {
			data = float64(0)
		}
		r = uint32(math.Trunc(data.(float64)))
	case "uint64":
		if data == nil {
			data = float64(0)
		}
		r = uint64(math.Trunc(data.(float64)))
	case "bool":
		if data == nil {
			data = false
		}
		r = data.(bool)
	case "time":
		if data == nil {
			data = time.Now()
		}
		r, _ = time.ParseInLocation("2006-01-02T15:04:05+08:00", data.(string), time.Local)
	default:
		r = data
	}

	return
}

func IsDefault(i interface{}) bool {
	var b bool

	switch i.(type) {
	case string:
		b = "" == i.(string)
	case int8:
		b = 0 == i.(int8)
	case int16:
		b = 0 == i.(int16)
	case int32:
		b = 0 == i.(int32)
	case int64:
		b = 0 == i.(int64)
	case int:
		b = 0 == i.(int)
	case uint8:
		b = 0 == i.(uint8)
	case uint16:
		b = 0 == i.(uint16)
	case uint32:
		b = 0 == i.(uint32)
	case uint64:
		b = 0 == i.(uint64)
	case float32:
		b = 0 == i.(float32)
	case float64:
		b = 0 == i.(float64)
	case bool:
		b = false == i.(bool)
	case time.Time:
		b = i.(time.Time).IsZero()
	case []interface{}:
		b = nil == i.([]interface{})
	case map[interface{}]interface{}:
		b = nil == i.(map[interface{}]interface{})
	}

	return b
}

func ToString(i interface{}) string {
	var v string

	switch i.(type) {
	case string:
		v = i.(string)
	case int8:
		v = strconv.FormatInt(int64(i.(int8)), 10)
	case int16:
		v = strconv.FormatInt(int64(i.(int16)), 10)
	case int32:
		v = strconv.FormatInt(int64(i.(int32)), 10)
	case int64:
		v = strconv.FormatInt(i.(int64), 10)
	case int:
		v = strconv.Itoa(i.(int))
	case uint8:
		v = strconv.FormatUint(uint64(i.(uint8)), 10)
	case uint16:
		v = strconv.FormatUint(uint64(i.(uint16)), 10)
	case uint32:
		v = strconv.FormatUint(uint64(i.(uint32)), 10)
	case uint64:
		v = strconv.FormatUint(i.(uint64), 10)
	case float32:
		v = strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32)
	case float64:
		v = strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case bool:
		v = strconv.FormatBool(i.(bool))
	case time.Time:
		v = i.(time.Time).Format("2006-01-02T15:04:05+08:00")
	}

	return v
}

func TypeString(i interface{}) string {
	var valType string

	switch i.(type) {
	case string:
		valType = "string"
	case float32:
		valType = "float32"
	case float64:
		valType = "float64"
	case int:
		valType = "int"
	case int8:
		valType = "int8"
	case int16:
		valType = "int16"
	case int32:
		valType = "int32"
	case int64:
		valType = "int64"
	case uint8:
		valType = "byte"
	case uint16:
		valType = "uint16"
	case uint32:
		valType = "uint32"
	case uint64:
		valType = "uint64"
	case bool:
		valType = "bool"
	case time.Time:
		valType = "time.Time"
	}

	return valType
}
