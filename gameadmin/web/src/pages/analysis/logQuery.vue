<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                :model="uiData"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :hase-all="true" :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>
                    <el-form-item :label="$t('操作用户')">
                        <el-input v-model.trim="searchList.OpName" clearable :placeholder="$t('请输入')" />
                    </el-form-item>
<!--                    <el-form-item :label="$t('变化值')">
                        <el-input v-model="searchList.Change" clearable type="number"
                                  :placeholder="$t('请输入')" />
                    </el-form-item>-->
                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="searchList.AnimUserPid" clearable oninput="value=value.replace(/[^\d]/g,'')"
                                  :placeholder="$t('请输入')" />
                    </el-form-item>
                    <el-form-item :label="$t('操作类型')">
                        <el-select v-model.number="searchList.Type"
                                   :options="fieldNameData"
                                   style="width: 150px"
                                   clearable :placeholder="$t('请输入')" >
                            <el-option
                                v-for="item in fieldNameData"
                                :key="item.Key"
                                :label="item.Desc"
                                :value="item.Key"
                            />
                        </el-select>

                    </el-form-item>
                    <el-form-item :label="$t('记录时间')">
                        <el-date-picker
                            v-model="searchList.timesChange"
                            type="datetimerange"
                            value-format="x"
                            format="YYYY-MM-DD HH:mm:ss"
                            :range-separator="$t('至')"
                            :shortcuts="shortcuts"
                            :start-placeholder="$t('开始时间') "
                            :end-placeholder="$t('结束时间') "
                        />
                    </el-form-item>
                    <el-form-item :label="$t('核对时间')">
                        <el-date-picker
                            v-model="searchList.times"
                            type="datetimerange"
                            value-format="x"
                            format="YYYY-MM-DD HH:mm:ss"
                            :range-separator="$t('至')"
                            :shortcuts="shortcuts"
                            :start-placeholder="$t('开始时间') "
                            :end-placeholder="$t('结束时间') "
                        />
                    </el-form-item>

                    <el-button type="primary" @click="queryList" style="margin-left: 10px">{{ $t('搜索') }}</el-button>
                </el-space>
            </el-form>


        </div>

        <div class="page_table_context">
        <customtable
            table-name="logQuery_list"
            v-loading="loading"
            :table-data="GuessTime"
            :table-header="tableHeader"
            :page="searchList.PageNumber"
            :page-size="searchList.PageSize"
            :count="searchList.Count"
            @refresh-table="queryList"
            @pageChange="pageChange">
            <template #Currency="scope">
                {{scope.scope.Currency || '/'}}
            </template>
            <template #Change="scope">
                    <template v-if="scope.scope.Type === 4">
                        /
                    </template>
                    <template v-else>
                        <!--                        <el-icon style="cursor: pointer" @click="copyToClipboardNum(scope.row.Change)"><CopyDocument /></el-icon>-->
                        {{ut.toNumberWithComma(scope.scope.Change)}}
                    </template>
            </template>
            <template #OldGold="scope">
                <template v-if="scope.scope.Type === 4">
                    {{ statusList[scope.scope.OldGold] }}
                </template>
                <template v-else>
                    {{ut.toNumberWithComma(scope.scope.OldGold)}}
                </template>
            </template>
            <template #NewGold="scope">
                <template v-if="scope.scope.Type === 4">
                    {{ statusList[scope.scope.NewGold] }}
                </template>
                <template v-else>
                    {{ut.toNumberWithComma(scope.scope.NewGold)}}
                </template>
            </template>
            <template #EnsureStatus="scope">
                <el-tag type="success" effect="light" >
                    {{EnsureStatusType[scope.scope.EnsureStatus]}}
                </el-tag>
            </template>
        </customtable>
        </div>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, h, shallowRef, nextTick} from 'vue';
