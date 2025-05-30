package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

var (
	streamCmd   *exec.Cmd
	streamMutex sync.Mutex
	isStreaming bool
	token       string
	url         string
)

func main() {
	token, url = loadEnv()
	ctx := context.Background()

	b, err := bot.New(token)
	if err != nil {
		panic(err)
	}

	// Register individual handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypeExact, stopHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact, statusHandler)

	println("INFO: Bot started...")
	b.Start(ctx)
}

func loadEnv() (string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		println("ERROR: Could not load .env file")
		panic(err)
	}

	return os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("YOUTUBE_STREAMING_URL")
}

func startStream() error {
	streamMutex.Lock()
	defer streamMutex.Unlock()

	if isStreaming {
		return fmt.Errorf("stream is already running")
	}

	var stdout, stderr io.ReadCloser
	var err error
	streamCmd, stdout, stderr, err = runFfmpeg(url)
	if err != nil {
		return err
	}

	isStreaming = true

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			println("FFMPEG OUT:", scanner.Text())
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			println("FFMPEG ERR:", scanner.Text())
		}
	}()

	// Wait for command to finish
	go func() {
		streamCmd.Wait()
		streamMutex.Lock()
		isStreaming = false
		streamMutex.Unlock()
		println("INFO: Stream ended")
	}()

	return nil
}

func stopStream() error {
	streamMutex.Lock()
	defer streamMutex.Unlock()

	if !isStreaming || streamCmd == nil {
		return fmt.Errorf("no stream is running")
	}

	if err := streamCmd.Process.Kill(); err != nil {
		return fmt.Errorf("error stopping stream: %v", err)
	}

	isStreaming = false
	return nil
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var text string
	if err := startStream(); err != nil {
		text = fmt.Sprintf("Failed to start stream: %v", err)
	} else {
		text = "🔴 Stream started!"
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		println("Error sending message: %v", err)
	}
}

func stopHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var text string
	if err := stopStream(); err != nil {
		text = fmt.Sprintf("Failed to stop stream: %v", err)
	} else {
		text = "⏹️ Stream stopped!"
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		println("Error sending message: %v", err)
	}
}

func statusHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	streamMutex.Lock()
	var text string
	if isStreaming {
		text = "🔴 Stream is running"
	} else {
		text = "⏹️ Stream is not running"
	}
	streamMutex.Unlock()

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		println("Error sending message: %v", err)
	}
}
