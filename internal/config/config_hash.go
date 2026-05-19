package config

import (
	"hash/crc64"

	"gopkg.in/yaml.v3"
)

var hashTable = crc64.MakeTable(crc64.ECMA)

func getHash[T ConfigType](value T) uint64 {
	content, err := yaml.Marshal(value)
	if err != nil {
		return 0
	}

	crc := crc64.New(hashTable)

	_, err = crc.Write(content)
	if err != nil {
		return 0
	}

	return crc.Sum64()
}
