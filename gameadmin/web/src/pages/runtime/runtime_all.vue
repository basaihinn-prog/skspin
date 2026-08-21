<template>
    <div>
        <div class="searchView gameList" >
            <el-form
                style="max-width: 100%"
                @keyup.enter="initSearch"
            >
                <el-space wrap>
                    <el-form-item :label="$t('日期')">
                        <el-date-picker
                            v-model="uiData.Date"
                            locale="zh-cn" type="date" value-format="YYYYMMDD"
                        />
                    </el-form-item>
                    <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange"></operator_container>


                    <currency_container :defaultOperatorEvent="defaultCurrencyEvent"
                                        @select-operatorInfo="currencyListChange"></currency_container>
                    <el-form-item :label="$t('游戏名称')">
                        <el-select-v2
                            v-model="uiData.GameListId"
                            style="width: 220px"
                            filterable
                            clearable
                            autocomplete="on"
                            :options="options"
                            :placeholder="$t('请输入')"
                            @change="selectGameList"
                        >
                            <template #default="{ item }">
                                <el-text style="width: 90%" truncated>
                                    {{ item.label }}
                                </el-text>
                                <span class="value" style="width: 10px;float: right;overflow: hidden">
                                    <el-tag :type="item.Status?item.Status===1?'warning':'info': 'success'" round effect="dark"
                                            style="width: 5px;height: 5px;padding: 0;position: relative;top: -5px"></el-tag>
                                </span>

                            </template>
                        </el-select-v2>
                    </el-form-item>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>
        <div class="page_table_context">
            <customTable :table-data="tableData" :table-header="tableHeader">
                <template #LMSystemBenefits="scope">
                    <el-text :type="scope.scope?.LMSystemBenefits < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope?.LMSystemBenefits) }}</el-text>
                </template>
                <template #LMTotalReturnRate="scope">
                    {{ percentFormatter(0, 0, scope.scope?.LMTotalReturnRate) }}
                </template>
                <template #NMSystemBenefits="scope">

                    <el-text :type="scope.scope?.NMSystemBenefits < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope?.NMSystemBenefits) }}</el-text>

                </template>
                <template #NMTotalReturnRate="scope">{{ percentFormatter(0, 0, scope.scope?.NMTotalReturnRate) }}
                </template>
                <template #SDBetAmount="scope">

                    {{ ut.toNumberWithComma(scope.scope?.SDBetAmount) }}

                </template>
                <template #SDWinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.SDWinAmount) }}

                </template>
                <template #SDSystemBenefits="scope">
                    <el-text :type="scope.scope?.SDSystemBenefits < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope?.SDSystemBenefits) }}</el-text>
                </template>
                <template sortable :label="$t('七日回报率')" #SDTotalReturnRate="scope">
                    {{ percentFormatter(0, 0, scope.scope?.SDTotalReturnRate) }}
                </template>
                <template #NDBetAmount="scope">

                    {{ ut.toNumberWithComma(scope.scope?.NDBetAmount) }}

                </template>
                <template #NDWinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.NDWinAmount) }}
                </template>
                <template #NDSystemBenefits="scope">
                    <el-text :type="scope.scope?.NDSystemBenefits < 0 ? 'danger' : 'success'" size="small">{{ ut.toNumberWithComma(scope.scope?.NDSystemBenefits) }}</el-text>
                </template>
                <template #NDTotalReturnRate="scope">
                    {{ percentFormatter(0, 0, scope.scope?.NDTotalReturnRate) || '0.00%' }}
                </template>
                <template #Operator="scope">
                    <template v-if="scope.scope?.GameName">
                        <el-button link type="primary" size="small" style="margin-bottom: 0"
                                   @click="configFileClick(scope.scope?.GameName)">
                            <el-icon>
                                <SetUp/>
                            </el-icon>
                        </el-button>
                    </template>

                </template>

            </customTable>
        </div>

        <el-dialog v-model="configFileDialog" title="Shipping address" width="800">
            <div class="codeView">
                <codemirror v-if="fileExt !== 'csv' && fileExt !== ''" v-model="json" :placeholder="$t('请输入')"
                            :indent-with-tab="true" :extensions="[javascript(), oneDark]" :tabSize="4"
                            style="min-height: 400px;"/>
                <div v-else></div>
            </div>
            <el-button type="primary" @click="SaveConfig">{{ $t('提交') }}</el-button>
        </el-dialog>
    </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, reactive, nextTick} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {AdminStatsRpc} from '@/api/stats/stats';
