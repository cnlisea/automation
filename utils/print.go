package utils

import (
	"bytes"
	"fmt"
	ghttp "net/http"

	"github.com/cnlisea/automation/http"

	"github.com/cnlisea/automation/constant"
)

func PrintRequestParams(params map[string]interface{}) string {
	var b bytes.Buffer

	b.WriteString("request:\n")

	b.WriteString(fmt.Sprintf("\turl: %s\n", params[constant.RequestHttpUrl]))

	method := params[constant.RequestHttpMethod].(string)
	b.WriteString(fmt.Sprintf("\tmethod: %s\n", method))

	switch method {
	case ghttp.MethodPost:
		fallthrough
	case ghttp.MethodPut:
		fallthrough
	case ghttp.MethodPatch:
		bodyType := params[constant.RequestHttpBodyType].(int)
		switch bodyType {
		case http.HttpBodyTypeJson:
			b.WriteString("\tContent-Type: JSON\n")
		case http.HttpBodyTypeXml:
			b.WriteString("\tbody type: XML\n")
		case http.HttpBodyTypePostForm:
			b.WriteString("\tbody type PostForm\n")
		}
		b.WriteString(fmt.Sprintf("\tbody: %s\n", params[constant.RequestHttpBody]))
	}

	var token string
	if t, ok := params[constant.RequestHttpAuthToken]; ok {
		token = t.(string)
	}

	authType := params[constant.RequestHttpAuthType].(int)
	switch authType {
	case http.BearerTokenAuth:
		b.WriteString(fmt.Sprintf("\tAuthorization: Bearer %s\n", token))
	}

	return b.String()
}

func PrintResponseParams(data map[string]interface{}) string {
	var b bytes.Buffer

	b.WriteString("reqsponse:\n")
	for k, v := range data {
		if k == constant.ResponseHttpStatusCode {
			continue
		}

		switch v.(type) {
		case map[string]interface{}:
			b.WriteString(fmt.Sprintf("\t%s: ", k))
			b.WriteString("{\n")
			for key, val := range v.(map[string]interface{}) {
				b.WriteString(fmt.Sprintf("\t\t%s: %v\n", key, val))
			}
			b.WriteString("\t}\n")
		default:
			b.WriteString(fmt.Sprintf("\t%s: %v\n", k, v))
		}
	}

	return b.String()
}

func PrintLine() string {
	return "================================================================================================================\n"
}
