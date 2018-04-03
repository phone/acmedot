package main

import (
	"9fans.net/go/acme"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unicode/utf8"
)

func isWordRune(c rune) bool {
	if c <= ' ' {
		return false
	}
	if 0x7F <= c && c <= 0xA0 {
		return false
	}
	for _, r := range "!\"#$%&'()*+,-./:;<=>?@[\\]^`{|}~" {
		if c == r {
			return false
		}
	}
	return true
}

func getWindowBody(win *acme.Win) (string, error) {
	// Stolen from Rog Peppe's godef
	var body []byte
	buf := make([]byte, 8000)
	for {
		n, err := win.Read("body", buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		body = append(body, buf[0:n]...)
	}
	return string(body), nil
}

func main() {
	swinid := os.Getenv("winid")
	if swinid == "" {
		log.Fatal("winid not set")
	}
	winid, err := strconv.Atoi(swinid)
	if err != nil {
		log.Fatalf("Bad winid: %s", winid)
	}
	win, err := acme.Open(winid, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer win.CloseFiles()
	_, _, err = win.ReadAddr()
	if err != nil {
		log.Fatal(err)
	}
	err = win.Ctl("addr=dot")
	if err != nil {
		log.Fatal(err)
	}
	q0, q1, err := win.ReadAddr()
	if err != nil {
		log.Fatal(err)
	}
	var word string
	if q0 == q1 {
		// We just have one point, so try to walk left and right until we hit wordish boundaries.
		body, err := getWindowBody(win)
		if err != nil {
			log.Fatalf("error getting window body: %s\n", err)
		}
		runelen := utf8.RuneCountInString(body)
		for q0 > 0 {
			rune, width := utf8.DecodeRuneInString(body[q0-1:])
			if isWordRune(rune) {
				q0 -= width
			} else {
				break
			}
		}
		for q1 < runelen {
			rune, width := utf8.DecodeRuneInString(body[q1:])
			if isWordRune(rune) {
				q1 += width
			} else {
				break
			}
		}
		word = body[q0:q1]
	} else {
		body, err := getWindowBody(win)
		if err != nil {
			log.Fatal(err)
		}
		word = body[q0:q1]
	}
	fmt.Printf("%s", word)
}
