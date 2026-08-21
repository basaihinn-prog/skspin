package mux

const (
	ClassHttp    = "http"
	ClassRpc     = "rpc"
	ClassMsgProc = "msgproc"
)

type PHandler struct {
	Path         string
	Handler      interface{}
	Desc         string
	Kind         string
	ParamsSample interface{}

	// http, rpc, msgproc
	Class string

	// 仅开发环境有效
	OnlyDev bool

	GetArg0 GetArg0Fn
}

func (h *PHandler) SetOnlyDev() *PHandler {
	h.OnlyDev = true
	return h
}
