<template>

  <!-- 添加弹框 -->
    <el-dialog v-model="props.modelValue" :title="$t('游戏设置')"
               @open-auto-focus="openDialog"
               destroy-on-close
               :width="store.viewModel === 2 ? '100%' : '650px'" @close="emits('update:modelValue') + ':'">
        <el-form ref="addFormRef" :model="gameInfo" label-width="120px" :inline="true" class="dialog__form">

            <el-form-item :label="$t('游戏名称')">
                <el-select v-model="gameInfo.GameList" @change="selectGameChange"
                           style="width: 300px;" :placeholder="$t('请选择游戏')"
                           @clear="clearGame"
                           multiple filterable collapse-tags clearable>
                    <el-option v-for="(item, index) in games" :key="index" :label="$t(item.GameName)"
                               :value="item.GameId">

                        <div style="display: flex;align-items: center;justify-content: space-between">
                            <span style="float: left">{{ $t(item.GameName) }}</span>
                            <el-tag
                                    v-if="item.ID != 'ALL'"
                                    type="primary"
                                    size="small"
                                    style="float: right">

                                {{ item.GameManufacturer }}
                            </el-tag>
                        </div>
                    </el-option>
                </el-select>

            </el-form-item>
            <el-form-item :label="$t('RTP设置')">
                <el-select
                        v-model="gameInfo.RTP"
                        :placeholder="$t('请选择')"
                        style="width: 300px;"
                        prefix="%"
                >
                    <template v-for="item in RTP_VALUE">
                        <el-option
                            v-if="props.gameSettingInfoData.HighRTPOff == 0 ||
                                 (props.gameSettingInfoData.HighRTPOff > 0
                                 && (Number(item) >= OperatorRTP[props.gameSettingInfoData.HighRTPOff].split('-')[0]
                                 && Number(item) <= OperatorRTP[props.gameSettingInfoData.HighRTPOff].split('-')[1]))"
                            :key="Number(item)"
                            :label="Number(item) == 93 ? `93% ${$t('(常用)')}` : `${item}%`"
                            :value="Number(item)"
                        />
                    </template>

                </el-select>
            </el-form-item>
            <el-form-item :label="$t('购买免费游戏RTP')" v-if="((props.gameSettingInfoData.BuyRTPOff == 0 && store.AdminInfo.GroupId <= 1) || props.gameSettingInfoData.BuyRTPOff == 1)">
                <el-select
                        v-model="gameInfo.BuyRTP"
                        :placeholder="$t('请选择')"
                        style="width: 300px;"
                        prefix="%"
                >
                    <template v-for="item in BuyRTP">
                        <el-option
                            :label="item == 0 ? $t('无') : `${item}%`"
                            :value="Number(item)"
                        />
                    </template>

                </el-select>
            </el-form-item>

            <el-form-item>

                <template #label>
                    <el-tooltip
                        class="box-item"
                        effect="dark"

                        placement="top"
                    >

                        <template #content>
                            <div>{{ $t("心跳型：中奖率约5%-15%，中奖和小奖概率低，有概率出现超大奖，波动很大。")}}</div>
                            <div>{{ $t("波动型：中奖率约15%-20%，中奖和小奖概率低，有概率出现超大奖，波动大。")}}</div>
                            <div>{{ $t("仿正版（默认）：中奖率在20%-30%之间，与官方原版体验相同。")}}</div>
                            <div>{{ $t("混合型：中奖率约25%-35%，大奖概率略低，波动较小。")}}</div>
                            <div>{{ $t("稳定型：中奖率约35%-45%，大奖概率低，数值体验相对稳定。")}}</div>
                            <div>{{ $t("高中奖率：中奖率约45%-60%，大奖概率很低，中奖和小奖的中奖率高，数值稳定性高。")}}</div>
                            <div>{{ $t("超高中奖率：中奖率约40%-70%左右，大奖概率极低，中奖和小奖的中奖率较高，数值稳定性较高。")}}</div>
                        </template>
                        <el-icon><QuestionFilled /></el-icon>
                    </el-tooltip>



                    {{ $t('游戏类型') + ':' }}
                </template>
                <el-select
                        v-model="gameInfo.GamePattern"
                        :placeholder="$t('请选择')"
                        style="width: 300px;"
                >
                    <el-option
                            v-for="(item, index) in GameType"
                            :key="index"
                            :label="item.value == 5 ? $t(item.label) + '(' + $t('默认') + ')' : $t(item.label) "
                            :value="item.value"
                    />
                </el-select>
            </el-form-item>


                <el-form-item :label="$t('赢取最高钱数') + ':'" style="margin-bottom: 0">
                    <el-input v-model.number.trim="gameInfo.MaxWinPoints"
                              style="width: 300px;"
                              :disabled="!((props.gameSettingInfoData.MaxWinPointsOff == 0 && store.AdminInfo.GroupId <= 1) || props.gameSettingInfoData.MaxWinPointsOff == 1)"
                              :placeholder="$t('请输入')" @input="MaxWinPointsInput($event)">
                    </el-input>


                    <el-text type="danger" style="text-align: left;width: 300px;display: block">
                        {{ $t('赢取最高钱数可设置区间{Num}', {Num:"1-1,000,000"})}}
                    </el-text>

                </el-form-item>



            <el-form-item :label="$t('赢取最高倍数') + ':'" style="margin-bottom: 0">
                <!--                      :disabled="gameInfo.merchantInfo.MaxMultipleOff || store.AdminInfo.GroupId <= 1"-->
                <el-input v-model.number.trim="gameInfo.MaxMultiple"
                          style="width: 300px;"
                          :disabled="!((props.gameSettingInfoData.MaxMultipleOff == 0 && store.AdminInfo.GroupId <= 1) || props.gameSettingInfoData.MaxMultipleOff == 1)"
                          :placeholder="$t('请输入')" @blur="maxnultiple($event)">
                </el-input>

                <el-text type="danger" style="text-align: left;width: 300px;display: block">
                    {{ $t('赢取最高倍数可设置{Num}',{Num:"30~10000"}) }}
                </el-text>


            </el-form-item>

        </el-form>

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
import {Ref, ref} from "vue";
import {RTP_VALUE} from "@/lib/RTP_config";
import {Client} from "@/lib/client";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {BatchEditGameRTPParam, Merchant} from "@/api/adminpb/merchant";
import {tip} from "@/lib/tip";
import {ElLoading, ElMessageBox} from "element-plus";
import {OperatorRTP, GameType} from "@/enum";
import {QuestionFilled} from "@element-plus/icons-vue";


