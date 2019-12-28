package dotmap_test

import (
	"fmt"
	"github.com/zerjioang/dotmap"
	"testing"
)

func TestNoop(t *testing.T) {
	t.Run("usaage-example", func(t *testing.T) {
		mm := dotmap.New()
		mm.Reset(map[string]interface{}{
			"foo": "bar",
			"enableDebug": false,
			"version": 1.0,
			"config": map[string]interface{}{
				"http": 2,
			},
		})
		v, found := dotmap.GetDotMap(mm, "config.http")
		fmt.Println("key found: ", found)
		fmt.Println("key value: ", v)
		v2, found2 := dotmap.GetDotMap(mm, "config.key")
		fmt.Println("key found: ", found2)
		fmt.Println("key value: ", v2)
	})
}

