<template>
    <div class="customTableContainer">
        <div class="tableToolContainer" v-if="!props.hideTableHandel" :style="{background : visibleTableHeaderBackground ? '#f7f7f7' : '#ffffff'}">
            <div class="tableTool">
                <div class="tableTool_item" style="overflow-x: auto;white-space: nowrap">
                    <slot name="handleTools"></slot>
                </div>
                <div class="tableTool_item2">
                    <el-breadcrumb separator="|">
                        <el-breadcrumb-item>
                            <div class="tableHandleSwitch">
                                <span>{{ $t("边框") }}</span>
                                <el-switch size="small" v-model="tableBorder"/>
                            </div>

                        </el-breadcrumb-item>
                        <el-breadcrumb-item>
                            <div class="tableHandleSwitch">
                                <span>{{ $t("表头") }}</span>
                                <el-switch size="small" v-model="visibleTableHeaderBackground"/>
                            </div>

                        </el-breadcrumb-item>

                        <el-breadcrumb-item>
                            <el-tooltip
                                class="box-item"
                                effect="dark"
                                :content="tableTipsAllContext && $t(tableTipsAllContext['表格刷新'])"
                                :open-delay="openDelay || 300"
                                :hide-after="hideAfter || 300"
                                placement="top-start"
                            >
                                <el-button size="small" :icon="RefreshRight" @click="refresh"
                                           style="width: 30px;height:30px;padding: 8px;background: #f7f7f7"/>

                            </el-tooltip>


                        </el-breadcrumb-item>
                        <el-breadcrumb-item>


                            <el-dropdown trigger="click">
                                <div>
                                <el-tooltip
                                        class="box-item"
                                        effect="dark"
                                        :content="tableTipsAllContext && $t(tableTipsAllContext['表格松紧'])"
                                        :open-delay="openDelay"
                                        :hide-after="hideAfter"
                                        placement="top-start"
                                >

                                    <el-button size="small"
                                               style="width: 30px;height:30px;padding: 8px;background: #f7f7f7">
                                        <el-image :src="lashenIcon" style="width: 18px;"></el-image>
                                    </el-button>

                                </el-tooltip>
                                </div>

                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item @click="changeTableLine('large')">
                                            <el-text :type="lineStatus == 'large' ? 'primary' : ''">{{
                                                $t("宽松")
                                                }}
                                            </el-text>
                                        </el-dropdown-item>
                                        <el-dropdown-item @click="changeTableLine('default')">
                                            <el-text :type="lineStatus == 'default' ? 'primary' : ''">{{
                                                $t("中等")
                                                }}
                                            </el-text>
                                        </el-dropdown-item>
                                        <el-dropdown-item @click="changeTableLine('small')">
                                            <el-text :type="lineStatus == 'small' ? 'primary' : ''">{{
                                                $t("紧凑")
                                                }}
                                            </el-text>
                                        </el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>


                        </el-breadcrumb-item>

                        <el-breadcrumb-item>
                            <el-tooltip
                                class="box-item"
                                effect="dark"
                                :content="tableTipsAllContext && $t(tableTipsAllContext['表格设置'])"
                                :open-delay="openDelay"
                                :hide-after="hideAfter"
                                placement="top-start"
                            >
                                <el-button size="small" :icon="Setting" @click="visibleTableState = true"
                                           style="width: 30px;height:30px;padding: 8px;background: #f7f7f7">
                                </el-button>


                            </el-tooltip>
                        </el-breadcrumb-item>

                    </el-breadcrumb>
                </div>
            </div>
        </div>
        <div class="customTable">
            <el-table :header-cell-style="{ background: '#f7f7f7', color: '#333333' }"
                      :data="tableDataList"
                      :fit="true"
                      style="width: 100%"
                      :border="tableBorder"
                      @cell-click="props.cellClick"
                      :height="props.height ? props.height : '500px'"
                      @selection-change="handleSelectionChange"
                      :size="lineStatus">
                <template v-for="item in tableHeaderList">
                    <template v-if="checkedCities.indexOf(item.value) != -1">


                        <el-table-column
                                :label="$t(item.label)"
                                :prop="item.value"
                                style="padding: 5px 0"
                                show-overflow-tooltip
                                align="center"
                                :sortable="item.sortable"
                                :column-key="item.value"
                                :fixed="item.fixed"
                                :width="item.width"
                                v-if="item.type == 'custom'">

                            <template #header="scope">

                                <table-tips :tips="tableTipsContext && tableTipsContext[item.label]" :tips-slot-name="item.headerName" :openDelay="openDelay" :hideAfter="hideAfter"/>
                                {{ $t(item.label) }}
                            </template>

                            <template #default="scope">
                                <slot :name="item.value" :scope="scope.row" :index="scope.$index"/>
                            </template>
                        </el-table-column>
                        <el-table-column style="padding: 5px 0" show-overflow-tooltip align="center"
                                         :sortable="item.sortable" :column-key="item.value" :fixed="item.fixed"
                                         :type="item.type" :prop="item.value" :label="$t(item.label)"
                                         :width="item.width"
                                         v-else-if="item.type && item.type != 'custom'">

                            <template #header="scope">

                                <table-tips :tips="tableTipsContext && tableTipsContext[item.label]" :tips-slot-name="item.headerName" :openDelay="openDelay" :hideAfter="hideAfter"/>
                                {{ $t(item.label) }}
                            </template>
                        </el-table-column>
                        <el-table-column style="padding: 5px 0" show-overflow-tooltip align="center"
                                         :sortable="item.sortable" :column-key="item.value" :fixed="item.fixed"
                                         :prop="item.value" :label="$t(item.label)" :width="item.width"
                                         :formatter="item.format" v-else>

                            <template #header="scope">

                                <table-tips :tips="tableTipsContext && tableTipsContext[item.label]" :tips-slot-name="item.headerName" :openDelay="openDelay" :hideAfter="hideAfter"/>
                                {{ $t(item.label) }}
                            </template>
                        </el-table-column>
                    </template>
                </template>
            </el-table>
        </div>

        <div class="table_pagination">
            <el-pagination
                    v-if="currentPage && pageSize"
                    v-model:current-page="currentPage"
                    v-model:page-size="dataSize"
                    :page-sizes="[10, 20, 50, 100]"
                    :disabled="disabled"
                    :background="background"
                    layout="->,slot,prev, pager,next,sizes,jumper,->"
                    :total="dataCount"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    :size="lineStatus"
            >

                <span>{{ $t('共计') }}{{ dataCount }}{{ $t('条数据') }}</span>


            </el-pagination>
        </div>


        <!--   对表格头部设置的dialog     -->
        <TableHeaderVisibleDialog :table-header="tableHeaderList" v-model="visibleTableState"
                                  @controlTableHeader="getVisibleTableHeader"/>

        <uploadExcelData v-model="excelDialogVisible" @uploadData="getUploadData"/>
    </div>
