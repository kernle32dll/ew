package gr

import (
	"github.com/kernle32dll/ew/internal"

	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	Tags tags `json:"tags"`
}

type tags map[string][]string

// ParseConfigFromGr tries to parse an config file from the
// original mixu/gr json config.
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
		Source:     internal.JsonSrc,
		LoadedFrom: filepath.Dir(filename),
		Tags:       map[string][]string(config.Tags),
	}, nil
}
