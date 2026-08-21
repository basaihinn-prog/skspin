package api

import (
	"fmt"
	"net/http"
)

// https://api.kafa010.com/web-api/game-proxy/v2/GameName/Get?traceId=YSFGZL13
func gameName(w http.ResponseWriter, r *http.Request) {
	// lang=en&btt=2&atk=4952EF30-CADC-4A0F-A2AF-5D6FEAFCBB7D&pf=1&gid=39
	// ret := `{"dt":{"31":"Baccarat Deluxe"},"err":null}`
	//ret := `{"dt":{"0":"Lobby","1":"Honey Trap of Diao Chan","2":"Gem Saviour","3":"Fortune Gods","6":"Medusa 2: The Quest of Perseus","7":"Medusa 1: The Curse of Athena","18":"Hood vs Wolf","20":"Reel Love","24":"Win Win Won","25":"Plushie Frenzy","26":"Tree of Fortune","28":"Hotpot","29":"Dragon Legend","33":"Hip Hop Panda","34":"Legend of Hou Yi","35":"Mr. Hallow-Win","36":"Prosperity Lion","37":"Santa's Gift Rush","38":"Gem Saviour Sword","39":"Piggy Gold","40":"Jungle Delight","41":"Symbols Of Egypt","42":"Ganesha Gold","44":"Emperor's Favour","48":"Double Fortune","50":"Journey to the Wealth","53":"The Great Icescape","54":"Captain's Bounty","57":"Dragon Hatch","58":"Vampire's Charm","59":"Ninja vs Samurai","60":"Leprechaun Riches","61":"Flirting Scholar","62":"Gem Saviour Conquest","63":"Dragon Tiger Luck","64":"Muay Thai Champion","65":"Mahjong Ways","67":"Shaolin Soccer","68":"Fortune Mouse","69":"Bikini Paradise","70":"Candy Burst","71":"Cai Shen Wins","73":"Egypt's Book of Mystery","74":"Mahjong Ways 2","75":"Ganesha Fortune","79":"Dreams of Macau","80":"Circus Delight","82":"Phoenix Rises","83":"Wild Fireworks","84":"Queen of Bounty","85":"Genie's 3 Wishes","86":"Galactic Gems","87":"Treasures of Aztec","88":"Jewels of Prosperity","89":"Lucky Neko","90":"Secrets of Cleopatra","91":"Guardians of Ice & Fire","92":"Thai River Wonders","93":"Opera Dynasty","94":"Bali Vacation","95":"Majestic Treasures","97":"Jack Frost's Winter","98":"Fortune Ox","100":"Candy Bonanza","101":"Rise of Apollo","102":"Mermaid Riches","103":"Crypto Gold","104":"Wild Bandito","105":"Heist Stakes","106":"Ways of the Qilin","107":"Legendary Monkey King","108":"Buffalo Win","110":"Jurassic Kingdom","112":"Oriental Prosperity","113":"Raider Jane's Crypt of Fortune","114":"Emoji Riches","115":"Supermarket Spree","117":"Cocktail Nights","118":"Mask Carnival","119":"Spirited Wonders","120":"The Queen's Banquet","121":"Destiny of Sun & Moon","122":"Garuda Gems","123":"Rooster Rumble","124":"Battleground Royale","125":"Butterfly Blossom","126":"Fortune Tiger","127":"Speed Winner","128":"Legend of Perseus","129":"Win Win Fish Prawn Crab","130":"Lucky Piggy","132":"Wild Coaster","135":"Wild Bounty Showdown","1312883":"Prosperity Fortune Tree","1338274":"Totem Wonders","1340277":"Asgardian Rising","1368367":"Alchemy Gold","1372643":"Diner Delights","1381200":"Hawaiian Tiki","1397455":"Fruity Candy","1402846":"Midas Fortune","1418544":"Bakery Bonanza","1420892":"Rave Party Fever","1432733":"Mystical Spirits","1448762":"Songkran Splash","1451122":"Dragon Hatch2","1473388":"Cruise Royale","1489936":"Ultimate Striker","1513328":"Super Golf Drive","1529867":"Ninja Raccoon Frenzy","1543462":"Fortune Rabbit","1555350":"Forge of Wealth","1568554":"Wild Heist Cashout","1572362":"Gladiator's Glory","1580541":"Mafia Mayhem","1594259":"Safari Wilds","1601012":"Lucky Clover Lady","1615454":"Werewolf's Hunt","1655268":"Tsar Treasures","1671262":"Gemstones Gold","1682240":"Cash Mania","1695365":"Fortune Dragon"},"err":null}`
	ret := `{"dt":{"31":"Baccarat Deluxe"},"err":null}`
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ret))
}

func gameRule(w http.ResponseWriter, r *http.Request) {
	// https://api.pg-demo.com/web-api/game-proxy/v2/GameRule/Get?traceId=TGVFVQ16
	// btt=2&gid=39&atk=73A1554D-35C4-4DA0-B6BC-478FABC1C6AD&pf=1
	r.ParseForm()
	gid := r.PostForm.Get("gid")
	rtp := "96.81"
	if len(gid) != 0 {
		rtp = gameRTPMap[gid]
	}
	//ret := `{"dt":{"rtp":{"Default":{"min":133.29,"max":133.29}},"ows":{"itare":false,"tart":0,"igare":false,"gart":0},"jws":null},"err":null}`
	ret := fmt.Sprintf(`{"dt":{"rtp":{"Default":{"min":%v,"max":%v}},"ows":{"itare":false,"tart":0,"igare":false,"gart":0},"jws":null},"err":null}`, rtp, rtp)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ret))
}

