<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                    :model="alarmSearchForm"
                    style="max-width: 100%"
            >
                <el-space wrap>


                    <operator_container :hase-all="true" :is-init="true" :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange"/>
                    <el-form-item :label="$t('游戏名称')">
                        <el-select v-model="alarmSearchForm.GameId" @clear="clearGame" filterable clearable style="width: 150px">
                            <el-option
                                    v-for="item in games"
                                    :key="item.value"
                                    :label="$t(item.label)"
                                    :value="item.value"
                            />
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="alarmSearchForm.Pid" clearable
                                  oninput="value=value.replace(/[^\d]/g,'')"
                                  :placeholder="$t('请输入')"/>
                    </el-form-item>
                    <el-form-item :label="$t('触发时间')">
                        <el-date-picker
                                v-model="timeRange"
                                type="datetimerange"
                                value-format="x"
                                format="YYYY-MM-DD HH:mm:ss"
                                :range-separator="$t('至')"
                                :start-placeholder="$t('开始时间')"
                                :end-placeholder="$t('结束时间')"
                        />
                    </el-form-item>


                </el-space>
            </el-form>

            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>


        <div class="page_table_context">


            <customTable
                v-loading="loading"
                table-name="returnRateAlarm_list"
                :table-data="tableData"
                :table-header="tableHeader"
                :page="alarmSearchForm.Page"
                :page-size="alarmSearchForm.PageSize"
                :count='Count'
                @refresh-table="queryList"
                @page-change="changePage">

                <template #handleTools>
                    <el-button type="primary" plain @click="checkAllRead">{{ $t('全部已读') }}</el-button>
                    <el-button type="warning" plain @click="RickRuleDialog = true">{{ $t('设置规则') }}</el-button>
                </template>


                <template #OrderId="scope">

                    {{ scope.scope.OrderId }}
                </template>

                <template #ProfitAndLoss="scope">

                    <el-text :type="scope.scope.ProfitAndLoss.indexOf('-') != -1 ? 'danger' : 'success'" size="small">
                        {{ scope.scope.ProfitAndLoss }}
                    </el-text>
                </template>
                <template #Operator="scope">

                    <el-button plain :type="scope.scope.ReadStatus == 1 ? 'primary' : 'warning'" size="small"
                               @click="checkCurrent(scope.scope)">{{ scope.scope.ReadStatus == 1 ? $t("已读") : $t("未读") }}
                    </el-button>

                </template>
            </customTable>

        </div>

        <rickControlDialog v-model="RickRuleDialog" title="回报率预警" rickType="returnRate"></rickControlDialog>

    </div>
</template>
<script setup lang="ts">

import type {Ref} from "vue";
import type {AlarmRequest, AlarmListResponse, AlarmData, EditRequest} from "@/api/adminpb/alarm";
import {ref} from "vue";
import {useI18n} from "vue-i18n";

import {Alarm, AlarmItem} from "@/api/adminpb/alarm";
import {Client} from "@/lib/client";
import Gamelist_container from "@/components/gamelist_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import rickControlDialog from "@/pages/systemAlarm/component/riskControlDialog.vue";
import ut from "@/lib/util";
import {tip} from "@/lib/tip";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {ElMessageBox} from "element-plus";
import Operator_container from "@/components/operator_container.vue";
import {useStore} from "@/pinia";
import {storeToRefs} from "pinia";

const {t} = useI18n()
const store = useStore();
const { AdminInfo } = storeToRefs(store)
const loading = ref(false)
const Count = ref(0)
const timeRange = ref([])

const RickRuleDialog = ref(false)

const alarmSearchForm: Ref<AlarmRequest> = ref(<AlarmRequest>{
    Pid: null,
    GameId: "",
    AppId: "",
    AlarmTimeStart: 0,
    Manufacturer: "",
    AlarmTimeEnd: 0,
    Page: 1,
    PageSize: 20,
    Type: "turnRate"
})

// Client.send("/mq/AdminInfo/Interior/PushRickRuleTransferOut", {AppID:"qwe456"})
// Client.send("/mq/AdminInfo/Interior/PushRickRuleReturnRate", {AppID:"qwe456"})
const defaultGameEvent = ref({})
const defaultOperatorEvent = ref({})
const games = ref([])
let DefaultManufacturerOn =AdminInfo.value.Businesses.DefaultManufacturerOn

