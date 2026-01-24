package main

import (
	"log"
	"sync"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

type VideoPlayer struct {
	pipeline    *gst.Pipeline
	loop        *glib.MainLoop
	mu          sync.Mutex
	currentFile string
}

func NewVideoPlayer() *VideoPlayer {
	return NewVideoPlayerPreloaded("")
}

func NewVideoPlayerPreloaded(location string) *VideoPlayer {
	pipeline, err := gst.NewPipelineFromString("filesrc name=src ! decodebin name=decode decode. ! queue ! audioconvert ! audioresample ! autoaudiosink decode. ! queue ! videoconvert ! videoscale ! video/x-raw,width=800,height=480 ! kmssink")
	if err != nil {
		log.Fatal(err.Error())
	}

	player := &VideoPlayer{
		pipeline:    pipeline,
		currentFile: location,
		loop:        glib.NewMainLoop(glib.MainContextDefault(), false),
		mu:          sync.Mutex{},
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		player.pipeline.GetBus().AddWatch(func(msg *gst.Message) bool {
			switch msg.Type() {
			case gst.MessageEOS:
				player.pipeline.BlockSetState(gst.StatePaused)
				player.pipeline.SeekSimple(0, gst.FormatTime, gst.SeekFlagFlush|gst.SeekFlagAccurate)
				player.pipeline.BlockSetState(gst.StatePlaying)
			}

			return true
		})

		player.loop.RunError()
	}()

	return player
}

func (player *VideoPlayer) IsPlaying() bool {
	return player.loop.IsRunning()
}

func (player *VideoPlayer) Pause() {
	player.mu.Lock()
	defer player.mu.Unlock()
	player.pipeline.BlockSetState(gst.StatePaused)
}

func (player *VideoPlayer) Stop() {
	player.mu.Lock()
	defer player.mu.Unlock()
	player.pipeline.BlockSetState(gst.StateNull)
	player.pipeline.SendEvent(gst.NewFlushStartEvent())
	src, err := player.pipeline.GetElementByName("src")

	if err != nil {
		log.Fatal(err.Error())
	}

	src.SetProperty("location", "")
}

func (player *VideoPlayer) Play() {
	player.pipeline.BlockSetState(gst.StatePlaying)
}

func (player *VideoPlayer) Reset() {
	player.pipeline.BlockSetState(gst.StateNull)
	player.pipeline.SeekSimple(0, gst.FormatTime, gst.SeekFlagFlush|gst.SeekFlagAccurate)
}
func (player *VideoPlayer) Dispose() {
	player.mu.Lock()
	defer player.mu.Unlock()
	player.loop.Quit()
	player.currentFile = ""
	player.pipeline.BlockSetState(gst.StateNull)
	player.pipeline.SendEvent(gst.NewFlushStartEvent())
}

func (player *VideoPlayer) PlayFile(file string) {
	player.mu.Lock()
	defer player.mu.Unlock()
	src, err := player.pipeline.GetElementByName("src")

	if err != nil {
		log.Fatal(err.Error())
	}

	currentFile, err := src.GetProperty("location")
	if err == nil {
		if currentFile == file {
			return
		}
	}
	player.currentFile = file
	player.pipeline.BlockSetState(gst.StateNull)
	player.pipeline.SendEvent(gst.NewFlushStartEvent())
	player.pipeline.SendEvent(gst.NewFlushStopEvent(true))
	src.SetProperty("location", file)
	src.SyncStateWithParent()
	gst.NewFlushStartEvent()
	player.pipeline.SendEvent(gst.NewReconfigureEvent())
	player.pipeline.BlockSetState(gst.StatePlaying)
}
