package utils

import (
	"errors"
	"os"
	"strconv"
)

func GetenvInt64(key string) (int64, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func GetenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")
