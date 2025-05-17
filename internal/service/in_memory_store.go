package service

import (
	"context"
	"errors"
)

var (
	KVStore     map[string]string
	NotFoundErr = "key not found"
)

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Get(ctx context.Context, key string) (string, error) {

	if v, ok := KVStore[key]; ok {
		return v, nil
	}

	return "", errors.New(NotFoundErr)

}

// TODO: how to handle collisions? Probably just overwrite for now tbh
func Set(ctx context.Context, key, val string) error {
	KVStore[key] = val

	return nil
}
