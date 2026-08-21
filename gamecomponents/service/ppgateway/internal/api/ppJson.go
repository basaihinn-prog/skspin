package api

import (
	_ "embed"
	"net/http"
)

//go:embed ppdata/PromoActive.json
var promoActiveJson []byte

//go:embed  ppdata/MinilobbyGames.json
var minilobbyGamesJson []byte

//go:embed  ppdata/PromoTournamentDetails.json
var promoTournamentDetailsJson []byte

//go:embed  ppdata/PromoRaceDetails.json
var promoRaceDetailsJson []byte

//go:embed  ppdata/PromoRacePrizes.json
var promoRacePrizesJson []byte

func PromoFrbAvailable(w http.ResponseWriter, r *http.Request) {
	buf := `<PromoDisabledResponse><error>4</error><description>Promo system is disabled</description></PromoDisabledResponse>`
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/xml;charset=utf-8")
	w.Write([]byte(buf))
}

func PromoActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(promoActiveJson)
}

func MinilobbyGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(minilobbyGamesJson)
}

func PromoTournamentDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(promoTournamentDetailsJson)
}

func PromoRaceDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(promoRaceDetailsJson)
}

func PromoTournamentScores(w http.ResponseWriter, r *http.Request) {
	buf := `{"error":0,"description":"OK"}`
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(buf))
}

func PromoRacePrizes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(promoRacePrizesJson)
}

func Unread(w http.ResponseWriter, r *http.Request) { // 先写死
	buf := `<EmptyGetUnreadAnnouncementsResponseDTO><error>0</error><description>OK</description><announcements/></EmptyGetUnreadAnnouncementsResponseDTO>`
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/xml;charset=utf-8")
	w.Write([]byte(buf))
}
