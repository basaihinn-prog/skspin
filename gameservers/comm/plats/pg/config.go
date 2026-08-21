package pg

// type Config struct {
// 	PgSoftAPIDomain   string
// 	DataGrabAPIDomain string
// 	LaunchURL         string
// 	OperatorToken     string
// 	SecretKey         string
// 	Lang              string
// 	OperatorLaunchUrl string
// }

type PGConfig struct {
	PgSoftAPIDomain   string
	DataGrabAPIDomain string
	LaunchURL         string
	// 爬虫专用
	ClientApiURL      string
	OperatorToken     string
	SecretKey         string
	Lang              string
	OperatorLaunchUrl string
	Currency          string
}

// var pgcfg = &PGConfig{
// 	PgSoftAPIDomain:   "https://api.pg-bo.me/external",
// 	DataGrabAPIDomain: "https://api.pg-bo.me/external-datagrabber",
// 	LaunchURL:         "https://m.pg-demo.com",
// 	ClientApiURL:      "https://api.pg-demo.com",
// 	OperatorToken:     "18aa254b73981ec0dc42e24e6b9f03ce",
// 	SecretKey:         "976de6c68077668be9dba2f64c61ee7d",
// 	Lang:              "en",
// 	OperatorLaunchUrl: "https://plats.kafa010.com/pg/launch.html",
// }

var pgcfg = &PGConfig{
	PgSoftAPIDomain:   "https://api.pg-bo.me/external",
	DataGrabAPIDomain: "https://api.pg-bo.me/external-datagrabber",
	LaunchURL:         "https://m.pg-demo.com",
	ClientApiURL:      "https://api.pg-demo.com",
	OperatorToken:     "I-0dd08edd9bef4fd18d26f4e4ed747662",
	SecretKey:         "8fccb08ecede4fa990d2320c7536003e",
	Lang:              "en",
	OperatorLaunchUrl: "https://plats.rpgamestest.com/pg/launch.html",
	Currency:          "BRL",
}

func GetConfig() (cfg *PGConfig) {
	// cfg = &gamedata.Get().PG
	return pgcfg
}
