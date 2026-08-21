<template>
    <div>
        <div class="searchView">
            <el-form
                    :model="uiData"
                    style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange" :hase-all="true"></operator_container>
                    <el-form-item :label="$t('日期')">
                        <el-date-picker
                                v-model="uiData.Date"
                                locale="zh-cn" type="month" value-format="YYYYMM"
                        />
                    </el-form-item>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>


            </el-space>

        </div>
        <!-- 数据 -->


        <div class="page_table_context">





                <tableChildrenComponent
                        table-name="statment_list"
                        v-loading="loading"
                        :table-data="tableData"
                        :table-header="tableHeader"
                        :page="uiData.Page"
                        :page-size="uiData.PageSize"
                        :count="Count"
                        @page-change="handleSizeChange"
                        @refreshTable="refreshTable"
                    >
                    <template #AppId="scope">
                        <span class="ddcs">{{ scope.scope.AppId }}</span>
                    </template>

                    <template #Type="scope">
                        <span v-if="scope.scope.Type == 2">SLOT</span>
                    </template>

                    <template #Bet="scope">
                        <span v-if="scope.scope.Type == 2">{{ ut.toNumberWithComma(scope.scope.Bet) }}</span>

                    </template>
                    <template #CooperationType="scope">
                        <template v-if="scope.scope.Type == 1">
                            <span style="height: 40px;display: block;line-height: 40px" v-if="scope.scope.Type == 1">{{ $t('收益分成') }}</span>
                            <div style="height: 40px;display: block;line-height: 40px" v-if="scope.scope.Type == 1">{{ $t('流水分成') }}</div>
                        </template>
                        <template v-else>
                            <span style="height: 40px;display: block;line-height: 40px" v-if="scope.scope.CooperationType == 1">{{ $t('收益分成') }}</span>
                            <div style="height: 40px;display: block;line-height: 40px" v-if="scope.scope.CooperationType == 2">{{ $t('流水分成') }}</div>
                        </template>
                    </template>

                    <template #Win="scope">
                        <span v-if="scope.scope.Type == 2">{{ ut.toNumberWithComma(scope.scope.Win) }}</span>
                    </template>

                    <template #Profit="scope">
                        <el-text v-if="scope.scope.Type == 2" size="small"
                                 :type="scope.scope.Profit < 0 ? 'danger' : 'success'">
                            {{ ut.toNumberWithComma(scope.scope.Profit) }}
                        </el-text>
                    </template>

                    <template #PlantRate="scope">

                        <div style="height: 40px;display: flex;align-items: center;justify-content: center"  v-if="scope.scope.Type == 1 || scope.scope.CooperationType == 1">
                            <el-input-number style="width: 120px" :controls="false"
                                             v-model.trim="scope.scope.PlantRate"
                                             @blur="changeRate(scope.scope, scope.scope.index)"
                                             clearable :placeholder="$t('请输入')" size="default">
                                <template #suffix>
                                    %
                                </template>
                            </el-input-number>
                            <el-button plain type="primary" size="small" @click="setExchangeRate(scope.scope)"
                                       style="margin-left: 10px">
                                {{ $t('换算') }}
                            </el-button>
                        </div>
                        <div style="height: 40px;display: flex;align-items: center;justify-content: center" v-if="scope.scope.Type == 1 || scope.scope.CooperationType == 2">
                            <el-input-number style="width: 120px" :controls="false"
                                             v-model.trim="scope.scope.TurnoverPay"
                                             @blur="changeRate(scope.scope, scope.scope.index)"
                                             clearable :placeholder="$t('请输入')" size="default">
                                <template #suffix>
                                    %
                                </template>
                            </el-input-number>
                            <el-button plain type="primary" size="small" @click="setExchangeRate(scope.scope)"
                                       style="margin-left: 10px">
                                {{ $t('换算') }}
                            </el-button>
                        </div>

                    </template>
                    <template #BalanceHeader>

                        <el-tooltip
                            class="box-item"
                            effect="dark"
                            :open-delay="2000"
                            :hide-after="2000"
                            placement="top-start"
                            trigger="click"
                        >

                            <template #content>
                                {{ $t('使用币种汇率 * 目标币种汇率 = 应收外币') }}
                            </template>
                            <el-icon size="20">
                                <QuestionFilled/>
                            </el-icon>
                        </el-tooltip>
                    </template>


                    <template #Receivable="scope">
                        <span v-if="scope.scope.Type == 2">{{ ut.toNumberWithComma(scope.scope.Receivable) }}</span>
                    </template>

                    <template #onlineProfit="scope">
                        <span v-if="scope.scope.Type == 2">{{ ut.toNumberWithComma(scope.scope.onlineProfit) }}</span>
                    </template>

                    <template #exchangeRate="scope">
                        <el-space wrap style="margin: 10px auto"  v-if="scope.scope.Type == 2">
                            <el-input-number style="width: 120px" :controls="false"
                                             v-model.trim="scope.scope.exchangeRate"
                                             clearable :placeholder="$t('请输入')" size="default"/>
                            <el-button plain type="primary" size="small" @click="setExchangeRate(scope.scope)"
                                       style="margin: 0">
                                {{ $t('换算') }}
                            </el-button>
                        </el-space>

                    </template>

                    <template #ChangeReceivable="scope">
                        <div :style="scope.scope.Type == 1 ? {fontWeight: '600'}: ''">
                            <div style="height: 40px;line-height: 40px" v-if="scope.scope.Type == 1 || scope.scope.CooperationType == 1">
                                <span v-if="scope.scope.Type == 1">{{ $t('收益分成') }}：</span>
                                {{ ut.toNumberWithComma(scope.scope.ChangePlantRateReceivable) }}
                            </div>
                            <div style="height: 40px;line-height: 40px" v-if="scope.scope.Type == 1 || scope.scope.CooperationType == 2">
                                <span v-if="scope.scope.Type == 1">{{ $t('流水分成') }}：</span>
                                {{ ut.toNumberWithComma(scope.scope.ChangeTurnoverPayReceivable) }}
                            </div>
                        </div>
                    </template>

                </tableChildrenComponent>

        </div>

    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, Reactive, Ref, watch} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {AdminGroup, MonthBill} from "@/api/adminpb/group";
