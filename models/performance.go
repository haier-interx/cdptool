package models

// PerformanceTiming (https://w3c.github.io/navigation-timing/)
// PerformanceResourceTiming (https://w3c.github.io/resource-timing/#sec-resource-timing)
type PerformanceTiming struct {
	ConnectEnd                 float64     `json:"connectEnd"`
	ConnectStart               float64     `json:"connectStart"`
	DecodedBodySize            int64       `json:"decodedBodySize"`
	DomComplete                float64     `json:"domComplete"`
	DomContentLoadedEventEnd   float64     `json:"domContentLoadedEventEnd"`
	DomContentLoadedEventStart float64     `json:"domContentLoadedEventStart"`
	DomInteractive             float64     `json:"domInteractive"`
	DomainLookupEnd            float64     `json:"domainLookupEnd"`
	DomainLookupStart          float64     `json:"domainLookupStart"`
	Duration                   float64     `json:"duration"`
	EncodedBodySize            int64       `json:"encodedBodySize"`
	EntryType                  string      `json:"entryType"`
	FetchStart                 float64     `json:"fetchStart"`
	InitiatorType              string      `json:"initiatorType"`
	LoadEventEnd               float64     `json:"loadEventEnd"`
	LoadEventStart             float64     `json:"loadEventStart"`
	Name                       string      `json:"name"`
	NextHopProtocol            string      `json:"nextHopProtocol"`
	RedirectCount              int64       `json:"redirectCount"`
	RedirectEnd                float64     `json:"redirectEnd"`
	RedirectStart              float64     `json:"redirectStart"`
	RequestStart               float64     `json:"requestStart"`
	ResponseEnd                float64     `json:"responseEnd"`
	ResponseStart              float64     `json:"responseStart"`
	SecureConnectionStart      float64     `json:"secureConnectionStart"`
	ServerTiming               interface{} `json:"serverTiming"`
	StartTime                  float64     `json:"startTime"`
	TransferSize               int64       `json:"transferSize"`
	ShotType                   string      `json:"type"`
	UnloadEventEnd             float64     `json:"unloadEventEnd"`
	UnloadEventStart           float64     `json:"unloadEventStart"`
	WorkerStart                int         `json:"workerStart"`
}

// metrics: (https://developer.mozilla.org/en-US/docs/Web/API/Resource_Timing_API/Using_the_Resource_Timing_API)
func (pt *PerformanceTiming) Metrics() *PerformanceMetric {
	pm := new(PerformanceMetric)

	//重定向时间
	pm.RedirectTime = pt.RedirectEnd - pt.RedirectStart

	//dns查询耗时
	pm.DnsTime = pt.DomainLookupEnd - pt.DomainLookupStart

	//TTFB 读取页面第一个字节的时间
	//pm.TtfbTime = pt.ResponseStart - pt.NavigationStart
	pm.TtfbTime = pt.ResponseStart - pt.FetchStart

	//DNS 缓存时间
	pm.AppcacheTime = pt.DomainLookupStart - pt.FetchStart

	//卸载页面的时间
	pm.UnloadTime = pt.UnloadEventEnd - pt.UnloadEventStart

	//tcp连接耗时
	pm.TcpTime = pt.ConnectEnd - pt.ConnectStart

	//request请求耗时
	pm.ReqTime = pt.ResponseEnd - pt.ResponseStart

	//解析dom树耗时
	pm.AnalysisTime = pt.DomComplete - pt.DomInteractive

	//白屏时间
	//if pt.DomInteractive > 0 {
	pm.BlankTime = pt.DomInteractive - pt.FetchStart
	//} else {
	//	pm.BlankTime = pt.DomLoading - pt.FetchStart
	//}

	//domReadyTime
	pm.DomReadyTime = pt.DomContentLoadedEventEnd - pt.FetchStart

	return pm
}

type PerformanceMetric struct {
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
