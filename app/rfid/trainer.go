package rfid

import (
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"slices"
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

type RFIDTrainer struct {
	rfid *mfrc522.Dev
	port spi.PortCloser
}

func NewRfidTrainer() *RFIDTrainer {
	trainer := &RFIDTrainer{}
	arch := runtime.GOARCH
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
		trainer.port = port

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
		trainer.rfid = rfid

		fmt.Println("Started rfid reader.")

		runtime.AddCleanup(trainer, func(rfid *mfrc522.Dev) {
			rfid.Halt()
		}, trainer.rfid)
		runtime.AddCleanup(trainer, func(port spi.PortCloser) {
			port.Close()
		}, trainer.port)
	}

	return trainer
}

func (trainer *RFIDTrainer) TrainCards() {
	trainedCards := make([]string, 0)
	for {
		log.Printf("Card %d of 10. Scan card to train.", len(trainedCards)+1)

		data, err := trainer.rfid.ReadUID(time.Hour)

		if err != nil {
			continue
		} else {
			id := hex.EncodeToString(data)
			found := slices.Contains(trainedCards, id)

			if found {
				log.Print("This card has already been trained. Try again.")
				time.Sleep(time.Second * 2)
				continue
			}

			err = trainer.rfid.WriteCard(5*time.Second, byte(commands.PICC_AUTHENT1B), 2, 0, stringIntoByte16(fmt.Sprintf("%d", len(trainedCards)+1)), mfrc522.DefaultKey)
			if err != nil {
				continue
			} else {
				fmt.Println("Write successful")
				trainedCards = append(trainedCards, id)
			}

			time.Sleep(time.Second * 2)
			if len(trainedCards) == 10 {
				break
			}
		}
	}

	log.Printf("Training complete!")
}

func stringIntoByte16(str string) [16]byte {
	var data [16]byte
	copy(data[:], str) // copy already checks length of str
	return data
}
