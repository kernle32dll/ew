package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig_GetTagsSorted(t *testing.T) {
	tests := []struct {
		name string
		tags Tags
		want []string
	}{
		{name: "no tags exist", tags: Tags{}, want: []string{}},
		{name: "tags exist", tags: Tags{"a": []string{}, "c": []string{}, "b": []string{}}, want: []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Source: YamlSrc,
				Tags:   tt.tags,
			}
			if got := c.GetTagsSorted(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetPathsOfTagSorted(t *testing.T) {
	tests := []struct {
		name  string
		given Tags
		when  string
		want  []string
	}{
		{name: "no tags configured", given: Tags{}, when: "does-not-exist", want: []string{}},
		{name: "tags configured, but not found", given: Tags{"exists": {"a"}}, when: "does-not-exist", want: []string{}},
		{name: "tags configured, found, but no paths", given: Tags{"exists": {}}, when: "exists", want: []string{}},
		{name: "tags configured, found, and has paths", given: Tags{"exists": {"c", "a", "b"}}, when: "exists", want: []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Source: YamlSrc,
				Tags:   tt.given,
			}
			if got := c.GetPathsOfTagSorted(tt.when); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPathsOfTagSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetPathsOfTagsSorted(t *testing.T) {
	tests := []struct {
		name  string
		given Tags
		when  []string
		want  []string
	}{
		// Fast-path no tags requested
		{name: "no tags configured, none requested", given: Tags{}, when: []string{}, want: []string{}},

		{name: "no tags configured, two requested", given: Tags{}, when: []string{"does-not-exist", "does-not-exist-either"}, want: []string{}},
		{name: "tags configured, multiple requested", given: Tags{
			"exists1": {"c", "b"}, "exists2": {}, "exists3": {"a"},
		}, when: []string{"does-not-exist", "exists1", "exists2", "exists3"}, want: []string{
			"a", "b", "c",
		}},

		// Fast-path to single-tag func
		{name: "no tags configured, one requested", given: Tags{}, when: []string{"does-not-exist"}, want: []string{}},
		{name: "tags configured, one requested, but not found", given: Tags{"exists": {"one"}}, when: []string{"does-not-exist"}, want: []string{}},
		{name: "tags configured, one requested, found, but no paths", given: Tags{"exists": {}}, when: []string{"exists"}, want: []string{}},
		{name: "tags configured, one requested, found, and has paths", given: Tags{"exists": {"c", "a", "b"}}, when: []string{"exists"}, want: []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Source: YamlSrc,
				Tags:   tt.given,
			}
			if got := c.GetPathsOfTagsSorted(tt.when...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPathsOfTagsSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetTagsOfPathSorted(t *testing.T) {
	tests := []struct {
		name string
		tags Tags
		path string
		want []string
	}{
		{name: "no tags configured", tags: Tags{}, path: "does-not-exist", want: []string{}},
		{name: "matches no tags", tags: Tags{"tag": []string{"some-path"}}, path: "does-not-exist", want: []string{}},
		{name: "matches one tag", tags: Tags{"tag": []string{"unrelated-path", "some-path"}}, path: "some-path", want: []string{"tag"}},
		{name: "matches two tags",
			tags: Tags{"tag1": []string{"some-path", "unrelated-path"}, "tag2": []string{}, "tag3": []string{"some-path"}},
			path: "some-path",
			want: []string{"tag1", "tag3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Source: YamlSrc,
				Tags:   tt.tags,
			}
			if got := c.GetTagsOfPathSorted(tt.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsOfPathSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WriteConfig(t *testing.T) {
	tags := Tags{
		"some-tag":    []string{"path1", "path2"},
		"another-tag": []string{"path4", "path3"},
	}

	folder := t.TempDir()

	// --------------

	wantJson := `{
  "tags": {
    "another-tag": [
      "path4",
      "path3"
    ],
    "some-tag": [
      "path1",
      "path2"
    ]
  }
}
`

	wantYml := `tags:
    another-tag:
        - path4
        - path3
    some-tag:
        - path1
        - path2
`

	tests := []struct {
		name        string
		source      ReadSource
		path        string
		wantPath    string
		wantContent string
		wantErr     bool
	}{
		{name: "cannot create", source: YamlSrc, path: "does-not-exist", wantPath: "", wantErr: true, wantContent: ""},
		{name: "create json", source: JsonSrc, path: folder, wantPath: filepath.Join(folder, ".ewconfig.json"), wantErr: false, wantContent: wantJson},
		{name: "create yaml", source: YamlSrc, path: folder, wantPath: filepath.Join(folder, ".ewconfig.yml"), wantErr: false, wantContent: wantYml},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Source: tt.source,
				Tags:   tags,
			}
			got, err := c.WriteConfig(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantPath {
				t.Errorf("WriteConfig() path got = %v, want %v", got, tt.wantPath)
			}

			if got != "" {
				gotContent, err := os.ReadFile(got)
				if err != nil {
					t.Error(err)
					return
				}

				if string(gotContent) != tt.wantContent {
					t.Errorf("WriteConfig() content got = %v, want %v", string(gotContent), tt.wantContent)
				}
			}
		})
	}
}

func writeTempFile(t *testing.T, filename string, fileString string) string {
	folder := t.TempDir()

	if err := os.WriteFile(filepath.Join(folder, filename), []byte(fileString), 0600); err != nil {
		t.Fatal(err.Error())
	}

	return folder
}

func TestParseConfigFromFolder(t *testing.T) {
	t.Run("no config found", func(t *testing.T) {
		want := Config{Source: YamlSrc, LoadedFrom: "does-not-exist", Tags: map[string][]string{}}
		wantOutput := ""

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, "does-not-exist"); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})

	t.Run("json config found, and ok", func(t *testing.T) {
		fileString := `{"tags": {"some-tag": ["path1", "path2"]}}`

		folder := writeTempFile(t, ".ewconfig.json", fileString)

		want := Config{Source: JsonSrc, LoadedFrom: folder, Tags: Tags{"some-tag": []string{"path1", "path2"}}}
		wantOutput := ""

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, folder); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})

	t.Run("json config found, but borked", func(t *testing.T) {
		folder := writeTempFile(t, ".ewconfig.json", `{]`)

		want := Config{Source: YamlSrc, LoadedFrom: folder, Tags: map[string][]string{}}
		wantOutput := fmt.Sprintf(
			"Failed to read json config in %s: invalid character ']' looking for beginning of object key string\n",
			folder,
		)

		// Support for go 1.16 on master
		wantOutputEither := fmt.Sprintf(
			"Failed to read json config in %s: json: invalid character ']' looking for beginning of object key string\n",
			folder,
		)

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, folder); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput && got != wantOutputEither {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})

	t.Run("json config found, but EOF", func(t *testing.T) {
		folder := writeTempFile(t, ".ewconfig.json", ``)

		want := Config{Source: YamlSrc, LoadedFrom: folder, Tags: map[string][]string{}}
		wantOutput := fmt.Sprintf(
			"Skipping empty json config in %s\n",
			folder,
		)

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, folder); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})

	t.Run("yaml config found, and ok", func(t *testing.T) {
		fileString := `
tags:
  some-tag:
  - path1
  - path2
`

		folder := writeTempFile(t, ".ewconfig.yml", fileString)

		want := Config{Source: YamlSrc, LoadedFrom: folder, Tags: Tags{"some-tag": []string{"path1", "path2"}}}
		wantOutput := ""

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, folder); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})

	t.Run("yaml config found, but borked", func(t *testing.T) {
		folder := writeTempFile(t, ".ewconfig.yml", `t{]`)

		want := Config{Source: YamlSrc, LoadedFrom: folder, Tags: map[string][]string{}}
		wantOutput := fmt.Sprintf(
			"Failed to read yaml config in %s: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `t{]` into internal.Config\n",
			folder,
		)

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, folder); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})

	t.Run("yaml config found, but EOF", func(t *testing.T) {
		folder := writeTempFile(t, ".ewconfig.yml", ``)

		want := Config{Source: YamlSrc, LoadedFrom: folder, Tags: map[string][]string{}}
		wantOutput := fmt.Sprintf(
			"Skipping empty yaml config in %s\n",
			folder,
		)

		output := &bytes.Buffer{}
		if got := ParseConfigFromFolder(output, folder); !reflect.DeepEqual(got, want) {
			t.Errorf("ParseConfigFromFolder() = %v, want %v", got, want)
		}

		if got := output.String(); got != wantOutput {
			t.Errorf("ParseConfigFromFolder() output = %v, want %v", got, wantOutput)
		}
	})
}
