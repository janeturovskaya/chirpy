package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"clear": {
			input: "does not contain any bad words",
			want:  "does not contain any bad words",
		},
		"contains_two_bad": {
			input: "This is a kerfuffle opinion sharbert",
			want:  "This is a **** opinion ****",
		},
		"contains_upper!": {
			input: "This is a Kerfuffle opinion",
			want:  "This is a **** opinion",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := clearProfineWords(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf("%v", diff)
			}
		})
	}
}
