<template>
    <div>
        <el-form :inline="true" class="demo-form-inline" style="text-align: center;">
            <el-form-item label="全局奖池">
                <el-input placeholder="全局奖池" v-model="uiData.PoolValue"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button @click="setAwardPool(0)" size="small" type="primary">设置</el-button>
            </el-form-item>
            <el-form-item>
                <el-button @click="exportExcelAction" type="primary">三日玩家导出</el-button>
            </el-form-item>
        </el-form>
        <el-row>
            <el-col :span="8" style="text-align: center;">7日总流水:{{ ut.fmtGold(uiData.TotalExpense7days) }}</el-col>
            <el-col :span="8" style="text-align: center;">7日总产出:{{ ut.fmtGold(uiData.TotalIncome7days) }}</el-col>
            <el-col :span="8" style="text-align: center;">7日总回报率:{{ (uiData.TotalIncome7days /
                uiData.TotalExpense7days).toFixed(4)
            }}</el-col>
        </el-row>
        <el-row>
            <el-col :span="8" style="text-align: center;">总流水:{{ goldFormater(0, 0, uiData.TotalExpense) }}</el-col>
            <el-col :span="8" style="text-align: center;">总产出:{{ goldFormater(0, 0, uiData.TotalIncome) }}</el-col>
            <el-col :span="8" style="text-align: center;">总回报率:{{ (uiData.TotalIncome / uiData.TotalExpense).toFixed(4)
            }}</el-col>
        </el-row>
        <el-table :data="uiData.tableData" style="width: 100%">
            <el-table-column label="日期" prop="Date" width="180"></el-table-column>
            <el-table-column label="DAU" prop="Dau"></el-table-column>
            <el-table-column label="有效下注人数" prop="BetPlayerCount"></el-table-column>
            <el-table-column :formatter="goldFormater" label="平局下注分数" prop="AvgBet"></el-table-column>
            <el-table-column :formatter="goldFormater" label="总押分" prop="Bet"></el-table-column>
            <el-table-column :formatter="goldFormater" label="总赢分" prop="Profit"></el-table-column>
            <el-table-column :formatter="goldFormater" label="系统净输赢" prop="WinLose"></el-table-column>
            <el-table-column label="回报率" prop="RateOfReturn"></el-table-column>
        </el-table>
    </div>
</template>

<script lang="ts" setup>

import { onMounted, ref, reactive } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { excel } from '@/lib/excel';
import { AdminHiloRpc } from '@/api/hilo/admin_rpc';
import { AdminZhaoCaiMaoRpc } from '@/api/zhaocaimao/admin_rpc';
import { AdminBaoZangRpc } from '@/api/baozang/admin_rpc';
import { AdminNiubiRpc } from '@/api/niubi/admin_rpc';
import { AdminRomaRpc } from '@/api/roma/admin_rpc';
import { AdminRomaXRpc } from '@/api/romax/admin_rpc';
import ut from '@/lib/util'
import type { FormInstance, FormRules } from 'element-plus'
import { useStore } from '@/pinia/index';
import { useI18n } from 'vue-i18n';

let uiData = reactive({
    PoolValue: 0,
    tableData: [],
    TotalExpense: 0,
    TotalIncome: 0,
    TotalExpense7days: 0,
    TotalIncome7days: 0,
})

onMounted(() => {
    refresh()
})
const {t} = useI18n()
const exportExcelAction = async () => {

    let [resp, err] = await Client.Do(AdminHiloRpc.RuntimeData, {})
    if (err) {
        return tip.e(err)
    }

    let data = JSON.parse(resp.Json)

    let fmtGold = (v) => ut.fmtGold(v)

    excel.dump(data.List, "hilo", [
        { key: "Date", name: t("日期") },
        { key: "Dau", name: t("DAU") },
        { key: "BetPlayerCount", name: t("有效下注人数") },
        { key: "AvgBet", name: t("平局下注分数"), fmt: fmtGold },
        { key: "Bet", name: t("总押分"), fmt: fmtGold },
        { key: "Profit", name: t("总赢分"), fmt: fmtGold },
        { key: "WinLose", name: t("系统净输赢"), fmt: fmtGold },
        { key: "RateOfReturn", name: t("回报率"), },
    ])
}

const refresh = async () => {
    let [resp, err] = await Client.Do(AdminHiloRpc.RuntimeData, {})
    if (err) {
        return tip.e(err)
    }

    let data = JSON.parse(resp.Json)

    uiData.TotalExpense = data.Bet
    uiData.TotalIncome = data.Profit
    uiData.tableData = data.List
    uiData.PoolValue = data.AwardPool

    let income7 = 0
    let expense7 = 0

    for (let i = 0; i < 7; i++) {
        income7 += uiData.tableData[i].Income
        expense7 += uiData.tableData[i].Expense
    }

    uiData.TotalIncome7days = income7
    uiData.TotalExpense7days = expense7
}

const setAwardPool = async (flag) => { //0,1,2,3 //0是全局奖池
    if (ut.isNull(uiData.PoolValue)) {
        return tip.e("设置的值不合法")
    }

    let value = uiData.PoolValue

    let [_, err] = await Client.Do(AdminHiloRpc.SetAwardPool, {
        Val: parseInt(value.toString()),
    })

    if (err) {
        return tip.e(err)
    }

    tip.s("成功")
}
</script>
