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
			Want:    &RunReport{},
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
		"two anonymous tests, one with diagnostics": {
			Input: "# Got output:\n#     foo\nok\nok",
			Want: &RunReport{
				Tests: []*Report{
					{
						Num:    1,
						Result: Pass,
						Diagnostics: []string{
							"Got output:",
							"    foo",
						},
					},
					{Num: 2, Result: Pass},
				},
			},
		},
		"one skipped test with no plan": {
			Input: "ok 1 thingy # skipped no server is available",
			Want: &RunReport{
				Tests: []*Report{
					{
						Num:        1,
						Result:     Skip,
						Name:       "thingy",
						SkipReason: "no server is available",
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
						TodoReason: "implement",
					},
				},
			},
		},
		"bail out after one test with no plan": {
			Input: "ok 1 boop\nBail out! Database not available",
			Want: &RunReport{
				Tests: []*Report{
					{
						Num:    1,
						Result: Pass,
						Name:   "boop",
					},
				},
			},
			WantErr: `testing aborted: Database not available`,
		},
		"one test, planned before": {
			Input: "1..1\nok 1",
			Want: &RunReport{
				Plan: &Plan{Min: 1, Max: 1},
				Tests: []*Report{
					{Num: 1, Result: Pass},
				},
			},
		},
		"one test, planned after": {
			Input: "ok 1\n1..1",
			Want: &RunReport{
				Plan: &Plan{Min: 1, Max: 1},
				Tests: []*Report{
					{Num: 1, Result: Pass},
				},
			},
		},
		"planned two but reported only one": {
			Input: "1..2\nok 1",
			Want: &RunReport{
				Plan: &Plan{Min: 1, Max: 2},
				Tests: []*Report{
					{Num: 1, Result: Pass},
					nil,
				},
			},
			WantErr: `no result for 2`,
		},
		"planned one but reported two": {
			Input: "1..1\nok 1\nok 2",
			Want: &RunReport{
				Plan: &Plan{Min: 1, Max: 1},
				Tests: []*Report{
					{Num: 1, Result: Pass},
				},
			},
			WantErr: `unexpected extra result for 2`,
		},
		"planned five but reported only one": {
			Input: "1..5\nok 1",
			Want: &RunReport{
				Plan: &Plan{Min: 1, Max: 5},
				Tests: []*Report{
					{Num: 1, Result: Pass},
					nil,
					nil,
					nil,
					nil,
				},
			},
			WantErr: `no result for 2-5`,
		},
		"planned five but reported only two, non-contiguous": {
			Input: "1..5\nok 1\nok 3",
			Want: &RunReport{
				Plan: &Plan{Min: 1, Max: 5},
				Tests: []*Report{
					{Num: 1, Result: Pass},
					nil,
					{Num: 3, Result: Pass},
					nil,
					nil,
				},
			},
			WantErr: `no result for 2, 4-5`,
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
			}

			if !cmp.Equal(got, test.Want) {
				t.Fatalf("wrong result\n%s", cmp.Diff(test.Want, got))
			}
		})
	}
}
