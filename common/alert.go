package common

import (
	"strconv"
)

//alert title
const (
	AlertWarning = iota
	AlertError
	AlertFail
)

// alert Content
const (
	AlertParamsError = iota
	AlertSaveFail
	AlertDelFail
	AlertGetDataFail
	AlertLoginFail
	AlertUserError
	AlertPassError
	AlertDBFail
	AlertCheckTokenError
	AlertPathVisitError
	AlertDialogLengthError
	AlertParentIDError
	AlertRepeatError
	AlertVoucherError
	AlertVoucherIIError
	AlertFileEmptyError
	AlertFileFormatError
	AlertFileSizeError
	AlertBatchVoucherFail
	AlertRangeDateError
	AlertRangeDateErrorII
)

var (
	alertTitle   = []string{"警告", "错误", "失败"}
	alertContent = []string{
		"提交参数错误！请检查是否提交了空白参数，或参数的内容及格式是否正确。",
		"提交/保存失败，具体请查看日志.",
		"删除失败，具体请查看日志.",
		"获取数据失败，或无数据记录，具体请查看日志.",
		"登录失败，请检查账户和密码.",
		"用户名只允许4-12位数字+大小写字母+下划线组成",
		"密码必须6位以上",
		"数据库查询失败，请联系管理员查看日志",
		"登录失效，登陆已超时或账户已在其他地方登录",
		"你无权访问本页面，请重新登录",
		"对话长度不能超过十五个字！",
		"不能设置自己为父级",
		"已有该种类，不可重复添加同一类型",
		"您无法编辑非本店铺的兑换券",
		"您没有相关权限,无法新增/批量发放兑换券",
		"文件为空或上传失败",
		"只允许上传jpg,png,gif,jpeg格式",
		"上传图片大小不得超过3M",
		"您选择的地区中没有店铺或者店铺没有绑定用户，批量发送取消",
		"选择日期范围错误，结束时间不得早于开始时间",
		"选择日期范围错误，结束时间不得早于当前时间",
	}
)

// GetAlertMsg 获取消息
func GetAlertMsg(t, c string) (string, string) {
	if t == "" || c == "" {
		return "", ""
	}
	ti, err := strconv.Atoi(t)
	ci, err2 := strconv.Atoi(c)
	if err != nil || err2 != nil {
		return "", ""
	}
	if ti >= len(alertTitle) || ci >= len(alertContent) {
		return "", ""
	}
	return alertTitle[ti], alertContent[ci]
}

// GetAlertCentent 获取消息内容
func GetAlertCentent(c int) string {
	if c >= len(alertContent) {
		return ""
	}
	return alertContent[c]
}
