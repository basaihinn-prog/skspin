
<template>
    <div >
        <div class="searchView">
            <el-form
                :model="param"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operatorInfo="operatorListChange" :hase-all="true"></operator_container>
                    <el-form-item :label="$t('唯一标识')">
                        <el-input v-model.number.trim="param.Pid" clearable maxlength="12" oninput="value=value.replace(/[^\d]/g,'')"
                                  :placeholder="$t('请输入')" />
                    </el-form-item>
                    <el-form-item :label="$t('玩家ID')">
                        <el-input v-model.trim="param.Uid" clearable :placeholder="$t('请输入')" />
                    </el-form-item>

                    <Currency_container :defaultOperatorEvent ="defaultCurrencyEvent" @select-operatorInfo="currencyListChange" :hase-all="true"></Currency_container>
                    <el-form-item :label="$t('玩家状态')">
                        <el-select v-model.number="param.Status"
                                   style="width:140px"
                                   clearable
                                   :placeholder="$t('请输入')" >
                            <el-option
                                v-for="(item,index) in statusList"
                                :key="item"
                                :label="$t(item.name)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('在线状态')">
                        <el-select v-model.number="param.OnLineStatus"
                                   style="width:140px"
                                   clearable
                                   :placeholder="$t('请选择')" >
                            <el-option
                                v-for="(item,index) in onlineStatusList"
                                :key="item"
                                :label="$t(item.name)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>

                  <el-form-item>
                        <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                    </el-form-item>
                </el-space>
            </el-form>
        </div>
        <!-- 数据 -->


        <div class="page_table_context">


            <customTable
                table-name="playerAnalysis_list"
                :table-data="uiData.tableData"
                v-loading="loading"
                @refresh-table="queryList"
                :table-header="tableHeader" :page="param.PageIndex"
                :page-size="param.PageSize" :count='uiData.Count' @page-change="changePage">

                <template #handleTools>

                    <el-button type="warning" plain @click="cancelRTPControl"  v-if="visibleRtpControl">{{ $t('批量修改RTP') }}</el-button>
                </template>



                <template #OnLineStatus="scope">
                    <el-tag :type="scope.scope.OnLineStatus == 1 ? 'success' : 'danger' ">{{
                            scope.scope.OnLineStatus == 1 ?
                                $t('在线') :
                                $t('离线')
                        }}</el-tag>
                </template>
                <template #Status="scope">
                    <el-tag :type="scope.scope.Status == 1 ? 'success' : 'danger' ">{{
                            scope.scope.Status == 1 ?
                                $t('启用') :
                                $t('关闭')
                        }}</el-tag>
                </template>
                <template #Bet="scope">
                    {{ ut.toNumberWithComma(scope.scope.Bet) }}
                </template>
                <template #Win="scope">
                    {{ ut.toNumberWithComma(scope.scope.Win) }}
                </template>
                <template #x="scope">
                    {{ scope.scope.Bet > 0 ? percentFormatter(0, 0, scope.scope.Win / scope.scope.Bet) : '0.00%' }}
                </template>
                <template #operator="scope">
                    <el-button type="primary" size="small" plain @click="checkUserInfo(scope.scope)">{{ $t('基础信息') }}</el-button>
                </template>
            </customTable>
        </div>

        <info v-model="userInfoDialog" :player-id="playerId"></info>
        <RTPDialogPanel v-model="RTPControlDialog" :player-id="playerId" @RTPControlData="getRTPPreSubmitData"></RTPDialogPanel>
        <submitRTPControlDialog v-model="commitRTPDialog" :PreSubmitPlayer="RTPPrePlayer" @commitController="closeDialog" :highRTP="HighRTP" marginBottom="0"/>
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { useStore } from '@/pinia/index';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { useI18n } from 'vue-i18n';
import ut from '@/lib/util'
import { AdminAnalysis, SearchType } from '@/api/gamepb/analysis';
import {AdminInfo} from "@/api/adminpb/info";
import {AdminGameCenter} from "@/api/gamepb/admin";
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import Currency_container from "@/components/currency_container.vue";
import Info from "@/pages/analysis/component/Info.vue";
import RTPDialogPanel from "@/pages/analysis/component/cancelRTPDialog.vue";
import da from "element-plus/es/locale/lang/DA";
import SubmitRTPControlDialog from "@/pages/analysis/component/submitRTPControlDialog.vue";
const defaultOperatorEvent = ref({})
const defaultCurrencyEvent = ref({})
const { t } = useI18n()
const store = useStore()
const userInfoDialog = ref(false)
let RTPControlDialog = ref(false)
let commitRTPDialog = ref(false)
let param = reactive({
    Pid: null,
    Uid: "",
    OperatorId: 0,
    CurrencyCode: "",
    times: null,
    StartTime:null,
    EndTime:null,
    Status: -1,
    OnLineStatus: -1,
    PageIndex:1,
    PageSize:20,
    SortBet:0,
    SortWin:0,
})

let visibleRtpControl = ref(false)
let loading = ref(false)
let CurrencyCode = ref("")

