package tracking

import (
	"time"
)

// public structs
type Client struct {
	Token string
}
type UserEvent struct {
	DistinctId string
	Name       string
	Time       time.Time
	Ip         string
	Properties map[string]interface{}
}
type EventOptions struct {
	Ip       int    // if 1, use incoming IP as distinct id
	Redirect string // if present, tracking request will redirect to value
	Img      int    // if 1, tracking request returns 1x1 transparent pixel
	Callback string // if present, tracking request will return JS callback
	Verbose  int    // if 1, tracking request will return status + error
}
type EngageOptions struct {
	Time       int64  `json:"$time,omitempty"`
	Ip         string `json:"$ip,omitempty"`
	IgnoreTime bool   `json:"$ignore_time,omitempty"`
}

// private structs
type eventData struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}
type engageData struct {
	Token      string `json:"$token"`
	DistinctId string `json:"$distinct_id"`
	EngageOptions

	Set     map[string]interface{} `json:"$set,omitempty"`
	SetOnce map[string]interface{} `json:"$set_once,omitempty"`
	Add     map[string]interface{} `json:"$add,omitempty"`
	Append  map[string]interface{} `json:"$append,omitempty"`
	Union   map[string]interface{} `json:"$union,omitempty"`
	Unset   []string               `json:"$unset,omitempty"`
	Delete  string                 `json:"$delete,omitempty"`
}
