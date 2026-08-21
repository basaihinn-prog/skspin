<template>


    <el-dialog
            v-model="dialogVisible"
            :title="$t('玩家数据')"
            :width="store.viewModel === 2 ? '85%' : '950px'"
            @close="closeDialog"
            @open-auto-focus="openDialog"

    >
        <div v-loading="loading">

            <div class="playerInfo">
<!--                <span>{{ $t('玩家名称') }}：</span>{{ playData.UName }}-->
                <span>{{ $t('唯一标识') }}：</span>{{ playData.Pid }}<span>{{ $t('账号') }}：</span>{{ playData.Uid }}</div>

            <div class="playerInfoContent">
                <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
                    <el-tab-pane :label="$t('会员信息')" name="info">


                        <el-row>
                            <el-col :span="14">


                                <el-form label-width="125px">
                                    <el-row>
                                        <el-col :span="12">
                                            <el-form-item :label="`${$t('账号状态')}:`" label-width="110px">
                                                <el-tag
                                                        :type="playData.Status == 0 ? 'danger': 'success'">{{
                                                    playData.Status == 0 ? $t('关闭') : $t('开启')
                                                    }}
                                                </el-tag>

                                            </el-form-item>


                                        </el-col>
                                        <el-col :span="12">
                                            <el-form-item :label="`${$t('调控游戏')}:`">
                                                <div style="width: 100%">
                                                    <el-switch
                                                            :style="{float: 'right'}"
                                                            :disabled="!(store.AdminInfo.GroupId <= 1 || (store.AdminInfo.GroupId == 3 && store.AdminInfo.Businesses.PlayerRTPSettingOff == 1))"
                                                            v-model="userStatusSwitch"
                                                            @click="userStatusChange"
                                                            :active-value="1"
                                                            :inactive-value="0"

                                                    />
                                                </div>
                                            </el-form-item>

                                        </el-col>
                                        <el-col :span="24">
                                            <el-form-item :label="`${$t('在线状态')}:`" style="width: 50%"
                                                          label-width="110px">
                                                <el-tag :type="!playData.OnlineStatus ? 'info': 'success'">{{
                                                    !playData.OnlineStatus ? $t('离线') : $t('在线')
                                                    }}
                                                </el-tag>

                                            </el-form-item>

                                        </el-col>

                                        <el-col :span="12">
                                            <el-form-item :label="$t('所属商户') + ':'" label-width="110px">
                                                <el-input disabled v-model="playData.AppID"></el-input>
                                            </el-form-item>
                                            <el-form-item :label="`${$t('所属线路商')}:`"
                                                          v-if="store.AdminInfo.GroupId <= 3" label-width="110px">
                                                <el-input disabled v-model="playData.onlineAppID"></el-input>
                                            </el-form-item>
                                        </el-col>
                                        <el-col :span="12">
                                        </el-col>

                                        <el-col :span="24">
                                            <el-form-item :label="`${$t('总输赢')}:`" style="width: 50%"
                                                          label-width="110px">
                                                <el-input disabled v-model="playData.totalWin"></el-input>
                                            </el-form-item>

                                        </el-col>

                                        <el-col :span="12">
                                            <el-form-item :label="`${$t('注册时间')}:`" label-width="110px">
                                                <el-input disabled v-model="playData.CreateAt"></el-input>
                                            </el-form-item>

                                        </el-col>
                                        <el-col :span="12">
                                        </el-col>
                                        <el-col :span="12">
                                            <el-form-item :label="`${$t('最后登录时间')}:`"  label-width="110px">
                                                <el-input disabled v-model="playData.LoginAt"></el-input>
                                            </el-form-item>

                                        </el-col>
                                        <el-col :span="12">
                                        </el-col>
                                        <!--                                    -->
                                        <!--                                    <el-col :span="24">-->
                                        <!--                                        <el-form-item :label="`${$t('注册时间')}:`" style="width: 50%" label-width="100px">-->
                                        <!--                                            <el-input disabled v-model="playData.CreateAt"></el-input>-->
                                        <!--                                        </el-form-item>-->

                                        <!--                                    </el-col>-->
                                        <!--                                    <el-col :span="24">-->
                                        <!--                                        <el-form-item :label="`${$t('最后登录时间')}:`" style="width: 50%" label-width="100px">-->
                                        <!--                                            <el-input disabled v-model="playData.LoginAt"></el-input>-->
                                        <!--                                        </el-form-item>-->

                                        <!--                                    </el-col>-->

                                    </el-row>
                                </el-form>

                            </el-col>
                            <el-col :span="10" v-if="userStatusSwitch">

                                <PlayerControl ref="PlayerControlRef" class="playerDialogControl"
                                               :SelectPlayInfo="playData" :GameList="playData.GameID"
                                               :PlayerRTPDisabled="!(store.AdminInfo.GroupId <= 1 || (store.AdminInfo.GroupId == 3 && store.AdminInfo.Businesses.PlayerRTPSettingOff == 1))"
                                               :PlayerBuyRTPDisabled="(store.AdminInfo.GroupId <= 1 || (store.AdminInfo.GroupId == 3 && store.AdminInfo.Businesses.BuyRTPOff == 1))"></PlayerControl>
                            </el-col>
                        </el-row>

                    </el-tab-pane>
                    <!--                <el-tab-pane label="账户交易" name="payment">-->

                    <!--                    <el-descriptions :column="2">-->
                    <!--                        <el-descriptions-item :label="item.label" v-for="(item, index) in playUserInfoAll" :key="index"></el-descriptions-item>-->
                    <!--                    </el-descriptions>-->

                    <!--                    <customTable :table-header="paymentHeader" :table-data="paymentData"></customTable>-->

                    <!--                </el-tab-pane>-->
                    <!--                <el-tab-pane label="投注统计" name="bet">-->

                    <!--                    <el-descriptions :column="2">-->
                    <!--                        <el-descriptions-item :label="item.label" v-for="(item, index) in playUserInfoAll" :key="index"></el-descriptions-item>-->
                    <!--                    </el-descriptions>-->

                    <!--                    <customTable :table-header="betHeader" :table-data="betData"></customTable>-->

                    <!--                </el-tab-pane>-->
                    <!--                <el-tab-pane label="登入日志" name="login">-->

                    <!--                    <el-descriptions :column="2">-->
                    <!--                        <el-descriptions-item :label="item.label" v-for="(item, index) in playUserInfoAll" :key="index"></el-descriptions-item>-->
                    <!--                    </el-descriptions>-->

                    <!--                    <customTable :table-header="loginHeader" :table-data="loginData"></customTable>-->

                    <!--                </el-tab-pane>-->
                    <!--                <el-tab-pane label="后台变更日志" name="change">-->

                    <!--                    <el-descriptions :column="2">-->
                    <!--                        <el-descriptions-item :label="item.label" v-for="(item, index) in playUserInfoAll" :key="index"></el-descriptions-item>-->
                    <!--                    </el-descriptions>-->

                    <!--                    <customTable :table-header="changeHeader" :table-data="changeData"></customTable>-->

                    <!--                </el-tab-pane>-->
                </el-tabs>
            </div>
        </div>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="closeDialog">{{ $t('关闭') }}</el-button>
                <el-button type="primary" style="margin-left: 10px" @click="commitRTPControl" v-if="store.AdminInfo.GroupId <= 1 || (store.AdminInfo.GroupId == 3 && store.AdminInfo.Businesses.PlayerRTPSettingOff == 1)">
                    {{ $t('提交') }}
                </el-button>
            </div>


        </template>
    </el-dialog>

