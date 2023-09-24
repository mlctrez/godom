package demos

import (
	_ "embed"
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
	"github.com/mlctrez/godom/gfet"
	"html"
	"os/exec"
)

//go:embed diff.html
var diffHtml string

type diffExample struct {
	Diff          godom.Element `go:"diff"`
	ShowDiff      godom.Element `go:"showDiff"`
	CommitMessage godom.Element `go:"commitMessage"`
	Commit        godom.Element `go:"commit"`
}

func Diff(ctx *app.Context) godom.Element {
	de := &diffExample{}
	doc := ctx.Doc.WithCallback(callback.Reflect(de))
	result := doc.H(diffHtml)

	clearChildNodes := func() {
		for _, node := range de.Diff.ChildNodes() {
			de.Diff.RemoveChild(node.This())
		}
	}

	de.ShowDiff.AddEventListener("click", func(event godom.Value) {
		go func() {
			req := &gfet.Request{URL: "/diff/api"}
			res, err := req.Fetch()
			clearChildNodes()
			if err != nil {
				de.Diff.AppendChild(doc.T(fmt.Sprintf("error : %s\n", err.Error())))
			} else {
				de.Diff.AppendChild(doc.H(string(res.Body)))
			}
		}()
	})
	de.Commit.AddEventListener("click", func(event godom.Value) {
		go func() {
			req := &gfet.Request{
				URL: "/diff/api", Method: gfet.MethodPost,
				Headers: map[string]string{"commit-message": de.CommitMessage.This().Get("value").String()},
			}
			res, err := req.Fetch()
			clearChildNodes()
			if err != nil {
				de.Diff.AppendChild(doc.T(fmt.Sprintf("error : %s\n", err.Error())))
			} else {
				de.Diff.AppendChild(doc.H(string(res.Body)))
			}
		}()

	})
	return result
}

func DiffServe(request app.Request, response app.Response) {
	switch request.Method() {
	case "GET":
		doGet(response)
	case "POST":
		doPost(request, response)
	}
}

func doPost(request app.Request, response app.Response) {
	commitMessage := request.Headers()["Commit-Message"]
	doc := godom.Document().DocApi()
	fragment := doc.El("pre")
	wouldHaveRun := fmt.Sprintf("# would have run:\ngit add .\ngit commit -m %q\n", commitMessage)
	fragment.AppendChild(doc.T(html.EscapeString(wouldHaveRun)))
	_, _ = response.Write([]byte(fragment.String()))
}

func doGet(response app.Response) {
	cmd := exec.Command("git", "diff")
	output, err := cmd.CombinedOutput()
	doc := godom.Document().DocApi()
	fragment := doc.El("pre")

	if err != nil {
		response.WriteHeader(500)
		fragment.AppendChild(doc.T(fmt.Sprintf("error : %s\n", err.Error())))
	}
	if output != nil {
		fragment.AppendChild(doc.T(html.EscapeString(string(output))))
	}
	_, _ = response.Write([]byte(fragment.String()))
}
