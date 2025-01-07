package configs

type serverConfig struct {
	API_Prefix string
}

var Server serverConfig

func init() {
	Server.API_Prefix = "/api/v1"
}
