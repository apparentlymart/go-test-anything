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
		"anonymous test success with no plan": {
			Input: `ok`,
			Want: &RunReport{
				Tests: []*Report{
					{Num: 1, Result: Pass},
				},
			},
		},
		"two anonymous test successes with no plan": {
			Input: "ok\nok",
			Want: &RunReport{
				Tests: []*Report{
					{Num: 1, Result: Pass},
					{Num: 2, Result: Pass},
				},
			},
		},
		"two anonymous tests, second failing, with no plan": {
			Input: "ok\nnot ok",
			Want: &RunReport{
				Tests: []*Report{
					{Num: 1, Result: Pass},
					{Num: 2, Result: Fail},
				},
			},
		},
		"one skipped test with no plan": {
			Input: "ok 1 thingy # skipped because no server is available",
			Want: &RunReport{
				Tests: []*Report{
					{
						Num:        1,
						Result:     Skip,
						Name:       "thingy",
						SkipReason: "skipped because no server is available",
					},
				},
			},
		},
		"one todo test with no plan": {
			Input: "not ok 1 wotsit # TODO: implement",
			Want: &RunReport{
				Tests: []*Report{
					{
						Num:        1,
						Result:     Fail,
						Name:       "wotsit",
						Todo:       true,
						TodoReason: "TODO: implement",
					},
				},
			},
		},
		"bail out after one test with no plan": {
			Input:   "ok 1 boop\nBail out! Database not available",
			WantErr: `testing aborted: Database not available`,
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
