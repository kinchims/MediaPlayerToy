package main

import (
	"context"
	"main/main/api"
	"main/main/states"
	"main/main/video"
	"os/signal"
	"syscall"

	"github.com/go-gst/go-gst/gst"
)

var (
	state         states.State
	ctx           context.Context
	postProcessor *video.PostProcessor
	SrcFiles      string = "/home/sean/Development/MediaPlayerToy/files2-nc"
	PlayableFiles string = "/home/sean/Development/MediaPlayerToy/files2-c"
)

func main() {
	gst.Init(nil)
	ctx, _ = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	postProcessor := video.NewPostProcessor(ctx, video.VideoOptions{Width: 800, Height: 400, OutputDir: PlayableFiles, Transform: "scale=800:480,vflip,hflip"})
	controller := api.NewMPTAPI(api.MPTAPIOptions{FilesDirectory: SrcFiles, PlayableFileDirectory: PlayableFiles})
	go func() {
		controller.Start()
	}()

	go func() {
		for {
			select {
			case <-controller.Operation:
				SetState(nil)
			case <-controller.Done:
				SetState(states.NewNormalState())
			case path := <-controller.FileUploaded:
				postProcessor.ProcessFile(path)
			case complete := <-postProcessor.ProcessingComplete:
				controller.SendSSEMessage("completed", complete)
			case started := <-postProcessor.ProcessingStarted:
				controller.SendSSEMessage("started", started)
			}

		}
	}()

	SetState(states.NewNormalState())
	<-ctx.Done() //exit

	SetState(nil)
	state.Dispose()
	gst.Deinit()
}

func SetState(next states.State) {
	if state != nil {
		state.Dispose()
	}

	state = next
	if state != nil {
		state.Run(ctx)
	}
}
