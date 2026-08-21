<template>
    <div class="adminList">
        <div class="searchView">
            <el-form :inline="true" :model="param" class="demo-form-inline">
                <el-form-item :label="$t('商户类型')" v-if="store.AdminInfo.GroupId <= 1">
                    <el-select
                        v-model="param.OperatorType"
                        :placeholder="$t('请选择商户类型')"
                        style="width: 150px"
                        clearable
                    >
                        <el-option v-for="item in operatorType" :label="$t(item.label)" :value="item.value"/>
                    </el-select>
                </el-form-item>

                <operator_container
                    ref="operatorContainerRef"
                    :hase-all="true"
                    :defaultOperatorEvent="defaultOperatorEvent"
                    @select-operatorInfo="operatorListChange"
                ></operator_container>
                <el-form-item :label="$t('国家')">

                    <el-select
                        v-model="param.BelongingCountry"
                        :placeholder="$t('请选择国家')"
                        filterable
                        clearable
                        @clear="clearCountry"
                        style="width: 150px"
                    >

                        <el-option v-for="item in BelongingCountryList" :label="$t(item.label)" :value="item.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('状态')">

                    <el-select
                        v-model="param.Status"
                        :placeholder="$t('请选择商户状态')"
                        style="width: 150px"
                    >

                        <el-option v-for="item in operatorStatus" :label="$t(item.label)" :value="item.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('创建时间')">

                    <el-date-picker
                        v-model="timeRange"
                        type="datetimerange"
                        :start-placeholder="$t('开始时间')"
                        :end-placeholder="$t('结束时间')"
                        format="YYYY-MM-DD HH:mm:ss"
                        date-format="MMM DD, YYYY"
                        time-format="HH:mm"
                    />
                </el-form-item>

                <el-form-item>
                    <el-button type="primary" @click="onSearchNotify">{{ $t('查询') }}</el-button>
                    <el-button style="margin-left: 10px" @click="resetData">{{ $t('重置') }}</el-button>
                </el-form-item>
            </el-form>
        </div>

        <div class="page_table_context">

            <table-children-component
                ref="tableChildRef"
                table-name="operatorMaintenance_list"
                v-loading="loading"
                row-key="AppID"
                :load-child="loadChildren"
                :table-data="tableData"
                :table-header="tableHeader"
                :page="param.PageIndex"
                :page-size="param.PageSize"
                :count="param.Count"
                @page-change="pageChange"
                @refreshTable="onSearchNotify"
            >

                <template #handleTools>
                    <el-button type="warning" plain @click="addDialog = true" v-if="store.AdminInfo.GroupId <= 1  || (store.AdminInfo.GroupId == 2  && store.AdminInfo.Businesses.ReviewType != 2)">
                        {{ $t('新增') }}
                    </el-button>
                </template>

                <template #AppID="scope">
                    <span :style="{fontWeight: scope.scope.OperatorType == 1 ? 'bolder' : ''}">
                        {{ scope.scope.AppID}}
                    </span>

                </template>

                <template #BalanceHeader>

                    <el-tooltip
                        class="box-item"
                        effect="dark"
                        :open-delay="300"
                        :hide-after="300"
                        placement="top-start"
                    >

                        <template #content>
                            <div>
                                <div>{{ $t('商户的钱包余额') }}</div>
                                <div><span style="color: green">{{ $t('绿色') }}</span>：{{ $t('预存余额正常') }}</div>
                                <div><span style="color: blue">{{ $t('蓝色') }}</span>：{{ $t('余额低于阈值') }}</div>
                                <div><span style="color: red">{{ $t('红色') }}</span>：{{ $t('余额不足') }}</div>

                            </div>
                        </template>
                        <el-icon size="20">
                            <QuestionFilled/>
                        </el-icon>
                    </el-tooltip>
                </template>


                <template #Status="scope">
                    <template v-if="scope.scope.ReviewStatus != 0">
                        <el-text type="danger" v-if="scope.scope.Status == 0">{{ $t('禁用') }}</el-text>
                        <el-text type="success" v-if="scope.scope.Status == 1">{{ $t('开启') }}</el-text>
                    </template>
                    <template v-else>
                        <el-text type="warning">{{ $t('未审批') }}</el-text>
                    </template>

                </template>
                <template #Balance="scope">
                    <template v-if="scope.scope.OperatorType == 2">
                        <el-text type="danger" v-if="scope.scope.Balance < 0">{{ ut.toNumberWithComma(scope.scope.Balance * 1e4) }}</el-text>
                        <el-text type="success" v-else-if="scope.scope.Balance > scope.scope.BalanceThreshold">{{ ut.toNumberWithComma(scope.scope.Balance * 1e4) }}</el-text>
                        <el-text type="primary" v-else-if="scope.scope.Balance <= scope.scope.BalanceThreshold">{{ ut.toNumberWithComma(scope.scope.Balance * 1e4) }}</el-text>
                        <el-text type="danger" v-else>0.00</el-text>
                    </template>
                    <template v-else>
                        /
                    </template>

                </template>
                <template #CreateTime="scope">
                    {{  ut.fmtDateSecond(scope.scope.CreateTime) == "/" ? "/" :  ut.fmtDateSecond(ut.fmtSelectedUTCDateFormat(scope.scope.CreateTime)) }}
