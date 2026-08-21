package rtp

import (
	"fmt"
	"serve/fish_comm/rng"
	"serve/service_fish/domain/probability"
	PSFM_00004_95_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-95-1"
	PSFM_00004_96_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-96-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00004_98_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-98-1"
	"serve/service_fish/models"
	"strconv"
	"testing"
)

var mathModuleId_00003 string

const (
	game_id_00003         = models.PSF_ON_00003
	subgame_id_00003      = 0
	secWebSocketKey_00003 = "jerry"
	bet_00003             = 1
	rate_00003            = 100

	run_times_00003 = 100 * 100 * 100 * 1
	//mathModule_id_00003     = models.PSFM_00004_98_1

	// Weapon Bullet Bet
	betNormal  = 1
	betShotGun = 3
	betBazooka = 5
	betMortar  = 10
	// Rtp
	threw = 0
	bite  = 1

	switchWeaponBullet = 10
)

//func TestService_00003_MultiCalc(t *testing.T) {
//   var threw_fish_0 uint64 = 0
//   var threw_fish_1 uint64 = 0
//   var threw_fish_2 uint64 = 0
//   var threw_fish_3 uint64 = 0
//   var threw_fish_4 uint64 = 0
//   var threw_fish_5 uint64 = 0
//   var threw_fish_6 uint64 = 0
//   var threw_fish_7 uint64 = 0
//   var threw_fish_8 uint64 = 0
//   var threw_fish_9 uint64 = 0
//   var threw_fish_10 uint64 = 0
//   var threw_fish_11 uint64 = 0
//   var threw_fish_12 uint64 = 0
//   var threw_fish_13 uint64 = 0
//   var threw_fish_14 uint64 = 0
//   var threw_fish_15 uint64 = 0
//   var threw_fish_16 uint64 = 0
//   var threw_fish_17 uint64 = 0
//   var threw_fish_303 uint64 = 0
//
//   var bite_fish_0 uint64 = 0
//   var bite_fish_1 uint64 = 0
//   var bite_fish_2 uint64 = 0
//   var bite_fish_3 uint64 = 0
//   var bite_fish_4 uint64 = 0
//   var bite_fish_5 uint64 = 0
//   var bite_fish_6 uint64 = 0
//   var bite_fish_7 uint64 = 0
//   var bite_fish_8 uint64 = 0
//   var bite_fish_9 uint64 = 0
//   var bite_fish_10 uint64 = 0
//   var bite_fish_11 uint64 = 0
//   var bite_fish_12 uint64 = 0
//   var bite_fish_13 uint64 = 0
//   var bite_fish_14 uint64 = 0
//   var bite_fish_15 uint64 = 0
//   var bite_fish_16 uint64 = 0
//   var bite_fish_17 uint64 = 0
//   var bite_fish_303 uint64 = 0
//
//   var threw_hit_0 uint64 = 0
//   var threw_hit_1 uint64 = 0
//   var threw_hit_2 uint64 = 0
//   var threw_hit_3 uint64 = 0
//   var threw_hit_4 uint64 = 0
//   var threw_hit_5 uint64 = 0
//   var threw_hit_6 uint64 = 0
//   var threw_hit_7 uint64 = 0
//   var threw_hit_8 uint64 = 0
//   var threw_hit_9 uint64 = 0
//   var threw_hit_10 uint64 = 0
//   var threw_hit_11 uint64 = 0
//   var threw_hit_12 uint64 = 0
//   var threw_hit_13 uint64 = 0
//   var threw_hit_14 uint64 = 0
//   var threw_hit_15 uint64 = 0
//   var threw_hit_16 uint64 = 0
//   var threw_hit_17 uint64 = 0
//   var threw_hit_303 uint64 = 0
//
//   var bite_hit_0 uint64 = 0
//   var bite_hit_1 uint64 = 0
//   var bite_hit_2 uint64 = 0
//   var bite_hit_3 uint64 = 0
//   var bite_hit_4 uint64 = 0
//   var bite_hit_5 uint64 = 0
//   var bite_hit_6 uint64 = 0
//   var bite_hit_7 uint64 = 0
//   var bite_hit_8 uint64 = 0
//   var bite_hit_9 uint64 = 0
//   var bite_hit_10 uint64 = 0
//   var bite_hit_11 uint64 = 0
//   var bite_hit_12 uint64 = 0
//   var bite_hit_13 uint64 = 0
//   var bite_hit_14 uint64 = 0
//   var bite_hit_15 uint64 = 0
//   var bite_hit_16 uint64 = 0
//   var bite_hit_17 uint64 = 0
//   var bite_hit_303 uint64 = 0
//
//   var threw_win_0 uint64 = 0
//   var threw_win_1 uint64 = 0
//   var threw_win_2 uint64 = 0
//   var threw_win_3 uint64 = 0
//   var threw_win_4 uint64 = 0
//   var threw_win_5 uint64 = 0
//   var threw_win_6 uint64 = 0
//   var threw_win_7 uint64 = 0
//   var threw_win_8 uint64 = 0
//   var threw_win_9 uint64 = 0
//   var threw_win_10 uint64 = 0
//   var threw_win_11 uint64 = 0
//   var threw_win_12 uint64 = 0
//   var threw_win_13 uint64 = 0
//   var threw_win_14 uint64 = 0
//   var threw_win_15 uint64 = 0
//   var threw_win_16 uint64 = 0
//   var threw_win_17 uint64 = 0
//   var threw_win_303 uint64 = 0
//
//   var bite_win_0 uint64 = 0
//   var bite_win_1 uint64 = 0
//   var bite_win_2 uint64 = 0
//   var bite_win_3 uint64 = 0
//   var bite_win_4 uint64 = 0
//   var bite_win_5 uint64 = 0
//   var bite_win_6 uint64 = 0
//   var bite_win_7 uint64 = 0
//   var bite_win_8 uint64 = 0
//   var bite_win_9 uint64 = 0
//   var bite_win_10 uint64 = 0
//   var bite_win_11 uint64 = 0
//   var bite_win_12 uint64 = 0
//   var bite_win_13 uint64 = 0
//   var bite_win_14 uint64 = 0
//   var bite_win_15 uint64 = 0
//   var bite_win_16 uint64 = 0
//   var bite_win_17 uint64 = 0
//   var bite_win_303 uint64 = 0
//
//   var threw_bullet_0 uint64 = 0
//   var threw_bullet_1 uint64 = 0
//   var threw_bullet_2 uint64 = 0
//   var threw_bullet_3 uint64 = 0
//   var threw_bullet_4 uint64 = 0
//   var threw_bullet_5 uint64 = 0
//   var threw_bullet_6 uint64 = 0
//   var threw_bullet_7 uint64 = 0
//   var threw_bullet_8 uint64 = 0
//   var threw_bullet_9 uint64 = 0
//   var threw_bullet_10 uint64 = 0
//   var threw_bullet_11 uint64 = 0
//   var threw_bullet_12 uint64 = 0
//   var threw_bullet_13 uint64 = 0
//   var threw_bullet_14 uint64 = 0
//   var threw_bullet_15 uint64 = 0
//   var threw_bullet_16 uint64 = 0
//   var threw_bullet_17 uint64 = 0
//   var threw_bullet_303 uint64 = 0
//
//   var bite_bullet_0 uint64 = 0
//   var bite_bullet_1 uint64 = 0
//   var bite_bullet_2 uint64 = 0
//   var bite_bullet_3 uint64 = 0
//   var bite_bullet_4 uint64 = 0
//   var bite_bullet_5 uint64 = 0
//   var bite_bullet_6 uint64 = 0
//   var bite_bullet_7 uint64 = 0
//   var bite_bullet_8 uint64 = 0
//   var bite_bullet_9 uint64 = 0
//   var bite_bullet_10 uint64 = 0
//   var bite_bullet_11 uint64 = 0
//   var bite_bullet_12 uint64 = 0
//   var bite_bullet_13 uint64 = 0
//   var bite_bullet_14 uint64 = 0
//   var bite_bullet_15 uint64 = 0
//   var bite_bullet_16 uint64 = 0
//   var bite_bullet_17 uint64 = 0
//   var bite_bullet_303 uint64 = 0
//
//   for index := 0; index < 19; index++ {
//       fishId := int32(index)
//       switch fishId {
//       case 18:
//           fishId = 303
//       }
//
//       indexSecWebSocketKey := "jerry" + strconv.Itoa(int(fishId))
//
//       for i := 0; i < run_times_00003; i++ {
//           Service.Decrease(
//               game_id_00003, subgame_id_00003, mathModule_id_00003, indexSecWebSocketKey,
//               bet_00003 * rate_00003,
//           )
//           rtpId := Service.RtpId(indexSecWebSocketKey, game_id_00003, subgame_id_00003)
//           rtpState := Service.RtpState(indexSecWebSocketKey, game_id_00003, subgame_id_00003)
//
//           result := probability.Service.Calc(
//               game_id_00003,
//               mathModule_id_00003,
//               rtpId,
//               fishId,
//               -1,
//               0,
//               bet_00003 * rate_00003,
//           )
//
//           switch fishId {
//           case 0:
//               threw_fish_0, threw_hit_0, threw_win_0, bite_fish_0, bite_hit_0, bite_win_0, threw_bullet_0, bite_bullet_0 = getResult(
//                   rtpState, t, result,
//                   threw_fish_0, threw_hit_0, threw_win_0, bite_fish_0, bite_hit_0, bite_win_0,
//                   threw_bullet_0, bite_bullet_0,
//               )
//               threw_bullet_0, bite_bullet_0 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_0, bite_bullet_0)
//
//           case 1:
//               threw_fish_1, threw_hit_1, threw_win_1, bite_fish_1, bite_hit_1, bite_win_1, threw_bullet_1, bite_bullet_1 = getResult(
//                   rtpState, t, result,
//                   threw_fish_1, threw_hit_1, threw_win_1, bite_fish_1, bite_hit_1, bite_win_1,
//                   threw_bullet_1, bite_bullet_1,
//               )
//               threw_bullet_1, bite_bullet_1 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_1, bite_bullet_1)
//
//           case 2:
//               threw_fish_2, threw_hit_2, threw_win_2, bite_fish_2, bite_hit_2, bite_win_2, threw_bullet_2, bite_bullet_2 = getResult(
//                   rtpState, t, result,
//                   threw_fish_2, threw_hit_2, threw_win_2, bite_fish_2, bite_hit_2, bite_win_2,
//                   threw_bullet_2, bite_bullet_2,
//               )
//               threw_bullet_2, bite_bullet_2 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_2, bite_bullet_2)
//
//           case 3:
//               threw_fish_3, threw_hit_3, threw_win_3, bite_fish_3, bite_hit_3, bite_win_3, threw_bullet_3, bite_bullet_3 = getResult(
//                   rtpState, t, result,
//                   threw_fish_3, threw_hit_3, threw_win_3, bite_fish_3, bite_hit_3, bite_win_3,
//                   threw_bullet_3, bite_bullet_3,
//               )
//               threw_bullet_3, bite_bullet_3 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_3, bite_bullet_3)
//
//           case 4:
//               threw_fish_4, threw_hit_4, threw_win_4, bite_fish_4, bite_hit_4, bite_win_4, threw_bullet_4, bite_bullet_4 = getResult(
//                   rtpState, t, result,
//                   threw_fish_4, threw_hit_4, threw_win_4, bite_fish_4, bite_hit_4, bite_win_4,
//                   threw_bullet_4, bite_bullet_4,
//               )
//               threw_bullet_4, bite_bullet_4 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_4, bite_bullet_4)
//
//           case 5:
//               threw_fish_5, threw_hit_5, threw_win_5, bite_fish_5, bite_hit_5, bite_win_5, threw_bullet_5, bite_bullet_5 = getResult(
//                   rtpState, t, result,
//                   threw_fish_5, threw_hit_5, threw_win_5, bite_fish_5, bite_hit_5, bite_win_5,
//                   threw_bullet_5, bite_bullet_5,
//               )
//               threw_bullet_5, bite_bullet_5 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_5, bite_bullet_5)
//
//           case 6:
//               threw_fish_6, threw_hit_6, threw_win_6, bite_fish_6, bite_hit_6, bite_win_6, threw_bullet_6, bite_bullet_6 = getResult(
//                   rtpState, t, result,
//                   threw_fish_6, threw_hit_6, threw_win_6, bite_fish_6, bite_hit_6, bite_win_6,
//                   threw_bullet_6, bite_bullet_6,
//               )
//               threw_bullet_6, bite_bullet_6 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_6, bite_bullet_6)
//
//           case 7:
//               threw_fish_7, threw_hit_7, threw_win_7, bite_fish_7, bite_hit_7, bite_win_7, threw_bullet_7, bite_bullet_7 = getResult(
//                   rtpState, t, result,
//                   threw_fish_7, threw_hit_7, threw_win_7, bite_fish_7, bite_hit_7, bite_win_7,
//                   threw_bullet_7, bite_bullet_7,
//               )
//               threw_bullet_7, bite_bullet_7 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_7, bite_bullet_7)
//
//           case 8:
//               threw_fish_8, threw_hit_8, threw_win_8, bite_fish_8, bite_hit_8, bite_win_8, threw_bullet_8, bite_bullet_8 = getResult(
//                   rtpState, t, result,
//                   threw_fish_8, threw_hit_8, threw_win_8, bite_fish_8, bite_hit_8, bite_win_8,
//                   threw_bullet_8, bite_bullet_8,
//               )
//               threw_bullet_8, bite_bullet_8 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_8, bite_bullet_8)
//
//           case 9:
//               threw_fish_9, threw_hit_9, threw_win_9, bite_fish_9, bite_hit_9, bite_win_9, threw_bullet_9, bite_bullet_9 = getResult(
//                   rtpState, t, result,
//                   threw_fish_9, threw_hit_9, threw_win_9, bite_fish_9, bite_hit_9, bite_win_9,
//                   threw_bullet_9, bite_bullet_9,
//               )
//               threw_bullet_9, bite_bullet_9 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_9, bite_bullet_9)
//
//           case 10:
//               threw_fish_10, threw_hit_10, threw_win_10, bite_fish_10, bite_hit_10, bite_win_10, threw_bullet_10, bite_bullet_10 = getResult(
//                   rtpState, t, result,
//                   threw_fish_10, threw_hit_10, threw_win_10, bite_fish_10, bite_hit_10, bite_win_10,
//                   threw_bullet_10, bite_bullet_10,
//               )
//               threw_bullet_10, bite_bullet_10 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_10, bite_bullet_10)
//
//           case 11:
//               threw_fish_11, threw_hit_11, threw_win_11, bite_fish_11, bite_hit_11, bite_win_11, threw_bullet_11, bite_bullet_11 = getResult(
//                   rtpState, t, result,
//                   threw_fish_11, threw_hit_11, threw_win_11, bite_fish_11, bite_hit_11, bite_win_11,
//                   threw_bullet_11, bite_bullet_11,
//               )
//               threw_bullet_11, bite_bullet_11 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_11, bite_bullet_11)
//
//           case 12:
//               threw_fish_12, threw_hit_12, threw_win_12, bite_fish_12, bite_hit_12, bite_win_12, threw_bullet_12, bite_bullet_12 = getResult(
//                   rtpState, t, result,
//                   threw_fish_12, threw_hit_12, threw_win_12, bite_fish_12, bite_hit_12, bite_win_12,
//                   threw_bullet_12, bite_bullet_12,
//               )
//               threw_bullet_12, bite_bullet_12 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_12, bite_bullet_12)
//
//           case 13:
//               threw_fish_13, threw_hit_13, threw_win_13, bite_fish_13, bite_hit_13, bite_win_13, threw_bullet_13, bite_bullet_13 = getResult(
//                   rtpState, t, result,
//                   threw_fish_13, threw_hit_13, threw_win_13, bite_fish_13, bite_hit_13, bite_win_13,
//                   threw_bullet_13, bite_bullet_13,
//               )
//               threw_bullet_13, bite_bullet_13 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_13, bite_bullet_13)
//
//           case 14:
//               threw_fish_14, threw_hit_14, threw_win_14, bite_fish_14, bite_hit_14, bite_win_14, threw_bullet_14, bite_bullet_14 = getResult(
//                   rtpState, t, result,
//                   threw_fish_14, threw_hit_14, threw_win_14, bite_fish_14, bite_hit_14, bite_win_14,
//                   threw_bullet_14, bite_bullet_14,
//               )
//               threw_bullet_14, bite_bullet_14 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_14, bite_bullet_14)
//
//           case 15:
//               threw_fish_15, threw_hit_15, threw_win_15, bite_fish_15, bite_hit_15, bite_win_15, threw_bullet_15, bite_bullet_15 = getResult(
//                   rtpState, t, result,
//                   threw_fish_15, threw_hit_15, threw_win_15, bite_fish_15, bite_hit_15, bite_win_15,
//                   threw_bullet_15, bite_bullet_15,
//               )
//               threw_bullet_15, bite_bullet_15 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_15, bite_bullet_15)
//
//           case 16:
//               threw_fish_16, threw_hit_16, threw_win_16, bite_fish_16, bite_hit_16, bite_win_16, threw_bullet_16, bite_bullet_16 = getResult(
//                   rtpState, t, result,
//                   threw_fish_16, threw_hit_16, threw_win_16, bite_fish_16, bite_hit_16, bite_win_16,
//                   threw_bullet_16, bite_bullet_16,
//               )
//               threw_bullet_16, bite_bullet_16 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_16, bite_bullet_16)
//
//           case 17:
//               threw_fish_17, threw_hit_17, threw_win_17, bite_fish_17, bite_hit_17, bite_win_17, threw_bullet_17, bite_bullet_17 = getResult(
//                   rtpState, t, result,
//                   threw_fish_17, threw_hit_17, threw_win_17, bite_fish_17, bite_hit_17, bite_win_17,
//                   threw_bullet_17, bite_bullet_17,
//               )
//               threw_bullet_17, bite_bullet_17 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_17, bite_bullet_17)
//
//           case 303:
//               switch rtpState {
//               case threw:
//                   if result.Pay > 0 || result.Bullet > 0 {
//                       threw_hit_303++
//                       threw_win_303 += uint64(result.Pay * result.Multiplier)
//                   }
//                   threw_fish_303++
//
//               case bite:
//                   if result.Pay > 0 || result.Bullet > 0 {
//                       bite_hit_303++
//                       bite_win_303 += uint64(result.Pay * result.Multiplier)
//                   }
//                   bite_fish_303++
//
//               default:
//                   t.Fatal("RTP State Error")
//               }
//
//               threw_bullet_303, bite_bullet_303 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_303, bite_bullet_303)
//           }
//       }
//   }
//
//   threw_win_0, bite_win_0, threw_bullet_0, bite_bullet_0 = mercenaryBulletProcess(threw_win_0, bite_win_0, threw_bullet_0, bite_bullet_0)
//   threw_win_1, bite_win_1, threw_bullet_1, bite_bullet_1 = mercenaryBulletProcess(threw_win_1, bite_win_1, threw_bullet_1, bite_bullet_1)
//   threw_win_2, bite_win_2, threw_bullet_2, bite_bullet_2 = mercenaryBulletProcess(threw_win_2, bite_win_2, threw_bullet_2, bite_bullet_2)
//   threw_win_3, bite_win_3, threw_bullet_3, bite_bullet_3 = mercenaryBulletProcess(threw_win_3, bite_win_3, threw_bullet_3, bite_bullet_3)
//   threw_win_4, bite_win_4, threw_bullet_4, bite_bullet_4 = mercenaryBulletProcess(threw_win_4, bite_win_4, threw_bullet_4, bite_bullet_4)
//   threw_win_5, bite_win_5, threw_bullet_5, bite_bullet_5 = mercenaryBulletProcess(threw_win_5, bite_win_5, threw_bullet_5, bite_bullet_5)
//   threw_win_6, bite_win_6, threw_bullet_6, bite_bullet_6 = mercenaryBulletProcess(threw_win_6, bite_win_6, threw_bullet_6, bite_bullet_6)
//   threw_win_7, bite_win_7, threw_bullet_7, bite_bullet_7 = mercenaryBulletProcess(threw_win_7, bite_win_7, threw_bullet_7, bite_bullet_7)
//   threw_win_8, bite_win_8, threw_bullet_8, bite_bullet_8 = mercenaryBulletProcess(threw_win_8, bite_win_8, threw_bullet_8, bite_bullet_8)
//   threw_win_9, bite_win_9, threw_bullet_9, bite_bullet_9 = mercenaryBulletProcess(threw_win_9, bite_win_9, threw_bullet_9, bite_bullet_9)
//   threw_win_10, bite_win_10, threw_bullet_10, bite_bullet_10 = mercenaryBulletProcess(threw_win_10, bite_win_10, threw_bullet_10, bite_bullet_10)
//   threw_win_11, bite_win_11, threw_bullet_11, bite_bullet_11 = mercenaryBulletProcess(threw_win_11, bite_win_11, threw_bullet_11, bite_bullet_11)
//   threw_win_12, bite_win_12, threw_bullet_12, bite_bullet_12 = mercenaryBulletProcess(threw_win_12, bite_win_12, threw_bullet_12, bite_bullet_12)
//   threw_win_13, bite_win_13, threw_bullet_13, bite_bullet_13 = mercenaryBulletProcess(threw_win_13, bite_win_13, threw_bullet_13, bite_bullet_13)
//   threw_win_14, bite_win_14, threw_bullet_14, bite_bullet_14 = mercenaryBulletProcess(threw_win_14, bite_win_14, threw_bullet_14, bite_bullet_14)
//   threw_win_15, bite_win_15, threw_bullet_15, bite_bullet_15 = mercenaryBulletProcess(threw_win_15, bite_win_15, threw_bullet_15, bite_bullet_15)
//   threw_win_16, bite_win_16, threw_bullet_16, bite_bullet_16 = mercenaryBulletProcess(threw_win_16, bite_win_16, threw_bullet_16, bite_bullet_16)
//   threw_win_17, bite_win_17, threw_bullet_17, bite_bullet_17 = mercenaryBulletProcess(threw_win_17, bite_win_17, threw_bullet_17, bite_bullet_17)
//   threw_win_303, bite_win_303, threw_bullet_303, bite_bullet_303 = mercenaryBulletProcess(threw_win_303, bite_win_303, threw_bullet_303, bite_bullet_303)
//
//   t.Log("000", threw_fish_0, threw_hit_0, threw_win_0, threw_bullet_0, bite_fish_0, bite_hit_0, bite_win_0, bite_bullet_0)
//   t.Log("001", threw_fish_1, threw_hit_1, threw_win_1, threw_bullet_1, bite_fish_1, bite_hit_1, bite_win_1, bite_bullet_1)
//   t.Log("002", threw_fish_2, threw_hit_2, threw_win_2, threw_bullet_2, bite_fish_2, bite_hit_2, bite_win_2, bite_bullet_2)
//   t.Log("003", threw_fish_3, threw_hit_3, threw_win_3, threw_bullet_3, bite_fish_3, bite_hit_3, bite_win_3, bite_bullet_3)
//   t.Log("004", threw_fish_4, threw_hit_4, threw_win_4, threw_bullet_4, bite_fish_4, bite_hit_4, bite_win_4, bite_bullet_4)
//   t.Log("005", threw_fish_5, threw_hit_5, threw_win_5, threw_bullet_5, bite_fish_5, bite_hit_5, bite_win_5, bite_bullet_5)
//   t.Log("006", threw_fish_6, threw_hit_6, threw_win_6, threw_bullet_6, bite_fish_6, bite_hit_6, bite_win_6, bite_bullet_6)
//   t.Log("007", threw_fish_7, threw_hit_7, threw_win_7, threw_bullet_7, bite_fish_7, bite_hit_7, bite_win_7, bite_bullet_7)
//   t.Log("008", threw_fish_8, threw_hit_8, threw_win_8, threw_bullet_8, bite_fish_8, bite_hit_8, bite_win_8, bite_bullet_8)
//   t.Log("009", threw_fish_9, threw_hit_9, threw_win_9, threw_bullet_9, bite_fish_9, bite_hit_9, bite_win_9, bite_bullet_9)
//   t.Log("010", threw_fish_10, threw_hit_10, threw_win_10, threw_bullet_10, bite_fish_10, bite_hit_10, bite_win_10, bite_bullet_10)
//   t.Log("011", threw_fish_11, threw_hit_11, threw_win_11, threw_bullet_11, bite_fish_11, bite_hit_11, bite_win_11, bite_bullet_11)
//   t.Log("012", threw_fish_12, threw_hit_12, threw_win_12, threw_bullet_12, bite_fish_12, bite_hit_12, bite_win_12, bite_bullet_12)
//   t.Log("013", threw_fish_13, threw_hit_13, threw_win_13, threw_bullet_13, bite_fish_13, bite_hit_13, bite_win_13, bite_bullet_13)
//   t.Log("014", threw_fish_14, threw_hit_14, threw_win_14, threw_bullet_14, bite_fish_14, bite_hit_14, bite_win_14, bite_bullet_14)
//   t.Log("015", threw_fish_15, threw_hit_15, threw_win_15, threw_bullet_15, bite_fish_15, bite_hit_15, bite_win_15, bite_bullet_15)
//   t.Log("016", threw_fish_16, threw_hit_16, threw_win_16, threw_bullet_16, bite_fish_16, bite_hit_16, bite_win_16, bite_bullet_16)
//   t.Log("017", threw_fish_17, threw_hit_17, threw_win_17, threw_bullet_17, bite_fish_17, bite_hit_17, bite_win_17, bite_bullet_17)
//   t.Log("303", threw_fish_303, threw_hit_303, threw_win_303, threw_bullet_303, bite_fish_303, bite_hit_303, bite_win_303, bite_bullet_303)
//}

