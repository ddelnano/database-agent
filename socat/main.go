package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <in_pipe> <out_pipe>\n", os.Args[0])
		os.Exit(1)
	}
	inPath, outPath := os.Args[1], os.Args[2]

	// Open both FIFOs RDWR so the open() calls never block
	inPipe, err := os.OpenFile(inPath, os.O_RDWR, 0)
	must(err)
	defer inPipe.Close()
	outPipe, err := os.OpenFile(outPath, os.O_RDWR, 0)
	must(err)
	defer outPipe.Close()

	// ---------- 1. read one newline-terminated request from the client ----------
	reqReader := bufio.NewReader(os.Stdin)
	req, err := reqReader.ReadBytes('\n') // includes the '\n'
	if err != nil {
		log.Fatal("client closed before newline:", err)
	}

	// ---------- 2. send request to worker ----------
	_, err = inPipe.Write(req)
	must(err)
	_ = inPipe.Sync() // best-effort flush

	// ---------- 3. read one reply line from worker ----------
	replyReader := bufio.NewReader(outPipe)
	reply, err := replyReader.ReadBytes('\n') // includes the '\n'
	must(err)

	// ---------- 4. HTTP/1.1 framing ----------
	body := bytes.TrimRight(reply, "\n") // keep JSON intact, drop trailing LF for len
	fmt.Fprintf(os.Stdout,
		"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n",
		len(body))
	os.Stdout.Write(body)
	// add the newline back so client sees exactly what worker produced
	os.Stdout.Write([]byte{'\n'})
}