import Operator_container from "@/components/operator_container.vue";
import {useStore} from "@/pinia";
import ut from "@/lib/util";
import {excel} from "@/lib/excel";
import {useI18n} from "vue-i18n";
import {GameType} from "@/api/gamepb/customer";
import type {TableColumnCtx} from 'element-plus'

const defaultOperatorEvent = ref({})
const store = useStore();
import {Loading, QuestionFilled, RefreshRight} from '@element-plus/icons-vue'

let loading = ref(false)
let {t} = useI18n()
const sortOrder = ref('ascending');
const operatorDataList = ref(null);
const operatorData = ref(null);
const operatorMap = ref(null);
let uiData = ref({
    Date: '',
    Operator: null,
    Page: 1,
    PageSize: 20
})


const tableHeader = ref([


])

const tableUp = ref([
    {label: "商户AppID", value: "AppId", width: "150px", type: "custom"},
    {label: "商户类型", value: "OperatorType", width: "140px", align:"center", format:(row)=> row.Type == 1 ? t('线路商') : t('商户')},
    {label: "日期", value: "Date", align:"center", format:(row) => ut.fmtDate(row.Date, 'YYYY-MM'), width: "140px"},
    {label: "分类",type:"custom", value: "Type", align: "center", width: "140px"},
    {label: "合作模式",type:"custom", value: "CooperationType", align: "center", width: "140px"},
    {label: "投注金额",type:"custom", value: "Bet", align: "center", width: "140px"},
    {label: "收益",type:"custom", value: "Profit", align: "center", width: "140px"},
    {label: "收益比例",type:"custom", value: "PlantRate", width: store.lang != 'zh' ? "300px" : "240px", align: "center", hiddenVisible: true},
    {label: "应收",type:"custom", value: "Receivable", align: "center", width: "140px"},

])


const tableDown = ref([
    {label: "汇率",type:"custom", value: "exchangeRate", width:store.lang != 'zh' ? "300px" : "240px", align: "center", hiddenVisible: true},
    {label: "应收外币",type:"custom", value: "ChangeReceivable", width: "240px", fixed:"right", hiddenVisible: true},
])



const Count = ref(0)
const TypeMap: JsMap<GameType, string> = {
    [GameType.Slot]: t('电子'),
    // [GameType.Fish]: t('捕鱼游戏'),
    [GameType.Poker]: t('百人'),
    [GameType.Lotter]: t('彩票'),
};

let tableData: Ref<MonthBill[]> = ref<MonthBill[]>([])
let excelData = reactive([])

