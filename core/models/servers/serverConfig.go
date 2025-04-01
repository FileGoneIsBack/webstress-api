package servers

var (
	Servers map[string]*Server = make(map[string]*Server)
	Config  *Configuration
)

type Configuration struct {
	Listener int      `json:"listener"`
	Allowed  []string `json:"allowed"`
	Key      string   `json:"key"`
}
