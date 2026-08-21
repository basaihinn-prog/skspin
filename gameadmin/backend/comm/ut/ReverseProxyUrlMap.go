package ut

import (
	"cmp"
	"net/url"
	"slices"
	"strings"
)

type ReverseProxyUrlMapUnit struct {
	Prefix string
	Url    *url.URL
}

// 前缀匹配, 长度优先的反向代理选择器
type ReverseProxyUrlMap []*ReverseProxyUrlMapUnit

func NewReverseProxyUrlMap(c map[string]string) (m ReverseProxyUrlMap, err error) {
	m = make(ReverseProxyUrlMap, 0, len(c))
	for host, remote := range c {
		var u *url.URL
		u, err = url.Parse(remote)
		if err != nil {
			return
		}

		m = append(m, &ReverseProxyUrlMapUnit{
			Prefix: host,
			Url:    u,
		})
		// m[host] = u
	}

	// sort.Slice()
	slices.SortFunc(m, func(x, y *ReverseProxyUrlMapUnit) int {
		if len(x.Prefix) > len(y.Prefix) {
			return -1
		}
		if len(x.Prefix) < len(y.Prefix) {
			return 1
		}

		return cmp.Compare(x.Prefix, y.Prefix)
	})

	return
}

func (m ReverseProxyUrlMap) Get(hostname string) *url.URL {
	for _, v := range m {
		if strings.HasPrefix(hostname, v.Prefix) {
			var cp = *v.Url
			return &cp
		}
	}
	return nil
}
