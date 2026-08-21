<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                    :model="SystemLogListParam"
                    style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container ref="operatorContainer2Ref" :hase-all="true" :is-init="true" :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange2"/>


                    <el-form-item :label="$t('创建时间')">
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
                <el-button type="default" @click="resetSearch">{{ $t('重置') }}</el-button>
            </el-space>
        </div>


        <div class="page_table_context">


            <customTable
                    v-loading="loading"
                    table-name="operatorApproval_list"
                    :table-data="tableData"
                    :table-header="tableHeader"
                    :page="SystemLogListParam.Page"
                    :page-size="SystemLogListParam.PageSize"
                    :count='Count'
                    @refresh-table="queryList"
                    @page-change="changePage">

                <template #OrderId="scope">
                    {{ scope.scope.OrderId }}
                </template>

                <template #OperatorType="scope">

                    <el-text :type="scope.scope.OperatorType == 2 ? 'danger' : 'success'" size="small">
                        {{ scope.scope.OperatorType == 2 ? $t('扣除') : $t('充值')}}
                    </el-text>
                </template>

            </customTable>

        </div>



    </div>
</template>
<script setup lang="ts">

import type {Ref} from "vue";
import {ref} from "vue";
import {useI18n} from "vue-i18n";

import {Client} from "@/lib/client";
import customTable from "@/components/customTable/tableComponent.vue";

import Operator_container from "@/components/operator_container.vue";
import ut from "@/lib/util";
import {Log, LOG_OPERATOR_TYPE, SystemLog, SystemLogListParams} from "@/api/adminpb/log";

const {t} = useI18n()

const operatorContainerRef = ref(null)
const operatorContainer2Ref = ref(null)

const loading = ref(false)
const Count = ref(0)
const timeRange = ref([])
const SystemLogListParam: Ref<SystemLogListParams> = ref(<SystemLogListParams>{
    OperateMod:   0,
    StartTime:    0,
    EndTime:      0,
    OperateType:  LOG_OPERATOR_TYPE.BALANCE_EDIT,
    OperatorName: "",
    Page: 1,
    PageSize: 20,
})
const defaultGameEvent = ref({})
const defaultOperatorEvent = ref({})
const tableHeader = [
    {label: "商户AppID", value: "Operator"},
    {label: "操作类型", value: "OperatorType", type:"custom"},
    {label: "操作金额", value: "Balance"},
    {label: "操作前金额", value: "PreDeposit"},
    {label: "操作后金额", value: "AfterDeposit"},
    {label: "操作人", value: "OperatorName"},
    {label: "操作时间", value: "OperateTime", format:(row)=> ut.fmtSelectedUTCDateFormat(row.OperateTime), width: "180px"},
    {label: "备注", value: "Remark"},

]
const tableData: Ref<SystemLog[]> = ref(<SystemLog[]>[])




const init = async () => {

    loading.value = true


    let startTime = 0
    let endTime = 0


    if (timeRange.value[0] && timeRange.value[1]) {

        startTime = ut.fmtUTCDate(timeRange.value[0])
        endTime = ut.fmtUTCDate(timeRange.value[1])
    }


    const queryData = <SystemLogListParams>{
        ...SystemLogListParam.value,
    }

    if (startTime) {

        queryData.StartTime = ut.fmtSelectedUTCDate(startTime * 1000, "reduce")
    }
    if (endTime) {

        queryData.EndTime = ut.fmtSelectedUTCDate(endTime * 1000, "reduce")
    }

    if (queryData.AppID){

        queryData.OperateContent = `,"appId":"${queryData.AppID}"`
    }


    const [response, err] = await Client.Do(Log.GetSystemLog, queryData)
    loading.value = false
    if (!err) {
        Count.value = response.All
    }



    tableData.value = response.List.map(item => {

        let turnData = {
            ...item,
            Balance:0,
            OperatorType: 1,
            Remark:"",
            Operator:"",
            PreDeposit:0,
            AfterDeposit:0
        }
        if (item.OperateContent){

            let content = JSON.parse(item.OperateContent)

            if (!content.oldData){
                content.oldData = {
                    Balance:0
                }
            }

            turnData["Balance"] = content.Balance ? content.Balance : 0
            turnData["OperatorType"] = content.OperatorType
            turnData["Remark"] = content.Remark
            turnData["Operator"] = content.appId
            turnData["PreDeposit"] = Number(content.oldData.Balance).toFixed(2)
            turnData["AfterDeposit"] = Number(content.newData.Balance).toFixed(2)
        }

        return turnData
    })
}


const operatorListChange2 = (value) => {
    if (value) {

        SystemLogListParam.value.AppID = value.AppID
    } else {

        SystemLogListParam.value.AppID = ""
    }
}
const operatorListChange = (value) => {
    if (value) {

        SystemLogListParam.value.ParentAppID = value.AppID
    } else {

        SystemLogListParam.value.ParentAppID = ""
    }
}
const changePage = (page) => {
    SystemLogListParam.value.Page = page.currentPage
    SystemLogListParam.value.PageSize = page.dataSize
    init()
}

// 查询报警日志
const queryList = () => {
    tableData.value = []
    init()
}





init()

const resetSearch = () => {
    SystemLogListParam.value = {
        OperateMod:   0,
        StartTime:    0,
        EndTime:      0,
        OperateType:  LOG_OPERATOR_TYPE.BALANCE_EDIT,
        OperatorName: "",
        Page: 1,
        PageSize: 20,
    }
    operatorContainerRef.value.paramData = "ALL"
    operatorContainer2Ref.value.paramData = "ALL"
    timeRange.value = []
}
</script>

<style scoped lang="scss">
.el-link .el-icon--right.el-icon {
  vertical-align: text-bottom;
}
</style>
