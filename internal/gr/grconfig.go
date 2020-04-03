package gr

import (
	"github.com/kernle32dll/ew/internal"

	"encoding/json"
	"os"
)

type config struct {
	Tags tags `json:"tags"`
}

type tags map[string][]string

func ParseConfigFromGr(filename string) (internal.Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return internal.Config{}, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	var config config
	if err := decoder.Decode(&config); err != nil {
		return internal.Config{}, err
	}

	return internal.Config{
		Source: internal.JsonSrc,
		Tags:   map[string][]string(config.Tags),
	}, nil
}
