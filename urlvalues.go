package valuesmap

import (
	"encoding/json"
	"net/url"
	"strings"
)

type URLValuesGetter struct {
	src url.Values
	Key string
}

func (g *URLValuesGetter) Get() interface{} {
	return g.src.Get(g.Key)
}

func (g *URLValuesGetter) MarshalJSON() (data []byte, err error) {
	return json.Marshal(g.Get())
}

func ParseURLValues(src url.Values, prefix ...string) map[string]interface{} {
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
			m[k3] = &URLValuesGetter{src, k}
		} else {
			m[k2] = &URLValuesGetter{src, k}
		}
	}
	return result
}
