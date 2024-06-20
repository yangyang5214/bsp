package techan

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"strings"
)

type Single struct {
	regionPath string
	resultFile string
	log        *log.Helper
}

func NewSingle(regionPath string) *Single {
	return &Single{
		regionPath: regionPath,
		resultFile: "result.csv",
		log:        log.NewHelper(log.DefaultLogger),
	}
}

func (s *Single) Run() error {
	bytes, err := os.ReadFile(s.regionPath)
	if err != nil {
		return err
	}
	content := string(bytes)
	items := strings.Split(content, "\n")

	f, err := os.Create(s.resultFile)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	_ = writer.Write([]string{"省份", "城市", "特产"})

	for _, item := range items {
		var index *FullIndex
		err = json.Unmarshal([]byte(item), &index)
		if err != nil {
			return err
		}
		products, err := s.process(index)
		if err != nil {
			return err
		}
		_ = writer.Write(append([]string{index.Province, index.Region}, products...))
	}
	return nil
}

func (s *Single) process(item *FullIndex) ([]string, error) {
	urlStr := item.Url

	var products []string
	for i := 0; i < 10; i++ {
		if i != 0 {
			urlStr = item.Url + fmt.Sprintf("/page/%d", i)
		}
		s.log.Infof("start process url: %s", urlStr)
		doc, err := htmlquery.LoadURL(urlStr)
		if err != nil {
			return nil, err
		}

		nodes, err := htmlquery.QueryAll(doc, "//div[@class='card bg-white mb-5 shadow']/a")
		if err != nil {
			return nil, err
		}

		for _, node := range nodes {
			for _, attribute := range node.Attr {
				if attribute.Key == "title" {
					products = append(products, attribute.Val)
				}
			}
		}

		nextNode, err := htmlquery.Query(doc, "//span[@id='nextpage']")
		if err != nil {
			return nil, err
		}
		if nextNode == nil {
			break
		}
	}
	return products, nil
}
