<template>
    <div>
        <div class="searchView">
            <el-form
                    :model="queryForm"
                    style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange" :is-init="true"></operator_container>

                    <el-form-item :label="$t('唯一标识')">
                        <el-input maxlength="12" v-model.number.trim="queryForm.Pid"
                                  oninput="value=value.replace(/[^\d]/g,'')"
                                  clearable :placeholder="$t('请输入唯一标识')"/>
                    </el-form-item>
                    <gamelist_container
                        :hase-manufacturer="true"
                        :hase-all="false"
                        :is-init="true"
                        :defaultGameEvent="defaultGameEvent"
                        @select-operator="selectGameList"
                    />


                    <el-form-item :label="$t('控制时间')">
                        <el-date-picker
                                v-model="ControlTime"
                                type="datetimerange"
                                :start-placeholder="$t('开始时间')"
                                :end-placeholder="$t('结束时间')"
                                format="YYYY-MM-DD HH:mm:ss"
                                date-format="MMM DD, YYYY"
                                time-format="HH:mm"/>
                    </el-form-item>
<!--                    <el-form-item :label="$t('游戏名称')">-->
<!--                        <el-select v-model="queryForm.GameId" style="width: 150px" :placeholder="t('请选择游戏')"-->
<!--                                   filterable clearable>-->
<!--                            <el-option v-for="(item, index) in games" :key="index" :label="$t(item.Name)"-->
<!--                                       :value="item.ID"></el-option>-->
<!--                        </el-select>-->
<!--                    </el-form-item>-->


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

            <customTable
                table-name="playerRTPData_list"
                @refresh-table="queryList"
                v-loading="loading"
                :table-header="tableHeader"
                :table-data="queryRTPData"
                :count="Count"
                :page="queryForm.Page"
                :page-size="queryForm.PageSize"
                @page-change="pageChange"
                @selection-change="selectionChange">


                <template #handleTools>

                    <el-button type="warning" v-if="visibleRtpControl" plain @click="cancelRTPControl">{{ $t("批量取消") }}
                    </el-button>
                </template>

                <template #GameName="scope">
                    <div style="display: flex;align-items: center;">
                        <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>
                        <div style="margin: 0 10px">
                        <el-avatar :src="scope.scope.GameIcon" size="small" />
                        </div>
                            {{ scope.scope.GameName }}
                    </div>
                </template>
                <template #Operator="scope">

                    <el-button type="primary" plain @click="singleRTPCancel(scope.scope)" size="small">{{ $t("取消") }}</el-button>
                </template>
            </customTable>
        </div>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, Ref, ref} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {useStore} from "@/pinia";
import ut from "@/lib/util";

const {t} = useI18n()
import {useI18n} from "vue-i18n";
import {AdminPlayer, CancelControlRequest, PlayerRTPData, PlayerRTPListRequest} from "@/api/adminpb/adminPlayer";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {ElMessage, ElMessageBox} from "element-plus";
import Gamelist_container from "@/components/gamelist_container.vue";

const store = useStore();
let loading = ref(false)


let tableHeader = ref([
    {label: "游戏名称", value: "GameName", type: "custom", width: '240px'},
    {label: "控制时间", value: "ControlTime", format: (row) => ut.fmtSelectedUTCDateFormat(row.ControlTime)},
    {label: "商户AppID", value: "OperatorName"},
    {label: "唯一标识", value: "Pid"},
    {label: "游戏配置RTP", value: "GameRTP", format: (row) => row.GameRTP + "%", width: "140px"},
    {label: "玩家历史RTP", value: "PlayerHistoryRTP", format: (row) => row.PlayerHistoryRTP + "%", width: "140px"},
    {label: "控制RTP", value: "ControlRTP", format: (row) => row.ControlRTP + "%"},
    {label: "被控制时RTP", value: "ControllingRTP", format: (row) => row.ControllingRTP + "%", width:"140px"},
    {label: "自动解除RTP", value: "AutoRemoveRTP", format: (row) => row.AutoRemoveRTP + "%", width:"140px"},
    {label: "购买免费游戏RTP", value: "BuyRTP", format: (row) => row.BuyRTP + "%" },
])




