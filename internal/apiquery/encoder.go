package apiquery

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anomalyco/opencode-sdk-go/internal/param"
	"github.com/anomalyco/opencode-sdk-go/internal/timeformat"
)

var encoders sync.Map // map[reflect.Type]encoderFunc

type encoder struct {
	dateFormat string
	root       bool
	settings   QuerySettings
}

type encoderFunc func(key string, value reflect.Value) ([]Pair, error)

type encoderField struct {
	tag parsedStructTag
	fn  encoderFunc
	idx []int
}

type encoderEntry struct {
	reflect.Type
	dateFormat string
	root       bool
	settings   QuerySettings
}

type Pair struct {
	key   string
	value string
}

func (e *encoder) typeEncoder(t reflect.Type) encoderFunc {
	entry := encoderEntry{
		Type:       t,
		dateFormat: e.dateFormat,
		root:       e.root,
		settings:   e.settings,
	}

	if fi, ok := encoders.Load(entry); ok {
		return fi.(encoderFunc)
	}

	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoders.LoadOrStore(entry, encoderFunc(func(key string, v reflect.Value) ([]Pair, error) {
		wg.Wait()
		return f(key, v)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	// Use defer to ensure wg.Done() is called even if newTypeEncoder panics,
	// preventing deadlock for concurrent goroutines waiting on wg.Wait().
	f = e.newTypeEncoder(t)
	defer wg.Done()
	encoders.Store(entry, f)
	return f
}

func marshalerEncoder(key string, value reflect.Value) ([]Pair, error) {
	s, err := value.Interface().(json.Marshaler).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("apiquery: MarshalJSON failed for type %T: %w", value.Interface(), err)
	}
	return []Pair{{key, string(s)}}, nil
}

func (e *encoder) newTypeEncoder(t reflect.Type) encoderFunc {
	if t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
		return e.newTimeTypeEncoder(t)
	}
	if !e.root && t.Implements(reflect.TypeOf((*json.Marshaler)(nil)).Elem()) {
		return marshalerEncoder
	}
	e.root = false
	switch t.Kind() {
	case reflect.Pointer:
		encoder := e.typeEncoder(t.Elem())
		return func(key string, value reflect.Value) ([]Pair, error) {
			if !value.IsValid() || value.IsNil() {
				return nil, nil
			}
			return encoder(key, value.Elem())
		}
	case reflect.Struct:
		return e.newStructTypeEncoder(t)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return e.newArrayTypeEncoder(t)
	case reflect.Map:
		return e.newMapEncoder(t)
	case reflect.Interface:
		return e.newInterfaceEncoder()
	default:
		return e.newPrimitiveTypeEncoder(t)
	}
}

func (e *encoder) newStructTypeEncoder(t reflect.Type) encoderFunc {
	if t.Implements(reflect.TypeOf((*param.FieldLike)(nil)).Elem()) {
		return e.newFieldTypeEncoder(t)
	}

	encoderFields := []encoderField{}

	var collectEncoderFields func(r reflect.Type, index []int)
	collectEncoderFields = func(r reflect.Type, index []int) {
		for i := 0; i < r.NumField(); i++ {
			idx := append(index, i)
			field := t.FieldByIndex(idx)
			if !field.IsExported() {
				continue
			}
			if field.Anonymous {
				collectEncoderFields(field.Type, idx)
				continue
			}
			ptag, ok := parseQueryStructTag(field)
			if !ok {
				continue
			}

			if ptag.name == "-" && !ptag.inline {
				continue
			}

			dateFormat, ok := parseFormatStructTag(field)
			oldFormat := e.dateFormat
			if ok {
				switch dateFormat {
				case "date-time":
					e.dateFormat = time.RFC3339
				case "date":
					e.dateFormat = timeformat.Date
				}
			}
			encoderFields = append(encoderFields, encoderField{ptag, e.typeEncoder(field.Type), idx})
			e.dateFormat = oldFormat
		}
	}
	collectEncoderFields(t, []int{})

	return func(key string, value reflect.Value) ([]Pair, error) {
		var pairs []Pair
		for _, ef := range encoderFields {
			var subkey = e.renderKeyPath(key, ef.tag.name)
			if ef.tag.inline {
				subkey = key
			}

			field := value.FieldByIndex(ef.idx)
			subPairs, err := ef.fn(subkey, field)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, subPairs...)
		}
		return pairs, nil
	}
}

func (e *encoder) newMapEncoder(t reflect.Type) encoderFunc {
	keyEncoder := e.typeEncoder(t.Key())
	elementEncoder := e.typeEncoder(t.Elem())
	return func(key string, value reflect.Value) ([]Pair, error) {
		var pairs []Pair
		iter := value.MapRange()
		for iter.Next() {
			encodedKey, err := keyEncoder("", iter.Key())
			if err != nil {
				return nil, fmt.Errorf("encoding map key: %w", err)
			}
			if len(encodedKey) != 1 {
				return nil, fmt.Errorf("apiquery: map key must encode to exactly one pair, got %d (non-primitive map keys not supported)", len(encodedKey))
			}
			subkey := encodedKey[0].value
			keyPath := e.renderKeyPath(key, subkey)
			subPairs, err := elementEncoder(keyPath, iter.Value())
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, subPairs...)
		}
		return pairs, nil
	}
}

