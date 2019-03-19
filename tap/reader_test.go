package tap

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReader(t *testing.T) {
	tests := map[string]struct {
		Input   string
		Want    *RunReport
		WantErr string
	}{
		"empty": {
			Input:   ``,
			WantErr: `no tests`,
		},
		"anonymous test with no plan": {
			Input: `ok`,
			Want:  &RunReport{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rr := strings.NewReader(test.Input)
			r := NewReader(rr)
			got, err := r.ReadAll()

			if err != nil {
				if test.WantErr == "" {
					t.Fatalf("unexpected error\ngot:  %s\nwant: success", err.Error())
				}
				if got, want := err.Error(), test.WantErr; got != want {
					t.Fatalf("unexpected error\ngot:  %s\nwant: %s", got, want)
				}
				return
			}

			if !cmp.Equal(got, test.Want) {
				t.Fatalf("wrong result\n%s", cmp.Diff(test.Want, got))
			}
		})
	}
}