</template>

<script setup lang="ts">

import {computed, Ref, ref} from "vue";
import {Client} from "@/lib/client";
import {AdminPlayer, PlayerResponse} from "@/api/adminpb/adminPlayer";
import ut from "@/lib/util";
import {useI18n} from "vue-i18n";
import {RTP_VALUE} from "@/lib/RTP_config";
import PlayerControl from "@/pages/analysis/component/PlayerControl.vue";
import {useStore} from "@/pinia";
import {tip} from "@/lib/tip";
import {ElMessageBox} from "element-plus";
import {AdminInfo} from "@/api/adminpb/info";

const {t} = useI18n()


let store = useStore()
const props = defineProps(["modelValue", "playerId"])
const emits = defineEmits(["update:modelValue", "update:playerId"])

const activeName = ref("info")
const PlayerControlRef = ref(null)
const playData: Ref<PlayerResponse> = ref<PlayerResponse>({
    Pid: 0,
    Uid: "",
    GameID: "",
    AppID: "",
    CurrencyKey: "",
    CurrencyName: "",
    RTPControlStatus: 0,
    BuyRTP: 95,
    LoginAt: "",
    CreateAt: "",
    TypeInfo: [],
    PlayerRTPControllist:[],
    Bet: 0,
    Win: 0,
})
const huibao = ref(0)
const loading = ref(true)

