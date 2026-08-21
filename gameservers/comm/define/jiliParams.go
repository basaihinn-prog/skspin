package define

type JiliParams struct {
	Path string
	// Pid  int64
	Data []byte
}

func (ps JiliParams) GetPid() int64 {
	return 0
}
