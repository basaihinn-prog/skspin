package platpg

import (
	"game/service/pggateway/internal/gamedata"
)

// type Config struct {
// 	PgSoftAPIDomain   string
// 	DataGrabAPIDomain string
// 	LaunchURL         string
// 	OperatorToken     string
// 	SecretKey         string
// 	Lang              string
// 	OperatorLaunchUrl string
// }

var pgcfg = &gamedata.PG{
	PgSoftAPIDomain:   "https://api.pg-bo.me/external",
	DataGrabAPIDomain: "https://api.pg-bo.me/external-datagrabber",
	LaunchURL:         "https://m.pg-demo.com",
	ClientApiURL:      "https://api.pg-demo.com",
	OperatorToken:     "18aa254b73981ec0dc42e24e6b9f03ce",
	SecretKey:         "976de6c68077668be9dba2f64c61ee7d",
	Lang:              "en",
	OperatorLaunchUrl: "https://platpg.kafa010.com/pg/launch.html",
	PGLaunchUrl:       "https://m.pg-demo.com/39/index.html?ot=18aa254b73981ec0dc42e24e6b9f03ce&btt=1&l=en&ops=db76d6940594f28e1815d2558ac52c8a&or=23pqxqfz%3Dmd-abjl%3Dzlj&__hv=1fb13252",
}

func GetConfig() (cfg *gamedata.PG) {
	// cfg = &gamedata.Get().PG
	return pgcfg
}
