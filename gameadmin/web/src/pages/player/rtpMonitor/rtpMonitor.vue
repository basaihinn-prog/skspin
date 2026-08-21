<template>
    <div>
        <div class="searchView">
            <el-form
                    :model="queryForm"
                    style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange" :hase-all="false" :is-init="true"></operator_container>


                    <el-form-item>
                        <el-button type="primary" @click="queryList" style="margin-bottom: 3px">{{
                            $t('搜索')
                            }}
                        </el-button>
                    </el-form-item>
                </el-space>
            </el-form>

        </div>


        <div class="page_table_context">

            <div class="tableTitle">

                <el-tooltip
                    class="box-item"
                    effect="dark"

                    placement="top"
                >

                    <template #content>

                        <div>{{$t("本功能可以降低slots游戏的波动性，使玩家RTP（在此机台的总赢钱数/总投注）更快的接近配置的机台RTP")}}</div>
                        <div>{{$t("玩家的RTP高于配置RTP时，有概率（可配置）给玩家触发强制不中奖")}}</div>
                        <div>{{$t("玩家的RTP低于配置RTP时，有概率（可配置）给玩家触发额外小奖")}}</div>
                        <div>{{$t("例如玩家此时RTP为60%，机台配置RTP93%。此时玩家RTP过低，则有概率给玩家触发额外小奖")}}</div>
                        <div>{{$t("对于新玩家，可以使用不同的配置来提高新玩家的数值体验")}}</div>
                        <div>{{$t("注意：请慎重配置增加RTP的功能，该功能可能会使玩家总体RTP短期内略高")}}</div>
                    </template>
                    <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
                {{ $t("机台RTP自动控制") }}</div>
            <customTable
                table-name="playerRTPData_list"
                v-loading="loading"
                height="140px"
                :table-header="tableHeader"
                :table-data="queryRTPData"
                :hideTableHandel="true"
                >


                <template #funcStatus="scope">
                    <el-text :type="scope.scope.MoniterConfig && scope.scope.MoniterConfig.IsMoniter == 1 ? 'primary' : 'danger'">
                        {{ scope.scope.MoniterConfig && scope.scope.MoniterConfig.IsMoniter == 1 ? $t('开启') : $t('关闭') }}
                    </el-text>
                </template>
                <template #funcConfig="scope">
                    <div style="display: flex;align-items: center;justify-content: center">
                        <el-icon style="cursor: pointer" @click="checkConfigInfo(scope.scope, 'MoniterConfig')"><View /></el-icon>
                    </div>
                </template>
                <template #Operator="scope">
                    <el-button :type="scope.scope.MoniterConfig && scope.scope.MoniterConfig.IsMoniter == 1 ? 'danger' : 'primary'"
                               plain @click="taggleMoniterConfig(scope.scope, 'MoniterConfig')" size="small">
                        {{ scope.scope.MoniterConfig && scope.scope.MoniterConfig.IsMoniter == 1 ? $t('关闭') : $t('开启') }}
                    </el-button>

                    <el-button type="warning" plain @click="openConfigDialog(scope.scope, 'MoniterConfig')" size="small">{{ $t("编辑") }}</el-button>
                </template>
            </customTable>

            <div class="tableTitle">
                <div style="width: 100%;text-align: center">

                    <el-tooltip
                        class="box-item"
                        effect="dark"

                        placement="top"
                    >

                        <template #content>


                            <div>{{$t("本功能可以降低slots游戏的波动性，使玩家RTP（在此机台的总赢钱数/总投注）更快的接近配置的玩家个人RTP")}}</div>
                            <div>{{$t("玩家的RTP高于配置RTP时，有概率（可配置）给玩家触发强制不中奖")}}</div>
                            <div>{{$t("玩家的RTP低于配置RTP时，有概率（可配置）给玩家触发额外小奖")}}</div>
                            <div>{{$t("例如玩家此时RTP为60%，个人配置RTP93%。此时玩家RTP过低，则有概率给玩家触发额外小奖")}}</div>
                            <div>{{$t("对于新玩家，可以使用不同的配置来提高新玩家的数值体验")}}</div>
                            <div>{{$t("注意：请慎重配置增加RTP的功能，该功能可能会使玩家总体RTP短期内略高")}}</div>

                        </template>
                        <el-icon><QuestionFilled /></el-icon>
                    </el-tooltip>

                    {{ $t("玩家RTP自动控制") }}</div>

                <div style="width: 100%;text-align: center">（{{ $t("需要在玩家数据中开启调控游戏才会生效") }}）</div>
            </div>
            <customTable
                table-name="playerRTPData_list"

                v-loading="loading"
                :table-header="tableHeader2"
                :table-data="queryRTPData"
                :hideTableHandel="true"
                height="140px">



                <template #funcStatus="scope">
                    <el-text :type="scope.scope.PersonalMoniterConfig && scope.scope.PersonalMoniterConfig.IsMoniter == 1 ? 'primary' : 'danger'">
                        {{ scope.scope.PersonalMoniterConfig && scope.scope.PersonalMoniterConfig.IsMoniter == 1 ? $t('开启') : $t('关闭') }}
                    </el-text>
                </template>
                <template #funcConfig="scope">
                    <div style="display: flex;align-items: center;justify-content: center">
                        <el-icon style="cursor: pointer" @click="checkConfigInfo(scope.scope, 'PersonalMoniterConfig')"><View /></el-icon>
                    </div>
                </template>
                <template #Operator="scope">

                    <el-button :type="scope.scope.PersonalMoniterConfig && scope.scope.PersonalMoniterConfig.IsMoniter == 1 ? 'danger' : 'primary'"
                               plain @click="taggleMoniterConfig(scope.scope, 'PersonalMoniterConfig')" size="small">
                        {{ scope.scope.PersonalMoniterConfig && scope.scope.PersonalMoniterConfig.IsMoniter == 1 ? $t('关闭') : $t('开启') }}
                    </el-button>
                    <el-button type="warning" plain @click="openConfigDialog(scope.scope, 'PersonalMoniterConfig')" size="small">{{ $t("编辑") }}</el-button>
                </template>
            </customTable>
        </div>

        <monitorDialog v-model="visibleRtpControl" :MonitorData="dialogData" :type="dialogType" @update:modelValue="closeDialog"/>
        <monitorInfoDialog v-model="visibleRtpInfoControl" :MonitorData="dialogData" :type="dialogType" @update:modelValue="closeDialog"/>
    </div>
