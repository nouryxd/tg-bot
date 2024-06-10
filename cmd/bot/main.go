package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type config struct {
	AllowedUser       string
	UploadURL         string
	WolframAlphaAppID string
}

type application struct {
	Telegram *tele.Bot
	Config   *config
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &config{
		AllowedUser:       os.Getenv("ALLOWED_USER"),
		UploadURL:         os.Getenv("UPLOAD_URL"),
		WolframAlphaAppID: os.Getenv("WOLFRAMALPHA_APP_ID"),
	}

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Use(middleware.Logger())

	app := &application{
		Telegram: b,
		Config:   cfg,
	}

	app.Telegram.Handle(tele.OnText, func(c tele.Context) error {
		if c.Sender().Username != app.Config.AllowedUser {
			return fmt.Errorf("Unauthorized user. User: %v", c.Sender().Username)
		}

		if err := app.handleText(); err != nil {
			return err
		}

		return nil
	})

	app.Telegram.Handle(tele.OnVideo, func(c tele.Context) error {
		if c.Sender().Username != app.Config.AllowedUser {
			return fmt.Errorf("Unauthorized user. User: %v", c.Sender().Username)
		}

		if err := app.handleVideo(c); err != nil {
			return err
		}

		return nil
	})

	app.Telegram.Handle(tele.OnPhoto, func(c tele.Context) error {
		if c.Sender().Username != app.Config.AllowedUser {
			return fmt.Errorf("Unauthorized user. User: %v", c.Sender().Username)
		}

		if err := app.handlePhoto(c); err != nil {
			return err
		}

		return nil
	})

	app.Telegram.Handle(tele.OnDocument, func(c tele.Context) error {
		if c.Sender().Username != app.Config.AllowedUser {
			return fmt.Errorf("Unauthorized user. User: %v", c.Sender().Username)
		}

		if err := app.handleDocument(c); err != nil {
			return err
		}

		return nil
	})

	app.Telegram.Handle(tele.OnAnimation, func(c tele.Context) error {
		if c.Sender().Username != app.Config.AllowedUser {
			return fmt.Errorf("Unauthorized user. User: %v", c.Sender().Username)
		}

		if err := app.handleAnimation(c); err != nil {
			return err
		}

		return nil
	})

	app.Telegram.Start()
}
