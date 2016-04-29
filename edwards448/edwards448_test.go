package edwards448

import "testing"

func TestBasePointIsOnCurve(t *testing.T) {
	ed448 := Ed448()
	assert_true(t, ed448.IsOnCurve(ed448.Params().Gx, ed448.Params().Gy))
}

func TestAdd(t *testing.T) {
	ed448 := Ed448()

	x2, y2 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	x4, y4 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, x2, y2)

	assert_true(t, ed448.IsOnCurve(x2, y2))
	assert_true(t, ed448.IsOnCurve(x4, y4))
}

func TestDouble(t *testing.T) {
	ed448 := Ed448()

	xd2, yd2 := ed448.Double(ed448.Params().Gx, ed448.Params().Gy)
	xd4, yd4 := ed448.Double(xd2, yd2)

	assert_true(t, ed448.IsOnCurve(xd2, yd2))
	assert_true(t, ed448.IsOnCurve(xd4, yd4))
}

func TestOperationsAreEquivalent(t *testing.T) {
	ed448 := Ed448()

	addX, addY := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	doubleX, doubleY := ed448.Double(ed448.Params().Gx, ed448.Params().Gy)

	assert_equals(t, addX, doubleX)
	assert_equals(t, addY, doubleY)
}
