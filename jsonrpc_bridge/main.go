package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	inPipe  = "/tmp/mcp-pipes/in"
	outPipe = "/tmp/mcp-pipes/out"
)

func main() {
	// open FIFOs RDWR once per process; share across all requests
	in, err := os.OpenFile(inPipe, os.O_RDWR, 0)
	must(err)
	defer in.Close()
	out, err := os.OpenFile(outPipe, os.O_RDWR, 0)
	must(err)
	defer out.Close()

	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handler(in, out),
	}

	log.Println("JSON-RPC bridge on :8080")
	log.Fatal(srv.ListenAndServe())
}

func handler(in, out *os.File) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.Method, "POST") {
			http.Error(w, "only POST accepted", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MiB guard
		fmt.Println("read body:", string(body), "err:", err)
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}

		// send payload terminated by newline
		if len(body) == 0 || body[len(body)-1] != '\n' {
			body = append(body, '\n')
		}
		if _, err := in.Write(body); err != nil {
			http.Error(w, "backend write error", http.StatusBadGateway)
			return
		}

		// pick a context deadline tied to the client's connection
		ctx := r.Context()
		reply, err := readReply(ctx, out)
		if err != nil {
			http.Error(w, "backend timeout", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprint(len(reply)))
		w.WriteHeader(http.StatusOK)
		w.Write(reply)
	}
}

// read one newline-terminated reply, but cancel if the client goes away
func readReply(ctx context.Context, out *os.File) ([]byte, error) {
	reader := bufio.NewReader(out)
	type result struct {
		b   []byte
		err error
	}
	ch := make(chan result, 1)

	go func() {
		b, err := reader.ReadBytes('\n')
		fmt.Println("readBytes:", string(b), "err:", err)
		ch <- result{b: bytes.TrimRight(b, "\n"), err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		return res.b, res.err
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
