package api

import (
	"fmt"
	"io"
	"log"
	"main/main/rfid"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MPTAPI struct {
	options      MPTAPIOptions
	Operation    chan struct{}
	Done         chan struct{}
	FileUploaded chan string
	clients      []chan SSEEvent
	mutex        sync.Mutex
}

type MPTAPIOptions struct {
	FilesDirectory        string
	PlayableFileDirectory string
	Handlers              []Handler
}

type PlayableFile struct {
	Id   string `json:"id"`
	Size int64  `json:"size"`
}

type SSEEvent struct {
	Topic string
	Data  string
}

func NewMPTAPI(options MPTAPIOptions) *MPTAPI {
	return &MPTAPI{options: options, Operation: make(chan struct{}), Done: make(chan struct{}), clients: []chan SSEEvent{}, mutex: sync.Mutex{}, FileUploaded: make(chan string)}
}

func (api *MPTAPI) Start() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowWildcard = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "DELETE", "PATCH"}
	router.Use(cors.New(config))

	router.GET("/files", api.getFileList)
	router.GET("/files/:id", api.downloadFile)
	router.POST("/files/:id", api.scanCard)
	router.DELETE("/files/:id", api.deleteFile)
	router.POST("/files", api.uploadFile)
	router.GET("/sse", api.initializeSSE)

	for _, v := range api.options.Handlers {
		router.Any(v.Endpoint(), v.Action)
	}

	router.Run(":7000")
}

func (api *MPTAPI) initializeSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()
	api.mutex.Lock()
	msgChannel := make(chan SSEEvent)
	api.clients = append(api.clients, msgChannel)
	api.mutex.Unlock()

	c.Stream(func(w io.Writer) bool {
		select {
		case m := <-msgChannel:
			c.SSEvent(m.Topic, m.Data)

		case <-c.Request.Context().Done():
			idx := -1
			for i, v := range api.clients {
				if v == msgChannel {
					idx = i
				}
			}
			if idx > -1 {
				api.mutex.Lock()
				api.clients = slices.Delete(api.clients, idx, 1)
				api.mutex.Unlock()
			}

			return false
		}

		return true
	})

}

func (api *MPTAPI) SendSSEMessage(topic string, msg string) {
	for _, v := range api.clients {
		v <- SSEEvent{Topic: topic, Data: msg}
	}
}

func (api *MPTAPI) getFileList(c *gin.Context) {
	files, err := os.ReadDir(api.options.FilesDirectory)

	if err != nil {
		log.Fatal(err)
	}

	objs := []PlayableFile{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		stat, _ := os.Stat(fmt.Sprintf("%s/%s", api.options.FilesDirectory, file.Name()))
		objs = append(objs, PlayableFile{Id: file.Name(), Size: stat.Size()})
	}

	c.JSON(200, objs)
}

func (api *MPTAPI) uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dst := filepath.Join(api.options.FilesDirectory, filepath.Base(file.Filename))
	c.SaveUploadedFile(file, dst)
	api.FileUploaded <- dst
}

func (api *MPTAPI) scanCard(c *gin.Context) {
	id := c.Param("id")
	api.Operation <- struct{}{}
	scanner := rfid.NewRFIDReader()
	success := scanner.WriteCard(c, id)
	api.Done <- struct{}{}

	if success != nil {
		c.Status(400)
	} else {
		c.Status(200)
	}
}

func (api *MPTAPI) downloadFile(c *gin.Context) {
	id := c.Param("id")
	path := fmt.Sprintf("%s/%s", api.options.FilesDirectory, id)

	c.File(path)
}

func (api *MPTAPI) deleteFile(c *gin.Context) {
	id := c.Param("id")
	path := fmt.Sprintf("%s/%s", api.options.FilesDirectory, id)
	os.Remove(path)
	path = fmt.Sprintf("%s/%s", api.options.PlayableFileDirectory, id)
	os.Remove(path)

	c.Status(200)
}
