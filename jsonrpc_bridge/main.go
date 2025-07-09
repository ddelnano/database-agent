package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	inPipe  = "/tmp/mcp-pipes/in"
	outPipe = "/tmp/mcp-pipes/out"
)

var (
	mcpMu  sync.Mutex    // serialize access to MCP
	outRdr *bufio.Reader // shared reader over outPipe
)

func main() {
	in, err := os.OpenFile(inPipe, os.O_RDWR, 0)
	must(err)
	defer in.Close()

	out, err := os.OpenFile(outPipe, os.O_RDWR, 0)
	must(err)
	defer out.Close()

	outRdr = bufio.NewReader(out)

	// optional: pprof on :6060
	go func() { log.Println(http.ListenAndServe(":6060", nil)) }()

	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handler(in),
	}

	log.Println("JSON-RPC bridge on :8080")
	log.Fatal(srv.ListenAndServe())
}

func handler(in *os.File) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}

		//------------------------------------------------------------------
		// Critical section: one request → one reply on the shared FIFOs
		//------------------------------------------------------------------
		mcpMu.Lock()
		defer mcpMu.Unlock()

		// 1) send request (always terminated by '\n' for MCP)
		if len(body) == 0 || body[len(body)-1] != '\n' {
			body = append(body, '\n')
		}
		if _, err := in.Write(body); err != nil {
			http.Error(w, "backend write", http.StatusBadGateway)
			return
		}

		// 2) read exactly one JSON value back
		raw, err := decodeOneJSON(r.Context())
		if err != nil {
			http.Error(w, "backend read/timeout", http.StatusBadGateway)
			return
		}

		//------------------------------------------------------------------

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprint(len(raw)))
		w.WriteHeader(http.StatusOK)
		w.Write(raw)
	}
}

// decodeOneJSON reads continuously from outRdr until it successfully
// decodes one complete JSON value or ctx is cancelled.
func decodeOneJSON(ctx context.Context) ([]byte, error) {
	dec := json.NewDecoder(outRdr)
	for {
		var raw json.RawMessage
		err := dec.Decode(&raw)
		if err == nil {
			return raw, nil // ✔ got the reply
		}
		if err == io.EOF {
			// worker closed pipe unexpectedly
			return nil, io.ErrUnexpectedEOF
		}

		// err is a *syntax* error: bytes weren’t JSON.
		// Discard one byte and retry so we resync on the next '{'.
		if serr, ok := err.(*json.SyntaxError); ok {
			if _, derr := outRdr.Discard(int(serr.Offset)); derr != nil {
				return nil, derr
			}
			continue
		}

		// Any other error: give up.
		return nil, err
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
