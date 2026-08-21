<template>

  <!-- 添加弹框 -->
    <el-dialog v-model="addDialog"
               :title="`${GameManufacturer}${$t('游戏设置')}`"
               @open-auto-focus="openDialog"
               destroy-on-close
               :align-center="true"
               :width="store.viewModel === 2 ? '100%' : '750px'" @close="emits('update:modelValue') + ':'">

        <gameSet ref="gameSetRef" :merchantInfo="gameInfo.merchantInfo" :gameInfo="gameInfo.gameInfo"/>


        <gameBet ref="gameBetRef" :merchantInfo="gameInfo.merchantInfo" :gameInfo="gameInfo.gameInfo"
                 :table-data="tableData"/>
        <gameDefault ref="gameDefaultRef" :merchantInfo="gameInfo.merchantInfo" :gameInfo="gameInfo.gameInfo"
                     :table-data="tableData" v-if="GameManufacturer =='PG' || GameManufacturer =='HACKSAW'"/>

        <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('确认') }}</el-button>
                </span>
        </template>
    </el-dialog>


</template>


<script setup lang="ts">

import {computed, onMounted, Ref, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {Merchant, MerchantGameInterface, updateGameConfigParams} from "@/api/adminpb/merchant";
import {Client} from "@/lib/client";
import {tip} from "@/lib/tip";
import {RTP_VALUE} from "@/lib/RTP_config";
import {ElMessageBox} from "element-plus";
import GameSet from "@/pages/systemManagement/component/gameSettingInfoComponent/gameSet.vue";
import GameBet from "@/pages/systemManagement/component/gameSettingInfoComponent/gameBet.vue";
import GameDefault from "@/pages/systemManagement/component/gameSettingInfoComponent/gameDefault.vue";


let {t} = useI18n()
const store = useStore()

const props = defineProps(["modelValue", "gameSettingInfoData"])
const emits = defineEmits(['update:modelValue', 'update:gameSettingInfoData'])

const gameSetRef = ref(null)
const gameBetRef = ref(null)

const GameManufacturer = ref("")
const linNum = ref(1)

const openDialog = () => {
    addForm.value = {
        GameId: "",
        Gametype: 0,
        GameName: "",
        ConfigPath: "",
        RTP: RTP_VALUE[0],
        StopLoss: 1,
        MaxMultipleOff: 0,
        MaxMultiple: 10000,
        BetBase: "",
        GamePattern: 1,
        Preset: 0,
        GameOn: 0,
        FreeGameOff: 0,
        GameDemo: "",
    }
    tableData.value = []
}

const gameInfo = computed({
    get() {



        let data = {...props.gameSettingInfoData}

        tableData.value = []
        let BetBase = []

        if (data.gameInfo.BetBase != "") {

            BetBase = data.gameInfo.BetBase.split(",")

        }
        GameManufacturer.value = data.gameInfo.GameManufacturer

        linNum.value = store.LineMap[data.gameInfo.GameId] || 1


        for (const i in BetBase) {

            BetBase[i] = parseFloat(BetBase[i])
            if (GameManufacturer.value.toLowerCase() == "jili" || GameManufacturer.value.toLowerCase() == "tada") {
                BetBase[i] = (parseFloat(BetBase[i]) / linNum.value).toFixed(4)
            }


            tableData.value.push({
                BetValue: parseFloat(BetBase[i]),
                BetInfo: [
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 1,
                        result: parseFloat(BetBase[i]) * linNum.value * 1,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 2,
                        result: parseFloat(BetBase[i]) * linNum.value * 2,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 3,
                        result: parseFloat(BetBase[i]) * linNum.value * 3,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 5,
                        result: parseFloat(BetBase[i]) * linNum.value * 5,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 6,
                        result: parseFloat(BetBase[i]) * linNum.value * 6,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 7,
                        result: parseFloat(BetBase[i]) * linNum.value * 7,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 8,
                        result: parseFloat(BetBase[i]) * linNum.value * 8,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 9,
                        result: parseFloat(BetBase[i]) * linNum.value * 9,
                        linNum: linNum.value
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 10,
                        result: parseFloat(BetBase[i]) * linNum.value * 10,
                        linNum: linNum.value
                    },
                ]
            })
        }


        data.gameInfo.DefaultCs = data.gameInfo.DefaultCs || tableData.value[0].BetValue
        data.gameInfo.CrashRate = data.gameInfo.CrashRate || 1
        data.gameInfo.DefaultBetLevel = data.gameInfo.DefaultBetLevel || tableData.value[0].BetInfo[0].multipleValue

        data.gameInfo.totalBet = (Number(data.gameInfo.DefaultCs) * Number(data.gameInfo.DefaultBetLevel) * Number(linNum.value)).toFixed(2)


        return data
    },
    set(value) {


        addForm.value = value
    }
})
const addDialog = computed(() => {
    return Boolean(props.modelValue)
})