</template>

<script setup lang="ts">

import {useI18n} from "vue-i18n";

import {computed, ref, watchEffect, unref, onMounted, watch} from "vue";
import lashenIcon from "@/assets/login/lashen.png"
import { RefreshRight, Setting} from "@element-plus/icons-vue";
import {useStore} from "@/pinia";
import TableHeaderVisibleDialog from "@/components/customTable/tableHeaderVisibleDialog.vue";
import tableTips from "./tableTips.vue";

import {useRoute} from "vue-router";
import UploadExcelData from "@/components/customTable/uploadExcelData.vue";
import {OperatorRTP} from "@/enum";

const emits = defineEmits(["selectionChange", "pageChange", "pageSizeChange", "refreshTable", "update:tableData", "generatorMockData"])
const props = defineProps(["tableHeader", "tableData", "page", "pageSize", "count", "height", "checkData", "cellClick", "tableName", "hideTableHandel"])

const tableDataList = ref([])
const excelDialogVisible = ref(false)
const visibleTableState = ref(false)
const visibleTableHeaderList = ref([])
const currentPage = computed(() => parseInt(props.page))
const dataSize = computed(() => parseInt(props.pageSize))
const dataCount = computed(() => parseInt(props.count) || 0)
const checkData = computed(() => props.checkData)
const store = useStore()
const route = useRoute()
const tableHeaderList = computed({
    get() {
        visibleTableHeaderList.value = props.tableHeader.map(item => ({...item, visible: true}))
        checkedCities.value = props.tableHeader.map(item => item.value)
        return visibleTableHeaderList.value
    },
    set(value) {
        visibleTableHeaderList.value = value
        visibleTableHeaderList.value = value
    }
})

