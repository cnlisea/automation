package generate

import (
	"bytes"
	"encoding/json"
	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/http"
	"github.com/cnlisea/automation/utils"
	ghttp "net/http"
	gurl "net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func InterfaceFormat(name string, req map[string]interface{}, res map[string]interface{}) (string, error) {
	var (
		title      string
		author     string
		method     string
		url        *gurl.URL
		function   string
		explain    []string
		mReqSet    = make(map[string]interface{}, 0)
		nReqSet    = make(map[string]interface{}, 0)
		reqDescSet = make(map[string]string, 0)
		reqTypeSet = make(map[string]string, 0)
		resDescSet = make(map[string]string, 0)
		resTypeSet = make(map[string]string, 0)
		err        error
	)

	// title
	title = name
	if v, ok := req[constant.Title+constant.DescSuffix]; ok {
		if t, ok := v.(string); ok {
			title = t
		}
	}

	// author
	author = constant.DefaultAuthor
	if v, ok := req[constant.Author+constant.DescSuffix]; ok {
		if a, ok := v.(string); ok {
			author = a
		} else {
			author = constant.DefaultAuthor
		}
	}

	// function
	function = title
	if v, ok := req[constant.Function+constant.DescSuffix]; ok {
		if f, ok := v.(string); ok {
			function = f
		}
	}

	// method
	method = req[constant.RequestHttpMethod].(string)
	// url
	url, err = gurl.Parse(req[constant.RequestHttpUrl].(string))
	if nil != err {
		return "", err
	}

	// explain
	if v, ok := req[constant.Explain+constant.DescSuffix]; ok {
		switch v.(type) {
		case []string:
			explain = v.([]string)
		case string:
			explain = make([]string, 1)
			explain[0] = v.(string)
		case []interface{}:
			for _, v := range v.([]interface{}) {
				explain = append(explain, utils.ToString(v))
			}
		}
	}

	// request param
	for k, v := range req {
		switch k[len(k)-2:] {
		case "_m":
			mReqSet[k[:len(k)-2]] = v
		case "_n":
			nReqSet[k[:len(k)-2]] = v
		case constant.DescSuffix:
			reqDescSet[k[:len(k)-2]] = utils.ToString(v)
		case constant.TypeSuffix:
			reqTypeSet[k[:len(k)-2]] = utils.ToString(v)
		}
	}

	// response param
	for k, v := range res {
		switch k[len(k)-2:] {
		case constant.DescSuffix:
			resDescSet[k] = utils.ToString(v)
		case constant.TypeSuffix:
			resTypeSet[k] = utils.ToString(v)
		default:
		}
	}

	// create buffer
	var b bytes.Buffer

	// write markdown doc
	b.WriteString("## " + title + "\n")
	b.WriteString("* 作    者: " + author + "\n")
	b.WriteString("* 修改时间: " + time.Now().Format("2006年01月02日15点04分") + "\n")
	b.WriteString("\n")

	b.WriteString("#### 方法\n")
	b.WriteString("`" + method + "` `" + url.Path + "`\n")
	b.WriteString("\n")

	b.WriteString("#### 功能\n")
	b.WriteString("* " + function + "\n")
	b.WriteString("\n")

	if 0 < len(explain) {
		b.WriteString("#### 说明\n")
		for _, v := range explain {
			b.WriteString("* " + v + "\n")
		}
		b.WriteString("\n")
	}

	// 请求参数
	b.WriteString("#### 请求参数\n")
	b.WriteString("\n")
	b.WriteString("|字段|类型|必选|说明|\n")
	b.WriteString("|:--|:--|:--|:--|\n")
	// 必选请求字段
	for k, v := range mReqSet {
		desc := k
		if d, ok := reqDescSet[k]; ok {
			desc = d
		}
		b.WriteString("|" + k + "|")
		if t, ok := reqTypeSet[k]; ok {
			b.WriteString(t)
		} else {
			b.WriteString(utils.TypeString(v))
		}

		b.WriteString("|是|" + desc + "|\n")
	}

	// 可选请求字段
	for k, v := range nReqSet {
		desc := k
		if d, ok := reqDescSet[k]; ok {
			desc = d
		}
		b.WriteString("|" + k + "|")
		if t, ok := reqTypeSet[k]; ok {
			b.WriteString(t)
		} else {
			b.WriteString(utils.TypeString(v))
		}
		b.WriteString("|否|" + desc + "|\n")
	}
	b.WriteString("\n")

	// 请求实例
	b.WriteString("#### 请求实例\n")
	b.WriteString("\n")
	b.WriteString("```bash\n")
	b.WriteString("curl ")
	switch req[constant.RequestHttpAuthType].(int) {
	case http.BearerTokenAuth:
		var token string
		if t, ok := req[constant.RequestHttpAuthToken]; ok {
			token = utils.ToString(t)
		}
		b.WriteString("-H \"Authorization: Bearer " + token + "\" ")
	}

	switch req[constant.RequestHttpMethod].(string) {
	case ghttp.MethodPost:
		b.WriteString("-d " + string(req[constant.RequestHttpBody].([]byte)) + " -X POST ")
	case ghttp.MethodPut:
		b.WriteString("-d " + string(req[constant.RequestHttpBody].([]byte)) + " -X PUT ")
	case ghttp.MethodPatch:
		b.WriteString("-d " + string(req[constant.RequestHttpBody].([]byte)) + " -X PATCH ")
	case ghttp.MethodGet:
		b.WriteString("-X GET ")
	case ghttp.MethodHead:
		b.WriteString("-X HEAD ")
	case ghttp.MethodDelete:
		b.WriteString("-X DELETE ")
	case ghttp.MethodOptions:
		b.WriteString("-X OPTIONS")
	}

	b.WriteString(url.String() + "\n")
	b.WriteString("```\n")
	b.WriteString("\n")

	b.WriteString("#### 返回结果\n")
	b.WriteString("```json\n")
	dels := make([]string, 0)
	if _, ok := res[constant.ResponseHttpStatusCode]; ok {
		dels = append(dels, constant.ResponseHttpStatusCode)
	}
	for k, _ := range resDescSet {
		dels = append(dels, k)
	}
	for k, _ := range resTypeSet {
		dels = append(dels, k)
	}

	dataMap := utils.DeleteMap(res, dels)

	resData, err := json.MarshalIndent(&dataMap, "", "\t")
	if nil != err {
		return "", err
	}
	b.Write(resData)
	b.WriteString("\n```\n")

	b.WriteString("#### 返回参数说明\n")
	b.WriteString("\n")
	b.WriteString("|参数名|类型|说明|\n")
	b.WriteString("|:----|:----|:----|\n")
	for k, v := range dataMap {
		b.WriteString("|" + k + "|")

		if d, ok := resTypeSet[k+constant.TypeSuffix]; ok {
			b.WriteString(d)
		} else {
			b.WriteString(utils.TypeString(v))
		}
		b.WriteString("|")

		if t, ok := resDescSet[k+constant.DescSuffix]; ok {
			b.WriteString(t)
		} else {
			b.WriteString(k)
		}
		b.WriteString("|\n")
	}
	b.WriteString("\n")

	b.WriteString("错误代码请见错误代码对照表\n")
	b.WriteString("\n")
	b.WriteString("----\n")

	return b.String(), nil
}