<!--                    {{  ut.fmtDateSecond(scope.scope.CreateTime) }}-->
                </template>
                <template #TokenExpireAt="scope">

                    {{  ut.fmtDateSecond(scope.scope.TokenExpireAt) == "/" ? "/" :  ut.fmtDateSecond(ut.fmtSelectedUTCDateFormat(scope.scope.TokenExpireAt)) }}
<!--                    {{  ut.fmtDateSecond(scope.scope.TokenExpireAt) }}-->
                </template>

                <template #OperatorType="scope">
                        <span v-for="(item, index) in operatorType" :key="index">

                            <text v-if="scope.scope.OperatorType == item.value">{{ $t(item.label) }}</text>
                        </span>

                </template>

                <template #ThisMonthProfit="scope">
                    <template v-if="scope.scope.Name != 'admin'">
                        <el-text :type="scope.scope.ThisMonthProfit < 0 ? 'danger' : 'success'">
                            {{ ut.toNumberWithComma(scope.scope.ThisMonthProfit) }}
                        </el-text>
                    </template>
                    <template v-else>
                        /
                    </template>

                </template>

                <template #LastMonthProfit="scope">
                    <template v-if="scope.scope.Name != 'admin'">

                        <el-text :type="scope.scope.LastMonthProfit < 0 ? 'danger' : 'success'">
                            {{ ut.toNumberWithComma(scope.scope.LastMonthProfit) }}
                        </el-text>
                    </template>
                    <template v-else>
                        /
                    </template>

                </template>

                <template #operator="scope">
                    <div style="width: 100%; display: flex;justify-content: center">
                        <el-button size="small" type="primary" plain @click="checkOperator(scope.scope)" style="margin-bottom: 3px">
                            {{ $t('信息') }}
                        </el-button>
                        <el-button size="small" type="danger" plain @click="openGameSetting(scope.scope)" style="margin-bottom: 3px"
                                   :style="{visibility: scope.scope.OperatorType == 1 || store.AdminInfo.GroupId == 2 ? 'hidden' : 'visible'}">{{ $t('设置') }}
                        </el-button>
                    </div>


                </template>
            </table-children-component>

        </div>

        <!--  商户添加      -->
        <OperatorDialog v-model="addDialog"
                        :operatorType="operatorType"
                        :walletOptions="walletOptions"
                        :operatorStatus="operatorStatus"
                        :incomeType="incomeType"
                        :currencyRate="currencyRate"
                        :belongingCountryList="BelongingCountryList"
                        @update:modelValue="getList"
                        @addOperator="addOperator"

        />


        <!--   基础信息     -->
        <EditOperatorDialog v-model="editDialog"
                            :operatorData="editData"
                            :operatorType="operatorType"
                            :walletOptions="walletOptions"
                            :operatorStatus="operatorStatus"
                            :incomeType="incomeType"
                            :currencyRate="currencyRate"
                            :belongingCountryList="BelongingCountryList"
                            @editData="editOperator"/>


        <!--  游戏设置      -->
        <GameSetting v-model="gameSettingDialog"
                     ref="gameSettingRef"
                     :gameSettingData="gameSettingData"
                     @openSettingInfo="openGameSettingDialog"
        />


        <!--   对应游戏信息     -->
        <GameSettingInfo v-model="gameSettingInfoDialog"
                         :gameSettingInfoData="gameSettingInfoData"
                         @update:model-value="editGameChange"
        />


    </div>
