package oauth

import (
	"bytes"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
)

// GetLoginMsg 获取登陆失败后页面上的提示语
func GetLoginMsg(resp *resty.Response) string {
	re := regexp.MustCompile(`<span\s+id="msg">(.+?)</span>`)
	matches := re.FindStringSubmatch(resp.String())
	if len(matches) == 0 {
		return ""
	}
	// 删除span内部的标签
	re = regexp.MustCompile(`<[^>]*>`)
	msg := re.ReplaceAllString(matches[1], "")
	return msg
}

// CheckLogin 用于判断登陆是否成功
func CheckLogin(resp *resty.Response) error {
	// 判断登陆是否成功
	destination := resp.RawResponse.Request.URL.String()
	if destination == PersonalCenterURL {
		// 登陆成功后会跳转到用户中心
		return nil
	}

	// 判断是否需要修改密码
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return err
	}
	title := doc.Find("title").Text()
	if title == "修改密码" {
		return oauthException.EditPasswordError
	}

	// 判断失败原因
	msg := GetLoginMsg(resp)
	switch msg {
	case WrongPasswordMsg:
		return oauthException.WrongPassword
	case WrongAccountMsg:
		return oauthException.WrongAccount
	case NotActivatedMsg:
		return oauthException.NotActivatedError
	}
	return oauthException.OtherError
}

// CheckIsClosed 判断统一是否关闭
func CheckIsClosed() error {
	resp, err := http.Get(PersonalCenterURL)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	title := doc.Find("title").Text()
	if title == "Error 403.6" {
		return oauthException.ClosedError
	}
	return nil
}
