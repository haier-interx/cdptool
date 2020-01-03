package pipeline

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/haier-interx/cdptool/action"
	"github.com/haier-interx/cdptool/models"
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
	STEP_SCREENSHOTS   = "_screenshots_"
	STEP_PERFORMANCE   = "_performance_"
	STEP_JAVASCRIPT    = "_javascript_"
	STEP_DUMP          = "_dump_"
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

	Screenshots *Screenshots `json:"screenshots"`
	Input       string       `json:"input"`
	JavaScript  string       `json:"javascript"`
	Screen      *Screen      `json:"screen"`
	Language    string       `json:"language"`
	Url         string       `json:"url"`
}

type Screenshots struct {
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

func (s *Step) ActionWithCtx(ctx context.Context, ret *Result, cds *CustomDefinitions) chromedp.Tasks {
	tasks, err := s.Action(ret, cds)
	if err != nil {
		ret.Failed(err)
		return nil
	}

	ts := make([]chromedp.Action, len(tasks))
	for i, task := range tasks {
		ts[i] = action.Wrap(ctx, task)
	}

	return ts
}

func (s *Step) Action(ret *Result, cds *CustomDefinitions) (tasks chromedp.Tasks, err error) {
	defer func() {
		if err != nil {
			ret.Failed(err)
		}
	}()

	tasks = make([]chromedp.Action, 0)

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
			err = ERR_NAVIGATE_URL_REQUIRED
			return
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
			err = ERR_SCREEN_CONFIG_INVALID
			return
		}
		tasks = append(tasks, action.DeviceScreen(s.Screen.Width, s.Screen.Height, s.Screen.Mobile))

	case STEP_SCREENSHOTS: // 截屏
		if s.Screenshots == nil {
			s.Screenshots = &Screenshots{80}
		}
		if s.Screenshots.Quality <= 0 {
			err = ERR_SCREENSOTS_CONFIG_INVALID
			return
		}
		if s.Sel != "" {
			tasks = append(tasks, chromedp.WaitReady(s.Sel, queryOpt))
		}
		filename := fmt.Sprintf("%s.%d.jpg", s.id, time.Now().UnixNano())
		tasks = append(tasks,
			action.FullScreenshot(s.Screenshots.Quality, filename),
			chromedp.ActionFunc(func(ctx context.Context) error {
				ret.ScreenshotsFileName = append(ret.ScreenshotsFileName, filename)
				return nil
			}),
		)

	case STEP_PERFORMANCE: // performance
		pr := new(models.PerformanceTiming)
		ret.Performances = append(ret.Performances, pr)
		tasks = append(tasks, action.Performance(pr))

	case STEP_JAVASCRIPT: //run JavaScript
		jsr := new([]byte)
		ret.JavaScriptResult = append(ret.JavaScriptResult, jsr)
		tasks = append(tasks, chromedp.Evaluate(s.JavaScript, jsr))

	case STEP_DUMP:
		tasks = append(tasks, action.Dump())

	default: // 非内置
		if cds == nil {
			err = ERR_STEPTYPE_INVALID
			return
		}

		// custom defined step type 自定义步骤
		sg, found := cds.Steps(s.Type)
		if !found {
			err = ERR_STEPTYPE_INVALID
			return
		}

		// every step in step group
		for sub_idx, sub_step := range sg.Steps {
			sub_es := NewStepResult(sg.Id, sub_idx, sub_step.Type)

			var sub_tasks chromedp.Tasks
			sub_tasks, err = sub_step.Action(ret, cds)
			if err != nil {
				return
			}

			// pre action
			tasks = append(tasks,
				chromedp.ActionFunc(func(ctx context.Context) error {
					ret.SetStepStarted(sub_es)
					return nil
				}),
			)

			// action
			tasks = append(tasks, sub_tasks...)

			// post action
			// calculate the step duration
			tasks = append(tasks,
				chromedp.ActionFunc(func(ctx context.Context) error {
					ret.SetStepOver(sub_es)
					return nil
				}),
			)
		}
	}

	return
}

func IsBuildInStep(id string) bool {
	switch id {
	case STEP_NAVIGATE, STEP_WAIT, STEP_INPUT, STEP_CLICK, STEP_LANGUAGE, STEP_DEVICE_SCREEN, STEP_SCREENSHOTS, STEP_PERFORMANCE, STEP_JAVASCRIPT, STEP_DUMP:
		return true
	default:
		return false
	}
}
