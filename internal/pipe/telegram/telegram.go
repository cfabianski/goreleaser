package telegram

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/caarlos0/log"
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/context"
)

const defaultMessageTemplate = `{{ .ProjectName }} {{ .Tag }} is out! Check it out at {{ .ReleaseURL }}`

type Pipe struct{}

func (Pipe) String() string                 { return "telegram" }
func (Pipe) Skip(ctx *context.Context) bool { return !ctx.Config.Announce.Telegram.Enabled }

type Config struct {
	ConsumerToken string `env:"TELEGRAM_TOKEN,notEmpty"`
}

func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Announce.Telegram.MessageTemplate == "" {
		ctx.Config.Announce.Telegram.MessageTemplate = defaultMessageTemplate
	}
	return nil
}

func (Pipe) Announce(ctx *context.Context) error {
	msg, err := tmpl.New(ctx).Apply(ctx.Config.Announce.Telegram.MessageTemplate)
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("telegram: %w", err)
	}

	log.Infof("posting: '%s'", msg)
	bot, err := api.NewBotAPI(cfg.ConsumerToken)
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}

	tm := api.NewMessage(ctx.Config.Announce.Telegram.ChatID, msg)
	tm.ParseMode = "MarkdownV2"
	_, err = bot.Send(tm)
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}
	log.Debug("message sent")
	return nil
}
