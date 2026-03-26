package video

import (
	"fmt"
	"log"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

type VideoPlayer struct {
	pipeline    *gst.Pipeline
	loop        *glib.MainLoop
	CurrentFile string
}

func NewVideoPlayer() *VideoPlayer {
	return NewVideoPlayerPreloaded("")
}

func NewVideoPlayerPreloaded(location string) *VideoPlayer {
	pipeline, err := gst.NewPipelineFromString(fmt.Sprintf("filesrc location=%s name=src ! decodebin name=decode decode. ! queue ! audioconvert ! audioresample ! alsasink device=hw:CARD=wm8960soundcard,DEV=0 decode. ! queue ! videoconvert ! kmssink connector-id=33", location))
	if err != nil {
		log.Fatal(err.Error())
	}

	player := &VideoPlayer{
		pipeline:    pipeline,
		CurrentFile: location,
		loop:        glib.NewMainLoop(glib.MainContextDefault(), false),
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
	player.pipeline.BlockSetState(gst.StatePaused)
}

func (player *VideoPlayer) Play() {
	player.pipeline.SetState(gst.StatePlaying)
}

func (player *VideoPlayer) Reset() {
	player.pipeline.BlockSetState(gst.StateNull)
	player.pipeline.SeekSimple(0, gst.FormatTime, gst.SeekFlagFlush|gst.SeekFlagAccurate)
}
func (player *VideoPlayer) Dispose() {
	player.loop.Quit()
	player.CurrentFile = ""
	player.pipeline.BlockSetState(gst.StateNull)
	player.pipeline.SendEvent(gst.NewFlushStartEvent())
}
