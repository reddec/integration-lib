package telegram

import "testing"

func TestTelegram_SendText(t *testing.T) {
	tg := Default()
	err := tg.SendText("Hell in world!")
	if err != nil {
		t.Error(err)
	}
}
