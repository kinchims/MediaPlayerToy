package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/devices/v3/mfrc522/commands"
	"periph.io/x/host/v3"
)

type RFIDReader struct {
	rfid     *mfrc522.Dev
	port     spi.PortCloser
	CardData chan string
}

const ()

func NewRFIDReader() *RFIDReader {
	arch := runtime.GOARCH
	reader := &RFIDReader{CardData: make(chan string)}
	if strings.Contains(arch, "arm") {
		resetPin := "P1_22" // GPIO 25
		irqPin := "P1_18"   // GPIO 24

		if _, err := host.Init(); err != nil {
			log.Fatal(err)
		}

		port, err := spireg.Open("")
		if err != nil {
			log.Fatal(err)
		}
		reader.port = port

		var gpioResetPin gpio.PinOut = gpioreg.ByName(resetPin)
		if gpioResetPin == nil {
			log.Fatalf("Failed to find %v", resetPin)
		}

		var gpioIRQPin gpio.PinIn = gpioreg.ByName(irqPin)
		if gpioIRQPin == nil {
			log.Fatalf("Failed to find %v", irqPin)
		}

		rfid, err := mfrc522.NewSPI(port, gpioResetPin, gpioIRQPin, mfrc522.WithSync())
		if err != nil {
			log.Fatal(err)
		}

		// setting the antenna signal strength, signal strength from 0 to 7
		rfid.SetAntennaGain(6)
		reader.rfid = rfid

		fmt.Println("Started rfid reader.")

		runtime.AddCleanup(reader, func(rfid *mfrc522.Dev) {
			rfid.Halt()
		}, reader.rfid)
		runtime.AddCleanup(reader, func(port spi.PortCloser) {
			port.Close()
		}, reader.port)
	}

	return reader
}

func (reader *RFIDReader) Dispose() {
	reader.rfid.Halt()
	reader.port.Close()
}

func (reader *RFIDReader) Start() {
	arch := runtime.GOARCH
	if strings.Contains(arch, "arm") {
		go func() {
			for {
				data, err := reader.rfid.ReadCard(time.Second*2, commands.PICC_AUTHENT1B, 2, 0, mfrc522.DefaultKey)
				if err == nil {
					str := string(data[:clen(data)])

					select {
					case reader.CardData <- str:
					default:
					}
				} else {
					str := err.Error()
					log.Print(str)

					if strings.Contains(str, "timeout waiting for IRQ edge") {
						select {
						case reader.CardData <- "":
						default:
						}
					}
				}
			}
		}()
	}
}

// find first null byte
func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