func TestService_00003_Calc(t *testing.T) {
	var threw_fish_0 uint64 = 0
	var threw_fish_1 uint64 = 0
	var threw_fish_2 uint64 = 0
	var threw_fish_3 uint64 = 0
	var threw_fish_4 uint64 = 0
	var threw_fish_5 uint64 = 0
	var threw_fish_6 uint64 = 0
	var threw_fish_7 uint64 = 0
	var threw_fish_8 uint64 = 0
	var threw_fish_9 uint64 = 0
	var threw_fish_10 uint64 = 0
	var threw_fish_11 uint64 = 0
	var threw_fish_12 uint64 = 0
	var threw_fish_13 uint64 = 0
	var threw_fish_14 uint64 = 0
	var threw_fish_15 uint64 = 0
	var threw_fish_16 uint64 = 0
	var threw_fish_17 uint64 = 0
	var threw_fish_303 uint64 = 0

	var bite_fish_0 uint64 = 0
	var bite_fish_1 uint64 = 0
	var bite_fish_2 uint64 = 0
	var bite_fish_3 uint64 = 0
	var bite_fish_4 uint64 = 0
	var bite_fish_5 uint64 = 0
	var bite_fish_6 uint64 = 0
	var bite_fish_7 uint64 = 0
	var bite_fish_8 uint64 = 0
	var bite_fish_9 uint64 = 0
	var bite_fish_10 uint64 = 0
	var bite_fish_11 uint64 = 0
	var bite_fish_12 uint64 = 0
	var bite_fish_13 uint64 = 0
	var bite_fish_14 uint64 = 0
	var bite_fish_15 uint64 = 0
	var bite_fish_16 uint64 = 0
	var bite_fish_17 uint64 = 0
	var bite_fish_303 uint64 = 0

	var threw_hit_0 uint64 = 0
	var threw_hit_1 uint64 = 0
	var threw_hit_2 uint64 = 0
	var threw_hit_3 uint64 = 0
	var threw_hit_4 uint64 = 0
	var threw_hit_5 uint64 = 0
	var threw_hit_6 uint64 = 0
	var threw_hit_7 uint64 = 0
	var threw_hit_8 uint64 = 0
	var threw_hit_9 uint64 = 0
	var threw_hit_10 uint64 = 0
	var threw_hit_11 uint64 = 0
	var threw_hit_12 uint64 = 0
	var threw_hit_13 uint64 = 0
	var threw_hit_14 uint64 = 0
	var threw_hit_15 uint64 = 0
	var threw_hit_16 uint64 = 0
	var threw_hit_17 uint64 = 0
	var threw_hit_303 uint64 = 0

	var bite_hit_0 uint64 = 0
	var bite_hit_1 uint64 = 0
	var bite_hit_2 uint64 = 0
	var bite_hit_3 uint64 = 0
	var bite_hit_4 uint64 = 0
	var bite_hit_5 uint64 = 0
	var bite_hit_6 uint64 = 0
	var bite_hit_7 uint64 = 0
	var bite_hit_8 uint64 = 0
	var bite_hit_9 uint64 = 0
	var bite_hit_10 uint64 = 0
	var bite_hit_11 uint64 = 0
	var bite_hit_12 uint64 = 0
	var bite_hit_13 uint64 = 0
	var bite_hit_14 uint64 = 0
	var bite_hit_15 uint64 = 0
	var bite_hit_16 uint64 = 0
	var bite_hit_17 uint64 = 0
	var bite_hit_303 uint64 = 0

	var threw_win_0 uint64 = 0
	var threw_win_1 uint64 = 0
	var threw_win_2 uint64 = 0
	var threw_win_3 uint64 = 0
	var threw_win_4 uint64 = 0
	var threw_win_5 uint64 = 0
	var threw_win_6 uint64 = 0
	var threw_win_7 uint64 = 0
	var threw_win_8 uint64 = 0
	var threw_win_9 uint64 = 0
	var threw_win_10 uint64 = 0
	var threw_win_11 uint64 = 0
	var threw_win_12 uint64 = 0
	var threw_win_13 uint64 = 0
	var threw_win_14 uint64 = 0
	var threw_win_15 uint64 = 0
	var threw_win_16 uint64 = 0
	var threw_win_17 uint64 = 0
	var threw_win_303 uint64 = 0

	var bite_win_0 uint64 = 0
	var bite_win_1 uint64 = 0
	var bite_win_2 uint64 = 0
	var bite_win_3 uint64 = 0
	var bite_win_4 uint64 = 0
	var bite_win_5 uint64 = 0
	var bite_win_6 uint64 = 0
	var bite_win_7 uint64 = 0
	var bite_win_8 uint64 = 0
	var bite_win_9 uint64 = 0
	var bite_win_10 uint64 = 0
	var bite_win_11 uint64 = 0
	var bite_win_12 uint64 = 0
	var bite_win_13 uint64 = 0
	var bite_win_14 uint64 = 0
	var bite_win_15 uint64 = 0
	var bite_win_16 uint64 = 0
	var bite_win_17 uint64 = 0
	var bite_win_303 uint64 = 0

	var threw_bullet_0 uint64 = 0
	var threw_bullet_1 uint64 = 0
	var threw_bullet_2 uint64 = 0
	var threw_bullet_3 uint64 = 0
	var threw_bullet_4 uint64 = 0
	var threw_bullet_5 uint64 = 0
	var threw_bullet_6 uint64 = 0
	var threw_bullet_7 uint64 = 0
	var threw_bullet_8 uint64 = 0
	var threw_bullet_9 uint64 = 0
	var threw_bullet_10 uint64 = 0
	var threw_bullet_11 uint64 = 0
	var threw_bullet_12 uint64 = 0
	var threw_bullet_13 uint64 = 0
	var threw_bullet_14 uint64 = 0
	var threw_bullet_15 uint64 = 0
	var threw_bullet_16 uint64 = 0
	var threw_bullet_17 uint64 = 0
	var threw_bullet_303 uint64 = 0

	var bite_bullet_0 uint64 = 0
	var bite_bullet_1 uint64 = 0
	var bite_bullet_2 uint64 = 0
	var bite_bullet_3 uint64 = 0
	var bite_bullet_4 uint64 = 0
	var bite_bullet_5 uint64 = 0
	var bite_bullet_6 uint64 = 0
	var bite_bullet_7 uint64 = 0
	var bite_bullet_8 uint64 = 0
	var bite_bullet_9 uint64 = 0
	var bite_bullet_10 uint64 = 0
	var bite_bullet_11 uint64 = 0
	var bite_bullet_12 uint64 = 0
	var bite_bullet_13 uint64 = 0
	var bite_bullet_14 uint64 = 0
	var bite_bullet_15 uint64 = 0
	var bite_bullet_16 uint64 = 0
	var bite_bullet_17 uint64 = 0
	var bite_bullet_303 uint64 = 0

	var weapon_1 uint64 = 0
	var weapon_3 uint64 = 0
	var weapon_5 uint64 = 0
	var weapon_10 uint64 = 0

	bet := betNormal

	fmt.Print("MathModuleId(95、96、97、98), RunTimes(x 萬), FishId -> ")
	var rtpId, run_times int
	var inputFishId string
	fmt.Scanf("%d, %d, %s", &rtpId, &run_times, &inputFishId)

	switch rtpId {
	case 95:
		mathModuleId_00003 = models.PSFM_00004_95_1
	case 96:
		mathModuleId_00003 = models.PSFM_00004_96_1
	case 97:
		mathModuleId_00003 = models.PSFM_00004_97_1
	case 98:
		mathModuleId_00003 = models.PSFM_00004_98_1
	}
	run_times = run_times * 10000

	for i := 0; i < run_times; i++ {
		// Process Weapon
		bet, weapon_1, weapon_3, weapon_5, weapon_10 = chooseWeapon(i, weapon_1, weapon_3, weapon_5, weapon_10)

		//fishId := rngFishId_00003()
		//fishId := int32(0)

		var fishId int32
		if inputFishId != "" {
			tempFish, _ := strconv.Atoi(inputFishId)
			fishId = int32(tempFish)
		} else {
			fishId = rngFishId_00003()
		}

		Service.Decrease(
			game_id_00003, subgame_id_00003, mathModuleId_00003, secWebSocketKey_00003,
			uint64(bet)*rate_00003,
		)
		rtpId := Service.RtpId(secWebSocketKey_00003, game_id_00003, subgame_id_00003)
		rtpState := Service.RtpState(secWebSocketKey_00003, game_id_00003, subgame_id_00003)

		for j := 0; j < bet; j++ {
			result := probability.Service.Calc(
				game_id_00003,
				mathModuleId_00003,
				rtpId,
				fishId,
				-1,
				0,
				bet_00003*rate_00003,
			)

			// Mercenary Bullet Collection Process
			//Service.IncreaseBulletCollection(secWebSocketKey_00003, game_id_00003, subgame_id_00003, uint64(result.Bullet))
			//bulletCollection, _ := Service.MercenaryInfo(secWebSocketKey_00003, game_id_00003, subgame_id_00003)
			//mercenaryPay := mercenaryProcess(secWebSocketKey_00003, bulletCollection, rtpId, fishId, t)

			switch fishId {
			case 0:
				threw_fish_0, threw_hit_0, threw_win_0, bite_fish_0, bite_hit_0, bite_win_0, threw_bullet_0, bite_bullet_0 = getResult(
					rtpState, t, result,
					threw_fish_0, threw_hit_0, threw_win_0, bite_fish_0, bite_hit_0, bite_win_0,
					threw_bullet_0, bite_bullet_0,
				)
				threw_bullet_0, bite_bullet_0 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_0, bite_bullet_0)

			case 1:
				threw_fish_1, threw_hit_1, threw_win_1, bite_fish_1, bite_hit_1, bite_win_1, threw_bullet_1, bite_bullet_1 = getResult(
					rtpState, t, result,
					threw_fish_1, threw_hit_1, threw_win_1, bite_fish_1, bite_hit_1, bite_win_1,
					threw_bullet_1, bite_bullet_1,
				)
				threw_bullet_1, bite_bullet_1 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_1, bite_bullet_1)

			case 2:
				threw_fish_2, threw_hit_2, threw_win_2, bite_fish_2, bite_hit_2, bite_win_2, threw_bullet_2, bite_bullet_2 = getResult(
					rtpState, t, result,
					threw_fish_2, threw_hit_2, threw_win_2, bite_fish_2, bite_hit_2, bite_win_2,
					threw_bullet_2, bite_bullet_2,
				)
				threw_bullet_2, bite_bullet_2 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_2, bite_bullet_2)

			case 3:
				threw_fish_3, threw_hit_3, threw_win_3, bite_fish_3, bite_hit_3, bite_win_3, threw_bullet_3, bite_bullet_3 = getResult(
					rtpState, t, result,
					threw_fish_3, threw_hit_3, threw_win_3, bite_fish_3, bite_hit_3, bite_win_3,
					threw_bullet_3, bite_bullet_3,
				)
				threw_bullet_3, bite_bullet_3 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_3, bite_bullet_3)

			case 4:
				threw_fish_4, threw_hit_4, threw_win_4, bite_fish_4, bite_hit_4, bite_win_4, threw_bullet_4, bite_bullet_4 = getResult(
					rtpState, t, result,
					threw_fish_4, threw_hit_4, threw_win_4, bite_fish_4, bite_hit_4, bite_win_4,
					threw_bullet_4, bite_bullet_4,
				)
				threw_bullet_4, bite_bullet_4 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_4, bite_bullet_4)

			case 5:
				threw_fish_5, threw_hit_5, threw_win_5, bite_fish_5, bite_hit_5, bite_win_5, threw_bullet_5, bite_bullet_5 = getResult(
					rtpState, t, result,
					threw_fish_5, threw_hit_5, threw_win_5, bite_fish_5, bite_hit_5, bite_win_5,
					threw_bullet_5, bite_bullet_5,
				)
				threw_bullet_5, bite_bullet_5 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_5, bite_bullet_5)

			case 6:
				threw_fish_6, threw_hit_6, threw_win_6, bite_fish_6, bite_hit_6, bite_win_6, threw_bullet_6, bite_bullet_6 = getResult(
					rtpState, t, result,
					threw_fish_6, threw_hit_6, threw_win_6, bite_fish_6, bite_hit_6, bite_win_6,
					threw_bullet_6, bite_bullet_6,
				)
				threw_bullet_6, bite_bullet_6 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_6, bite_bullet_6)

			case 7:
				threw_fish_7, threw_hit_7, threw_win_7, bite_fish_7, bite_hit_7, bite_win_7, threw_bullet_7, bite_bullet_7 = getResult(
					rtpState, t, result,
					threw_fish_7, threw_hit_7, threw_win_7, bite_fish_7, bite_hit_7, bite_win_7,
					threw_bullet_7, bite_bullet_7,
				)
				threw_bullet_7, bite_bullet_7 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_7, bite_bullet_7)

			case 8:
				threw_fish_8, threw_hit_8, threw_win_8, bite_fish_8, bite_hit_8, bite_win_8, threw_bullet_8, bite_bullet_8 = getResult(
					rtpState, t, result,
					threw_fish_8, threw_hit_8, threw_win_8, bite_fish_8, bite_hit_8, bite_win_8,
					threw_bullet_8, bite_bullet_8,
				)
				threw_bullet_8, bite_bullet_8 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_8, bite_bullet_8)

			case 9:
				threw_fish_9, threw_hit_9, threw_win_9, bite_fish_9, bite_hit_9, bite_win_9, threw_bullet_9, bite_bullet_9 = getResult(
					rtpState, t, result,
					threw_fish_9, threw_hit_9, threw_win_9, bite_fish_9, bite_hit_9, bite_win_9,
					threw_bullet_9, bite_bullet_9,
				)
				threw_bullet_9, bite_bullet_9 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_9, bite_bullet_9)

			case 10:
				threw_fish_10, threw_hit_10, threw_win_10, bite_fish_10, bite_hit_10, bite_win_10, threw_bullet_10, bite_bullet_10 = getResult(
					rtpState, t, result,
					threw_fish_10, threw_hit_10, threw_win_10, bite_fish_10, bite_hit_10, bite_win_10,
					threw_bullet_10, bite_bullet_10,
				)
				threw_bullet_10, bite_bullet_10 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_10, bite_bullet_10)

			case 11:
				threw_fish_11, threw_hit_11, threw_win_11, bite_fish_11, bite_hit_11, bite_win_11, threw_bullet_11, bite_bullet_11 = getResult(
					rtpState, t, result,
					threw_fish_11, threw_hit_11, threw_win_11, bite_fish_11, bite_hit_11, bite_win_11,
					threw_bullet_11, bite_bullet_11,
				)
				threw_bullet_11, bite_bullet_11 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_11, bite_bullet_11)

			case 12:
				threw_fish_12, threw_hit_12, threw_win_12, bite_fish_12, bite_hit_12, bite_win_12, threw_bullet_12, bite_bullet_12 = getResult(
					rtpState, t, result,
					threw_fish_12, threw_hit_12, threw_win_12, bite_fish_12, bite_hit_12, bite_win_12,
					threw_bullet_12, bite_bullet_12,
				)
				threw_bullet_12, bite_bullet_12 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_12, bite_bullet_12)

			case 13:
				threw_fish_13, threw_hit_13, threw_win_13, bite_fish_13, bite_hit_13, bite_win_13, threw_bullet_13, bite_bullet_13 = getResult(
					rtpState, t, result,
					threw_fish_13, threw_hit_13, threw_win_13, bite_fish_13, bite_hit_13, bite_win_13,
					threw_bullet_13, bite_bullet_13,
				)
				threw_bullet_13, bite_bullet_13 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_13, bite_bullet_13)

			case 14:
				threw_fish_14, threw_hit_14, threw_win_14, bite_fish_14, bite_hit_14, bite_win_14, threw_bullet_14, bite_bullet_14 = getResult(
					rtpState, t, result,
					threw_fish_14, threw_hit_14, threw_win_14, bite_fish_14, bite_hit_14, bite_win_14,
					threw_bullet_14, bite_bullet_14,
				)
				threw_bullet_14, bite_bullet_14 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_14, bite_bullet_14)

			case 15:
				threw_fish_15, threw_hit_15, threw_win_15, bite_fish_15, bite_hit_15, bite_win_15, threw_bullet_15, bite_bullet_15 = getResult(
					rtpState, t, result,
					threw_fish_15, threw_hit_15, threw_win_15, bite_fish_15, bite_hit_15, bite_win_15,
					threw_bullet_15, bite_bullet_15,
				)
				threw_bullet_15, bite_bullet_15 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_15, bite_bullet_15)

			case 16:
				threw_fish_16, threw_hit_16, threw_win_16, bite_fish_16, bite_hit_16, bite_win_16, threw_bullet_16, bite_bullet_16 = getResult(
					rtpState, t, result,
					threw_fish_16, threw_hit_16, threw_win_16, bite_fish_16, bite_hit_16, bite_win_16,
					threw_bullet_16, bite_bullet_16,
				)
				threw_bullet_16, bite_bullet_16 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_16, bite_bullet_16)

			case 17:
				threw_fish_17, threw_hit_17, threw_win_17, bite_fish_17, bite_hit_17, bite_win_17, threw_bullet_17, bite_bullet_17 = getResult(
					rtpState, t, result,
					threw_fish_17, threw_hit_17, threw_win_17, bite_fish_17, bite_hit_17, bite_win_17,
					threw_bullet_17, bite_bullet_17,
				)
				threw_bullet_17, bite_bullet_17 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_17, bite_bullet_17)

			case 303:
				switch rtpState {
				case threw:
					if result.Pay > 0 || result.Bullet > 0 {
						threw_hit_303++
						threw_win_303 += uint64(result.Pay * result.Multiplier)
					}
					threw_fish_303++

				case bite:
					if result.Pay > 0 || result.Bullet > 0 {
						bite_hit_303++
						bite_win_303 += uint64(result.Pay * result.Multiplier)
					}
					bite_fish_303++

				default:
					t.Fatal("RTP State Error")
				}

				threw_bullet_303, bite_bullet_303 = bulletResult(rtpState, t, uint64(result.Bullet), threw_bullet_303, bite_bullet_303)
			}
		}
	}

	threw_win_0, bite_win_0, threw_bullet_0, bite_bullet_0 = mercenaryBulletProcess(threw_win_0, bite_win_0, threw_bullet_0, bite_bullet_0)
	threw_win_1, bite_win_1, threw_bullet_1, bite_bullet_1 = mercenaryBulletProcess(threw_win_1, bite_win_1, threw_bullet_1, bite_bullet_1)
	threw_win_2, bite_win_2, threw_bullet_2, bite_bullet_2 = mercenaryBulletProcess(threw_win_2, bite_win_2, threw_bullet_2, bite_bullet_2)
	threw_win_3, bite_win_3, threw_bullet_3, bite_bullet_3 = mercenaryBulletProcess(threw_win_3, bite_win_3, threw_bullet_3, bite_bullet_3)
	threw_win_4, bite_win_4, threw_bullet_4, bite_bullet_4 = mercenaryBulletProcess(threw_win_4, bite_win_4, threw_bullet_4, bite_bullet_4)
	threw_win_5, bite_win_5, threw_bullet_5, bite_bullet_5 = mercenaryBulletProcess(threw_win_5, bite_win_5, threw_bullet_5, bite_bullet_5)
	threw_win_6, bite_win_6, threw_bullet_6, bite_bullet_6 = mercenaryBulletProcess(threw_win_6, bite_win_6, threw_bullet_6, bite_bullet_6)
	threw_win_7, bite_win_7, threw_bullet_7, bite_bullet_7 = mercenaryBulletProcess(threw_win_7, bite_win_7, threw_bullet_7, bite_bullet_7)
	threw_win_8, bite_win_8, threw_bullet_8, bite_bullet_8 = mercenaryBulletProcess(threw_win_8, bite_win_8, threw_bullet_8, bite_bullet_8)
	threw_win_9, bite_win_9, threw_bullet_9, bite_bullet_9 = mercenaryBulletProcess(threw_win_9, bite_win_9, threw_bullet_9, bite_bullet_9)
	threw_win_10, bite_win_10, threw_bullet_10, bite_bullet_10 = mercenaryBulletProcess(threw_win_10, bite_win_10, threw_bullet_10, bite_bullet_10)
	threw_win_11, bite_win_11, threw_bullet_11, bite_bullet_11 = mercenaryBulletProcess(threw_win_11, bite_win_11, threw_bullet_11, bite_bullet_11)
	threw_win_12, bite_win_12, threw_bullet_12, bite_bullet_12 = mercenaryBulletProcess(threw_win_12, bite_win_12, threw_bullet_12, bite_bullet_12)
	threw_win_13, bite_win_13, threw_bullet_13, bite_bullet_13 = mercenaryBulletProcess(threw_win_13, bite_win_13, threw_bullet_13, bite_bullet_13)
	threw_win_14, bite_win_14, threw_bullet_14, bite_bullet_14 = mercenaryBulletProcess(threw_win_14, bite_win_14, threw_bullet_14, bite_bullet_14)
	threw_win_15, bite_win_15, threw_bullet_15, bite_bullet_15 = mercenaryBulletProcess(threw_win_15, bite_win_15, threw_bullet_15, bite_bullet_15)
	threw_win_16, bite_win_16, threw_bullet_16, bite_bullet_16 = mercenaryBulletProcess(threw_win_16, bite_win_16, threw_bullet_16, bite_bullet_16)
	threw_win_17, bite_win_17, threw_bullet_17, bite_bullet_17 = mercenaryBulletProcess(threw_win_17, bite_win_17, threw_bullet_17, bite_bullet_17)
	threw_win_303, bite_win_303, threw_bullet_303, bite_bullet_303 = mercenaryBulletProcess(threw_win_303, bite_win_303, threw_bullet_303, bite_bullet_303)

	t.Log("Weapon", weapon_1, weapon_3, weapon_5, weapon_10)

	total_threw_times := threw_fish_0 + threw_fish_1 + threw_fish_2 + threw_fish_3 + threw_fish_4 + threw_fish_5 + threw_fish_6 + threw_fish_7 +
		threw_fish_8 + threw_fish_9 + threw_fish_10 + threw_fish_11 + threw_fish_12 + threw_fish_13 + threw_fish_14 + threw_fish_15 +
		threw_fish_16 + threw_fish_17 + threw_fish_303
	total_bite_times := bite_fish_0 + bite_fish_1 + bite_fish_2 + bite_fish_3 + bite_fish_4 + bite_fish_5 + bite_fish_6 + bite_fish_7 +
		bite_fish_8 + bite_fish_9 + bite_fish_10 + bite_fish_11 + bite_fish_12 + bite_fish_13 + bite_fish_14 + bite_fish_15 +
		bite_fish_16 + bite_fish_17 + bite_fish_303
	total_threw_win := threw_win_0 + threw_win_1 + threw_win_2 + threw_win_3 + threw_win_4 + threw_win_5 + threw_win_6 + threw_win_7 +
		threw_win_8 + threw_win_9 + threw_win_10 + threw_win_11 + threw_win_12 + threw_win_13 + threw_win_14 + threw_win_15 +
		threw_win_16 + threw_win_17 + threw_win_303
	total_bite_win := bite_win_0 + bite_win_1 + bite_win_2 + bite_win_3 + bite_win_4 + bite_win_5 + bite_win_6 + bite_win_7 +
		bite_win_8 + bite_win_9 + bite_win_10 + bite_win_11 + bite_win_12 + bite_win_13 + bite_win_14 + bite_win_15 +
		bite_win_16 + bite_win_17 + bite_win_303
	total_threw_hit := threw_hit_0 + threw_hit_1 + threw_hit_2 + threw_hit_3 + threw_hit_4 + threw_hit_5 + threw_hit_6 + threw_hit_7 +
		threw_hit_8 + threw_hit_9 + threw_hit_10 + threw_hit_11 + threw_hit_12 + threw_hit_13 + threw_hit_14 + threw_hit_15 +
		threw_hit_16 + threw_hit_17 + threw_hit_303
	total_bite_hit := bite_hit_0 + bite_hit_1 + bite_hit_2 + bite_hit_3 + bite_hit_4 + bite_hit_5 + bite_hit_6 + bite_hit_7 +
		bite_hit_8 + bite_hit_9 + bite_hit_10 + bite_hit_11 + bite_hit_12 + bite_hit_13 + bite_hit_14 + bite_hit_15 +
		bite_hit_16 + bite_hit_17 + bite_hit_303
	total_threw_bullet := threw_bullet_0 + threw_bullet_1 + threw_bullet_2 + threw_bullet_3 + threw_bullet_4 + threw_bullet_5 + threw_bullet_6 + threw_bullet_7 +
		threw_bullet_8 + threw_bullet_9 + threw_bullet_10 + threw_bullet_11 + threw_bullet_12 + threw_bullet_13 + threw_bullet_14 + threw_bullet_15 +
		threw_bullet_16 + threw_bullet_17 + threw_bullet_303
	total_bite_bullet := bite_bullet_0 + bite_bullet_1 + bite_bullet_2 + bite_bullet_3 + bite_bullet_4 + bite_bullet_5 + bite_bullet_6 + bite_bullet_7 +
		bite_bullet_8 + bite_bullet_9 + bite_bullet_10 + bite_bullet_11 + bite_bullet_12 + bite_bullet_13 + bite_bullet_14 + bite_bullet_15 +
		bite_bullet_16 + bite_bullet_17 + bite_bullet_303

	fmt.Println("FishId", ":", "ThrewTimes", "ThrewHit", "ThrewWin", "ThrewBullet", "BiteTimes", "BiteHit", "BiteWin", "BitBullet",
		"ThrewHitRate", "ThrewPayRate", "BiteHitRate", "BitePayRate", "TotalPayRate")
	fmt.Println("000", ":", threw_fish_0, threw_hit_0, threw_win_0, threw_bullet_0, bite_fish_0, bite_hit_0, bite_win_0, bite_bullet_0,
		getRate(threw_hit_0, threw_fish_0), getRate(threw_win_0+threw_bullet_0, threw_fish_0),
		getRate(bite_hit_0, bite_fish_0), getRate(bite_win_0+bite_bullet_0, bite_fish_0),
		getRate(threw_win_0+threw_bullet_0+bite_win_0+bite_bullet_0, threw_fish_0+bite_fish_0))
	fmt.Println("001", ":", threw_fish_1, threw_hit_1, threw_win_1, threw_bullet_1, bite_fish_1, bite_hit_1, bite_win_1, bite_bullet_1,
		getRate(threw_hit_1, threw_fish_1), getRate(threw_win_1+threw_bullet_1, threw_fish_1),
		getRate(bite_hit_1, bite_fish_1), getRate(bite_win_1+bite_bullet_1, bite_fish_1),
		getRate(threw_win_1+threw_bullet_1+bite_win_1+bite_bullet_1, threw_fish_1+bite_fish_1))
	fmt.Println("002", ":", threw_fish_2, threw_hit_2, threw_win_2, threw_bullet_2, bite_fish_2, bite_hit_2, bite_win_2, bite_bullet_2,
		getRate(threw_hit_2, threw_fish_2), getRate(threw_win_2+threw_bullet_2, threw_fish_2),
		getRate(bite_hit_2, bite_fish_2), getRate(bite_win_2+bite_bullet_2, bite_fish_2),
		getRate(threw_win_2+threw_bullet_2+bite_win_2+bite_bullet_2, threw_fish_2+bite_fish_2))
	fmt.Println("003", ":", threw_fish_3, threw_hit_3, threw_win_3, threw_bullet_3, bite_fish_3, bite_hit_3, bite_win_3, bite_bullet_3,
		getRate(threw_hit_3, threw_fish_3), getRate(threw_win_3+threw_bullet_3, threw_fish_3),
		getRate(bite_hit_3, bite_fish_3), getRate(bite_win_3+bite_bullet_3, bite_fish_3),
		getRate(threw_win_3+threw_bullet_3+bite_win_3+bite_bullet_3, threw_fish_3+bite_fish_3))
	fmt.Println("004", ":", threw_fish_4, threw_hit_4, threw_win_4, threw_bullet_4, bite_fish_4, bite_hit_4, bite_win_4, bite_bullet_4,
		getRate(threw_hit_4, threw_fish_4), getRate(threw_win_4+threw_bullet_4, threw_fish_4),
		getRate(bite_hit_4, bite_fish_4), getRate(bite_win_4+bite_bullet_4, bite_fish_4),
		getRate(threw_win_4+threw_bullet_4+bite_win_4+bite_bullet_4, threw_fish_4+bite_fish_4))
	fmt.Println("005", ":", threw_fish_5, threw_hit_5, threw_win_5, threw_bullet_5, bite_fish_5, bite_hit_5, bite_win_5, bite_bullet_5,
		getRate(threw_hit_5, threw_fish_5), getRate(threw_win_5+threw_bullet_5, threw_fish_5),
		getRate(bite_hit_5, bite_fish_5), getRate(bite_win_5+bite_bullet_5, bite_fish_5),
		getRate(threw_win_5+threw_bullet_5+bite_win_5+bite_bullet_5, threw_fish_5+bite_fish_5))
	fmt.Println("006", ":", threw_fish_6, threw_hit_6, threw_win_6, threw_bullet_6, bite_fish_6, bite_hit_6, bite_win_6, bite_bullet_6,
		getRate(threw_hit_6, threw_fish_6), getRate(threw_win_6+threw_bullet_6, threw_fish_6),
		getRate(bite_hit_6, bite_fish_6), getRate(bite_win_6+bite_bullet_6, bite_fish_6),
		getRate(threw_win_6+threw_bullet_6+bite_win_6+bite_bullet_6, threw_fish_6+bite_fish_6))
	fmt.Println("007", ":", threw_fish_7, threw_hit_7, threw_win_7, threw_bullet_7, bite_fish_7, bite_hit_7, bite_win_7, bite_bullet_7,
		getRate(threw_hit_7, threw_fish_7), getRate(threw_win_7+threw_bullet_7, threw_fish_7),
		getRate(bite_hit_7, bite_fish_7), getRate(bite_win_7+bite_bullet_7, bite_fish_7),
		getRate(threw_win_7+threw_bullet_7+bite_win_7+bite_bullet_7, threw_fish_7+bite_fish_7))
	fmt.Println("008", ":", threw_fish_8, threw_hit_8, threw_win_8, threw_bullet_8, bite_fish_8, bite_hit_8, bite_win_8, bite_bullet_8,
		getRate(threw_hit_8, threw_fish_8), getRate(threw_win_8+threw_bullet_8, threw_fish_8),
		getRate(bite_hit_8, bite_fish_8), getRate(bite_win_8+bite_bullet_8, bite_fish_8),
		getRate(threw_win_8+threw_bullet_8+bite_win_8+bite_bullet_8, threw_fish_8+bite_fish_8))
	fmt.Println("009", ":", threw_fish_9, threw_hit_9, threw_win_9, threw_bullet_9, bite_fish_9, bite_hit_9, bite_win_9, bite_bullet_9,
		getRate(threw_hit_9, threw_fish_9), getRate(threw_win_9+threw_bullet_9, threw_fish_9),
		getRate(bite_hit_9, bite_fish_9), getRate(bite_win_9+bite_bullet_9, bite_fish_9),
		getRate(threw_win_9+threw_bullet_9+bite_win_9+bite_bullet_9, threw_fish_9+bite_fish_9))
	fmt.Println("010", ":", threw_fish_10, threw_hit_10, threw_win_10, threw_bullet_10, bite_fish_10, bite_hit_10, bite_win_10, bite_bullet_10,
		getRate(threw_hit_10, threw_fish_10), getRate(threw_win_10+threw_bullet_10, threw_fish_10),
		getRate(bite_hit_10, bite_fish_10), getRate(bite_win_10+bite_bullet_10, bite_fish_10),
		getRate(threw_win_10+threw_bullet_10+bite_win_10+bite_bullet_10, threw_fish_10+bite_fish_10))
	fmt.Println("011", ":", threw_fish_11, threw_hit_11, threw_win_11, threw_bullet_11, bite_fish_11, bite_hit_11, bite_win_11, bite_bullet_11,
		getRate(threw_hit_11, threw_fish_11), getRate(threw_win_11+threw_bullet_11, threw_fish_11),
		getRate(bite_hit_11, bite_fish_11), getRate(bite_win_11+bite_bullet_11, bite_fish_11),
		getRate(threw_win_11+threw_bullet_11+bite_win_11+bite_bullet_11, threw_fish_11+bite_fish_11))
	fmt.Println("012", ":", threw_fish_12, threw_hit_12, threw_win_12, threw_bullet_12, bite_fish_12, bite_hit_12, bite_win_12, bite_bullet_12,
		getRate(threw_hit_12, threw_fish_12), getRate(threw_win_12+threw_bullet_12, threw_fish_12),
		getRate(bite_hit_12, bite_fish_12), getRate(bite_win_12+bite_bullet_12, bite_fish_12),
		getRate(threw_win_12+threw_bullet_12+bite_win_12+bite_bullet_12, threw_fish_12+bite_fish_12))
	fmt.Println("013", ":", threw_fish_13, threw_hit_13, threw_win_13, threw_bullet_13, bite_fish_13, bite_hit_13, bite_win_13, bite_bullet_13,
		getRate(threw_hit_13, threw_fish_13), getRate(threw_win_13+threw_bullet_13, threw_fish_13),
		getRate(bite_hit_13, bite_fish_13), getRate(bite_win_13+bite_bullet_13, bite_fish_13),
		getRate(threw_win_13+threw_bullet_13+bite_win_13+bite_bullet_13, threw_fish_13+bite_fish_13))
	fmt.Println("014", ":", threw_fish_14, threw_hit_14, threw_win_14, threw_bullet_14, bite_fish_14, bite_hit_14, bite_win_14, bite_bullet_14,
		getRate(threw_hit_14, threw_fish_14), getRate(threw_win_14+threw_bullet_14, threw_fish_14),
		getRate(bite_hit_14, bite_fish_14), getRate(bite_win_14+bite_bullet_14, bite_fish_14),
		getRate(threw_win_14+threw_bullet_14+bite_win_14+bite_bullet_14, threw_fish_14+bite_fish_14))
	fmt.Println("015", ":", threw_fish_15, threw_hit_15, threw_win_15, threw_bullet_15, bite_fish_15, bite_hit_15, bite_win_15, bite_bullet_15,
		getRate(threw_hit_15, threw_fish_15), getRate(threw_win_15+threw_bullet_15, threw_fish_15),
		getRate(bite_hit_15, bite_fish_15), getRate(bite_win_15+bite_bullet_15, bite_fish_15),
		getRate(threw_win_15+threw_bullet_15+bite_win_15+bite_bullet_15, threw_fish_15+bite_fish_15))
	fmt.Println("016", ":", threw_fish_16, threw_hit_16, threw_win_16, threw_bullet_16, bite_fish_16, bite_hit_16, bite_win_16, bite_bullet_16,
		getRate(threw_hit_16, threw_fish_16), getRate(threw_win_16+threw_bullet_16, threw_fish_16),
		getRate(bite_hit_16, bite_fish_16), getRate(bite_win_16+bite_bullet_16, bite_fish_16),
		getRate(threw_win_16+threw_bullet_16+bite_win_16+bite_bullet_16, threw_fish_16+bite_fish_16))
	fmt.Println("017", ":", threw_fish_17, threw_hit_17, threw_win_17, threw_bullet_17, bite_fish_17, bite_hit_17, bite_win_17, bite_bullet_17,
		getRate(threw_hit_17, threw_fish_17), getRate(threw_win_17+threw_bullet_17, threw_fish_17),
		getRate(bite_hit_17, bite_fish_17), getRate(bite_win_17+bite_bullet_17, bite_fish_17),
		getRate(threw_win_17+threw_bullet_17+bite_win_17+bite_bullet_17, threw_fish_17+bite_fish_17))
	fmt.Println("303", ":", threw_fish_303, threw_hit_303, threw_win_303, threw_bullet_303, bite_fish_303, bite_hit_303, bite_win_303, bite_bullet_303,
		getRate(threw_hit_303, threw_fish_303), getRate(threw_win_303+threw_bullet_303, threw_fish_303),
		getRate(bite_hit_303, bite_fish_303), getRate(bite_win_303+bite_bullet_303, bite_fish_303),
		getRate(threw_win_303+threw_bullet_303+bite_win_303+bite_bullet_303, threw_fish_303+bite_fish_303))
	fmt.Println("Total", total_threw_times, total_threw_hit, total_threw_win, total_threw_bullet,
		total_bite_times, total_bite_hit, total_bite_win, total_bite_bullet,
		getRate(total_threw_win+total_threw_bullet+total_bite_win+total_bite_bullet, total_threw_times+total_bite_times))
}

