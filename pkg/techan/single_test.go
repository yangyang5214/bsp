package techan

import "testing"

func TestParserSingle(t *testing.T) {
	item := &FullIndex{
		Url:      "https://www.zhtechan.cn/haerbin/",
		Province: "",
		Region:   "",
	}
	r, err := NewSingle("").process(item)
	if err != nil {
		panic(err)
	}
	t.Log(r)
}
