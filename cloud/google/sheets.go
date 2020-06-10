package google

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Creates default sheets config
//
//     GOOGLE_APPLICATION_CREDENTIALS - (mandatory) path to service json file (https://cloud.google.com/docs/authentication/production#creating_a_service_account)
func DefaultSheets() *SheetsConfig {
	cfg := &SheetsConfig{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	if cfg.CredentialsFile == "" {
		panic("credentials file location required")
	}
	return cfg
}

// Google spread sheets config.
//
// Don't forget to share your spreadsheet with service account email
type SheetsConfig struct {
	CredentialsFile string `env:"GOOGLE_APPLICATION_CREDENTIALS"` // (mandatory) path to service json file (https://cloud.google.com/docs/authentication/production#creating_a_service_account)
	SheetId         string `env:"GOOGLE_SHEET_ID"`                // default spreadsheet id
}

// Append values as row to first empty row to default sheet
func (cfg SheetsConfig) AppendRow(values ...interface{}) error {
	return cfg.AppendRowToContext(context.Background(), cfg.SheetId, "A1", values...)
}

// Append values as row to first empty row to default sheet
func (cfg SheetsConfig) AppendRowContext(ctx context.Context, values ...interface{}) error {
	return cfg.AppendRowToContext(ctx, cfg.SheetId, "A1", values...)
}

// Append row to google spread sheets
func (cfg SheetsConfig) AppendRowToContext(ctx context.Context, sheetID, appendRange string, values ...interface{}) error {
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(cfg.CredentialsFile))
	if err != nil {
		return fmt.Errorf("sheets - append row: initialize google service: %w", err)
	}
	vals := sheets.ValueRange{
		Values: [][]interface{}{values},
	}
	_, err = srv.Spreadsheets.Values.Append(sheetID, appendRange, &vals).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("sheets - append row: %w", err)
	}
	return nil
}
