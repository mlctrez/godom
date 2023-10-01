package nap

import (
	"context"
	"fmt"
	"github.com/mlctrez/godom"
	"io"
	"log/slog"
	"net/url"
	"os"
	"regexp"
)

//TODO: change name of Config struct?
//		bring back in options configuration?
//		make some of these private ( or not based on where they are referenced )

type Config struct {
	Addr            string
	Context         context.Context
	Logger          *slog.Logger
	App             App
	Resources       []*Resource
	Pages           []*Page
	NotFoundFunc    NotFoundFunc
	PagesPath       string
	HandleSignals   bool
	JetBrainsReload bool
	BuildOutput     string
	isServer        bool
}

func New() (c *Config) {
	c = &Config{}
	c.Addr = "0.0.0.0:8080"
	c.Context = context.Background()
	c.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	c.HandleSignals = true
	return c
}

func (c *Config) WithApp(app App) *Config {
	c.App = app
	c.App.Init(c, godom.Document().DocumentElement())
	return c
}

func (c *Config) Setup(isServer bool) error {
	c.isServer = isServer
	if !isServer {
		godom.Console().Log("%c"+"⬜⬜⬜ startup ⬜⬜⬜", "font-size: 1.5em; color: white;")
	}
	if c.App == nil {
		return fmt.Errorf("missing App in configuration")
	}
	if c.Logger == nil {
		return fmt.Errorf("missing Logger in configuration")
	}
	c.App.Routes(c)
	return nil
}

func (c *Config) IsServer() bool {
	return c.isServer
}

type DocContext interface {
	DocApi() godom.DocApi
}

type App interface {
	Init(o *Config, html godom.Element)
	Routes(router Router)
}

// Router is used to add valid routes to the app.
type Router interface {
	Resource(re *regexp.Regexp, contentType string, resourceWriter ResourceWriter)
	Page(re *regexp.Regexp, pageFunc PageFunc)
	NotFound(fn NotFoundFunc)
}

type PageFunc func(ctx DocContext) godom.Element
type ResourceWriter func(writer io.Writer) (n int, err error)
type NotFoundFunc func(u *url.URL, ctx DocContext) godom.Element

type ResourceBytes []byte

func (rb ResourceBytes) Write(writer io.Writer) (n int, err error) {
	return writer.Write(rb)
}

type Resource struct {
	Regexp      *regexp.Regexp
	ContentType string
	Writer      ResourceWriter
}

func (c *Config) Resource(re *regexp.Regexp, contentType string, resourceWriter ResourceWriter) {
	if re == nil {
		panic("nil Regexp")
	}
	if contentType == "" {
		panic("empty contentType")
	}
	if resourceWriter == nil {
		panic("nil resourceWriter")
	}
	c.Logger.Debug("Router.Resource", "regex", re)
	r := &Resource{Regexp: re, ContentType: contentType, Writer: resourceWriter}
	c.Resources = append(c.Resources, r)
}

type Page struct {
	Regexp   *regexp.Regexp
	PageFunc PageFunc
}

func (c *Config) Page(re *regexp.Regexp, p PageFunc) {
	if re == nil {
		panic("nil regular expression")
	}
	if p == nil {
		panic("nil PageFunc")
	}
	c.Logger.Debug("Router.Page", "regex", re)
	c.Pages = append(c.Pages, &Page{Regexp: re, PageFunc: p})
}

func (c *Config) NotFound(renderer NotFoundFunc) {
	c.NotFoundFunc = renderer
}

type Event interface {
}

type Location struct {
	Event
	URL      *url.URL
	External bool
	PopState bool
}

func Re(p string) *regexp.Regexp {
	return regexp.MustCompile("^" + p + "$")
}
