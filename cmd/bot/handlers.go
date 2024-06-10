package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/nouryxd/tg-bot/pkg/commands"
	tele "gopkg.in/telebot.v3"
)

func (app *application) handleText() error {

	app.Telegram.Handle("ping", func(c tele.Context) error {
		return c.Send("Pong!")
	})

	app.Telegram.Handle("pong", func(c tele.Context) error {
		return c.Send("Ping!")
	})

	app.Telegram.Handle("/currency", func(c tele.Context) error {
		tags := c.Args()

		reply, err := commands.Currency(tags[0], tags[1], tags[3])
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/wa", func(c tele.Context) error {
		tags := c.Args()

		var query string
		for _, tag := range tags {
			query += tag + " "
		}

		reply, err := commands.WolframAlphaConv(query, app.Config.WolframAlphaAppID)
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/wolframalpha", func(c tele.Context) error {
		tags := c.Args()

		var query string
		for _, tag := range tags {
			query += tag + " "
		}

		reply, err := commands.WolframAlphaConv(query, app.Config.WolframAlphaAppID)
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/ytdl", func(c tele.Context) error {
		if c.Sender().Username != app.Config.AllowedUser {
			return fmt.Errorf("Unauthorized user. User: %v", c.Sender().Username)
		}

		tags := c.Args()

		reply, err := app.YtdlDownload(tags[0])
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/weather", func(c tele.Context) error {
		var location string

		for _, str := range c.Args() {
			location += str + " "
		}

		reply, err := commands.Weather(location)
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/xkcd", func(c tele.Context) error {

		reply, err := commands.Xkcd()
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/coinflip", func(c tele.Context) error {
		return c.Send(commands.Coinflip())
	})

	app.Telegram.Handle("/rxkcd", func(c tele.Context) error {
		reply, err := commands.RandomXkcd()
		if err != nil {
			return err
		}

		return c.Send(reply)
	})

	app.Telegram.Handle("/tags", func(c tele.Context) error {
		tags := c.Args()

		var reply string
		for _, tag := range tags {
			reply += tag + " "
		}

		return c.Send(reply)
	})

	return nil
}

func (app *application) handleVideo(c tele.Context) error {
	f := c.Message().Video
	path := fmt.Sprintf("%s", f.FileID)

	c.Send("Downloading video...")
	app.Telegram.Download(&f.File, path)

	c.Send("Uploading video...")
	link, err := app.upload(path)
	if err != nil {
		c.Send("Something went wrong")
		return fmt.Errorf("error during upload: %v", err)
	}

	return c.Send(link)
}

func (app *application) handlePhoto(c tele.Context) error {
	f := c.Message().Photo
	path := fmt.Sprintf("%s", f.FileID)

	app.Telegram.Download(&f.File, path)
	link, err := app.upload(path)
	if err != nil {
		c.Send("Something went wrong")
		return fmt.Errorf("error during upload: %v", err)
	}

	return c.Send(link)
}

func (app *application) handleDocument(c tele.Context) error {
	f := c.Message().Document
	path := fmt.Sprintf("%s", f.FileID)

	c.Send("Downloading...")
	app.Telegram.Download(&f.File, path)

	c.Send("Uploading...")
	link, err := app.upload(path)
	if err != nil {
		c.Send("Something went wrong")
		return fmt.Errorf("error during upload: %v", err)
	}

	return c.Send(link)
}

func (app *application) handleAnimation(c tele.Context) error {
	f := c.Message().Animation
	path := fmt.Sprintf("%s", f.FileID)

	c.Send("Downloading...")
	app.Telegram.Download(&f.File, path)

	c.Send("Uploading...")
	link, err := app.upload(path)
	if err != nil {
		c.Send("Something went wrong")
		return fmt.Errorf("error during upload: %v", err)
	}

	return c.Send(link)
}

func (app *application) upload(path string) (string, error) {
	defer os.Remove(path)
	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {

		defer pw.Close()

		err := form.WriteField("name", "xd")
		if err != nil {
			os.Remove(path)
			return
		}

		file, err := os.Open(path) // path to image file
		if err != nil {
			os.Remove(path)
			return
		}

		w, err := form.CreateFormFile("file", path)
		if err != nil {
			os.Remove(path)
			return
		}

		_, err = io.Copy(w, file)
		if err != nil {
			os.Remove(path)
			return
		}

		form.Close()
	}()

	req, err := http.NewRequest(http.MethodPost, app.Config.UploadURL, pr)
	if err != nil {
		return "Something went wrong", err
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	httpClient := http.Client{
		Timeout: 300 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "Something went wrong", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Something went wrong", err
	}

	var reply = string(body[:])

	return reply, nil
}
