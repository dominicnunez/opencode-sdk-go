package apijson

import (
	"reflect"
	"sync"
	"testing"

	"github.com/tidwall/gjson"
)

type testUnionA struct{ A string }
type testUnionB struct{ B string }
type testUnion interface{ testUnion() }

func (testUnionA) testUnion() {}
func (testUnionB) testUnion() {}

func TestRegisterUnionConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			variant := UnionVariant{
				TypeFilter: gjson.JSON,
				Type:       reflect.TypeOf(testUnionA{}),
			}
			if i%2 == 0 {
				variant.Type = reflect.TypeOf(testUnionB{})
			}
			RegisterUnion(reflect.TypeOf((*testUnion)(nil)).Elem(), "", variant)
		}(i)
	}
	wg.Wait()
}