import {AdminGameCenter} from "@/api/gamepb/admin";
import {AdminConfigFile} from '@/api/adminpb/json';
import {Codemirror} from "vue-codemirror";
import {javascript} from "@codemirror/lang-javascript";
import {oneDark} from "@codemirror/theme-one-dark";
import beautify from "js-beautify";
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {useI18n} from 'vue-i18n';
import { useStore } from '@/pinia/index';

const store = useStore()
const {t} = useI18n()
import Gamelist_container from "@/components/gamelist_container.vue";
const defaultGameEvent = ref({})
import ut from '@/lib/util'
import Currency_container from "@/components/currency_container.vue";
const defaultOperatorEvent = ref({})
const defaultCurrencyEvent = ref({})
let options = ref([])
let gameData = ref(null)

let uiData = reactive({
    Date: '',
    OperatorId: 0,
    GameListId: null,
    CurrencyCode: "",
    need7: true,
})
let tableHeader = [
    {sortable:true,width: "150px",label: "游戏名称", value: "GameID"},
    {width: "80px", label: "商户", value: "AppId"},
    {width: "80px", label: "币种", value: "CurrencyName"},
    {sortable:true,width: "120px",label: "上月系统收益", value: "LMSystemBenefits", type: "custom"},
    {sortable:true,width: "120px",label: "上月回报率", value: "LMTotalReturnRate", type: "custom"},
    {sortable:true,width: "120px",label: "本月系统收益", value: "NMSystemBenefits", type: "custom"},
    {sortable:true,width: "120px",label: "本月回报率", value: "NMTotalReturnRate", type: "custom"},
    {sortable:true,width: "120px",label: "七日总押分", value: "SDBetAmount", type: "custom"},
    {sortable:true,width: "120px",label: "七日总赢分", value: "SDWinAmount", type: "custom"},
    {sortable:true,width: "120px",label: "七日系统收益", value: "SDSystemBenefits", type: "custom"},
    {sortable:true,width: "120px",label: "七日回报率", value: "SDTotalReturnRate", type: "custom"},
    {sortable:true,width: "120px",label: "今日总押分", value: "NDBetAmount", type: "custom"},
    {sortable:true,width: "120px",label: "今日总赢分", value: "NDWinAmount", type: "custom"},
    {sortable:true,width: "120px",label: "今日系统收益", value: "NDSystemBenefits", type: "custom"},
    {sortable:true,width: "120px",label: "今日回报率", value: "NDTotalReturnRate", type: "custom"}
]
let tableData = ref([])
const value = ref([])
const loading = ref(false)
onMounted(async () => {
    await selectSearch()
});
const CurrencyCode = ref("")
const operatorListChange = (value) => {
    if (value){

        uiData.OperatorId = value.Id
        CurrencyCode.value = value.CurrencyKey
    }else{

        uiData.OperatorId = 0
        CurrencyCode.value = ""
    }
}
const currencyListChange = (value) => {
    if (value){

        uiData.CurrencyCode = value.CurrencyCode
    }else{

        uiData.CurrencyCode = ""
    }
}
let selectSearch = async () => {
    let [data, err] = await Client.Do(AdminGameCenter.GameList, {})
    if (err) {
        return tip.e(err)
    }
    options.value  = data.List.filter((item) => item.ID !== 'lottery')
        .map(item => ({
            value: item.ID,
            label: item.Name,
            Status: item.Status,
        }));
}

