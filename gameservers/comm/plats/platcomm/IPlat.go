package platcomm

import (
	"encoding/json"
	"path"

	"serve/comm/define"
	"serve/comm/mux"
)

type IPlat interface {
	LaunchGame(uid string, game, lang string, useProxy bool) (url string, err error)
	FundTransferIn(uid string, amount float64) (status string)

	GetBalance(uid string) (balance float64, err error)
	FundTransferOut(uid string) (amount float64, status string)

	GetGameList() (games HotGames, err error)

	Regist(uid string) (err error)
}

var Plats = map[string]IPlat{}
var PPGameRelation = map[string]string{}

func Start() {
	for k, plat := range Plats {
		regPlat("plat/"+k, plat)
	}
	//init pp map
	PPGameRelation["vs20olympx"] = "4ae52ed2e1a8c353878ba65ed7791ac4"
	PPGameRelation["vs20fruitsw"] = "6dcaf78e4e23929cbe2deb3d1210928c"
	PPGameRelation["vs20sugarrush"] = "6aeb8420d03722cfd4b66d91e17758ba"
	PPGameRelation["vs20starlight"] = "be6b6890587ed84289fad941d99a3613"
	PPGameRelation["vs20starlightx"] = "09d08939279289a03b89f2f146a7f817"
	PPGameRelation["vs25scarabqueen"] = "6ea631cad8e04b5596da2ab771a81198"
	PPGameRelation["vs5aztecgems"] = "05d693169ee6aeec45d6328c7fc57a43"
	PPGameRelation["vs10bblotgl"] = "aaea3a95f3e5ce36e4beea92c71ed912"
	PPGameRelation["vs20fruitsw"] = "6dcaf78e4e23929cbe2deb3d1210928c"
	PPGameRelation["vs20doghouse2"] = "82bead4e84621b9b553f6b7c0c5faa73"
	PPGameRelation["vs20fruitparty"] = "05a9876c537ff79898d661c5c01c5f9d"
	PPGameRelation["vs20procount"] = "bda320d9d63d12432d578b1f97855087"
	PPGameRelation["vs20sugarrushx"] = "3b2630042452d3e6b79054e4e96a2bf3"
	PPGameRelation["vs20sugrux"] = "8170339cc716c66d5b52779119fd718f"
	PPGameRelation["vs25goldparty"] = "fb6886b6f0ee1379fcadbdda398d3dc7"
	PPGameRelation["vs25pandagold"] = "4d530e5363907c3799f15bee86023202"
	PPGameRelation["vs40wildwest"] = "16ed3068536a6bdae156e3d98bb9a365"
	PPGameRelation["vswaysbufking"] = "0ea5af839234de02f6ddcf4c51f31f59"
	PPGameRelation["vs20santawonder"] = "d327a0dea94bd8a5fc090b4855051dff"
	PPGameRelation["vswaysfirewmw"] = "32a9b5e9958d37e859f8ae6f0705c095"
	PPGameRelation["vs20wolfie"] = "f6ac7090e2eda2d579dc6650acec42ad"
	PPGameRelation["vs20shootstars"] = "8e0755ec28b776e7c8d8763924152920"

	PPGameRelation["vs10bbsplxmas"] = "d78b41ecb69fe8370a561d2d478d980f"
	PPGameRelation["vs10bhallbnza2"] = "045f290df6c578256adb632ea1da485f"
	PPGameRelation["vs10fangfree"] = "3cde6f0025377dd4725f37a7e5e72ead"
	PPGameRelation["vs10jokerhot"] = "db243d8097f9e469dd308699dfad5c6a"
	PPGameRelation["vs10noodles"] = "dfdef0baca51fd11a62636cc9a197e5a"
	PPGameRelation["vs10txbigbass"] = "a1dfe25af8aa67b40e3f5c2b090a7f66"
	PPGameRelation["vs12bbb"] = "533b95b0d17cda431521139970d39073"
	PPGameRelation["vs12bbbxmas"] = "98f093dcaa2b0fa713d35f4594e9451e"
	PPGameRelation["vs20aztecgates"] = "3424ac1587cce058dea273a4e5a3cbf8"
	PPGameRelation["vs20clreacts"] = "3293f239dccef44f742e8d5241eb89c0"
	PPGameRelation["vs20drgbless"] = "7213d6d23f67e21cb0843c95ef4d5614"
	PPGameRelation["vs20mustanggld2"] = "5ef2e439010aeca067f3a6916c901f07"
	PPGameRelation["vs20procountx"] = "97a7bed99262b8c950be4bf1b173e23d"
	PPGameRelation["vs40wildrun"] = "b210bf0a7ae17ae28c9ca2ebd6f76b05"
	PPGameRelation["vs50juicyfr"] = "77e051470b364993d017c7cede10120e"
	PPGameRelation["vs5jjwild"] = "00ab11e01d9f85c9b74e0e52233e9a85"
	PPGameRelation["vs5joker"] = "97bd5b1e85dcd00b3cc141a5f3527529"
	PPGameRelation["vswaysfirewmw"] = "32a9b5e9958d37e859f8ae6f0705c095"
	PPGameRelation["vswaysfreezet"] = "7e20b309f0e0bf88a751b41f54e35545"
	PPGameRelation["vswaysmfreya"] = "78b425dd73ed6af262dbafaf6da89fcc"
	PPGameRelation["vswayssevenc"] = "1e93f95916665108cb61065274cb46a4"

	PPGameRelation["vs10bxmasbnza"] = "5b22ab96368a5d3aed8a1f541ead932f"
	PPGameRelation["vswayschilheat"] = "7361d621d9c309724411620f049e715b"
	PPGameRelation["vs117649starz"] = "c829549200d90e7a6b7a48eee41e5a3c"
	PPGameRelation["vs20sh"] = "7aa04b0b3a719bb6e8856dbd06af594e"

	PPGameRelation["vs20fourmc"] = "c6f5ee2ea499da7f10b4ce92e434bb9e"
	PPGameRelation["vs25ultwolgol"] = "7605e4f0d6ed5812280f6ce3a20fd7d3"
	PPGameRelation["vs20sbxmas"] = "be44f1e2c7d4dd7ba95f39abcef2cb62"
	PPGameRelation["vs25xmasparty"] = "e330d77b3eb2cfae299f981c692f18af"

	PPGameRelation["vs5himalaw"] = "cbf9e465f33f7abb892e8a3286b2af79"
	PPGameRelation["vs20schristmas"] = "3140a7726214521c34804325fe35f7a1"
	PPGameRelation["vs15godsofwar"] = "df0492724f0fc4ae4ec183d0c0176d46"

	PPGameRelation["vs40bigjuan"] = "7cdac8e7102f1d8959650c5a76aa9370"
	PPGameRelation["vs25bkofkngdm"] = "9186d388bb0eaec935e5c324d9220422"
	PPGameRelation["vs25btygold"] = "74ad5ad66f45a9d1b8136ef704bbdf26"

	PPGameRelation["vswaysmonkey"] = "ac700413f2a4f9241e306e622a5bf935"

	PPGameRelation["vs10bbdoubled"] = "03ca2d1bef1ca213826750ee1e38b6fd"
	PPGameRelation["vs20clustcol"] = "e342b8ebea1dae1e7dc78a6b6edd7e3b"

	PPGameRelation["vs1dragon8"] = "50fb95404b19b978a8983f2ff6eafc6d"
	PPGameRelation["vs1tigers"] = "be0c7e72bfc9b2a9186793e005e57ee6"
	PPGameRelation["vs20olympdice"] = "1896e4aac745e9e3bc233db2e06512de"
	PPGameRelation["vs9aztecgemsdx"] = "4e88eb8562a4a9fff6d8673b61153a8d"
	PPGameRelation["vs10bbfloats"] = "64d75dfe52f40c443aeae3c5e9d5f4d8"
	PPGameRelation["vs10bburger"] = "ee6844b10cc2295104d72c0554af94b3"
	PPGameRelation["vs10cowgold"] = "18d7dbff8fd377f709d0aa0655e45522"
	PPGameRelation["vs10firestrike"] = "b04d9cd89be5d93602d9a484f8bfe575"
	PPGameRelation["vs10threestar"] = "55180df7fa1b780975c9c0a19b54d8ee"

	PPGameRelation["vs15diamond"] = "a2c7dd829aa83dc2115e969471da9eda"
	PPGameRelation["vs20bnnzdice"] = "2c39d0d7684c6562da4c1fc5b3c07990"
	PPGameRelation["vs20candyblitz"] = "3f90b6a8e9b089615d514a25c1b29323"
	PPGameRelation["vs20candybltz2"] = "5c5ecf28395360512e10ba9b12eb736d"
	PPGameRelation["vs20cbrhst"] = "cb5d7e55024f427d18297db3d3959cb4"
	PPGameRelation["vs20cjcluster"] = "1584763738f901c73e643c55c7f2de24"
	PPGameRelation["vs20cleocatra"] = "6dd63841677a36bf07ce069ec54bc456"
	PPGameRelation["vs20clustext"] = "8895e06c770c160e3a864b31e01fd285"

	PPGameRelation["vs20doghouse"] = "c182b3e577b60ec8d4e411bcbca3f81d"
	PPGameRelation["vs20doghousemh"] = "9198d01543114fa961ad20f0aff63672"
	PPGameRelation["vs20forge"] = "5ac3ad3f0072e0b114a8301f556a8bc5"
	PPGameRelation["vs20gatotgates"] = "4344bf39cbeaafca9edfbd724330d284"
	PPGameRelation["vs20gatotx"] = "33ec04f802880123694d7a223542ed6c"
	PPGameRelation["vs20goldclust"] = "c708cac68bf68fa5de05977fe7f11153"
	PPGameRelation["vs20goldfever"] = "27cdc3e5e12964d38df093e9c65569eb"
	PPGameRelation["vs20multiup"] = "44ab0ddf84f32cc79f6bdcf6309b9710"

	PPGameRelation["vs20olympgate"] = "9afa884f169dfd11a5cd39da32bd4df1"
	PPGameRelation["vs20pbonanza"] = "ab841b96a216b2321baa11d6121185a3"
	PPGameRelation["vs20portals"] = "563f07d59be196ad3cec2533d7c2fbd0"
	PPGameRelation["vs20saiman"] = "09427b3980780ea41a9f15dda6b2726c"
	PPGameRelation["vs20sbxmas"] = "be44f1e2c7d4dd7ba95f39abcef2cb62"
	PPGameRelation["vs20schristmas"] = "3140a7726214521c34804325fe35f7a1"
	PPGameRelation["vs20stickypos"] = "07b39ef4436dea6438ac2b7d5cfad996"
	PPGameRelation["vs20tweethouse"] = "614e8ca458bbc48331b6c9a5900510e1"
	PPGameRelation["vs25checaishen"] = "d7c6f416bfff4113677afc25d050bf33"
	PPGameRelation["vs40stckwldlvl"] = "c4ed86503269c27ef923b68ef0040bb6"
	PPGameRelation["vs50dmdcascade"] = "daec52818a6b49138cc584cef0b3a028"
	PPGameRelation["vs1024fortune"] = "a15d8418eaff591cfa5856a441f6b3eb"
	PPGameRelation["vs1024mahjwins"] = "b1fc7b189ba6f2ccf74e8bcf4b55a686"
	PPGameRelation["vs1024mjwinbns"] = "7c9e9decfcd0b4d23050c9d8531e26bf"
	PPGameRelation["vs7776aztec"] = "e6eaed7085796b0c706cbde5221d05cf"
	PPGameRelation["vswaysmadame"] = "63c8f96398e1f798e40d44e44282b0fc"
	PPGameRelation["vswaysmegahays"] = "17905b61e70edc1cb0e90bd84669dd73"
	PPGameRelation["vswaysoldminer"] = "90746f39bdca21a346389573de57f9d8"
	PPGameRelation["vswaysstampede"] = "39f34ced7888fd518c18262a6bab409d"
	PPGameRelation["vswaysstrlght"] = "bfb8b9e84ff096a54adffb8688a2bf45"
	PPGameRelation["vswayswest"] = "22b9cf05f4d62fec02aca39fb60030c4"
	PPGameRelation["vswayswildgang"] = "b17af566fd988c3d7a06ee12facf0c60"
	PPGameRelation["vswayswildwest"] = "f0f54e0bbb300845990c0586fabd9bed"
	PPGameRelation["vs20bonzgold"] = "45b73bc24d304f030808d138bf1a824e"

	PPGameRelation["vs1fufufu"] = "05fb524b636b6e8a8b89f90b4550ddac"
	PPGameRelation["vs5jokerdice"] = "625498971ffcdb4e4c828c60159a5636"
	PPGameRelation["vs20framazon"] = "20a5902cf138c5ebdfaa313677a0369f"
	PPGameRelation["vs20gravity"] = "43f280df4cdb1b1489d6744d7b1d6332"
	PPGameRelation["vs20hotzone"] = "c47e7fda3a1d706d21fc215b4f080551"
	PPGameRelation["vswaysmoneyman"] = "ef907480a0cebef4e072c9d15ceb1d84"

	PPGameRelation["vswayszombcarn"] = "b66957542324a360e19b6f7b7415ef0d"

	PPGameRelation["vs20hstgldngt"] = "7d09a0335b51fd879867402be7f80a43"
	PPGameRelation["vs20jewelparty"] = "bd44e960f4ef716c0ad0702c4a8beee8"
	PPGameRelation["vs20lvlup"] = "1bf67bfe1e66d5dee133f69542501e5b"
	PPGameRelation["vs20mparty"] = "5ffaacab7e6808fd3b84e8d054c51c7c"
	PPGameRelation["vs20piggybank"] = "c47913d1477cf25524760a3a07dbeedd"

	PPGameRelation["vs25holiday"] = "28e4d704cf080bb6cae517365b27c1fa"
	PPGameRelation["vs25pandatemple"] = "9185ca53157a194f5442894fda79967a"
	PPGameRelation["vs40demonpots"] = "315a7c0ff22e9b9c8f000ba72626b1eb"
	PPGameRelation["vs40rainbowr"] = "5f969b7d1b8027ea52ba96d92d5c6948"

	PPGameRelation["vs20aladdinsorc"] = "da84df20357e880c2fd6058165534cf8"
	PPGameRelation["vs10goldfish"] = "1ee85d001705ebaf63cc0ccae61538ca"
	PPGameRelation["vs10floatdrg"] = "1be0946d4cdc5c26c5922b96c7e154f0"
	PPGameRelation["vs10eyestorm"] = "b9ee66020af593bf5c5145c1b2b24a34"
	PPGameRelation["vs10egypt"] = "ed25af2d1a4744a55e5087112bf23eec"
	PPGameRelation["vs10amm"] = "67ad028a5b9fdc3a1af99351e3bd07ee"
	PPGameRelation["vs9hotroll"] = "8c51d3f8b59cf20bdf4110c6c5baa206"
	PPGameRelation["vs10mayangods"] = "df5b1181bad0795df7f991750e97a09d"

	PPGameRelation["vs20amuleteg"] = "8060b9627603e47cd5381997bf8e54e9"
	PPGameRelation["vs20chickdrop"] = "85c42fe280b53f00702b07b3335356b4"
	PPGameRelation["vs20daydead"] = "39e5d6e45fce6f88c23be485d0eff5f9"
	PPGameRelation["vs20egypttrs"] = "e20027b755b02a3939f3136315e97aba"
	PPGameRelation["vs20eking"] = "97c448a2fafe0a1cdbff7c7831f9dd4a"
	PPGameRelation["vs20ekingrr"] = "77d20f8e6ff87397200f694026914bbf"
	PPGameRelation["vs20emptybank"] = "6f0279bab4e1e4812280743ea720b549"
	PPGameRelation["vs20fparty2"] = "6be857d19ef7b4fac0ba9c375538a7d9"
}

