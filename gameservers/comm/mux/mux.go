package mux

import (
	"log"
	"log/slog"
	"net/http"
	"serve/comm/lazy"
	"sort"
)

var (
	DefaultRpcMux = &Mux{
		handlers: map[string]*PHandler{},
		ptrMap:   map[uintptr]HttpHandler{},
		ServeMux: http.DefaultServeMux,
	}
)

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

type GetArg0Fn = func(r *http.Request) (any, error)

type Mux struct {
	handlers map[string]*PHandler
	ptrMap   map[uintptr]HttpHandler
	arr      []*PHandler
	ServeMux *http.ServeMux
}

func NewMux() *Mux {
	return &Mux{
		handlers: map[string]*PHandler{},
		ptrMap:   map[uintptr]HttpHandler{},
		ServeMux: http.NewServeMux(),
	}
}

func (m *Mux) Add(data *PHandler) *PHandler {
	Assert(len(m.arr) == 0)

	_, ok := m.handlers[data.Path]
	Assert(!ok, "重复的path:"+data.Path)

	m.handlers[data.Path] = data
	return data
}

func (m *Mux) ToArr() []*PHandler {
	if len(m.arr) == 0 {
		arr := make([]*PHandler, 0, len(m.handlers))

		for _, v := range m.handlers {
			arr = append(arr, v)
		}

		sort.Slice(arr, func(i, j int) bool {
			x := arr[i]
			y := arr[j]
			if x.Kind == y.Kind {
				return x.Path < y.Path
			}
			return x.Kind < y.Kind
		})

		m.arr = arr
	}
	return m.arr
}

// var DefaultRpcMux = NewMux()

/*
func RegObj(ns string, obj any) {
	v := reflect.ValueOf(obj)

	typ := v.Type()
	nmethod := typ.NumMethod()

	for i := 0; i < nmethod; i++ {
		method := typ.Method(i)
		methodv := v.Method(i)

		data := &PHandler{
			Path:    path.Join("/", ns, method.Name),
			Handler: methodv.Interface(),
			Kind:    ns,
			Class:   "rpc",
		}

		DefaultRpcMux.Add(data)

	}
}
*/

func (m *Mux) StartHttpServer(addr string) {
	// r := http.NewServeMux()
	r := m.ServeMux
	fs := http.FileServer(http.Dir("./api"))
	r.Handle("/api/", http.StripPrefix("/api/", fs))
	r.HandleFunc("/list_api", m.list_api)

	for _, h := range m.ToArr() {
		slog.Info("regist",
			"path", h.Path,
			"kind", h.Kind,
			"desc", h.Desc,
			"class", h.Class,
		)
		r.HandleFunc(h.Path, httpRpcWrapper(h))
	}

	slog.Info("listen",
		"addr", addr,
	)

	go func() {
		lazy.Init("")
		//证书
		certfile, _ := lazy.RouteFile.Get("certfile")
		keyfile, _ := lazy.RouteFile.Get("keyfile")
		err := http.ListenAndServeTLS(addr, certfile, keyfile, r)
		//err := http.ListenAndServe(addr, r)
		log.Printf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
	}()
}

func StartHttpServer(addr string) {
	DefaultRpcMux.StartHttpServer(addr)
}
