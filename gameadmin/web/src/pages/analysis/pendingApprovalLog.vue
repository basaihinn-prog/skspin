<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form

                :model="searchList"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :hase-all="true" :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>

                    <el-form-item :label="$t('操作用户')">
                        <el-input v-model.trim="searchList.OpName" clearable :placeholder="$t('请输入')" />
                    </el-form-item>
                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="searchList.AnimUserPid" clearable
                                  :placeholder="$t('请输入')" />
                    </el-form-item>
                    <el-form-item :label="$t('操作类型')">
                        <el-select style="width: 150px" v-model.number="searchList.Type"
                                   :options="fieldNameData"
                                   clearable :placeholder="$t('请输入')" >
                            <el-option
                                v-for="item in fieldNameData"
                                :key="item.Key"
                                :label="item.Desc"
                                :value="item.Key"
                            />
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('日期查询')">
                        <el-date-picker
                            v-model="searchList.times"
                            type="datetimerange"
                            value-format="x"
                            format="YYYY-MM-DD HH:mm:ss"
                            :range-separator="$t('至')"
                            :shortcuts="shortcuts"
                            :start-placeholder="$t('开始时间')"
                            :end-placeholder="$t('开始时间')"
                        />
                        <el-button type="primary" @click="queryList" style="margin-left: 10px">{{ $t('搜索') }}</el-button>
                    </el-form-item>

                </el-space>
            </el-form>

            <el-space wrap>

            </el-space>
        </div>

        <div class="page_table_context">
        <customTable
            table-name="pendingApprovalLog_list"
            v-loading="loading"
            :table-header="tableHeader"
            :table-data="GuessTime"
            :count="searchList.Count"
            :page="searchList.PageNumber"
            :page-size="searchList.PageSize"
            @refresh-table="queryList"
            @page-change="changePage">

            <template #Type="scope">
                {{Object.entries(fieldNameData).find(i=>i[1].Key === scope.scope.Type)[1]?.Desc}}
            </template>

            <template #Currency="scope">
                {{scope.scope.Currency || '/'}}
            </template>

            <template #OldGold="scope">
                <template v-if="scope.scope.Type === 4">
                    {{ statusList[scope.scope.OldGold] }}
                </template>
                <template v-else>
                    {{ut.toNumberWithComma(scope.scope.OldGold)}}
                </template>
            </template>

            <template #Change="scope">
                <template v-if="scope.scope.Type === 4">
                    /
                </template>
                <template v-else>
                    {{ut.toNumberWithComma(scope.scope.Change)}}
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
                <el-tag type="warning" effect="light">
                    {{ scope.scope.EnsureStatus?'':$t('未核对') }}
                </el-tag>
            </template>
            <template #Remark1="scope">
                <el-space>
                    <el-input
                        v-model="scope.scope.Remark"
                        style="width: 120px"
                        maxlength="100"
                        show-word-limit
                        type="textarea"
                    />
                </el-space>
            </template>
            <template #Remark2="scope">
                <el-button type="primary" size="small" plain @click="changeStatus(scope.scope,false)">
                    {{$t('审核')}}
                </el-button>
            </template>
        </customTable>

        </div>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, h, shallowRef, nextTick} from 'vue';
import type {ElInput, FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {AdminStatsRpc} from "@/api/stats/stats";
import ut from '@/lib/util'
import {
    Finished,
} from '@element-plus/icons-vue'
import {ElMessageBox} from "element-plus";
const {t} = useI18n()
const store = useStore()
const formRef = ref<FormInstance>()
let tableHeader = [

    {label: "记录时间", value: "Time", width: "180px", format:(row)=>ut.fmtSelectedUTCDateFormat(row.Time)},
    {label: "操作类型", value: "Type", type: "custom"},
    {label: "唯一标识", value: "AnimUserPid"},
    {label: "所属商户", value: "AppID"},
    {label: "币种", value: "Currency", type: "custom"},
    {label: "调整前数据", value: "OldGold", type: "custom", width: "160px"},
    {label: "变化值", value: "Change", type: "custom"},
    {label: "调整后数据", value: "NewGold", type: "custom", width: "160px"},
    {label: "状态", value: "EnsureStatus", type: "custom"},
    {label: "操作用户", value: "OpName"},
    {label: "备注", value: "Remark1",type: "custom", width:"240px", hiddenVisible: true},
    {label: "操作", value: "Remark2",type: "custom", width:"180px", fixed:"right", hiddenVisible: true},
]
let GuessTime = ref([])
let loading = ref(false)
let fieldNameData = ref(null)
const defaultOperatorEvent = ref({})
const searchList = reactive({
    OperatorId:null,
    StartTime:null,
    EndTime:null,
    times:null,
    OpName:'',
    // Change:null,
    AnimUserPid:null,
    Type:null,
    EnsureStatus:0,
    EnsureStartTime:null,
    EnsureEndTime:null,
    PageSize:20,
    PageNumber:1,
    Count:null,
})
const statusList = [
    t('禁用'),
    t('启用'),
]
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
const operatorListChange = (value) =>{
    console.log(value);
    searchList.OperatorId = value.value
}

const changeStatus = async (datas,bool) => {

    ElMessageBox.confirm(
        t('确认审核当前预警'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            let [data, err] = await Client.Do(AdminStatsRpc.GetDoEnsureslotPoolHistory,
                {
                    ID: datas.ID,
                    Remark: datas.Remark
                }
            )
            if (err) {
                queryList()
                return tip.e(err)
            }
            tip.s(t('操作成功'))
            queryList()

    })

}

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
    if (searchList.times && searchList.times.length) {
        searchList.StartTime = ut.fmtSelectedUTCDate(ut.fmtUTCDate(searchList.times[0]) * 1000, "reduce")
        searchList.EndTime = ut.fmtSelectedUTCDate(ut.fmtUTCDate(searchList.times[1]) * 1000, "reduce")
    } else {
        searchList.StartTime = null
        searchList.EndTime = null
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
    searchList.Count = data.Count
}
const changePage = (page) => {
    searchList.PageNumber = page.currentPage
    searchList.PageSize = page.dataSize
    queryList()
}
const rowbackground = (row,rowIndex) => {
    if(row.row.IsFinal) return 'rowbackground'
}
onMounted(()=>{
    fieldNameReq()
    queryList()
});

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
