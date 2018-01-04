package core

import (
	"encoding/json"
	gxml "encoding/xml"
	"fmt"
	"strings"

	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/utils"
)

func AnalysisResponse(data []byte, resStruct map[string]interface{}) (map[string]interface{}, error) {
	/*if len(data) == 0 {
		return nil, errors.New("bad data parameter")
	}*/

	if len(data) == 0 || len(resStruct) == 0 {
		return map[string]interface{}{}, nil
	}

	var (
		//PostFormSet map[string]interface{}
		jsonSet = make(map[string]interface{}, 0)
		xmlSet  = make(map[string]interface{}, 0)
		tSet    = make(map[string]interface{}, 0)
		dSet    = make(map[string]interface{}, 0)
		mnSet   = make([]string, 0)
		mSet    = make([]string, 0)
		nSet    = make([]string, 0)
	)

	for k, v := range resStruct {
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

		// is param type and describe?
		switch ks[len(ks)-1] {
		case "t":
			tSet[key] = v
			continue
		case "d":
			dSet[key] = v
			continue
		}

		// type
		switch ks[0] {
		case "json":
			jsonSet[key] = v
		case "xml":
			xmlSet[key] = v
		default:
			fmt.Println("Invalid parameter", k)
			continue
		}

		// m/n
		switch ks[len(ks)-1] {
		case "mn":
			mnSet = append(mnSet, key)
		case "m":
			mSet = append(mSet, key)
		case "n":
			nSet = append(nSet, key)

		}

	}

	var result map[string]interface{}
	if len(jsonSet) > 0 {
		var temp map[string]interface{}
		if err := json.Unmarshal(data, &temp); nil != err {
			return nil, err
		}

		// check call success or failure
		status := true
		for _, v := range mnSet {
			switch jsonSet[v].(type) {
			case []interface{}:
				continue
			case map[string]interface{}:
				continue
			}

			if utils.ToString(temp[v]) != utils.ToString(jsonSet[v]) {
				status = false
				break
			}
		}

		// success
		if status {
			jsonSet = utils.DeleteMap(jsonSet, nSet)
		} else { // failure
			jsonSet = utils.DeleteMap(jsonSet, mSet)
		}

		result = make(map[string]interface{}, len(jsonSet))
		for k, v := range jsonSet {
			result[k] = utils.ToType(v, temp[k])
		}

	} else if len(xmlSet) > 0 {
		var temp map[string]interface{}
		if err := gxml.Unmarshal(data, &temp); nil != err {
			return nil, err
		}

		// check call success or failure
		status := true
		for _, v := range mnSet {
			switch xmlSet[v].(type) {
			case []interface{}:
				continue
			case map[string]interface{}:
				continue
			}

			if utils.ToString(temp[v]) != utils.ToString(xmlSet[v]) {
				status = false
				break
			}
		}

		// success
		if status {
			jsonSet = utils.DeleteMap(xmlSet, nSet)
		} else { // failure
			jsonSet = utils.DeleteMap(xmlSet, mSet)
		}

		result = make(map[string]interface{}, len(xmlSet))
		for k, v := range xmlSet {
			result[k] = utils.ToType(v, temp[k])
		}
	} else {
		return map[string]interface{}{}, nil
	}

	// 参数描述
	for k, v := range dSet {
		if _, ok := result[k]; ok {
			result[k+constant.DescSuffix] = v
		}
	}

	// 参数类型
	for k, v := range tSet {
		if _, ok := result[k]; ok {
			result[k+constant.TypeSuffix] = v
		}
	}

	return result, nil
}
