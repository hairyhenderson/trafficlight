package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	pinBuzzer = "15" // BCM 22, pin 15
	pinGreen  = "16" // BCM 23, pin 16
	pinYellow = "18" // BCM 24, pin 18
	pinRed    = "22" // BCM 25, pin 22
)

type pins struct {
	red, yellow, green *gpio.LedDriver
	buzzer             *gpio.BuzzerDriver
}

func initPins() pins {
	r := raspi.NewAdaptor()
	err := r.Connect()
	if err != nil {
		panic(err)
	}
	p := pins{}
	p.green = gpio.NewLedDriver(r, pinGreen)
	err = p.green.Start()
	if err != nil {
		panic(err)
	}
	p.yellow = gpio.NewLedDriver(r, pinYellow)
	err = p.yellow.Start()
	if err != nil {
		panic(err)
	}
	p.red = gpio.NewLedDriver(r, pinRed)
	err = p.red.Start()
	if err != nil {
		panic(err)
	}
	p.buzzer = gpio.NewBuzzerDriver(r, pinBuzzer)
	err = p.buzzer.Start()
	if err != nil {
		panic(err)
	}
	return p
}

func main() {
	p := initPins()

	ch := make(chan rune)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case r := <-ch:
				switch r {
				case 'r':
					p.red.Toggle()
				case 'y':
					p.yellow.Toggle()
				case 'g':
					p.green.Toggle()
				}
			case <-done:
				fmt.Printf("done\n")
				p.green.Off()
				p.yellow.Off()
				p.red.Off()
				break
			default:
			}
		}
	}()

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		} else if key == keyboard.KeyEsc {
			done <- struct{}{}
			break
		}
		fmt.Printf("%q", char)
		ch <- char
	}
}
