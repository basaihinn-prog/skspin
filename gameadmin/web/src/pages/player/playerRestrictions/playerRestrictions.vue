
<template>
<div>
    <div class="searchView">
        <el-form
                :model="queryForm"
                style="max-width: 100%"
        >
            <el-space wrap>
                <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                    @select-operatorInfo="operatorListChange" :hase-all="true" :is-init="true"></operator_container>

                <el-form-item :label="$t('唯一标识')">
                    <el-input maxlength="12" v-model.number.trim="queryForm.Pid" :placeholder="$t('请输入')"
                              clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('玩家ID')">
                    <el-input clearable v-model.trim="queryForm.Uid" :placeholder="$t('请输入')"></el-input>
                </el-form-item>
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
            v-loading="loading"
            table-name="playerRestrictions_list"
            :table-data="tableData"
            :table-header="tableHeader"
            :page="queryForm.Page"
            :page-size="queryForm.PageSize"
            :count='Count'
            @refresh-table="queryList"
            @page-change="pageChange">


            <template #handleTools>
                <el-button @click="dialogVisible = true"  type="warning" plain>

                    {{ $t("新增") }}
                </el-button>
            </template>
            <template #Operator="scope">
                <el-button size="small" plain type="primary" @click="editRestriction(scope.scope)">{{ $t('编辑') }}</el-button>
                <el-button size="small" plain type="danger" @click="cancelRestriction(scope.scope)">{{ $t('取消') }}</el-button>
            </template>
        </customTable>


        <restrictionDialog v-model="dialogVisible" @updateFrom="queryList" @update:modelValue="RestrictionsData = null" :RestrictionData="RestrictionsData"/>

    </div>
</div>
</template>

<script setup lang="ts">

import Operator_container from "@/components/operator_container.vue";
import customTable from "@/components/customTable/tableComponent.vue"
import restrictionDialog from "./component/restrictionDialog.vue"
import {ref} from "vue";
import {useI18n} from "vue-i18n";
import ut from "@/lib/util";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {Client} from "@/lib/client";
import {ElMessageBox} from "element-plus";
import {tip} from "@/lib/tip";

const {t} = useI18n()
const loading = ref(false)
const dialogVisible = ref(false)
const RestrictionsData = ref(null)


const defaultOperatorEvent = ref({})
const tableHeader = [
    {label: "商户AppID", value: "AppID"},
    {label: "唯一标识", value: "Pid"},
    {label: "玩家终生最大赢钱总数", value: "RestrictionsMaxWin", width:"180px", format:(row)=> ut.toNumberWithCommaNormal(row.RestrictionsMaxWin)},
    {label: "玩家终生赢取倍数", value: "RestrictionsMaxMulti", width:"180px", format:(row)=> ut.toNumberWithCommaNormal(row.RestrictionsMaxMulti)},
    {label: "当前赢钱总数", value: "Win", format:(row)=> ut.toNumberWithCommaNormal(row.Win / 10000)},
    {label: "当前赢取倍数", value: "Multi", format:(row)=> ut.toNumberWithCommaNormal(row.Multi)},
    {label: "设置时间", value: "RestrictionsTime", format: (row) => ut.fmtSelectedUTCDateFormat(row.RestrictionsTime)},
    {label: "操作", value: "Operator", type: "custom", hiddenVisible: true},
]
const tableData = ref([])

const Count = ref(0)
const queryForm = ref({
    AppID:null,
    Pid:null,
    Uid:"",
    Page: 1,
    PageSize: 20,
})

const operatorListChange = (value) => {
    if (value){

        queryForm.value.AppID = value.AppID == "ALL" ? null : value.AppID
    }else{

        queryForm.value.AppID = null
    }
}

const editRestriction = (row) => {
    dialogVisible.value = true
    RestrictionsData.value = row
}
const cancelRestriction = (row) => {

    ElMessageBox.confirm(
        t('确认取消限制'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }).then(async () => {

            loading.value = true
        let requestData = {
            Pid: row.Pid
        }

        const [response, err] = await Client.Do(AdminGameCenter.CancelPlayerRestrictions, requestData)
        queryList()
        if (err != null) {
            tip.e(t("调整失败"))
            return
        }

    })
}

const queryList = async () => {
    loading.value = true
    tableData.value = []
    let requestData = {
        ...queryForm.value,
        Pid: Number(queryForm.value.Pid)
    }
   const [response, err] = await Client.Do(AdminGameCenter.GetPlayerRestrictionsList, requestData)
    loading.value = false
    if (err != null) return

    tableData.value = response.List
    Count.value = response.AllCount
}
queryList()

const pageChange = (page) => {
    queryForm.value.PageSize = page.dataSize
    queryForm.value.Page = page.currentPage



    queryList();

}
</script>



<style scoped lang="scss">

</style>
