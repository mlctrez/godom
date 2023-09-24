package demos

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/gfet"
	"runtime"
)

func Diff(ctx *app.Context) godom.Element {
	if runtime.GOARCH == "wasm" {
		fmt.Println("wasm yes")
		req := gfet.Request{URL: "/git/diff"}
		res, err := req.Fetch()
		if err != nil {
			return ctx.Doc.H(fmt.Sprintf("<div>%s</div>", err))
		}
		scanner := bufio.NewScanner(bytes.NewBuffer(res.Body))
		pre := ctx.Doc.El("pre")
		for scanner.Scan() {
			pre.AppendChild(ctx.Doc.T(scanner.Text() + "\n"))
		}
		return pre
	} else {
		return ctx.Doc.H("<div>not implemented</div>")
	}

}
