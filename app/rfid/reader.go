package rfid

import (
	"context"
	"errors"
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
		rfid.SetAntennaGain(5)
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
	if reader.rfid != nil {
		reader.rfid.Halt()
	}

	if reader.port != nil {
		reader.port.Close()
	}
}

func (reader *RFIDReader) ReadCard(context context.Context) string {
	data, err := reader.rfid.ReadCard(0, commands.PICC_AUTHENT1B, 2, 0, mfrc522.DefaultKey)
	if err == nil {
		return string(data[:clen(data)])
	} else {
		str := err.Error()

		if strings.Contains(str, "timeout waiting for IRQ edge") {
			return reader.ReadCard(context)
		}

		return ""
	}
}

func (reader *RFIDReader) WriteCard(ctx context.Context, data string) error {
	now := time.Now()
	c, cf := context.WithDeadline(ctx, now.Add(time.Minute*2))
	defer cf()

	for {
		select {
		case <-c.Done():
			if c.Err() == context.DeadlineExceeded {
				return errors.New("request timeout")
			}

			return errors.New("request cancelled")
		default:
			err := reader.rfid.WriteCard(time.Second*2, byte(commands.PICC_AUTHENT1B), 2, 0, stringIntoByte16(data), mfrc522.DefaultKey)
			if err != nil {
				return err
			}
			fmt.Println("Write successful")

			return nil
		}
	}
}

func (reader *RFIDReader) Start(context context.Context) {
	arch := runtime.GOARCH
	if strings.Contains(arch, "arm") {
		go func() {
			for {
				if context.Err() != nil {
					return
				}

				data, err := reader.rfid.ReadCard(time.Second*2, commands.PICC_AUTHENT1B, 2, 0, mfrc522.DefaultKey)
				if err == nil {
					str := string(data[:clen(data)])

					select {
					case reader.CardData <- str:
					default:
					}
				} else {
					str := err.Error()

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
