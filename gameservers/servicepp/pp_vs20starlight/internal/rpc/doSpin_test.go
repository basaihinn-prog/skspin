package rpc

import (
	"os"
	"serve/servicepp/ppcomm"
	"testing"
)

func TestDoSpin(t *testing.T) {
	file, _ := os.Create("doSpin.log")
	ppcomm.PrettyPrintVas(file, "action=doSpin&symbol=vs20starlight&c=0.5&l=20&bl=0&index=2&counter=3&repeat=0&mgckey=AUTHTOKEN@6318b8f52aa27a735236b247dafd0e27037578637940d87ad54250bfeed0e431~stylename@hllgd_hollygod~SESSION@7dbfc452-ef7a-4b14-8069-a857e5dc165b~SN@e9d7b550")
	ppcomm.PrettyPrintVas(file, "tw=0.00&rid=38693909614091&balance=999,990.00&accm=cp&acci=0&index=2&balance_cash=999,990.00&reel_set=1&balance_bonus=0.00&na=s&accv=0&bl=0&stime=1725258160393&sa=6,6,3,5,9,4&sb=3,7,8,11,10,5&sh=5&c=0.50&sver=5&counter=4&ntp=-10.00&l=20&s=10,7,5,8,3,3,9,8,10,8,3,3,10,6,5,8,11,9,8,11,7,10,11,9,8,9,7,10,8,5&w=0.00")

	ppcomm.PrettyPrintVas(file, "action=doSpin&symbol=vs20starlight&c=0.05&l=20&bl=0&index=12&counter=23&repeat=0&mgckey=AUTHTOKEN@00e44af101baa8906dcc093944050bfd5c15a8d0f238ab3155fd25146163e676~stylename@hllgd_hollygod~SESSION@640ab73d-d307-4c3f-8823-ddcd75d4dc8e~SN@f6b8c376")
	ppcomm.PrettyPrintVas(file, "tw=0.40&tmb=1,10~2,10~3,10~4,10~6,10~10,10~26,10~29,10&rid=38694566884091&balance=999,954.80&accm=cp&acci=0&index=12&balance_cash=999,954.80&reel_set=1&balance_bonus=0.00&na=s&accv=0&rs=mc&tmb_win=0.40&l0=0~0.40~1~2~3~4~6~10~26~29&rs_p=0&bl=0&stime=1725259451853&sa=5,9,7,7,9,9&sb=6,6,3,3,6,5&rs_c=1&sh=5&rs_m=1&c=0.05&sver=5&counter=24&ntp=-35.20&l=20&s=7,10,10,10,10,9,10,3,6,3,10,8,9,6,7,11,11,8,1,7,1,3,7,7,6,8,10,11,4,10&w=0.40")

	ppcomm.PrettyPrintVas(file, "action=doSpin&symbol=vs20starlight&c=0.05&l=20&bl=0&index=13&counter=25&repeat=0&mgckey=AUTHTOKEN@00e44af101baa8906dcc093944050bfd5c15a8d0f238ab3155fd25146163e676~stylename@hllgd_hollygod~SESSION@640ab73d-d307-4c3f-8823-ddcd75d4dc8e~SN@f6b8c376")
	ppcomm.PrettyPrintVas(file, "tw=0.40&rid=38694566884091&tmb_res=0.40&balance=999,954.80&accm=cp&acci=0&index=13&balance_cash=999,954.80&reel_set=1&balance_bonus=0.00&na=c&accv=0&rs_t=1&tmb_win=0.40&bl=0&stime=1725259455643&sa=5,7,10,8,7,11&sb=6,6,3,3,6,5&sh=5&c=0.05&sver=5&counter=26&ntp=-35.20&l=20&s=5,9,9,7,1,9,7,3,7,3,9,9,9,6,6,11,11,8,1,7,7,3,7,8,6,8,1,11,4,7&w=0.00")

	ppcomm.PrettyPrintVas(file, "symbol=vs20starlight&action=doCollect&index=14&counter=27&repeat=0&mgckey=AUTHTOKEN@00e44af101baa8906dcc093944050bfd5c15a8d0f238ab3155fd25146163e676~stylename@hllgd_hollygod~SESSION@640ab73d-d307-4c3f-8823-ddcd75d4dc8e~SN@f6b8c376")
	ppcomm.PrettyPrintVas(file, "balance=999,955.20&index=14&balance_cash=999,955.20&balance_bonus=0.00&na=s&stime=1725259455784&sver=5&counter=28&ntp=-34.80")

}

func TestFetch(t *testing.T) {

}
