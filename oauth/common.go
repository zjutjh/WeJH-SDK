package oauth

// 用户中心
const (
	MeZjutURL         = "http://www.me.zjut.edu.cn"
	PersonalCenterURL = MeZjutURL + "/personal-center"
	UserInfoApi       = MeZjutURL + "/api/basic/info"
)

// 统一登陆
const (
	BaseUrl      = "https://oauth.zjut.edu.cn/cas"
	PublicKeyUrl = BaseUrl + "/v2/getPubKey"
	LoginUrl     = BaseUrl + "/login"
)

// 登陆错误对应在页面的提示信息
const (
	WrongPasswordMsg = "用户名或密码错误" // #nosec G101
	WrongAccountMsg  = "当前账号无权登录"
	NotActivatedMsg  = "账号未激活，请激活后再登录"
)

// UserInfo 用户信息
type UserInfo struct {
	College      string `json:"bmmc"`         // 学院
	Name         string `json:"username"`     // 姓名
	StudentID    string `json:"yhm"`          // 学号
	UserType     string `json:"jsdm"`         // 用户类型 [id]
	UserTypeDesc string `json:"jsmc"`         // 用户类型 [description]
	Gender       string `json:"xb"`           // 性别
	Avatar       string `json:"headPortrait"` // 头像链接
}
