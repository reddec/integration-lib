package google

import (
	"testing"
	"time"
)

func TestSheetsConfig_AppendRowToContext(t *testing.T) {
	srv := DefaultSheets()
	err := srv.AppendRow(time.Now().Format("2006-01-02"), 1, 2, 3, 4)
	if err != nil {
		t.Error(err)
	}
}
