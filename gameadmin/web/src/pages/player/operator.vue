/*
*@Author: 西西米
*@Date: 2022-12-24 15:53:22
*@Description: 用户详情
*/
<template>
    <div class='userDetails' v-loading="loading">
        <!-- 搜索 -->
        <div class="searchView">
            <el-form
                :model="param"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.trim.number="param.Pid" clearable oninput="value=value.replace(/[^\d]/g,'')"
                                  :placeholder="$t('请输入')" />
                    </el-form-item>
                    <el-form-item :label="$t('玩家ID')">
                        <el-input v-model.trim="param.Uid" clearable :placeholder="$t('请输入')" />
                    </el-form-item>
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operatorInfo="operatorListChange"></operator_container>

                    <el-form-item>
                        <el-button type="primary" @click="getUserDetails">{{ $t('搜索') }}</el-button>
                    </el-form-item>

                </el-space>
            </el-form>

        </div>
        <template v-if="selectUserData.length>1">
            <el-radio-group v-model="selectUserModel" v-for="(item,index) in selectUserData" @change="selectUserChange">
                <el-radio-button  :label="item.Pid" >
                    {{$t('玩家ID')+'：'+item.Uid+'（'+ (Number(index+1)) +'）'}}
                </el-radio-button>
            </el-radio-group>
        </template>
        <el-row :gutter="10">
            <el-col :lg="6" :xs="24">
                <div class="descriptionsView" >
                    <el-descriptions :title="$t('玩家信息')" :column="1" :border="true" size="small">
                        <el-descriptions-item :label="$t('唯一标识')">{{ userData.Pid }}
                            <el-tag class="copy" @click="copyText(userData.Pid)">{{ $t('复制') }}</el-tag>
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('玩家ID')">{{ userData.Uid }}
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('商户')">{{ userData.AppID }}
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('玩家状态')" v-if="userData.Uid">
                            <el-switch
                                v-model="userStatus"
                                @click="userStatusChange"
                                size="small"
                                :active-value="1"
                                :inactive-value="0"
                                :active-text="$t('启用')"
                                :inactive-text="$t('禁用')"
                            />

                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('总投注')">{{ ut.toNumberWithComma(userData.Bet) }} </el-descriptions-item>
                        <el-descriptions-item :label="$t('总赢分')">{{ ut.toNumberWithComma(userData.Win) }} </el-descriptions-item>
                        <el-descriptions-item :label="$t('总盈亏')">{{ ut.toNumberWithComma(userData.Win-userData.Bet) }} </el-descriptions-item>
                        <el-descriptions-item :label="$t('回报率')">{{ percentFormatter(0, 0, userData.Win/userData.Bet || 0) }} </el-descriptions-item>

                        <el-descriptions-item :label="$t('最后登陆时间')">
                            {{ dateFormater(0, 0, userData.LoginAt) }}
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('注册时间')">
                            {{ dateFormater(0, 0, userData.CreateAt) }}
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('未结算金额')">{{ userData.Unsettled?.toFixed(2) }} </el-descriptions-item>
                        <el-descriptions-item :label="$t('未转移金额')">{{ userData.Balance?.toFixed(2) }} </el-descriptions-item>
                    </el-descriptions>

                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive, getCurrentInstance, watch, watchEffect } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
// import router from '@/router';
import { useStore } from '@/pinia/index';
import { AdminGameCenter, PlayerInfoReq } from '@/api/gamepb/admin';
import { AdminSlotsRpc, } from '@/api/slots/admin_rpc';
import { PoolType } from '@/api/slots/pool';
import ut from "@/lib/util";
import { useRoute, useRouter } from 'vue-router';
import {AdminInfo} from "@/api/adminpb/info";
import {ElMessage} from "element-plus";
import {useI18n} from "vue-i18n";
import Operator_container from "@/components/operator_container.vue";
const store = useStore()
const { proxy } = getCurrentInstance()
let param = reactive({
    Pid: null,
    Uid: "",
    OperatorId: null,
    OperatorName: null,
})
let operatorParam = reactive({
    PageIndex: 1,
    PageSize: 10000
})
const route = useRoute();
let userData = ref({})
let userStatus = ref(0)
let selectUserData = ref([])
let selectUserModel = ref(null)
let rhs = reactive({
    tabName: "tab1",
    rewardPool: 0,
})
let { t } = useI18n()
let loading = ref(false)

