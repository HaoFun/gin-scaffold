package locales

import "fmt"

type ResCode int

type CErrors struct {
	Code ResCode
}

func (ce *CErrors) Error() string {
	return fmt.Sprintf("CodeError: %d", ce.Code)
}

const (
	CodeSuccess   ResCode = 20000
	CodeUserExist ResCode = 40000 + iota
	CodeMobileExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeNeedLogin
	CodeInvalidToken
	CodeInvalidGenToken
	CodeInvalidLogout

	CodeDBCreateError
)

const (
	CodeInvalidParamBase         = 42200
	CodeInvalidParam     ResCode = CodeInvalidParamBase + iota
)

const (
	CodeServerBase         = 50000
	CodeServerBusy ResCode = CodeServerBase + iota
	CodeGuardNameNotExist
)

var codeMsgMap = map[ResCode]map[string]string{
	CodeSuccess: {
		"en":    "Success",
		"zh_tw": "成功",
	},
	CodeUserExist: {
		"en":    "Username already exists",
		"zh_tw": "用戶名已存在",
	},
	CodeMobileExist: {
		"en":    "Mobile number already exists",
		"zh_tw": "手機號碼已存在",
	},
	CodeUserNotExist: {
		"en":    "Username does not exist",
		"zh_tw": "用戶名不存在",
	},
	CodeInvalidPassword: {
		"en":    "Invalid username or password",
		"zh_tw": "用戶名或密碼錯誤",
	},
	CodeNeedLogin: {
		"en":    "Login required",
		"zh_tw": "需要登入",
	},
	CodeInvalidToken: {
		"en":    "Invalid token",
		"zh_tw": "無效的token",
	},
	CodeInvalidGenToken: {
		"en":    "Error generating token",
		"zh_tw": "生成Token錯誤",
	},
	CodeInvalidLogout: {
		"en":    "Logout failed",
		"zh_tw": "登出失敗",
	},
	CodeDBCreateError: {
		"en":    "Error occurred during DB creation",
		"zh_tw": "DB新增發生錯誤",
	},

	CodeInvalidParam: {
		"en":    "Invalid request parameter",
		"zh_tw": "請求參數錯誤",
	},

	CodeServerBusy: {
		"en":    "Service busy",
		"zh_tw": "服務繁忙",
	},

	CodeGuardNameNotExist: {
		"en":    "GuardName does not exist",
		"zh_tw": "GuardName不存在",
	},
}

func (c ResCode) Msg(lang string) string {
	msgMap, ok := codeMsgMap[c]
	if !ok {
		msgMap = codeMsgMap[CodeServerBusy]
	}

	msg, ok := msgMap[lang]
	if !ok {
		msg = msgMap["en"]
	}

	return msg
}

func IsCErrors(err error, code ResCode) bool {
	cErr, ok := err.(*CErrors)
	if !ok {
		return false
	}
	return cErr.Code == code
}