const addForm: Ref<MerchantGameInterface> = ref(<MerchantGameInterface>{
    GameId: "",
    Gametype: 0,
    GameName: "",
    ConfigPath: "",
    GameManufacturer: "",
    RTP: RTP_VALUE[0],
    StopLoss: 1,
    MaxMultipleOff: 0,
    MaxWinPointsOff: 0,
    MaxMultiple: 10000,
    MaxWinPoints: 1000000,
    BetBase: "",
    GamePattern: 1,
    Preset: 0,
    GameOn: 0,
    FreeGameOff: 0,          // 免费游戏开关
    ShowNameAndTimeOff: 0,   // 商户是否显示游戏名称
    ShowNameAndTime: 0,   // 游戏是否显示游戏名称
    ShowBagOff: 0,             // 商户是是否显示背包
    ShowBag: 0,             // 游戏是否显示背包
    GameDemo: "",             // 试玩链接
    BuyRTP: 95,             // 试玩链接
    DefaultCs:0,
    DefaultLevel:0,
    OnlineUpNum:0,
    OnlineDownNum:0,
    })

let commitData: Ref<updateGameConfigParams> = ref<updateGameConfigParams>({
    AppID: "",
    UserName: "",
    GameId: "",
    GameOn: 0,
    Preset: 0,
    StopLoss: 1,
    GamePattern: 0,
    FreeGameOff: 0,
    RTP: RTP_VALUE[0],
    MaxMultiple: 10000,
    ShowBag: 0,
    ShowNameAndTimeOff: 0,
    MaxWinPoints: 1000000,
    BetBase: "",
    DefaultCs: 0,
    DefaultBetLevel: 0,
})


const tableData = ref([])

