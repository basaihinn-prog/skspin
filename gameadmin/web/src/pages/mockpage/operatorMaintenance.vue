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
                    <el-button type="primary">{{ $t('查询') }}</el-button>
                    <el-button style="margin-left: 10px">{{ $t('重置') }}</el-button>
                </el-form-item>
            </el-form>
        </div>

        <div class="page_table_context">

            <table-children-component
                    table-name="operatorMaintenance_list"
                    v-loading="loading"
                    @page-change="pageChange"
                    :table-data="tableData"
                    :table-header="tableHeader"
                    :page="param.PageIndex"
                    :page-size="param.PageSize"
                    :count="param.Count"

            >

                <template #handleTools>
                    <el-space wrap>
                        <el-button type="warning" plain
                                   v-if="store.AdminInfo.GroupId <= 1  || (store.AdminInfo.GroupId == 2  && store.AdminInfo.Businesses.ReviewType != 2)">
                            {{ $t('新增') }}
                        </el-button>

                        <upload-excel @uploadFile="resolveExcel" file-name="商户列表"/>
                    </el-space>

                </template>

                <template #AppID="scope">
                    <span :style="{fontWeight: scope.scope.OperatorType == 1 ? 'bolder' : ''}">
                        {{ scope.scope.AppID }}
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
                        <el-text type="danger" v-if="scope.scope.Status == '禁用'">{{ $t('禁用') }}</el-text>
                        <el-text type="success" v-if="scope.scope.Status == '开启'">{{ $t('开启') }}</el-text>
                    </template>
                    <template v-else>
                        <el-text type="warning">{{ $t('未审批') }}</el-text>
                    </template>

                </template>

                <template #Balance="scope">

                        <el-text type="danger" v-if="scope.scope.Balance < 0">
                            {{ ut.toNumberWithComma(scope.scope.Balance * 1e4) }}
                        </el-text>

                        <el-text type="success" v-else-if="scope.scope.Balance >= 0">
                            {{ ut.toNumberWithComma(scope.scope.Balance * 1e4) }}
                        </el-text>
                        <el-text type="danger" v-else>0.00</el-text>



                </template>

                <template #CreateTime="scope">

                    {{ dateFormater(0, 0, scope.scope.CreateTime) }}
                </template>

                <template #TokenExpireAt="scope">
                    {{ dateFormater(0, 0, scope.scope.TokenExpireAt) }}
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
                        <el-button size="small" type="primary" plain style="margin-bottom: 3px">
                            {{ $t('信息') }}
                        </el-button>
                        <el-button size="small" type="danger" plain style="margin-bottom: 3px"
                                   :style="{visibility: scope.scope.OperatorType == 1 || store.AdminInfo.GroupId == 2 ? 'hidden' : 'visible'}">
                            {{ $t('设置') }}
                        </el-button>
                    </div>


                </template>
            </table-children-component>

        </div>


    </div>
</template>


<script lang='ts' setup>
import {ref, onMounted} from 'vue'
import {Client} from '@/lib/client';
import {useStore} from '@/pinia';
import {useI18n} from 'vue-i18n';

import {AdminInfo} from "@/api/adminpb/info";
import TableChildrenComponent from "@/components/customTable/tableChildrenComponent.vue";
import Operator_container from "@/components/operator_container.vue";
import {COUNTRY_LIST} from "@/lib/config";
import ut from "@/lib/util";

