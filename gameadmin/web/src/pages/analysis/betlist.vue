<template>
    <div>
        <div class="searchView">
            <el-form

                    :model="uiData"
                    style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :hase-all="true" :is-init="true" :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange"></operator_container>
                    <currency_container :defaultGameEvent="defaultCurrencyEvent" :hase-all="true"
                                        @select-operatorInfo="selectCurrencyList"></currency_container>

                    <el-form-item :label="$t('交易订单号')">
                        <el-input clearable v-model.trim="searchList.OrderId" :placeholder="$t('请输入')"></el-input>
                    </el-form-item>

                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="searchList.Pid" :placeholder="$t('请输入')"
                                  clearable></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('玩家ID')">
                        <el-input clearable v-model.trim="searchList.UserID" :placeholder="$t('请输入')"></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('游戏类型')">

                        <el-select v-model="searchList.GameType" style="width: 150px" clearable @clear="clearLogType" @change="GameClassChange">
                            <el-option v-for='(value, key) in GameClass' :label="value"
                                       :value="parseInt(key as unknown as string)"/>
                        </el-select>

                    </el-form-item>
                    <gamelist_container
                            :hase-manufacturer="true"
                            :hase-all="true"
                            :defaultGameEvent="defaultGameEvent"
                            :is-init="true"
                            @select-operator="selectGameList"/>
                    <!--                    <el-form-item :label="$t('货币')">
                                            <el-input clearable v-model="uiData.Pid" :placeholder="$t('Pid')"></el-input>
                                        </el-form-item>-->
                    <el-form-item :label="$t('投注金额')">
                        <!--                        <el-input clearable v-model.trim="searchList.Bet"  maxlength="15" :placeholder="$t('请输入')" @input="NotNegative('Bet', $event)"></el-input>-->
                        <el-input clearable v-model.trim="searchList.Bet" maxlength="15" :placeholder="$t('请输入')"
                                  onkeyup="value=value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1')">
                            <template #prefix>≥
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('总赢分')">
                        <!--                        <el-input clearable v-model.trim="searchList.Win"  maxlength="15" :placeholder="$t('请输入')"  @input="NotNegative('Win',$event)"></el-input>-->
                        <el-input clearable v-model.trim="searchList.Win" maxlength="15" :placeholder="$t('请输入')"
                                  onkeyup="value=value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1')">
                            <template #prefix>≥
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('玩家输赢金额')">
                        <!--                        <el-input clearable v-model.trim="searchList.WinLose" :placeholder="$t('请输入')" maxlength="15"  @input="keyupListen('WinLose',$event)"></el-input>-->
                        <el-input clearable v-model.trim="searchList.WinLose" :placeholder="$t('请输入')" maxlength="15"
                                  onkeyup="value=value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1')">
                            <template #prefix>≥
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('玩家余额')">
                        <!--                        <el-input clearable v-model.trim="searchList.Balance" :placeholder="$t('请输入')" maxlength="15" @input="NotNegative('Balance', $event)"></el-input>-->
                        <el-input clearable v-model.trim="searchList.Balance" :placeholder="$t('请输入')" maxlength="15"
                                  onkeyup="value=value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1')">
                            <template #prefix>≥
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('事件类型')">


                        <el-select v-model="searchList.LogType" :placeholder="$t('请输入')" style="width: 150px" clearable @clear="clearLogType">
                            <el-option v-for="(item, index) in LogTypeList" :key="index" :label="$t(item.label)" :value="item.value"></el-option>
                        </el-select>

                    </el-form-item>
                    <el-form-item :label="$t('投注时间')">
                        <el-date-picker
                                v-model="searchList.times"
                                type="datetimerange"
                                value-format="x"
                                format="YYYY-MM-DD HH:mm:ss"
                                :range-separator="$t('至')"
                                :shortcuts="shortcuts"
                                :clearable="false"
                                :start-placeholder="$t('开始时间')"
                                :end-placeholder="$t('结束时间')"
                        />
                        <el-button type="primary" style="margin-left: 10px" @click="queryListSearch">
                            {{ $t('搜索') }}
                        </el-button>
                    </el-form-item>

                </el-space>
            </el-form>
            <!--            <el-button type="primary" @click="queryList" :disabled="disButton">{{ $t('搜索') }}</el-button>-->
            <el-space>

                <el-tooltip placement="top" effect="light">
                    <template #content>
                        {{ $t('必填项：') }}
                        <br/>
                        {{ $t('商户，玩家ID，投注时间') }}
                        <br/>
                        {{ $t('投注的起止时间不能超过7天') }}
                    </template>

                </el-tooltip>

            </el-space>
        </div>

        <div class="page_table_context">
            <customTable
                    v-loading="loading"
                    table-name="betLog_list"
                    :table-header="tableHeader"
                    :table-data="uiData.tableData"
                    :page="searchList.PageIndex"
                    :page-size="searchList.PageSize"
                    :count="Count"
                    :cell-click="clickCell"
                    @refresh-table="queryListSearch"
                    @page-change="pageChange">
                <template #ID="scope">
                    <el-link type="primary" v-if="scope.scope.LogType == 0 && scope.scope.GameManufacturerName != 'SPRIBE'">{{ scope.scope.ID }}</el-link>
                    <div type="primary" v-else>{{ scope.scope.ID }}</div>

                    <template v-if="scope.scope.GameManufacturerName != 'SPRIBE'">
                        <p v-if="scope.scope.GameID === 'lottery' || 'lottery_1'">{{ scope.scope.RoundID }}</p>
                        <p v-else>{{  scope.scope.ManufacturerName == 'pp' ? scope.scope.RoundID : scope.scope.PGBetID }}</p>
                    </template>
                </template>
                <template #InsertTime="scope">
                    <div style="display: flex;align-items: center;">
                        <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>

                        <div style="margin: 0 10px">
                            <el-avatar :src="scope.scope.GameIcon" size="small"/>
                        </div>
                        {{ scope.scope.Game_id }}
                    </div>
                </template>
                <template #GameID="scope">
                    <template v-if="scope.scope.LogType == 0 || (scope.scope.GameManufacturerName == 'SPRIBE')">
                        <div style="display: flex;align-items: center;">
                            <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>

                            <div style="margin: 0 10px">
                                <el-avatar :src="scope.scope.GameIcon" size="small"/>
                            </div>
                            {{ scope.scope.Game_id }}
                        </div>
                        {{ findGameName(scope.scope.Game_id?.startsWith('lottery') ? 'lottery' : scope.scope.Game_id) }}
                    </template>
                    <template v-else>
                        /
                    </template>

                </template>

                <template #LogType="scope">
                    <el-tag :type="scope.scope.LogType ? scope.scope.LogType === 1 ? 'success':'danger':'info'">
                        {{ LogType[scope.scope.LogType] }}
                    </el-tag>
                </template>
                <template #TransferAmount="scope">
                    {{ ut.toNumberWithComma(scope.scope.TransferAmount) }}
                </template>
                <template #Bet="scope">
                    {{ ut.toNumberWithComma(scope.scope.Bet) }}
                </template>
                <template #Win="scope">
                    {{ ut.toNumberWithComma(scope.scope.Win) }}
                </template>
                <template #WinLose="scope">
                    {{ ut.toNumberWithComma(scope.scope.WinLose) }}
                </template>
                <template #Balance="scope">
                    {{ ut.toNumberWithComma(scope.scope.Balance) }}
                </template>
            </customTable>
        </div>

        <el-dialog v-model="pageList.dialogVisible" class="slot-dialog" width="90%"
                   :style="pageList.isPg?{'max-width': '520px'}:''" :title="$t('投注信息')" append-to-body>
            <div class="dialog-top py-10" v-show="pageList.isPg && pageList.whitchPg !== 'lottery'">
                <el-row :gutter="20">
                    <el-col :span="6">
                        <div class="top">
                            <p>
                                {{ pageList.row.ID }}
                            </p>
                        </div>
                        <div class="color-text">{{ $t('交易订单号') }}</div>
                    </el-col>
                    <el-col :span="6">
                        <div class="top">
                            <p>
                                {{ ut.toNumberWithComma(pageList.row.Bet) }}
                            </p>
                        </div>
                        <div class="color-text">{{ $t('投注金额') }}</div>
                    </el-col>
                    <el-col :span="6">
                        <div class="top">
                            <p>
                                {{ ut.toNumberWithComma(pageList.row.WinLose) }}
                            </p>
                        </div>
                        <div class="color-text">{{ $t('玩家输赢金额') }}</div>
                    </el-col>
                    <el-col :span="6">
                        <div class="top">
                            <p>
                                {{ ut.toNumberWithComma(pageList.row.Balance) }}
                            </p>
                        </div>
                        <div class="color-text">{{ $t('玩家余额') }}</div>
                    </el-col>
                </el-row>
            </div>
            <BetListCommentShow :oneLog="pageList.row"/>
        </el-dialog>
        <el-dialog v-model="isPgDialog" top="4vh" class="slot-dialog" style="max-width:620px;height: 880px"
                   :title="$t('投注信息')"
                   append-to-body>
            <iframe :src="betDetailsUrl" ref="iframeScale"
                    style=" width: 100%;height: 100%;position: absolute;left: 0"></iframe>
        </el-dialog>
    </div>
