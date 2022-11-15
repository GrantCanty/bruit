package env

import (
	"errors"

	"github.com/joho/godotenv"
)

func Read() (map[string]string, error) {
	env, err := godotenv.Read()
	if err != nil {
		return nil, err
	}
	if file, found := env["READ"]; !found {
		return nil, errors.New("Could not find 'READ' field in .env file")
	} else if file == ".env" {
		return env, nil
	} else {
		return godotenv.Read(file)
	}
}
