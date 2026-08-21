<template>
    <div class="customTableContainer">
        <div class="tableToolContainer" :style="{background: visibleTableHeader ? '#f7f7f7f7' : '#ffffff'}">
            <div class="tableTool">
                <div class="tableTool_item">
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
                                <el-switch size="small" v-model="visibleTableHeader"/>
                            </div>

                        </el-breadcrumb-item>

                        <el-breadcrumb-item>
                            <el-tooltip
                                class="box-item"
                                effect="dark"
                                :content="tableTipsAllContext && $t(tableTipsAllContext['表格刷新'])"
                                :open-delay="openDelay"
                                :hide-after="hideAfter"
                                placement="top-start"
                            >

                                <el-button size="small" :icon="RefreshRight" @click="refresh" style="width: 30px;height:30px;padding: 8px;background: #f7f7f7"/>

                            </el-tooltip>


                        </el-breadcrumb-item>

                        <el-breadcrumb-item>


                            <el-dropdown trigger="click">
                                <div>
                                <el-tooltip
                                    class="box-item"
                                    effect="dark"
                                    :content="tableTipsAllContext && $t(tableTipsAllContext['表格松紧'])"
                                    :open-delay="openDelay || 300"
                                    :hide-after="hideAfter || 300"
                                    placement="top-start"
                                >

                                    <el-button size="small" style="width: 30px;height:30px;padding: 8px;background: #f7f7f7">
                                        <el-image :src="lashenIcon" style="width: 18px;"></el-image>
                                    </el-button>
                                </el-tooltip>
                                </div>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item @click="changeTableLine('large')"><el-text :type="lineStatus == 'large' ? 'primary' : ''" >{{ $t("宽松") }}</el-text></el-dropdown-item>
                                        <el-dropdown-item @click="changeTableLine('default')"><el-text :type="lineStatus == 'default' ? 'primary' : ''" >{{ $t("中等") }}</el-text></el-dropdown-item>
                                        <el-dropdown-item @click="changeTableLine('small')"><el-text :type="lineStatus == 'small' ? 'primary' : ''" >{{ $t("紧凑") }}</el-text></el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>



                        </el-breadcrumb-item>
                        <el-breadcrumb-item>
                            <el-tooltip
                                class="box-item"
                                effect="dark"
                                :content="tableTipsAllContext && $t(tableTipsAllContext['表格设置'])"
                                :open-delay="openDelay || 300"
                                :hide-after="hideAfter || 300"
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
            <el-table
                ref="tableRef"
                height="450px"
                v-loading="loading"
                :data="tableDataList"
                :row-key="props.rowKey ? props.rowKey : 'Id'"
                :lazy="props.loadChild != undefined && props.loadChild != 'undefined' && props.loadChild != 'null' && props.loadChild != null"
                :size="lineStatus"
                :border="tableBorder"
                :tree-props="props.treeProps ? props.treeProps : { children: 'children', hasChildren: 'hasChildren' }"
                :load="load"
                :fit="true"
                :header-cell-style="{ background: '#f7f7f7', color: '#333333' }">
                <template v-for="item in tableHeaderList">

                    <template v-if="checkedCities.indexOf(item.value) != -1">


                        <el-table-column style="padding: 5px 0" :fixed="item.fixed" show-overflow-tooltip :align="item.align"
                                         :prop="item.value" :label="$t(item.label)" :width="item.width"
                                         v-if="item.type == 'custom'">

                            <template #header="scope" v-if="tableTipsContext && tableTipsContext[item.label]">

                                <template v-if="item.headerCustom">

                                    <slot :name="item.headerName"/>
                                    {{ $t(item.label) }}
                                </template>
                                <template v-else>

                                    <table-tips :tips="tableTipsContext && tableTipsContext[item.label]"
                                                :tips-slot-name="item.headerName"
                                                :openDelay="openDelay"
                                                :hideAfter="hideAfter"/>
                                    {{ $t(item.label) }}
                                </template>
                            </template>
                            <template #default="scope">

                                <slot :name="item.value" :scope="scope.row"/>

                            </template>
                        </el-table-column>
                        <el-table-column style="padding: 5px 0" show-overflow-tooltip :fixed="item.fixed"
                                         :align="item.align" type="index" :prop="item.value" :label="$t(item.label)"
                                         :width="item.width"
                                         v-else-if="item.type == 'index'">
                            <template #header="scope" v-if="tableTipsContext && tableTipsContext[item.label]">

                                <table-tips :tips="tableTipsContext && tableTipsContext[item.label]" :tips-slot-name="item.headerName"
                                                               :openDelay="openDelay"
                                                               :hideAfter="hideAfter"/>
                                {{$t(item.label)}}
                            </template>
                        </el-table-column>
                        <el-table-column style="padding: 5px 0" show-overflow-tooltip :fixed="item.fixed"
                                         :align="item.align" :prop="item.value" :label="$t(item.label)"
                                         :width="item.width" :formatter="item.format" v-else>
                            <template #header="scope" v-if="tableTipsContext && tableTipsContext[item.label]">

                                <table-tips :tips="tableTipsContext && tableTipsContext[item.label]" :tips-slot-name="item.headerName"
                                            :openDelay="openDelay"
                                            :hideAfter="hideAfter"/>
                                {{$t(item.label)}}
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
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50, 100]"
                :size="lineStatus"
                :disabled="disabled"
                :background="background"
                layout="slot,->, prev, pager, next, sizes,jumper"
                :total="dataCount"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
        >

            <span>{{ $t('共计') }}{{dataCount}}{{ $t('条数据') }}</span>


        </el-pagination>
        </div>
        <!--   对表格头部设置的dialog     -->
        <TableHeaderVisibleDialog :table-header="tableHeaderList" v-model="visibleTableState" @controlTableHeader="getVisibleTableHeader"/>
    </div>
