package pp

type PPConfig struct {
	ApiUrl      string
	IconUrl     string
	SecureLogin string
	SecretKey   string
}

// var ppConfig = PPConfig{
// 	ApiUrl:      "https://api.prerelease-env.biz/IntegrationService/v3/http/CasinoGameAPI",
// 	IconUrl:     "https://happycasino.prerelease-env.biz/game_pic/square/200/%s.png",
// 	SecureLogin: "hllgd_hollygod",
// 	SecretKey:   "C1Bd51884d654e10",
// }

var ppConfig = PPConfig{
	ApiUrl:      "https://jsgame.live/game/v1",
	IconUrl:     "https://happycasino.prerelease-env.biz/game_pic/square/200/%s.png",
	SecureLogin: "810f7f8748d39cdb8265fb95fa0ad462",
	SecretKey:   "e586168f7b781fa802020f4935974bb4",
}