func GenerateDoc(filename string, title string, inters map[string]string, wg *sync.WaitGroup) error {
	if !strings.Contains(filename, ".md") {
		// 添加默认文件后缀名
		filename = filename + ".md"
	}

	// 判断文件是否存在
	if utils.FileExist(filename) {
		filename = strings.Replace(filename, ".md", "", -1) + time.Now().Format("20060102150405") + ".md"
	}

	// 创建并打开文件
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return err
	}

	f.WriteString("<secret>\n")
	f.WriteString("\n")
	f.WriteString("# " + title + "\n")
	f.WriteString("\n")

	if v, ok := inters["info"]; ok && (-1 != strings.Index(v, `data": `)) {
		f.WriteString("## " + title + "结构\n")
		f.WriteString("\n")
		f.WriteString("#### json结构\n")
		f.WriteString("\n")
		f.WriteString("```json\n")
		jsonBegin := strings.Index(v, `data": `) + 7
		jsonEnd := strings.LastIndex(v, "\t}") + 2
		jsonBody := strings.Replace(v[jsonBegin:jsonEnd], "\n\t", "\n", -1)
		f.WriteString(jsonBody + "\n")
		f.WriteString("```\n")
		f.WriteString("\n")

		infoMap := make(map[string]interface{}, 0)
		if err = json.Unmarshal([]byte(jsonBody), &infoMap); nil != err {
			return err
		}

		f.WriteString("#### 字段说明\n")
		f.WriteString("\n")
		f.WriteString("|字段|类型|说明|\n")
		f.WriteString("|:--|:--|:--|\n")
		for k, vv := range infoMap {
			f.WriteString("|" + k + "|")
			vvType := utils.TypeString(vv)
			switch vvType {
			case "string":
				if _, err = time.ParseInLocation("2006-01-02T15:04:05+08:00", utils.ToString(vv), time.Local); nil == err {
					vvType = "time"
				}
			case "float64":
				if _, err := strconv.Atoi(strconv.FormatFloat(vv.(float64), 'f', -1, 64)); nil == err {
					vvType = "int"
				}
			}
			f.WriteString(vvType)
			f.WriteString("|" + k + "|\n")
		}
		f.WriteString("\n")
		f.WriteString("----\n")
		f.WriteString("\n")
	}

	for _, v := range inters {
		f.WriteString(v)
		f.WriteString("\n")
	}

	f.WriteString("</secret>\n")

	f.Close()

	wg.Done()
	return nil
}