func mercenaryBulletProcess(threw_win, bite_win, threw_bullet, bite_bullet uint64) (threwWin, biteWin, threwBullet, biteBullet uint64) {
	for {
		threw_win, bite_win, threw_bullet = mercenary(threw, threw_bullet, threw_win, bite_win)
		if threw_bullet < 60 {
			break
		}
	}

	for {
		threw_win, bite_win, bite_bullet = mercenary(bite, bite_bullet, threw_win, bite_win)
		if bite_bullet < 60 {
			break
		}
	}

	return threw_win, bite_win, threw_bullet, bite_bullet
}

func rngFishId_00003() int32 {
	options := make([]rng.Option, 0, 19)

	for i := 0; i < 18; i++ {
		options = append(options, rng.Option{Weight: 1, Item: i})
	}
	options = append(options, rng.Option{Weight: 1, Item: 303})

	return int32(rng.New(options).Item.(int))
}

func shot_Drill(rtpId string, bullets int) (pay, bullet uint64) {
	var fishId int32 = -1

	for i := 0; i < bullets; i++ {
		// Random Fish ID
		for {
			fishId = rngFishId_00003()

			if fishId != 12 && fishId != 13 && fishId != 14 &&
				fishId != 15 && fishId != 16 && fishId != 17 {
				break
			}
		}

		result := probability.Service.Calc(
			game_id_00003,
			mathModuleId_00003,
			rtpId,
			fishId,
			200,
			Service.RtpBudget(secWebSocketKey_00003, game_id_00003, subgame_id_00003, mathModuleId_00003),
			0,
		)

		if result.Pay > 0 {
			pay += uint64(result.Pay * result.Multiplier)

			// Slot
			if fishId == 6 || fishId == 7 || fishId == 8 ||
				fishId == 9 || fishId == 10 || fishId == 11 {
				pay += uint64(result.ExtraData[0].(int))
			}
		}

		if result.Bullet > 0 {
			bullet += uint64(result.Bullet)
		}
	}

	return pay, bullet
}

