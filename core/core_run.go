package core

import (
	"errors"
	"fmt"
	"github.com/cnlisea/automation/constant"
	"github.com/cnlisea/automation/generate"
	"github.com/cnlisea/automation/http"
	ghttp "net/http"
	"sync"
)

func (i *Instance) Run() error {
	if len(i.Runs) == 0 {
		return errors.New("parameter not resolved, please running before Parse Config")
	}

	// auth
	if err := i.auth(); nil != err {
		return errors.New("interface auth fail, err: " + err.Error())
	}

	wg := new(sync.WaitGroup)
	for menuName, menuSet := range i.Runs {
		docs := make(map[string]string, len(menuSet))
		for k, v := range menuSet {
			req, res, err := i.InterfaceTest(v)
			if nil != err {
				fmt.Printf("%s.%s api test fail, err: %s\n", menuName, k, err)
				continue
			}

			resHttpStatus, ok := res[constant.ResponseHttpStatusCode].(int)
			if !ok {
				fmt.Printf("%s.%s api test error, not find http status code", menuName, k)
				continue
			}

			if resHttpStatus != ghttp.StatusOK {
				fmt.Printf("%s.%s api test fail, status: %d", menuName, k, resHttpStatus)
				continue
			}

			// request and response data format
			doc, err := generate.InterfaceFormat(menuName+"_"+k, req, res)
			if nil != err {
				fmt.Println("%s.%s api doc format fail, error: %s", menuName, k, err)
				continue
			}
			docs[k] = doc

			//fmt.Printf("%s.%s api test success, res: %v\n", menuName, k, res)
			fmt.Printf("%s.%s api test success\n", menuName, k)
		}

		if 0 < len(docs) {
			if "user" == menuName && http.NoAuth != i.AuthType {
				continue
			}
			// generate markdown doc
			wg.Add(1)
			go generate.GenerateDoc(menuName, menuName, docs, wg)
		}
	}

	wg.Wait()
	return nil
}

/*
func (t JubaoTest) Run2() error {
	if len(t.Runs) == 0 {
		return errors.New("parameter not resolved, please running before Parse Config")
	}
	for menuName, menuSet := range t.Runs {
		for k, v := range menuSet {
			if k == "login" {
				continue
			}
			fmt.Println("================================================================================================================")
			req, res, err := t.runOne(v, AuthType)
			if nil != err {
				fmt.Printf("%s.%s api test fail, err: %s\n", menuName, k, err)
				continue
			}
			fmt.Println("req:", req)
			fmt.Println("res:", res)
			fmt.Printf(t.PrintRequestParams(req))

			resHttpStatus, ok := res[constant.ResponseHttpStatusCode].(int)
			if !ok {
				return errors.New(fmt.Sprintf("%s.%s api test error, not find http status code", menuName, k))
			}

			if resHttpStatus != http.StatusOK {
				return errors.New(fmt.Sprintf("%s.%s api test fail, status: %d", menuName, k, resHttpStatus))
			}

			t.GenerateDoc1(menuName+"_"+k, req, res)

			//fmt.Printf("%s.%s api test success, res: %v\n", menuName, k, res)
			fmt.Printf("%s.%s api test success\n", menuName, k)
			fmt.Printf(t.PrintResponseParams(res))
			fmt.Println("================================================================================================================")
		}
	}
	return nil
}
*/