</template>
<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick} from 'vue';
import {useStore} from '@/pinia/index';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import ut from '@/lib/util'
import {useI18n} from 'vue-i18n';
import {AdminGameCenter} from '@/api/gamepb/admin';
import Operator_container from "@/components/operator_container.vue";
import BetListCommentShow from "@/pages/slots/BetListCommentShow.vue"
import Gamelist_container from "@/components/gamelist_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {storeToRefs} from 'pinia'
import moment from 'moment';
import {GameClass} from "@/enum";
const store = useStore()
const {activeTabs, language, Token} = storeToRefs(store)
const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})
const defaultCurrencyEvent = ref({})

let loading = ref(false)
let disButton = ref(false)
let LogTypeList = ref([])


let SlotLogTypeList = [
    {label: "全部", value: -1},
    {label: "电子投注/结算(游戏)", value: 0},
    {label: "转入", value: 1},
    {label: "转出", value: 2},
    {label: "admin设置", value: 3},
]

let MiniLogTypeList = [
    {label: "全部", value: -1},
    {label: "电子投注/结算(游戏)", value: 0},
    {label: "退还", value: 1},
]

let tableHeader = [
    {label: "交易订单号", value: "ID", type: "custom", width: "220px", fixed: "left"},
    {label: "游戏名称", value: "GameID", type: "custom", width: "300px"},
    {label: "商户AppID", value: "AppID", width: "150px"},
    {label: "币种", value: "CurrencyName", width: "100px"},
    {label: "唯一标识", value: "Pid", width: "120px"},
    {label: "事件类型", value: "LogType", type: "custom", width: "150px"},
    {label: "转移金额", value: "TransferAmount", type: "custom", width: "150px"},
    {label: "投注金额", value: "Bet", type: "custom", width: "150px"},
    {label: "总赢分", value: "Win", type: "custom", width: "150px"},
    {label: "玩家输赢金额", value: "WinLose", type: "custom", width: "150px"},
    {label: "玩家余额", value: "Balance", type: "custom", width: "150px"},
    {label: "玩家ID", value: "UserID", width: "300px"},
    {label: "记录生成时间", value: "InsertTime", width: "200px", format: (row) => ut.fmtSelectedUTCDateFormat(row.InsertTime)},
]

