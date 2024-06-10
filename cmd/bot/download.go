package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/wader/goutubedl"
)

func (app *application) YtdlDownload(link string) (string, error) {
	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(context.Background(), link, goutubedl.Options{})
	if err != nil {
		return "", err
	}

	// For some reason youtube links return webm as result.Info.Ext but
	// are in reality mp4.
	var rExt string
	if strings.HasPrefix(link, "https://www.youtube.com/") || strings.HasPrefix(link, "https://youtu.be/") {
		rExt = "mp4"
	} else {
		rExt = result.Info.Ext
	}

	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		return "", err
	}

	identifier := uuid.NewString()
	fileName := fmt.Sprintf("./public/%s.%s", identifier, rExt)
	f, err := os.Create(fileName)

	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = io.Copy(f, downloadResult); err != nil {
		return "", err
	}

	downloadResult.Close()
	f.Close()

	// Upload
	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		err := form.WriteField("name", "xd")
		if err != nil {
			os.Remove(fileName)
			return
		}

		file, err := os.Open(fileName) // path to image file
		if err != nil {
			os.Remove(fileName)
			return
		}

		w, err := form.CreateFormFile("file", fileName)
		if err != nil {
			os.Remove(fileName)
			return
		}

		_, err = io.Copy(w, file)
		if err != nil {
			os.Remove(fileName)
			return
		}

		form.Close()
	}()

	req, err := http.NewRequest(http.MethodPost, app.Config.UploadURL, pr)
	if err != nil {
		os.Remove(fileName)
		return "", err
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	httpClient := http.Client{
		Timeout: 300 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		os.Remove(fileName)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		os.Remove(fileName)
		return "", err
	}

	var reply = string(body[:])

	os.Remove(fileName)
	return reply, nil
}
