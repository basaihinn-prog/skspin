
<template>
    <div>
        <el-space wrap class="small-font">
            <span>刷新列表</span>
            <el-button :icon="Refresh" circle @click="refreshData"></el-button>
            总回报率：{{ uiData.value['outtotalbeishubi'] == '' ? "" : parseFloat(uiData.value['outtotalbeishubi']).toFixed(4) }}%
            购买游戏比例: {{ uiData.value['gameTotalBeiShuBi'] == '' ? "" :
            parseFloat(uiData.value['gameTotalBeiShuBi']).toFixed(4) }}%
                <el-switch @change="changeAutoCreated" active-text="开启自动生成数据" inactive-text="关闭自动生成数据" inline-prompt
                           v-model="uiData.value.IsAutoCreated"></el-switch>

                <el-switch active-text="生成购买小游戏" inactive-text="生成正常数据"
                           style="--el-switch-on-color: #409eff; --el-switch-off-color: #13ce66"
                           v-model="uiData.value.BuyGame" inline-prompt></el-switch>
        </el-space>
        <el-table :data="uiData.tableData" style="width: 100%" v-loading="loading">
            <el-table-column label="区间">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        <!-- {{ item.Flag + "(" + item.MinRate + ", " + item.MaxRate + "]" }} -->
                        {{ item.Flag}}
                    </div>
                </template>
            </el-table-column>

            <el-table-column label="库存参考">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.ExpectCount + "(" + item.CurrentCount + ")" }}</div>
                </template>
            </el-table-column>
            <el-table-column label="剩余个数">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">{{
                            item.NotUseCombineCount }}</div>
                </template>
            </el-table-column>
            <el-table-column label="组合用掉个数">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">{{ item.UseCombineCount }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="配置个数" width="180">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        <el-input placeholder="配置个数" v-model.number="item.DstCount"></el-input>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="分布比例">
                <template #default="scope">
                    <div style="display: flex;justify-content: space-around;align-items: flex-end;line-height: 40px;">
                        <div>
                            <div :key="index" v-for="(item, index) in scope.row">{{ item.fenbubili.toFixed(4) }}%</div>
                        </div>
                        <div style="padding-left: 20px">{{ scope.row.totalfenbubili.toFixed(4) }}%</div>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="中奖倍数" width="280px">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
