import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// 发送Http请求: POST
func DoHttp(ctx *srfs.Context, url, HttpReqBody string) (rspBodyByte []byte, errMsg string, err error) {
	// 准备: Http请求
	reqBody := strings.NewReader(HttpReqBody)
	HttpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		errMsg = fmt.Sprintf("NewRequest fail, url: %s, HttpReqBody: %+v, err: %+v", url, HttpReqBody, err)
		err = errors.New(errMsg)
		ctx.Error(errMsg)
		return
	}
	HttpReq.Header.Add("Content-Type", "application/json")
	// DO: Http请求
	HttpRsp, err := http.DefaultClient.Do(HttpReq)
	if err != nil {
		errMsg = fmt.Sprintf("do Http fail, url: %s, HttpReqBody: %+v, err:%+v", url, HttpReqBody, err)
		err = errors.New(errMsg)
		ctx.Error(errMsg)
		return
	}
	defer HttpRsp.Body.Close()
	// Read: Http结果
	rspBodyByte, err = ioutil.ReadAll(HttpRsp.Body)
	if err != nil {
		errMsg = fmt.Sprintf("ReadAll failed, url: %s, HttpReqBody: %+v, err: %+v", url, HttpReqBody, err)
		err = errors.New(errMsg)
		ctx.Error(errMsg)
		return
	}
	return
}