func getResult(rtpState int, t *testing.T, result *probability.Probability,
	threw_fish, threw_hit, threw_win, bite_fish, bite_hit, bite_win, threw_bullet, bite_bullet uint64,
) (threwFish, threwHit, threwWin, biteFish, biteHit, biteWin, threwBullet, biteBullet uint64,
) {
	switch rtpState {
	case threw:
		biteFish = bite_fish
		biteHit = bite_hit
		biteWin = bite_win
		biteBullet = bite_bullet

		if result.Pay > 0 {
			pay, bullet := getWin(result, threw_win)

			threwHit = threw_hit + 1
			threwWin = pay
			threwBullet = threw_bullet + bullet
		} else {
			threwHit = threw_hit
			threwWin = threw_win
			threwBullet = threw_bullet
		}
		threwFish = threw_fish + 1

	case bite:
		threwFish = threw_fish
		threwHit = threw_hit
		threwWin = threw_win
		threwBullet = threw_bullet

		if result.Pay > 0 {
			pay, bullet := getWin(result, bite_win)

			biteHit = bite_hit + 1
			biteWin = pay
			biteBullet = bite_bullet + bullet
		} else {
			biteHit = bite_hit
			biteWin = bite_win
			biteBullet = bite_bullet
		}
		biteFish = bite_fish + 1

	default:
		t.Fatal("RTP State Error")
	}

	return threwFish, threwHit, threwWin, biteFish, biteHit, biteWin, threwBullet, biteBullet
}

