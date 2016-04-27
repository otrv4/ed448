package edwards448

import "testing"

func TestOnCurve(t *testing.T) {
	ed448 := Ed448()
	assert_true(t, ed448.IsOnCurve(ed448.Params().Gx, ed448.Params().Gy))
}

func TestDouble(t *testing.T) {
	ed448 := Ed448()
	x2, y2 := ed448.Double(ed448.Params().Gx, ed448.Params().Gy)
	x3, y3 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	assert_equals(t, x2, x3)
	assert_equals(t, y2, y3)
}
