package auth

import (
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/google/go-cmp/cmp"
)

func TestHashPassword(t *testing.T) {

	tests := map[string]struct {
		input string
		want  bool
	}{
		"simple": {
			input: "somePassword",
			want:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hash, _ := HashPassword(tc.input)
			got, _ := argon2id.ComparePasswordAndHash(tc.input, hash)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf("%v", diff)
			}
		})
	}
}
