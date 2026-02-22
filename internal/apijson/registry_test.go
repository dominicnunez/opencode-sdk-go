package apijson

import (
	"reflect"
	"sync"
	"testing"

	"github.com/tidwall/gjson"
)

func TestRegisterUnionConcurrent(t *testing.T) {
	type TestUnionA struct{ A string }
	type TestUnionB struct{ B string }
	type TestUnion interface{ testUnion() }
	func(TestUnionA) testUnion() {}
	func(TestUnionB) testUnion() {}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			variant := UnionVariant{
				TypeFilter: gjson.JSON,
				Type:       reflect.TypeOf(TestUnionA{}),
			}
			if i%2 == 0 {
				variant.Type = reflect.TypeOf(TestUnionB{})
			}
			RegisterUnion(reflect.TypeOf((*TestUnion)(nil)).Elem(), "", variant)
		}(i)
	}
	wg.Wait()
}
