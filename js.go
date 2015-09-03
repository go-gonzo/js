package js

// github.com/yosssi/gcss binding for gonzo.
// No Configuration required.

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

func Minify() gonzo.Stage {
	return func(ctx context.Context, in <-chan gonzo.File, out chan<- gonzo.File) error {

		for {
			select {
			case file, ok := <-in:
				if !ok {
					return nil
				}

				buff := new(bytes.Buffer)
				name := strings.TrimSuffix(file.FileInfo().Name(), ".js") + ".min.js"
				ctx.Infof("Compiling %s to %s", file.FileInfo().Name(), name)
				m := minify.New()
				m.AddFunc("text/css", js.Minify)
				err := m.Minify("text/css", buff, file)
				if err != nil {
					return err
				}

				file = gonzo.NewFile(ioutil.NopCloser(buff), file.FileInfo())
				file.FileInfo().SetSize(int64(buff.Len()))
				file.FileInfo().SetName(name)

				out <- file
			case <-ctx.Done():
				return nil //ctx.Err()
			}
		}
	}
}