const dialogVisible = computed(() => {

    if (props.modelValue) {
        getPlayerInfo()
    } else {
        if (PlayerControlRef.value && PlayerControlRef.value.init) {
            PlayerControlRef.value.init()
        }


        emits('update:playerId', "")
    }

    return props.modelValue
})

const userStatusSwitch = ref(0)
const userStatusChange = async () => {

}

const openDialog = () => {
    playUserInfo.value = [
        {label: "账号状态：", value: "Status", span: 1, type: "custom"},
        {label: "调控游戏：", value: "control", span: 1, type: "custom"},
        {label: "在线状态：", value: "OnlineStatus", span: 2, type: "custom"},
        // {label: "风控状态：", value: "", span: 2, type: ""},
        {label: "所属商户：", value: "AppID", span: 1, type: ""},
        {label: "所属商户余额：", value: "PlatformPay", span: 1, type: ""},
        {label: "所属线路商：", value: "onlineAppID", span: 1, type: ""},
        {label: "所属线路商余额：", value: "onlinePlatformPay", span: 1, type: ""},
        {label: "总输赢：", value: "totalWin", span: 2, type: ""},
        {label: "注册时间：", value: "CreateAt", span: 2, type: ""},
        {label: "最后登录时间：", value: "LoginAt", span: 1, type: ""},
    ]
}


let paymentData = ref([])
let paymentHeader = [
    {label: "账户类型", value: ""},
    {label: "币种", value: ""},
    {label: "交易前金额", value: ""},
    {label: "交易金额", value: ""},
    {label: "交易后金额", value: ""},
    {label: "交易时间", value: ""},
]
let betData = ref([])
let betHeader = [
    {label: "订单编号", value: ""},
    {label: "游戏名称", value: ""},
    {label: "投注时间", value: ""},
    {label: "输赢", value: ""},
    {label: "投注钱金额", value: ""},
    {label: "投注后金额", value: ""},
    {label: "投注时间", value: ""},
]
let loginData = ref([])
let loginHeader = [
    {label: "最后登录时间", value: ""},
    {label: "本次登录累计投注", value: ""},
]
let changeData = ref([])
let changeHeader = [
    {label: "时间", value: ""},
    {label: "操作状态", value: ""},
    {label: "操作人员", value: ""},
]


let playUserInfo = ref([
    {label: "账号状态：", value: "Status", span: 1, type: "custom"},
    {label: "调控游戏：", value: "control", span: 1, type: "custom"},
    {label: "在线状态：", value: "OnlineStatus", span: 2, type: "custom"},
    {label: "风控状态：", value: "", span: 2, type: ""},
    {label: "所属商户：", value: "AppID", span: 1, type: ""},
    {label: "所属商户余额：", value: "PlatformPay", span: 1, type: ""},
    {label: "所属线路商：", value: "onlineAppID", span: 1, type: ""},
    {label: "所属线路商余额：", value: "onlinePlatformPay", span: 1, type: ""},
    {label: "总输赢：", value: "totalWin", span: 2, type: ""},
    {label: "注册时间：", value: "CreateAt", span: 1, type: ""},
    {label: "最后登录时间：", value: "LoginAt", span: 1, type: ""},
])
let playUserInfoAll = [
    {label: "账号状态：", value: ""},
    {label: "所属商户：", value: ""},
    {label: "所属线路商：", value: ""},
    {label: "在线状态：", value: ""},
    {label: "风控状态：", value: ""},
    {label: "可用余额：", value: ""},
]

const handleClick = (tab) => {
    activeName.value = tab
}

