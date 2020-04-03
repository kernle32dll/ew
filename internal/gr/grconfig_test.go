package gr_test

import (
	"github.com/kernle32dll/ew/internal"
	"github.com/kernle32dll/ew/internal/gr"

	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParseConfigFromGr(t *testing.T) {
	// Prepare test files
	brokenFile := testFileWithContent(t, "[}")
	defer os.Remove(brokenFile.Name())
	workingFile := testFileWithContent(t, `{"tags": {"tag1": ["path1a", "path1b"], "tag2": ["path2a"], "tag3": []}}`)
	defer os.Remove(workingFile.Name())

	tests := []struct {
		name     string
		filename string
		want     internal.Config
		wantErr  bool
	}{
		{name: "not existing", filename: "does-not-exist", want: internal.Config{}, wantErr: true},
		{name: "broken file", filename: brokenFile.Name(), want: internal.Config{}, wantErr: true},
		{name: "working file", filename: workingFile.Name(), want: internal.Config{
			Source: internal.JsonSrc,
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
				t.Errorf("ParseConfigFromGr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func testFileWithContent(t *testing.T, content string) *os.File {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err.Error())
	}

	return f
}