func (e *encoder) renderKeyPath(key string, subkey string) string {
	if len(key) == 0 {
		return subkey
	}
	if e.settings.NestedFormat == NestedQueryFormatDots {
		return fmt.Sprintf("%s.%s", key, subkey)
	}
	return fmt.Sprintf("%s[%s]", key, subkey)
}

func (e *encoder) newArrayTypeEncoder(t reflect.Type) encoderFunc {
	switch e.settings.ArrayFormat {
	case ArrayQueryFormatComma:
		innerEncoder := e.typeEncoder(t.Elem())
		return func(key string, v reflect.Value) ([]Pair, error) {
			elements := []string{}
			for i := 0; i < v.Len(); i++ {
				pairs, err := innerEncoder("", v.Index(i))
				if err != nil {
					return nil, err
				}
				for _, pair := range pairs {
					elements = append(elements, pair.value)
				}
			}
			if len(elements) == 0 {
				return nil, nil
			}
			return []Pair{{key, strings.Join(elements, ",")}}, nil
		}
	case ArrayQueryFormatRepeat:
		innerEncoder := e.typeEncoder(t.Elem())
		return func(key string, value reflect.Value) ([]Pair, error) {
			var pairs []Pair
			for i := 0; i < value.Len(); i++ {
				subPairs, err := innerEncoder(key, value.Index(i))
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, subPairs...)
			}
			return pairs, nil
		}
	case ArrayQueryFormatIndices:
		return func(key string, value reflect.Value) ([]Pair, error) {
			return nil, fmt.Errorf("apiquery: array format 'indices' is not supported")
		}
	case ArrayQueryFormatBrackets:
		innerEncoder := e.typeEncoder(t.Elem())
		return func(key string, value reflect.Value) ([]Pair, error) {
			var pairs []Pair
			for i := 0; i < value.Len(); i++ {
				subPairs, err := innerEncoder(key+"[]", value.Index(i))
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, subPairs...)
			}
			return pairs, nil
		}
	default:
		panic(fmt.Sprintf("apiquery: unknown ArrayFormat value: %d", e.settings.ArrayFormat))
	}
}

func (e *encoder) newPrimitiveTypeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {
	case reflect.Pointer:
		inner := t.Elem()

		innerEncoder := e.newPrimitiveTypeEncoder(inner)
		return func(key string, v reflect.Value) ([]Pair, error) {
			if !v.IsValid() || v.IsNil() {
				return nil, nil
			}
			return innerEncoder(key, v.Elem())
		}
	case reflect.String:
		return func(key string, v reflect.Value) ([]Pair, error) {
			return []Pair{{key, v.String()}}, nil
		}
	case reflect.Bool:
		return func(key string, v reflect.Value) ([]Pair, error) {
			if v.Bool() {
				return []Pair{{key, "true"}}, nil
			}
			return []Pair{{key, "false"}}, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(key string, v reflect.Value) ([]Pair, error) {
			return []Pair{{key, strconv.FormatInt(v.Int(), 10)}}, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(key string, v reflect.Value) ([]Pair, error) {
			return []Pair{{key, strconv.FormatUint(v.Uint(), 10)}}, nil
		}
	case reflect.Float32, reflect.Float64:
		return func(key string, v reflect.Value) ([]Pair, error) {
			return []Pair{{key, strconv.FormatFloat(v.Float(), 'f', -1, 64)}}, nil
		}
	case reflect.Complex64, reflect.Complex128:
		bitSize := 64
		if t.Kind() == reflect.Complex128 {
			bitSize = 128
		}
		return func(key string, v reflect.Value) ([]Pair, error) {
			return []Pair{{key, strconv.FormatComplex(v.Complex(), 'f', -1, bitSize)}}, nil
		}
	default:
		return func(key string, v reflect.Value) ([]Pair, error) {
			return nil, nil
		}
	}
}

func (e *encoder) newFieldTypeEncoder(t reflect.Type) encoderFunc {
	f, _ := t.FieldByName("Value")
	enc := e.typeEncoder(f.Type)

	return func(key string, value reflect.Value) ([]Pair, error) {
		present := value.FieldByName("Present")
		if !present.Bool() {
			return nil, nil
		}
		null := value.FieldByName("Null")
		if null.Bool() {
			return nil, nil
		}
		raw := value.FieldByName("Raw")
		if !raw.IsNil() {
			return e.typeEncoder(raw.Type())(key, raw)
		}
		return enc(key, value.FieldByName("Value"))
	}
}

func (e *encoder) newTimeTypeEncoder(t reflect.Type) encoderFunc {
	format := e.dateFormat
	return func(key string, value reflect.Value) ([]Pair, error) {
		return []Pair{{
			key,
			value.Convert(reflect.TypeOf(time.Time{})).Interface().(time.Time).Format(format),
		}}, nil
	}
}

func (e encoder) newInterfaceEncoder() encoderFunc {
	return func(key string, value reflect.Value) ([]Pair, error) {
		value = value.Elem()
		if !value.IsValid() {
			return nil, nil
		}
		return e.typeEncoder(value.Type())(key, value)
	}
}
