<template>
    <div>
      <div class="searchView gameList">
          <el-space wrap>
              <el-form
                      style="max-width: 100%"
                      @keyup.enter="initSearch"
              >
                  <el-space wrap>
                      <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                          @select-operatorInfo="operatorListChange" :is-init="true" :hase-all="true"/>

                      <el-form-item :label="$t('时间')">
                          <el-date-picker
                              v-model="timeRange"
                              type="daterange"
                              value-format="x"
                              format="YYYY-MM-DD"
                              :range-separator="$t('至')"
                              :shortcuts="shortcuts"
                              :clearable="false"
                              :start-placeholder="$t('开始时间')"
                              :end-placeholder="$t('结束时间')"
                              :disabled-date="option"
                          />

                      </el-form-item>
                      <el-form-item>


                          <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
                      </el-form-item>
                  </el-space>
              </el-form>
          </el-space>
      </div>

        <div class="page_table_context">
            <customTable
                v-loading="loading"
                table-name="gameUserRetainedList"
                :table-header="tableHeader"
                :table-data="tableData"
                :page-size="uiData.PageSize"
                :page="uiData.Page"
                :count="count"
                @refresh-table="initSearch"
                @pageChange="PageChange"
            >

                <template #Operator="scope">
                    <el-button size="small" type="primary" plain @click="checkInfo(scope.scope)">{{ $t('详情') }}</el-button>
                </template>
            </customTable>
        </div>




        <InfoDialog v-model="dialog" :ReportId="ReportId"></InfoDialog>


    </div>
</template>

<script setup lang="ts">

import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {onMounted, Reactive, reactive, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import {Client} from "@/lib/client";
import {GameUserRetained, GameUserRetainedReq} from "@/api/adminpb/gameUserRetained";
import ut from "@/lib/util";
import InfoDialog from "./component/InfoDialog.vue";
import {excel} from "@/lib/excel";
import {useStore} from "@/pinia";
const {t} = useI18n()
const store = useStore()

const ReportId = ref("")
const Operator = ref(0)
const count = ref(0)
const loading = ref(false)
const dialog = ref(false)
const timeRange = ref([])

watch(dialog, (newData)=>{

    if (!newData){
        ReportId.value = ""
    }
})



const option = (value) => {


    return value > new Date()
}

const formatTime = (time) => {
    let timeStr = ut.fmtDate(new Date(time).getTime()/1000)
    return ut.fmtDate(timeStr, "YYYY-MM-DD")
}

// 留存数据列表头
const tableHeader = ref([
    {label: "商户AppID", value: "AppID"},
    {label: "日期", value: "Date", format: (row) => formatTime(row.Date)},
    {label: "投注用户数", value: "RetentionPlayerCount"},
    {label: "次留", value: "RetentionPlayer1d", format:(row) => `${row.RetentionPlayer1d.toFixed(2)}%`},
    {label: "三日留存", value: "RetentionPlayer3d", format:(row) => `${row.RetentionPlayer3d.toFixed(2)}%`},
    {label: "七日留存", value: "RetentionPlayer7d", format:(row) => `${row.RetentionPlayer7d.toFixed(2)}%`},
    {label: "十四日留存", value: "RetentionPlayer14d", format:(row) => `${row.RetentionPlayer14d.toFixed(2)}%`},
    {label: "三十日留存", value: "RetentionPlayer30d", format:(row) => `${row.RetentionPlayer30d.toFixed(2)}%`},
    {label: "操作", value: "Operator", type: "custom", hiddenVisible: true},
])

// 留存数据列表
const tableData = ref([])

// 时间查询
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

const uiData:Reactive<GameUserRetainedReq> = reactive<GameUserRetainedReq>({
    Operator: 0,
    StartTime: 0,
    EndTime: 0,
    PageSize:20,
    Page:1,
})

onMounted(()=>{

    const today = new Date();
    const year = today.getFullYear();
    const month = today.getMonth();
    const date = today.getDate();

    const startTime = new Date(year, month, date, 0, 0, 0, 0o00).getTime();
    const endTime = new Date(year, month, date, 23, 59, 59, 0o00).getTime();

    timeRange.value = [startTime, endTime]
    initSearch()
})

const defaultOperatorEvent = ref()
const operatorListChange = (value) => {
    if (value){

        uiData.Operator = value.Id == "ALL" ? null : value.Id
    }else{

        uiData.Operator = 0
    }
}


const initSearch = async () => {

    const requestData = {
        ...uiData,
    }

    requestData.StartTime = ut.fmtUTCDate(timeRange.value[0])
    requestData.EndTime = ut.fmtUTCDate(timeRange.value[1])


    tableData.value = []
    loading.value = true
    let [data, err] = await Client.Do(GameUserRetained.GameUserRetainedList, requestData);
    loading.value = false



    if (data.All) {
        tableData.value = data.List
    }
    count.value = data.All
}

const PageChange = (page) => {
    uiData.Page = page.currentPage
    uiData.PageSize = page.dataSize
    initSearch()
}

const checkInfo = (row) => {
    ReportId.value = row.ID
    dialog.value = true
}


</script>

<style scoped lang="scss">

</style>
