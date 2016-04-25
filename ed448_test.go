package ed448

import "testing"

func TestOnCurve(t *testing.T) {
	if !ed448.IsOnCurve(ed448.Params().Gx, ed448.Params().Gy) {
		t.Errorf("FAIL")
	}
}
