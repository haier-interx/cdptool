package models

import (
	"fmt"
	"github.com/chromedp/cdproto/network"
	"time"
)

type NetworkLogs struct {
	Items      []*NetworkLog
	requestIds map[string]*NetworkLog
}

func NewNetworkLogs() *NetworkLogs {
	return &NetworkLogs{
		Items:      make([]*NetworkLog, 0),
		requestIds: make(map[string]*NetworkLog),
	}
}

func (nls *NetworkLogs) PutEvent(ev interface{}) error {
	var reqid string
	switch v := ev.(type) {
	case *network.EventRequestWillBeSentExtraInfo:
		reqid = v.RequestID.String()

	case *network.EventRequestWillBeSent:
		reqid = v.RequestID.String()

	case *network.EventResourceChangedPriority:
		reqid = v.RequestID.String()

	case *network.EventResponseReceivedExtraInfo:
		reqid = v.RequestID.String()

	case *network.EventResponseReceived:
		reqid = v.RequestID.String()

	case *network.EventDataReceived:
		reqid = v.RequestID.String()

	case *network.EventLoadingFinished:
		reqid = v.RequestID.String()

	case *network.EventLoadingFailed:
		reqid = v.RequestID.String()

	case *network.EventRequestServedFromCache:
		reqid = v.RequestID.String()

	default:
		return fmt.Errorf("unknown event type: %T", ev)
	}

	nl, found := nls.requestIds[reqid]
	if !found {
		nl = &NetworkLog{ReqID: reqid}
		nls.requestIds[reqid] = nl
		nls.Items = append(nls.Items, nl)
	}

	return nl.PutEvent(ev)
}

type NetworkLog struct {
	ReqID           string
	DocURL          string
	Name            string
	Method          string
	Status          int64
	Type            string
	Initiator       string
	Size            int64
	Time            int64
	Events          NetworkEvents
	ServedFromCache bool
}

type NetworkEvents struct {
	*network.EventRequestWillBeSentExtraInfo
	*network.EventRequestWillBeSent
	*network.EventResourceChangedPriority
	*network.EventResponseReceivedExtraInfo
	*network.EventResponseReceived
	*network.EventDataReceived
	*network.EventLoadingFinished
	*network.EventLoadingFailed
}

type NetworkMetrics struct {
	Dns             float64 `json:"dns"`
	Connection      float64 `json:"connection"`
	Proxy           float64 `json:"proxy"`
	SSL             float64 `json:"ssl"`
	RequestSent     float64 `json:"requestSent"`
	TTFB            float64 `json:"ttfb"`
	ContentDownload float64 `json:"contentDownload"`
}

func (nl *NetworkLog) PutEvent(ev interface{}) error {
	switch v := ev.(type) {
	case *network.EventRequestWillBeSentExtraInfo:
		nl.Events.EventRequestWillBeSentExtraInfo = v

	case *network.EventRequestWillBeSent:
		nl.Events.EventRequestWillBeSent = v
		nl.DocURL = v.DocumentURL
		nl.Name = v.Request.URL
		nl.Method = v.Request.Method
		nl.Type = v.Type.String()
		nl.Initiator = v.Initiator.Type.String()

	case *network.EventResourceChangedPriority:
		nl.Events.EventResourceChangedPriority = v

	case *network.EventResponseReceivedExtraInfo:
		nl.Events.EventResponseReceivedExtraInfo = v

	case *network.EventResponseReceived:
		nl.Events.EventResponseReceived = v
		nl.Status = v.Response.Status
		nl.Time = v.Timestamp.Time().Sub(nl.Events.EventRequestWillBeSent.Timestamp.Time()).Milliseconds()

	case *network.EventDataReceived:
		nl.Events.EventDataReceived = v
		if v.DataLength > 0 {
			nl.Size = v.DataLength
		}
		nl.Time = v.Timestamp.Time().Sub(nl.Events.EventRequestWillBeSent.Timestamp.Time()).Milliseconds()

	case *network.EventLoadingFinished:
		nl.Events.EventLoadingFinished = v
		nl.Time = v.Timestamp.Time().Sub(nl.Events.EventRequestWillBeSent.Timestamp.Time()).Milliseconds()

	case *network.EventLoadingFailed:
		nl.Events.EventLoadingFailed = v

	case *network.EventRequestServedFromCache:
		nl.ServedFromCache = true

	default:
		return fmt.Errorf("unknown event type: %T", ev)
	}

	return nil
}

func (nl *NetworkLog) IsFailed() bool {
	if nl.Error() != nil {
		return true
	}
	if nl.Status >= 400 {
		return true
	}
	if nl.IsFinished() && nl.Status <= 0 {
		return true
	}
	return false
}

func (nl *NetworkLog) IsFinished() bool {
	return nl.Events.EventLoadingFinished != nil
}

func (nl *NetworkLog) Error() error {
	if nl.Events.EventLoadingFailed == nil {
		return nil
	}

	if nl.Events.EventLoadingFailed.ErrorText != "" {
		return fmt.Errorf(nl.Events.EventLoadingFailed.ErrorText)
	} else if nl.Events.EventLoadingFailed.Canceled {
		return fmt.Errorf("canceled")
	} else {
		return fmt.Errorf("unknown error: %+v", nl.Events.EventLoadingFailed)
	}
}

func (nl *NetworkLog) Metrics() (nm *NetworkMetrics) {
	nm = new(NetworkMetrics)

	if nl.Events.EventResponseReceived != nil && nl.Events.EventResponseReceived.Response != nil && nl.Events.EventResponseReceived.Response.Timing != nil {
		nm.Dns = nl.Events.EventResponseReceived.Response.Timing.DNSEnd - nl.Events.EventResponseReceived.Response.Timing.DNSStart
		nm.Connection = nl.Events.EventResponseReceived.Response.Timing.ConnectEnd - nl.Events.EventResponseReceived.Response.Timing.ConnectStart
		nm.Proxy = nl.Events.EventResponseReceived.Response.Timing.ProxyEnd - nl.Events.EventResponseReceived.Response.Timing.ProxyStart
		nm.SSL = nl.Events.EventResponseReceived.Response.Timing.SslEnd - nl.Events.EventResponseReceived.Response.Timing.SslStart
		nm.RequestSent = nl.Events.EventResponseReceived.Response.Timing.SendEnd - nl.Events.EventResponseReceived.Response.Timing.SendStart
		nm.TTFB = nl.Events.EventResponseReceived.Response.Timing.ReceiveHeadersEnd - nl.Events.EventResponseReceived.Response.Timing.SendEnd
	}

	if nl.Events.EventDataReceived != nil && nl.Events.EventResponseReceived != nil {
		nm.ContentDownload = float64(nl.Events.EventDataReceived.Timestamp.Time().Sub(nl.Events.EventResponseReceived.Timestamp.Time())) / float64(time.Millisecond)
	}

	return
}