const selectGameList = (value) =>{
    uiData.GameListId = value
}

let initSearch = async () => {

    if (uiData.OperatorId == 0 && uiData.CurrencyCode == ""){

        tip.e(t("商户和币种至少一个必填"))
        return
    }

    let [data, err] = await Client.Do(AdminStatsRpc.GetBetWinTotal, uiData)
    if (err) {
        return tip.e(err)
    }
    let arrLastMonth = returnDataMap(data.LastMonth)
    let arrNowDay = returnDataMap(data.NowDay)
    let arrNowMonth = returnDataMap(data.NowMonth)
    let arrSevenDay = returnDataMap(data.SevenDay)
    let arr = []
    let mostArr = new Map([...arrLastMonth, ...arrNowDay, ...arrNowMonth, ...arrSevenDay])
    for (let key of mostArr) {
        console.log("-----------------key---------------------")
        console.log(key)
        if (!arr.includes(key[1].GameName)) {
            arr.push(key[1].GameName)
        }
    }

    tableData.value = []
    let tableDataValue = arr.filter(i=>i!=='lottery_guess')?.map(item => {
        let list = options.value.find(o => o.value === item)
        if (!list) {
            console.log(item);
        }
        let lm = arrLastMonth.get(list?.label)
        let nd = arrNowDay.get(list?.label)
        let sd = arrSevenDay.get(list?.label)
        let nm = arrNowMonth.get(list?.label)
        return {
            GameID: list?.label,
            GameName: list?.value,
            LMTotalReturnRate: lm?.BetAmount?(lm?.WinAmount / lm?.BetAmount):0 || 0,
            LMSystemBenefits: (lm?.BetAmount - lm?.WinAmount) || 0,
            NDBetAmount: nd?.BetAmount,
            NDWinAmount: nd?.WinAmount,
            NDTotalReturnRate: nd?.BetAmount?(nd?.WinAmount / nd?.BetAmount):0 || 0,
            NDSystemBenefits: (nd?.BetAmount - nd?.WinAmount) || 0,
            NMTotalReturnRate: nm?.BetAmount?(nm?.WinAmount / nm?.BetAmount):0 || 0,
            NMSystemBenefits: (nm?.BetAmount - nm?.WinAmount) || 0,
            SDBetAmount: sd?.BetAmount,
            SDWinAmount: sd?.WinAmount,
            SDTotalReturnRate: sd?.BetAmount?(sd?.WinAmount / sd?.BetAmount):0 || 0,
            SDSystemBenefits: (sd?.BetAmount - sd?.WinAmount) || 0,
            CurrencyName: lm?.CurrencyName || nd?.CurrencyName || nm?.CurrencyName || sd?.CurrencyName,
            AppId: lm?.AppId || nd?.AppId || nm?.AppId || sd?.AppId,

        }
    }) || []
    tableData.value = tableDataValue
    if (uiData.GameListId){
        let list = tableDataValue.find(list => list.GameName === uiData.GameListId)
        if (list) {
            tableData.value = []
            tableData.value.push(list)
        }
    }
}
let returnDataMap = (data) => {
    let newMap = new Map()
    data?.map((item) => {
        let name = options.value.find(o => o.value === item.GameID)
        if (!name) {
            name = {
                label: 'lottery_guess',
                value: 'lottery_guess'
            }
        }
        if (newMap.has(name?.label)) {
            const existingEntry = newMap.get(name?.label);
            existingEntry['BetAmount'] += item.Bet;
            existingEntry['WinAmount'] += item.Win;

            if (name?.value.slice(0, 2) === 'pg') {
                existingEntry.GameName = name?.value;
                existingEntry['name'] = name?.label;
            }

            existingEntry['CurrencyName'] = item.CurrencyName;
            existingEntry['AppId'] = item.AppId;
            newMap.set(name?.label, existingEntry)
        } else {
            newMap.set(name?.label, {
                BetAmount: item.Bet,
                WinAmount: item.Win,

                CurrencyName: item.CurrencyName,
                AppId: item.AppId,
                GameName: name?.value
            })
        }

        // return {
        //     GameID:name?.label || 'lottery_guess',
        //     GameName:name?.value || 'lottery_guess',
        //     BetAmount: (item.Bet), // 总赢分
        //     WinAmount: (item.Win), // 总下注
        //     // TotalReturnRate: (item.Win / item.Bet) || 0, // 回报率
        //     // SystemBenefits: (item.Bet - item.Win) || 0, // 系统收益
        // }
    })
    return newMap
}

