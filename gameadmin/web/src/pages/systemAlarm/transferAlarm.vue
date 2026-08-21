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
                    table-name="transferRateAlarm_list"
                    :table-data="tableData"
                    :table-header="tableHeader"
                    :page="alarmSearchForm.Page"
                    :page-size="alarmSearchForm.PageSize"
                    :count='Count'
                    @refresh-table="queryList"
                    @page-change="changePage">

                <template #handleTools>
                    <el-button type="primary"  @click="checkAllRead" plain>{{ $t('全部已读') }}</el-button>
                    <el-button type="warning" plain @click="RickRuleDialog = true">{{ $t('设置规则') }}</el-button>
                </template>

                <template #OrderId="scope">
                    {{ scope.scope.OrderId }}
                </template>

                <template #ProfitAndLoss="scope">

                    <el-text :type="scope.scope.ProfitAndLoss[0] == '-' ? 'danger' : 'success'" size="small">
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


        <rickControlDialog v-model="RickRuleDialog" title="转账预警" rickType="transfer"></rickControlDialog>

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

import rickControlDialog from "@/pages/systemAlarm/component/riskControlDialog.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import ut from "@/lib/util";
import {tip} from "@/lib/tip";
import Operator_container from "@/components/operator_container.vue";

const {t} = useI18n()

const loading = ref(false)
const Count = ref(0)
const timeRange = ref([])
const alarmSearchForm: Ref<AlarmRequest> = ref(<AlarmRequest>{
    Pid: null,
    GameId: "",
    AppID: "",
    AlarmTimeStart: 0,
    AlarmTimeEnd: 0,
    Page: 1,
    PageSize: 20,
    Type: "transfer",
})
const defaultGameEvent = ref({})
const defaultOperatorEvent = ref({})
const RickRuleDialog = ref(false)
const tableHeader = [
    {label: "交易订单号", value: "OrderId", width: "250px", type: "custom"},
    {label: "唯一标识", value: "Pid"},
    {label: "商户AppID", value: "MerchantId"},
    {label: "转出金额", value: "TotalMoney"},
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

    if (value) {
        alarmSearchForm.value.GameId = value
    } else {
        alarmSearchForm.value.GameId = ""
    }
}


const init = async () => {

    loading.value = true


    let startTime = 0
    let endTime = 0


    if (timeRange.value) {
    if (timeRange.value[0] && timeRange.value[1]) {

        startTime = ut.fmtUTCDate(timeRange.value[0]) * 1000
        endTime = ut.fmtUTCDate(timeRange.value[1]) * 1000
    }
    }


    const queryData = {
        ...alarmSearchForm.value,
        Pid: alarmSearchForm.value.Pid ? alarmSearchForm.value.Pid : null
    }

    if (startTime) {

        queryData.AlarmTimeStart = ut.fmtSelectedUTCDate(startTime, "reduce")
    }
    if (endTime) {

        queryData.AlarmTimeEnd = ut.fmtSelectedUTCDate(endTime, "reduce")
    }


    let [response, err] = await Client.Do(Alarm.GetAlarmHistory, queryData)
    loading.value = false
    if (!err) {
        Count.value = response.Count
    }
    let responseData = []
    response.List.forEach(item => {

        responseData.push({
            ...item,
            OrderId: item.OrderId,
            TotalMoney: item.Amount,
            MerchantId: item.AppId,

        })


    })


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

const readAlert = async (data: EditRequest) => {
    data.Type = alarmSearchForm.value.Type
    const [response, err] = await Client.Do(Alarm.SetAlarmHistory, data)

    if (err) {
        return tip.e(t(err))
    }

    init()

}

init()
</script>

<style scoped lang="scss">
.el-link .el-icon--right.el-icon {
  vertical-align: text-bottom;
}
</style>