let uiData = reactive({
    tableData: [],
    Pid: "",
    GameID: "",
    Operator: 0,
    CurrencyCode: "",
    CurrencyName: "",
    Manu: "",
    GameList: []
})
let pageList = reactive({
    row: null,
    dialogVisible: false,
    isPg: false,
    whitchPg: '',
})
const Count = ref(0)
const CurrentPage = ref(0)
const NextPage = ref(0)
let searchList = reactive({
    OperatorId: 0,
    OrderId: '',
    CurrencyCode: '',
    GameType: 0,
    CurrencyName: '',
    Manufacturer: '',
    NextId: '',
    NextStartTime: null,
    UserID: null,
    PageIndex: 1,
    PageSize: 20,
    Pid: undefined,
    times: null,
    StartTime: null,
    EndTime: null,
    GameID: '',
    Bet: null,
    Win: null,
    WinLose: null,
    Balance: null,
    Lottery_Peroid: null,
    LogType: -1,
})
let {t} = useI18n()

const initSearchTime = () => {
    const today = new Date();
    const year = today.getUTCFullYear();
    const month = today.getUTCMonth();
    const date = today.getUTCDate();

    const startTime = new Date(year,month,date, 0,0,0).getTime()
    const endTime =  new Date(year,month,date, 23,59,59).getTime()


    searchList.times = [startTime, endTime];
}
initSearchTime()

