package states

import (
	"context"
	"fmt"
	"log"
	"main/main/rfid"
	"main/main/video"
	"os"
	"os/exec"
	"strings"
	"time"
)

type NormalState struct {
	player *video.VideoPlayer
	reader *rfid.RFIDReader
}

func NewNormalState() *NormalState {
	return &NormalState{}
}

func (state *NormalState) Run(ctx context.Context) {
	state.player = video.NewVideoPlayer()
	state.reader = rfid.NewRFIDReader()

	inactiveTimer := time.Duration(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Minute):
				if !state.player.IsPlaying() {
					inactiveTimer += 1 * time.Second
				} else {
					inactiveTimer = 0
				}

				if inactiveTimer > time.Minute*5 {
					cmd := exec.Command("shutdown", "-h", "now")
					cmd.Run()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case id := <-state.reader.CardData:
				if id == "" {
					state.player.Dispose()
				} else {
					fileName := getFileByName(id)
					if fileName == "" {
						state.player.Dispose()
					} else if state.player.CurrentFile != fileName {
						state.player.Dispose()
						state.player = video.NewVideoPlayerPreloaded(fileName)
						state.player.Play()
					} else {
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	state.reader.Start(ctx)
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

func (state *NormalState) Dispose() {
	state.reader.Dispose()
	state.player.Dispose()
}
