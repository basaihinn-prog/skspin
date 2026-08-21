
<template>
    <div >
        <div class="searchView">
            <AnalysisPlatform :uiData='uiData' @query-list="queryList"></AnalysisPlatform>
        </div>
        <!-- 数据 -->
        <el-table :data="uiData.tableDataSlice" class="elTable" v-loading="loading"
                  :header-cell-style="{ background: '#F5F7FA', color: '#333333' }">
            <el-table-column type="index" :label="$t('序号')" />
            <el-table-column prop="Date" :label="$t('投注日期')"/>
            <el-table-column prop="AppID" :label="$t('商户')"/>
            <el-table-column prop="EnterPlrCount" :label="$t('玩家数量')" />
            <el-table-column prop="RegistPlrCount" :label="$t('新玩家数量')" />
            <el-table-column prop="SpinCount" :label="$t('游戏总人次')" ></el-table-column>
            <el-table-column prop="Flow">
                <template #header="scope">
                    <div class="sort-table-div">
                        {{ $t('流水') }}
                        <span class="sort-table-column">
                            <el-icon @click="columnSort('Flow','ascending')" :class="sortBool['Flow'].top?'sort-border-color-top':''"><CaretTop /></el-icon>
                            <el-icon @click="columnSort('Flow','descending')" :class="sortBool['Flow'].down?'sort-border-color-bottom':''"><CaretBottom /></el-icon>
                        </span>
                    </div>
                </template>
                <template #default="scope">
                    {{ ut.toNumberWithComma(scope.row.Flow) }}
                </template>
            </el-table-column>
            <el-table-column prop="playersWinLose" :label="$t('玩家输赢')" >
                <template #header="scope">
                    <div class="sort-table-div">
                        {{ $t('玩家输赢') }}
                        <span class="sort-table-column">
                            <el-icon @click="columnSort('playersWinLose','ascending')" :class="sortBool['playersWinLose'].top?'sort-border-color-top':''"><CaretTop /></el-icon>
                            <el-icon @click="columnSort('playersWinLose','descending')" :class="sortBool['playersWinLose'].down?'sort-border-color-bottom':''"><CaretBottom /></el-icon>
                        </span>
                    </div>
                </template>
                <template #default="scope">
                    {{ ut.toNumberWithComma(scope.row.playersWinLose) }}
                </template>
            </el-table-column>
            <el-table-column prop="companyWinsLoses" :label="$t('公司输赢')">
                <template #header="scope">
                    <div class="sort-table-div">
                        {{ $t('公司输赢') }}
                        <span class="sort-table-column">
                            <el-icon @click="columnSort('companyWinsLoses','ascending')" :class="sortBool['companyWinsLoses'].top?'sort-border-color-top':''"><CaretTop /></el-icon>
                            <el-icon @click="columnSort('companyWinsLoses','descending')" :class="sortBool['companyWinsLoses'].down?'sort-border-color-bottom':''"><CaretBottom /></el-icon>
                        </span>
                    </div>
                </template>
                <template #default="scope">
                    {{ ut.toNumberWithComma(scope.row.companyWinsLoses) }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('公司输赢占比')">
                <template #default="scope">
                    {{ percentFormatter(0,0, ((scope.row.BetAmount-scope.row.WinAmount)/scope.row.BetAmount)) }}
                    <!--                    {{ percentFormatter(0,0, ((scope.row.WinAmount-scope.row.BetAmount)/scope.row.BetAmount)) }}-->
                </template>
            </el-table-column>
        </el-table>
        <!-- 分页 -->
        <!-- 分页 -->
        <el-pagination
            v-model:current-page="uiData.PageIndex"
            v-model:page-size="uiData.PageSize"
            :page-sizes="[20, 50, 100]"
            small="small"
            background="true"
            layout="total, sizes, prev, pager, next, jumper"
            :total="uiData.Count"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />


    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { useStore } from '@/pinia/index';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { useI18n } from 'vue-i18n';
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import AnalysisPlatform from '@/components/analysis_platform.vue';
import moment from "moment";
const { t } = useI18n()
const store = useStore()
import ut from "@/lib/util";
let loading = ref(false)
let uiData = reactive({
    tableData: [],
    tableDataSlice: [],
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    SearchType: SearchType.Day,
    Pid: null,
    times: "",
})
let sortBool = reactive({
    BetAmount: {
        top:false,
        down:false
    },
    WinAmount: {
        top:false,
        down:false
    },
    Flow: {
        top:false,
        down:false
    },
    playersWinLose: {
        top:false,
        down:false
    },
    companyWinsLoses: {
        top:false,
        down:false
    },
})
const queryList = async () => {
    loading.value=true
    let [data, err] = await Client.Do(AdminAnalysis.GetPokerOperator, {
        PageIndex: uiData.PageIndex,
        PageSize: uiData.PageSize,
        Operator: uiData.Pid,
        Date: !uiData.times?'':moment(uiData.times[0]).format('YYYYMMDD'),
        EndDate: !uiData.times?'':moment(uiData.times[1]).format('YYYYMMDD'),
        Type: '',
    })
    loading.value=false
    if (err) {
        return tip.e(err)
    }
    uiData.Count = data.All
    uiData.tableData = data.List?.map(t=>{
        return {
            ...t,
            Flow:(t.Flow),
            playersWinLose:(t.WinAmount-t.BetAmount),
            companyWinsLoses:(t.BetAmount-t.WinAmount),
        }
    }) || []
    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData.slice(startIndex, endIndex);
}
const changePage = () => {
    const startIndex = (uiData.PageIndex - 1) * uiData.PageSize;
    const endIndex = startIndex + uiData.PageSize;
    uiData.tableDataSlice = uiData.tableData.slice(startIndex, endIndex);
}
const handleCurrentChange = async (page: number) => {
    uiData.PageIndex = page
    changePage()
}
const handleSizeChange = (size: number) => {
    uiData.PageSize = size
    changePage()
}
const columnSort = (name, sorts) => {
    const isAscending = sorts === 'ascending';
    const isDescending = !isAscending;
    sortBool = {
        BetAmount: {
            top:false,
            down:false
        },
        WinAmount: {
            top:false,
            down:false
        },
        Flow: {
            top:false,
            down:false
        },
        playersWinLose: {
            top:false,
            down:false
        },
        companyWinsLoses: {
            top:false,
            down:false
        },
    }
    if (isAscending) {
        sortBool[name].top = !sortBool[name].top;
        sortBool[name].down = false;
    } else {
        sortBool[name].down = !sortBool[name].down
        sortBool[name].top = false;
    }

    if ((isAscending && sortBool[name].top) || (isDescending && sortBool[name].down)) {
        uiData.tableData = uiData.tableData.sort((a, b) => (isAscending ? a[name] - b[name] : b[name] - a[name]));
        changePage();
    } else {
        queryList();
    }
};
onMounted(() => {
    queryList()
});
</script>
