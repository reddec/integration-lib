package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/reddec/integration-lib/internal/utils"
)

// Default telegram configuration from env variables
//
// TELEGRAM_TOKEN - mandatory, bot token (from @BotFather), panic if not set
func Default() *Telegram {
	cfg := &Telegram{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	if cfg.Token == "" {
		panic("telegram token required")
	}
	return cfg
}

// Telegram endpoint
type Telegram struct {
	Token  string `env:"TELEGRAM_TOKEN"`   // Telegram bot token
	ChatID int64  `env:"TELEGRAM_CHAT_ID"` // Default telegram chat ID (https://t.me/getidsbot)
}

// Send basic text message to default chat ID (see SendTextToContext)
func (cfg Telegram) SendText(text string) error {
	return cfg.SendTextToContext(context.Background(), text, cfg.ChatID)
}

// Send basic text message to default chat ID (see SendTextToContext)
func (cfg Telegram) SendTextContext(ctx context.Context, text string) error {
	return cfg.SendTextToContext(ctx, text, cfg.ChatID)
}

// Send basic text message to custom chat ID (see SendTextToContext)
func (cfg Telegram) SendTextTo(text string, chatId int64) error {
	return cfg.SendTextToContext(context.Background(), text, chatId)
}

// Send basic text message to custom chat ID. Bot should have access to the chat, otherwise Bad Request will be raised
//
// API: https://api.telegram.org/
func (cfg Telegram) SendTextToContext(ctx context.Context, text string, chatId int64) error {
	var baseURL = "https://api.telegram.org/bot" + cfg.Token + "/sendMessage"
	var payload struct {
		Text   string `json:"text"`
		ChatID int64  `json:"chat_id"`
	}
	payload.Text = text
	payload.ChatID = chatId

	var reply json.RawMessage
	err := utils.PostJSON(ctx, baseURL, payload, &reply)
	if err != nil {
		return fmt.Errorf("telegram - send text: %w", err)
	}
	return nil
}
