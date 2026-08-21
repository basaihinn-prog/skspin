<template>
    <div style="width: 100%;">

        <div class="searchView gameList">
            <el-form
                    style="max-width: 100%"
                    @keyup.enter="initSearch"
            >
                <el-space wrap>

                    <operator_container
                        :is-init="true"
                        :defaultOperatorEvent="defaultOperatorEvent"
                        @select-operatorInfo="operatorListChange"/>


                    <gamelist_container
                        :hase-manufacturer="true"
                        :hase-all="true"
                        :is-init="true"
                        :defaultGameEvent="defaultGameEvent"
                        @select-operator="selectGameList"
                    />

                    <el-form-item :label="$t('日期')" style="margin-left: -10px">
                        <el-date-picker
                            :placeholder="$t('请选择')"
                            v-model="uiData.Date"
                            locale="zh-cn" type="date" value-format="YYYYMMDD"
                        />
                    </el-form-item>
                    <el-form-item>

                        <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
                    </el-form-item>
                </el-space>
            </el-form>
        </div>

        <div class="page_table_context">
            <customTable v-loading="loading" :table-data="tableData" :table-header="tableHeader" @refreshTable="initSearch" tableName="dataSummary_list">

                <!--     今日总押分       -->
                <template #GameID="scope">

                    <div style="display: flex;align-items: center;">
                    <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>

                        <div style="margin: 0 10px">
                        <el-avatar :src="scope.scope.GameIcon" size="small"/>
                        </div>
                            {{ scope.scope.GameID }}
                    </div>
                </template>
                <!--     今日总押分       -->
                <template #NDBetAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.NDBetAmount) }}
                </template>
                <!--     今日总赢分       -->
                <template #NDWinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.NDWinAmount) }}
                </template>
                <!--     今日系统收益       -->
                <template #NDSystemBenefits="scope">
                    <el-text size="small" :type=" scope.scope?.NDSystemBenefits < 0 ? 'danger' : 'success'">
                        {{ ut.toNumberWithComma(scope.scope?.NDSystemBenefits) }}
                    </el-text>
                </template>
                <!--     今日回报率       -->
                <template #NDTotalReturnRate="scope">
                    {{ percentFormatter(0, 0, scope.scope?.NDTotalReturnRate) || '0.00%' }}
                </template>


                <!--     本月总押分       -->
                <template #NMBetAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.NMBetAmount) }}
                </template>
                <!--     本月总赢分       -->
                <template #NMWinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.NMWinAmount) }}
                </template>
                <!--     本月系统收益       -->
                <template #NMSystemBenefits="scope">
                    <el-text size="small" :type=" scope.scope?.NMSystemBenefits < 0 ? 'danger' : 'success'">
                        {{ ut.toNumberWithComma(scope.scope?.NMSystemBenefits) }}
                    </el-text>
                </template>
                <!--     本月回报率       -->
                <template #NMTotalReturnRate="scope">
                    {{ percentFormatter(0, 0, scope.scope?.NMTotalReturnRate) }}
                </template>



                <!--     上月总押分       -->
                <template #LMBetAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.LMBetAmount) }}
                </template>
                <!--     上月总赢分       -->
                <template #LMWinAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope?.LMWinAmount) }}
                </template>
                <!--     上月系统收益       -->
                <template #LMSystemBenefits="scope">
                    <el-text size="small" :type=" scope.scope?.LMSystemBenefits < 0 ? 'danger' : 'success'">
                        {{ ut.toNumberWithComma(scope.scope?.LMSystemBenefits) }}
                    </el-text>
                </template>
                <!--     上月回报率       -->
                <template #LMTotalReturnRate="scope">
                    {{ percentFormatter(0, 0, scope.scope?.LMTotalReturnRate) }}
                </template>
            </customTable>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, reactive, nextTick} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {AdminStatsRpc} from '@/api/stats/stats';
import {AdminGameCenter} from "@/api/gamepb/admin";
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {useI18n} from 'vue-i18n';
import Gamelist_container from "@/components/gamelist_container.vue";
const defaultGameEvent = ref({})
const {t} = useI18n()
import ut from '@/lib/util'

const defaultOperatorEvent = ref({})
const defaultCurrencyEvent = ref({})
let options = ref([])
let gameData = ref(null)

