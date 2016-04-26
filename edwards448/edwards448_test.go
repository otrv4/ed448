package edwards448

import "testing"

func TestOnCurve(t *testing.T) {
	ed448 := Ed448()
	if !ed448.IsOnCurve(ed448.Params().Gx, ed448.Params().Gy) {
		t.Errorf("ed448 base point is not on curve")
	}
}

func TestDouble(t *testing.T) {
	ed448 := Ed448()
	x2, y2 := ed448.Double(ed448.Params().Gx, ed448.Params().Gy)
	x3, y3 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	if x2.Cmp(x3) != 0 || y2.Cmp(y3) != 0 {
		t.Errorf("x2:%v ,y2:%v!= x3:%v ,y3:%v", x2, y2, x3, y3)
	}
}
