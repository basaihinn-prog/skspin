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
                    <el-form-item :label="$t('日期')">
                        <el-date-picker
                            v-model="uiData.openDay" :editable="false" clearable
                            type="date"
                            format="YYYY-MM-DD"
                            value-format="YYYY-MM-DD"
                        >
                            <template #default="cell">
                                <div class="cell" :class="{ current: cell.isCurrent }">
                                    <span class="text">{{ cell.text }}</span>
                                    <span v-if="isHoliday(cell)" class="holiday" />
                                </div>
                            </template>
                        </el-date-picker>
                    </el-form-item>
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>
        <el-row :gutter="12">
            <el-col :span="4" class="elCard" v-for="item in OpenTime.List" v-loading="loading">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <el-text style="font-size: 1rem;margin-right: 1rem" type="primary">
                                {{ item.Day }}
                            </el-text>
                            <el-text class="mx-1" style="float: right" type="info">{{ $t('商户') }}：{{ item.AppID }}</el-text>
                        </div>
                    </template>
                    <el-space wrap>
                        <el-statistic :title="$t('DAU')" :value="item.BuyPersonCount" />
                        <el-divider direction="vertical" border-style="dashed" />
                        <el-statistic :title="$t('销售数量')" :value="item.SoldCount" />
                        <el-divider direction="vertical" border-style="dashed" />
                        <el-statistic :title="$t('销售金额')" :value="item.SoldGold" />
                    </el-space>
                </el-card>
            </el-col>
        </el-row>
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

let holidays = ref([])
const defaultOperatorEvent = ref({})

const uiData = reactive({
    openDay: moment().format('YYYY-MM-DD'),
    Operator:null
})
const isHoliday = ({ dayjs }) => {
    return holidays.value.includes(dayjs.format('YYYY-MM-DD'))
}
let OpenTime = reactive({
    List: [],
})

let loading = ref(false)
const operatorListChange = (value) =>{
    console.log(value);
    uiData.Operator = value.value
}
async function callFunc(args: any) {
    // return Client.send('mq/lotter/admin/getOpenDaysByYear', {Year:datePickerValue.value})
    return Client.send(`mq/lotter/admin/reportOfPeriod`, args)
}
import ut from '@/lib/util'
const queryList = async () => {
    loading.value = true
    // let [data, err] = await callFunc({openDay: datePickerValue.value})
    let [data, err] = await callFunc(uiData)
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    OpenTime.List = data.Result?.map(v => {
        return {
            ...v,
            SoldGold:ut.toNumberWithComma(v.SoldGold).trim(),
        }
    })
    console.log(OpenTime.List);
}

const OpenPrizeButton = async () => {
    let [data,err] = await Client.send(`mq/lotter/admin/reportOfGetPeriodDays`, {})
    if (err) {
        return tip.e(err)
    }
    holidays.value = data.openDays
}
onMounted(() => {

    OpenPrizeButton()

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
</style>