</template>


<script lang='ts' setup>
import {ref, reactive, onMounted, nextTick} from 'vue'
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import {DatetimeFormat, useI18n} from 'vue-i18n';

import customTable from "@/components/customTable/tableComponent.vue"
import {AdminInfo} from "@/api/adminpb/info";

import {saveAs} from 'file-saver';
import OperatorDialog from "@/pages/systemManagement/component/operatorDialog.vue";
import EditOperatorDialog from "@/pages/systemManagement/component/editOperatorDialog.vue";
import GameSetting from "@/pages/systemManagement/component/gameSettingDialog.vue";
import GameSettingInfo from "@/pages/systemManagement/component/gameSettingInfo.vue";
import {stat} from "fs";
import TableChildrenComponent from "@/components/customTable/tableChildrenComponent.vue";
import Operator_container from "@/components/operator_container.vue";
import {COUNTRY_LIST} from "@/lib/config";
import ut, {PasswordRSAEncryption} from "@/lib/util";
import {ElMessageBox} from "element-plus";
import copy from "copy-to-clipboard";
import {QuestionFilled} from "@element-plus/icons-vue";
import {OperatorRTP} from "@/enum";


let {t} = useI18n()
const store = useStore()

// Client.send("/mq/AdminInfo/Interior/operatorInfo", {AppID:"hml789"})

let tableHeaderUp = [

    {
        label: "商户AppID",
        value: "AppID",
        type: "custom",
        width: "150px"
    },
    {
        label: "商户序号",
        value: "Id",
        align: "center",
        width: "110px"
    },
    {
        label: "商户编码",
        value: "AppCode",
        align:"center",
        width: "330px"
    },
    {
        label: "商户账号",
        value: "UserName",
        align:"center",
        width: "110px"
    },
    {
        label: "商户类型",
        value: "OperatorType",
        type: "custom",
        align:"center",
        width: "110px"
    },
    {
        label: "归属国家",
        value: "BelongingCountry",
        align:"center",
        width: "110px",
        format:(row)=> {
            let item = BelongingCountryList.value.find(item => item.value == row.BelongingCountry)

            if (item && item.label){

                return t(item.label)
            }else{
                return ""
            }
        }
    },
    {
        label: "币种",
        value: "CurrencyName",
        align:"center",
        width: "110px",

    },
    {
        label: "钱包类型",
        value: "WalletMode",
        width: "110px",
        format: (row) => {
            return ["", t("转账钱包"), t("单一钱包")][row.WalletMode]
        },
        align:"center"

    },
    {
        label: "合作模式",
        value: "WalletMode",
        width: "110px",
        format: (row) => {
            return row.OperatorType == 2 ? ["/", t("收益分成"), t("流水分成")][row.CooperationType] : "/"
        },
        align:"center"

    },
    {
        label: "余额",
        value: "Balance",
        type: "custom",
        width: "120px",
        headerCustom: true,
        headerName: "BalanceHeader",
        tips:true,
        align:"center"
    },
    {
        label: "状态",
        value: "Status",
        type: "custom",
        align:"center"
    }
]
let tableHeaderDown = [
    {
        label: "创建时间",
        value: "CreateTime",
        type: "custom",
        align:"center",
        width: "200px"
    },
    {
        label: "最后登录时间",
        value: "TokenExpireAt",
        type: "custom",
        align:"center",
        width: "200px"
    },
    {
        label: "备注",
        value: "Remark",
        align:"center",
        width: "220px"
    },
    {
        label: "操作",
        value: "operator",
        type: "custom",
        width: "200px",
        align:"center",
        fixed:"right",
        hiddenVisible: true
    }
]
let tableChildRef = ref(null)
let tableHeader = ref([


])

