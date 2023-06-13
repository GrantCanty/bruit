package kraken_data

var (
	ApiKey     string
	PrivateKey string
)

func LoadKeys(env map[string]string) {
	ApiKey = env["KRAKENAPIKEY"]
	PrivateKey = env["KRAKENPRIVATEKEY"]
}
