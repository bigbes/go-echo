package server

// Configuration holds configuration for instrumentation server.
type Configuration struct {
	ListenAddress string `insconfig:"0.0.0.0:9000| Address to listen for instrumentation server on"`
}

// NewConfiguration creates new default configuration for server.
func NewConfiguration() Configuration {
	return Configuration{
		ListenAddress: "0.0.0.0:9090",
	}
}
