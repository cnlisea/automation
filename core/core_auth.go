package core

import (
	"errors"
	"fmt"
	"github.com/cnlisea/automation/config"
	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/http"
	"github.com/cnlisea/automation/utils"
	ghttp "net/http"
	"strings"
)

func (i *Instance) auth() error {
	// Judge auth
	switch i.AuthType {
	case http.BearerTokenAuth:
		// on config setting token

		if 0 < len(config.Token) {
			i.Token = config.Token
			break
		}

		if _, ok := i.Runs["user"]["login"]; !ok {
			return errors.New("user.login auth api not find for bearer auth")
		}
		req, res, err := i.InterfaceTest(i.Runs["user"]["login"])
		if nil != err {
			return errors.New(fmt.Sprintln("user.login api test runOne error: ", err))
		}
		fmt.Print(utils.PrintLine())
		fmt.Print(utils.PrintRequestParams(req))

		resHttpStatus, ok := res[constant.ResponseHttpStatusCode].(int)
		if !ok {
			return errors.New("user.login api test error, not find http status code")
		}

		if resHttpStatus != ghttp.StatusOK {
			return errors.New(fmt.Sprintln("user.login api test fail, status: ", resHttpStatus))
		}

		fmt.Print(utils.PrintResponseParams(res))
		fmt.Print(utils.PrintLine())

		// get token
		tks := strings.Split(config.TokenContain, "_")
		var temp interface{} = res
		for _, v := range tks {
			if val, ok := temp.(map[string]interface{}); ok {
				temp = val[v]
			}
		}

		if val, ok := temp.(map[string]interface{}); ok {
			for k, v := range val {
				if strings.Contains(k, config.TokenKey) {
					if token, ok := v.(string); ok {
						i.Token = token
					}
				}
			}
		}

		if len(i.Token) == 0 {
			return errors.New("user.login api get token fail")
		}

		fmt.Println(utils.PrintLine())
		fmt.Println("login success, Authorization Bearer token: ", i.Token)
		fmt.Println(utils.PrintLine())
	case http.NoAuth:
		fallthrough
	default:
		i.Token = ""
	}

	return nil
}