/*
func (t JubaoTest) GenerateDoc1(name string, req map[string]interface{}, res map[string]interface{}) error {
	var (
		docname    string
		title      string
		author     string
		method     string
		url        *gurl.URL
		function   string
		explain    []string
		mReqSet    = make(map[string]interface{}, 0)
		nReqSet    = make(map[string]interface{}, 0)
		reqDescSet = make(map[string]string, 0)
		reqTypeSet = make(map[string]string, 0)
		resDescSet = make(map[string]string, 0)
		resTypeSet = make(map[string]string, 0)
		err        error
	)

	// doc file name
	docname = name
	if v, ok := req[constant.Docname+constant.DescSuffix]; ok {
		if n, ok := v.(string); ok {
			docname = n
		} else {
			docname = name
		}
	}

	// title
	title = name
	if v, ok := req[constant.Title+constant.DescSuffix]; ok {
		if t, ok := v.(string); ok {
			title = t
		}
	}

	// author
	author = constant.DefaultAuthor
	if v, ok := req[constant.Author+constant.DescSuffix]; ok {
		if a, ok := v.(string); ok {
			author = a
		} else {
			author = constant.DefaultAuthor
		}
	}

	// function
	function = title
	if v, ok := req[constant.Function+constant.DescSuffix]; ok {
		if f, ok := v.(string); ok {
			function = f
		}
	}

	// method
	method = req[constant.RequestHttpMethod].(string)
	// url
	url, err = gurl.Parse(req[constant.RequestHttpUrl].(string))
	if nil != err {
		return err
	}

	// explain
	if v, ok := req[constant.Explain+constant.DescSuffix]; ok {
		switch v.(type) {
		case []string:
			explain = v.([]string)
		case string:
			explain = make([]string, 1)
			explain[0] = v.(string)
		}
	}

	// request param
	for k, v := range req {
		switch k[len(k)-2:] {
		case "_m":
			mReqSet[k[:len(k)-2]] = v
		case "_n":
			nReqSet[k[:len(k)-2]] = v
		case constant.DescSuffix:
			reqDescSet[k[:len(k)-2]] = utils.ToString(v)
		case constant.TypeSuffix:
			reqTypeSet[k[:len(k)-2]] = utils.ToString(v)
		}
	}

	// response param
	for k, v := range res {
		switch k[len(k)-2:] {
		case constant.DescSuffix:
			resDescSet[k] = utils.ToString(v)
		case constant.TypeSuffix:
			resTypeSet[k] = utils.ToString(v)
		default:
		}
	}

	if !strings.Contains(docname, ".md") {
		// 添加默认文件后缀名
		docname = docname + ".md"
	}

	// 判断文件是否存在
	if utils.FileExist(docname) {
		docname = docname + time.Now().Format("20060102150405")
	}

	// 创建并打开文件
	f, err := os.OpenFile(docname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return err
	}

	// write markdown doc
	f.WriteString("## " + title + "\n")
	f.WriteString("* 作    者: " + author + "\n")
	f.WriteString("* 修改时间: " + time.Now().Format("2006年01月02日15点04分") + "\n")
	f.WriteString("\n")

	f.WriteString("#### 方法\n")
	f.WriteString("`" + method + "` `" + url.Path + "`\n")
	f.WriteString("\n")

	f.WriteString("#### 功能\n")
	f.WriteString("* " + function + "\n")
	f.WriteString("\n")

	if 0 < len(explain) {
		f.WriteString("#### 说明\n")
		for _, v := range explain {
			f.WriteString("* " + v + "\n")
		}
		f.WriteString("\n")
	}

	// 请求参数
	f.WriteString("#### 请求参数\n")
	f.WriteString("\n")
	f.WriteString("|字段|类型|必选|说明|\n")
	f.WriteString("|:--|:--|:--|:--|\n")
	// 必选请求字段
	for k, v := range mReqSet {
		desc := k
		if d, ok := reqDescSet[k]; ok {
			desc = d
		}
		f.WriteString("|" + k + "|" + utils.TypeString(v) + "|是|" + desc + "|\n")
	}

	// 可选请求字段
	for k, v := range nReqSet {
		desc := k
		if d, ok := reqDescSet[k]; ok {
			desc = d
		}
		f.WriteString("|" + k + "|" + utils.TypeString(v) + "|否|" + desc + "|\n")
	}
	f.WriteString("\n")

	// 请求实例
	f.WriteString("#### 请求实例\n")
	f.WriteString("\n")
	f.WriteString("```bash\n")
	f.WriteString("curl ")
	switch req[constant.RequestHttpAuthType].(int) {
	case http.BearerTokenAuth:
		f.WriteString("-H \"Authorization: Bearer " + t.Token + "\" ")
	}

	switch req[constant.RequestHttpMethod].(string) {
	case ghttp.MethodPost:
		f.WriteString("-d " + req[constant.RequestHttpBody].(string) + " -X POST ")
	case ghttp.MethodPut:
		f.WriteString("-d " + req[constant.RequestHttpBody].(string) + " -X PUT ")
	case ghttp.MethodPatch:
		f.WriteString("-d " + req[constant.RequestHttpBody].(string) + " -X PATCH ")
	}

	f.WriteString(url.String() + "\n")
	f.WriteString("```\n")

	f.WriteString("#### 返回结果\n")
	f.WriteString("```json\n")
	dels := make([]string, 0)
	if _, ok := res[constant.ResponseHttpStatusCode]; ok {
		dels = append(dels, constant.ResponseHttpStatusCode)
	}
	for k, _ := range resDescSet {
		dels = append(dels, k)
	}
	for k, _ := range resTypeSet {
		dels = append(dels, k)
	}

	dataMap := utils.DeleteMap(res, dels)

	resData, _ := json.MarshalIndent(&dataMap, "", "\t")
	f.Write(resData)
	f.WriteString("\n```\n")

	f.WriteString("#### 返回参数说明\n")
	f.WriteString("\n")
	f.WriteString("|参数名|类型|说明|\n")
	f.WriteString("|:----|:----|:----|\n")
	for k, v := range dataMap {
		f.WriteString("|" + k + "|")

		if d, ok := resTypeSet[k+constant.TypeSuffix]; ok {
			f.WriteString(d)
		} else {
			f.WriteString(utils.ToString(v))
		}
		f.WriteString("|")

		if t, ok := resDescSet[k+constant.DescSuffix]; ok {
			f.WriteString(t)
		} else {
			f.WriteString(k)
		}
		f.WriteString("|\n")
	}
	f.WriteString("\n")

	f.WriteString("错误代码请见错误代码对照表\n")
	f.WriteString("\n")
	f.WriteString("----\n")
	return f.Close()
}
*/
