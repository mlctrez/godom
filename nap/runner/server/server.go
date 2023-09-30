package server

import (
	"context"
	"errors"
	"github.com/mlctrez/cmdrunner"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/nap"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Run(c *nap.Config) (err error) {
	s := &server{o: c}
	if s.build() {
		return
	}
	if c.HandleSignals {
		s.ctx, s.cancel = signal.NotifyContext(s.o.Context, os.Interrupt, os.Kill)
	} else {
		s.ctx, s.cancel = context.WithCancel(s.o.Context)
	}
	defer s.cancel()

	if err = s.startHttp(); err != nil {
		return err
	}
	<-s.ctx.Done()
	err = s.stopHttp()
	return err
}

type server struct {
	o      *nap.Config
	ctx    context.Context
	cancel context.CancelFunc
	server *http.Server
}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	request.URL.Path = strings.TrimPrefix(request.URL.Path, s.o.PagesPath)
	s.o.Logger.Debug("ServeHTTP", "url", request.URL.String())
	for _, r := range s.o.Pages {
		if r.Regexp.MatchString(request.URL.Path) {
			writer.Header().Set("Content-Type", "text/html")
			pageHtml := r.PageFunc(godom.Document())
			if _, err := writer.Write([]byte(pageHtml.String())); err != nil {
				s.o.Logger.Error("error writing response", "error", err)
			}
			return
		}
	}
	for _, r := range s.o.Resources {
		if r.Regexp.MatchString(request.URL.Path) {
			writer.Header().Set("Content-Type", r.ContentType)
			if _, err := r.Writer(writer); err != nil {
				s.o.Logger.Error("error writing response", "error", err)
			}
			return
		}
	}
	if s.o.NotFoundFunc != nil {
		writer.Header().Set("Content-Type", "text/html")
		pageHtml := s.o.NotFoundFunc(request.URL, godom.Document())
		if _, err := writer.Write([]byte(pageHtml.String())); err != nil {
			s.o.Logger.Error("error writing response", "error", err)
		}
		return
	}
	// fall back to go's default
	http.DefaultServeMux.ServeHTTP(writer, request)
}

func (s *server) startHttp() error {
	s.o.Logger.Debug("startHttp")
	// TODO: revisit this for correct handling of errors
	listen, err := net.Listen("tcp", s.o.Addr)
	if err != nil {
		s.o.Logger.Error("startHttp", "error", err)
		return err
	}

	s.server = &http.Server{Handler: s}

	go func() {
		serverErr := s.server.Serve(listen)
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			s.o.Logger.Error("error http.Server.Serve: ", serverErr)
		}
		s.o.Logger.Debug("serverExit")
	}()

	return nil
}

func (s *server) stopHttp() (err error) {
	if s.server == nil {
		return nil
	}
	s.o.Logger.Info("stopHttp")

	stopContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = s.server.Shutdown(stopContext); err != nil {
		s.o.Logger.Error("server.Shutdown", "error", err)
		os.Exit(-1)
	}
	return nil
}

func (s *server) build() bool {
	if len(os.Args) < 2 || os.Args[1] != "build" {
		return false
	}
	s.o.Logger.Info("building static content")
	if s.o.BuildOutput == "" {
		log.Fatal("BuildOutput must be set correctly")
	}

	err := os.MkdirAll(s.o.BuildOutput, 0755)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range s.o.Resources {
		out := strings.TrimPrefix(strings.Trim(r.Regexp.String(), "^$"), "/")
		if strings.HasSuffix(out, ".wasm") {
			var mainFile string
			if mainFile, err = findMain(); err != nil {
				log.Fatal(err)
			}
			wasmFile := filepath.Join(s.o.BuildOutput, out)
			if err = os.RemoveAll(wasmFile); err != nil {
				log.Fatal(err)
			}
			command := exec.Command("go", "build", "-o", wasmFile, mainFile)
			command.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
			runner := cmdrunner.NewCmdRunner(command)
			err = runner.Start(func(out *cmdrunner.CmdOutput) {
				if out.Channel == cmdrunner.CmdStderr {
					s.o.Logger.Error(out.Text)
				} else {
					s.o.Logger.Info(out.Text)
				}
			})
			if err != nil {
				log.Fatal("build failed", err)
			}
			if runner.WaitExit() != 0 {
				log.Fatal("build failed, check log output above")
			}
			stat, err := os.Stat(wasmFile)
			if err != nil {
				log.Fatal(err)
			}
			s.o.Logger.Info(out, "size", stat.Size(), "created", stat.ModTime().Format(time.RFC3339))
			continue
		}
		finalPath := filepath.Join(s.o.BuildOutput, out)
		if err = os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			log.Fatal(err)
		}
		var create *os.File
		if create, err = os.Create(filepath.Join(s.o.BuildOutput, out)); err != nil {
			log.Fatal(err)
		}
		if _, err = r.Writer(create); err != nil {
			log.Fatal(err)
		}
		if err = create.Close(); err != nil {
			log.Fatal(err)
		}
	}
	for _, r := range s.o.Pages {
		outParts := strings.Split(r.Regexp.String(), "|")
		for _, out := range outParts {
			out = strings.Trim(out, "^$")
			if out == "/" {
				continue
			}

			var create *os.File
			create, err = os.Create(filepath.Join(s.o.BuildOutput, out+".html"))
			if err != nil {
				log.Fatal(err)
			}
			if _, err = create.Write([]byte(r.PageFunc(godom.Document()).String())); err != nil {
				log.Fatal(err)
			}
			if err = create.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}

	return true
}

func findMain() (file string, err error) {
	for i := 0; i < 10; i++ {
		if pc, f, _, ok := runtime.Caller(i); ok {
			if runtime.FuncForPC(pc).Name() == "main.main" {
				return f, nil
			}
		}
	}
	return "", errors.New("unable to determine main file name")
}