func getWin(result *probability.Probability, inputWin uint64) (win, bullet uint64) {
	switch result.FishTypeId {
	case 6, 7, 8, 9, 10, 11:
		win = inputWin + uint64(result.Pay*result.Multiplier) + uint64(result.ExtraData[0].(int))
		bullet = 0
	case 12, 13, 14, 15, 16, 17:
		pay, drillBullet := shot_Drill(getDrillRtpID(result.FishTypeId), result.BonusPayload.(int))
		win = inputWin + uint64(result.Pay*result.Multiplier) + pay
		bullet = drillBullet
	default:
		win = inputWin + uint64(result.Pay*result.Multiplier)
		bullet = 0
	}

	return win, bullet
}

func getDrillRtpID(fishId int) string {
	switch mathModuleId_00003 {
	case models.PSFM_00004_95_1:
		switch fishId {
		case 12:
			return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
		case 13:
			return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie2Drill.UseRtp
		case 14:
			return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie3Drill.UseRtp
		case 15:
			return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie4Drill.UseRtp
		case 16:
			return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie5Drill.UseRtp
		case 17:
			return PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Zombie6Drill.UseRtp
		default:
			return ""
		}

	case models.PSFM_00004_96_1:
		switch fishId {
		case 12:
			return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
		case 13:
			return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie2Drill.UseRtp
		case 14:
			return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie3Drill.UseRtp
		case 15:
			return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie4Drill.UseRtp
		case 16:
			return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie5Drill.UseRtp
		case 17:
			return PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Zombie6Drill.UseRtp
		default:
			return ""
		}

	case models.PSFM_00004_97_1:
		switch fishId {
		case 12:
			return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
		case 13:
			return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie2Drill.UseRtp
		case 14:
			return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie3Drill.UseRtp
		case 15:
			return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie4Drill.UseRtp
		case 16:
			return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie5Drill.UseRtp
		case 17:
			return PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Zombie6Drill.UseRtp
		default:
			return ""
		}

	case models.PSFM_00004_98_1:
		switch fishId {
		case 12:
			return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie1Drill.UseRtp
		case 13:
			return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie2Drill.UseRtp
		case 14:
			return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie3Drill.UseRtp
		case 15:
			return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie4Drill.UseRtp
		case 16:
			return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie5Drill.UseRtp
		case 17:
			return PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Zombie6Drill.UseRtp
		default:
			return ""
		}

	default:
		return ""
	}
}

