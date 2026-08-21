<template>
    <div>
      <div class="searchView gameList">
          <el-space wrap>
              <el-form
                      style="max-width: 100%"

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

                <template #handleTools>
                    <el-space>
                        <upload-excel @uploadFile="resolveData" fileName="游戏留存"/>

                        <el-button type="primary" plain @click="generatorExcel">生成Excel</el-button>

                    </el-space>
                </template>
                <template #Operator="scope">
                    <el-button size="small" type="primary" plain @click="checkInfo(scope.scope)">{{ $t('详情') }}</el-button>
                </template>
            </customTable>
        </div>






    </div>
</template>

<script setup lang="ts">

import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue";
import {nextTick, onMounted, Reactive, reactive, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import {GameUserRetained, GameUserRetainedReq} from "@/api/adminpb/gameUserRetained";
import ut from "@/lib/util";
import UploadExcel from "./component/uploadExcel.vue";

import {excel} from "@/lib/excel";
import {Client} from "@/lib/client";
const {t} = useI18n()

const ReportId = ref("")
const Operator = ref(0)
const count = ref(0)
const loading = ref(false)
const dialog = ref(false)
const uploadStatus = ref(false)
const timeRange = ref([new Date().getTime(), new Date().getTime()])


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
const tableHeaderMap = ref(null)
const tableHeader = ref([
    {label: "商户AppID", value: "AppID"},
    {label: "日期", value: "Date", format: (row) => formatTime(row.Date)},
    {label: "投注用户数", value: "RetentionPlayerCount"},
    {label: "次留", value: "RetentionPlayer1d", format:(row) => `${Number(row.RetentionPlayer1d).toFixed(2)}%`},
    {label: "三日留存", value: "RetentionPlayer3d", format:(row) => `${Number(row.RetentionPlayer3d).toFixed(2)}%`},
    {label: "七日留存", value: "RetentionPlayer7d", format:(row) => `${Number(row.RetentionPlayer7d).toFixed(2)}%`},
    {label: "十四日留存", value: "RetentionPlayer14d", format:(row) => `${Number(row.RetentionPlayer14d).toFixed(2)}%`},
    {label: "三十日留存", value: "RetentionPlayer30d", format:(row) => `${Number(row.RetentionPlayer30d).toFixed(2)}%`},
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


    const startTime = new Date(timeRange.value[0]);
    startTime.setHours(0)
    startTime.setMinutes(0)
    startTime.setSeconds(0)
    startTime.setMilliseconds(0)

    const endTime = new Date(timeRange.value[1]);
    endTime.setHours(23)
    endTime.setMinutes(59)
    endTime.setSeconds(59)
    endTime.setMilliseconds(0)

    requestData.StartTime = startTime.getTime() / 1000
    requestData.EndTime = endTime.getTime() / 1000

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

    if (uploadStatus.value){
        uploadStatus.value = false
        return
    }

    initSearch()
}

const checkInfo = (row) => {
    ReportId.value = row.ID
    dialog.value = true
}
// 生成Excel
const generatorExcel = (row) => {

    let generatorHeader = [...tableHeader.value].filter(item => item.label != '操作')


    excel.DataGeneratorExcel(generatorHeader, tableData.value, `游戏留存`)
}
// 解析Excel
const resolveData = (data) => {
    uploadStatus.value = true
    if (!tableHeaderMap.value) {
        tableHeaderMap.value = {}
        for (const i in tableHeader.value) {

            tableHeaderMap.value[tableHeader.value[i].label] = tableHeader.value[i].value
        }
    }

    count.value = data.data.length
    let renderData = data.data.slice(0, uiData.PageSize)
    uiData.Page = 1

    nextTick(()=>{
        tableData.value = excel.excelResolveData(tableHeaderMap.value, renderData)
        uploadStatus.value = false
    })

}
</script>

<style scoped lang="scss">

</style>
