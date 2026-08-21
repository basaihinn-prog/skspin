
<template>
    <div >
        <div class="searchView">
            <el-form
                :model="uiData"
                style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container :is-init="true" :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>
                    <currency_container :hase-all="true" :defaultOperatorEvent ="defaultCurrencyEvent" @select-operatorInfo="currencyListChange"></currency_container>
                    <el-form-item :label="$t('日期')">
                        <el-date-picker
                            v-model="uiData.Date"
                            placeholder="请选择"
                            type="date" value-format="YYYYMMDD"
                        />
                    </el-form-item>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                <el-button type="primary" @click="queryList">
                    {{ $t('刷新') }}
                </el-button>
            </el-space>
        </div>

        <el-row :gutter="16" style="margin-top: 15px">
            <el-col
                v-for="item in cardData"
                :xl="6" :lg="8" :md="12"  class="card"
            >
                <el-card :body-style="{ padding: '0px' }" shadow="hover">
                    <div class="flex-around" style="padding-top: 1rem">
                        <div class="statistic-card img" style="width: 50%">
                            <el-image :src="item.IconUrl" fit="fill"></el-image>
                        </div>
                        <div class="statistic-card card-text" style="width: 50%; text-align: center">
                            <p>{{ $t('游戏ID') }}：{{item.Game}}</p>
                            <p>
                                {{ $t('游戏名称') }}：
                                {{ item.GameName }}
                            </p>
                            <p>{{ $t('游戏难度') }}：{{ $t('困难') }}</p>
                            <p>{{ $t('币种') }}：{{ item.CurrencyName ? $t(item.CurrencyName) : "" }}</p>
                        </div>
                    </div>
                    <div style="padding: 14px" class="flex-around">
                        <div class="statistic-card">
                            <el-statistic :value="ut.toNumberWithComma(item.BetAmount) ">
                                <template #title>
                                    <div style="display: inline-flex; align-items: center">
                                        {{ $t('投注金额') }}
                                    </div>
                                </template>
                            </el-statistic>
                            <div class="statistic-footer">
                                <div class="footer-item">
                                    <span>{{ $t('增幅') }}</span>
                                    <template v-if="item.BetAmountROI">
                                        <span :class="item.BetAmountROI < 0?'red':'green'">
                                            {{ut.toNumberWithCommaNormal(ut.fmtGold(item.BetAmountROI))}}
                                            <el-icon>
                                                <template v-if="item.BetAmountROI > 0">
                                                    <CaretTop />
                                                </template>
                                                <template v-else>
                                                    <CaretBottom />
                                                </template>
                                            </el-icon>
                                        </span>
                                    </template>
                                    <template v-else>
                                        <span>
                                            0
                                        </span>
                                    </template>
                                </div>
                            </div>
                        </div>
                        <div class="statistic-card">


