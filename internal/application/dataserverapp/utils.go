package dataserverapp

import (
	"errors"
	"net/url"
	"strings"
)

const (
	assetPosition = 2

	bearerSuffix = "Bearer "
)

var (
	errBadURL      = errors.New("bad URL")
	errIsNotBearer = errors.New("is not beaere")
)

func getAsset(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", errBadURL
	}
	vals := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(vals) != 3 {
		return "", errBadURL
	}
	return vals[assetPosition], nil
}

func getBearerToken(source string) (string, error) {
	if !isBearer(source) {
		return "", errIsNotBearer
	}

	token := strings.TrimPrefix(source, bearerSuffix)
	return token, nil
}

func isBearer(source string) bool {
	return strings.HasPrefix(source, bearerSuffix)
}
