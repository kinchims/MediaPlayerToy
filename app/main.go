package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-gst/go-gst/gst"
)

func main() {
	gst.Init(nil)
	sigs := make(chan os.Signal, 1)
	if len(os.Args) > 1 {
		args := os.Args[1:]
		mode := args[0]

		log.Printf("using mode %s", mode)
		if mode == "trainer" {
			trainer := NewRfidTrainer()

			trainer.TrainCards()

			os.Exit(0)
		}
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	player := NewVideoPlayer()
	reader := NewRFIDReader()

	inactiveTimer := time.Duration(0)
	go func() {
		for {
			time.Sleep(time.Second)

			if !player.IsPlaying() {
				inactiveTimer += 1 * time.Second
			} else {
				inactiveTimer = 0
			}

			if inactiveTimer > time.Minute*5 {
				cmd := exec.Command("shutdown", "-h", "now")
				cmd.Run()
			}
		}
	}()

	go func() {
		for {
			id := <-reader.CardData

			if id == "" {
				player.Dispose()
			} else {
				fileName := getFileByName(id)
				if fileName == "" {
					player.Dispose()
				} else if player.currentFile != fileName {
					player.Dispose()
					player = NewVideoPlayerPreloaded(fileName)
					player.Play()
				} else {
				}
			}
		}
	}()

	reader.Start()

	<-sigs //exit

	reader.Dispose()
	player.Dispose()
	gst.Deinit()
}

func getFileByName(fileId string) string {
	directory := "/config/data"
	files, err := os.ReadDir(directory)

	if err != nil {
		log.Fatal(err.Error())
	}

	for i, v := range files {
		if strings.Contains(strings.TrimSuffix(v.Name(), ".mp4"), fileId) {
			return fmt.Sprintf("%s/%s", directory, files[i].Name())
		}
	}

	return ""
}
