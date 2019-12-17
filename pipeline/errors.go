package pipeline

import "errors"

var (
	ERR_STEPTYPE_INVALID            = errors.New("the step type invalid")
	ERR_PIPELINE_ID_REQUIRED        = errors.New("pipeline required")
	ERR_NAVIGATE_URL_REQUIRED       = errors.New("url required")
	ERR_ELEMENT_NOTFOUND            = errors.New("element not found")
	ERR_ELEMENT_NOTFOUND_OR_TIMEOUT = errors.New("element not found or execute timeout")
	ERR_SCREEN_CONFIG_INVALID       = errors.New("screen configure invalid")
	ERR_SCREENSOT_CONFIG_INVALID    = errors.New("screenshot configure invalid")
	ERR_STEPDEFINED_REPEAT          = errors.New("step definitions repeat")
)

func ErrorCN(err error) string {
	if err == nil {
		return ""
	}

	switch err {
	case ERR_PIPELINE_ID_REQUIRED:
		return "缺少Id配置"
	case ERR_ELEMENT_NOTFOUND:
		return "元素未找到"
	case ERR_ELEMENT_NOTFOUND_OR_TIMEOUT:
		return "元素未找到或者任务执行太慢而超时"
	case ERR_SCREEN_CONFIG_INVALID:
		return "设备屏幕像素配置错误"
	case ERR_SCREENSOT_CONFIG_INVALID:
		return "截屏参数配置错误"
	case ERR_STEPTYPE_INVALID:
		return "步骤类型错误"
	default:
		err_parent := errors.Unwrap(err)
		if err_parent == nil {
			return err.Error()
		}

		return ErrorCN(err_parent)
	}
}
