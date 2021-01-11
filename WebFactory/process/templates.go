package process

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/pretty"
)

const (
	templateName    = `$TEMPLATE`
	templateContent = `$CONTENT`

	cssTemplate = `

{{define "$TEMPLATE"}}

<!---->

<style>

$CONTENT

</style>

<!---->

{{end}}

	`

	jsTemplate = `

{{define "$TEMPLATE"}}

<!---->

<script>

$CONTENT

</script>

<!---->

{{end}}

	`
)

type TemplateOperator struct {
	CleanSrc bool
}

func (t *TemplateOperator) Operate(src string) error {
	ext := filepath.Ext(src)

	if ext != ".js" && ext != ".css" {
		return fmt.Errorf("unknown file extension: %s", ext)
	}

	var err error
	dst := t.toDstName(src)

	pretty.ColorPrintln(pretty.PaperBlue300, fmt.Sprintf("processing file: %s -> %s", src, dst))

	err = t.makeTemplate(src, dst)
	if err != nil {
		return err
	}

	if t.CleanSrc {
		err = os.Remove(src)
		if err != nil {
			return err
		}

		if filepath.Ext(src) == ".js" {
			err = os.Remove(fmt.Sprintf("%s.deps", src))
			if err != nil {
				return err
			}

			err = os.Remove(fmt.Sprintf("%s.map", src))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TemplateOperator) makeTemplate(src, dst string) error {
	name := t.toTemplateName(src)
	ft := filepath.Ext(src)[1:]

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	content = t.process(content, name, ft)
	if err != nil {
		return err
	}

	nf, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = nf.Write(content)
	if err != nil {
		return err
	}

	nf.Sync()
	nf.Close()

	return nil
}

func (t *TemplateOperator) process(content []byte, name, fileType string) []byte {
	var template []byte

	switch fileType {
	case "css":
		template = []byte(cssTemplate)
	case "js":
		template = []byte(jsTemplate)
	default:
		panic("invalid file type")
	}

	template = bytes.ReplaceAll(template, []byte(templateName), []byte(name))
	template = bytes.ReplaceAll(template, []byte(templateContent), content)
	template = bytes.TrimSpace(template)

	return template
}

func (t *TemplateOperator) toTemplateName(src string) string {
	fn := strings.ReplaceAll(filepath.Base(src), filepath.Ext(src), "")
	ft := filepath.Ext(src)[1:]

	s := strings.Split(strings.ReplaceAll(fn, "-", "_"), "_")
	for i, v := range s {
		s[i] = strings.Title(v)
	}

	return fmt.Sprintf("%s%s", strings.Join(s, ""), strings.ToUpper(ft))
}

func (t *TemplateOperator) toDstName(src string) string {
	const regex = "[^_]+_(?:js|css).html"
	re := regexp.MustCompile(regex)

	if re.MatchString(src) {
		return src
	}

	name := strings.ReplaceAll(strings.ReplaceAll(src, "-", "_"), filepath.Ext(src), "")
	ft := filepath.Ext(src)[1:]

	return fmt.Sprintf("%v_%v%v", name, ft, ".html")
}
