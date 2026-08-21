<template>

  <!-- 添加弹框 -->
    <el-dialog v-model="addDialog" :title="$t('游戏设置')"
               destroy-on-close
               @open-auto-focus="openDialog"
               :width="store.viewModel === 2 ? '100%' : '950px'" @close="emits('update:modelValue')">
        <el-form ref="addFormRef" :model="gameSearch" label-width="120px" :inline="true" class="dialog__form">

            <el-form-item :label="$t('游戏分类') + ':'">
                <el-radio-group v-model="activeName">
                    <el-radio-button  size="large" value="SLOT">SLOT</el-radio-button>
                </el-radio-group>

            </el-form-item>
            <el-form-item :label="$t('游戏厂商') + ':'">
                <el-select
                    v-model="gameSearch.GameManufacturer"
                    :placeholder="$t('选择游戏状态')"
                    style="width: 150px;"
                    @change="changeForm"

                >
                    <el-option
                        v-for="(item, key) in ManufacturerList"
                        :key="key"
                        :label="$t(item.ManufacturerName)"
                        :value="item.ManufacturerCode"
                    />
                </el-select>

                <div style="display: flex;align-items: center">
                    <span class="el-form-item__label" style="margin-left: 15px;font-size: 14px">{{
                            $t("游戏名称")
                        }}:</span>
                    <el-input v-model.lazy="gameSearch.GameName" style="width: 150px" @input="changeForm"
                              :placeholder="$t('请输入游戏名称')" clearable/>

                    <span class="el-form-item__label" style="margin-left: 15px;font-size: 14px">{{
                            $t("开启状态")
                        }}:</span>
                    <el-select
                        v-model="gameSearch.GameOn"
                        :placeholder="$t('选择游戏状态')"
                        style="width: 150px;"
                        @change="changeForm"
                    >
                        <el-option
                            v-for="(item, key) in gameState"
                            :key="key"
                            :label="$t(item.label)"
                            :value="item.value"
                        />
                    </el-select>


                </div>

            </el-form-item>
        </el-form>


        <div v-loading="tableDataLoading">
        <customTable
            table-name="operatorMaintenance_game_list"
            :table-header="tableHeader"
            :table-data="tableData"
            :page="gameSearch.PageIndex"
            :page-size="gameSearch.PageSize"
            @refresh-table="getGameList"
            :count="Count"
            height="400px"
            @selection-change="getSelectGame"
            @page-change="pageChange">

            <template #handleTools>

                <el-button class="handleBtn" type="primary" plain @click="batchGameStatus(0)">{{ $t("批量开启") }}</el-button>
                <el-button class="handleBtn" type="danger" plain @click="batchGameStatus(1)">{{ $t("批量关闭") }}</el-button>
                <el-button class="handleBtn" type="warning" plain @click="settingDialog = true" :disabled="!((store.AdminInfo.Businesses.RTPOff == 0 && store.AdminInfo.GroupId <= 1) ||  store.AdminInfo.Businesses.RTPOff == 1)">{{ $t("批量修改RTP") }}</el-button>
                <el-button class="handleBtn" type="primary" plain @click="setGameSetting = true">{{ $t("Slots游戏设置") }}</el-button>
                <el-button class="handleBtn" type="primary" plain @click="setBetDialog = true" style="margin-right: 30px">{{ $t("批量修改Bet") }}</el-button>
            </template>

            <template #StopLoss="scope">
                <el-text type="success" style="cursor: auto;" v-if="scope.scope.StopLoss == 0">{{ $t('开启') }}</el-text>
                <el-text type="danger"  style="cursor: auto;" v-if="scope.scope.StopLoss == 1">{{ $t('关闭') }}</el-text>

            </template>
            <template #operator="scope">
                <el-button type="primary" size="small" plain @click="openGameSetting(scope.scope)">{{ $t('设置') }}</el-button>
                <el-button type="success"  v-if="scope.scope.GameOn == 1" size="small" plain @click="editGameStatus(scope.scope)">
                    {{ $t('开启') }}</el-button>
                <el-button type="danger"  v-if="scope.scope.GameOn == 0" size="small" plain @click="editGameStatus(scope.scope)">
                    {{ $t('关闭') }}</el-button>
<!--                <el-button type="warning" size="small" plain>试玩</el-button>-->

            </template>

        </customTable>
        </div>
        <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('关闭') }}</el-button>
                </span>
        </template>

    </el-dialog>


    <RTPSettingByOperator v-model="settingDialog" :gameSettingInfoData="props.gameSettingData" @commit="getGameList()"/>
