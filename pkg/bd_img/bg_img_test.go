package bd_img

import (
	"testing"
)

func TestName(t *testing.T) {
	bd, err := NewBdImg("/Users/beer/beer/bsp/pkg/bd_img/shetou.txt")
	if err != nil {
		panic(err)
	}
	err = bd.Run()
	if err != nil {
		panic(err)
	}
}

func TestDownloadImg(t *testing.T) {
	bd, err := NewBdImg("/Users/beer/beer/bsp/pkg/bd_img/shetou.txt")
	if err != nil {
		panic(err)
	}
	err = bd.downloadImg("/tmp", "https://img0.baidu.com/it/u=4111734781,2214031660&fm=253&fmt=auto&app=138&f=JPEG?w=696&h=500")
	if err != nil {
		panic(err)
	}

}
