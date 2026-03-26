package main

import (
	"context"
	"encoding/json"
	"main/main/api"
	"main/main/native"
	"main/main/states"
	"main/main/video"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-gst/go-gst/gst"
)

var (
	state         states.State
	ctx           context.Context
	postProcessor *video.PostProcessor
	SrcFiles      string  = "/home/sean/Development/MediaPlayerToy/files2-nc"
	PlayableFiles string  = "/home/sean/Development/MediaPlayerToy/files2-c"
	configPath    string  = "/etc/jukebox/config.txt"
	iface         string  = "wlp9s0" //wlan0
	config        *Config = &Config{}
)

func main() {
	gst.Init(nil)
	ctx, _ = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	loadConfig()
	native.SetVolume(config.Volume)
	native.BroadcastSSID(iface, config.WifiAPName, config.WifiPassword)

	postProcessor := video.NewPostProcessor(ctx, video.VideoOptions{Width: 800, Height: 400, OutputDir: PlayableFiles, Transform: "scale=800:480,vflip,hflip"})
	controller := api.NewMPTAPI(api.MPTAPIOptions{FilesDirectory: SrcFiles, PlayableFileDirectory: PlayableFiles, Handlers: []api.Handler{SettingsHandler{}, PlayFileHandler{}}})
	go func() {
		controller.Start()
	}()

	go func() {
		for {
			select {
			case <-controller.Operation:
				SetState(nil)
			case <-controller.Done:
				SetState(NewNormalState())
			case path := <-controller.FileUploaded:
				postProcessor.ProcessFile(path)
			case complete := <-postProcessor.ProcessingComplete:
				controller.SendSSEMessage("completed", complete)
			case started := <-postProcessor.ProcessingStarted:
				controller.SendSSEMessage("started", started)
			}

		}
	}()

	SetState(NewNormalState())
	<-ctx.Done() //exit
	native.StopBroadcasting()

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

func SetConfig(c Config) {
	if c.WifiAPName != "" || c.WifiPassword != "" {
		updated := false
		if c.WifiAPName != "" {
			if c.WifiAPName != config.WifiAPName {
				updated = true
				config.WifiAPName = c.WifiAPName
			}
		}

		if c.WifiPassword != "" {
			config.WifiPassword = c.WifiPassword
			updated = true
		}

		if updated {
			native.StopBroadcasting()
			native.BroadcastSSID(iface, config.WifiAPName, config.WifiPassword)
		}
	}

	if c.Volume != config.Volume {
		config.Volume = c.Volume
		native.SetVolume(config.Volume)
	}

	saveConfig()
}

func loadConfig() {
	file, err := os.ReadFile(configPath)

	if err != nil {
		config = &Config{
			WifiPassword: "UYGsBvyr",
			WifiAPName:   "jukebox",
			Volume:       50,
		}
	} else {
		config = &Config{}
		json.Unmarshal([]byte(file), config)
	}
}

func saveConfig() {
	data, err := json.Marshal(config)

	if err != nil {
		return
	}

	os.WriteFile(configPath, data, 0655)
}

type Config struct {
	WifiAPName   string `json:"wifiAPName"`
	WifiPassword string `json:"wifiPassword"`
	Volume       int32  `json:"volume"`
}
