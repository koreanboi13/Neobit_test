package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	TestCmd  = "/test"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from chat %d", text, chatID)

	switch text {
	case TestCmd:
		go p.testSpeed(context.Background(), chatID)
		return nil
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) testSpeed(ctx context.Context, chatID int) error {
	if err := p.tg.SendMessage(ctx, chatID, msgSpeedTestStart); err != nil {
		return fmt.Errorf("failed to send starting message: %w", err)
	}

	result, err := p.speedtestClient.SpeedTest(int64(chatID))
	if err != nil {
		log.Printf("speed test failed: %v", err)
		return p.tg.SendMessage(ctx, chatID, msgSpeedTestFailed)
	}

	responseText := fmt.Sprintf(
		"Speed Test Results:\n\nUpload: %.2f Мбит/с\nDownload: %.2f Мбит/с",
		result.UploadSpeedMbps,
		result.DownloadSpeedMbps,
	)

	return p.tg.SendMessage(ctx, chatID, responseText)
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}
