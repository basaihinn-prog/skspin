package maintain

type Maintain struct {
	ShutDown       bool   `json:"shutdown"`
	Message        string `json:"message"`
	DatabaseReload bool   `json:"database_reload"`
}