func rngWeapon() int {
	options := make([]rng.Option, 0, 4)

	options = append(options, rng.Option{Weight: 1, Item: betNormal})
	options = append(options, rng.Option{Weight: 1, Item: betShotGun})
	options = append(options, rng.Option{Weight: 1, Item: betBazooka})
	options = append(options, rng.Option{Weight: 1, Item: betMortar})

	return rng.New(options).Item.(int)
}

func rngOpenMercenary(mercenarySize int) (isOpen bool, mercenaryType int) {
	options := make([]rng.Option, 0, 2)
	options = append(options, rng.Option{Weight: 1, Item: true})
	options = append(options, rng.Option{Weight: 1, Item: false})

	mercenaryOptions := make([]rng.Option, 0, mercenarySize)
	for i := 0; i < mercenarySize; i++ {
		mercenaryOptions = append(mercenaryOptions, rng.Option{Weight: 1, Item: i})
	}

	return rng.New(options).Item.(bool), rng.New(mercenaryOptions).Item.(int)
}

func bulletResult(rtpState int, t *testing.T, bullet uint64, threwBullet, biteBullet uint64) (threwResult, biteResult uint64) {
	switch rtpState {
	case threw:
		threwResult = threwBullet + bullet
		biteResult = biteBullet
	case bite:
		threwResult = threwBullet
		biteResult = biteBullet + bullet
	default:
		t.Fatal("RTP State Error")
	}

	return threwResult, biteResult
}

