package action

import (
	"github.com/chromedp/chromedp"
)

const (
	perfromance_js = `let times = {};
let t = window.performance.timing;

// 优先使用 navigation v2  https://www.w3.org/TR/navigation-timing-2/
if (typeof window.PerformanceNavigationTiming === 'function') {
  try {
    var nt2Timing = performance.getEntriesByType('navigation')[0]
    if (nt2Timing) {
      t = nt2Timing
    }
  } catch (err) {
  }
}

//重定向时间
times.redirectTime = t.redirectEnd - t.redirectStart;

//dns查询耗时
times.dnsTime = t.domainLookupEnd - t.domainLookupStart;

//TTFB 读取页面第一个字节的时间
times.ttfbTime = t.responseStart - t.navigationStart;

//DNS 缓存时间
times.appcacheTime = t.domainLookupStart - t.fetchStart;

//卸载页面的时间
times.unloadTime = t.unloadEventEnd - t.unloadEventStart;

//tcp连接耗时
times.tcpTime = t.connectEnd - t.connectStart;

//request请求耗时
times.reqTime = t.responseEnd - t.responseStart;

//解析dom树耗时
times.analysisTime = t.domComplete - t.domInteractive;

//白屏时间 
times.blankTime = (t.domInteractive || t.domLoading) - t.fetchStart;

//domReadyTime
times.domReadyTime = t.domContentLoadedEventEnd - t.fetchStart;

times;`
)

type PerformanceResult struct {
	// 重定向时间
	RedirectTime float64 `json:"redirectTime"`

	// dns查询耗时
	DnsTime float64 `json:"dnsTime"`

	// TTFB 读取页面第一个字节的时间
	TtfbTime float64 `json:"ttfbTime"`

	//DNS 缓存时间
	AppcacheTime float64 `json:"appcacheTime"`

	//卸载页面的时间
	UnloadTime float64 `json:"unloadTime"`

	//tcp连接耗时
	TcpTime float64 `json:"tcpTime"`

	//request请求耗时
	ReqTime float64 `json:"reqTime"`

	//解析dom树耗时
	AnalysisTime float64 `json:"analysisTime"`

	//白屏时间
	BlankTime float64 `json:"blankTime"`

	//domReadyTime
	DomReadyTime float64 `json:"domReadyTime"`
}

func Performance(ret *PerformanceResult) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Evaluate(perfromance_js, &ret),
	}
}