const tableHeader = [
    {label: "交易订单号", value: "OrderId", width: "250px", type: "custom"},
    {label: "唯一标识", value: "Pid"},
    {label: "商户AppID", value: "MerchantId"},
    {label: "币种", value: "Currency"},
    {label: "总回报率", value: "TotalRate"},
    {label: "总投注", value: "BetMount", width: "150px"},
    {label: "总赢分", value: "ProfitAndLoss", type: "custom", width: "150px"},
    {label: "日志描述", value: "LogInfo", width: "220px"},
    {
        label: "触发时间",
        value: "createTime",
        format: (row) => row.CreateTime == 0 ? '' : `${ut.fmtSelectedUTCDateFormat(row.CreateTime)}`,
        width: "220px"
    },
    {label: "操作", value: "Operator", type: "custom", fixed: "right", width: "100px", hiddenVisible: true},

]
const tableData: Ref<AlarmItem[]> = ref(<AlarmItem[]>[])


const selectGameList = (value, lotteryValue) => {

    if (value.gameData) {
        alarmSearchForm.value.GameId = value.gameData
    } else {
        alarmSearchForm.value.GameId = ""
    }
    if (value.manufacturer || value.manufacturer == null){
        alarmSearchForm.value.Manufacturer = value.manufacturer
    } else {
        alarmSearchForm.value.Manufacturer = ""
    }
}


const init = async () => {

    loading.value = true

    let startTime = 0
    let endTime = 0

    if (timeRange.value){
        if (timeRange.value[0] && timeRange.value[1]) {

            startTime = ut.fmtUTCDate(timeRange.value[0])
            endTime = ut.fmtUTCDate(timeRange.value[1])
        }

    }


    const queryData = {
        ...alarmSearchForm.value,
        Pid: alarmSearchForm.value.Pid ? alarmSearchForm.value.Pid : null
    }

    if (startTime) {

        queryData.AlarmTimeStart = ut.fmtSelectedUTCDate(startTime * 1000, "reduce")
    }
    if (endTime) {

        queryData.AlarmTimeEnd = ut.fmtSelectedUTCDate(endTime * 1000, "reduce")
    }


    let [response, err] = await Client.Do(Alarm.GetAlarmHistory, queryData)
    loading.value = false

    if (!err) {
        Count.value = response.Count
    }
    let responseData = []
    if (response.List) {


        response.List.forEach(item => {


            responseData.push({
                ...item,
                OrderId: item.OrderId,
                ProfitAndLoss: item.TotalWinLoss,
                BetMount: item.TotalBet,
                TotalRate: item.WinRate,
                Currency: item.Currency,
                MerchantId: item.AppId,
                LogInfo: t("在 {GameId} 盈利 {WinMoney}", {GameId: item.GameId, WinMoney: item.WinMoney}),
            })


        })
    }

    tableData.value = responseData
}

const operatorListChange = (value) => {
    if (value){

        alarmSearchForm.value.AppId = value.AppID
    }else{

        alarmSearchForm.value.AppId = ""
    }
}


const changePage = (page) => {
    alarmSearchForm.value.Page = page.currentPage
    alarmSearchForm.value.PageSize = page.dataSize
    init()
}

// 查询报警日志
const queryList = () => {
    tableData.value = []
    init()
}


const checkCurrent = (row) => {
    if (row.ReadStatus) {
        return
    }
    readAlert(<EditRequest>{Id: row.Id})
}

const checkAllRead = () => {
    readAlert(<EditRequest>{Range: "all"})
}

const getGameList = async () => {
    let [data, err] = await Client.Do(AdminGameCenter.GameList, {})

    if (DefaultManufacturerOn){
        data.List = data.List.filter(item => DefaultManufacturerOn.indexOf(item.ManufacturerName) != -1)
    }


    data.List.unshift({
        ID: "ALL",
        Name: "全部",
        Status: "0",
    })

    alarmSearchForm.value.GameId = "ALL"


    games.value = data.List
        .map(item => ({
            value: item.ID,
            label: item.Name,
            Status: item.Status,
        }));


}

const readAlert = async (data: EditRequest) => {


    data["Type"] = alarmSearchForm.value.Type
    const [response, err] = await Client.Do(Alarm.SetAlarmHistory, data)

    if (err) {
        return tip.e(t(err))
    }

    init()

}

init()
getGameList()

const clearGame = () => {

    alarmSearchForm.value.GameId = "ALL"
}
</script>

<style scoped lang="scss">
.el-link .el-icon--right.el-icon {
  vertical-align: text-bottom;
}
</style>