import {QuestionFilled} from "@element-plus/icons-vue";
import UploadExcel from "./component/uploadExcel.vue";


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
        label: "商户ID",
        value: "Id",
        align: "center",
        width: "110px"
    },
    {
        label: "商户编码",
        value: "AppCode",
        align: "center",
        width: "330px"
    },
    {
        label: "商户账号",
        value: "UserName",
        align: "center",
        width: "110px"
    },
    {
        label: "商户类型",
        value: "OperatorType",
        align: "center",
        width: "110px"
    },
    {
        label: "归属国家",
        value: "BelongingCountry",
        align: "center",
        width: "110px",
    },
    {
        label: "币种",
        value: "CurrencyName",
        align: "center",
        width: "110px",

    },
    {
        label: "钱包类型",
        value: "WalletMode",
        width: "110px",
        align: "center"

    },
    {
        label: "合作模式",
        value: "CooperationType",
        width: "110px",
        align: "center"

    },
    {
        label: "余额",
        value: "Balance",
        type: "custom",
        width: "120px",
        headerCustom: true,
        headerName: "BalanceHeader",
        tips: true,
        align: "center"
    },
    {
        label: "状态",
        value: "Status",
        type: "custom",
        align: "center"
    }
]
let tableHeaderDown = [
    {
        label: "创建时间",
        value: "CreateTime",
        type: "custom",
        align: "center",
        width: "200px"
    },
    {
        label: "最后登录时间",
        value: "TokenExpireAt",
        type: "custom",
        align: "center",
        width: "200px"
    },
    {
        label: "备注",
        value: "Remark",
        align: "center",
        width: "220px"
    },
    {
        label: "操作",
        value: "operator",
        type: "custom",
        width: "200px",
        align: "center",
        fixed: "right",
        hiddenVisible: true
    }
]
let tableHeader = ref([])

let tableHeaderMap = ref(null)

const operatorContainerRef = ref(null)
let tableData = ref([])
let resolveData = ref([])
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
    ReviewStatus: 0,
    PageIndex: 1,
    PageSize: 20,
    Count: 0
})

let operatorType = [
    {label: "全部", value: 0},
    {label: "线路商", value: 1},
    {label: "商户", value: 2}
]
let operatorStatus = [
    {label: "全部", value: -1},
    {label: "开启", value: 1},
    {label: "禁用", value: 0},
    {label: "未审批", value: 2}
]

let BelongingCountryList = ref(COUNTRY_LIST)


if (store.language != 'zh') {
    BelongingCountryList.value.forEach(item => {
        item.label = item.value
    })
}

let loading = ref(false)

// 添加面板是否显示
let addDialog = ref(false)
let timeRange = ref([])


const defaultOperatorEvent = ({})
const operatorListChange = (value) => {
    if (value) {
        param.value.AppID = value.AppID
    } else {
        param.value.AppID = ""
    }
}

const getList = async () => {

}
const pageChange = (page) => {
    param.value.PageIndex = page.currentPage
    param.value.PageSize = page.dataSize
    let start = (param.value.PageIndex - 1) * param.value.PageSize
    let end = param.value.PageSize * param.value.PageIndex
    tableData.value = resolveData.value.slice(start, end)
}


const resolveExcel = (data) => {


    if (!tableHeaderMap.value) {
        tableHeaderMap.value = {}
        for (const i in tableHeader.value) {
            tableHeaderMap.value[tableHeader.value[i].label] = tableHeader.value[i].value
        }
    }


    let generatorData = []
    for (const i in data.data) {

        let dataItem = data.data[i]
        let generatorItem = {}

        for (const key in dataItem) {
            generatorItem[tableHeaderMap.value[key]] = dataItem[key]
        }
        generatorData.push(generatorItem)

    }
    param.value.Count = generatorData.length
    resolveData.value = generatorData
    param.value.PageIndex = 1
    tableData.value = resolveData.value.slice(0, param.value.PageSize)
}

    onMounted(() => {
        tableHeader.value.push(...tableHeaderUp)


        if (store.AdminInfo.GroupId != 3) {


            tableHeader.value.push(...[
                {
                    label: "线路商上月收益",
                    value: "LastMonthProfit",
                    width: "150px",
                    type: "custom",
                    align: "center"
                },
                {
                    label: "线路商本月收益",
                    value: "ThisMonthProfit",
                    width: "150px",
                    type: "custom",
                    align: "center"
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