watchEffect(() => {
    tableDataList.value = props.tableData
})

// 是否显示表格线
const openDelay = ref(300)
const hideAfter = ref(300)
const tableBorder = ref()
const checkedCities = ref([])
// 是否显示表头
const visibleTableHeaderBackground = ref(true)

const lineStatus = ref("default")

interface pageInterface {
    currentPage: number
    dataSize: number
}

const tableTipsContext = ref(null)
const tableTipsAllContext = ref(null)

onMounted(() => {
    let tableSetConfig = store.SystemConfig.tableBorderConfig
    let flag = false

    if (props.tableName) {
        flag = tableSetConfig && tableSetConfig[props.tableName]
    }
    tableBorder.value = flag


    lineStatus.value = localStorage.getItem("lineStatus") || "default"




    tableTipsContext.value = store.tipsMap[route.meta.title]
    tableTipsAllContext.value = store.tipsMap['表格通用']

})

watch(tableBorder, (newData) => {

    if (props.tableName) {

        let tableSetConfig = store.SystemConfig.tableBorderConfig
        if (!tableSetConfig) {

            tableSetConfig = {}
        }
        tableSetConfig[props.tableName] = newData
        localStorage.setItem("tableSetConfig", JSON.stringify(tableSetConfig))
    }

})


let {t} = useI18n()

const background = ref(true)
const disabled = ref(false)
const loading = ref(false)


const getVisibleTableHeader = (data: string[]) => {
    checkedCities.value = data

    tableHeaderList.value.forEach(item => {
        item.visible = checkedCities.value.includes(item.value)
    })

}

const getUploadData = (data) => {
    emits("generatorMockData", data)
}

const handleSizeChange = (value) => {
    tableDataList.value = []
    emits("pageChange", <pageInterface>{
        currentPage: 1,
        dataSize: value,
    })
}
const handleCurrentChange = (value) => {
    tableDataList.value = []
    emits("pageChange", <pageInterface>{
        currentPage: value,
        dataSize: dataSize.value,
    })
}

const handleSelectionChange = (value) => {
    emits("selectionChange", value)
}

const changeTableLine = (selectLine) => {
    lineStatus.value = selectLine
    localStorage.setItem("lineStatus", selectLine)
}

const refresh = () => {
    emits("refreshTable")
}

</script>

<style scoped lang="scss">
.customTable {
  width: 100%;
  height: auto;
  border-top: 1px solid #dcdfe6;
  border-bottom: none;
  margin-bottom: 15px;
}

.tableToolContainer {
  width: 100%;
  height: auto;
  background: #f7f7f7;
}


.tableTool {
  width: 98%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  margin: auto;
}

.customTableContainer {
  width: 98%;
  margin: 0 auto;
  border-radius: 5px;
  border: 1px solid #e5e5e5;
}

.tableHandleSwitch {
  height: 30px;
  display: flex;
  align-items: center;

}

.tableHandleSwitch > span {
  margin-right: 5px;
}
</style>

<style>

.customTable .cell, .customTable .cell .el-icon {
    vertical-align: middle;
}

.customTable .cell, .customTable .cell * {
    font-size: 14px;
}

.el-table .caret-wrapper {
    width: auto;
}

.checkBoxItem:first-child {
    border-bottom: 1px solid #eeeeee;
}

.checkBoxItem {
    width: 100%;
    height: auto;
}

.header_check_item {
    width: 200px;
    padding-left: 10px;
}

.table_pagination {
    width: 100%;
    height: auto;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 10px;
}
.tableTool_item{
   flex: 6;

}
.tableTool_item2{
   flex: 2.2!important;

    display: flex;
    align-items: center;
    justify-content: right;
}
</style>
