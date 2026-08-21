<template>


    <el-dialog
            v-model="dialogVisible"
            :title="$t('玩家数据')"
            top="2vh"
            :width="store.viewModel === 2 ? '85%' : '950px'"
            @close="closeDialog"
            destroy-on-close
    >
        <div class="searchView">
            <el-form
                    :model="queryData"
                    style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container
                        ref="operatorRef"
                        :is-init="true"
                        :defaultOperatorEvent="defaultOperatorEvent"
                        @select-operator="operatorListChanges"
                        @select-operatorInfo="operatorListChange"></operator_container>
                    <el-form-item >


                        <el-select
                                v-model="queryData.DateType"
                                placeholder="Select"
                                style="width: 100px"
                        >
                            <el-option
                                    :label="$t('登录时间')"
                                    :value="1"
                            />
                            <el-option
                                    :label="$t('创建时间')"
                                    :value="2"
                            />
                        </el-select>
                        <el-date-picker
                            v-model="timeRange"
                            type="datetimerange"
                            :start-placeholder="$t('开始时间')"
                            :end-placeholder="$t('结束时间')"
                            format="YYYY-MM-DD HH:mm:ss"
                            date-format="YYYY/MM/DD ddd"
                            time-format="A hh:mm:ss"
                        />
                    </el-form-item>



                    <el-form-item :label="$t('今日下注')">
                        <div style="padding-right: 10px">
                            <el-input v-model.trim="queryData.NowBet" clearable
                                      onkeyup="value=value.replace(/[^\d\.\d{0,2}]/g, '')" :placeholder="$t('请输入')">
                                <template #prefix>≥</template>
                            </el-input>
                        </div>
                    </el-form-item>

                    <el-form-item :label="$t('历史下注')">
                        <div style="padding-right: 10px">
                            <el-input v-model.trim="queryData.historyBet" clearable
                                      onkeyup="value=value.replace(/[^\d\.\d{0,2}]/g, '')" :placeholder="$t('请输入')">
                                <template #prefix>≥</template>
                            </el-input>
                        </div>
                    </el-form-item>
                    <el-form-item :label="$t('今日输赢')">
                        <div style="padding-right: 10px">
                            <el-input v-model.trim="queryData.NowWin" clearable
                                      onkeyup="value=value.replace(/[^\d\.\d{0,2}]/g, '')" :placeholder="$t('请输入')">
                                <template #prefix>≥</template>
                            </el-input>
                        </div>
                    </el-form-item>

                    <el-form-item :label="$t('历史输赢')">
                        <div style="padding-right: 10px">
                            <el-input v-model.trim="queryData.HistoryWin" clearable
                                      onkeyup="value=value.replace(/[^\d\.\d{0,2}]/g, '')" :placeholder="$t('请输入')">
                                <template #prefix>≥</template>
                            </el-input>
                        </div>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                    </el-form-item>
                </el-space>
            </el-form>
        </div>
        <div>

            <customTable
                table-name="playerAnalysis_RTPDialogPanel_list"
                height="300" :table-header="paymentHeader" :table-data="paymentData"
                :count="Count" :page="queryData.Page" :page-size="queryData.PageSize"
                @refresh-table="queryList"
                @selection-change="SelectPlayer" @page-change="pageChange" v-loading="loadding">


                <template #handleTools>

                    <el-button type="warning" plain @click="PreSubmit">{{ $t('批量修改RTP') }}</el-button>
                </template>

                <template #NowBet="scope">
                        {{ut.toNumberWithCommaNormal(scope.scope.NowBet / 10000)}}
                </template>
                <template #historyBet="scope">
                        {{ut.toNumberWithCommaNormal(scope.scope.HistoryBet / 10000)}}
                </template>
                <template #NowWin="scope">
                        {{ut.toNumberWithCommaNormal(scope.scope.NowWin / 10000)}}
                </template>
                <template #HistoryWin="scope">
                        {{ut.toNumberWithCommaNormal(scope.scope.HistoryWin / 10000)}}
                </template>
                <template #Balance="scope">

                    <div style="width: 100%;white-space: normal;text-align: center">
                        {{ut.toNumberWithCommaNormal(scope.scope.Balance / 10000)}}
                    </div>

                </template>

            </customTable>

        </div>

        <template #footer>
            <div class="dialog-footer">
                <el-button @click="closeDialog">{{ $t('关闭') }}</el-button>
            </div>
        </template>
    </el-dialog>

</template>

<script setup lang="ts">

import customTable from "@/components/customTable/tableComponent.vue"
import {computed, Ref, ref} from "vue";
import {useI18n} from "vue-i18n";
import {AdminPlayer, PlayerResponse, PlayRTP, SelectPlayRTPRequest} from "@/api/adminpb/adminPlayer";
import Operator_container from "@/components/operator_container.vue";
import ut from "@/lib/util";
import {Client} from "@/lib/client";
import {tip} from "@/lib/tip";
import {useStore} from "@/pinia";

