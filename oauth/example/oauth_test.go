package oauth

import (
	"errors"
	"testing"

	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
)

func TestLogin(t *testing.T) {
	for i, tc := range []struct {
		inputUsername string
		inputPassword string
		expect        error
	}{
		{
			inputUsername: "example",
			inputPassword: "i m password",
			expect:        nil,
		},
		{
			inputUsername: "wrong username",
			inputPassword: "this password is correct",
			expect:        oauthException.WrongAccount,
		},
		{
			inputUsername: "MangoGovo",
			inputPassword: "wrong password",
			expect:        oauthException.WrongPassword,
		},
	} {
		cookies, _, err := oauth.GetUserInfo(tc.inputUsername, tc.inputPassword)
		if err != nil || len(cookies) == 0 {
			t.Errorf("测试点 %d (%s,%s): 遇到了异常: %v", i+1, tc.inputUsername, tc.inputPassword, err)
			continue
		}

		if !errors.Is(err, tc.expect) {
			t.Errorf("测试点 %d (%s,%s): 期望值:=%s 实际值=%s", i+1, tc.inputUsername, tc.inputPassword, tc.expect, err)
		}
	}
}
