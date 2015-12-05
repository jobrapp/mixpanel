package tracking

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mixpanel"
	"net/http"
	"net/url"
	"strconv"
)

// track event
func (mp *Client) Track(event UserEvent, opts ...EventOptions) error {
	u, err := mp.CreateTrackingUrl(event, opts...)
	if err != nil {
		return err
	}
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (mp *Client) CreateTrackingUrl(event UserEvent, opts ...EventOptions) (string, error) {
	data := &eventData{
		Event: event.Name,
		Properties: map[string]interface{}{
			"token":       mp.Token,
			"distinct_id": event.DistinctId,
			"time":        mixpanel.Now(), // default to now
		},
	}
	if event.Ip != "" {
		data.Properties["ip"] = event.Ip
	}
	if event.Time.IsZero() == false {
		data.Properties["time"] = mixpanel.TimeToMPFmt(event.Time)
	}
	for k, v := range event.Properties {
		data.Properties[k] = v
	}
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("data", base64.StdEncoding.EncodeToString(marshaledData))
	if len(opts) > 0 {
		opt := opts[0]
		if opt.Ip != 0 {
			params.Set("ip", strconv.Itoa(opt.Ip))
		}
		if opt.Redirect != "" {
			params.Set("redirect", opt.Redirect)
		}
		if opt.Img != 0 {
			params.Set("img", strconv.Itoa(opt.Img))
		}
		if opt.Callback != "" {
			params.Set("callback", opt.Callback)
		}
		if opt.Verbose != 0 {
			params.Set("verbose", strconv.Itoa(opt.Verbose))
		}
	}
	return fmt.Sprintf("%s/%s/?%s", host, trackPath, params.Encode()), nil
}

// engage
func (mp *Client) Set(distinctId string, fields map[string]interface{}, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.Set = fields
	return mp.sendEngage(data)
}
func (mp *Client) SetOnce(distinctId string, fields map[string]interface{}, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.SetOnce = fields
	return mp.sendEngage(data)
}
func (mp *Client) Add(distinctId string, fields map[string]interface{}, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.Add = fields
	return mp.sendEngage(data)
}
func (mp *Client) Append(distinctId string, fields map[string]interface{}, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.Append = fields
	return mp.sendEngage(data)
}
func (mp *Client) Union(distinctId string, fields map[string]interface{}, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.Union = fields
	return mp.sendEngage(data)
}
func (mp *Client) Unset(distinctId string, fields []string, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.Unset = fields
	return mp.sendEngage(data)
}
func (mp *Client) Delete(distinctId string, opts ...EngageOptions) error {
	data := mp.setBaseData(distinctId, opts)
	data.Delete = "delete"
	return mp.sendEngage(data)
}

// create alias
func (mp *Client) CreateAlias(distinctId, newId string) error {
	return mp.Track(UserEvent{
		DistinctId: distinctId,
		Name:       "$create_alias",
		Properties: map[string]interface{}{
			"alias": newId,
		},
	})
}
