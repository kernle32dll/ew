package cmd

import (
	"github.com/fatih/color"

	"testing"
)

func Test_colorPath(t *testing.T) {
	// Force color for test
	color.NoColor = false

	want := "[90m/test/[0mpath"
	if got := colorPath("/test/path"); got != want {
		t.Errorf("colorPath() = %v, want %v", got, want)
	}
}

func Test_determinateNoPathsErrorMessage(t *testing.T) {
	// Force color for test
	color.NoColor = false

	tests := []struct {
		name    string
		forTags []string
		want    string
	}{
		{name: "no tags", forTags: nil, want: "[31mNo paths found[0m"},
		{name: "zero tags", forTags: []string{}, want: "[31mNo paths found[0m"},
		{name: "one tag", forTags: []string{"one"}, want: "[31mNo paths found for tag @one[0m"},
		{name: "two tags", forTags: []string{"one", "two"}, want: "[31mNo paths found for tags @one, @two[0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determinateNoPathsErrorMessage(tt.forTags); got != tt.want {
				t.Errorf("determinateNoPathsErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