let {t} = useI18n()

const store = useStore()
const props = defineProps(["modelValue", "gameSettingInfoData"])
const emits = defineEmits(['update:modelValue', 'update:gameSettingInfoData', "commit"])
const defaultGameEvent = ref({})
const batDefaultRTP = ref({})

const gameInfo:Ref<BatchEditGameRTPParam> = ref(<BatchEditGameRTPParam>{
    GameList : [],
    AppID:"",
    RTP:93,
    BuyRTP:90,
    GamePattern:5,
    MaxWinPoints:1000000,
    MaxMultiple:100,
})

const games = ref([])

const BuyRTP = import.meta.env.VITE_BUY_RTP.split(",")
const initGameList = async () => {
    const [gameList, err] = await Client.Do(Merchant.GetMerchantGames, {
        UserName: props.gameSettingInfoData.AppID,
        GameManufacturer: props.gameSettingInfoData.DefaultManufacturerOn,
        PageIndex:1,
        PageSize:1000
    } as any)


    games.value = gameList.List.sort((a,b)=> a.Id - b.Id)


        games.value.unshift({
            GameName: "全部",
            ID: "ALL",
            GameId: "ALL"
        })
    gameInfo.value.GameList = ["ALL"]

}


const openDialog = () => {
    let localBet = localStorage.getItem("batDefaultRTP") || "{}"
    batDefaultRTP.value = JSON.parse(localBet)

    let APPConfig = {
        RTP: 93,
        BuyRTP: 90,
        GamePattern: 5,
        MaxWinPoints: 1000000,
        MaxMultiple: 100,
    }
    if(!batDefaultRTP.value[props.gameSettingInfoData.AppID]){
        batDefaultRTP.value[props.gameSettingInfoData.AppID] = APPConfig
    }


    gameInfo.value.GameList = []
    gameInfo.value.AppID = props.gameSettingInfoData.AppID
    gameInfo.value.RTP = batDefaultRTP.value[props.gameSettingInfoData.AppID]["RTP"]
    gameInfo.value.BuyRTP = batDefaultRTP.value[props.gameSettingInfoData.AppID]["BuyRTP"]
    gameInfo.value.GamePattern = batDefaultRTP.value[props.gameSettingInfoData.AppID]["GamePattern"]
    gameInfo.value.MaxWinPoints = batDefaultRTP.value[props.gameSettingInfoData.AppID]["MaxWinPoints"]
    gameInfo.value.MaxMultiple = batDefaultRTP.value[props.gameSettingInfoData.AppID]["MaxMultiple"]
    initGameList()

}


