package edwards448

import "testing"

func TestOnCurve(t *testing.T) {
	ed448 := Ed448()
	assert_true(t, ed448.IsOnCurve(ed448.Params().Gx, ed448.Params().Gy))
}

func TestAdd(t *testing.T) {
	ed448 := Ed448()
	x2, y2 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	assert_true(t, ed448.IsOnCurve(x2, y2))
	x4, y4 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, x2, y2)
	assert_true(t, ed448.IsOnCurve(x4, y4))
}

func TestDouble(t *testing.T) {
	ed448 := Ed448()
	xd2, yd2 := ed448.Double(ed448.Params().Gx, ed448.Params().Gy)
	xa2, ya2 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	xd4, yd4 := ed448.Double(xd2, yd2)
	xa4, ya4 := ed448.Add(xa2, ya2, xa2, ya2)
	assert_equals(t, xd2, xa2)
	assert_equals(t, yd2, ya2)
	assert_equals(t, xd4, xa4)
	assert_equals(t, yd4, ya4)
	assert_true(t, ed448.IsOnCurve(xd2, yd2))
	assert_true(t, ed448.IsOnCurve(xd4, yd4))
}
