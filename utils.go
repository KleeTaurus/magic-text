package summaryit

import (
	"fmt"
	"path"
	"strings"
)

func getOutfile(infile, ext string) string {
	dirname := path.Dir(infile)
	basename := path.Base(infile)

	return fmt.Sprintf("%s/%s%s", dirname, strings.Split(basename, ".")[0], ext)
}
