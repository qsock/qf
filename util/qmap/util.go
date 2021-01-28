package qmap

import (
	"fmt"
	"github.com/qsock/qf/util/qcast"
	"reflect"
	"strings"
)

// MergeStringMap merge two map
func MergeStringMap(dest, src map[string]interface{}) {
	for sk, sv := range src {
		tv, ok := dest[sk]
		if !ok {
			// val不存在时，直接赋值
			dest[sk] = sv
			continue
		}

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if svType != tvType {
			fmt.Println("continue, type is different")
			continue
		}

		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			tsv := sv.(map[interface{}]interface{})
			ssv := ToMapStringInterface(tsv)
			stv := ToMapStringInterface(ttv)
			MergeStringMap(stv, ssv)
			dest[sk] = stv
		case map[string]interface{}:
			MergeStringMap(ttv, sv.(map[string]interface{}))
			dest[sk] = ttv
		default:
			dest[sk] = sv
		}
	}
}

// ToMapStringInterface cast map[interface{}]interface{} to map[string]interface{}
func ToMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

// InsensitiviseMap insensitivise map
func InsensitiviseMap(m map[string]interface{}) {
	for key, val := range m {
		switch val.(type) {
		case map[interface{}]interface{}:
			InsensitiviseMap(qcast.ToStringMap(val))
		case map[string]interface{}:
			InsensitiviseMap(val.(map[string]interface{}))
		}

		lower := strings.ToLower(key)
		if key != lower {
			delete(m, key)
		}
		m[lower] = val
	}
}

// DeepSearchInMap deep search in map
func DeepSearchInMap(m map[string]interface{}, paths ...string) map[string]interface{} {
	//深度拷贝
	mtmp := make(map[string]interface{})
	for k, v := range m {
		mtmp[k] = v
	}
	for _, k := range paths {
		m2, ok := mtmp[k]
		if !ok {
			m3 := make(map[string]interface{})
			mtmp[k] = m3
			mtmp = m3
			continue
		}

		m3, err := qcast.ToStringMapE(m2)
		if err != nil {
			m3 = make(map[string]interface{})
			mtmp[k] = m3
		}
		// continue search
		mtmp = m3
	}
	return mtmp
}