func getSocialInitConfig(w http.ResponseWriter, r *http.Request) {
	// https://api.pg-demo.com/web-api/game-proxy/v2/GameRule/Get?traceId=TGVFVQ16
	// btt=2&gid=39&atk=73A1554D-35C4-4DA0-B6BC-478FABC1C6AD&pf=1
	ret := `{"dt":{"countries":null,"achievements":null,"levelActionPermissions":null,"gameChatEmoticonTemplateDatas":null,"gameSpecificThresholdId":null,"configurationSetting":null,"onlinePlayerCounts":null,"playerFavouriteGamesDatas":null,"playerProfileInfo":null,"levelBackground":null,"onlinePlayerCount":null},"err":null}`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ret))
}

// https://api.pg-demo.com/web-api/game-proxy/v2/Resources/GetByReferenceIdsResourceTypeIds?traceId=OMHCFR21
// atk=83F8C938-1EAD-4E81-8BBE-CC8A6C304A68&pf=1&btt=1&du=https%3A%2F%2Fstatic-pg.kafa010.com&otk=18aa254b73981ec0dc42e24e6b9f03ce&rids=39&rtids=19
//{"dt":[{"rid":39,"rtid":19,"url":"https://static-pg.kafa010.com/pages/static/image/en/SocialGameExtraSmall/39/Piggy_Gold-0d711f06.png","l":"en-US","ut":"2019-09-27T10:57:10"}],"err":null}

func getByReferenceIdsResourceTypeIds(w http.ResponseWriter, r *http.Request) {
	// ret := `{"dt":[{"rid":%v,"rtid":%v,"url":"","l":"en-US","ut":"2019-09-27T10:57:10"}],"err":null}`

	// ret = fmt.Sprintf(ret, r.FormValue("rids"), r.FormValue("rtids"))
	ret := `{"dt":[],"err":null}`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ret))
}

// https://api.pg-demo.com/web-api/game-proxy/v2/Resources/GetByResourcesTypeIds?traceId=GFHQCC21
func getByResourcesTypeIds(w http.ResponseWriter, r *http.Request) {
	ret := `{"dt":[],"err":null}`

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ret))
}

// https://api.pg-demo.com/web-api/game-proxy/v2/Resources/GetByReferenceIdsResourceTypeIds?traceId=HACHJT21

var gameRTPMap = map[string]string{
	"1":       "97.13",
	"2":       "97.05",
	"3":       "96.99",
	"6":       "97.17",
	"7":       "97.19",
	"17":      "96.84",
	"18":      "96.85",
	"20":      "97.01",
	"24":      "97.11",
	"25":      "96.87",
	"26":      "96.91",
	"28":      "96.93",
	"29":      "96.84",
	"33":      "96.84",
	"34":      "97.06",
	"35":      "96.81",
	"36":      "96.99",
	"37":      "97.03",
	"38":      "96.97",
	"39":      "96.89",
	"40":      "97.04",
	"41":      "96.95",
	"42":      "97.11",
	"44":      "97.05",
	"48":      "97.08",
	"50":      "96.98",
	"53":      "96.91",
	"54":      "96.83",
	"57":      "97.08",
	"58":      "97.19",
	"59":      "97.15",
	"60":      "96.93",
	"61":      "97.06",
	"62":      "96.98",
	"63":      "96.82",
	"64":      "96.92",
	"65":      "96.81",
	"67":      "97.14",
	"68":      "97.06",
	"69":      "97.19",
	"70":      "96.95",
	"73":      "96.99",
	"74":      "97.11",
	"75":      "96.86",
	"79":      "96.89",
	"80":      "96.99",
	"82":      "97.12",
	"83":      "97.06",
	"84":      "96.91",
	"85":      "97.03",
	"86":      "97.02",
	"87":      "96.95",
	"88":      "97.12",
	"89":      "97.06",
	"90":      "97.18",
	"91":      "96.97",
	"92":      "96.99",
	"93":      "96.83",
	"94":      "96.87",
	"95":      "97.06",
	"96":      "96.97",
	"98":      "97.08",
	"100":     "96.99",
	"101":     "96.86",
	"102":     "97.14",
	"103":     "97.11",
	"104":     "96.82",
	"105":     "97.13",
	"106":     "97.11",
	"107":     "96.82",
	"108":     "97.12",
	"109":     "97.14",
	"110":     "97.06",
	"111":     "96.88",
	"112":     "97.11",
	"113":     "97.11",
	"114":     "96.93",
	"115":     "97.09",
	"116":     "97.16",
	"117":     "96.95",
	"118":     "96.93",
	"119":     "96.81",
	"120":     "97.11",
	"125":     "97.16",
	"126":     "96.99",
	"127":     "97.11",
	"128":     "97.12",
	"129":     "97.08",
	"130":     "97.08",
	"132":     "96.88",
	"135":     "97.11",
	"1312883": "96.85",
	"1338274": "96.86",
	"1340277": "96.81",
	"1368367": "97.06",
	"1372643": "97.14",
	"1381200": "97.11",
	"1397455": "97.16",
	"1402846": "97.18",
	"1418544": "97.18",
	"1420892": "97.08",
	"1432733": "96.93",
	"1448762": "96.86",
	"1451122": "96.88",
	"1473388": "97.11",
	"1489936": "97.19",
	"1508783": "97.11",
	"1513328": "97.06",
	"1529867": "97.09",
	"1543462": "96.97",
	"1555350": "96.85",
	"1568554": "97.06",
	"1572362": "96.82",
	"1580541": "96.94",
	"1594259": "97.03",
	"1601012": "97.15",
	"1615454": "96.99",
	"1655268": "97.11",
	"1671262": "96.92",
	"1682240": "97.11",
	"1695365": "96.96",
	"1879752": "96.75",
}
