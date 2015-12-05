package tracking

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func (mp *Client) setBaseData(distinctId string, opts []EngageOptions) engageData {
	data := engageData{
		Token:      mp.Token,
		DistinctId: distinctId,
	}
	if len(opts) > 0 {
		data.EngageOptions = opts[0]
	}
	return data
}

func (mp *Client) sendEngage(e engageData) error {
	marshalledData, err := json.Marshal(e)
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
