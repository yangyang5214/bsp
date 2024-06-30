package bd_img

import (
	"bsp/pkg"
	"crypto/md5"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type BdImg struct {
	filepath   string
	httpClient *http.Client
	log        *log.Helper
	chrome     *pkg.ChromePool
}

func NewBdImg(filepath string) (*BdImg, error) {
	chrome, err := pkg.NewChromePool()
	if err != nil {
		return nil, err
	}
	return &BdImg{
		filepath:   filepath,
		httpClient: http.DefaultClient,
		log:        log.NewHelper(log.DefaultLogger),
		chrome:     chrome,
	}, nil
}

func (s *BdImg) Run() error {
	defer s.chrome.Clone()
	byteDatas, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(byteDatas), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		urlStr := buildUrl(line)
		s.log.Info(urlStr)
		err = s.process(urlStr, line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BdImg) process(urlStr string, keyword string) error {
	html, err := s.chrome.NavigateUrl(urlStr, func(p *rod.Page) error {
		for i := 0; i < 5; i++ {
			_, _ = p.Eval(`() => {
			window.scrollBy(0, document.body.scrollHeight);
		}`)
			time.Sleep(time.Second)
		}
		return nil
	})
	if err != nil {
		return err
	}
	imgUrls, err := parseImgs(html)
	if err != nil {
		return err
	}
	s.log.Infof("for key %s, urls size %d", keyword, len(imgUrls))

	dir := path.Join("", keyword)

	_ = os.MkdirAll(dir, 0755)

	var wg sync.WaitGroup
	wg.Add(len(imgUrls))
	for _, imgUrl := range imgUrls {
		go func(url string) {
			defer wg.Done()
			err = s.downloadImg(dir, url)
			if err != nil {
				panic(err)
			}
		}(imgUrl)
	}
	wg.Wait()
	return nil
}

func (s *BdImg) downloadImg(dir string, imgUrl string) error {
	resp, err := s.httpClient.Get(imgUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ext string
	for k, v := range resp.Header {
		if k == "Content-Type" {
			ext = v[0]
			ext = strings.Split(ext, "/")[1]
		}
	}
	p := path.Join(dir, fmt.Sprintf("%x.%s", md5.Sum(bytes), ext))
	return os.WriteFile(p, bytes, 0755)
}

func buildUrl(key string) string {
	vals := url.Values{}
	vals.Add("tn", "baiduimage")
	vals.Add("ie", "utf-8")
	vals.Add("word", key)

	urlStr := "https://image.baidu.com/search/index" + "?" + vals.Encode()
	return urlStr
}

func parseImgs(html string) ([]string, error) {
	xpath := "//img[@class='main_img img-hover']"
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	nodes, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, node := range nodes {
		val := htmlquery.SelectAttr(node, "src")
		if !strings.HasPrefix(val, "https") {
			continue
		}
		urls = append(urls, val)
	}
	return urls, nil
}
