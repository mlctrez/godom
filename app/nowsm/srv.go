package nowsm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/cskr/pubsub"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/watcher"
	"github.com/mlctrez/wasmexec"
	"github.com/rjeczalik/notify"
	"nhooyr.io/websocket"
)

func Run(h app.Handler) {

	sc := &app.ServerContext{}
	h.Prepare(sc)

	// TODO: separate out build and watcher into a component
	var err error
	if err = BuildWasm(sc); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := &Server{h: h, sc: sc, pubSub: pubsub.New(10), clientNumber: 0}

	var w *watcher.Watcher
	if w, err = watcher.New(s.fileChange, sc.Watch...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go w.Run()

	server := &http.Server{Handler: s, Addr: sc.Address}
	go func() {
		var message = "dev server running on http" + "://"
		if strings.HasPrefix(sc.Address, ":") {
			message += "localhost" + sc.Address
		} else {
			message += sc.Address
		}
		fmt.Println(message)
		err = server.ListenAndServe()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	ctx := contextWithSigterm(context.Background())
	<-ctx.Done()
	_ = server.Shutdown(context.TODO())

}

func (s *Server) fileChange(info notify.EventInfo) {
	//fmt.Printf("%s file changed %s\n", time.Now().Format(time.RFC3339Nano), info.Path())
	if err := BuildWasm(s.sc); err != nil {
		fmt.Println(strings.TrimSpace(err.Error()))
		return
	}
	s.pubSub.Pub("wasm", "build")
}

func BuildWasm(sc *app.ServerContext) error {
	command := exec.Command("go", "build", "-o", sc.Output, sc.Main)
	command.Env = append(os.Environ(), "GOARCH=wasm", "GOOS=js")
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error building wasm: %s\n%s\n", err, string(output))
	}
	if sc.ShowWasmSize {
		var s os.FileInfo
		s, err = os.Stat(sc.Output)
		if err != nil {
			return fmt.Errorf("error getting wasm size: %s\n", err)
		}
		fmt.Println("wasm size:", s.Size())
	}
	return nil
}

func contextWithSigterm(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}

type Server struct {
	h            app.Handler
	clientNumber int
	pubSub       *pubsub.PubSub
	sc           *app.ServerContext
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	// TODO: revisit if api.Handler should be able to override these
	//   or provide additional endpoints

	switch request.URL.Path {
	case "/app.js":
		AppJs(writer, request)
	case "/favicon.ico":
		writer.WriteHeader(http.StatusNotFound)
	case "/app.wasm":
		s.Wasm(writer, request)
	case "/ws":
		s.Echo(writer, request)
	default:
		s.defaultRoute(writer, request)
	}
}

func (s *Server) defaultRoute(writer http.ResponseWriter, request *http.Request) {
	document := godom.Document()
	doc := document.DocApi()

	html := doc.El("html", doc.At("lang", "en"))
	head := doc.H(`<head>
	<meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>TODO</title>
    <script type="application/javascript" src="app.js"></script>
</head>`)
	ctx := &app.Context{Doc: doc, URL: request.URL, Events: nil}
	html.AppendChild(head)
	s.h.Headers(ctx, head)
	html.AppendChild(s.h.Body(ctx))
	buf := &bytes.Buffer{}
	enc := godom.NewEncoder(buf)
	html.Marshal(enc).Flush()
	buf.WriteString("\n")
	writer.Header().Set("Content-Type", "text/html")
	_, _ = writer.Write(buf.Bytes())
}

func (s *Server) Wasm(writer http.ResponseWriter, request *http.Request) {
	file, err := os.ReadFile(s.sc.Output)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/wasm")
	if strings.Contains(request.Header.Get("Accept-Encoding"), "br") {
		writer.Header().Set("Content-Encoding", "br")
		//start := time.Now()
		brWriter := brotli.NewWriterLevel(writer, brotli.DefaultCompression)
		_, _ = brWriter.Write(file)
		_ = brWriter.Flush()
		//fmt.Println("compression took", time.Since(start).String())
	} else {
		_, _ = writer.Write(file)
	}
}

func (s *Server) Echo(writer http.ResponseWriter, request *http.Request) {
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

	s.clientNumber = s.clientNumber + 1

	sub := s.pubSub.Sub("build")
	defer func() {
		s.pubSub.Unsub(sub, "build")
		fmt.Println("pubSub.Unsub client", s.clientNumber)
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
				time.Now().Format(time.RFC3339Nano), s.clientNumber, msgString)
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

func AppJs(writer http.ResponseWriter, _ *http.Request) {
	wasmexec.WriteLauncher(writer)
}