</template>

<script setup lang='ts'>
import {Ref, ref} from 'vue';
import {useStore} from "@/pinia";
import {useI18n} from "vue-i18n";
import { PlayerRTPData, PlayerRTPListRequest} from "@/api/adminpb/adminPlayer";

import MonitorDialog from "@/pages/player/rtpMonitor/component/monitorDialog.vue";
import monitorInfoDialog from "@/pages/player/rtpMonitor/component/monitorInfoDialog.vue";
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {AdminInfo} from "@/api/adminpb/info";
import {Client} from "@/lib/client";
import {Monitor} from "@/api/adminpb/monitor";
import {tip} from "@/lib/tip";

const {t} = useI18n()
const store = useStore();
let loading = ref(true)


let tableHeader = ref([
    {label: "商户AppID", value: "AppID"},
    {label: "商户序号", value: "Id"},
    {label: "功能状态", value: "funcStatus", type:"custom", width: "140px"},
    // {label: "新玩家游戏局数", value: "MoniterNewbieNum", format: (row) => t("{Num}局", {Num: row.MoniterNewbieNum }), width: "140px"},
    // {label: "游戏内统计数据周期", value: "MoniterNumCycle", format: (row) => t("{Num}局", {Num: row.MoniterNumCycle })},
    {label: "新玩家游戏局数", value: "MoniterNewbieNum", format: (row) => row.MoniterConfig.MoniterNewbieNum, width: "140px"},
    {label: "游戏内统计数据周期", value: "MoniterNumCycle", format: (row) => row.MoniterConfig.MoniterNumCycle},
    {label: "玩家RTP误差范围", value: "MoniterRTPErrorValue", format: (row) => row.MoniterConfig.MoniterRTPErrorValue + "%", width:"140px"},
    {label: "功能自动控制详情", value: "funcConfig",type:"custom", width:"140px"},
    {label: "操作", value: "Operator", type:"custom" },
])


