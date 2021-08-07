package main

import (
	"fmt"
	"os"

	"github.com/MarinX/keylogger"
)

func main() {
	keyboard := keylogger.FindKeyboardDevice()
	k, err := keylogger.New(keyboard)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer k.Close()

	pause := false
	pressed := ks{}

	fmt.Println("F5 for /hidout\nL_CTRL+D to pause")
	for e := range k.Read() {
		switch e.Type {
		case keylogger.EvKey:
			if e.KeyPress() {
				pressed[e.KeyString()] = struct{}{}
				if !pause && has(pressed, "F5") {
					vs := []string{"enter", "/", "h", "l", "g", "k", ";", "i", "f", "enter"}
					for _, v := range vs {
						k.WriteOnce(v)
					}
				}
				if has(pressed, "L_CTRL", "D") {
					pause = !pause
					fmt.Println("paused:", pause)
				}
			}
			if e.KeyRelease() {
				delete(pressed, e.KeyString())
			}
		}
	}
}

type ks = map[string]struct{}

func has(pressed ks, k ...string) bool {
	if len(pressed) != len(k) {
		return false
	}
	for _, a := range k {
		if _, ok := pressed[a]; !ok {
			return false
		}
	}
	return true
}