const operatorContainerRef = ref(null)
let tableData = ref([])
let param = ref({
    BelongingCountry: "ALL",
    OperatorType: 0,
    Name: "",
    Status: -1,
    CreatedTime: [],
    createdStartTime: 0,
    createdEndTime: 0,
    AppID: "",
    Id: "",
    ReviewStatus:0,
    PageIndex: 1,
    PageSize: 20,
    Count: 0
})

let gameSettingRef = ref(null)

let currencyRate = ref(["商户处理", "POP官方处理"])
let operatorType = [
    {label:"全部", value: 0},
    {label:"线路商", value: 1},
    {label:"商户", value: 2}
]
let operatorStatus = [
    {label: "全部", value: -1},
    {label: "开启", value: 1},
    {label: "禁用", value: 0},
    {label: "未审批", value: 2}
]

let BelongingCountryList = ref(COUNTRY_LIST)


let incomeType = [
    {label: "收益分成", value: 1},
    {label: "流水分成", value: 2},
]


const walletOptions = [
    {
        label: '单一钱包',
        value: 2
    },
    {
        label: '转账钱包',
        value: 1
    }
]


if(store.language != 'zh'){
    BelongingCountryList.value.forEach(item=>{
        item.label = item.value
    })
}

let loading = ref(false)

// 添加面板是否显示
let addDialog = ref(false)
let timeRange = ref([])

// 修改面板相关信息
let editDialog = ref(false)
let editData = ref({})

// 游戏设置
// 针对商户游戏设置信息
let gameSettingDialog = ref(false)
let gameSettingData = ref({})

// 针对游戏设置信息
let gameSettingInfoDialog = ref(false)
let gameSettingInfoData = ref({})


const defaultOperatorEvent = ({})
const operatorListChange = (value) => {
    if (value){
        param.value.AppID = value.AppID
    }else{
        param.value.AppID = ""
    }
}

const getList = async () => {
    tableData.value = []
    const searchData = {...param.value}


    searchData.createdStartTime = 0
    searchData.createdEndTime = 0
    if (timeRange.value && timeRange.value.length) {


        if (timeRange.value[0] && timeRange.value[1]) {


            let start = ut.fmtUTCDate(timeRange.value[0].getTime()) * 1000
            let end = ut.fmtUTCDate(timeRange.value[1].getTime()) * 1000


            searchData.createdStartTime = ut.fmtSelectedUTCDate(start, "reduce") * 1000
            searchData.createdEndTime = ut.fmtSelectedUTCDate(end, "reduce") * 1000

        }
    }

    if (searchData.Status == 2){
        searchData.ReviewStatus = -1
        searchData.Status = -1
    }


    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.NewGetOperatorList, searchData as any)
    loading.value = false
    if (err) {
        return
    }
    param.value.Count = data?.AllCount
    if (data.List){
        data.List.forEach(item => {
            item.hasChildren = item.HasChildren
        })

        tableData.value = data.List
    }
}

onMounted(()=>{
    tableHeader.value.push(...tableHeaderUp)


    if (store.AdminInfo.GroupId != 3){


        tableHeader.value.push(...[
            {
                label: "线路商上月收益",
                value: "LastMonthProfit",
                width: "150px",
                type: "custom",
                align:"center"
            },
            {
                label: "线路商本月收益",
                value: "ThisMonthProfit",
                width: "150px",
                type:"custom",
                align:"center"
            }
        ])
    }
    tableHeader.value.push(...tableHeaderDown)
})



