package email

import "testing"

func TestMail_SendText(t *testing.T) {
	mail := Default()
	err := mail.SendText("Test", "Hello world")
	if err != nil {
		t.Error(err)
	}
}