const shortcuts = [
    {
        text: t('过去7天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            return [Date.parse(start.toString()), Date.parse(end.toString())]
        },
    },
    {
        text: t('今天'),
        value: () => {
            const today = new Date();
            const year = today.getFullYear();
            const month = today.getMonth();
            const date = today.getDate();

            const startTime = new Date(year, month, date, 0, 0, 0, 0o00).getTime();
            const endTime = new Date(year, month, date, 23, 59, 59, 0o00).getTime();

            return [startTime, endTime];
        },
    },
]

// 0 游戏, 1 转入, 2 转出, 3 admin设置, 4 彩票投注, 5 彩票结算, 6 百人投注/结算
let LogType = ref([])

let SlotLogTypeMap = ref([
    t('电子投注/结算'),
    t('+ 转入'),
    t('- 转出'),
    t('admin设置'),
    t('彩票投注'),
    t('彩票结算'),
    t('百人投注/结算'),
])
let MINILogTypeMap = ref([
    t('电子投注/结算'),
    t('+ 退回'),
])
let initGameList = ref(null)
let initSearch = async () => {
    loading.value = true
    let [datas, errs] = await Client.Do(AdminGameCenter.GameList, {})
    loading.value = false
    if (errs) {
        return tip.e(errs)
    }
    initGameList.value = datas.List
}
const findGameName = (gameID) => {
    let game = initGameList.value?.find(t => t.ID === gameID);
    return game ? game.Name : ''
}
const findGame = (gameID) => {
    let game = initGameList.value?.find(t => t.ID === gameID);

    return game ? game : null
}
const isPgDialog = ref(false)
const clickCell = (row, column, cell, event) => {
    if (column.property != "ID") return false;

    if (row.LogType == 3) {
        return false
    }
    if (row.GameManufacturerName == "SPRIBE"){
        return false
    }

    pageList.row = null
    pageList.row = row
    // 每一个游戏都需要加入一个投注金额=> scope.row.Bet
    let pgDenominator = {
        'YingCaiShen': '20',
        'BaoZang': '20',
        'XingYunXiang': '20',
        'JinNiu': '10',
        'MaJiang': '20',
        'MaJiang2': '20',
        'TuZi': '10',
        'ZhaoCaiMao': '20',
        'Roma': '15',
        'RomaX': '15',
        'NiuBi': '1',
    }
    const pgArray = [
        'YingCaiShen', 'BaoZang', 'XingYunXiang',
        'SugarRush', 'JinNiu', 'MaJiang', 'MaJiang2',
        'TuZi', 'Hilo', 'Bonaza', 'Olympus', 'lottery',
        'ZhaoCaiMao', 'CowBoy', 'pg_1489936', 'Olympus1000', 'Starlight1000', 'Starlight', 'StarlightChristmas']
    let isPg = pgArray.includes(pageList.row.GameID)
    pageList.isPg = isPg
    pageList.whitchPg = pageList.row.GameID

    if (pageList.row.GameID.split('_').length > 1 &&
        (
            pageList.row.GameID.startsWith('pg_') ||
            pageList.row.GameID.startsWith('jili_') ||
            pageList.row.GameID.startsWith('pp_')
        )) {
        pgDialogQuery(pageList.row)

        pageList.dialogVisible = false
    } else {
        pageList.dialogVisible = true
        isPgDialog.value = false
    }
}

