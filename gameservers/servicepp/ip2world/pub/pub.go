package ip2worldpub

import (
	"serve/comm/mq"
	"serve/comm/mux"
)

type Endpoint struct {
	IP   string
	Port int64
}

func GetEndpoint() (endpoint Endpoint, err error) {
	err = mq.Invoke("/ip2world/getEndpoint", mux.EmptyParams{}, &endpoint)
	return
}
