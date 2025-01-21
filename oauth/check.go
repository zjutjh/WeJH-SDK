package oauth

import (
	"bytes"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
)

// getLoginMsg 获取登陆失败后页面上的提示语
func getLoginMsg(resp *resty.Response) string {
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

// checkLogin 用于判断登陆是否成功
func checkLogin(resp *resty.Response) error {
	// 判断登陆是否成功
	destination := resp.RawResponse.Request.URL.String()
	if destination == PersonalCenterURL {
		// 登陆成功后会跳转到用户中心
		return nil
	}

	// 判断失败原因
	msg := getLoginMsg(resp)
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

// checkIsClosed 判断统一是否关闭
func checkIsClosed(resp *resty.Response) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return err
	}

	title := doc.Find("title").Text()
	if title == "Error 403.6" {
		return oauthException.ClosedError
	}
	return nil
}