const selectUserChange = async () => {
    userData.value = selectUserData.value.find(m=>m.Pid === selectUserModel.value)
    // userStatus.value = userData.value['Status']
    param.Pid = userData.value.Pid
    param.Uid = userData.value.Uid
    await getUserPool(PoolType.Slots, "SlotsPool")
    await getUserPool(PoolType.BaiRen, "BaiRenPool")
    await getPlayerBalance()
}

const defaultOperatorEvent = ref({})
const operatorListChange = (value) => {
    if (value){

        param.OperatorId = value.Id
    }else{

        param.OperatorId = ""
    }
}

const getUserDetails = async () => {
    if (!param.Pid && !param.Uid) {
        return tip.e(t('玩家ID和唯一标识必填其中一项'))
    }
    if (param.Pid === '') {
        param.Pid = null
    }
    loading.value = true
    let [data, err] = await Client.Do(AdminGameCenter.NewAppList, param)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    if (!data.List) userData.value = {}
    if (data.List.length > 1) {
        selectUserData.value = data.List
        userData.value = {}
    } else {
        userData.value = data.List[0]
        selectUserData.value = []
    }
    if (data.List.length === 1) {

        try {
            await getUserPool(PoolType.Slots, "SlotsPool")
            await getUserPool(PoolType.BaiRen, "BaiRenPool")
            await getPlayerBalance()
        }catch (e) {

        }

        // userStatus.value = data.List[0]['Status']
    }
}

const userStatusChange = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetUpdatePlayerStatus, {
        Pid:userData.value['Pid'],
        Status:Number(userStatus.value)
    })
    if (err) {
        return tip.e(err)
    }
    tip.s(t(`玩家${userData.value['Pid']}状态已修改`))
    await getUserDetails()
}
const getPlayerBalance = async () => {
    if (!param.Pid) {
        if (selectUserData.value.length>0){
            if (!selectUserModel.value) {
                ElMessage.error(t('请选择玩家ID'))
                return
            } else {
                param.Pid = selectUserModel.value
            }
        }
        param.Pid = userData.value.Pid
    }
    let datas = {
        Pid:param.Pid
    }
    let [data, err] = await Client.Do(AdminGameCenter.GetPlayerBalance, datas)

    if (err) {
        return tip.e(err)
    }
    let list = {
        Balance:data.Balance,
        Unsettled:data.Unsettled
    }
    Object.assign(userData.value, list);
    console.log(userData.value,'data.Balance');
}
const getUserPool = async (type, key) => {
    console.log(userData.value,'userData');
    // 获取奖池
    let datas = {
        ...param,
        type,
        Gold:null
    }
    datas.Pid = Number(datas.Pid)
    datas.Gold = userData.value[key] === 0?'0.00':datas?.Gold/10000
    let [data, err] = await Client.Do(AdminSlotsRpc.GetSelfSlotsPool, datas)

    if (err) {
        return tip.e(err)
    }
    userData.value[key] = data?.Gold === 0?'0.00':data?.Gold/10000
}
const setUserPool = async (type,key) => {
    let arr = userData.value
    userData.value.Pid = Number(userData.value.Pid)
    let datas = {
        ...arr,
        type,
        Gold: Math.floor(Number(userData.value[key]) * 10000)
    }
    let [_, err] = await Client.Do(AdminSlotsRpc.SetSelfSlotsPool, datas)
    if (err) {
        return tip.e(err)
    }
    tip.s(t("设置成功"))
    getUserDetails()
}

const setPlayerPool_BaiRen = async () => {
    setUserPool(PoolType.BaiRen,'BaiRenPool')
}

const setPlayerPool_Slot = async () => {
    setUserPool(PoolType.Slots,'SlotsPool')
}




</script>
<style scoped lang='scss'>
.userDetails {
    min-height: calc(100vh - 200px);
}

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
    .el-descriptions {
        .labelItem {
            width: 130px;
        }

        .el-descriptions__body tr {
            height: 40px;
        }
    }
}
</style>