<!--                        {{ item.DstTimes.toFixed(4)-item.PoolCost.toFixed(4) }},原始值：{{ item.DstTimes1.toFixed(4) }}-->
                        {{ item.DstTimesMinusCost.toFixed(4) }},原始值：{{ item.DstTimes1.toFixed(4) }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="中奖倍数比">
                <template #default="scope">
                    <div style="display: flex;align-items: flex-end;line-height: 40px;">
                        <div>
                            <div :key="index" v-for=" (item, index) in scope.row">{{ item.beishubi }}%</div>
                        </div>

                        <div style="padding-left: 20px">{{ scope.row.totalbeishubi.toFixed(4) }}%</div>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button @click="ExtractAction(scope.$index)" style="vertical-align: bottom;"
                               type="primary">提取</el-button>
                </template>
            </el-table-column>
            <el-table-column label="生成时间" prop="GenerateTime">
                <template #default="scope">
                    <div>
                        <div :key="index" v-for=" (item, index) in scope.row">{{ item.GenerateTime }}</div>
                    </div>
                </template>
            </el-table-column>
        </el-table>
        <div style="text-align:left">组合数组设置</div>
        <el-table :data="uiData.dataCombineConfigTableData">
            <!-- <el-table-column label="组合" prop="name"></el-table-column> -->
            <el-table-column label="组合" prop="Name"></el-table-column>
            <el-table-column label="配置组数" width="180">
                <template #default="scope">
                    <!-- <el-input placeholder="配置组数" v-model.number="scope.row.NeedCombine"></el-input> -->
                    <el-input placeholder="配置组数" v-model.number="scope.row.Count"></el-input>
                    <!-- <el-button @click="syncdataCombineConfigTableData" size="small" type="primary">设置</el-button> -->
                </template>
            </el-table-column>
        </el-table>

        <div style="margin: 20px;display:flex;">
            <div>
                <el-button @click="genCombineDataAction" type="primary">生成组合数据</el-button>
                <!-- {{ uiData.value.GenCombineDataSuccess }} -->
            </div>
            <!-- <div>
                <el-button @click="generatAction" type="primary">序列化生成数据</el-button>
                {{ uiData.value.GenerateDataSuccess }}
            </div> -->
        </div>
    </div>
</template>
<script lang="ts" setup>
import { onMounted, ref, reactive } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { AdminConfigFile } from '@/api/adminpb/json';
import ut from '@/lib/util'
import type { FormInstance, FormRules } from 'element-plus'
import { useStore } from '@/pinia/index';
import { useI18n } from 'vue-i18n';
import {Refresh} from '@element-plus/icons-vue'
class RateItem {
    buyGame
    DstCount
    DstTimes1
    DstTimes
    Minus
    fenbubili
    beishubi
    GenerateTime
}

class ValudData {
    GenerateDataSuccess: boolean
    GenCombineDataSuccess: boolean
    IsAutoCreated: boolean
    BuyGame: boolean
    RateStatus: Array<Array<RateItem>>
}

let uiData = reactive({
    times: 0,
    value: new ValudData(),
    tableData: [],
    dataCombineConfigTableData: [],
    gameName: "MaJiang"
})

let loading = ref(false)

onMounted(() => {
    refreshData()
})

async function callFunc(api: string, args: any) {
    return Client.send(`mq/${uiData.gameName}/AdminSlotsRpc/${api}`,  args)
}

const changeAutoCreated = async () => {
    loading.value = true
    let [data, err] = await callFunc("SwitchAutoCreateData", {
        Status: uiData.value.IsAutoCreated,
        BuyGame: uiData.value.BuyGame,
    })

    loading.value = false
    if (err) {
        console.error(err)
        return tip.e(err)
    }
    tip.s("success")
}

const { t } = useI18n()

async function refreshCombine() {
    let [combine, err]  = await callFunc("combine/get", {})
    if (err) {
        return tip.e(err)
    }
    console.log(combine)
    uiData.dataCombineConfigTableData = combine
}

async function refreshData() {
    await refreshCombine()
    let [data, err] = await callFunc("GetPoolStatus", {})
    loading.value = false
    if (err) {
        console.error(err)
        return tip.e(err)
    }
    // let data = JSON.parse(jsondata.Data)

    uiData.value = data
    let tableData = uiData.value.RateStatus

    let commonTotalCount = 0
    let gameTotalCount = 0

    for (let i = 0; i < tableData.length; i++) {
        let tt = tableData[i]
        for (let j = 0; j < tt.length; j++) {
            let item = tt[j]
            if (item.buyGame) {
                gameTotalCount += item.DstCount
            } else {
                commonTotalCount += item.DstCount
            }
        }
    }

    let commonTotalBeiShu = 0
    let gameTotalBeiShu = 0
    let buyGameMulti = 50 // 购买小游戏需要50倍的押注
    switch (uiData.gameName) {
        case "BaoZang":
            buyGameMulti = 50 // 购买小游戏需要50倍的押注
            break
        case "ZhaoCaiMao":
            buyGameMulti = 75 // 购买小游戏需要75倍的押注
            break
        case "XingYunXiang":
            buyGameMulti = 75 // 购买小游戏需要75倍的押注
            break
        case "YingCaiShen":
            buyGameMulti = 75 // 购买小游戏需要75倍的押注
            break
        default:
            // return tip.e("尚未设置buyGameMulti")
    }

    for (let i = 0; i < tableData.length; i++) {
        let tt = tableData[i]
        tt['totalfenbubili'] = 0
        tt['totalbeishubi'] = 0
        for (let j = 0; j < tt.length; j++) {
            let o = tt[j]
            o.DstTimes1 = o.DstTimes
            o.DstTimes = o.DstTimes - (o.Minus * o.DstCount)

            if (o.buyGame) {
                o.fenbubili = Number.parseFloat((o.DstCount / gameTotalCount).toFixed(6)) * 100
                o.beishubi = Number.parseFloat((o.DstTimes / gameTotalCount / buyGameMulti).toFixed(6)) * 100
            } else {
                o.fenbubili = Number.parseFloat((o.DstCount / commonTotalCount).toFixed(6)) * 100
                o.beishubi = Number.parseFloat((o.DstTimesMinusCost / commonTotalCount).toFixed(6)) * 100
            }

            tt['totalfenbubili'] += o.fenbubili
            tt['totalbeishubi'] += o.beishubi

            if (o.buyGame) {
                gameTotalBeiShu += o.beishubi
            } else {
                commonTotalBeiShu += o.beishubi
            }

            o.beishubi = o.beishubi.toFixed(4)
            o.GenerateTime = o.GenerateTime
        }
    }
    uiData.value['outtotalbeishubi'] = commonTotalBeiShu
    uiData.value['gameTotalBeiShuBi'] = gameTotalBeiShu
    uiData.tableData = tableData
}

const ExtractAction = async (index) => {
    if (uiData.value.IsAutoCreated) {
        return tip.e('请先关闭自动生成数据')
    }

    let datas = uiData.tableData[index]
    console.log(datas)
    let args = []
    for (let i = 0; i < datas.length; i++) {
        let data = datas[i]
        if (ut.isNull(data.DstCount)) {
            data.DstCount = 0
        }

        if (data.CurrentCount <= 0 && data.DstCount > 0) {
            return tip.e(`区间(${data.MinRate},${data.MaxRate}]，没有生成的，不可配置 `)
        }
        args.push({index: data.Index, count: data.DstCount})
    }
    loading.value = true
    let [data, err] = await callFunc("ExtractAction",  args)
    loading.value = false

    console.log(data)
    if (err) {
        return tip.e(err)
    } else {
        refreshData()
        return tip.s("succss")
    }

}
async function  genCombineDataAction() {
    await callFunc("combine/gen", uiData.dataCombineConfigTableData)
}

// const generatAction = async () => {
//     await callFunc("GenerateData", {})
// }


</script>
<style lang="less" scoped>
.small-font{
    font-size: 12px;
    color: rgba(18, 31, 62, 0.74);
    margin-bottom: 0 !important;
    :deep(.el-button){
        margin-bottom: 0 !important;
    }

}
</style>
