package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})
