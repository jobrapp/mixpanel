package tracking

// public - engage - special keys
// applies to Set, SetOnce, Add, Append, Union, Unset
const (
	FirstName = "$first_name"
	LastName  = "$last_name"
	Name      = "$name"
	Created   = "$created"
	Email     = "$email"
	Phone     = "$phone"
)

// private - url
const (
	host       = "https://api.mixpanel.com"
	trackPath  = "track"
	engagePath = "engage"
)
