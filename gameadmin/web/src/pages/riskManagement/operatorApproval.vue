<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                    :model="ReviewListParams"
                    style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container ref="operatorContainerRef" :hase-all="true" :is-init="true" :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange" :operator-type="1"/>
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
                    :page="ReviewListParams.Page"
                    :page-size="ReviewListParams.PageSize"
                    :count='Count'
                    @refresh-table="queryList"
                    @page-change="changePage">

                <template #OrderId="scope">
                    {{ scope.scope.OrderId }}
                </template>

                <template #ProfitAndLoss="scope">

                    <el-text :type="scope.scope.ProfitAndLoss[0] == '-' ? 'danger' : 'success'" size="small">
                        {{ scope.scope.ProfitAndLoss }}
                    </el-text>
                </template>
                <template #Operator="scope">

                    <el-button plain :type="scope.scope.Reviewer ? 'warning' : 'primary'"  size="small"
                               @click="checkCurrent(scope.scope)">{{ scope.scope.Reviewer ? $t("已审核") : $t("审核") }}
                    </el-button>

                </template>
            </customTable>

        </div>



    </div>
</template>
<script setup lang="ts">

import type {Ref} from "vue";
import type { EditRequest} from "@/api/adminpb/alarm";
import {ref} from "vue";
import {useI18n} from "vue-i18n";

import {Alarm, AlarmItem} from "@/api/adminpb/alarm";
import {Client} from "@/lib/client";
import customTable from "@/components/customTable/tableComponent.vue";
import {tip} from "@/lib/tip";
import Operator_container from "@/components/operator_container.vue";
import {Merchant, Review, ReviewListParams, ReviewListRes, ReviewStatusParams} from "@/api/adminpb/merchant";
import ut from "@/lib/util";
import {ElMessage, ElMessageBox} from "element-plus";

const {t} = useI18n()

const operatorContainerRef = ref(null)
const operatorContainer2Ref = ref(null)

const loading = ref(false)
const Count = ref(0)
const timeRange = ref([])
const ReviewListParams: Ref<ReviewListParams> = ref(<ReviewListParams>{
    AppID:       "",
    ParentAppID: "",
    StartTime:   0,
    EndTime:     0,
    Page:        1,
    PageSize:    20
})
const defaultGameEvent = ref({})
const defaultOperatorEvent = ref({})
const tableHeader = [
    {label: "线路商", value: "Name", width: "250px"},
    {label: "商户AppID", value: "AppID"},
    {label: "商户序号", value: "Id"},
    {label: "平台费率", value: "PlatformPay", format:(row)=> row.PlatformPay + "%"},
    {label: "钱包类型", value: "WalletMode", format:(row)=> row.WalletMode == 1 ? t("转账钱包") : t("单一钱包")},
    {label: "合作模式", value: "CooperationType", format:(row)=> row.CooperationType == 1 ? t("收益分成") : t("流水分成")},
    {label: "币种", value: "CurrencyKey"},
    {label: "创建时间", value: "CreateTime", format:(row)=> ut.fmtSelectedUTCDateFormat(new Date(row.CreateTime).getTime()), width: "180px"},
    {label: "备注", value: "Remark"},
    {label: "核对人", value: "Reviewer"},
    {label: "操作", value: "Operator", type: "custom", fixed: "right", width: "140px", hiddenVisible: true},
]
const tableData: Ref<Review[]> = ref(<Review[]>[])




const init = async () => {

    loading.value = true


    let startTime = 0
    let endTime = 0


    if (timeRange.value[0] && timeRange.value[1]) {

        startTime = ut.fmtUTCDate(timeRange.value[0]) * 1000
        endTime = ut.fmtUTCDate(timeRange.value[1]) * 1000
    }


    const queryData = <ReviewListParams>{
        ...ReviewListParams.value,
    }

    if (startTime) {

        queryData.StartTime = ut.fmtSelectedUTCDate(startTime, "reduce")
    }
    if (endTime) {

        queryData.EndTime = ut.fmtSelectedUTCDate(endTime, "reduce")
    }


    let [response, err] = await Client.Do(Merchant.ApprovalOperatorList, queryData)
    loading.value = false
    if (!err) {
        Count.value = response.All
    }


    tableData.value = response.List
}


const operatorListChange2 = (value) => {
    if (value) {

        ReviewListParams.value.AppID = value.AppID
    } else {

        ReviewListParams.value.AppID = ""
    }
}
const operatorListChange = (value) => {
    if (value) {

        ReviewListParams.value.ParentAppID = value.AppID
    } else {

        ReviewListParams.value.ParentAppID = ""
    }
}
const changePage = (page) => {
    ReviewListParams.value.Page = page.currentPage
    ReviewListParams.value.PageSize = page.dataSize
    init()
}

// 查询报警日志
const queryList = () => {
    tableData.value = []
    init()
}


const checkCurrent = (row) => {

    if (row.Reviewer) {
        return
    }
    readAlert(<ReviewStatusParams>{OperatorID: row.Id})
}


const readAlert = async (data: ReviewStatusParams) => {

    ElMessageBox.confirm(
        t('是否通过该商户审批'),
        t('提示'),
        {
            confirmButtonText: t('同意'),
            cancelButtonText: t('拒绝'),
            type: 'warning',
        }
    )
        .then(async () => {

            const [response, err] = await Client.Do(Merchant.ApprovalOperator, data)

            if (err) {
                return tip.e(t(err))
            }

            init()
        })


}

init()

const resetSearch = () => {
    ReviewListParams.value = {
        AppID:       "",
        ParentAppID: "",
        StartTime:   0,
        EndTime:     0,
        Page:        1,
        PageSize:    20
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
