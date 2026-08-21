<template>
    <div class='userDetails' v-loading="loading">
        <!-- 搜索 -->
        <div class="searchView">

            <el-form
                    :model="param"
                    @keyup.enter="getUserDetails"
                    style="max-width: 100%"
            >
                <el-space wrap>
                    <el-form-item :label="$t('商户')">
                        <el-select v-model="param.OperatorId" filterable clearable
                                   :disabled="store.AdminInfo.GroupId == 3"
                                   :placeholder="$t('请选择')"
                                   style="min-width: 147px">
                            <el-option v-for='item in operatorData' :key="item.Id" :label="item.AppID"
                                       :value="item.Id"/>
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="param.Pid" clearable
                                  oninput="value=value.replace(/[^\d]/g,'')"
                                  :placeholder="$t('请输入')"/>
                    </el-form-item>
                    <el-form-item :label="$t('玩家ID')">
                        <el-input v-model.trim="param.Uid" clearable :placeholder="$t('请输入')"/>
                    </el-form-item>

                    <el-form-item>
                        <el-button type="primary" @click="getUserDetails">{{ $t('搜索') }}</el-button>
                    </el-form-item>

                </el-space>
            </el-form>

        </div>


        <div class="page_table_context">
            <div style="width: 98%; margin: 0 auto">
                <template v-if="selectUserData.length>1">
                    <el-radio-group v-model="selectUserModel" v-for="(item,index) in selectUserData"
                                    @change="selectUserChange">
                        <el-radio-button :label="item.Pid">
                            {{ $t('玩家ID') + '：' + item.Uid + '（' + (Number(index + 1)) + '）' }}
                        </el-radio-button>
                    </el-radio-group>
                </template>

                <el-row :gutter="12">
                    <el-col :lg="12" :xs="24">
                        <div class="descriptionsView">
                            <el-descriptions :title="$t('玩家信息')" :column="1" :border="true" size="small" >


                                <el-descriptions-item maxlength="12" :label="$t('唯一标识')">
                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                            <table-tips :tips="tipsMap['唯一标识']"/>
                                            {{$t('唯一标识')}}
                                        </div>

                                    </template>
                                    {{ userData.Pid }}
                                    <el-tag class="copy" @click="copyText(userData.Pid)">{{ $t('复制') }}</el-tag>
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('玩家ID')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['玩家ID']"/>
                                        {{$t('玩家ID')}}
                                        </div>
                                    </template>

                                    <div style="width: 330px">{{ userData.Uid }}</div>
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('所属商户')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['所属商户']"/>
                                        {{$t('所属商户')}}
                                        </div>
                                    </template>
                                    {{ userData.AppID }}</el-descriptions-item>
