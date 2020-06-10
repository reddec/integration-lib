package rates

import "testing"

func TestFiat(t *testing.T) {
	rates, err := Fiat("usd")
	if err != nil {
		t.Error(err)
		return
	}
	f, _ := rates["USD"].Float64()
	if f != 1.0 {
		t.Errorf("got rate %f", f)
		return
	}
}

func TestCrypto(t *testing.T) {
	rates, err := Crypto()
	if err != nil {
		t.Error(err)
		return
	}
	f, _ := rates["USDC"].Float64()
	if f != 1.0 {
		t.Errorf("got rate %f", f)
		return
	}
}
