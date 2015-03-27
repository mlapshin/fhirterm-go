package fhirterm

import (
	"fmt"
)

type Storage interface {
	FindValueSetById(id string) (*ValueSet, error)
}

type storageFactoryFunc func(cfg JsonObject) (Storage, error)

var storageFactories = map[string]storageFactoryFunc{
	"rest": MakeRestStorage,
}

var storage Storage = nil

func makeStorage(cfg JsonObject) (Storage, error) {
	if cfg == nil {
		return nil, fmt.Errorf("passed config is null or not a JSON object")
	}

	storageType, ok := cfg["type"].(string)
	if !ok {
		return nil, fmt.Errorf("passed config does not contain string attribute 'type'")
	}

	factory, found := storageFactories[storageType]
	if !found {
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}

	return factory(cfg)
}

func InitStorage(cfg JsonObject) error {
	var err error
	storage, err = makeStorage(cfg)

	if err != nil {
		return err
	}

	return nil
}

func GetStorage() Storage {
	return storage
}
