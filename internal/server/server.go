package main

import (
	"context"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/cskr/pubsub"
	"github.com/mlctrez/wasmexec"
	"github.com/rjeczalik/notify"
	"log"
	"magefiles/watcher"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"os/exec"
	"strings"
	"time"
)

var pubSub = pubsub.New(10)

func main() {
	err := DevServer()
	if err != nil {
		log.Fatal(err)
	}
}

func DevServer() (err error) {

	if err = BuildWasm(); err != nil {
		return err
	}

	var w *watcher.Watcher
	if w, err = watcher.New(fileChange, "app"); err != nil {
		return err
	}
	go w.Run()

	fmt.Println("dev server running on http://localhost:8080")
	Server()

	return nil
}

func fileChange(info notify.EventInfo) {
	//fmt.Printf("%s file changed %s\n", time.Now().Format(time.RFC3339Nano), info.Path())
	if err := BuildWasm(); err != nil {
		fmt.Println(strings.TrimSpace(err.Error()))
		return
	}
	pubSub.Pub("wasm", "build")
}

func BuildWasm() error {
	command := exec.Command("go", "build", "-o", "build/app.wasm", "app/main.go")
	command.Env = append(os.Environ(), "GOARCH=wasm", "GOOS=js")
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error building wasm: %s\n%s\n", err, string(output))
	}
	var stat os.FileInfo
	stat, err = os.Stat("build/app.wasm")
	if err != nil {
		return fmt.Errorf("error getting wasm size: %s\n", err)
	}
	fmt.Println("wasm size", stat.Size())
	return nil
}

func Server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/app.js", AppJs)
	mux.HandleFunc("/app.wasm", Wasm)
	mux.HandleFunc("/ws", Echo)
	_ = http.ListenAndServe(":8080", mux)
}

func Wasm(writer http.ResponseWriter, request *http.Request) {

	file, err := os.ReadFile("build/app.wasm")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/wasm")
	if strings.Contains(request.Header.Get("Accept-Encoding"), "sbr") {
		writer.Header().Set("Content-Encoding", "br")
		start := time.Now()
		brWriter := brotli.NewWriterLevel(writer, brotli.DefaultCompression)
		_, _ = brWriter.Write(file)
		_ = brWriter.Flush()
		fmt.Println("compression took", time.Since(start).String())
	} else {
		_, _ = writer.Write(file)
	}

}

var clientNumber int

func Echo(writer http.ResponseWriter, request *http.Request) {
	options := buildAcceptOptions(request)
	var client *websocket.Conn
	var err error
	if client, err = websocket.Accept(writer, request, options); err != nil {
		fmt.Println("websocket.Accept", err)
		return
	}
	defer func() {
		_ = client.Close(websocket.StatusServiceRestart, "server exit")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientNumber = clientNumber + 1

	sub := pubSub.Sub("build")
	defer func() {
		pubSub.Unsub(sub, "hello")
		fmt.Println("pubSub.Unsub client", clientNumber)
	}()

	go func() {
		for {
			var typ websocket.MessageType
			var bytes []byte
			if typ, bytes, err = client.Read(ctx); err != nil {
				cancel()
				return
			}
			if err = client.Write(ctx, typ, bytes); err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			msgString := msg.(string)
			typ := websocket.MessageBinary
			if err = client.Write(ctx, typ, []byte(msgString)); err != nil {
				cancel()
				return
			}
			fmt.Printf("%s pub message for client %d : %s\n",
				time.Now().Format(time.RFC3339Nano), clientNumber, msgString)
		}
	}
}

func buildAcceptOptions(request *http.Request) *websocket.AcceptOptions {
	var options *websocket.AcceptOptions
	// https://github.com/gorilla/websocket/issues/731
	// Compression in certain Safari browsers is broken, turn it off
	if strings.Contains(request.UserAgent(), "Safari") {
		options = &websocket.AcceptOptions{CompressionMode: websocket.CompressionDisabled}
	}
	return options
}

func Index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	_, _ = writer.Write([]byte(indexHtml))
}

var indexHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Index</title>
    <style>
        body {
            color: white;
            background-color: black;
        }
    </style>
    <script type="application/javascript" src="app.js"></script>
</head>
<body><p>loading</p></body>
</html>
`

func AppJs(writer http.ResponseWriter, request *http.Request) {
	wasmexec.WriteLauncher(writer)
}
