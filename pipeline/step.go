package pipeline

import (
	"cdptool/action"
	"cdptool/models"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

type StepType string

const (
	STEP_NAVIGATE      = "_navigate_"
	STEP_WAIT          = "_wait_"
	STEP_INPUT         = "_input_"
	STEP_CLICK         = "_click_"
	STEP_LANGUAGE      = "_language_"
	STEP_DEVICE_SCREEN = "_deviceScreen_"
	STEP_SCREENSHOT    = "_screenshot_"
	STEP_PERFORMANCE   = "_performance_"
	STEP_JAVASCRIPT    = "_javascript_"
	STEP_DUMP          = "_dump_"
	STEP_NETWORK       = "_network_"
)

var (
	defaultDeviceScreen = &Screen{
		1600,
		900,
		false,
	}
)

type Step struct {
	id string

	Type    string `json:"type" yaml:"type"`
	Sel     string `json:"sel" yaml:"sel"`
	NodeIdx int    `json:"node_idx" yaml:"node_idx"`

	Screenshot *Screenshot `json:"screenshot"`
	Input      string      `json:"input"`
	JavaScript string      `json:"javascript"`
	Screen     *Screen     `json:"screen"`
	Language   string      `json:"language"`
	Url        string      `json:"url"`
}

type Screenshot struct {
	Quality int64 `json:"quality"` //compression quality from range [0..100] (jpeg only)
}

type Screen struct {
	Width  int64 `json:"width"`  // Overriding width value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Height int64 `json:"height"` // Overriding height value in pixels (minimum 0, maximum 10000000). 0 disables the override.
	Mobile bool  `json:"mobile"` // Whether to emulate mobile device. This includes viewport meta tag, overlay scrollbars, text autosizing and more.
}

func (s *Step) SetId(id string) {
	s.id = id
}

func (s *Step) Id() string {
	return s.id
}

func (s *Step) ActionWithCtx(ctx context.Context, ret *Result, cds *CustomDefinitions) (chromedp.Tasks, error) {
	tasks, err := s.Action(ret, cds)
	if err != nil {
		return nil, err
	}

	ts := make([]chromedp.Action, len(tasks))
	for i, task := range tasks {
		ts[i] = action.Wrap(ctx, task)
	}

	return ts, nil
}

func (s *Step) Action(ret *Result, cds *CustomDefinitions) (chromedp.Tasks, error) {
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

	case STEP_WAIT:
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
		filename := fmt.Sprintf("%s.%d.png", s.id, time.Now().UnixNano())
		tasks = append(tasks, action.FullScreenshot(s.Screenshot.Quality, filename))

	case STEP_PERFORMANCE: // performance
		pr := new(models.PerformanceTiming)
		ret.Performances = append(ret.Performances, pr)
		tasks = append(tasks, action.Performance(pr))

	case STEP_JAVASCRIPT: //run JavaScript
		jsr := new([]byte)
		ret.JavaScriptResult = append(ret.JavaScriptResult, jsr)
		tasks = append(tasks, chromedp.Evaluate(s.JavaScript, jsr))

	case STEP_NETWORK:

	case STEP_DUMP:
		tasks = append(tasks, action.Dump())

	default: // 非内置
		if cds == nil {
			return nil, ERR_STEPTYPE_INVALID
		}

		// custom defined step type 自定义步骤
		sg, found := cds.Steps(s.Type)
		if !found {
			return nil, ERR_STEPTYPE_INVALID
		}

		// every step in step group
		for _, step_tmp := range sg.Steps {
			tasks_tmp, err := step_tmp.Action(ret, cds)
			if err != nil {
				return nil, err
			}

			tasks = append(tasks, tasks_tmp...)
		}
	}

	return tasks, nil
}