</template>

<script setup lang="ts">

import {useI18n} from "vue-i18n";

import {computed, ref, watchEffect, onMounted, watch} from "vue";
import {QuestionFilled, RefreshRight, Setting} from "@element-plus/icons-vue";
import lashenIcon from "@/assets/login/lashen.png";
import {useStore} from "@/pinia";
import TableHeaderVisibleDialog from "@/components/customTable/tableHeaderVisibleDialog.vue";
import TableTips from "@/components/customTable/tableTips.vue";
import {useRoute} from "vue-router";

const props = defineProps(["tableHeader", "tableData", "page", "pageSize", "count", "loadChild", "tableName", "rowKey", "treeProps"])
const emits = defineEmits(["pageChange", "update:tableData", "refreshTable"])

const store = useStore()
const currentPage = computed(() => props.page)
const pageSize = computed(() => props.pageSize)
const dataCount = computed(() => props.count || 0)
const tableDataList = ref([])
const tableRef = ref(null)


watchEffect(()=>{
    tableDataList.value = props.tableData
})

// 是否显示表格线
const tableBorder = ref(true)
// 是否显示表头
const visibleTableHeader = ref(true)

let lineStatus = ref("default")


let {t} = useI18n()
let route = useRoute()

const openDelay = ref(300)
const hideAfter = ref(300)
const background = ref(true)
const disabled = ref(false)
const loading = ref(false)

const visibleTableState = ref(false)
const tableTipsContext = ref({})
const tableTipsAllContext = ref(null)
const checkedCities = ref([])
const visibleTableHeaderList = ref([])

onMounted(()=>{
    let tableSetConfig = store.SystemConfig.tableBorderConfig

    let flag = false
    if (tableSetConfig && props.tableName){
        flag = tableSetConfig[props.tableName]
    }
    tableBorder.value = flag


    lineStatus.value = localStorage.getItem("lineStatus")



    tableTipsContext.value = store.tipsMap[route.meta.title]
    tableTipsAllContext.value = store.tipsMap["表格通用"]
})

watch(tableBorder, (newData)=>{

    let tableSetConfig = store.SystemConfig.tableBorderConfig
    if (tableSetConfig && props.tableName){

        tableSetConfig[props.tableName] = newData
        localStorage.setItem("tableSetConfig", JSON.stringify(tableSetConfig))
    }

})

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

class pageInterface {
}

const refresh = () => {
    emits("refreshTable")
}

const handleSizeChange = (value) => {
    tableDataList.value = []
      emits("pageChange", <pageInterface>{
        currentPage: currentPage.value,
        dataSize: value,
    })

}
const handleCurrentChange = (value) => {

    tableDataList.value = []
       emits("pageChange", <pageInterface>{
        currentPage: value,
        dataSize: pageSize.value,
    })

}

const load = async (   row: any,
                       treeNode: unknown,
                       resolve: (data: any[]) => void) => {


    let [resp, err] = await props.loadChild(row.AppID)

    if (!resp.List){
        resp.List = []
    }



    resolve(resp.List)




}

const changeTableLine = (changeStatus) => {
    lineStatus.value = changeStatus

    localStorage.setItem("lineStatus", changeStatus)
}

const getVisibleTableHeader = (data: string[]) => {

    checkedCities.value = data

    tableHeaderList.value.forEach(item=>{
        item.visible = checkedCities.value.includes(item.value)
    })

}


defineExpose({
    tableRef
})


</script>

<style scoped lang="scss">
.customTable {
  width: 100%;
    height: auto;
  border: 1px solid #dcdfe6;
  border-bottom: none;
  margin-bottom: 15px;
}

.tableHandleSwitch{
    height: 30px;
    display: flex;
    align-items: center;

}
.tableHandleSwitch>span{
    margin-right: 5px;
}



.tableTool {
    width: 98%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 0;
    margin: auto;
}

.tableToolContainer{
    background: #f7f7f7
}
</style>

<style>


.table_pagination {
    width: 100%;
    height: auto;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 10px;
}

.customTableContainer {
    width: 98%;
    margin: 0 auto;
    border-radius: 5px;
    border: 1px solid #e5e5e5;
}
.tableTool_item{
    flex: 4;

}
.tableTool_item2{
    flex: 6;
    display: flex;
    align-items: center;
    justify-content: right;
}
</style>
