package env

import (
	"errors"

	"github.com/joho/godotenv"
)

func Read(path string) (map[string]string, error) {
	env, err := godotenv.Read()
	if err != nil {
		return nil, err
	}
	if file, found := env[path]; !found {
		return nil, errors.New("Could not find 'CLIENT' field in .env file")
	} else if file == ".env" {
		return env, nil
	} else {
		return godotenv.Read(file)
	}
}

/*func ReadConfig() (map[string]string, error) {
	env, err := godotenv.Read()
	if err != nil {
		return nil, err
	}
	if file, found := env["CONFIG"]; !found {
		return nil, errors.New("Could not find 'CONFIG' field in .env file")
	} else if file == ".env" {
		return env, nil
	} else {
		return godotenv.Read(file)
	}
}
*/
