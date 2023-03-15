package gr_test

import (
	"github.com/kernle32dll/ew/internal"
	"github.com/kernle32dll/ew/internal/gr"

	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseConfigFromGr(t *testing.T) {
	// Prepare test files
	brokenFileName := testFileWithContent(t, "[}")
	workingFileName := testFileWithContent(t, `{"tags": {"tag1": ["path1a", "path1b"], "tag2": ["path2a"], "tag3": []}}`)

	tests := []struct {
		name     string
		filename string
		want     internal.Config
		wantErr  bool
	}{
		{name: "not existing", filename: "does-not-exist", want: internal.Config{}, wantErr: true},
		{name: "broken file", filename: brokenFileName, want: internal.Config{}, wantErr: true},
		{name: "working file", filename: workingFileName, want: internal.Config{
			Source:     internal.JsonSrc,
			LoadedFrom: filepath.Dir(workingFileName),
			Tags: map[string][]string{
				"tag1": {"path1a", "path1b"},
				"tag2": {"path2a"},
				"tag3": {},
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gr.ParseConfigFromGr(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfigFromGr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfigFromGr/%s got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func testFileWithContent(t *testing.T, content string) string {
	folder := t.TempDir()

	f, err := os.CreateTemp(folder, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err.Error())
	}

	if err := f.Close(); err != nil {
		t.Fatal(err.Error())
	}

	return f.Name()
}