<!--                            <el-statistic :value="isNaN((Number(item.WinAmount) / Number(item.BetAmount))) ? 0 : (Number(item.WinAmount) / Number(item.BetAmount)).toFixed(2) + '%' || $t('暂无')">-->
                            <el-statistic :value="item.BetAmount <= 0 ? '∞' : isNaN(item.WinAmount/ item.BetAmount) ? 0 : percentFormatter(0,0, item.WinAmount/ item.BetAmount) || $t('暂无')">
                                <template #title>
                                    <div style="display: inline-flex; align-items: center">
                                        {{ $t('回报率') }}
                                    </div>
                                </template>
                            </el-statistic>
                            <div class="statistic-footer">
                                <div class="footer-item">
                                    <span>{{ $t('增幅') }}</span>
                                    <template v-if="item.profitabilityROI">
                                        <span :class="item.profitabilityROI < 0?'red':'green'">
                                            {{ut.fmtPercentFour(item.profitabilityROI, 2) == 'NaN%' ? '∞' : ut.fmtPercentFour(item.profitabilityROI, 2) }}
                                            <el-icon>
                                                <template v-if="item.profitabilityROI > 0">
                                                    <CaretTop />
                                                </template>
                                                <template v-else>
                                                    <CaretBottom />
                                                </template>
                                            </el-icon>
                                        </span>
                                    </template>
                                    <template v-else>
                                        <span>
                                            0
                                        </span>
                                    </template>
                                </div>
                            </div>
                        </div>
                        <div class="statistic-card">
                            <el-statistic :value="ut.toNumberWithComma( (item.BetAmount-item.WinAmount))">
                                <template #title>
                                    <div style="display: inline-flex; align-items: center">
                                        {{ $t('收益') }}
                                    </div>
                                </template>
                            </el-statistic>
                            <div class="statistic-footer">
                                <div class="footer-item">
                                    <span>{{ $t('增幅') }}</span>
                                    <template v-if="item.pointsROI">
                                        <span :class="item.pointsROI<0?'red':'green'">
                                            {{ ut.toNumberWithCommaNormal(ut.fmtGold(item.pointsROI))}}
                                            <el-icon>
                                                <template v-if="item.pointsROI>0">
                                                    <CaretTop />
                                                </template>
                                                <template v-else>
                                                    <CaretBottom />
                                                </template>
                                            </el-icon>
                                        </span>
                                    </template>
                                    <template v-else>
                                        <span>
                                            0
                                        </span>
                                    </template>
                                </div>
                            </div>
                        </div>
                        <div class="statistic-card">
                            <el-statistic :value="item.EnterPlrCount">
                                <template #title>
                                    <div style="display: inline-flex; align-items: center">
                                        {{ $t('活跃用户') }}
                                    </div>
                                </template>
                            </el-statistic>
                            <div class="statistic-footer">
                                <div class="footer-item">
                                    <span>{{ $t('增幅') }}</span>
                                    <template v-if="item.EnterPlrCountROI">
                                        <span :class="item.EnterPlrCountROI<0?'red':'green'">
                                            {{ut.toNumberWithCommaNormal(item.EnterPlrCountROI)}}
                                            <el-icon>
                                                <template v-if="item.EnterPlrCountROI>0">
                                                    <CaretTop />
                                                </template>
                                                <template v-else>
                                                    <CaretBottom />
                                                </template>
                                            </el-icon>
                                        </span>
                                    </template>
                                    <template v-else>
                                        <span>
                                            0
                                        </span>
                                    </template>
                                </div>
                            </div>
                        </div>
                        <div class="statistic-card">
                            <el-statistic :value="ut.toNumberWithComma(item.WinAmount)">
                                <template #title>
                                    <div style="display: inline-flex; align-items: center">
                                        {{ $t('总赢分') }}
                                    </div>
                                </template>
                            </el-statistic>
                            <div class="statistic-footer">
                                <div class="footer-item">
                                    <span>{{ $t('增幅') }}</span>
                                    <template v-if="item.WinAmountROI">
                                        <span :class="item.WinAmountROI<0?'red':'green'">
                                            {{ut.toNumberWithCommaNormal(ut.fmtGold(item.WinAmountROI))}}
                                            <el-icon>
                                                <template v-if="item.WinAmountROI>0">
                                                    <CaretTop />
                                                </template>
                                                <template v-else>
                                                    <CaretBottom />
                                                </template>
                                            </el-icon>
                                        </span>
                                    </template>
                                    <template v-else>
                                        <span>
                                            0
                                        </span>
                                    </template>
                                </div>
                            </div>
                        </div>
                        <div class="statistic-card">
                            <el-statistic :value="item.SpinPlrCount">
                                <template #title>
                                    <div style="display: inline-flex; align-items: center">
                                        {{ $t('有效用户') }}
                                    </div>
                                </template>
                            </el-statistic>
                            <div class="statistic-footer">
                                <div class="footer-item">
                                    <span>{{ $t('增幅') }}</span>
                                    <template v-if="item.SpinPlrCountROI">
                                        <span :class="item.SpinPlrCountROI<0?'red':'green'">
                                            {{ut.toNumberWithCommaNormal(item.SpinPlrCountROI)}}
                                            <el-icon>
                                                <template v-if="item.SpinPlrCountROI>0">
                                                    <CaretTop />
                                                </template>
                                                <template v-else>
                                                    <CaretBottom />
                                                </template>
                                            </el-icon>
                                        </span>
                                    </template>
                                    <template v-else>
                                        <span>
                                            0
                                        </span>
                                    </template>
                                </div>
                            </div>
                        </div>
                    </div>
                </el-card>
            </el-col>
        </el-row>
        <el-backtop :right="100" :bottom="100" />

    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick} from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import { dyImport } from "@/pages/analysis/allimage";
import {useStore} from "@/pinia";
import { useI18n } from 'vue-i18n';
import {AdminInfo} from "@/api/adminpb/info";
const { t } = useI18n()
const store = useStore()
let loading = ref(false)
import Operator_container from "@/components/operator_container.vue";
import moment from 'moment';
import ut from "@/lib/util";
import {AdminGameCenter} from "@/api/gamepb/admin";
import Currency_container from "@/components/currency_container.vue";
import tr from "element-plus/es/locale/lang/TR";
// 获取当前日期
const today = moment().format('YYYYMMDD');
console.log(today);
// 格式化日期
let uiData = reactive({
    PageIndex: 1,
    PageSize: 20,
    Operator: null,
    CurrencyCode: null,
    Date: today,
    NeedLast: true,
})
let GameList = reactive([])
let cardData = ref([])
let form = reactive({

})
const defaultOperatorEvent = ref({})
const defaultCurrencyEvent = ref({})


const getGameList = async () => {
    let [data, err] = await Client.Do(AdminGameCenter.GameList, {})
    if (err) {
        return tip.e(err)
    }
    GameList = data.List || []
    nextTick(() => {
        cardData.value.forEach(item=>{
            GameList.forEach(s=> {
                if (s.ID === item.Game) {
                    item['GameName'] = t(s.OriginName)
                    item['IconUrl'] = s.IconUrl
                }
            })
        })
    })

}


// Client.send("/mq/AdminInfo/Interior/getGame", {AppID: "faketrans", Language: "zh"})

