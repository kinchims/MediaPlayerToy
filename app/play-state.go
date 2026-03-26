package main

import (
	"context"
	"main/main/rfid"
	"main/main/video"
)

type PlayState struct {
	player *video.VideoPlayer
	reader *rfid.RFIDReader
	file   string
}

func NewPlayState(file string) *PlayState {
	return &PlayState{file: file}
}

func (state *PlayState) Run(ctx context.Context) {
	state.player = video.NewVideoPlayerPreloaded(state.file)
	state.reader = rfid.NewRFIDReader()

	state.player.Play()

	go func() {
		for {
			select {
			case <-state.reader.CardData:
				SetState(NewNormalState())
			case <-ctx.Done():
				return
			}
		}
	}()

	state.reader.Start(ctx)
}

func (state *PlayState) Dispose() {
	state.player.Dispose()
	state.reader.Dispose()
}