func mercenary(rtpState int, bullets uint64, threw_win, bite_win uint64) (threwWin, biteWin, leftBullets uint64) {
	rtpId := ""
	switch mathModuleId_00003 {
	case models.PSFM_00004_95_1:
		rtpId = PSFM_00004_95_1.RTP95BS.PSF_ON_00003_1_BsMath.Icons.Bullets.UseRtp
	case models.PSFM_00004_96_1:
		rtpId = PSFM_00004_96_1.RTP96BS.PSF_ON_00003_1_BsMath.Icons.Bullets.UseRtp
	case models.PSFM_00004_97_1:
		rtpId = PSFM_00004_97_1.RTP97BS.PSF_ON_00003_1_BsMath.Icons.Bullets.UseRtp
	case models.PSFM_00004_98_1:
		rtpId = PSFM_00004_98_1.RTP98BS.PSF_ON_00003_1_BsMath.Icons.Bullets.UseRtp
	}

	times := bullets / 60
	leftBullet := bullets - (times * 60)
	var mercenaryGetBullet uint64 = 0
	var threwPay, bitePay uint64 = 0, 0

	for i := 0; uint64(i) < (times * 60); i++ {
		fishId := rngFishId_00003()

		result := probability.Service.Calc(
			game_id_00003,
			mathModuleId_00003,
			rtpId,
			fishId,
			200,
			Service.RtpBudget(secWebSocketKey_00003, game_id_00003, subgame_id_00003, mathModuleId_00003),
			0,
		)

		switch rtpState {
		case threw:
			if result.Pay > 0 {
				threwPay += uint64(result.Pay * result.Multiplier)
				// Slot
				if fishId == 6 || fishId == 7 || fishId == 8 ||
					fishId == 9 || fishId == 10 || fishId == 11 {
					threwPay += uint64(result.ExtraData[0].(int))
				}

				// Drill
				if fishId == 12 || fishId == 13 || fishId == 14 ||
					fishId == 15 || fishId == 16 || fishId == 17 {
					pay, drillBullet := shot_Drill(getDrillRtpID(result.FishTypeId), result.BonusPayload.(int))
					threwPay += pay
					mercenaryGetBullet += drillBullet
				}
			}

			if result.Bullet > 0 {
				mercenaryGetBullet += uint64(result.Bullet)
			}

		case bite:
			if result.Pay > 0 {
				bitePay += uint64(result.Pay * result.Multiplier)
				// Slot
				if fishId == 6 || fishId == 7 || fishId == 8 ||
					fishId == 9 || fishId == 10 || fishId == 11 {
					bitePay += uint64(result.ExtraData[0].(int))
				}

				// Drill
				if fishId == 12 || fishId == 13 || fishId == 14 ||
					fishId == 15 || fishId == 16 || fishId == 17 {
					pay, drillBullet := shot_Drill(getDrillRtpID(result.FishTypeId), result.BonusPayload.(int))
					bitePay += pay
					mercenaryGetBullet += drillBullet
				}
			}

			if result.Bullet > 0 {
				mercenaryGetBullet += uint64(result.Bullet)
			}
		}

	}

	threwWin = threw_win + threwPay
	biteWin = bite_win + bitePay

	return threwWin, biteWin, leftBullet + mercenaryGetBullet
}

func chooseWeapon(index int, weapon1, weapon3, weapon5, weapon10 uint64) (bet int, rweapon1, rweapon3, rweapon5, rweapon10 uint64) {
	bet = betNormal

	if (index+1)%switchWeaponBullet == 0 {
		bet = rngWeapon()

		switch bet {
		case 1:
			rweapon1 = weapon1 + 1
		case 3:
			rweapon3 = weapon3 + 1
		case 5:
			rweapon5 = weapon5 + 1
		case 10:
			rweapon10 = weapon10 + 1
		}
	}

	return bet, rweapon1, rweapon3, rweapon5, rweapon10
}
