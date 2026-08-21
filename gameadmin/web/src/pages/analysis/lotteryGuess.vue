<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                label-position="top"
                label-width="100px"
                :model="uiData"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>
                </el-space>
            </el-form>

            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>
        <el-table :data="GuessTime" class="elTable" :header-cell-style="{ background: '#F5F7FA', color: '#333333' }"
                  v-loading="loading" :row-class-name="rowbackground">
            <el-table-column prop="BuyDay" :label="$t('日期')"/>
            <el-table-column prop="AppID" :label="$t('商户')"/>
            <el-table-column prop="OpenDay" :label="$t('开奖期数')"/>
            <el-table-column prop="EnterCount" :label="$t('进入次数')"/>
            <el-table-column prop="BuyPersonCount" :label="$t('有效投注DAU')"/>
            <el-table-column prop="SoldGuessGold" :label="$t('销售金额')"/>
            <el-table-column prop="PeriodSoldGuessGold" :label="$t('总销售金额')"/>
            <el-table-column prop="PeriodGuessProfits" :label="$t('开奖金额')"/>
            <el-table-column prop="PlayersWin" :label="$t('玩家输赢')"/>
            <el-table-column prop="CompanyWin" :label="$t('公司输赢')"/>
            <el-table-column prop="CompanyWinProps" :label="$t('公司输赢占比')"/>
        </el-table>
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
const {t} = useI18n()
const store = useStore()
const formRef = ref<FormInstance>()
const defaultOperatorEvent = ref({})
const uiData = reactive({
    Operator:null,
})
const operatorListChange = (value) =>{
    console.log(value);
    uiData.Operator = value.value
}

let GuessTime = ref([])
let loading = ref(false)

async function callGuessFunc(args: any) {
    // return Client.send('mq/lotter/admin/getOpenDaysByYear', {Year:datePickerValue.value})
    return Client.send(`mq/lotter/admin/guessReportOfDay`, args)
}
import ut from '@/lib/util'

const queryList = async () => {
    loading.value = true
    let [data, err] = await callGuessFunc(uiData)
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    GuessTime.value = data.Result?.map(v => {
        return {
            ...v,
            SoldGuessGold:ut.toNumberWithComma(v.SoldGuessGold).trim(),
            PeriodSoldGuessGold:ut.toNumberWithComma(v.PeriodSoldGuessGold).trim(),
            PeriodGuessProfits:ut.toNumberWithComma(v.PeriodGuessProfits).trim(),
            CompanyWin:ut.toNumberWithComma(v.CompanyWin).trim(),
            PlayersWin:ut.toNumberWithComma(v.PlayersWin).trim(),
            CompanyWinProps:v.CompanyWinProps || t('暂无'),
        }
    })
    console.log(GuessTime.value);
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
