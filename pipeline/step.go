package pipeline

import (
	"cdptool/action"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"time"
)

type StepType string

const (
	STEP_NAVIGATE      = "navigate"
	STEP_SHOW          = "show"
	STEP_INPUT         = "input"
	STEP_CLICK         = "click"
	STEP_LANGUAGE      = "language"
	STEP_DEVICE_SCREEN = "deviceScreen"
	STEP_SCREENSHOT    = "screenshot"
	STEP_PERFORMANCE   = "performance"
	STEP_JAVASCRIPT    = "javascript"
	STEP_DUMP          = "dump"
)

var (
	defaultDeviceScreen = &Screen{
		1600,
		900,
		false,
	}
)

type Step struct {
	Type    string `json:"type" yaml:"type"`
	Sel     string `json:"sel" yaml:"sel"`
	NodeIdx int    `json:"node_idx" yaml:"node_idx"`

	Screenshot *Screenshot `json:"screenshot"`
	Input      string      `json:"input"`
	JavaScript string      `json:"javascript"`
	Screen     *Screen     `json:"screen"`
	Language   string      `json:"language"`
	Url        string      `json:"url"`

	performance *action.PerformanceResult
}

type Screenshot struct {
	Quality int64 `json:"quality"` //compression quality from range [0..100] (jpeg only)
}

type Screen struct {
	Width  int64 `json:"width"`  // Overriding width value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Height int64 `json:"height"` // Overriding height value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Mobile bool  `json:"mobile"` // Whether to emulate mobile device. This includes viewport meta tag, overlay scrollbars, text autosizing and more.
}

func (s *Step) ActionWithCtx(ctx context.Context, id string, data interface{}) (chromedp.Tasks, error) {
	tasks, err := s.Action(id, data)
	if err != nil {
		return nil, err
	}

	ret := make([]chromedp.Action, len(tasks))
	for i, task := range tasks {
		ret[i] = action.Wrap(ctx, task)
	}

	return ret, nil
}

func nodeIdSel(nodeIDs []cdp.NodeID, idx int) interface{} {
	if len(nodeIDs) < idx+1 {
		return ""
	}
	return nodeIDs[idx]
}

func (s *Step) Action(id string, data interface{}) (chromedp.Tasks, error) {
	tasks := make([]chromedp.Action, 0)

	queryOpt := chromedp.ByQuery
	if s.Sel != "" {
		if s.NodeIdx > 0 {
			s.Sel = fmt.Sprintf(`document.querySelectorAll('%s')[%d]`, s.Sel, s.NodeIdx)
			queryOpt = chromedp.ByJSPath
		}
	}

	switch s.Type {
	case STEP_NAVIGATE:
		if s.Url == "" {
			return tasks, ERR_NAVIGATE_URL_REQUIRED
		}
		if s.Sel != "" {
			tasks = append(tasks, chromedp.WaitReady(s.Sel, queryOpt))
		}
		tasks = append(tasks, chromedp.Navigate(s.Url))
	case STEP_SHOW:
		tasks = append(tasks, chromedp.WaitReady(s.Sel, queryOpt))
	case STEP_INPUT:
		tasks = append(tasks, chromedp.SendKeys(s.Sel, s.Input, queryOpt))
	case STEP_CLICK:
		tasks = append(tasks, chromedp.Click(s.Sel, queryOpt))
	case STEP_LANGUAGE: // 语言，比如中文
		tasks = append(tasks, action.UserAgent("", s.Language))
	case STEP_DEVICE_SCREEN: // 设备屏幕分辨率
		if s.Screen == nil {
			s.Screen = defaultDeviceScreen
		}
		if s.Screen.Width <= 0 || s.Screen.Height <= 0 {
			return tasks, ERR_SCREEN_CONFIG_INVALID
		}
		tasks = append(tasks, action.DeviceScreen(s.Screen.Width, s.Screen.Height, s.Screen.Mobile))
	case STEP_SCREENSHOT: // 截屏
		if s.Screenshot == nil {
			s.Screenshot = &Screenshot{80}
		}
		if s.Screenshot.Quality <= 0 {
			return tasks, ERR_SCREENSOT_CONFIG_INVALID
		}
		if s.Sel != "" {
			tasks = append(tasks, chromedp.WaitReady(s.Sel, queryOpt))
		}
		filename := fmt.Sprintf("%s.%d.png", id, time.Now().UnixNano())
		tasks = append(tasks, action.FullScreenshot(s.Screenshot.Quality, filename))

	case STEP_PERFORMANCE:
		pr := data.(*action.PerformanceResult)
		tasks = append(tasks, action.Performance(pr))
	case STEP_JAVASCRIPT:
		var no_receive bool
		if data == nil {
			data = make([]byte, 0)
			no_receive = true
		}
		tasks = append(tasks, chromedp.Evaluate(s.JavaScript, data))
		if no_receive {
			tasks = append(tasks, action.Debug(data))
		}
	case STEP_DUMP:
		tasks = append(tasks, action.Dump())
	default:
		return nil, ERR_STEPTYPE_INVALID
	}

	return tasks, nil
}

func (s *Step) Performance() *action.PerformanceResult {
	return s.performance
}
