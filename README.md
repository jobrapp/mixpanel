mixpanel
========

golang mixpanel lib

need to implement all of this v
https://mixpanel.com/help/reference/http

usage
```go

mixpanel.New(token).Track(userId, event_name, propsMap, [queryParams])

```

```go
mixpanel.New(token).Engage(userId, propsMap, ip)
```


## Basic Usage
```
import (
  "mixpanel"
  "os"
)
.
.
.
// event properties 
var propsMap map[string]interface{}{
	"purchaseAmt": 9,
	"productName": "Shoe"
}

// optional query parameters
qparams := map[string]interface{}{
	"img":      1,
	"ip":       1,
	"callback": "someFuncName",
	"redirect": "https://www.google.com/someUrl?p=true",
}
// Send data to Mixpanel
mixpanel.New(os.getEnv("mixpanel_token")).Track(1, "purchased", propsMap, qparams)
```

if you specify query parameters without events properties, pass nil for propsMap.


to do: 
- test adding params work (string, int, big integer, bad values) -> done.
- test engage
- abstract params into helper
- apply params() function in engage