const betDetailsUrl = ref('')
const pgDialogQuery = async (n) => {

    let datas = {
        Gid: n.GameID,
        BetID: n.PGBetID,

        Lang: language.value,
        Token: Token.value
    }
    if (n.GameID.startsWith('pg_')) {
        datas['SID'] = n.ID
        datas['BetID'] = n.PGBetID
        let [data, err] = await Client.Do(AdminGameCenter.GetBetDetailsUrl, datas)
        loading.value = false
        if (err) {
            return tip.e(err)
        }
        betDetailsUrl.value = data?.Url
    }
    if (n.GameID.startsWith('pp_')) {
        datas['SID'] = n.PGBetID
        datas['BetID'] = n.RoundID
        let [data, err] = await Client.Do(AdminGameCenter.GetPPBetDetailsUrl, datas)
        loading.value = false
        if (err) {
            return tip.e(err)
        }
        betDetailsUrl.value = data?.Url
        const iframeDocument = document.querySelector("iframe");
        console.log(iframeDocument);
    }

    if (n.GameID.startsWith('jili_') || n.GameID.startsWith('tada_')) {
        datas['BetID'] = n.ID
        let [data, err] = await Client.Do(AdminGameCenter.GetJiliBetDetailsUrl, datas)
        loading.value = false
        if (err) {
            return tip.e(err)
        }
        betDetailsUrl.value = data?.Url
    }
    isPgDialog.value = true
}
const rowbackground = (row, rowIndex) => {
    if (row.row.Frb) return 'rowbackground'
}
const keyupListen = (key, value) => {

    let upValue = value.replace(/[^-\d\.]/g, '')
    if (value.indexOf(".") != -1 && value.split(".")[1].length > 3) {
        upValue = `${value.split(".")[0]}.${value.split(".")[1].slice(0, 3)}`
        upValue = parseFloat(upValue)
    }

    searchList[key] = upValue

}
const NotNegative = (key, value) => {
    searchList[key] = value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g, '$1')
}
const selectGameList = (value, lotteryValue) => {
    searchList.Lottery_Peroid = lotteryValue.value
    searchList.GameID = value.gameData
    searchList.Manufacturer = value.manufacturer == "ALL" ? null : value.manufacturer


}
const selectCurrencyList = (value, lotteryValue) => {
    if (value) {

        searchList.CurrencyCode = value.CurrencyCode == "All" ? "" : value.CurrencyCode
        uiData.CurrencyCode = value.CurrencyCode
        uiData.CurrencyName = value.CurrencyName
    } else {
        searchList.CurrencyCode = ""
        uiData.CurrencyCode = ""
        uiData.CurrencyName = ""
    }
}
const changePageNum = ref(null)
let initPage = ref(0)
const queryListSearch = async () => {
    searchList.PageIndex = 1
    CurrentPage.value = 1
    NextPage.value = 1
    await queryList()
}
let searchId = null;
const queryList = async () => {
    // if (!searchList.OperatorId) {
    //     return tip.e(t('请选择商户'))
    // }





    loading.value = true
    disButton.value = true
    let dataList = {
        ...searchList,
        Manufacturer: searchList.Manufacturer ? searchList.Manufacturer.toLowerCase(): null  ,
        Pid: Number(searchList.Pid),
        PageSize: Number(searchList.PageSize),
        StartTime: ut.fmtSelectedUTCDate(searchList.times ? ut.fmtUTCDate(searchList.times[0]) * 1000 : null, "reduce"),
        EndTime: ut.fmtSelectedUTCDate(searchList.times ? ut.fmtUTCDate(searchList.times[1]) * 1000 : null, "reduce"),
        Bet: Number(searchList.Bet) ? parseFloat((Number(searchList.Bet) * 10000).toFixed(2)) : null,
        Win: Number(searchList.Win) ? parseFloat((Number(searchList.Win) * 10000).toFixed(2)) : null,
        WinLose: Number(searchList.WinLose) ? parseFloat((Number(searchList.WinLose) * 10000).toFixed(2)) : null,
        Balance: Number(searchList.Balance) ? parseFloat((Number(searchList.Balance) * 10000).toFixed(2)) : null,
    }

    dataList = JSON.parse(JSON.stringify(dataList))



    if (uiData.tableData.length){

        let firstTime = ut.fmtDateSecond(uiData.tableData[0].InsertTime)
        let firstID = uiData.tableData[0].ID
        let lastTime = ut.fmtDateSecond(uiData.tableData[uiData.tableData.length - 1].InsertTime)
        let lastID = uiData.tableData[uiData.tableData.length - 1].ID

            if(CurrentPage.value < NextPage.value){
                dataList["UpOrderId"] = firstID
                dataList["UpEndTime"] = ut.fmtUTCDate(new Date(firstTime).getTime())
            }else if(CurrentPage.value > NextPage.value){
                dataList["NextOrderId"] = lastID
                dataList["NextStartTime"] = ut.fmtUTCDate(new Date(lastTime).getTime())
                dataList.PageIndex =  Math.abs(CurrentPage.value - NextPage.value) || 1
            }



    }


    if (typeof changePageNum.value === "object" && changePageNum.value === null) {
        dataList['LastId'] = null
        initPage.value = 0
    } else {
        if (!!uiData.tableData) {
            if (changePageNum.value === 0) {
                searchId = dataList['LastId'] = uiData.tableData[0].ID

                dataList['Optional'] = 1
                initPage.value--
            } else {
                searchId = dataList['LastId'] = uiData.tableData[uiData.tableData.length - 1].ID
                dataList['Optional'] = 0
                initPage.value++
            }
        }

    }


    let [data, err] = await Client.Do(AdminGameCenter.BetLogList, dataList)
    loading.value = false
    disButton.value = false
    uiData.tableData = []
    if (err) {
        return tip.e(err)
    }
    Count.value = data.Count
    if (data.Count != 0) {


        data.List.forEach(item => {
            let GameIcon = ""
            let GameManufacturerName = ""

            let gameItem = findGame(item.GameID?.startsWith('lottery') ? 'lottery' : item.GameID)
            if (gameItem) {


                GameIcon = gameItem.IconUrl
                GameManufacturerName = gameItem.ManufacturerName
            }


            item.GameIcon = GameIcon
            item.GameManufacturerName = GameManufacturerName

            item.Game_id = findGameName(item.GameID?.startsWith('lottery') ? 'lottery' : item.GameID)


        })


    }


    uiData.tableData = data.Count === 0 ? [] : data.List

    if (!uiData.tableData) {
        dataList['LastId'] = searchId
    }
}
const changePage = (num) => {
    changePageNum.value = Number(num)
    queryList()
}
onMounted(() => {
    LogTypeList.value = SlotLogTypeList
    LogType.value = SlotLogTypeMap.value
    // queryList()
    initSearch()
});