const queryList = async () => {
    cardData.value = []

    if (!uiData.Operator && !uiData.CurrencyCode){
        return tip.e(t("商户和币种至少一个必填"))
    }



    if (!uiData.Date){
        return tip.e(t("请选择日期"))
    }

    loading.value=true
    const [data, err] = await Client.Do(AdminAnalysis.Game, uiData);
    if (err) {
        return tip.e(err);
    }
    let yesterday = data.YesterdayList

    let chooseDay = data.List
    const convertToROI = (dataList, yesterdayData) => {
        return dataList?.map((t) => {
            if (yesterdayData) {
                const matchingYesterday = yesterdayData.find((y) => t.Game === y.Game);
                let BetAmountROI = null
                let WinAmountROI = null
                let EnterPlrCountROI = null
                let SpinPlrCountROI = null
                let pointsROI = null
                let profitabilityROI = null
                if (matchingYesterday) {
                    BetAmountROI = Number(t.BetAmount) - Number(matchingYesterday?.BetAmount);
                    WinAmountROI = Number(t.WinAmount) - Number(matchingYesterday?.WinAmount);
                    EnterPlrCountROI = Number(t.EnterPlrCount) - Number(matchingYesterday?.EnterPlrCount);
                    SpinPlrCountROI = Number(t.SpinPlrCount) - Number(matchingYesterday?.SpinPlrCount);

                    pointsROI = (Number(t.BetAmount) - Number(t.WinAmount)) - (Number(matchingYesterday?.BetAmount) - Number(matchingYesterday?.WinAmount));


                    let profitabilityValue = (t.WinAmount / (Number(t.BetAmount) ? Number(t.BetAmount) : 1)  - (matchingYesterday?.WinAmount / (Number(matchingYesterday?.BetAmount) ? Number(matchingYesterday?.BetAmount) : 1)))

                    profitabilityROI = t.BetAmount == 0 || matchingYesterday?.BetAmount == 0 ? '∞' : profitabilityValue

                }


                return {
                    ...t,
                    BetAmountROI,
                    WinAmountROI,
                    EnterPlrCountROI,
                    SpinPlrCountROI,
                    pointsROI,
                    profitabilityROI,
                };
            }
            console.log(t)
            return t;
        });
    };

    let updatedData = convertToROI(chooseDay, yesterday);
    cardData.value = updatedData || [];
    if ((uiData.Date) && !data.List) {
        data.List = []
        data.YesterdayList = []
        updatedData = []
        console.log(data.List);
        tip.e(t('当前日期无数据'))
        return
    }
    console.log(updatedData,'updatedData','cardData.value');
    loading.value=false
    getGameList()
}

const operatorListChange = (value) =>{

    if (value){

        uiData.Operator = value.value
    }else{

        uiData.Operator = ""
    }
}
const currencyListChange = (value) =>{

    if (value){

        uiData.CurrencyCode = value.CurrencyCode
    }else{
        uiData.CurrencyCode = ""
    }
}
const changePage = (PageIndex, PageSize) => {
    uiData.PageIndex = PageIndex
    uiData.PageSize = PageSize
    queryList()
}
</script>
<style lang="scss">

.el-statistic {
    --el-statistic-content-font-size: 28px;
}

.statistic-card {
    height: 100%;
    padding: 20px;
    border-radius: 4px;
    background-color: var(--el-bg-color-overlay);
    margin-bottom: 1rem;
}
@media screen and (max-width: 968px) {
    .statistic-card{
        margin-bottom: .5rem;
    }
}

.statistic-footer {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    font-size: 12px;
    color: var(--el-text-color-regular);
    margin-top: .5rem;
}

.statistic-footer .footer-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.statistic-footer .footer-item span:last-child {
    display: inline-flex;
    align-items: center;
    margin-left: 4px;
}

.green {
    color: var(--el-color-success);
}
.red {
    color: var(--el-color-error);
}
.card-text{
    display: flex !important;
    flex-direction: column;
    align-content: center;
    justify-content: center;
    align-items: center;
    font-size: 1rem;
    text-align: right;
    p{
        margin-bottom: 1rem;
    }
}
@media screen and (max-width: 968px) {
    .card-text{
        font-size: .85rem;
        p{
            margin-bottom: .5rem;
        }
    }
}
.flex-around{
    display: flex;
    justify-content: space-around;
    flex-wrap: wrap;
    align-items: center;
    align-content: center;

    .el-image{
        width: 90%;
        height: auto;
        border-radius:.5rem;
        min-height:238px
    }

    /*.img{
        !*max-width: 200px;
        max-height: 200px;
        width: 100%;
        height: 100%;*!
        width: 200px;
        height: 200px;
        margin-bottom: 0;

    }*/
}
.statistic-card{
    padding: 0;
    width: 50%;
    text-align: center;
    font-size: 15px;
}
@media screen and (max-width: 968px) {
    .statistic-card{
        width: calc(100%/3);
    }
    .el-statistic__number{
        font-size: 1.2rem;
    }
}
.card{
    margin-bottom: 1rem;
    max-width: 600px;
    .el-card{
        border-radius: 1rem;

        p{
            width: 90%;
            margin-left:auto;
            margin-right:auto;

            text-align: left;

        }
    }
}
</style>
