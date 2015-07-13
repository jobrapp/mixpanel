package tracking

import (
	"bytes"
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
	EngageSet     = "set"
	EngageSetOnce = "set_once"
	EngageAdd     = "add"
	EngageAppend  = "append"
	EngageUnion   = "union"
	EngageUnset   = "unset"
	EngageDelete  = "delete"
)

type client struct {
	token string
}

type eventData struct {
	Event string                 `json:"event"`
	Props map[string]interface{} `json:"properties"`
}

type UserEvent struct {
	DistinctId int64
	Name       string
	Props      map[string]interface{}
}

type UserEventStringed struct {
	DistinctId string
	Name       string
	Props      map[string]interface{}
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

type engageDataStringed struct {
	Token   string      `json:"$token"`
	Time    int64       `json:"$time"`
	Id      string      `json:"$distinct_id"`
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

func (mp *client) Track(event UserEvent, queryParams ...map[string]interface{}) error {
	u, err := mp.CreateLink(event, queryParams...)
	if err != nil {
		return err
	}

	// send request
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (mp *client) CreateLink(event UserEvent, queryParams ...map[string]interface{}) (string, error) {
	data := &eventData{
		Event: event.Name,
		Props: map[string]interface{}{
			"token":       mp.token,
			"distinct_id": strconv.FormatInt(event.DistinctId, 10),
		},
	}
	for k, v := range event.Props {
		data.Props[k] = v
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("%s/%s/?data=%s", host, trackPath,
		base64.StdEncoding.EncodeToString(marshaledData))

	parameters := url.Values{}
	// iterate over any query parameters
	for _, val := range queryParams {
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

	return u, nil
}

func (mp *client) TrackBatch(events []UserEvent, queryParams ...map[string]interface{}) error {
	// enforce max batch size to 50 (mixpanel)
	maxBatch := 50
	if len(events) > maxBatch {
		moreEvents := events[maxBatch:]
		events = events[:maxBatch]
		defer mp.TrackBatch(moreEvents, queryParams...)
	}
	data := make([]*eventData, 0, maxBatch)
	for _, event := range events {
		d := &eventData{
			Event: event.Name,
			Props: map[string]interface{}{
				"time":        time.Now().Unix(), // default to now, can be overwritten by props
				"token":       mp.token,
				"distinct_id": strconv.FormatInt(event.DistinctId, 10),
			},
		}
		for k, v := range event.Props {
			d.Props[k] = v
		}
		data = append(data, d)
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	u := fmt.Sprintf("%s/%s", host, trackPath)
	encodedData := "data=" + base64.StdEncoding.EncodeToString(marshaledData)

	parameters := url.Values{}
	// iterate over any query parameters
	for _, val := range queryParams {
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
		u += "?" + qs
	}
	// send request
	resp, err := http.Post(u, "application/x-www-form-urlencoded", bytes.NewBufferString(encodedData))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (mp *client) Engage(uid int64, p map[string]interface{}, ip string) error {
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
		return err
	}

	url := fmt.Sprintf("%s/%s/?data=%s", host, engagePath, base64.StdEncoding.EncodeToString(marshalledData))

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (mp *client) EngageString(id string, p map[string]interface{}, ip string) error {
	profileData := &engageDataStringed{
		Id:    id,
		Token: mp.token,
		Time:  time.Now().Unix(),
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
		return err
	}

	url := fmt.Sprintf("%s/%s/?data=%s", host, engagePath, base64.StdEncoding.EncodeToString(marshalledData))

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (mp *client) TrackString(event UserEventStringed, queryParams ...map[string]interface{}) error {
	data := &eventData{
		Event: event.Name,
		Props: map[string]interface{}{
			"time":        time.Now().Unix(), // default to now, can be overwritten by props
			"token":       mp.token,
			"distinct_id": event.DistinctId,
		},
	}
	for k, v := range event.Props {
		data.Props[k] = v
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	u := fmt.Sprintf("%s/%s/?data=%s", host, trackPath,
		base64.StdEncoding.EncodeToString(marshaledData))

	parameters := url.Values{}
	// iterate over any query parameters
	for _, val := range queryParams {
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
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