uiData.value.Date = ut.fmtDate(ut.fmtDateSecond(new Date().getTime()/1000), "YYYYMM")
const queryList = async () => {
    loading.value = true

    if (store.AdminInfo.GroupId == 3) {
        uiData.value.Operator = null
    }
    let [mergedTableData, err] = await Client.Do(AdminGroup.GetAppGameInfo, uiData.value)
    loading.value = false

    if (err) {
        return tip.e(err)
    }

    tableData.value = mergedTableData.List;


    if (store.AdminInfo.GroupId == 3) {
        tableData.value = mergedTableData.List[0].children || mergedTableData.List
    }

    Count.value = mergedTableData.Count

    ExchangeRateCalc()
}

const ExchangeRateCalc = () => {
    tableData.value.forEach(item => {



        item.Id = item.AppId

        if (item.CooperationType == 1){
            item.Receivable = item.Profit < 0 ? 0 : item.Profit * (item.PlantRate/100)
        }else{
            item.Receivable = item.Bet < 0 ? 0 : item.Bet * (item.PlantRate/100)
        }
        if (!item.exchangeRate) {

            item.exchangeRate = 1
        }

        item.onlineProfit = 0

        item.ChangePlantRateReceivable =  item.Receivable * item.exchangeRate
        item.ChangeTurnoverPayReceivable =  item.Receivable * item.exchangeRate
        if (item.children && item.children.length) {
            item.ChangePlantRateReceivable = 0
            item.ChangeTurnoverPayReceivable = 0
            item.children.forEach(childItem => {


                childItem.Id = childItem.AppId

                if (!childItem.exchangeRate) {
                    childItem.exchangeRate = 1
                }

                let profitNum = childItem.Profit < 0 ? 0 : childItem.Profit

                if (childItem.CooperationType == 1){


                    if (item.PlantRate > childItem.PlantRate ){

                        childItem.PlantRate = item.PlantRate
                    }




                    childItem.Receivable = profitNum * (childItem.PlantRate / 100)   // 计算收益





                    childItem.onlineProfit = profitNum * parseFloat(Math.abs(item.PlantRate - childItem.PlantRate) / 100)

                    const onlineProfit = childItem.onlineProfit

                    childItem.ChangePlantRateReceivable = childItem.Receivable * childItem.exchangeRate
                    item.ChangePlantRateReceivable +=  profitNum * parseFloat(Math.abs(item.PlantRate - childItem.PlantRate) / 100) * parseFloat(childItem.exchangeRate)

                    item.onlineProfit +=  parseFloat(onlineProfit)
                }else{

                    childItem.Receivable = childItem.Bet < 0 ? 0 : childItem.Bet * (childItem.TurnoverPay / 100)   // 计算收益


                    childItem.onlineProfit = childItem.Bet * parseFloat(Math.abs(item.TurnoverPay - childItem.TurnoverPay) / 100)

                    const onlineProfit = childItem.onlineProfit

                    childItem.ChangeTurnoverPayReceivable = childItem.Receivable  * childItem.exchangeRate
                    item.ChangeTurnoverPayReceivable += childItem.Bet * parseFloat(Math.abs(item.TurnoverPay - childItem.TurnoverPay) / 100) * parseFloat(childItem.exchangeRate)

                    item.onlineProfit +=  parseFloat(onlineProfit)

                }
            })

        }
    })

}


onMounted(async () => {

    tableHeader.value.push(...tableUp.value)

    if (store.AdminInfo.GroupId != 3){
        tableHeader.value.push({label: "线路商收益",type:"custom", value: "onlineProfit", align: "center", width:"140px"},)
    }

    tableHeader.value.push(...tableDown.value)
    await operatorList()

    queryList()
});

const setFeePaymentRate = (row) => {
    row.receivable = ut.toNumberWithComma((row.Win - row.Bet) * row.feePaymentRate)
}

const refreshTable = () => {

    queryList()
}

const changeRate = () => {
    tableData.value.forEach(item => {
        if (item.children && item.children.length) {

            item.children.forEach(childItem => {
                if (childItem.CooperationType == 1 && item.PlantRate > childItem.PlantRate){

                    childItem.PlantRate = item.PlantRate
                }
                if (childItem.CooperationType == 2 && item.TurnoverPay > childItem.TurnoverPay){

                    childItem.TurnoverPay = item.TurnoverPay
                }
            })


        }
    })
}

