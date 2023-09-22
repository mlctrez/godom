package main

import (
	"context"
	"fmt"
	"github.com/mlctrez/cmdrunner"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var buildDir string

func Test(ctx context.Context) (err error) {
	if err = setupBuildDir(); err != nil {
		return err
	}
	steps := []struct {
		args []string
		wasm bool
	}{
		{args: []string{"test", "-race", "-covermode", "atomic", "./...", "-args", pathArg("test.gocoverdir")}},
		{args: []string{"test", "-covermode", "atomic", "./...", "-args", pathArg("test.gocoverdir")}, wasm: true},
		{args: []string{"tool", "covdata", "textfmt", pathArg("i"), pathArg("o", "coverage.out")}},
		{args: []string{"tool", "cover", pathArg("html", "coverage.out"), pathArg("o", "coverage.html")}},
	}

	outputSink := func(out *cmdrunner.CmdOutput) {
		fmt.Println(out.Text)
	}
	for _, step := range steps {
		if err = goCmd(outputSink, step.wasm, step.args...); err != nil {
			return err
		}
	}

	if err = removeCovFiles(); err != nil {
		return err
	}
	if err = summarizeCoverage(); err != nil {
		return err
	}
	return nil
}

func pathArg(arg string, paths ...string) string {
	p := filepath.Join(append([]string{buildDir}, paths...)...)
	return fmt.Sprintf("-%s=%s", arg, p)
}

// summarizeCoverage prints lines from go tool cover -func that are not 100% when total is < 100%.
func summarizeCoverage() (err error) {
	var modulePrefix string
	outputSink := func(out *cmdrunner.CmdOutput) {
		modulePrefix = strings.TrimSpace(out.Text) + "/"
	}
	if err = goCmd(outputSink, false, "list", "-m"); err != nil {
		return err
	}

	var badCoverage []string
	var totalNotGood bool

	outputSink = func(out *cmdrunner.CmdOutput) {
		outLine := strings.TrimPrefix(out.Text, modulePrefix)
		if strings.HasPrefix(outLine, "total:") {
			totalNotGood = !strings.HasSuffix(outLine, "100.0%")
		}
		if !strings.HasSuffix(outLine, "100.0%") {
			badCoverage = append(badCoverage, outLine)
		}
	}
	if err = goCmd(outputSink, false, "tool", "cover", pathArg("func", "coverage.out")); err != nil {
		return nil
	}

	if totalNotGood {
		for _, s := range badCoverage {
			fmt.Println(s)
		}
	}

	return nil

}

func setupBuildDir() (err error) {
	if err = os.MkdirAll("build", 0755); err != nil {
		return err
	}
	if buildDir, err = filepath.Abs("build"); err != nil {
		return err
	}
	if err = removeCovFiles(); err != nil {
		return err
	}
	return nil
}

func removeCovFiles() (err error) {
	if err = rmFiles(
		filepath.Join(buildDir, "covcounters*"),
		filepath.Join(buildDir, "covmeta*"),
	); err != nil {
		return err
	}
	return nil
}

func rmFiles(patterns ...string) (err error) {
	for _, pattern := range patterns {
		var matches []string
		matches, err = filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, match := range matches {
			err = os.Remove(match)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func goCmd(outputSink cmdrunner.OutputSink, wasm bool, args ...string) (err error) {
	command := exec.Command("go", args...)
	if wasm {
		command.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	}
	fmt.Println(strings.Join(command.Args, " "), fmt.Sprintf("(wasm=%t)", wasm))
	return cmd(command, outputSink)
}

func cmd(command *exec.Cmd, outputSink cmdrunner.OutputSink) (err error) {
	r := cmdrunner.NewCmdRunner(command)
	if err = r.Start(outputSink); err != nil {
		return err
	}
	if exit := r.WaitExit(); exit != 0 {
		return fmt.Errorf("%+v exit code %d", command.Args, exit)
	}
	return nil
}
