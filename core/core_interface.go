package core

import (
	"github.com/cnlisea/automation/http"
	"encoding/json"
	gxml "encoding/xml"
	"errors"
	"fmt"
	ghttp "net/http"
	"net/url"
	"strings"

	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/utils"
)

func (i Instance) InterfaceTest(params map[string]interface{}) (map[string]interface{}, map[string]interface{}, error) {
	uri, ok := params[constant.ConfigRequestUrl].(string)
	if !ok {
		return nil, nil, errors.New("bad url parameter ")
	}

	method, ok := params[constant.ConfigRequestMethod].(string)
	if !ok {
		return nil, nil, errors.New("bad method parameter ")
	}

	// request params
	reqParams, ok := params[constant.ConfigRequest].(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("bad request parameter ")
	}

	// response params
	resParams, ok := params[constant.ConfigResponse].(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("bad response parameter ")
	}

	// auth
	authType := http.NoAuth
	if t, ok := reqParams[constant.Auth]; ok && !utils.IsDefault(t) {
		authType = i.AuthType
	}

	var (
		mQuerys = make(map[string]interface{}, 0)
		nQuerys = make(map[string]interface{}, 0)
		mJsons  = make(map[string]interface{}, 0)
		nJsons  = make(map[string]interface{}, 0)
		mXmls   = make(utils.XmlInterface, 0)
		nXmls   = make(utils.XmlInterface, 0)
		types   = make(map[string]string, 0)
		descs   = make(map[string]interface{}, 0)
	)
	for k, v := range reqParams {
		switch k {
		case constant.Author:
			fallthrough
		case constant.Title:
			fallthrough
		case constant.Function:
			fallthrough
		case constant.Explain:
			descs[k] = v
			continue
		}

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

		// m/n
		style := ks[len(ks)-1]
		switch style {
		case "d": // is description?
			descs[key] = v
			continue
		case "t": // is type ?
			types[key] = utils.ToString(v)
			continue
		}

		switch strings.ToLower(ks[0]) {
		case "json":
			switch style {
			case "m":
				mJsons[key] = v
			case "n":
				nJsons[key] = v
			}
		case "query":
			switch style {
			case "m":
				mQuerys[key] = v
			case "n":
				nQuerys[key] = v
			}
		case "xml":
			switch style {
			case "m":
				mXmls[key] = v
			case "n":
				nXmls[key] = v
			}
		default:
			fmt.Println("Invalid parameter", k)
			continue
		}
	}

	// url query join
	if len(mQuerys) > 0 {
		queryArr := make([]string, 0, len(mQuerys))
		for k, v := range mQuerys {
			if k == "" || utils.ToString(v) == "" {
				continue
			}
			queryArr = append(queryArr, url.QueryEscape(k)+"="+url.QueryEscape(utils.ToString(v)))
		}
		uri = uri + "?" + strings.Join(queryArr, "&")
	}

	// request body
	var (
		body     []byte
		bodyType int
	)
	if len(mJsons) > 0 {
		data, err := json.MarshalIndent(mJsons, "\t", "\t")
		if nil != err {
			return nil, nil, err
		}

		body, bodyType = data, http.HttpBodyTypeJson
	} else if len(mXmls) > 0 {
		data, err := gxml.MarshalIndent(mXmls, "\t", "\t")
		if nil != err {
			return nil, nil, err
		}

		body, bodyType = data, http.HttpBodyTypeXml
	}

	status, resBody, err := http.HttpRequest(uri, body, bodyType, method, authType, i.Token)
	if nil != err {
		return nil, nil, err
	}

	// request url、method、body、bodyType
	reqData := map[string]interface{}{
		constant.RequestHttpUrl:       uri,
		constant.RequestHttpMethod:    method,
		constant.RequestHttpBody:      body,
		constant.RequestHttpBodyType:  bodyType,
		constant.RequestHttpAuthType:  authType,
		constant.RequestHttpAuthToken: i.Token,
	}

	// 必选的参数
	for k, v := range utils.MapMerge(mQuerys, mJsons, mXmls) {
		reqData[k+"_m"] = v
	}

	// 可选的参数
	for k, v := range utils.MapMerge(nQuerys, nJsons, nXmls) {
		reqData[k+"_n"] = v
	}

	// 描述
	for k, v := range descs {
		reqData[k+constant.DescSuffix] = v
	}

	// 参数类型
	for k, v := range types {
		reqData[k+constant.TypeSuffix] = v
	}

	resData := make(map[string]interface{}, 1)
	if status == ghttp.StatusOK {
		resData, err = AnalysisResponse(resBody, resParams)
		if nil != err {
			return nil, nil, err
		}
	}

	// response status code
	resData[constant.ResponseHttpStatusCode] = status

	return reqData, resData, nil
}