const games = ref([])
const visibleRtpControl = ref(false)
const CurrencyCode = ref("")
const ControlTime = ref([])
const tableData = ref([])
const queryRTPData: Ref<PlayerRTPData[]> = ref<PlayerRTPData[]>([])
const queryForm: Ref<PlayerRTPListRequest> = ref<PlayerRTPListRequest>({
    OperatorId: 0,               // 商户ID
    Pid: "",                      // 用户登录账号
    ControlTimeStart: "",         // 控制时间开始
    ControlTimeEnd: "",           // 控制时间结束
    Manufacturer: "",           // 控制时间结束
    GameId: "",                  // 控制时间结束
    Page: 1,
    PageSize: 20
})
const CancelControlIds: Ref<CancelControlRequest> = ref<CancelControlRequest>({
    Ids: ""
})

const Count = ref(0)
const init = ref(true)
const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})




const operatorListChange = (value) => {
    const OperatorId = value ? value.Id : null
    queryForm.value.OperatorId = OperatorId
    if (init.value){

        queryList()
        init.value = false
    }
}

const selectGameList = (value) => {

    if (value.gameData){

        queryForm.value.GameId = value.gameData || null
    }

    if (value.manufacturer || value.manufacturer == null){
        queryForm.value.Manufacturer = value.manufacturer
    }
}


const singleRTPCancel = (row) => {
    CancelControlIds.value.Ids = row.Id.toString()

    cancelRTPControl()
}


const cancelRTPControl = async () => {

    if (CancelControlIds.value.Ids == "") {
        return
    }

    ElMessageBox.confirm(
        t('确认取消当前的RTP控制'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {

            const [res, err] = await Client.Do(AdminPlayer.CancelPlayerRTPControl, CancelControlIds.value)
            if (err) {
                tip.e(t("RTP控制失败"))
                return
            }
            queryList()

        })

}


const initGameList = async () => {
    games.value = []
    const [gameList, err] = await Client.Do(AdminGameCenter.GameList, null)

    if (!err) {

        games.value = gameList.List
    }

}


onMounted(() => {
    if (store.AdminInfo.GroupId <= 1) {
        visibleRtpControl.value = true
    }
    if (store.AdminInfo.GroupId === 3) {

        visibleRtpControl.value = store.AdminInfo.Businesses["PlayerRTPSettingOff"] == 1
    }


    if(visibleRtpControl.value){

        tableHeader.value.unshift({type: "selection", label: "", width:"40px",fixed: "left", hiddenVisible: true})
        tableHeader.value.push({label: "操作", value: "Operator", type: "custom", fixed: "right", hiddenVisible: true})
    }


})

const queryList = async () => {


    let queryData = {...queryForm.value,
        Uid: Number(queryForm.value.Pid)}
    let startTime = 0
    let endTime = 0

    if (ControlTime.value){
        if (ControlTime.value[0] && ControlTime.value[1]) {

            startTime = (ControlTime.value[0].getTime())
            endTime = (ControlTime.value[1].getTime())
        }
    }



    if (startTime) {

        queryData.ControlTimeStart = ut.fmtSelectedUTCDateFormat(startTime, "reduce") == "/" ? "" : ut.fmtSelectedUTCDateFormat(startTime, "reduce")
    }
    if (endTime) {

        queryData.ControlTimeEnd = ut.fmtSelectedUTCDateFormat(endTime, "reduce") == "/" ? "" : ut.fmtSelectedUTCDateFormat(endTime, "reduce")
    }



    loading.value = true
    let [res, err] = await Client.Do(AdminPlayer.GetPlayerRTPList, queryData)

    loading.value = false
    if (err) {
        return
    }

    Count.value = res.Count
    queryRTPData.value = res.List
}


initGameList()
const pageChange = (page) => {
    queryForm.value.Page = page.currentPage
    queryForm.value.PageSize = page.dataSize

    queryList()
}
const selectionChange = (data) => {
    let ids = []
    ids = data.map(item => item.Id)
    CancelControlIds.value.Ids = ids.join(",")
}
</script>
