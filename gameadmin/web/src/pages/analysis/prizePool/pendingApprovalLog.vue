<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                label-position="top"
                label-width="100px"
                :model="searchList"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>
                    <el-form-item :label="$t('日期查询')">
                        <el-date-picker
                            v-model="searchList.times"
                            type="datetimerange"
                            value-format="x"
                            format="YYYY-MM-DD HH:mm:ss"
                            range-separator="To"
                            :shortcuts="shortcuts"
                            start-placeholder="start time"
                            end-placeholder="end time"
                        />
                    </el-form-item>
                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="searchList.Pid" clearable
                                  :placeholder="$t('请输入')" />
                    </el-form-item>
                </el-space>
            </el-form>

            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>
        <el-table :data="GuessTime" class="elTable" :header-cell-style="{ background: '#F5F7FA', color: '#333333' }"
                  v-loading="loading" :row-class-name="rowbackground" scrollbar-always-on highlight-current-row>
            <el-table-column prop="Time" :label="$t('记录时间')" :formatter="dateSecondFormater"/>
            <el-table-column prop="Pid" :label="$t('唯一标识')"/>
            <el-table-column prop="AppID" :label="$t('所属商户')"/>
            <el-table-column prop="Currency" :label="$t('币种')" width="80">
                <template #default="scope">
                    {{scope.row.Currency || '/'}}
                </template>
            </el-table-column>
            <el-table-column prop="OriginWin" :label="$t('中奖金额')">
                <template #default="scope">
                    {{ut.toNumberWithComma(scope.row.OriginWin) || '/'}}
                </template>
            </el-table-column>
            <el-table-column prop="RealWin" :label="$t('实际派发金额')">
                <template #default="scope">
                    {{ut.toNumberWithComma(scope.row.RealWin) || '/'}}
                </template>
            </el-table-column>
            <el-table-column prop="Desc" :label="$t('中奖时触发条件')">
            </el-table-column>
            <el-table-column prop="EnsureStatus" :label="$t('状态')">
                <template #default="scope">
                    <el-tag type="warning" effect="light">
                        {{ scope.row.EnsureStatus?'':$t('未核对') }}
                    </el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="Remark" :label="$t('备注')" width="180">
                <template #default="scope">
                    <el-space>
                        <el-input
                            v-model="scope.row.Remark"
                            style="width: 120px"
                            maxlength="100"
                            show-word-limit
                            type="textarea"
                        />
                    </el-space>
                </template>
            </el-table-column>
            <el-table-column prop="Remark" :label="$t('操作')" width="180">
                <template #default="scope">
                    <el-button type="primary" size="small" @click="changeStatus(scope.row,false)">
                        {{$t('审核')}}
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
        <Pagination :Count='searchList.Count' @paginationEmit='changePage'></Pagination>
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
import {AdminStatsRpc} from "@/api/stats/stats";
import ut from '@/lib/util'
const {t} = useI18n()
const store = useStore()
const formRef = ref<FormInstance>()
let GuessTime = ref([])
let loading = ref(false)
const defaultOperatorEvent = ref({})
const searchList = reactive({
    OperatorId:null,
    StartTime:null,
    EndTime:null,
    times:null,
    // Change:null,
    Pid:null,
    EnsureName:null,
    EnsureStatus:0,
    EnsureStartTime:null,
    EnsureEndTime:null,
    PageSize:50,
    PageNumber:1,
    Count:null,
})
const statusList = [
    t('启用'),
    t('禁用'),
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
    let [data, err] = await Client.Do(AdminStatsRpc.GetDoEnsureSlotWinLoseLimit,
        {
            ID:datas.ID,
            Remark:datas.Remark
        }
    )
    if (err) {
        return tip.e(err)
    }
    tip.s(t('操作成功'))
    queryList()
}

const queryList = async () => {
    loading.value = true
    if (searchList.times && searchList.times.length) {
        searchList.StartTime = searchList.times[0]/1000
        searchList.EndTime = searchList.times[1]/1000
    } else {
        searchList.StartTime = null
        searchList.EndTime = null
    }
    if (!searchList.Pid) searchList.Pid = null
    // if (!searchList.Change) searchList.Change = null
    let newSearchList = {
        ...searchList,
        // Change:searchList.Change ? Number(searchList.Change) * 10000 : null
    }
    let [data, err] = await Client.Do(AdminStatsRpc.GetSlotWinLoseLimitList, searchList)
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    console.log(data);
    GuessTime.value = data.Count === 0 ? [] : data.List
    searchList.Count = data.Count
}
const changePage = (PageNumber, PageSize) => {
    searchList.PageNumber = PageNumber
    searchList.PageSize = PageSize
    queryList()
}
const rowbackground = (row,rowIndex) => {
    if(row.row.IsFinal) return 'rowbackground'
}
onMounted(()=>{
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