import type {ElInput, FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {WebSideMap} from '@/api/gamepb/customer';
import {AdminGameCenter} from '@/api/gamepb/admin';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import moment from "moment";
import Operator_container from "@/components/operator_container.vue";
import customtable from "@/components/customTable/tableComponent.vue";
const {t} = useI18n()
const store = useStore()
const formRef = ref<FormInstance>()
const defaultOperatorEvent = ref({})
const searchList = reactive({
    Operator:null,
    OperatorId:null,
    StartTime:null,
    EndTime:null,
    times:null,
    timesChange:null,
    OpName:'',
    // Change:null,
    AnimUserPid:null,
    Type:null,
    EnsureStatus:1,
    EnsureStartTime:null,
    EnsureEndTime:null,
    PageSize:20,
    PageNumber:1,
    Count:null,
})
const shortcuts = [
    {
        text: t('过去7天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            return [Date.parse(start.toString()), Date.parse(end.toString())]
        },
    },
    {
        text: t('今天'),
        value: () => {
            const today = new Date();
            const year = today.getFullYear();
            const month = today.getMonth();
            const date = today.getDate();

            const startTime = new Date(year, month, date, 0, 0, 0, 0o00).getTime();
            const endTime = new Date(year, month, date, 23, 59, 59, 0o00).getTime();

            return [startTime, endTime];
        },
    },
]
let fieldNameData = ref(null)
const statusList = [
    t('启用'),
    t('禁用'),
]
const EnsureStatusType = [
    t('未核对'),
    t('已核对'),
]
const operatorListChange = (value) =>{
    console.log(value);
    searchList.OperatorId = value.value
}

let GuessTime = ref([])
let loading = ref(false)

let tableHeader = [
    {label: "记录时间", value: "Time", format:(row)=>ut.fmtSelectedUTCDateFormat(row.Time)},
    {label: "操作类型", value: "Type", format:(row)=> fieldNameData.value.filter((item) => row.Key === item.Type)[0].Desc},
    {label: "唯一标识", value: "AnimUserPid"},
    {label: "所属商户", value: "AppID"},
    {label: "币种", value: "Currency", type: "custom"},
    {label: "变化值", value: "Change", type: "custom"},
    {label: "调整前数据", value: "OldGold", type: "custom"},
    {label: "状态", value: "EnsureStatus", type: "custom"},
    {label: "操作用户", value: "OpName"},
    {label: "备注", value: "Remark"},
    {label: "核对时间", value: "EnsureTime", format:(row)=>ut.fmtSelectedUTCDateFormat(row.EnsureTime)},
    {label: "核对人", value: "EnsureOpName"},
]




async function callGuessFunc(args: any) {
    // return Client.send('mq/lotter/admin/getOpenDaysByYear', {Year:datePickerValue.value})
    return Client.send(`mq/lotter/admin/guessReportOfDay`, args)
}
import ut from '@/lib/util'
import {AdminStatsRpc} from "@/api/stats/stats";
const fieldNameReq = async () => {
    let [data, err] = await Client.Do(AdminStatsRpc.GetSlotPoolHistoryType, {})
    if (err) {
        return tip.e(err)
    }
    fieldNameData.value = data.Types
    console.log(data);
    let a = Object.entries(fieldNameData.value).find(t=>{
        console.log(t[1]);
    })
    console.log(a);
}
const queryList = async () => {
    loading.value = true
    if (searchList.timesChange && searchList.timesChange.length) {
        searchList.StartTime = ut.fmtSelectedUTCDate(ut.fmtUTCDate(searchList.timesChange[0]) * 1000, "reduce")
        searchList.EndTime = ut.fmtSelectedUTCDate(ut.fmtUTCDate(searchList.timesChange[1]) * 1000, "reduce")
    } else {
        searchList.StartTime = null
        searchList.EndTime = null
    }
    if (searchList.times && searchList.times.length) {
        searchList.EnsureStartTime = ut.fmtSelectedUTCDate(ut.fmtUTCDate(searchList.times[0]) * 1000, "reduce")
        searchList.EnsureEndTime = ut.fmtSelectedUTCDate(ut.fmtUTCDate(searchList.times[1]) * 1000, "reduce")
    } else {
        searchList.EnsureStartTime = null
        searchList.EnsureEndTime = null
    }
    if (!searchList.AnimUserPid) searchList.AnimUserPid = null
    if (!searchList.Type) searchList.Type = null
    // if (!searchList.Change) searchList.Change = null
    let newSearchList = {
        ...searchList,
        // Change:searchList.Change ? Number(searchList.Change) * 10000 : null
    }


    if (newSearchList.OperatorId == "ALL"){
        newSearchList.OperatorId = null
    }
    let [data, err] = await Client.Do(AdminStatsRpc.GetEnsureSlotPoolHistoryList, newSearchList)
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    console.log(data);
    GuessTime.value = data.Count === 0 ? [] : data.List
    console.log(GuessTime.value);
    searchList.Count = data.Count
}

const rowbackground = (row,rowIndex) => {
    if(row.row.IsFinal) return 'rowbackground'
}
const copyToClipboardNum = (i) => {
    ut.copyToClipboardNum(i/10000)
}
onMounted(()=>{
    fieldNameReq()
    queryList()
});

const pageChange = (page) => {
    searchList.PageNumber = page.currentPage
    searchList.PageSize = page.dataSize

    queryList()
}

</script>
<style scoped lang='scss'>
.elCard {
    border: none;
    margin-bottom: 1rem;
    .el-card {
        border-radius: .5rem;
        height: 100%;
        :deep(.el-card__header){
            padding: 1rem;
        }
        :deep(.el-card__body){
            padding: 1rem;
        }
    }

    .topCard {
        //padding-left: 1.5rem;
        position: relative;

        span {
            //margin-right: .5rem;
        }

        .topCardId {
            font-size: 10px;
            color: #999999;
            margin-bottom: .2rem;
        }

    }

    .bottom {
        display: flex;
        justify-content: space-between;
        width: 100%;
        position: relative;
        //padding-left: 1.5rem;
        font-size: 13px;

        div {
            padding: 0 3px;
            line-height: 24px;
        }
    }
}

.dialog-openList-body {
    .el-card {
        margin-bottom: 1rem;

        .d-flex {
            padding: 0.2rem 0.4rem;
            border: 1px solid #ccc;
            border-radius: 1rem;
            min-width: 83px;
            text-align: center;
            p{
                font-size: 12px;
            }
        }
    }
}
.cell {
    height: 30px;
    padding: 3px 0;
    box-sizing: border-box;
}
.cell .text {
    width: 24px;
    height: 24px;
    display: block;
    margin: 0 auto;
    line-height: 24px;
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    border-radius: 50%;
}
.cell.current .text {
    background: #626aef;
    color: #fff;
}
.cell .holiday {
    position: absolute;
    width: 6px;
    height: 6px;
    background: var(--el-color-danger);
    border-radius: 50%;
    bottom: 0px;
    left: 50%;
    transform: translateX(-50%);
}
.el-descriptions {
    :deep(.el-descriptions__body){
        box-shadow: none;
    }
}
:deep(.rowbackground){
    background-color: #8EC5FC;
    background-image: linear-gradient(62deg, #8EC5FC 0%, #E0C3FC 100%);
}
</style>
