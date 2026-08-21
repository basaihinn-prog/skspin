package comm

const (
	Operator_Status_Normal = iota // 正常
	Operatar_Status_Stop          // 停用
)

const (
	User_Status_Normal = iota // 正常
	User_Status_Delete        // 删除
)

const (
	Day   = "Day"
	Month = "Month"
)

const (
	Jili = "jili"
	PG   = "pg"
	PP   = "pp"
)

const (
	GameStatus_Open        = 0 // 正常
	GameStatus_Maintenance = 1 // 维护中（列表中出现）
	GameStatus_Hide        = 2 // 隐藏（列表中不出现）
	GameStatus_Stop        = 3 // 关闭
)

const (
	GameType_Slot    = 0 // 拉霸游戏
	GameType_Fish    = 1 // 捕鱼游戏
	GameType_Poker   = 3 // 棋牌游戏
	GameType_CaiPiao = 4 // 彩票游戏
)

var CurrencyMap = map[string]string{
	"USD": "$",
	"EUR": "欧元",
	"JPY": "日元",
	"GBP": "英镑",
	"AUD": "澳元",
	"CNY": "¥",
	"INR": "印度卢比",
	"CAD": "加拿大元",
	"CHF": "瑞士法郎",
	"MXN": "墨西哥比索",
}

var PPGameWebTitle = map[string]string{
	"vs20fruitsw":     "Sweet Bonanza",
	"vs20sugarrush":   "Sugar Rush",
	"vs20olympx":      "Gates of Olympus 1000",
	"vs20starlightx":  "Starlight Princess 1000",
	"vs20starlight":   "Starlight Princess",
	"vs25scarabqueen": "John Hunter and the Tomb of the Scarab Queen",
	"vs20olympgate":   "Gates of Olympus",
	"vs20sbxmas":      "Sweet Bonanza Xmas",
	"vs9aztecgemsdx":  "Aztec Gems Deluxe",
	"vs40wildwest":    "Wild West Gold",
	"vswaysbufking":   "Buffalo King Megaways",
	"vs20ninjapower":  "Power of Ninja",
	"vs20procount":    "Wisdom of Athena",
	"vs5aztecgems":    "Aztec Gems",
	"vs20saiman":      "Saiyan Mania",
	"vs1024mahjwins":  "Mahjong Wins",
	"vs20candybltz2":  "Candy Blitz Bombs",
	"vswayswildwest":  "Wild West Gold Megaways",
	"vs20pbonanza":    "Pyramid Bonanza",
	"vs20olympdice":   "Gates of Olympus Dice",
	"vs10threestar":   "Three Star Fortune",
	"vs1024mjwinbns":  "Mahjong Wins Bonus",
	"vs7776aztec":     "Aztec Bonanza",
	"vs20portals":     "Fire Portals",
	"vs20bison":       "Release the Bison",
	"vs40stckwldlvl":  "Ripe Rewards",
	"vs20gatotx":      "Gates of Gatot Kaca 1000",
	"vs20doghouse2":   "The Dog House - Dog or Alive",
	"vs20tweethouse":  "The Tweety House",
	"vs20sugrux":      "Sugar Rush Xmas",
	"vs10firestrike":  "Fire Strike",
	"vs20bonzgold":    "Bonanza Gold",
	"vs10bblotgl":     "Big Bass - Secrets of the Golden Lake",
	"vs15godsofwar":   "Zeus vs Hades - Gods of War",
	"vswayswest":      "Mystic Chief",
	"vswaysmegahays":  "Barnyard Megahays Megaways",
	"vs20cleocatra":   "Cleocatra",
	"vs20fruitparty":  "Fruit Party",
	"vs20cbrhst":      "Cyber Heist",
	"vs20stickypos":   "Ice Lobster",
	"vs20doghousemh":  "The Dog House Multihold",
	"vs20gatotgates":  "Gates of Gatot Kaca",
	"vswaysmadame":    "Madame Destiny Megaways",
	"vs15samurai4":    "Rise of Samurai 4",
	"vs10bburger":     "Big Burger Load it up with Xtra Cheese",
	"vs50dmdcascade":  "Diamond Cascade",
	"vs10bbfloats":    "Big Bass Floats my Boat",
	"vs20clustext":    "Gears of Horus",
	"vs1tigers":       "Triple Tigers",
	"vs25pandagold":   "Panda's Fortune",
	"vs20schristmas":  "Starlight Christmas",
	"vswayswildgang":  "The Wild Gang",
	"vs10cowgold":     "Cowboys Gold",
	"vs20cjcluster":   "Candy Jar Clusters",
	"vs20goldclust":   "Rabbit Garden",
	"vswaysstampede":  "Fire Stampede",
	"vs20multiup":     "Wheel O'Gold",
	"vs20candyblitz":  "Candy Blitz",
	"vswaysstrlght":   "Fortunes of Aztec",
	"vswaysoldminer":  "Old Gold Miner Megaways",
	"vs20bnnzdice":    "Sweet Bonanza Dice",
	"vs25goldparty":   "Gold Party",
	"vs25xmasparty":   "Penguins Christmas Party Time",
}

var OldPPGame = map[string]string{
	"vs20fruitsw":    "Sweet Bonanza",
	"vs20sugarrush":  "Sugar Rush",
	"vs20olympx":     "Gates of Olympus 1000",
	"vs20starlightx": "Starlight Princess 1000",
	"vs20starlight":  "Starlight Princess",
}
