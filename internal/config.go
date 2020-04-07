package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v2"
)

// ReadSource determinate from which source a given config was read,
//and/or in which format it should be persisted.
type ReadSource int

const (
	JsonSrc ReadSource = iota
	YamlSrc
)

type encodable interface {
	Encode(v interface{}) error
}

// Config contains all runtime configuration for ew, such as
// available tags.
type Config struct {
	Source ReadSource `json:"-" yaml:"-"`
	Tags   Tags       `json:"tags" yaml:"tags"`
}

// Tags is a convenience wrapper around map[string][]string
type Tags map[string][]string

// GetTagsSorted returns a sorted list of configured tags.
func (c Config) GetTagsSorted() []string {
	if len(c.Tags) == 0 {
		return []string{}
	}

	tags := make([]string, len(c.Tags))
	i := 0
	for tag := range c.Tags {
		tags[i] = tag
		i++
	}

	sort.Strings(tags)

	return tags
}

// GetPathsOfTagSorted returns a sorted list of paths of the given tag.
func (c Config) GetPathsOfTagSorted(tag string) []string {
	tagPaths, exists := c.Tags[tag]
	if !exists {
		return []string{}
	}

	if len(tagPaths) == 0 {
		return []string{}
	}

	paths := make([]string, len(tagPaths))
	copy(paths, c.Tags[tag])

	sort.Strings(paths)

	return paths
}

// GetPathsOfTagSorted returns a sorted list of paths of the given tags.
// Note, paths are sorted tag agnostic, so mixing up might occur.
// Duplicates are also filtered out.
func (c Config) GetPathsOfTagsSorted(tags ...string) []string {
	// Fast-path
	if len(tags) == 0 {
		return []string{}
	} else if len(tags) == 1 {
		return c.GetPathsOfTagSorted(tags[0])
	}

	pathSet := make(map[string]struct{})

	for _, tag := range tags {
		pathTags, exists := c.Tags[tag]

		if exists {
			for _, folder := range pathTags {
				pathSet[folder] = struct{}{}
			}
		}
	}

	// Another fast path
	if len(pathSet) == 0 {
		return []string{}
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

// GetTagsOfPathSorted returns a sorted list of tags by which
// the given path is tagged with.
func (c Config) GetTagsOfPathSorted(path string) []string {
	var tags []string

	for tag, paths := range c.Tags {
		if contains(paths, path) {
			tags = append(tags, tag)
		}
	}

	if len(tags) == 0 {
		return []string{}
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

// ParseConfigFromFolder parses the given folder for
// a valid ew config, or returns the default (empty)
// config if none can be found.
func ParseConfigFromFolder(path string) Config {
	yamlConf, err := parseConfigFromYaml(path)
	if err == nil {
		return yamlConf
	}

	jsonConf, err := parseConfigFromJson(path)
	if err == nil {
		return jsonConf
	}

	// If no config is found, use default yaml
	return Config{Source: YamlSrc}
}

func parseConfigFromYaml(path string) (Config, error) {
	f, err := os.Open(filepath.Join(path, ".ewconfig.yml"))
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
	f, err := os.Open(filepath.Join(path, ".ewconfig.json"))
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

// WriteConfig writes the config to the given folder.
// Naming of the file is derived from the read source of
// the config.
func (c *Config) WriteConfig(path string) (string, error) {
	// default to yml
	filename := ".ewconfig.yml"

	switch c.Source {
	case JsonSrc:
		filename = ".ewconfig.json"
	case YamlSrc:
		filename = ".ewconfig.yml"
	}

	f, err := os.Create(filepath.Join(path, filename))
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
