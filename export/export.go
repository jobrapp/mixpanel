package export

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

const (
	host = "https://data.mixpanel.com"
	path = "api/2.0/export"
)

type client struct {
	key string
	secret string
}

/*
request options:

required	from_date	string	
The date in yyyy-mm-dd format (UTC) from which to begin querying for the event from. This date is inclusive.

required	to_date	string	
The date in yyyy-mm-dd format (UTC) from which to stop querying for the event from. This date is inclusive.

optional	event	array	
The event or events that you wish to get data for, encoded as a JSON array.
Example format: '["play song", "log in", "add playlist"]'

optional	where	string	
An expression to filter events by. See the expression section on the main data export API page.

example request:
from_date = '2012-02-14'
to_date = '2012-02-14'
event = 'Viewed report'
where = 'properties["$os"]=="Linux"'

https://data.mixpanel.com/api/2.0/export/?from_date=2012-02-14&expire=1329760783&sig=bbe4be1e144d6d6376ef5484745aac45
&to_date=2012-02-14&api_key=f0aa346688cee071cd85d857285a3464&
where=properties%5B%22%24os%22%5D+%3D%3D+%22Linux%22&event=%5B%22Viewed+report%22%5D
*/

func (mp *client) constructQueryString(data *Data) string {
	v := url.Values{}
	v.Set("expire", fmt.Sprintf("%d", time.Now().Unix() + 60)) // expiry = 1m from now
	v.Set("api_key", mp.key)
	v.Set("from_date", data.FromDate)
	v.Set("to_date", data.ToDate)
	if len(data.Event) > 0 {
		b, e := json.Marshal(data.Event)
		if e != nil {
			panic(e)
		}
		v.Set("event", string(b))
	}

	/*
	if data.Props != "" {
	}
	*/

	/*
	pseudo-code for `sig` calculation
	https://mixpanel.com/docs/api-documentation/data-export-api#auth-implementation
	-----------
	args = all query parameters going to be sent out with the request 
	(e.g. api_key, unit, interval, expire, format, etc.) excluding sig.

	args_sorted = sort_args_alphabetically_by_key(args)

	args_concat = join(args_sorted) 

	# Output: api_key=ed0b8ff6cc3fbb37a521b40019915f18event=["pages"]
	#         expire=1248499222format=jsoninterval=24unit=hour

	sig = md5(args_concat + api_secret)
	*/
	sortedKeys := make([]string, len(v))
	i := 0
	for k, _ := range v {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	argsConcat := ""
	for _, k := range sortedKeys {
		argsConcat += fmt.Sprintf("%s=%s", k, v.Get(k))
	}
	// wtf did i just write
	v.Set("sig", fmt.Sprintf("%x", md5.Sum([]byte(argsConcat + mp.secret))))

	return v.Encode()
}

type Data struct {
	FromDate string
	ToDate   string
	Event    []string
	Props    string
}

func New(key, secret string) *client {
	return &client{
		key: key,
		secret: secret,
	}
}

func (mp *client) Export(data *Data) ([]byte, error) {
	if data.FromDate == "" || data.ToDate == "" {
		return nil, errors.New("Calls to (*mixpanel/export/Data).Export must contain both the `Data.FromDate` and `Data.ToDate` fields")
	}
	url := fmt.Sprintf("%s/%s/?%s", host, path, mp.constructQueryString(data))
	// send request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
