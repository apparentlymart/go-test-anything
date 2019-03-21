package tap

import (
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	tests := map[string]struct {
		Steps func(w *Writer) error
		Want  string
	}{
		"no output": {
			func(w *Writer) error {
				return nil
			},
			"",
		},
		"one report with no plan": {
			func(w *Writer) error {
				return w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Pass,
				})
			},
			"ok 1 reticulates splines\n",
		},
		"two reports with no plan": {
			func(w *Writer) error {
				err := w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Fail,
				})
				if err != nil {
					return err
				}
				err = w.Report(&Report{
					Num:    2,
					Name:   "boops",
					Result: Pass,
				})
				if err != nil {
					return err
				}
				return nil
			},
			"not ok 1 reticulates splines\nok 2 boops\n",
		},
		"one report with prior plan": {
			func(w *Writer) error {
				err := w.Plan(&Plan{Min: 1, Max: 1})
				if err != nil {
					return err
				}
				return w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Pass,
				})
			},
			"1..1\nok 1 reticulates splines\n",
		},
		"one report with post-plan": {
			func(w *Writer) error {
				err := w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Pass,
				})
				if err != nil {
					return err
				}
				return w.Plan(&Plan{Min: 1, Max: 1})
			},
			"ok 1 reticulates splines\n1..1\n",
		},
		"two reports with late plan": {
			func(w *Writer) error {
				err := w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Pass,
				})
				if err != nil {
					return err
				}
				err = w.Plan(&Plan{Min: 1, Max: 1})
				if err != nil {
					return err
				}
				return w.Report(&Report{
					Num:    2,
					Name:   "boop",
					Result: Pass,
				})
			},
			"ok 1 reticulates splines\nok 2 boop\n1..1\n",
		},
		"auto-numbered reports": {
			func(w *Writer) error {
				err := w.Report(&Report{
					Name:   "reticulates splines",
					Result: Pass,
				})
				if err != nil {
					return err
				}
				err = w.Report(&Report{
					Name:   "boops",
					Result: Fail,
				})
				if err != nil {
					return err
				}
				return nil
			},
			"ok 1 reticulates splines\nnot ok 2 boops\n",
		},
		"skipped test with no reason": {
			func(w *Writer) error {
				return w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Skip,
				})
			},
			"ok 1 reticulates splines # SKIP\n",
		},
		"skipped test with reason": {
			func(w *Writer) error {
				return w.Report(&Report{
					Num:        1,
					Name:       "reticulates splines",
					Result:     Skip,
					SkipReason: "no splines",
				})
			},
			"ok 1 reticulates splines # SKIP: no splines\n",
		},
		"todo test with no reason": {
			func(w *Writer) error {
				return w.Report(&Report{
					Num:    1,
					Name:   "reticulates splines",
					Result: Fail,
					Todo:   true,
				})
			},
			"not ok 1 reticulates splines # TODO\n",
		},
		"todo test with reason": {
			func(w *Writer) error {
				return w.Report(&Report{
					Num:        1,
					Name:       "reticulates splines",
					Result:     Fail,
					Todo:       true,
					TodoReason: "not yet implemented",
				})
			},
			"not ok 1 reticulates splines # TODO: not yet implemented\n",
		},
		"bail out": {
			func(w *Writer) error {
				return w.BailOut("printer on fire")
			},
			"Bail out! printer on fire\n",
		},
		"diagnostic": {
			func(w *Writer) error {
				return w.Diagnostic("a is 2")
			},
			"# a is 2\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			w := NewWriter(&buf)
			err := test.Steps(w)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			err = w.Close()
			if err != nil {
				t.Fatalf("Unexpected error from Close: %s", err)
			}

			got := buf.String()
			want := test.Want

			if got != want {
				t.Errorf("wrong result\ngot:\n%s\n\nwant:\n%s", got, want)
			}
		})
	}
}
