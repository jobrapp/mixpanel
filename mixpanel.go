package mixpanel

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

type mixpanel struct {
	token string
}

/*
	example track event data:
  {
    "event": "Signed Up",
    "properties": {
      "distinct_id": "13793",
      "token": "e3bc4100330c35722740fb8c6f5abddc",
      "Referred By": "Friend"
    }
  }
	example engage event data:
	{
		"$token": "36ada5b10da39a1347559321baf13063",
		"$distinct_id": "13793",
		"$ip": "123.123.123.123",
		"$set": {
			"Address": "1313 Mockingbird Lane"
		}
	}
*/

type eventData struct {
	Event string                 `json:"event"`
	Props map[string]interface{} `json:"properties"`
}

type engageData struct {
	Token   string      `json:"token"`
	Time    int64       `json:"time"`
	Id      int64       `json:"distinct_id"`
	Ip      string      `json:"ip,omitempty"`
	Set     interface{} `json:"set,omitempty"`
	SetOnce interface{} `json:"set_once,omitempty"`
	Add     interface{} `json:"add,omitempty"`
	Append  interface{} `json:"append,omitempty"`
	Union   interface{} `json:"union,omitempty"`
	Unset   interface{} `json:"unset,omitempty"`
	Delete  interface{} `json:"delete,omitempty"`
}

func New(token string) *mixpanel {
	return &mixpanel{
		token: token,
	}
}

func (mp *mixpanel) Track(uid int64, e string, p map[string]interface{}, params ...map[string]interface{}) bool {
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

	u := fmt.Sprintf("%s/%s/?data=%s", host, trackPath, base64.StdEncoding.EncodeToString(marshaledData))

	parameters := url.Values{}
	for _, val := range params {
		for k, v := range val {
			parameters.Add(k, v.(string))
		}
	}
	if qs := parameters.Encode(); qs != "" {
		u += "&" + qs
	}

	_, err = http.Get(u)
	if err != nil {
		return false
	}
	return true
}

func (mp *mixpanel) Engage(uid int64, p map[string]interface{}, ip string) bool {
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
	for k, v := range p {
		switch k {
		case "set":
			profileData.Set = v
			break
		case "set_once":
			profileData.SetOnce = v
			break
		case "add":
			profileData.Add = v
			break
		case "append":
			profileData.Append = v
			break
		case "union":
			profileData.Union = v
			break
		case "unset":
			profileData.Unset = v
			break
		case "delete":
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