const AddAdminer = () => {


    let gameSet = gameSetRef.value.commit
    let gameBet = gameBetRef.value.tableData
    let gameBetOBj = gameBetRef.value

    let BetBaseArray = []

    let BetBase = ""
    let OnlineUpNum = 0
    let OnlineDownNum = 0

    if (gameInfo.value.gameInfo.GameManufacturer != "SPRIBE"){
        if (gameInfo.value.gameInfo.GameManufacturer.toLowerCase() == "jili" || gameInfo.value.gameInfo.GameManufacturer.toLowerCase() == "tada") {
            BetBaseArray = gameBet.map(item => ((parseFloat(item.BetValue)* 1 * linNum.value)).toFixed(2))

        } else {

            BetBaseArray = gameBet.map(item => (parseFloat(item.BetValue)).toFixed(2))
        }


        if (gameInfo.value.gameInfo.GameManufacturer.toLowerCase() == "pp") {

            let PPBet = []
            for (let j = 1; j <= 10; j++) {
                for (let i = 0; i < BetBaseArray.length; i++) {


                    PPBet.push((parseFloat(BetBaseArray[i]) * j).toFixed(2))
                }
            }
            BetBaseArray = [...PPBet]

        }
        // BetBaseArray = [...new Set(BetBaseArray)]
        BetBaseArray = BetBaseArray.sort((a, b) => parseFloat(a) - parseFloat(b))

        if (gameInfo.value.gameInfo.GameManufacturer.toLowerCase() == "pg" && BetBaseArray.indexOf(parseFloat(gameInfo.value.gameInfo.DefaultCs).toFixed(2)) == -1) {
            return tip.e("默认押注金额和当前投注额不符")
        }
        BetBase = BetBaseArray.join(",")
        OnlineUpNum =  gameBetOBj.onlineUpNum
        OnlineDownNum = gameBetOBj.onlineDownNum


    }else{
        let SPRIBE_AircraftData = gameBetRef.value.SPRIBE_AircraftData
        BetBase = SPRIBE_AircraftData.BetBase.join(",")
        OnlineUpNum = gameBetOBj.SPRIBE_AircraftData.OnlineUpNum
        OnlineDownNum = gameBetOBj.SPRIBE_AircraftData.OnlineDownNum

    }


    gameInfo.value.Preset = isNaN(gameInfo.value.gameInfo.Preset) ? 0 : Number(gameInfo.value.gameInfo.Preset)


    commitData.value = {
        GameOn: gameInfo.value.gameInfo.GameOn,
        AppID: gameInfo.value.merchantInfo.AppID,                                                             // 商户名称
        GameId: gameInfo.value.gameInfo.GameId,                                                               // 游戏编号
        Preset: Number(gameInfo.value.gameInfo.Preset),                                                       // 预设面额
        StopLoss: Number(gameSet.StopLoss),                                                                   // 止盈止损开关
        BuyRTP: Number(gameSet.BuyRTP), // Bet上线下限
        GamePattern: gameSet.GamePattern,                                                                     // 游戏类型
        FreeGameOff: Number(gameSet.FreeGameOff),                                                             // 免费游戏开关
        RTP: gameSet.RTP,                                                                                     // RTP设置
        MaxMultiple: gameSet.MaxMultiple.toString(),                                                          // 赢取最高押注倍数
        MaxWinPoints: Number(gameSet.MaxWinPoints),                                                           // 赢取最高押注倍数
        ShowBag: Number(gameSet.ShowBag),                                                                     // 赢取最高押注倍数
        ShowNameAndTimeOff: Number(gameSet.ShowNameAndTimeOff),                                                // 赢取最高押注倍数
        DefaultCs: Number(gameInfo.value.gameInfo.DefaultCs),                                              // 赢取最高押注倍数z
        DefaultBetLevel: Number(gameInfo.value.gameInfo.DefaultBetLevel) ,                                   // 赢取最高押注倍数
        OnlineUpNum: Number(OnlineUpNum),                                                        // Bet上线下限
        OnlineDownNum: Number(OnlineDownNum),                                                    // Bet上线下限
        Scale: Number(gameBetOBj.SPRIBE_AircraftData.Scale),                                                    // Bet上线下限
        CrashRate:Number(gameSet.CrashRate),
        ProfitMargin:Number(gameSet.ProfitMargin),
        // Bet上线下限
        // BetMult: Number(gameInfo.value.gameInfo.BetMult),                                                    // Bet上线下限
        BetBase                                                                                             // 游戏投注
    }

    if (commitData.value.GameId == "spribe_01"){

            if(commitData.value.OnlineDownNum < commitData.value.OnlineUpNum) {
                return
            }

        commitData.value["DefaultCs"] = gameBetOBj.SPRIBE_AircraftData.DefaultCs
    }


    ElMessageBox.confirm(
        t('确认修改游戏设置'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            Client.Do(Merchant.UpdateGameConfig, commitData.value).then(([data, err]) => {
                if (err) {
                    tip.e(t(err))
                    return
                }
                tip.s(t("设置成功"))
                emits("update:modelValue", false)
            })
        })
}


</script>

<style scoped>


.switchContainer {
    width: 140px !important;
    display: flex;
    align-items: center;
    justify-content: flex-end;

}
</style>

<style lang="scss">
.dialog__form {

  .el-form-item {


    width: 100%;
  }

  .el-select {
    width: 100%;

  }

  .el-form-item__content > * {
    width: 80%;

  }

}

.gameType__container {
  width: 100%;

}

.gameType__item {
  width: auto;
  height: auto;
  padding: 8px 10px;
  border: 1px solid #cccccc59;
}
</style>