const handleInput = (value) => {

}
const CurrencyCode = ref("")
const operatorListChange = (value) => {

    if (value) {
        searchList.OperatorId = value.Id == "ALL" ? null : value.Id
        CurrencyCode.value = value.CurrencyKey
    } else {
        searchList.OperatorId = null
        CurrencyCode.value = ""
    }
}

function isWithinSevenDays(startTime, endTime) {
    const diffTime = endTime - startTime;
    const sevenDaysInMilliseconds = 7 * 24 * 60 * 60 * 1000;
    return diffTime <= sevenDaysInMilliseconds;
}

const pageChange = (pageConfig) => {


    CurrentPage.value = pageConfig.currentPage
    NextPage.value = searchList.PageIndex

    searchList.PageSize = pageConfig.dataSize
    searchList.PageIndex = pageConfig.currentPage
    queryList()
}
const GameClassChange = (value) => {
    uiData.tableData = []
    Count.value = 0
    if (value == 0){

        LogTypeList.value = SlotLogTypeList
        LogType.value = SlotLogTypeMap.value
    }else{
        LogTypeList.value = MiniLogTypeList
        LogType.value = MINILogTypeMap.value
    }
}
const clearLogType = () => {
    searchList.LogType = -1
}
</script>
<style lang="scss" scoped>
:deep(.rowbackground) {
  background-color: #E0C3FC;
  //background-image: linear-gradient(62deg, #8EC5FC 0%, #E0C3FC 100%);
}

/*:deep(.rowbackground:hover){
    background-color: #852cda !important;
    //background-image: linear-gradient(62deg, #8EC5FC 0%, #E0C3FC 100%);
}*/
:deep(.rowbackground::after) {
  content: 'free round bonus';
  letter-spacing: 10px;
  position: absolute;
  font-size: 30px;
  opacity: .2;
  right: 15%;
  font-family: monospace;
  font-style: italic;
  font-weight: 100;
}

.el-table .cell {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dialog-top {
  background: #282834;
  font-size: 14px;
  color: hsla(0, 0%, 100%, .6);

  .el-col-6 {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    text-align: center;

    .top {
      display: flex;
      justify-content: center;
      flex-direction: column;
      height: 100%;
      word-wrap: break-word;
      color: rgba(255, 255, 255, 0.6);

      p {
        text-align: center;
      }
    }

    .color-text {
      text-align: center;
      word-wrap: break-word;
      color: rgb(88, 245, 109);
    }
  }

}

:deep .el-carousel__container {
  min-height: 600px;
  max-height: 900px;
  height: auto !important;
}

.py-10 {
  padding: .5rem;
}




</style>
