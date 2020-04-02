package internal

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type ReadSource int

const (
	JsonSrc ReadSource = iota
	YamlSrc
)

type encodable interface {
	Encode(v interface{}) error
}

type Config struct {
	Source ReadSource
	Tags   Tags `json:"tags" yaml:"tags"`
}

type Tags map[string][]string

func (c Config) GetTagsSorted() []string {
	tags := make([]string, len(c.Tags))
	i := 0
	for tag := range c.Tags {
		tags[i] = tag
		i++
	}

	sort.Strings(tags)

	return tags
}

func (c Config) GetPathsOfTagSorted(tag string) []string {
	paths := make([]string, len(c.Tags[tag]))
	copy(paths, c.Tags[tag])

	sort.Strings(paths)

	return paths
}

func (c Config) GetPathsOfTagsSorted(tags ...string) []string {
	// Fast-path
	if len(tags) == 0 {
		return c.GetPathsOfTagSorted(tags[0])
	}

	pathSet := make(map[string]struct{})

	for _, tag := range tags {
		for _, folder := range c.Tags[tag] {
			pathSet[folder] = struct{}{}
		}
	}

	paths := make([]string, len(pathSet))
	i := 0
	for path := range pathSet {
		paths[i] = path
		i++
	}

	sort.Strings(paths)

	return paths
}

func (c Config) GetTagsOfPathSorted(path string) []string {
	var tags []string

	for tag, paths := range c.Tags {
		if contains(paths, path) {
			tags = append(tags, tag)
		}
	}

	sort.Strings(tags)

	return tags
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ParseConfigFromFolder(path string) Config {
	cleanPath := strings.TrimRight(path, "/") + "/"

	yamlConf, err := parseConfigFromYaml(cleanPath)
	if err == nil {
		return yamlConf
	}

	jsonConf, err := parseConfigFromJson(cleanPath)
	if err == nil {
		return jsonConf
	}

	// If no config is found, use default yaml
	return Config{Source: YamlSrc}
}

func parseConfigFromYaml(path string) (Config, error) {
	f, err := os.Open(path + ".ewconfig.yml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)

	config := Config{
		Source: YamlSrc,
	}
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func parseConfigFromJson(path string) (Config, error) {
	f, err := os.Open(path + ".ewconfig.json")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	config := Config{
		Source: JsonSrc,
	}
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) WriteConfig(path string) (string, error) {
	cleanPath := strings.TrimRight(path, "/") + "/"

	// default to yml
	filename := ".ewconfig.yml"

	switch c.Source {
	case JsonSrc:
		filename = ".ewconfig.json"
	case YamlSrc:
		filename = ".ewconfig.yml"
	}

	f, err := os.Create(cleanPath + filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// default to yml
	var encoder encodable = yaml.NewEncoder(f)

	switch c.Source {
	case JsonSrc:
		jsonEncoder := json.NewEncoder(f)
		jsonEncoder.SetIndent("", "  ")
		encoder = jsonEncoder
	case YamlSrc:
		encoder = yaml.NewEncoder(f)
	}

	if err := encoder.Encode(c); err != nil {
		return "", err
	}

	return f.Name(), nil
}