<!--                                <el-descriptions-item :label="$t('玩家名称')">{{ userData.UName }}</el-descriptions-item>-->
                                <el-descriptions-item :label="$t('封禁状态')" v-if="store.AdminInfo.GroupId != 2">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['封禁状态']"/>
                                        {{$t('封禁状态')}}
                                        </div>
                                    </template>

                                    <el-switch
                                            v-model="userStatusSwitch"
                                            @click="userStatusChange"
                                            size="small"
                                            :disabled="!userData.Uid"
                                            :active-value="0"
                                            :inactive-value="1"
                                    />
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('总投注')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['总投注']"/>
                                        {{$t('总投注')}}
                                        </div>
                                    </template>

                                    {{ ut.toNumberWithComma(userData.Bet) }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('总赢分')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['总赢分']"/>
                                        {{$t('总赢分')}}
                                        </div>
                                    </template>

                                    {{ ut.toNumberWithComma(userData.Win) }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('总盈亏')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['总盈亏']"/>
                                        {{$t('总盈亏')}}
                                        </div>
                                    </template>

                                    {{ ut.toNumberWithComma(userData.Win - userData.Bet) }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('回报率')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['回报率']"/>
                                        {{$t('回报率')}}
                                        </div>
                                    </template>

                                    {{ percentFormatter(0, 0, userData.Win / userData.Bet || 0) }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('最后登录时间')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['最后登录时间']"/>
                                        {{$t('最后登录时间')}}
                                        </div>
                                    </template>
                                    {{
                                        userData.AppID ? ut.fmtDateSecond(ut.fmtSelectedUTCDateFormat(userData.LoginAt)) : "/"
                                    }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('注册时间')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['注册时间']"/>
                                        {{$t('注册时间')}}
                                        </div>
                                    </template>
                                    {{
                                        userData.AppID ? ut.fmtDateSecond(ut.fmtSelectedUTCDateFormat(userData.CreateAt)) : "/"
                                    }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('未结算金额')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['未结算金额']"/>
                                        {{$t('未结算金额')}}
                                        </div>
                                    </template>

                                    {{ ut.toNumberWithCommaNormal(userData.Unsettled) }}
                                </el-descriptions-item>

                                <el-descriptions-item v-if="currentUser.AppID == 'admin'" :label="$t('未转移金额（当前余额）')">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['未转移金额（当前余额）']"/>
                                        {{$t('未转移金额（当前余额）')}}
                                        </div>
                                    </template>

                                    <el-space wrap>
                                        <el-input v-model.trim="userData.Balance" maxlength="15" style="width:200px"
                                                  :placeholder="$t('请设置未转移金额')"/>
                                        <div class="table-icon" @click="setBalance" v-if="updateMoneyShow">
                                            <el-button type="primary" plain size="small">{{ $t('提交') }}</el-button>
                                        </div>
                                        <div class="table-icon" @click="setHistory(3,$t('未转移金额'))"
                                             v-if="setHistoryShow">
                                            <el-button type="danger" plain size="small">{{ $t('记录') }}</el-button>
                                        </div>
                                    </el-space>
                                </el-descriptions-item>
                                <el-descriptions-item v-if="currentUser.AppID == 'admin'">

                                    <template #label>

                                        <div class="flex" style="align-items: center;color:#000000">

                                        <table-tips :tips="tipsMap['Slots个人奖池']"/>
                                        {{$t('Slots个人奖池')}}
                                        </div>

                                    </template>

                                    <el-space wrap>
                                        <el-input v-model.trim="userData.SlotsPool" style="width:200px" maxlength="15"
                                                  :placeholder="$t('请输入')"/>
                                        <div class="table-icon" @click="setPlayerPool_Slot" v-if="updateMoneyShow">
                                            <el-button type="primary" plain size="small">{{ $t("提交") }}</el-button>
                                        </div>
                                        <div class="table-icon" @click="setHistory(1,$t('Slot奖池'))"
                                             v-if="setHistoryShow">
                                            <el-button type="danger" plain size="small">{{ $t("记录") }}</el-button>
                                        </div>
                                    </el-space>
                                </el-descriptions-item>

                                <!--            百人奖池                -->

                                <!--
                                <el-descriptions-item v-if="currentUser.PermissionId == 1" :label="$t('百人奖池')">
                                    <el-space wrap>
                                        <el-input v-model.trim="userData.BaiRenPool" :placeholder="$t('请输入')"/>
                                        <div class="table-icon" @click="setPlayerPool_BaiRen">
                                            <el-icon><Select/></el-icon>
                                        </div>
                                        <div class="table-icon" @click="setHistory(2,$t('百人奖池'))" v-if="setHistoryShow">
                                            <el-icon>
                                                <Tickets/>
                                            </el-icon>
                                        </div>
                                    </el-space>
                                </el-descriptions-item>
                                -->
                            </el-descriptions>
                        </div>
                    </el-col>
                    <el-col :lg="12" :xs="24" :style="{visibility:setHistoryCard ? 'visible' : 'hidden' }">


                        <customTable
                            table-name="playerInfoOperator_list"
                            :table-data="activities.list"
                            :table-header="tableHeader"
                            @refresh-table="setHistory(textref == $t('Slot奖池') ? 1 : 3 ,textref)"
                            height="480"
                        >
                            <template #handleTools>
                                <el-button @click="setHistoryCard = false" size="small" style="float: right">{{ $t('关闭') }}
                                </el-button>

                            </template>
                            <template #GameRTP="scope">
                                <el-tag effect="plain">
                                    {{
                                        scope.scope.Change > 0 ? '+ ' + ut.toNumberWithComma(scope.scope.Change) : ut.toNumberWithComma(scope.scope.Change)
                                    }}
                                </el-tag>
                            </template>
                        </customTable>


                    </el-col>
                </el-row>

                <template v-if="TypeInfo && TypeInfo[0]">
                    <el-descriptions
                        v-for="(info, index) in [TypeInfo[0]]"
                        :key="index"
                        :title="$t(gameType[0])"
                        direction="vertical"
                        :column="5"
                        border
                    >
                        <el-descriptions-item>


                            <template #label>
                                <div class="flex" style="align-items: center;color:#000000">

                                    <table-tips :tips="tipsMap['普通旋转（投注）金额']"/>
                                    {{$t('普通旋转（投注）金额')}}
                                </div>

                            </template>

                            <el-tag type="warning" effect="dark" round>
                                {{ ut.toNumberWithComma(info.Bet) }}
                            </el-tag>
                        </el-descriptions-item>
                        <el-descriptions-item>




                            <template #label>
                                <div class="flex" style="align-items: center;color:#000000">


                                    <table-tips :tips="tipsMap['普通旋转（投注）次数']"/>
                                    {{$t('普通旋转（投注）次数')}}
                                </div>
                            </template>


                            {{ info.SpinCnt }}
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('普通旋转（投注）平均数')">
                            <template #label>
                                <div class="flex" style="align-items: center;color:#000000">


                                    <table-tips :tips="tipsMap['普通旋转（投注）平均数']"/>
                                    {{$t('普通旋转（投注）平均数')}}
                                </div>
                            </template>


                            {{
                                (info.AvgBet / 10000) || 0.00
                            }}
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('普通旋转（投注）赢分金额')">
                            <template #label>
                                <div class="flex" style="align-items: center;color:#000000">


                                    <table-tips :tips="tipsMap['普通旋转（投注）赢分金额']"/>
                                    {{$t('普通旋转（投注）赢分金额')}}
                                </div>
                            </template>


                            <el-tag type="warning" effect="dark" round>
                                {{ ut.toNumberWithComma(info.Win) }}
                            </el-tag>
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('购买小游戏次数')">
                            <template #label>
                                <div class="flex" style="align-items: center;color:#000000">


                                    <table-tips :tips="tipsMap['购买小游戏次数']"/>
                                    {{$t('购买小游戏次数')}}
                                </div>
                            </template>

                            {{ info.BuyGame }}
                        </el-descriptions-item>
                    </el-descriptions>
                </template>

                <template v-if="TypeInfo && TypeInfo[1]">

                    <div style="display: flex">

                        <div style="flex: 3">

                            <el-descriptions
                                v-for="(info, index) in [TypeInfo[1]]"
                                :key="index"
                                :title="$t(gameType[1])"
                                direction="vertical"
                                :column="4"
                                border
                                style="margin-top: 15px"
                            >
                                <el-descriptions-item>


                                    <template #label>
                                        <div class="flex" style="align-items: center;color:#000000">

                                            <table-tips :tips="tipsMap['普通投注金额']"/>
                                            {{$t('普通投注金额')}}
                                        </div>

                                    </template>

                                    <el-tag type="warning" effect="dark" round>
                                        {{ ut.toNumberWithComma(info.Bet) }}
                                    </el-tag>
                                </el-descriptions-item>
                                <el-descriptions-item>




                                    <template #label>
                                        <div class="flex" style="align-items: center;color:#000000">


                                            <table-tips :tips="tipsMap['普通投注次数']"/>
                                            {{$t('普通投注次数')}}
                                        </div>
                                    </template>


                                    {{ info.SpinCnt }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('普通投注平均数')">
                                    <template #label>
                                        <div class="flex" style="align-items: center;color:#000000">


                                            <table-tips :tips="tipsMap['普通投注平均数']"/>
                                            {{$t('普通投注平均数')}}
                                        </div>
                                    </template>


                                    {{
                                        (info.AvgBet / 10000) || 0.00
                                    }}
                                </el-descriptions-item>
                                <el-descriptions-item :label="$t('普通投注赢分金额')">
                                    <template #label>
                                        <div class="flex" style="align-items: center;color:#000000">


                                            <table-tips :tips="tipsMap['普通投注赢分金额']"/>
                                            {{$t('普通投注赢分金额')}}
                                        </div>
                                    </template>


                                    <el-tag type="warning" effect="dark" round>
                                        {{ ut.toNumberWithComma(info.Win) }}
                                    </el-tag>
                                </el-descriptions-item>

                            </el-descriptions>
                        </div>

                        <div style="flex: 0.48">
                        </div>
                    </div>
                </template>

                <el-descriptions v-if="plrLotteryInfoList"
                                 direction="vertical"
                                 :column="5"
                                 :title="$t('彩票')"
                                 border
                >
                    <template #label>
                        <div class="flex" style="align-items: center;color:#000000">


                            <table-tips :tips="tipsMap['彩票']"/>
                            {{$t('彩票')}}
                        </div>
                    </template>
                    <el-descriptions-item :label="$t('投注金额')">
                        <template #label>
                            <div class="flex" style="align-items: center;color:#000000">


                                <table-tips :tips="tipsMap['投注金额']"/>
                                {{$t('投注金额')}}
                            </div>
                        </template>
                        {{ ut.toNumberWithComma(plrLotteryInfoList?.Bet) }}
                    </el-descriptions-item>
                    <el-descriptions-item :label="$t('中奖')">

                        <template #label>
                            <div class="flex" style="align-items: center;color:#000000">


                                <table-tips :tips="tipsMap['中奖']"/>
                                {{$t('中奖')}}
                            </div>
                        </template>
                        {{
                        ut.toNumberWithComma(plrLotteryInfoList?.Win)
                        }}
                    </el-descriptions-item>
                    <el-descriptions-item :label="$t('盈亏')">
                        <template #label>
                            <div class="flex" style="align-items: center;color:#000000">


                                <table-tips :tips="tipsMap['盈亏']"/>
                                {{$t('盈亏')}}
                            </div>
                        </template>
                        <el-tag type="warning" effect="dark" round>
                            {{ ut.toNumberWithComma(plrLotteryInfoList?.ProfitAndLoss) }}
                        </el-tag>
                    </el-descriptions-item>
                </el-descriptions>
            </div>
        </div>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, getCurrentInstance, watch, watchEffect, shallowRef} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia/index';
import {AdminGameCenter, PlayerInfoReq} from '@/api/gamepb/admin';
import {AdminSlotsRpc,} from '@/api/slots/admin_rpc';
import {PoolType} from '@/api/slots/pool';
import {useRoute, useRouter} from 'vue-router';
import {AdminInfo} from "@/api/adminpb/info";
const store = useStore()
const {proxy} = getCurrentInstance()
import {ElMessage} from 'element-plus'
import ut from "@/lib/util";
import {useI18n} from "vue-i18n";
import {Switch} from "@element-plus/icons-vue";
import copy from "copy-to-clipboard";
import customTable from "@/components/customTable/tableComponent.vue"
import TableTips from "@/components/customTable/tableTips.vue";

let currentUser = store.AdminInfo


const route = useRoute();
let operatorData = ref([])
let userData = ref({})
let userStatusSwitch = ref(0)
let selectUserData = ref([])
let selectUserModel = ref(null)
let loading = ref(false)
let TypeInfo = ref(null)
let {t} = useI18n()
const tableHeader = [
    {label: "时间", value: "time", format: (row) => ut.fmtDateSecond(row.time), width: "150px"},
    {label: "操作人", value: "OpName", width: "80px"},
    {label: "调整前数据", value: "Pid", format: (row)=> ut.toNumberWithComma(row.OldGold)},
    {label: "变化值", value: "Change", format: (row)=> ut.toNumberWithComma(row.Change)},
    {label: "调整后数据", value: "NewGold", format:(row)=> ut.toNumberWithComma(row.NewGold)},
]
const activities = ref(
    {list: null}
)
const textref = ref('')
const setHistoryCard = ref(false)
const setHistoryShow = ref(false)
const updateMoneyShow = ref(true)
const gameType = [
    'SLOTS游戏', 'MINI游戏', '', '百人游戏', '彩票'
]
let param = ref({
    Pid: null,
    Uid: "",
    OperatorId: null
})
let tipsMap = ref({})
let operatorParam = reactive({
    PageIndex: 1,
    PageSize: 10000,
    OperatorType: 2,
    Status: 1
})
const userStatus = ref(0)
let plrLotteryInfoList = ref(null)
onMounted(() => {
    operatorList()

    if (store.AdminInfo.GroupId == 3) {
        param.value.OperatorId = store.AdminInfo.BusinessesId
    }

        console.log(store.tipsMap)
        tipsMap.value = store.tipsMap[route.meta.title]

});

const operatorList = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, operatorParam)
    if (err) {
        return tip.e(err)
    }
    // operatorData.value = data.AllCount === 0 ? [] : data.List.filter(list => !list.Status)
    operatorData.value = data.AllCount === 0 ? [] : data.List

}
const selectUserChange = async () => {
    userData.value = selectUserData.value.find(m => m.Pid === selectUserModel.value)
    userStatus.value = userData.value['Status']
    param.value.Pid = userData.value['Pid']
    param.value.Uid = userData.value['Uid']

    await Promise.all([
        //获取Slots的奖池
        getUserPool(PoolType.Slots, "SlotsPool"),
        //获取百人的奖池
        getUserPool(PoolType.BaiRen, "BaiRenPool"),
        getPlayerBalance(),
        // plrLotteryInfo()
    ])
    userStatusSwitch.value = userData.value['Status']
    setHistoryShow.value = true
}
const userStatusChange = async () => {

    if (!userData.value['Uid']) {
        return
    }

    let [data, err] = await Client.Do(AdminInfo.GetUpdatePlayerStatus, {
        Uid: userData.value['Uid'],
        Pid: userData.value['Pid'],
        AppID: userData.value['AppID'],
        Status: Number(userStatusSwitch.value)
    })
    if (err) {
        return tip.e(err)
    }
    tip.s(t(`${t('玩家')}${userData.value['Pid']}${t('状态已修改')}`))

}

const getUserDetails = async () => {
    setHistoryShow.value = false
    if (!param.value.Pid && !param.value.Uid) {
        return tip.e(t('玩家ID和唯一标识必填其中一项'))
    }
    param.value.Pid = param.value.Pid || null;
    loading.value = true
    try {
        let [data, err] = await Client.Do(AdminGameCenter.NewAppList, param.value)
        if (err) return tip.e(err)
        if (!data.List) {
            tip.e(t("未找到玩家"))
            userData.value = {}
            TypeInfo.value = null
            return
        }


        const userDataNew = data.List.length === 1 ? data.List[0] : {};
        const selectUserDataNew = data.List.length > 1 ? data.List : [];

        userData.value = userDataNew;

        let [operatorData, operatorErr] = await Client.Do(AdminInfo.getOperatorByAppId, {AppID: userData.value.AppID})

        if (operatorErr) {

            tip.e(t(operatorErr))
            userData.value = {}
            return
        }

        if (operatorData.WalletMode == 2) {
            updateMoneyShow.value = false

        }


        selectUserData.value = selectUserDataNew;
        param.value.Pid = userData.value['Pid'] || null
        TypeInfo.value = userData.value['TypeInfo'] ? Object.entries(userData.value['TypeInfo']).map(([type, data]) => ({
            ...data,
            type
        })) : null;
        if (data.List.length === 1) {
            await Promise.all([
                getUserPool(PoolType.Slots, "SlotsPool"),
                getUserPool(PoolType.BaiRen, "BaiRenPool"),
                getPlayerBalance(),
                // plrLotteryInfo()


            ])
            userStatusSwitch.value = data.List[0]['Status']

            setHistoryShow.value = true
        }
    } catch (e) {
        tip.e(e);
    } finally {
        loading.value = false;
    }
}
let BeforeBalance = 0
const getPlayerBalance = async () => {
    if (!param.value.Pid) {
        if (selectUserData.value.length > 0) {
            if (!selectUserModel.value) {
                ElMessage.error(t('请选择玩家ID'))
                return
            } else {
                param.value.Pid = selectUserModel.value
            }
        }
        param.value.Pid = userData.value['Pid']
    }
    let [data, err] = await Client.Do(AdminGameCenter.GetPlayerBalance, {Pid: param.value.Pid})
    if (err) return tip.e(err)
    let list = {
        Balance: data.Balance,
        Unsettled: data.Unsettled
    }
    BeforeBalance = data.Balance
    Object.assign(userData.value, list);
}
const plrLotteryInfo = async () => {
    if (!param.value.Pid) {
        return
    }
    let [data, err] = await Client.Do(AdminGameCenter.GetPlrLotteryInfo, {Pid: param.value.Pid})
    if (err) return tip.e(err)
    plrLotteryInfoList.value = data
}
const hasPlusGreaterThan0 = (value) => value > 0 ? (/[+]/.test(value) ? value : '+' + value) : value;
const getUserPool = async (type, key) => {
    // 获取奖池
    let datas = {
        ...param.value,
        type,
        Pid: Number(param.value.Pid),
        Gold: null
    }
    datas.Gold = userData.value[key] === 0 ? '0.00' : datas?.Gold / 10000
    // datas.Gold = formatGold(userData.value[key]);
    try {
        let [data, err] = await Client.Do(AdminSlotsRpc.GetSelfSlotsPool, datas)
        if (err) return tip.e(err)
        // userData.value[key] = data?.Gold === 0?'0.00':data?.Gold/10000
        userData.value[key] = data?.Gold === 0 ? '0.00' : data?.Gold / 10000
        userData.value['SlotsPool'] = hasPlusGreaterThan0(userData.value['SlotsPool'])
        userData.value['BaiRenPool'] = hasPlusGreaterThan0(userData.value['BaiRenPool'])
    } catch (e) {
        tip.e(e || t('获取奖池失败'));
    }

}

const copyText = (e) => {

    copy(e)
    tip.s(t("复制成功"))
}

const setUserPool = async (type, key) => {
    if (!userData.value.Pid) {
        return tip.e(t("未找到玩家"));
    }
    try {
        const userDataValue = userData.value;
        const pid = Number(userDataValue['Pid']);
        const gold = Math.floor(Number(userDataValue[key]) * 10000);
        if (isNaN(pid) || isNaN(gold)) {
            return tip.e(t("Invalid data type"));
        }
        let [_, err] = await Client.Do(AdminSlotsRpc.SetSelfSlotsPool, {
            ...userDataValue,
            Type: type,
            Pid: pid,
            Gold: gold
        })
        if (err) return tip.e(err)
        tip.s(t("设置成功"))
        await getUserDetails()
    } catch (e) {
        tip.e(e.message);
    }
}

const setPlayerPool_BaiRen = async () => await setUserPool(PoolType.BaiRen, 'BaiRenPool')
const setPlayerPool_Slot = async () => await setUserPool(PoolType.Slots, 'SlotsPool')

const setBalance = async () => {
    if (!userData.value.Pid) {
        return tip.e(t("未找到玩家"));
    }
    try {
        let [data] = await Client.Do(AdminGameCenter.GetPlayerBalance, {Pid: param.value.Pid})
        let [responseData, err] = await Client.Do(AdminSlotsRpc.setBalance, {
            Pid: userData.value['Pid'],
            Balance: Number(userData.value['Balance']),
            BeforeBalance: data.Balance
        })
        if (err) return tip.e(err)
        tip.s(t("设置成功"))
        await getUserDetails()
    } catch (e) {
        tip.e(e)
    }
}
const debounce = (func, delay, immediate = false) => {
    let timer;
    const debounced = function (...args) {
        const context = this;
        const later = () => {
            timer = null;
            if (!immediate) func.apply(context, args);
        };
        const callNow = immediate && !timer;
        clearTimeout(timer);
        timer = setTimeout(later, delay);
        if (callNow) func.apply(context, args);
    };
    debounced.cancel = () => clearTimeout(timer);
    return debounced;
};
const setHistory = debounce(async (type, name) => {
    setHistoryCard.value = true
    if (!setHistoryCard.value) {
        return
    }
    try {
        const pid = Number(userData.value['Pid'])
        let [responseData, err] = await Client.Do(AdminSlotsRpc.setHistory, {
            Pid: pid,
            Type: type,

        })
        if (err) return tip.e(err)
        activities.value.list = responseData.List || []

        activities.value.list.forEach(item=>{
            item.time = ut.fmtDateSecond(new Date(item.time).getTime() / 1000)
        })

        textref.value = name
    } catch (e) {
        tip.e(e)
        setHistoryCard.value = false
    } finally {
        console.log('');
    }
}, 200);


</script>
<style scoped lang='scss'>
.descriptionsView {
  width: 100%;
  display: flex;
  justify-content: flex-start;
  flex-direction: row;
  align-items: flex-start;

  .el-descriptions {
    width: 100%;
  }
}


.userDetails {
  width: 100%;
  min-height: calc(100vh - 200px);

  .el-descriptions {
    .labelItem {
      width: 130px;
    }

    .el-descriptions__body tr {
      height: 40px;
    }
  }
}
.el-descriptions__header {
  padding-top: 1rem;
}

.table-icon {

  margin-left: 5px;
  margin-right: 5px;
  font-size: 13px;
}
</style>



<style>


.el-descriptions__cell{
    .el-icon{
        font-size: 16px!important;
        color: #000000;
    }
}

</style>
