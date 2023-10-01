package server

import (
	"context"
	"errors"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/nap"
	"github.com/mlctrez/godom/nap/runner/build"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Run(c *nap.Config) (err error) {
	s := &server{o: c}

	if build.ArgsMatch() {
		return build.Run(c)
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
