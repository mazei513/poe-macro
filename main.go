package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
	var oldClip []byte
	ctx := context.Background()

	fmt.Println("F5 for /hidout\nL_CTRL+D to pause")
	for e := range k.Read() {
		out, err := readClip(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !bytes.HasPrefix(out, itemPrefix) {
			oldClip = out
		}

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
				if has(pressed, "L_CTRL", "C") {
					time.Sleep(50 * time.Millisecond)
					out, err := readClip(ctx)
					if err != nil {
						fmt.Println(err)
						continue
					}

					if bytes.HasPrefix(out, itemPrefix) {
						err = writeClip(ctx, oldClip)
						if err != nil {
							fmt.Println(err)
							continue
						}

						i, err := parse(out)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Printf("Class: %s\nRarity: %s\nName: %s", i.class, i.rarity, i.name)
						if i.itemName != "" {
							fmt.Printf(" (%s)", i.itemName)
						}
						fmt.Println()
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
		if _, ok := pressed[strings.ToUpper(a)]; !ok {
			return false
		}
	}
	return true
}

func readClip(ctx context.Context) ([]byte, error) {
	readCmd := exec.CommandContext(ctx, "xclip", "-o", "-sel", "clip")
	return readCmd.Output()
}

func writeClip(ctx context.Context, v []byte) error {
	writeCmd := exec.CommandContext(ctx, "xclip", "-sel", "clip")
	in, err := writeCmd.StdinPipe()
	if err != nil {
		return err
	}
	err = writeCmd.Start()
	if err != nil {
		return err
	}
	_, err = in.Write(v)
	if err != nil {
		return err
	}
	err = in.Close()
	if err != nil {
		return err
	}
	return writeCmd.Wait()
}

var (
	itemPrefix   = []byte("Item Class: ")
	rarityPrefix = []byte("Rarity: ")
)

type item struct {
	class    string
	rarity   string
	name     string
	itemName string
}

func parse(b []byte) (item, error) {
	class, b, err := next(b, itemPrefix)
	if err != nil {
		return item{}, err
	}

	rarity, b, err := next(b, rarityPrefix)
	if err != nil {
		return item{}, err
	}

	name, b, err := next(b, []byte{})
	if err != nil {
		return item{}, err
	}

	var itemName string
	if rarity == "Rare" || rarity == "Unique" {
		itemName, _, err = next(b, []byte{})
		if err != nil {
			return item{}, err
		}
	}

	return item{
		class:    class,
		rarity:   rarity,
		name:     name,
		itemName: itemName,
	}, nil
}

func next(b, pre []byte) (string, []byte, error) {
	ba := b[len(pre):]
	idx := bytes.Index(ba, []byte("\n"))
	if idx < 0 {
		return "", b, errors.New("invalid structure")
	}
	return string(ba[:idx]), ba[idx+1:], nil
}
