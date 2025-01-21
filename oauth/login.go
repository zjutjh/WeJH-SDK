package oauth

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

// Login 统一登陆
func Login(username, password string) ([]*http.Cookie, error) {
	client := resty.New()
	// 使用cookieJar管理cookie
	cookieJar, _ := cookiejar.New(nil)
	client.SetCookieJar(cookieJar)

	// 1. 初始化请求
	resp, err := client.R().
		Get(LoginUrl)
	if err != nil {
		return nil, err
	}
	// 检查统一系统是否关闭
	if err = checkIsClosed(resp); err != nil {
		return nil, err
	}

	// 2. 登陆参数生成
	// 解析execution
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}
	execution := doc.
		Find("input[type=hidden][name=execution]").
		AttrOr("value", "")
	// 密码加密
	encPwd, err := getEncryptedPassword(client, password)

	loginParams := map[string]string{
		"username":  username,
		"password":  encPwd,
		"execution": execution,
		"_eventId":  "submit",
	}

	// 3. 发送登陆请求
	resp, err = client.R().
		SetFormData(loginParams).
		Post(LoginUrl)
	if err != nil {
		return nil, err
	}
	// 检查登陆信息
	if err = checkLogin(resp); err != nil {
		return nil, err
	}

	// 4. 提取指定域名下的session并构造cookie列表
	u, _ := url.Parse(MeZjutURL)
	return cookieJar.Cookies(u), nil
}

// GetUserInfo 登陆并获取用户信息
func GetUserInfo(username, password string) (cookies []*http.Cookie, userInfo UserInfo, err error) {
	cookies, err = Login(username, password)
	if err != nil {
		return cookies, userInfo, err
	}
	userData := struct {
		Data UserInfo `json:"data"`
	}{}
	_, err = resty.New().R().
		SetCookies(cookies).
		SetResult(&userData).
		Get(UserInfoApi)
	return cookies, userData.Data, err
}