let tableHeader = [
    {label: "唯一标识", value: "Pid", width: "120px"},
    // {label: "玩家名称", value: "UName", width: "100px"},
    {label: "玩家ID", value: "Uid", width: "300px"},
    {label: "商户AppID", value: "AppID", width: "120px"},
    {label: "玩家状态", value: "Status",type: "custom", width: "120px"},
    {label: "在线状态", value: "OnLineStatus",type: "custom", width: "120px"},
    {label: "冻结账号时间", value: "CloseTime",format:(row) => `${ut.fmtSelectedUTCDateFormat(row.CloseTime)}`,width: "160px"},
    {label: "币种", value: "CurrencyName"},
    {label: "投注金额", value: "Bet", type: "custom", width: "160px"},
    {label: "总赢分", value: "Win", type: "custom", width: "160px"},
    {label: "回报率", value: "x", type: "custom", width: "160px"},
    {label: "操作", value: "operator", type: "custom", fixed: "right", width: "120px", hiddenVisible: true},
    {label: "玩家注册时间", value: "CreateAt", format:(row) => `${ut.fmtSelectedUTCDateFormat(row.CreateAt)}`,width: "160px"},
    {label: "玩家最后登录时间", value: "LoginAt", format:(row) => `${ut.fmtSelectedUTCDateFormat(row.LoginAt)}`,width: "180px"},
]
const RTPPrePlayer = ref([])
const HighRTP = ref(1)
const playerId = ref("")
let uiData = reactive({
    tableData: [],
    Count: 0,
    PageIndex: 1,
    PageSize: 20,
    SearchType: SearchType.Day,
    Pid: null,
    Uid: "",
    AppID: ''
})
const statusList = [
    {name: t('全部') ,value: -1},
    {name: t('启用') ,value: 1},
    {name: t('关闭') ,value: 0},
]
const onlineStatusList = [
    {name: t('全部') ,value: -1},
    {name: t('在线') ,value: 1},
    {name: t('离线') ,value: 2},
]
const queryList = async () => {
    uiData.tableData = []
    loading.value = true
    let requestParam = {...param}
    if (requestParam.times && requestParam.times.length) {
        requestParam.StartTime = requestParam.times[0]/1000
        requestParam.EndTime = requestParam.times[1]/1000
    } else {
        requestParam.StartTime = null
        requestParam.EndTime = null
    }
    if (requestParam.Status === -1){
        requestParam.Status = null
    }
    if (!requestParam.Pid) {
        requestParam.Pid = null
    }
    let [data, err] = await Client.Do(AdminGameCenter.NewAppList, requestParam)
    loading.value = false
    if (err) {
        return
    }
    if (data.AllCount > 1){

        data.List = data.List.map(row=>({
            ...row,
            CloseTime:new Date(row.CloseTime).getTime(),
            CreateAt:new Date(row.CreateAt).getTime(),
            LoginAt:new Date(row.LoginAt).getTime(),
        }))
    }


    uiData.Count = data.AllCount
    uiData.tableData = data.AllCount === 0 ? [] : data.List
}
const setSortParam = (param, column, order) => {
    const orderMapping = {
        'descending': 1,
        'ascending': -1,
    };

    switch (column.property) {
        case 'Bet':
            param.SortBet = orderMapping[order] || 0;
            break;
        case 'Win':
            param.SortWin = orderMapping[order] || 0;
            break;
        // 如果有其他列需要处理，可以在这里添加 case 分支
        default:
            // 默认情况下不执行任何操作
            break;
    }
}
const tableDataSort = async ({ prop, order,column })=>{
    setSortParam(param, column, order)
    await queryList()
}
onMounted(() => {
    if (store.AdminInfo.GroupId <= 1){
        visibleRtpControl.value = true
    }
    if (store.AdminInfo.GroupId === 3){

        visibleRtpControl.value = store.AdminInfo.Businesses["PlayerRTPSettingOff"] == 1
    }
});

const cancelRTPControl = () => {
    RTPControlDialog.value = true
}

const checkUserInfo = (row) => {
    userInfoDialog.value = true
    playerId.value = row.Pid
}
const changePage = (Page) => {
    param.PageIndex = Page.currentPage
    param.PageSize = Page.dataSize
    queryList()
}
const operatorListChange = (value) =>{
    if (value){

        param.OperatorId = value.Id == 'ALL' ? null : value.Id
        CurrencyCode.value = value.CurrencyKey
    }else{
        param.OperatorId = null
        CurrencyCode.value = ""
    }
}
const currencyListChange = (value) =>{
    if (value){

        param.CurrencyCode = value.CurrencyCode == 'All' ?  '': value.CurrencyCode
    }else{

        param.CurrencyCode = ""
    }
}


const getRTPPreSubmitData = (data, highRTP) => {

    RTPPrePlayer.value = data.value
    HighRTP.value = highRTP.value
    commitRTPDialog.value = true
}
const closeDialog = (data) => {
    RTPControlDialog.value = false
    commitRTPDialog.value = false
}

</script>
