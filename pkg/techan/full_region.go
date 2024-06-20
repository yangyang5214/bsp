package techan

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"os"
	"strings"
)

type FullRegion struct {
	client   *http.Client
	outPath  string
	log      *log.Helper
	htmlPath string
}

func NewFullRegion(htmlPath string) *FullRegion {
	return &FullRegion{
		client:   http.DefaultClient,
		outPath:  "region.txt",
		log:      log.NewHelper(log.DefaultLogger),
		htmlPath: htmlPath,
	}
}

func (s *FullRegion) Run() error {
	doc, err := htmlquery.LoadDoc(s.htmlPath)
	if err != nil {
		return err
	}
	body := htmlquery.FindOne(doc, "//div[@class='row']")

	nodes := htmlquery.Find(body, "//div")

	var (
		indexs   []*FullIndex
		province string
	)

	for _, div := range nodes {
		text := htmlquery.InnerText(div)

		text = strings.TrimSpace(text)
		text = strings.Trim(text, "\n")
		text = strings.Trim(text, "特产")

		for _, attribute := range div.Attr {
			if attribute.Val == "col-12 mt-2" {
				province = text
				break
			}
		}

		if province == "" {
			return errors.New("province parser failed")
		}

		attrs := div.FirstChild.Attr
		if len(attrs) != 2 {
			continue //ugly
		}

		url := attrs[1].Val

		indexs = append(indexs, &FullIndex{
			Url:      url,
			Province: province,
			Region:   text,
		})
		s.log.Infof("%s-%s, %s", province, text, url)
	}

	var b bytes.Buffer
	for _, item := range indexs {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		_, err = b.WriteString(string(data) + "\n")
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(s.outPath, b.Bytes(), 0655)
	if err != nil {
		return err
	}
	return nil
}