const loadChildren = (AppID) => {
    return Client.Do(AdminInfo.GetOperatorChildList, {AppID: AppID} as any)
}

const clearCountry = () => {
    param.value.BelongingCountry = "ALL"
}
const onSearchNotify = () => {
    tableData.value = []
    getList()
}

getList()


const resetData = () => {
    param.value = {
        OperatorType: 0,
        Name: "",
        Status: -1,
        CreatedTime: [],
        createdStartTime: 0,
        createdEndTime: 0,
        AppId: "",
        PageIndex: 1,
        PageSize: 20,
        Count: 0
    }
    operatorContainerRef.value.paramData = null
    timeRange.value = []
    getList()
}

let xiazai = async () => {
    // let [data, err] = await Client.Do(AdminGroup.DownloadOperatorData, {})
    // if (err) {
    //     return tip.e(err)
    // }
    // let datas = data.List.map(t => {
    //     return {
    //         ...t,
    //         Url: window.location.href
    //     }
    // })
    //
    // const row = datas[0]
    //
    // const txt = JSON.stringify(datas, null, 2)
    // let strData = new Blob([txt], {type: 'text/plain;charset=utf-8'});
    // saveAs(strData, `商户信息-${row.AppId}.txt`)
    // let str = '王佳伟Vue字符串保存到txt文件下载到电脑案例'
    // let strData = new Blob([str], { type: 'text/plain;charset=utf-8' });
    saveAs([], "测试文件下载.txt");

}
let checkOperator = (row) => {
    editDialog.value = true
    editData.value = row
}

let openGameSetting = (row) => {
    gameSettingDialog.value = true

    gameSettingData.value = row
}


let openGameSettingDialog = (data) => {
    gameSettingInfoDialog.value = true
    gameSettingInfoData.value = {
        ...data,
    }
}


const editGameChange = () => {
    gameSettingRef.value.getGameList()
}

const pageChange = (page) => {
    param.value.PageIndex = page.currentPage
    param.value.PageSize = page.dataSize
    getList()
}