let json = ref()
let isDark = ref(true)
let fileEvent = ref(null)
let configFileDialog = ref(false)

let fileExt = ref('')
let settings = ref({
    language: 'zh-CN', // 官方汉化
    licenseKey: 'non-commercial-and-evaluation', //去除底部非商用声明
    currentRowClassName: 'currentRow', // 突出显示行
    currentColClassName: 'currentCol', // 突出显示列
    // colHeaders: ["Label", "File", "Type", "Group"],
    colHeaders: true,
    trimWhitespace: false, //去除空格
    rowHeaderWidth: 50, //单元格宽度
    stretchH: 'all',
    rowHeaders: true, // 行标题   布尔值 / 数组/  函数
    contextMenu: true, //右键菜单
    manualColumnResize: true,
    autoWrapRow: true, //自动换行
    width: "100%",
    height: "auto",
})

let configFileName = ref('')
let configFileClick = (i) => {
    console.log(i);
    configFileDialog.value = true
    configFileName.value = i
    getLoadConfig(configFileName.value + '_setting.yaml')
    fileExt.value = 'yaml'
}
const getLoadConfig = async (FileName) => {
    loading.value = true
    let [data, err] = await Client.Do(AdminConfigFile.LoadConfig, {FileName})
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    showData(data.Content)
}
const showData = (data) => {
    nextTick(() => {
        if (fileEvent.value === 'json') {
            json.value = beautify(data, {
                indent_size: 2,
                space_in_empty_paren: true,
            });
        } else {
            json.value = data
        }
    })
}
let FileHistory = ref([])
let isFileHistory = ref(false)
let FileHistoryIndex = ref(null)
const SaveConfig = async () => {
    let param = {
        FileName: configFileName.value + '_setting.yaml',
        Content: json.value
    }
    let [data, err] = await Client.Do(AdminConfigFile.SaveConfig, param)
    if (err) {
        return tip.e(err)
    }
    FileHistory.value = []
    FileHistoryIndex.value = null
    isFileHistory.value = false
    tip.s(t('保存成功'))
    configFileDialog.value = false
}

</script>
<style lang="scss" scoped>
.top-chart {
    .top-chart-col {
        display: flex;
        justify-content: center;
        min-height: 270px;

        :deep(.el-progress-circle) {
            width: calc(100% / 3 - 100px);
            min-width: 300px;
            margin: 0 auto;
        }

        :deep(.el-progress-circle__track) {
            stroke: darkgray;
        }

        :deep(.el-progress__text) {
            font-size: 14px !important;
            width: 60%;
            transform: translate(30%, -50%);

            .top-chart {
                .el-col {
                    margin-top: 1rem;
                }
            }
        }
    }
}

.bottom-chart {
    .bottom-chart-col {
        display: flex;
        justify-content: center;
        margin-top: 5rem;
    }
}

:deep(.el-descriptions__body) {
    margin-top: 1rem;
    padding: 12px 0 0 12px;

    .el-descriptions__label:not(.is-bordered-label) {
        color: rgba(18, 31, 62, 0.74);
    }

    .el-descriptions__content {
        font-weight: bold;
    }
}
</style>
