package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	url := "https://takuya-kimura.jp/poster/"

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	thumbnails := doc.Find("ul.thumbnail")
	if thumbnails.Children().Length() <= 0 {
		return errors.New("thumbnails are not found")
	}

	thumbnails.Children().Each(func(i int, s *goquery.Selection) {
		id := s.Find("dl").AttrOr("data-id", "")
		prefectures := s.Find("dt > span").Text()
		transportation := strings.TrimLeft(s.Find("dt").Text(), prefectures)
		url := imageURL(id)

		fmt.Printf("\"%s\", \"%s\", \"%s\", \"%s\"\n", id, prefectures, transportation, url)

		downloadImage(id)
	})

	return nil
}

func downloadImage(id string) error {
	resp, err := http.Get(imageURL(id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fp, err := os.Create(fmt.Sprintf("p%s.jpg", id))
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, resp.Body)
	return err
}

func imageURL(id string) string {
	return fmt.Sprintf("https://takuya-kimura.jp/poster/assets/images/poster/p%s.jpg", id)
}
