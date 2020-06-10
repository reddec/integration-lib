package tinkoff

import (
	"testing"
)

func TestInvest_Portfolio(t *testing.T) {
	investing := DefaultInvest()
	portfolio, err := investing.Portfolio()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", portfolio)
}
