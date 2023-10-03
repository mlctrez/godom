package build

import (
	"errors"
	"github.com/mlctrez/cmdrunner"
	"github.com/mlctrez/godom/nap"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func ArgsMatch() bool {
	return len(os.Args) > 1 && os.Args[1] == "build"
}

func Run(o *nap.Config) (err error) {

	o.Logger.Info("building static content")
	if o.BuildOutput == "" {
		return errors.New("Config.BuildOutput must be set")
	}

	if err = os.MkdirAll(o.BuildOutput, 0755); err != nil {
		return err
	}
	for _, r := range o.Resources {
		out := strings.TrimPrefix(strings.Trim(r.Regexp.String(), "^$"), "/")
		if strings.HasSuffix(out, ".wasm") {
			var mainFile string
			if mainFile, err = findMain(); err != nil {
				return err
			}
			wasmFile := filepath.Join(o.BuildOutput, out)
			if err = os.RemoveAll(wasmFile); err != nil {
				return err
			}
			command := exec.Command("go", "build", "-o", wasmFile, mainFile)
			command.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
			runner := cmdrunner.NewCmdRunner(command)
			err = runner.Start(func(out *cmdrunner.CmdOutput) {
				if out.Channel == cmdrunner.CmdStderr {
					o.Logger.Error(out.Text)
				} else {
					o.Logger.Info(out.Text)
				}
			})
			if err != nil {
				return err
			}
			if runner.WaitExit() != 0 {
				return errors.New("build failed, check log output above")
			}
			var stat os.FileInfo
			if stat, err = os.Stat(wasmFile); err != nil {
				return err
			}
			o.Logger.Info(out, "size", stat.Size(), "created", stat.ModTime().Format(time.RFC3339))
			continue
		}
		finalPath := filepath.Join(o.BuildOutput, out)
		if err = os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			log.Fatal(err)
		}
		var create *os.File
		if create, err = os.Create(filepath.Join(o.BuildOutput, out)); err != nil {
			log.Fatal(err)
		}
		if _, err = r.Writer(create); err != nil {
			log.Fatal(err)
		}
		if err = create.Close(); err != nil {
			log.Fatal(err)
		}
	}
	for _, r := range o.Pages {
		outParts := strings.Split(r.Regexp.String(), "|")
		for _, out := range outParts {
			out = strings.Trim(out, "^$")
			if out == "/" {
				continue
			}

			var create *os.File
			create, err = os.Create(filepath.Join(o.BuildOutput, out+".html"))
			if err != nil {
				log.Fatal(err)
			}

			ctx := nap.NewDocContext(&url.URL{Path: out}, nil)

			if _, err = create.Write([]byte(r.PageFunc(ctx).String())); err != nil {
				log.Fatal(err)
			}
			if err = create.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
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