<!--    <BetSettingByOperator v-model="setBetDialog" :gameSettingInfoData="props.gameSettingData" @commit="getGameList()"/>-->
    <BetMultiplesSettingByOperator v-model="setBetDialog" :gameSettingInfoData="props.gameSettingData" @commit="getGameList()"/>
    <GameSettingByOperator v-model="setGameSetting" :gameSettingInfoData="props.gameSettingData" @commit="getGameList()"/>

</template>


<script setup lang="ts">

import {computed, Ref, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import customTable from "@/components/customTable/tableComponent.vue"
import {Merchant, MerchantGameReq, updateGameConfigParams} from "@/api/adminpb/merchant";
import {Client} from "@/lib/client";
import {tip} from "@/lib/tip";
import {ElMessageBox} from "element-plus";
import {Throttle} from "@/lib/util";
import RTPSettingByOperator from "@/pages/systemManagement/component/RTPSettingByOperator.vue";
import BetMultiplesSettingByOperator from "@/pages/systemManagement/component/BetMultiplesSettingByOperator.vue";
import GameSettingByOperator from "@/pages/systemManagement/component/GameSettingByOperator.vue";
import {Manufacturer} from "@/api/gamepb/game";
import {debug} from "util";



let {t} = useI18n()

const store = useStore()

const props = defineProps(["modelValue", "gameSettingData"])
const emits = defineEmits(['update:modelValue', 'update:gameSettingData', "openSettingInfo"])

const settingDialog = ref(false)
const setBetDialog = ref(false)
const setGameSetting = ref(false)

const addDialog = computed(() => {


    return Boolean(props.modelValue)
})


const tableDataLoading = ref(true)
const activeName = ref("SLOT")
const Count = ref(0)
const selectGameId = ref([])

const tableHeader = [
    {width:"40px" ,label: "", value: "", type: "selection", align:"left", fixed:"left", hiddenVisible: true},
    { label: "游戏编号", value: "GameId", width: "180px"},
    { label: "游戏名称", value: "GameName", width: "180px"},
    { label: "RTP", value: "RTP", width: "80px", format:(row)=>row.RTP + "%"},
    { label: "止盈止损开关", value: "StopLoss", type: "custom", width: "140px"},
    { label: "赢取最高倍数", value: "MaxMultiple", width: "150px"},
    { label: "游戏投注", value: "BetBase",width: "180px", format:(row) => row.GameManufacturer.toLowerCase() != "pp" ? row.BetBase : `${row.OnlineUpNum},${row.OnlineDownNum}`},
    { label: "操作", value: "operator", type: "custom", fixed:"right",  width: "250px", hiddenVisible: true},

]
const tableData = ref([

])


const gameSearch:Ref<MerchantGameReq> = ref<MerchantGameReq>({
    Gametype    : 0,
    GameName    : "",
    GameOn      : -1,
    UserName    : "",
    GameManufacturer:props.gameSettingData.DefaultManufacturerOn,
    PageIndex   : 1,
    PageSize    : 20,
})

const openDialog = () => {
    tableData.value = []
    gameSearch.value = {
        Gametype    : 0,
        GameManufacturer: props.gameSettingData.DefaultManufacturerOn,
        GameName    : "",
        GameOn      : -1,
        UserName    : "",
        PageIndex   : 1,
        PageSize    : 20,
    }
    selectGameId.value = []
    getGameList()
    getManufacturerList()
}


const ManufacturerList = ref([])
const gameState = ref([
    {label:t("全部"), value: -1},
    {label:t("开启"), value: 0},
    {label:t("关闭"), value: 1},
])


const getSelectGame = (data) => {
    selectGameId.value = data.map(item=> item.GameId)
}
const batchGameStatus = (status) => {

    if (selectGameId.value.length < 1){
        tip.e(t('请选择游戏'))
        return
    }
    let queryParam = {
        AppID:props.gameSettingData.AppID,
        GameIds:selectGameId.value,
        GameOn:status,
    }
    ElMessageBox.confirm(
        t(`确定批量${status ?'关闭' : '开启' }，确认切换吗?`),
        t('提示'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('取消'),
            type: 'warning',
        }
    )
        .then(async () => {
            tableDataLoading.value = true
            const [data, err] = await Client.Do(Merchant.BatchSetGameOn, queryParam as any)

            tableDataLoading.value = false

            if (err){
                tip.e(t('修改失败'))
                return
            }
            tip.s(t('修改完成'))

            getGameList()
        })

}
const changeForm = () => {


    Throttle(() => {
       getGameList()
    }, 800)()


}

const getManufacturerList = async () => {

    if (props.gameSettingData.DefaultManufacturerOn){

        let data = props.gameSettingData.DefaultManufacturerOn.map(item=>({
            Id:"",
            ManufacturerName: item,
            ManufacturerCode: item,
        }))
        data.unshift({
            Id:"",
            ManufacturerName:"全部",
            ManufacturerCode:"ALL",
        })
        ManufacturerList.value = data
        gameSearch.value.GameManufacturer = "ALL"

    }else{
        const [ response, err ] = await Client.Do(Manufacturer.GetManufacturerList, {} as any)

        if (err){
            ManufacturerList.value = []
        }else{
            response.List.unshift({
                Id:"",
                ManufacturerName:"全部",
                ManufacturerCode:"ALL",
            })
            ManufacturerList.value = response.List
            gameSearch.value.GameManufacturer = "ALL"
        }
    }

}
const getGameList = async () => {
    tableDataLoading.value = true
    tableData.value = []
    let queryParam = {
        ...gameSearch.value
    }
    // 商户的Id
    queryParam.UserName = props.gameSettingData.AppID

    if (queryParam.GameManufacturer == null || queryParam.GameManufacturer == undefined){
        queryParam.GameManufacturer = props.gameSettingData.DefaultManufacturerOn
    }
     if (queryParam.GameManufacturer && queryParam.GameManufacturer.indexOf("ALL") != -1){
        queryParam.GameManufacturer = props.gameSettingData.DefaultManufacturerOn
    }

     if (typeof queryParam.GameManufacturer == 'string'){
         queryParam.GameManufacturer =  queryParam.GameManufacturer.split(",")
     }




    const [data, err] = await Client.Do(Merchant.GetMerchantGames, queryParam)
    tableDataLoading.value = false
    tableData.value = data.List
    Count.value = data.CountAll
}
const pageSizeChange = (page) => {

    getGameList()
}
const pageChange = (page) => {
    gameSearch.value.PageSize = page.dataSize
    gameSearch.value.PageIndex = page.currentPage
    getGameList()
}

const editGameStatus = (row) => {
    ElMessageBox.confirm(
        t("确认调整当前游戏运行状态"),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(() => {
            let data: updateGameConfigParams = <updateGameConfigParams>{
                ...row,
                AppID                   : props.gameSettingData.AppID,   // 商户名称
                GameId                  : row.GameId, // 商户名称
                GameOn                  : row.GameOn == 0 ? 1 : 0,
                Gametype                : row.Gametype,// 商户名称
                GameName                : row.GameName,// 游戏编号
                GameManufacturer        : row.GameManufacturer,// 预设面额
                ConfigPath              : row.ConfigPath,// 止盈止损开关
                RTP                     : row.RTP,// 游戏类型
                StopLoss                : Number(row.StopLoss),// 免费游戏开关
                MaxMultipleOff          : row.MaxMultipleOff,// RTP设置
                MaxMultiple             : row.MaxMultiple.toString(),// 赢取最高押注倍数
                BetBase                 : row.BetBase,// 游戏投注
                GamePattern             : row.GamePattern,// 游戏投注
                Preset                  : row.Preset,// 游戏投注
                FreeGameOff             : row.FreeGameOff,
                DefaultCs               : row.DefaultCs,                      // 游戏投注
                DefaultLevel            : row.DefaultLevel,                      // 游戏投注
                ShowBag                 : row.ShowBag,
                ShowNameAndTimeOff      : row.ShowNameAndTimeOff,
                OnlineUpNum             : row.OnlineUpNum,
                OnlineDownNum           : row.OnlineDownNum,
                MaxWinPoints            : row.MaxWinPoints
            }


            Client.Do(Merchant.UpdateGameConfig, data).then(([data, err]) => {
                if (err){
                    tip.e(t("修改失败"))

                }
                tip.s(t("修改成功"))
                getGameList()

            })
        })



}

const openGameSetting = (data) => {
    data["userName"] = props.gameSettingData.UserName

    let settingGames = {
        gameInfo: {...data},
        merchantInfo: {
            ...props.gameSettingData
        }
    }


    emits("openSettingInfo", settingGames)
}

const AddAdminer = () => {
    emits('update:modelValue')
}
defineExpose({
    getGameList
})
</script>


<style lang="scss">
.dialog__form {
    .el-form-item{


        width: 100%;
    }

    .el-select{
        width: 100%;

    }
    .el-form-item__content > *{
        width: 80%;

    }

}
.gameType__container{
  width: 100%;

}
.gameType__item{
  width: auto;
  height: auto;
  padding: 8px 10px;
  border: 1px solid #cccccc59;
}
.handleBtn{
    margin-top: 5px;
    margin-bottom: 5px;
}
</style>
