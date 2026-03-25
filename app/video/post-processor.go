package video

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"slices"
	"sync"
)

type PostProcessor struct {
	options             VideoOptions
	queue               []string
	currentlyProcessing string
	mutex               sync.Mutex
	processingComplete  chan string
	ProcessingStarted   chan string
	ProcessingComplete  chan string
	ctx                 context.Context
}

type VideoOptions struct {
	Width     int32
	Height    int32
	Transform string
	OutputDir string
}

func NewPostProcessor(c context.Context, options VideoOptions) *PostProcessor {
	p := &PostProcessor{options: options, queue: []string{}, mutex: sync.Mutex{}, processingComplete: make(chan string), ctx: c, ProcessingStarted: make(chan string), ProcessingComplete: make(chan string)}

	go func() {
		for {
			select {
			case <-p.processingComplete:
				if len(p.queue) > 0 {
					next := p.queue[0]
					p.mutex.Lock()
					p.currentlyProcessing = next
					p.queue = slices.Delete(p.queue, 0, 1)
					p.mutex.Unlock()
					p.ProcessingStarted <- next
					p.processFile(p.ctx, next)
				} else {
					p.mutex.Lock()
					p.currentlyProcessing = ""
					p.mutex.Unlock()
				}
			}
		}
	}()

	return p
}

func (processor *PostProcessor) ProcessFile(filePath string) {
	processor.mutex.Lock()
	defer processor.mutex.Unlock()

	if processor.currentlyProcessing != "" || len(processor.queue) > 0 {
		processor.queue = append(processor.queue, filePath)
		return
	}

	processor.processFile(processor.ctx, filePath)
}

func (processor *PostProcessor) processFile(c context.Context, filePath string) {
	go func() {
		fileName := path.Base(filePath)
		outputPath := fmt.Sprintf("%s/%s", processor.options.OutputDir, fileName)
		processor.ProcessingStarted <- fileName
		cmd := exec.CommandContext(c, "ffmpeg", "-i", filePath, "-vf", fmt.Sprintf("%s", processor.options.Transform), outputPath)
		cmd.Run()
		go func() {
			processor.processingComplete <- filePath
		}()

		go func() {
			processor.ProcessingComplete <- fileName
		}()
	}()
}
