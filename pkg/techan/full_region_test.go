package techan

import "testing"

func TestName(t *testing.T) {
	err := NewFullRegion("/Users/beer/beer/bsp/pkg/techan/main.html").Run()
	if err != nil {
		panic(err)
	}
}