const MaxWinPointsInput = (value) => {
    gameInfo.value.MaxWinPoints = Number(gameInfo.value.MaxWinPoints.toString().replace(/^(0+)|[^\d]+/g,''))

    if (value < 1) {
        gameInfo.value.MaxWinPoints = 1
    }
    if (value > 1000000) {
        gameInfo.value.MaxWinPoints = 1000000
    }

}


const maxnultiple = (value) => {
    value = gameInfo.value.MaxMultiple
    if (value < 30) {
        gameInfo.value.MaxMultiple = 30
    }
    if (value > 10000) {
        gameInfo.value.MaxMultiple = 10000
    }

}


const selectGameChange = async (value) => {
    gameInfo.value.RTP = 93
    gameInfo.value.BuyRTP = 90
    gameInfo.value.GamePattern = 5
    gameInfo.value.MaxWinPoints = 1000000
    gameInfo.value.MaxMultiple = 100
    if (value.length >= 1){

        let index = value.findIndex(item => item == "ALL")

        if (index > 0){
            value = ["ALL"]
        }else {

            value = value.filter(item => item != 'ALL')

            if (value.length == 1){
                let param = {
                    "Gametype": 0,
                    "GameName": value[0],
                    "GameOn": 0,
                    "UserName": props.gameSettingInfoData.AppID,
                    "PageIndex": 1,
                    "PageSize": 20
                }
                Client.Do(Merchant.GetMerchantGames, param as any).then(([res, err])=>{

                    gameInfo.value.RTP = res.List[0].RTP
                    gameInfo.value.BuyRTP = res.List[0].BuyRTP
                    gameInfo.value.GamePattern = res.List[0].GamePattern
                    gameInfo.value.MaxWinPoints = res.List[0].MaxWinPoints
                    gameInfo.value.MaxMultiple = res.List[0].MaxMultiple
                })

            }

        }

    }else{
        value = ["ALL"]
    }

    gameInfo.value.GameList = value
}

const AddAdminer = async () => {
    ElMessageBox.confirm(
        t('确认对选中的游戏进行设置RTP'),
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
            const params = {...gameInfo.value}
            params.GameList = params.GameList.join(",")


            const [response, err] = await Client.Do(Merchant.BatchEditGameRTP, params)
            loading.close()
            if (err) {

                tip.e(t(err))
                return
            }
            tip.s(t("修改RTP成功"))
            saveRTPToLocalStorage()
            emits("update:modelValue")
            emits("commit")
        })


}
const clearGame = () => {
    gameInfo.value.GameList = ["ALL"]
}

const saveRTPToLocalStorage = () => {


    let data = batDefaultRTP.value || {}
    data[gameInfo.value.AppID] = {

        RTP: gameInfo.value.RTP,
        BuyRTP: gameInfo.value.BuyRTP,
        GamePattern: gameInfo.value.GamePattern,
        MaxWinPoints: gameInfo.value.MaxWinPoints,
        MaxMultiple: gameInfo.value.MaxMultiple,
    }


    localStorage.setItem("batDefaultRTP", JSON.stringify(data))
}

</script>

<style scoped>


</style>
