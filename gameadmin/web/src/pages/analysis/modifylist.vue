
<template>
    <div>
        <div class="searchView">
            <div>
                <el-input clearable v-model="uiData.Pid" :placeholder="$t('Pid')"></el-input>
            </div>
            <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
        </div>
        <!-- 数据 -->
        <el-table :data="uiData.tableData" @cell-click="clickCell" class="elTable" v-loading="loading"
            :header-cell-style="{ background: '#F5F7FA', color: '#333333' }">
            <el-table-column prop="Pid" label="玩家ID" width="100px" />
            <el-table-column prop="Change" :label="$t('金币变化')" :formatter="goldFormater" >
                <template #default="scope">
                    <el-tag v-if="scope.row.Change > 0" type="success">{{ goldSignedFormater(0, 0, scope.row.Change)
                    }}</el-tag>
                    <el-tag v-else-if="scope.row.Change < 0" type="danger">{{ goldSignedFormater(0, 0,
                        scope.row.Change) }}</el-tag>
                    <div v-else>{{ goldSignedFormater(0, 0, scope.row.Change) }}</div>
                </template>
            </el-table-column>
            <el-table-column prop="Balance" label="余额" :formatter="goldFormater" />
            <el-table-column prop="Duration" label="请求耗时"  >
                <template #default="scope">
                    <div>
                        {{  moment(scope.row.RespTime).diff(moment(scope.row.ReqTime), 'seconds', true) }}s
                    </div>
                </template>
                </el-table-column>
            <el-table-column prop="Status" :label="$t('状态')">
                <template #default="scope">
                    <el-tag :type="(scope.row.Status == 0) ? 'success' : 'danger'" size="small">
                        {{ scope.row.Status == 0 ? $t('成功') : $t('失败') }}
                    </el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="ReqTime" label="请求时间" :formatter="dateFormater" width="200px" />
            <el-table-column prop="Comment" label="备注"> </el-table-column>
        </el-table>
        <!-- 分页 -->
        <Pagination :Count='uiData.Count' @paginationEmit='changePage'></Pagination>
        <el-dialog v-model="uiData.dialogVisible" width="90%">
            <BetListCommentShow :oneLog="uiData.row" />
        </el-dialog>
    </div>
</template>
<style>
.el-table .cell {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { useStore } from '@/pinia/index';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { useI18n } from 'vue-i18n';
import { AdminGameCenter } from '@/api/gamepb/admin';

import BetListCommentShow from "@/pages/slots/BetListCommentShow.vue"
import moment from 'moment';

let loading = ref(false)
let uiData = reactive({
    tableData: [],
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    Pid: "",
    dialogVisible: false,
    row: null
})

const clickCell = (row, column, cell, event) => {
    if (column.property != "Comment") {
        return;
    }
    uiData.dialogVisible = true
    uiData.row = row
}

const queryList = async () => {
    let [data, err] = await Client.Do(AdminGameCenter.ModifyLogList, {
        PageIndex: uiData.PageIndex,
        PageSize: uiData.PageSize,
        Pid: parseInt(uiData.Pid),
    })
    if (err) {
        return tip.e(err)
    }

    uiData.Count = data.Count
    uiData.tableData = data.List
}

onMounted(() => {
    queryList()
});

const changePage = (PageIndex, PageSize) => {
    uiData.PageIndex = PageIndex
    uiData.PageSize = PageSize
    queryList()
}

</script>