const getPlayerInfo = async () => {
    loading.value = true
    let [response, err] = await Client.Do(AdminPlayer.GetPlayerInfo, {Pid: props.playerId})


    let gameList = []
    let gameListIds = []
    if (!response) {
        tip.e(t("未找到玩家信息"))
        loading.value = false
        emits("update:modelValue")
        return
    }


    if (response.PlayerRTPControllist && response.PlayerRTPControllist.length) {
        for (const index in response.PlayerRTPControllist) {
            gameList.push({
                ID: response.PlayerRTPControllist[index].GameID,
                Name: response.PlayerRTPControllist[index].GameName
            })
            gameListIds.push(response.PlayerRTPControllist[index].GameID)
        }

    }


    if (store.AdminInfo.GroupId == 3) {
        playUserInfo.value = playUserInfo.value.filter(item => item.value != "onlineAppID" && item.value != "onlinePlatformPay")
    }

    playData.value = {...response}

    let [merchant] = await Client.Do(AdminInfo.GetOperatorInfo, {AppID: response.AppID} as any)


    userStatusSwitch.value = response.RTPControlStatus ? 1 : 0
    playData.value.CreateAt = ut.fmtSelectedUTCDateFormat(response.CreateAt)
    playData.value.LoginAt = ut.fmtSelectedUTCDateFormat(response.LoginAt)
    playData.value.Operator = merchant.OperatorInfo
    playData.value.PlatformPay = response.Operator.Advance
    playData.value.onlineAppID = response.OnlineOperator.AppID
    playData.value.onlinePlatformPay = response.OnlineOperator.Advance
    playData.value.totalWin = (Number(response.Win) - Number(response.Bet)) / 10000
    playData.value.gameList = gameList
    playData.value.gameListIds = gameListIds
    playData.value.ContrllRTP = response.PlayerRTPControllist ? response.PlayerRTPControllist[0].ControlRTP : 93
    playData.value.AutoRemoveRTP = response.PlayerRTPControllist ? response.PlayerRTPControllist[0].AutoRemoveRTP : 0
    playData.value.PersonWinMaxScore = response.PlayerRTPControllist ? response.PlayerRTPControllist[0].PersonWinMaxScore : 0
    playData.value.PersonWinMaxMult = response.PlayerRTPControllist ? response.PlayerRTPControllist[0].PersonWinMaxMult : 0
    playData.value.huibao = isNaN(response.Win / response.Bet) ? 0 : response.Win / response.Bet
    playData.value.BuyRTP = response.PlayerRTPControllist ? response.PlayerRTPControllist[0].BuyRTP : 90


    huibao.value = isNaN(response.Win / response.Bet) ? 0 : response.Win / response.Bet
    loading.value = false
}



const commitRTPControl = async () => {


    if (userStatusSwitch.value == 1) {
        ElMessageBox.confirm(
            t('确认对用户进行RTP设置'),
            t('是否确认'),
            {
                confirmButtonText: t('确定'),
                cancelButtonText: t('关闭'),
                type: 'warning',
            }
        ).then(async () => {
            const tag = await PlayerControlRef.value.RTPConfig([playData.value])
            if (tag) {
                emits('update:modelValue', false)
            }
        })
    } else {

        ElMessageBox.confirm(
            t('确认关闭该用户的游戏调控'),
            t('是否确认'),
            {
                confirmButtonText: t('确定'),
                cancelButtonText: t('关闭'),
                type: 'warning',
            }
        )
            .then(async () => {
                const [response, err] = await Client.Do(AdminPlayer.ClosePlayerRTP, {Pid: playData.value.Pid})
                if (!err) {
                    tip.s(t("关闭RTP调控成功"))
                    emits('update:modelValue', false)
                } else {
                    tip.e(t("关闭RTP调控失败"))
                }
            })
    }


}


const closeDialog = () => {
    emits('update:modelValue')
}
</script>


<style scoped lang="scss">
.playerInfo {
  font-size: 16px;
  color: #737373;
}

.playerInfo span {
  font-weight: bolder;
  color: #000000;
  margin-left: 10px;
  display: inline-block;
}

.RTPControlGameList {
  width: 100%;
  height: 150px;

}
</style>
<style>

.playerInfoContent .el-descriptions__body {
    box-shadow: none !important;
}

.playerDialogControl .el-form-item__label {
    color: var(--el-text-color-primary);
}

.playerDialogControl .el-space__item {
    color: var(--el-text-color-regular);
}

.el-dialog__body .el-form .el-form-item {
    margin-bottom: 18px;
}
</style>