const addOperator = (data) => {




    data.data.WalletModeText = ["", t("转账钱包"), t("单一钱包")][data.data.WalletMode]
    data.data.CooperationTypeText = ["", t("收益分成"), t("流水分成")][data.data.CooperationType]

    data.data.ContactList = ""
    let CopyStr = ""

    let Contact = JSON.parse(data.data.Contact)
    for (let i in Contact){
        data.data.ContactList += `<div style="width: 100%"><span style="width: 30%;display: inline-block; border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;margin-left: 5px;margin-bottom: 8px">${Contact[i].name}</span><span style="width: 60%;display: inline-block; border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;margin-left: 5px;margin-bottom: 8px">${Contact[i].value}</span></div>`
        CopyStr += `                  ${t('联系名称')}:${Contact[i].name}        ${t('联系信息')}：${Contact[i].value}\n`
    }

    let copyJson = `${t("商户AppID")}: ${data.data.AppID}
${t("钱包类型")}: ${t(data.data.WalletModeText)}
${t("商户账号")}: ${data.data.UserName}
${t('账号') + t('密码') }: ${data.data.Password}
${t("商户币种")}: ${data.data.CurrencyKey}${data.data.OperatorType == 2 ? `\n${t("基础货币单位")}: ${t(currencyRate.value[data.data.CurrencyCtrlStatus])}` : ``}
${t("合作模式")}: ${data.data.CooperationTypeText}
${data.data.CooperationType == 1 ? `${t("平台费率")}: ${data.data.PlatformPay}%` : `${t("流水比例")}: ${data.data.TurnoverPay}%`}
${store.AdminInfo.GroupId <= 1 ? `${t("RTP区间")}: ${OperatorRTP[data.data.HighRTPOff].split("-")[0]}%-${OperatorRTP[data.data.HighRTPOff].split("-")[1]}%}` : ``}
${t("联系方式")}:        \n${CopyStr}
    `

    ElMessageBox.confirm(
        `
                <div style="width: 630px">
                <span style="width: 100%;text-align: center;display: inline-block;font-weight: bold;font-size: 16px">${t('商户信息')}</span>
                    <div style="display: flex;align-items: center;justify-content: center;flex-wrap: wrap; width: 100%">


                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">
                            <span>${t('商户AppID') + ":"}</span>
                            <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">${data.data.AppID}</div>
                        </div>
                        <div style="width: 50%;margin-top:5px;margin-bottom:5px;visibility: ${ data.data.OperatorType == 1 ? 'hidden' : 'visible'}">
                            <span>${t('钱包类型') + ":"}</span>
                            <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">${data.data.WalletModeText}</div>
                        </div>
                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">
                            <span>${t('商户账号') + ":"}</span>
                            <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">${data.data.UserName}</div>
                        </div>
                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">
                            <span>${t('账号') + t('密码') + ":"}</span>
                            <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">${data.data.Password}</div>
                        </div>
                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">
                            <span>${t('商户币种') + ":"}</span>
                            <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">${data.data.CurrencyKey}</div>
                        </div>
                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">
                            ${data.data.OperatorType == 2 ?
                                `<span>${t('基础货币单位') + ":"}</span>
                                    <div style = "width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px" >
                                     ${currencyRate.value[data.data.CurrencyCtrlStatus]}
                                    </div>`
                                :  ''}
                        </div>

                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">


                         ${data.data.OperatorType == 2 ?
                                    `<span>${t('合作模式') + ":"}</span>
                                    <div style = "width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px" >
                                        ${t(data.data.CooperationTypeText)}
                                    </div>`
                                    : `<span>${t('平台费率') + ":"}</span>
                                    <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">
                                        ${data.data.PlatformPay}%
                                   </div>`}
                         </div>

                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">

                              ${data.data.CooperationType == 2 ||  data.data.OperatorType == 1 ?
                                    `<span>${t('流水比例') + ":"}</span>
                                    <div style = "width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px" >
                                        ${data.data.TurnoverPay}%
                                    </div>`
                                 : `<span>${t('平台费率') + ":"}</span>
                                    <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">
                                        ${data.data.PlatformPay}%
                                   </div>`}
                        </div>


                        ${store.AdminInfo.GroupId <= 1 ?
                        `<div style="width: 50%;margin-top:5px;margin-bottom:5px">
                            <span>${t('RTP区间') + ":"}</span>
                            <div style="width: 80%;border:1px solid #e4e7ed;background: #e1e1e1;padding: 5px 8px;border-radius: 5px;min-height:30px">${OperatorRTP[data.data.HighRTPOff].split("-")[0]}%-${OperatorRTP[data.data.HighRTPOff].split("-")[1]}%</div>
                        </div>
                        <div style="width: 50%;margin-top:5px;margin-bottom:5px">
                        </div>` : ''}


                    </div>
                        <div style="width: 100%;margin-top:5px;margin-bottom:5px; display: flex">
                            <span>${t('联系方式') + ":"}</span>
                            <div style="width: 70%;display: flex;align-items: center;flex-wrap: wrap">${data.data.ContactList}</div>
                        </div>
                </div>`,
        t('创建成功'),
        {
            customClass:"createSucee",
            dangerouslyUseHTMLString: true,
            confirmButtonText: t('复制'),
            cancelButtonText: t('关闭'),
        }
    )
        .then(async () => {
            copy(copyJson)
            tip.s(t("复制成功"))
        })
}

const editOperator = async (data) => {


    if (data.key == "admin"){
        getList()
    }else{

        loading.value = true
        const [a, err] =  await Client.Do(AdminInfo.GetOperatorChildList, {AppID: data.key} as any)
        loading.value = false
        tableChildRef.value.tableRef.updateKeyChildren(data.key, a.List)
    }

}


</script>
<style scoped lang='scss'>
.flex_child_end {
    margin-bottom: 15px;
}

.demo-form-inline > .el-form-item {
    margin-bottom: 15px;
}
</style>


<style>

</style>