let tableHeader2 = ref([
    {label: "商户AppID", value: "AppID"},
    {label: "商户序号", value: "Id"},
    {label: "功能状态", value: "funcStatus", type:"custom", width: "140px"},
    // {label: "新玩家游戏局数", value: "MoniterNewbieNum", format: (row) => t("{Num}局", {Num: row.MoniterNewbieNum }), width: "140px"},
    // {label: "游戏内统计数据周期", value: "MoniterNumCycle", format: (row) => t("{Num}局", {Num: row.MoniterNumCycle })},
    {label: "新玩家游戏局数", value: "MoniterNewbieNum", format: (row) => row.PersonalMoniterConfig.MoniterNewbieNum, width: "140px"},
    {label: "游戏内统计数据周期", value: "MoniterNumCycle", format: (row) => row.PersonalMoniterConfig.MoniterNumCycle},
    {label: "玩家RTP误差范围", value: "MoniterRTPErrorValue", format: (row) => row.PersonalMoniterConfig.MoniterRTPErrorValue + "%", width:"140px"},
    {label: "功能自动控制详情", value: "funcConfig",type:"custom", width:"140px"},
    {label: "操作", value: "Operator", type:"custom" },
])


const visibleRtpControl = ref(false)
const visibleRtpInfoControl = ref(false)
const dialogData = ref({})
const dialogType = ref("")
const tableData = ref([])
const queryRTPData: Ref<PlayerRTPData[]> = ref<PlayerRTPData[]>([])
const queryForm: Ref<PlayerRTPListRequest> = ref<PlayerRTPListRequest>({
    AppID:"",
    OperatorId: 0,               // 商户ID
    Pid: "",                      // 用户登录账号
    ControlTimeStart: "",         // 控制时间开始
    ControlTimeEnd: "",           // 控制时间结束
    Manufacturer: "",           // 控制时间结束
    GameId: "",                  // 控制时间结束
    Page: 1,
    PageSize: 20
})


const Count = ref(0)
const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})



const operatorListChange = (value) => {

    loading.value = true
    const OperatorId = value ? value.Id : null
    queryForm.value.AppID = value.AppID
    if (queryRTPData.value.length > 0){
        queryRTPData.value = [value]
    }


    loading.value = false
}



const taggleMoniterConfig = async (row,key) => {


    let commitData = {
        ...row[key],
        MoniterType: key == 'MoniterConfig' ? 0 : 1,
        IsMoniter: row[key].IsMoniter ? 0 : 1,
        AppID: row.AppID
    }

    let [resp, err] = await Client.Do(Monitor.SetGameMonitor, commitData)
    if (err){
        return tip.e(t('修改失败'))
    }
    queryList()
}

const checkConfigInfo = (row, key) => {
    visibleRtpInfoControl.value = true
    dialogData.value = row
    dialogType.value = key
}

const openConfigDialog = (row, key) => {
    visibleRtpControl.value = true
    dialogData.value = row
    dialogType.value = key

}

const closeDialog = (value) => {
    if (!value) {
        visibleRtpControl.value = false
        dialogData.value = {}
        dialogType.value = ""
        queryList()
    }
}

const queryList = async () => {

    loading.value = true

   let [response, err] = await Client.Do(AdminInfo.GetOperatorInfo, {AppID:queryForm.value.AppID})



    queryRTPData.value = [response.OperatorInfo]

    loading.value = false
}




</script>

<style scoped>
.tableTitle{
    width: 98%;
    height: 60px;
    margin: 15px auto 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--el-color-primary);
    color: #ffffff;
    border-radius: 5px 5px 0 0;
    flex-wrap: wrap;
}
</style>