const setExchangeRate = (row) => {


    ExchangeRateCalc()

}

const operatorListChange = (value) => {
    tableData.value = []
    if (value) {

        uiData.value.Operator = value.Id
    } else {

        uiData.value.Operator = null
    }
    queryList()
}
let operatorParam = reactive({
    PageIndex: 1,
    PageSize: 10000,
    Status: -1
})
const operatorList = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, operatorParam)
    if (err) {
        return tip.e(err)
    }
    operatorDataList.value = data.AllCount === 0 ? [] : data.List


}
import moment from 'moment';
import {AdminInfo} from "@/api/adminpb/info";
import lashenIcon from "@/assets/login/lashen.png";
import TableHeaderVisibleDialog from "@/components/customTable/tableHeaderVisibleDialog.vue";
import TableChildrenComponent from "@/components/customTable/tableChildrenComponent.vue";

const buttonLoading = ref(false)
let xiazai = () => {
    buttonLoading.value = true
    const ids = new Map();
    const summaryIds = new Map();
    const results = [];
    let summaryData = [];
    const options = [
        {key: "typeGame", name: t("分类")},
        {key: "commaBetNew", name: t("投注额")},
        {key: "incomeNew", name: t("收益")},
        {key: "feePaymentRate", name: t("费率")},
        {key: "receivableNew", name: t("应收")},
        {key: "exchangeRate", name: t("汇率")},
        {key: "foreignCurrencyReceivableNew", name: t("应收外币")},
    ];
    const groupByAppId = (data, map) => {
        data.forEach(item => {
            const appId = item.AppId;
            if (!map.has(appId)) {
                map.set(appId, []);
            }
            map.get(appId).push(item);
        });
    };

    groupByAppId(excelData, ids);
    groupByAppId(tableData, summaryIds);
    summaryIds.forEach((data, appId) => {
        data.forEach(row => {
            summaryData.push(options.map(opt => row[opt.key] || ""));
        });
    });

    ids.forEach((data, appId) => {
        const topData = {
            operator: appId, // 根据具体实现获取商户 ID
            billCreationDate: moment().format('YYYYMMDD'), // 账单创建日期
            billDate: uiData.value.Date, // 账期
            settlementCurrency: 'THB' // 结算货币
        };
        results.push({
            data: data,
            options: options,
            sheetName: appId,
            topData: topData
        });
    });

    // 创建汇总 sheet 的数据
    const summaryTopData = {
        operator: '商户合计', // 汇总商户 ID
        billCreationDate: moment().format('YYYYMMDD'), // 账单创建日期
        billDate: uiData.value.Date, // 账期
        settlementCurrency: 'THB' // 结算货币
    };
    const formattedSummaryData = summaryData.map(row => {
        return options.reduce((acc, opt, index) => {
            acc[opt.key] = row[index] || "";
            return acc;
        }, {});
    });
    results.unshift({
        data: formattedSummaryData,
        options: options,
        sheetName: "汇总",
        topData: summaryTopData
    });
    // 导出 Excel
    let knowDate = moment().format('YYYYMM');
    setTimeout(() => {
        buttonLoading.value = false
        // Handle fetched data here
        excel.dump_multi_sheet_new(knowDate, results)

    }, 2000);
}

const handleSizeChange = (value) => {
    uiData.value.PageSize = value.dataSize
    uiData.value.Page = value.currentPage
    queryList()
}
const handleCurrentChange = (value) => {
    uiData.value.Page = value
    queryList()
}

interface SpanMethodProps {
    row: any
    column: TableColumnCtx<any>
    rowIndex: number
    columnIndex: number
}


</script>

<style>

.table_pagination {
    width: 100%;
    height: auto;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 10px;
    margin-top: 10px;
}

.customTableContainer {
    width: 98%;
    margin: 0 auto;
    border-radius: 5px;
    border: 1px solid #e5e5e5;
}

.customTable {
    width: 100%;
    height: auto;
    border: 1px solid #dcdfe6;
    border-bottom: none;
    margin-bottom: 15px;
}

.tableHandleSwitch {
    height: 30px;
    display: flex;
    align-items: center;

}

.tableHandleSwitch > span {
    margin-right: 5px;
}


.tableTool {
    width: 98%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 0;
    margin: auto;
}

.tableToolContainer {
    background: #f7f7f7
}
</style>