let uiData = reactive({
    Date: '',
    OperatorId: 0,
    GameListId: null,
    Manufacturer: null,
    // CurrencyCode : ""
})
let tableHeader = [
    {sortable: true, width: "220px", label: "游戏名称", value: "GameID", align: "center", type: "custom", fixed: "left"},
    {width: "110px", label: "币种", value: "CurrencyName", align: "center"},
    {width: "110px", label: "商户AppID", value: "AppId", align: "center"},
    {sortable: true, width: "140px", label: "今日总投注", value: "NDBetAmount", type: "custom",align: "center"},
    {sortable: true, width: "130px", label: "今日总赢分", value: "NDWinAmount", type: "custom",align: "center"},
    {sortable: true, width: "160px", label: "今日系统收益", value: "NDSystemBenefits", type: "custom",align: "center"},
    {sortable: true, width: "140px", label: "今日回报率", value: "NDTotalReturnRate", type: "custom",align: "center"},
    {sortable: true, width:"160px",label: "本月总投注", value: "NMBetAmount", type: "custom",align: "center"},
    {sortable: true, width:"160px",label: "本月总赢分", value: "NMWinAmount", type: "custom",align: "center"},
    {sortable: true, width: "160px", label: "本月系统收益", value: "NMSystemBenefits", type: "custom",align: "center"},
    {sortable: true, width: "140px", label: "本月回报率", value: "NMTotalReturnRate", type: "custom",align: "center"},
    {sortable: true, width:"130px",label: "上月总投注", value: "LMBetAmount", type: "custom",align: "center"},
    {sortable: true, width:"130px",label: "上月总赢分", value: "LMWinAmount", type: "custom",align: "center"},
    {sortable: true, width: "160px", label: "上月系统收益", value: "LMSystemBenefits", type: "custom",align: "center"},
    {sortable: true, width:"130px",label: "上月回报率", value: "LMTotalReturnRate", type: "custom",align: "center"},
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

        uiData.OperatorId = value.Id == "ALL" ? null : value.Id
        CurrencyCode.value = value.CurrencyKey
    }else{

        uiData.OperatorId = 0
        CurrencyCode.value = ""
    }
}
const CurrencyListChange = (value) => {
    if (value){

        uiData.CurrencyCode  = value.CurrencyCode
    }else{

        uiData.CurrencyCode  = ""
    }
}
let selectSearch = async () => {
    let [data, err] = await Client.Do(AdminGameCenter.GameList, {})
    if (err) {
        return tip.e(err)
    }
    options.value = data.List.filter((item) => item.ID !== 'lottery')
        .map(item => ({
            ...item,
            value: item.ID,
            label: item.Name,
            Status: item.Status,
        }));
}

const selectGameList = (value) => {
    if (value.gameData){

        uiData.GameListId = value.gameData
    }else{
        uiData.GameListId = ""
    }

    if (value.manufacturer || value.manufacturer == null){
        uiData.Manufacturer = value.manufacturer
    }
}

let initSearch = async () => {


    if (!uiData.OperatorId ){
        tip.e(t("请选择商户"))
        return
    }

    loading.value = true

    let [data, err] = await Client.Do(AdminStatsRpc.GetBetWinTotal, uiData)

    loading.value = false
    if (err) {
        return tip.e(err)
    }
    tableData.value = []
    let arrLastMonth = returnDataMap(data.LastMonth)
    let arrNowDay = returnDataMap(data.NowDay)
    let arrNowMonth = returnDataMap(data.NowMonth)
    let arr = []
    let mostArr = new Map([...arrLastMonth, ...arrNowDay, ...arrNowMonth,])
    for (let key of mostArr) {
        if (!arr.includes(key[1].GameName)) {
            arr.push(key[1].GameName)
        }
    }
    let tableDataValue = arr.filter(i => i !== 'lottery_guess')?.map(item => {
        let list = options.value.find(o => o.value === item)
        if (!list) {
            console.log(item);
        }
        let lm = arrLastMonth.get(list?.label)
        let nd = arrNowDay.get(list?.label)
        let nm = arrNowMonth.get(list?.label)


        return {
            GameID: list?.label,
            GameName: list?.value,
            GameIcon: list?.IconUrl,
            GameManufacturerName: list?.ManufacturerName,
            LMTotalReturnRate: lm?.BetAmount ? (lm?.WinAmount / lm?.BetAmount) : 0 || 0,
            LMSystemBenefits: (lm?.BetAmount - lm?.WinAmount) || 0,
            LMBetAmount: lm?.BetAmount || 0,
            LMWinAmount: lm?.WinAmount || 0,
            NDBetAmount: nd?.BetAmount || 0,
            NDWinAmount: nd?.WinAmount || 0,
            NDTotalReturnRate: nd?.BetAmount ? (nd?.WinAmount / nd?.BetAmount) : 0 || 0,
            NDSystemBenefits: (nd?.BetAmount - nd?.WinAmount) || 0,
            NMTotalReturnRate: nm?.BetAmount ? (nm?.WinAmount / nm?.BetAmount) : 0 || 0,
            NMSystemBenefits: (nm?.BetAmount - nm?.WinAmount) || 0,
            NMBetAmount: nm?.BetAmount,
            NMWinAmount: nm?.WinAmount,
            CurrencyName: lm?.CurrencyName || nd?.CurrencyName || nm?.CurrencyName,
            AppId: lm?.AppId || nd?.AppId || nm?.AppId,

        }
    }) || []
    tableData.value = tableDataValue
    if (uiData.GameListId) {
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
            existingEntry['CurrencyName'] = item.CurrencyName;
            existingEntry['AppId'] = item.AppId;
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
                GameName: name?.value,
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
