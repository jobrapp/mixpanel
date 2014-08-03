package tracking

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	host       = "http://api.mixpanel.com"
	trackPath  = "track"
	engagePath = "engage"
)

// engage constants
const (
	EngageSet = "set"
	EngageSetOnce = "set_once"
	EngageAdd = "add"
	EngageAppend = "append"
	EngageUnion = "union"
	EngageUnset = "unset"
	EngageDelete = "delete"
)

type client struct {
	token string
}

type eventData struct {
	Event string                 `json:"event"`
	Props map[string]interface{} `json:"properties"`
}

type engageData struct {
	Token   string      `json:"$token"`
	Time    int64       `json:"$time"`
	Id      int64       `json:"$distinct_id"`
	Ip      string      `json:"$ip,omitempty"`
	Set     interface{} `json:"$set,omitempty"`
	SetOnce interface{} `json:"$set_once,omitempty"`
	Add     interface{} `json:"$add,omitempty"`
	Append  interface{} `json:"$append,omitempty"`
	Union   interface{} `json:"$union,omitempty"`
	Unset   interface{} `json:"$unset,omitempty"`
	Delete  interface{} `json:"$delete,omitempty"`
}

func New(token string) *client {
	return &client{
		token: token,
	}
}

func (mp *client) Track(uid int64, e string, p map[string]interface{}, params ...map[string]interface{}) bool {
	data := &eventData{
		Event: e,
		Props: map[string]interface{}{
			"time":  time.Now().Unix(),
			"token": mp.token,
		},
	}
	if uid != 0 {
		data.Props["distinct_id"] = strconv.Itoa(int(uid))
	}
	for k, v := range p {
		data.Props[k] = v
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		return false
	}

	u := fmt.Sprintf("%s/%s/?data=%s", host, trackPath,
		base64.StdEncoding.EncodeToString(marshaledData))

	parameters := url.Values{}
	// iterate over any query parameters
	for _, val := range params {
		for k, v := range val {
			if str, ok := v.(string); ok {
				/* act on str */
				parameters.Add(k, str)
			} else {
				/* not string - int? */
				if in, ok := v.(int); ok {
					parameters.Add(k, strconv.Itoa(in))
				} else {
					continue
				}
			}

		}
	}
	// append encoded params to url if any
	if qs := parameters.Encode(); qs != "" {
		u += "&" + qs
	}
	// send request
	_, err = http.Get(u)
	if err != nil {
		return false
	}
	return true
}

func (mp *client) Engage(uid int64, p map[string]interface{}, ip string) bool {
	profileData := &engageData{
		Token: mp.token,
		Time:  time.Now().Unix(),
	}
	if uid != 0 {
		profileData.Id = uid
	}
	if ip != "" {
		profileData.Ip = ip
	}
	// should probably just add separate methods for each of these
	for k, v := range p {
		switch k {
		case EngageSet:
			profileData.Set = v
			break
		case EngageSetOnce:
			profileData.SetOnce = v
			break
		case EngageAdd:
			profileData.Add = v
			break
		case EngageAppend:
			profileData.Append = v
			break
		case EngageUnion:
			profileData.Union = v
			break
		case EngageUnset:
			profileData.Unset = v
			break
		case EngageDelete:
			profileData.Delete = v
			break
		}
	}

	marshalledData, err := json.Marshal(profileData)
	if err != nil {
		return false
	}

	url := fmt.Sprintf("%s/%s/?data=%s", host, engagePath,
		base64.StdEncoding.EncodeToString(marshalledData))

	_, err = http.Get(url)
	if err != nil {
		return false
	}
	return true
}
