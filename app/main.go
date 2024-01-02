package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/muasx88/github-webhook-telegram-bot/app/lib/github"
	"github.com/muasx88/github-webhook-telegram-bot/app/lib/telegram"
	m "github.com/muasx88/github-webhook-telegram-bot/app/middleware"
	"github.com/muasx88/github-webhook-telegram-bot/app/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var (
	PORT = "8000"

	GithubSecret     string
	TelegramBotToken string
	TelegramChatId   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GithubSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	TelegramChatId = os.Getenv("TELEGRAM_BOT_CHAT_ID")

	if GithubSecret == "" || TelegramBotToken == "" || TelegramChatId == "" {
		log.Fatal("GitHub Secret, Telegram Bot Token, and Telegram Chat ID must be filled")
	}
}

func main() {

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{"GET", "POST"},
	}))
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})
	e.POST("/github/notification", handleGithubWebhook, m.VerifyGithubSecret)

	// Start server
	go func() {
		if err := e.Start(":" + PORT); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func handleGithubWebhook(c echo.Context) error {
	var payload github.Webhook
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	var message string
	message = fmt.Sprintf("New push in repository <a href='%s'>%s</a>", payload.Repository.HtmlUrl, utils.EscapeHTML(payload.Repository.FullName))

	if payload.Sender != nil {
		message += fmt.Sprintf(" from <a href='%s'>%s</a>", payload.Sender.HtmlUrl, utils.EscapeHTML(payload.Sender.Login))
	}

	if payload.Commits != nil && len(payload.Commits) > 0 {
		message += "\n<b>Commits:</>"
		for _, commit := range payload.Commits {
			message += "\n<i>Pushed By: </i>" + utils.EscapeHTML(commit.Committer.Username)
			message += fmt.Sprintf("\nUrl: <a href='%s'>%s</a>", commit.Url, utils.EscapeHTML(commit.Id))
			message += "\n<i>Message:</i> " + commit.Message

			if len(commit.Added) > 0 {
				message += "\n<i>Added: </i> " + utils.EscapeHTML(strings.Join(commit.Added, ", "))
			}

			if len(commit.Removed) > 0 {
				message += "\n<i>Removed: </i> " + utils.EscapeHTML(strings.Join(commit.Removed, ", "))
			}

			if len(commit.Modified) > 0 {
				message += "\n<i>Modified: </i> " + utils.EscapeHTML(strings.Join(commit.Modified, ", "))
			}
		}
	}

	bot := telegram.NewTelegramBot(TelegramBotToken)
	bot.SendMessage(TelegramChatId, message, "HTML")

	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