func regPlat(ns string, plat IPlat) {
	type LaunchGamePs struct {
		UID      string
		Game     string
		Lang     string
		UseProxy bool
	}
	type LaunchGameRet struct {
		Url string
	}

	mux.RegHttpWithSample(path.Join("/", ns, "LaunchGame"), "拉起游戏", ns, func(ps LaunchGamePs, ret *LaunchGameRet) (err error) {
		// ba plat.GetBalance(ps.UID)
		ret.Url, err = plat.LaunchGame(ps.UID, PPGameRelation[ps.Game], ps.Lang, ps.UseProxy)
		return
	}, LaunchGamePs{"123456", "39", "en", true})

	type FundTransferInPs struct {
		UID    string
		Amount float64
	}
	type FundTransferInRet struct {
		Status string
	}
	mux.RegHttpWithSample(path.Join("/", ns, "FundTransferIn"), "带入", ns, func(ps FundTransferInPs, ret *FundTransferInRet) (err error) {
		ret.Status = plat.FundTransferIn(ps.UID, ps.Amount)
		return
	}, FundTransferInPs{"123456", 1000.0})

	type FundTransferOutPs struct {
		UID string
	}
	type FundTransferOutRet struct {
		Amount float64
		Status string
	}
	mux.RegHttpWithSample(path.Join("/", ns, "FundTransferOut"), "带出", ns, func(ps FundTransferOutPs, ret *FundTransferOutRet) (err error) {
		ret.Amount, ret.Status = plat.FundTransferOut(ps.UID)
		return
	}, FundTransferOutPs{"123456"})

	type GetBalancePs struct {
		UID string
	}
	type GetBalanceRet struct {
		Balance float64
	}
	mux.RegHttpWithSample(path.Join("/", ns, "GetBalance"), "获取玩家余额", ns, func(ps GetBalancePs, ret *GetBalanceRet) (err error) {
		ret.Balance, err = plat.GetBalance(ps.UID)
		return
	}, GetBalancePs{"123456"})

	type GetGameListPs struct {
		// GameType string
	}
	// type GetGameListRet struct {
	// 	Games HotGames
	// }
	type GetGameListRet = json.RawMessage

	var gamelist json.RawMessage
	mux.RegHttpWithSample(path.Join("/", ns, "GetGameList"), "获取游戏列表", ns, func(ps GetGameListPs, ret *GetGameListRet) (err error) {
		if len(gamelist) != 0 {
			*ret = gamelist
			return
		}
		games, err := plat.GetGameList()
		if err != nil {
			return
		}

		buf, _ := json.Marshal(define.M{
			"List": games,
		})
		gamelist = json.RawMessage(buf)
		*ret = gamelist
		return
	}, GetGameListPs{})

	type RegistPs struct {
		UID string
	}
	type RegistRet struct {
	}
	mux.RegHttpWithSample(path.Join("/", ns, "Regist"), "注册", ns, func(ps RegistPs, ret *RegistRet) (err error) {
		err = plat.Regist(ps.UID)
		return
	}, RegistPs{"123456"})

}

type HotGames []*HotGame

type HotGame struct {
	Plat string
	ID   string
	Name string
	Type int
	Icon string
}
