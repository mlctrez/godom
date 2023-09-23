package watcher

// Wraps github.com/rjeczalik/notify to add recursive watches and de-duplication of events.

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
)

type Watcher struct {
	paths   []string
	done    chan bool
	noTilde chan notify.EventInfo
	noDupes chan notify.EventInfo
	handler func(info notify.EventInfo)
}

func (w *Watcher) Stop() {
	w.done <- true
	notify.Stop(w.noTilde)
}

func (w *Watcher) noTildeRunner() {
	for {
		select {
		case event, ok := <-w.noTilde:
			if !ok {
				w.done <- true
				return
			}
			if !strings.HasSuffix(event.Path(), "~") {
				w.noDupes <- event
			}
		}
	}
}

func (w *Watcher) noDupesRunner() {
	for {
		var eventOne notify.EventInfo = nil
		var eventTwo notify.EventInfo = nil
		var ok = false

		// read first event
		select {
		case <-w.done:
			return
		case eventOne, ok = <-w.noDupes:
			if !ok {
				w.done <- true
				return
			}
		}
		//fmt.Printf("%s read eventOne %s %s\n", time.Now().Format(time.RFC3339Nano),
		//	eventOne.Path(), eventOne.Event())

		// look for another event closely following this one, bail after a few ms
		select {
		case <-w.done:
			return
		case eventTwo, ok = <-w.noDupes:
			if !ok {
				w.done <- true
				return
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("%s timed out waiting on another event\n", time.Now().Format(time.RFC3339Nano))
			// nothing
		}
		w.handler(eventOne)
		if eventTwo != nil && eventTwo.Path() != eventOne.Path() {
			fmt.Printf("%s read eventTwo %s %s\n", time.Now().Format(time.RFC3339Nano),
				eventTwo.Path(), eventTwo.Event())
			w.handler(eventTwo)
		}
	}
}

func New(handler func(info notify.EventInfo), paths ...string) (w *Watcher, err error) {
	w = &Watcher{
		handler: handler,
		paths:   make([]string, 0),
		done:    make(chan bool),
		noTilde: make(chan notify.EventInfo),
		noDupes: make(chan notify.EventInfo),
	}

	for _, path := range paths {
		var stat os.FileInfo
		if stat, err = os.Stat(path); err != nil {
			return
		}
		if !stat.IsDir() {
			w.paths = append(w.paths, path)
			continue
		}
		err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				w.paths = append(w.paths, path)
			}
			return nil
		})
		if err != nil {
			return
		}
	}

	return
}

func (w *Watcher) Run() {
	go w.noTildeRunner()
	go w.noDupesRunner()
	for _, path := range w.paths {
		//fmt.Println("watching", path)
		err := notify.Watch(path, w.noTilde, notify.Write)
		if err != nil {
			panic(err)
		}
	}
	<-w.done
}
