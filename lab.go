package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type WriteCounter struct {
	Total int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	done := make(chan int64)
	n := len(p)
	wc.Total += int64(n)
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				print(done)
			}
		}
	}()
	done <- wc.Total
	return n, nil
}

func print(done chan int64) {
	fmt.Println("\rDownloading...", <-done/1024, "kilobytes complete ")
}

func main() {

	var fileUrl string
	fmt.Printf("Enter your link: ")
	fmt.Scanf("%s", &fileUrl)

	//fileUrl := "http://www.tsu.ru/upload/iblock/418/kalendarnyy-grafik-2018.pdf" //(мой пример)

	fmt.Println("Download Started\n")

	err := DownloadFile(fileUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("Download Finished")
}

func DownloadFile(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path.Base(resp.Request.URL.String()))
	if err != nil {
		return err
	}

	counter := &WriteCounter{}

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return err
	}
	fmt.Print("\n")

	out.Close()

	return nil
}
