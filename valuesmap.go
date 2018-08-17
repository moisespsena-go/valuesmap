package valuesmap

import (
	"encoding/json"
	"reflect"
	"strings"
)

var GetterType = reflect.TypeOf(new(Getter)).Elem()

type Getter interface {
	Get() interface{}
}

type MapGetter struct {
	src map[string]interface{}
	Key string
}

func (g *MapGetter) Get() interface{} {
	return (g.src)[g.Key]
}

func (g *MapGetter) MarshalJSON() (data []byte, err error) {
	return json.Marshal(g.Get())
}

func ParseMap(src map[string]interface{}, prefix ...string) map[string]interface{} {
	var prefx string
	if len(prefix) > 0 && prefix[0] != "" {
		prefx = prefix[0]
	}
	result := map[string]interface{}{}
	var (
		m      map[string]interface{}
		ok     bool
		k2, k3 string
		parts  []string
	)
	for k := range src {
		if prefx != "" {
			if strings.HasPrefix(k, prefx+".") {
				k2 = strings.TrimPrefix(k, prefx+".")
			} else {
				continue
			}
		} else {
			k2 = k
		}
		m = result
		parts = strings.Split(k2, ".")
		if lp := len(parts); lp > 1 {
			parts, k3 = parts[0:lp-1], parts[lp-1]
			for _, p := range parts {
				if _, ok = m[p]; !ok {
					m[p] = map[string]interface{}{}
				}
				m = m[p].(map[string]interface{})
			}
			m[k3] = &MapGetter{src, k}
		} else {
			m[k2] = &MapGetter{src, k}
		}
	}
	return result
}