const {t} = useI18n()
const store = useStore()

const props = defineProps(["modelValue", "playerId"])
const emits = defineEmits(["RTPControlData", "update:modelValue", "update:playerId"])


let paymentData: Ref<PlayRTP[]> = ref<PlayRTP[]>([])
let paymentHeader = [
    {width:"40px" ,label: "", value: "", type: "selection", align:"left", fixed:"left", hiddenVisible: true},
    {width:"120px" ,label: "唯一标识", value: "Pid"},
    {width:"130px" ,label: "今日下注额", value: "NowBet", type:"custom"},                 // 后端要求除1w
    {width:"130px" ,label: "历史下注额", value: "historyBet", type:"custom"},         // 后端要求除1w
    {width:"130px" ,label: "今日输赢", value: "NowWin", type:"custom"},                  // 后端要求除1w
    {width:"130px" ,label: "历史输赢", value: "HistoryWin", type:"custom"},          // 后端要求除1w
    {width:"130px" ,label: "余额", value: "Balance", type:"custom", align: "center"},
    {width:"130px" ,label: "所属商户", value: "AppID", format: (row) => row.AppID},
    {width:"150px" ,label: "创建时间", value: "CreateTime", format: (row) => ut.fmtDate(row.CreateTime)},
    {width:"150px" ,label: "最后登录时间", value: "LastLoginTime", format: (row) => ut.fmtDate(row.LoginTime)},
]

const timeRange = ref([])
const Count = ref(0)
const operatorRef = ref(null)

const queryData: Ref<SelectPlayRTPRequest> = ref<SelectPlayRTPRequest>({
    DateType: 1,
    AppID: "",
    StartTime: "",
    EndTime: "",
    NowBet: null,
    historyBet: null,
    NowWin: null,
    HistoryWin: null,
    Page: 1,
    PageSize: 20,
})


const loadding = ref(false)
const operatorData = ref("")
const HighRTPOff = ref(1)
const dialogVisible = computed(() => {

    if (!props.modelValue) {

        timeRange.value = []
        operatorData.value = ""
        queryData.value = {
            AppID: "",
            StartTime: "",
            EndTime: "",
            DateType: 1,
            NowBet: null,
            historyBet: null,
            NowWin: null,
            HistoryWin: null,
            Page: 1,
            PageSize: 20,
        }

    }

    return props.modelValue
})

const defaultOperatorEvent = ({})
const operatorListChanges = (value) => {
    if (value){

        operatorData.value = value.value

    }else{

        operatorData.value = null
    }
}
const operatorListChange = (value) => {
    if (value){

        queryData.value.AppID = value.AppID
        HighRTPOff.value = value.HighRTPOff
    }else{

        queryData.value.AppID = ""
    }
}

const queryList = async () => {

    paymentData.value = []
    let startTime = 0
    let endTime = 0



    if (timeRange.value[0] && timeRange.value[1]) {

        startTime = timeRange.value[0].getTime() / 1000
        endTime = timeRange.value[1].getTime() / 1000

    }

    if (startTime) {

        queryData.value.StartTime = ut.fmtDateSecond(startTime) == "/" ? "" : ut.fmtDateSecond(startTime)
    }
    if (endTime) {

        queryData.value.EndTime = ut.fmtDateSecond(endTime) == "/" ? "" : ut.fmtDateSecond(endTime)
    }




    let queryForm = {
        ...queryData.value,
        NowBet: Number(queryData.value.NowBet),
        historyBet: Number(queryData.value.historyBet),
        NowWin: Number(queryData.value.NowWin),
        HistoryWin: Number(queryData.value.HistoryWin),
    }

    if (queryForm.AppID == ""){
        return tip.e(t("商户不能为空"))
    }
    loadding.value = true

    const [response, err] = await Client.Do(AdminPlayer.PlayerPayInfoList, queryForm)

    loadding.value = false
    if (response){
        Count.value = response.Count
        paymentData.value = response.List

    }
}
const RTPPlayers = ref([])

const init = () => {

}

const SelectPlayer = (value) => {
    RTPPlayers.value = value
}

const pageChange = (page) => {
    queryData.value.Page = page.currentPage
    queryData.value.PageSize = page.dataSize
    queryList()
}

const PreSubmit = () => {


    if (RTPPlayers.value.length <= 0){

        tip.e(t("请选择用户"))
        return
    }

    emits("RTPControlData", RTPPlayers, HighRTPOff)
}

const closeDialog = () => {
    emits('update:modelValue')
}
</script>


<style scoped lang="scss">
.playerInfo {
  font-size: 20px;
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
</style>
