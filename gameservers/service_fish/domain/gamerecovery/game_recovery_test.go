package gamerecovery

import (
	"fmt"
	"os"
	game_recovery_proto "serve/service_fish/domain/gamerecovery/proto"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestService_Data(t *testing.T) {
	if len(os.Args) < 4 {
		fmt.Println("請輸入正確的參數(Host_id Member_id Game_id)")
		return
	}

	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&timeout=%ss",
			"gm-server", "MBpBy7RCNyH2HjUW", "35.234.13.198", "3306", "slot-game", "2"),
	)

	if err != nil {
		fmt.Println("Connect DB Failed.", err)
		return
	}

	dbResult := Builder().setHostExtId("test").build()
	data := &game_recovery_proto.GameRecovery{}
	if ok := db.
		Table("game_recovery").
		Select("host_id, member_id, game_id, subgame_id, game_data, accounting_sn, game_end").
		Where("member_id = ? AND host_id = ? AND game_id = ?", os.Args[2], os.Args[1], os.Args[3]).
		Scan(dbResult.recovery).
		RowsAffected; ok != 1 {
		return
	}

	if err := proto.Unmarshal(dbResult.recovery.GameData, data); err != nil {
		return
	}

	fmt.Println("分母:", data.Rtp.Denominator, "分子", data.Rtp.Budget)
	fmt.Println("水位:", data.Rtp, "傭兵:", data.Mercenary)

	db.Close()
}
