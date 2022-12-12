package kraken_data

var (
	ApiKey     string
	PrivateKey string
)

func LoadKeys(env map[string]string) {
	ApiKey = env["KRAKEN_API_KEY"]
	PrivateKey = env["KRAKEN_PRIVATE_KEY"]
}
