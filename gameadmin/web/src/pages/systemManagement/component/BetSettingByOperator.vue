<template>

  <!-- 添加弹框 -->
    <el-dialog v-model="props.modelValue" :title="$t('游戏设置')"

               destroy-on-close
               :width="store.viewModel === 2 ? '100%' : '650px'" @close="emits('update:modelValue') + ':'">
        <el-form ref="addFormRef" :model="gameInfo" label-width="120px" :inline="true" class="dialog__form">

            <gamelist_container
                    :hase-all="true"
                    :hase-manufacturer="true"
                    :visible-manufacturer="props.gameSettingInfoData.DefaultManufacturerOn"
                    :defaultGameEvent="defaultGameEvent"
                    @select-operator="selectGameList"/>




        </el-form>

        <gameBet v-if="gameBetVisible" ref="gameBetRef" :merchantInfo="merchantInfo" :gameInfo="merchantInfo" :tableData="merchantInfo.tableData" :propLinNum="merchantInfo.lineNum"/>

        <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('确认') }}</el-button>
                </span>
        </template>
    </el-dialog>


</template>


<script setup lang="ts">

import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {nextTick, Ref, ref} from "vue";
import {RTP_VALUE} from "@/lib/RTP_config";
import {Client} from "@/lib/client";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {BatchEditGameRTPParam, Merchant} from "@/api/adminpb/merchant";
import {tip} from "@/lib/tip";
import {ElLoading, ElMessageBox} from "element-plus";
import Gamelist_container from "./gamelist_container.vue";
import GameSet from "@/pages/systemManagement/component/gameSettingInfoComponent/gameSet.vue";
import GameBet from "@/pages/systemManagement/component/gameSettingInfoComponent/gameBet.vue";


let {t} = useI18n()

const store = useStore()
const props = defineProps(["modelValue", "gameSettingInfoData"])
const emits = defineEmits(['update:modelValue', 'update:gameSettingInfoData', "commit"])
const defaultGameEvent = ref({})

const gameInfo: Ref<BatchEditGameRTPParam> = ref(<BatchEditGameRTPParam>{
    GameList: [],
    AppID: "",
    RTP: Number(RTP_VALUE[0]),
    GamePattern: 1,
    MaxWinPoints: 1000000,
    MaxMultiple: 10000,
})
const gameBetVisible = ref(false)
const componentInit = ref(false)
const gameBetRef = ref(null)
const merchantInfo = ref({
    GameManufacturer:"",
    GameIds:[],
    OnlineUpNum:0.05,
    OnlineDownNum:6000,
    ChangeBetOff:1,
    tableData:[],
    lineNum:1
})

const selectGameList = (value) => {


    let old_GameManufacturer = merchantInfo.value.GameManufacturer

    merchantInfo.value.GameManufacturer = value.manufacturer
    merchantInfo.value.GameIds = value.gameData

    if (old_GameManufacturer == "" || old_GameManufacturer != value.manufacturer){
        gameBetVisible.value = false
        merchantInfo.value.tableData = []
        merchantInfo.value.ChangeBetOff = 1
        let BetBase = ["0.05","0.5","2.5","10"]

        let linNum = merchantInfo.value.lineNum

        if (value.manufacturer.toLowerCase() == "jili" || value.manufacturer.toLowerCase() == "tada") {
            BetBase = ["1","2","3","5","8","10","20","50","100","200","300","400","500","700","1000"]
        }

        if (value.manufacturer.toLowerCase() == "pp") {
            merchantInfo.value.ChangeBetOff = 0

        }
        for (const i in BetBase) {

            if (value.manufacturer.toLowerCase() == "jili" || value.manufacturer.toLowerCase() == "tada") {
                BetBase[i] = (parseFloat(BetBase[i]) / linNum).toFixed(2)
            }


            merchantInfo.value.tableData.push({
                BetValue: parseFloat(BetBase[i]),
                BetInfo: [
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 1,
                        result: parseFloat(BetBase[i]) * linNum * 1,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 2,
                        result: parseFloat(BetBase[i]) * linNum * 2,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 3,
                        result: parseFloat(BetBase[i]) * linNum * 3,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 5,
                        result: parseFloat(BetBase[i]) * linNum * 5,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 6,
                        result: parseFloat(BetBase[i]) * linNum * 6,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 7,
                        result: parseFloat(BetBase[i]) * linNum * 7,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 8,
                        result: parseFloat(BetBase[i]) * linNum * 8,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 9,
                        result: parseFloat(BetBase[i]) * linNum * 9,
                        linNum: linNum
                    },
                    {
                        BetValue: parseFloat(BetBase[i]),
                        multipleValue: 10,
                        result: parseFloat(BetBase[i]) * linNum * 10,
                        linNum: linNum
                    }
                ]
            })
        }

        setTimeout(()=>{

            gameBetVisible.value = true
        })
    }

}

const AddAdminer = async () => {
    ElMessageBox.confirm(
        t('确认对选中的游戏进行设置Bet'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            const loading = ElLoading.service({
                lock: true,
                text: 'Loading',
                background: 'rgba(0, 0, 0, 0.7)',
            })

            let gameBet = gameBetRef.value.tableData
            let gameBetOBj = gameBetRef.value

            let BetBaseArray = []
            if (merchantInfo.value.GameManufacturer.toLowerCase() == "jili" || merchantInfo.value.GameManufacturer.toLowerCase() == "tada") {
                BetBaseArray = gameBet.map(item => parseFloat(item.BetValue).toFixed(2))

            } else {

                BetBaseArray = gameBet.map(item => parseFloat(item.BetValue).toFixed(2))
            }


            if (merchantInfo.value.GameManufacturer.toLowerCase() == "pp") {

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


            let BetBase = BetBaseArray.join(",")

            const params = {
                AppID: props.gameSettingInfoData.AppID,
                Manufacturer: merchantInfo.value.GameManufacturer,
                GameIds: merchantInfo.value.GameIds.join(","),
                OnlineUpNum: Number(gameBetOBj.onlineUpNum),                         // Bet上线下限
                OnlineDownNum: Number(gameBetOBj.onlineDownNum),                     // Bet上线下限
                BetBase,
            }

            const [response, err] = await Client.Do(Merchant.BatchSetGameBet, params)
            loading.close()
            if (err) {

                tip.e(t("修改Bet失败"))
                return
            }
            tip.s(t("修改Bet成功"))

            emits("update:modelValue")
            emits("commit")
        })


}

</script>

<style scoped>


</style>
