<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-form
                    :model="RTPProtectForm"
                    style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container ref="operatorContainer2Ref" :hase-all="true" :is-init="true" :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operatorInfo="operatorListChange2"/>


                    <el-form-item :label="$t('功能状态')">
                        <el-select v-model="RTPProtectForm.IsProtection"  style="width: 150px" @change="changeIsProtection">
                            <el-option :label="$t('全部')" :value="-1"/>
                            <el-option v-for="item in ProtectionStatus" :label="$t(item.label)" :value="item.value"/>
                        </el-select>
                    </el-form-item>


                </el-space>
            </el-form>

            <el-space wrap>
                <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                <el-button type="default" @click="resetSearch">{{ $t('重置') }}</el-button>
            </el-space>
        </div>


        <div class="page_table_context">


            <customTable
                    v-loading="loading"
                    table-name="operatorApproval_list"
                    :table-data="tableData"
                    :table-header="tableHeader"
                    :page="RTPProtectForm.Page"
                    :page-size="RTPProtectForm.PageSize"
                    :count='Count'
                    @refresh-table="queryList"
                    @page-change="changePage">

                <template #OrderId="scope">
                    {{ scope.scope.OrderId }}
                </template>

                <template #funcStatus="scope">

                    <el-text :type="scope.scope.IsProtection == 0 ? 'danger' : 'success'" size="small">
                        {{ scope.scope.IsProtection == 0 ? $t('关闭') : $t('开启')}}
                    </el-text>
                </template>
                <template #Operator="scope">

                    <el-button type="primary" @click="openDialog(scope.scope)" size="small" plain>{{ $t('操作') }}</el-button>
                </template>

            </customTable>

        </div>


        <el-dialog :width="store.viewModel === 2 ? '85%' : '650px'" v-model="RTPProtectDialog" :title="$t('新玩家RTP保护')"
                   destroy-on-close>
            <el-form
                label-width="140px"
                ref="RTPProtectSetFormRef"
                :rules="RTPProtectSetFormRules"
                :model="RTPProtectSetForm"
                style="max-width: 100%"
            >


                    <el-form-item :label="$t('开启状态') + ':'" label-position="right" >

                        <el-select v-model="RTPProtectSetForm.IsProtection"  style="width: 240px" @change="changeIsProtection">
                            <el-option v-for="item in ProtectionStatus" :label="$t(item.label)" :value="item.value"/>
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('局数设定') + ':'" label-position="right" prop="ProtectionRotateCount">

                        <el-input v-model.number="RTPProtectSetForm.ProtectionRotateCount"
                                  @input="ProtectionRotateCountChange"
                                  :placeholder="$t('该局数为每个游戏单独计算')"
                                  style="width: 240px"
                                  :disabled="!RTPProtectSetForm.IsProtection"/>

                    </el-form-item>
                    <div style="color: var(--el-color-danger);width: 60%;font-size: 12px;margin-left: 140px;margin-bottom: 15px">{{ $t('设置区间为{Num}局', {Num: "1~200"}) }}</div>
                    <el-form-item :label="$t('上升RTP') + ':'" label-position="right">

                        <el-select v-model="RTPProtectSetForm.ProtectionRewardPercentLess" style="width: 240px" :disabled="!RTPProtectSetForm.IsProtection">


                            <el-option v-for="item in ProtectionRewardPercentLessList" :label="`${Number(item)}%`" :value="Number(item)"/>

                        </el-select>
                    </el-form-item>


            </el-form>


            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="save" type="primary">{{ $t('保存') }}</el-button>
                    <el-button @click="RTPProtectDialog = false">{{ $t('关闭') }}</el-button>
                </div>
            </template>
        </el-dialog>

    </div>
</template>
<script setup lang="ts">

import type {Ref} from "vue";
import {reactive, ref} from "vue";
import {useI18n} from "vue-i18n";

import {Client} from "@/lib/client";
import customTable from "@/components/customTable/tableComponent.vue";

import Operator_container from "@/components/operator_container.vue";
import {LOG_OPERATOR_TYPE, SystemLog} from "@/api/adminpb/log";
import {Merchant, RTPProtectInfo, RTPProtectListParams, RTPProtectListRep} from "@/api/adminpb/merchant";
import {useStore} from "@/pinia";
import {ElLoading, ElMessageBox, FormRules} from "element-plus";
import {tip} from "@/lib/tip";

const {t} = useI18n()
const store = useStore()

const operatorContainer2Ref = ref(null)

const RTPProtectDialog = ref(false)
const loading = ref(false)
const Count = ref(0)
const timeRange = ref([])
const RTPProtectForm: Ref<RTPProtectListParams> = ref(<RTPProtectListParams>{
    IsProtection: -1,
    AppID:"",
    Page: 1,
    PageSize: 20,
})
const ProtectionRotateCountValid = (rule: any, value: any, callback: any) => {

    if (RTPProtectSetForm.value.IsProtection == 1){

        if (value == null || value == "null" || value == undefined || value == "undefined" || value == ""){
            callback(t("局数不能为空"))
            return;
        }

        if(value < 1 || value > 200){
            callback(t("可设置值范围为1~200"))
            return;
        }
        callback()
    }else{
        callback()
    }

}
const RTPProtectSetFormRef = ref(null)

const RTPProtectSetForm: Ref<RTPProtectListParams> = ref(<RTPProtectListParams>{
    AppID: "",
    ProtectionRewardPercentLess: 4,
    IsProtection: 0,
    ProtectionRotateCount: 1,
})


const RTPProtectSetFormRules = reactive<FormRules<RTPProtectListParams>>({
    ProtectionRotateCount:[{validator:ProtectionRotateCountValid, trigger: 'blur'}],
})
const defaultGameEvent = ref({})
const defaultOperatorEvent = ref({})
const ProtectionStatus = [
    {label:"开启", value: 1},
    {label:"关闭", value: 0},
]
const staticLess = import.meta.env.VITE_PROTECTION_RTP.split(",")
const ProtectionRewardPercentLessList = ref([])
ProtectionRewardPercentLessList.value = staticLess
const tableHeader = [
    {label: "商户AppID", value: "AppID"},
    {label: "功能状态", value: "funcStatus", type:"custom"},
    {label: "新手局数", value: "ProtectionRotateCount"},
    {label: "上升RTP", value: "ProtectionRewardPercentLess", format:(row)=>row.ProtectionRewardPercentLess + "%"},
    {label: "操作", value: "Operator",fixed:"right", type:"custom", hiddenVisible: true}
]
const editData: Ref<RTPProtectInfo> = ref(<RTPProtectInfo>{})
const tableData: Ref<RTPProtectInfo[]> = ref(<RTPProtectInfo[]>[])


const init = async () => {
    tableData.value = []
    loading.value = true
    const queryData = <any>{
        ...RTPProtectForm.value,
    }


    const [response, err] = await Client.Do(Merchant.RTPProtectList, queryData)
    loading.value = false
    if (!err) {
        Count.value = response.All
    }



    tableData.value = response.List
}

const reviewData = () => {

    RTPProtectSetForm.value.ProtectionRewardPercentLess = editData.value.IsProtection == 1 ? editData.value.ProtectionRewardPercentLess : 0
    RTPProtectSetForm.value.AppID = editData.value.AppID
    RTPProtectSetForm.value.ProtectionRotateCount = editData.value.IsProtection == 1 ? editData.value.ProtectionRotateCount : null

}
const openDialog = (row) => {
    editData.value = row
    RTPProtectSetForm.value.IsProtection = editData.value.IsProtection
    reviewData()
    ProtectionRewardPercentLessList.value = editData.value.IsProtection == 1 ? staticLess : ['0']
    RTPProtectDialog.value = true
}




const changeIsProtection = (value) => {





    reviewData()
    ProtectionRewardPercentLessList.value = value == 0 ? ['0'] : staticLess
    if (value == 0){

        RTPProtectSetFormRef.value.validateField("ProtectionRotateCount")
        RTPProtectSetForm.value.ProtectionRotateCount = null
        RTPProtectSetForm.value.ProtectionRewardPercentLess = 0
    }else if (RTPProtectSetForm.value.ProtectionRewardPercentLess == 0){
        RTPProtectSetForm.value.ProtectionRotateCount = null
        RTPProtectSetForm.value.ProtectionRewardPercentLess = 4
    }
}

const save = async () => {


   await ElMessageBox.confirm(t('确认修改新手RTP保护配置'), t('是否确认'), {
        confirmButtonText: t('确定'),
        cancelButtonText: t('关闭'),
        type: 'warning',
    });


    RTPProtectSetFormRef.value.validate(async valid => {

        if (valid){
            const loading = ElLoading.service({
                lock: true,
                text: 'Loading',
                background: 'rgba(0, 0, 0, 0.7)',
            })

            let queryData = {...RTPProtectSetForm.value}
            const [response, err] = await Client.Do(Merchant.SetRTPProtect, queryData)
            loading.close()
            if (err) {
                tip.e(t(err))
                return
            }

            RTPProtectDialog.value = false
            init()
            tip.s(t("成功"))
        }

    })
}
const operatorListChange2 = (value) => {
    if (value) {

        RTPProtectForm.value.AppID = value.AppID
    } else {

        RTPProtectForm.value.AppID = ""
    }
}


const ProtectionRotateCountChange = (value) => {

    if (Number(value) === 0){
        RTPProtectSetForm.value.ProtectionRotateCount = 1
    }
    if (Number(value) > 200){
        RTPProtectSetForm.value.ProtectionRotateCount = 200
    }
}

const changePage = (page) => {
    RTPProtectForm.value.Page = page.currentPage
    RTPProtectForm.value.PageSize = page.dataSize
    init()
}

const queryList = () => {
    RTPProtectForm.value.Page = 1
    init()
}


init()

const resetSearch = () => {

    RTPProtectForm.value = {
        IsProtection: -1,
        OperatorName: "",
        AppID:"",
        Page: 1,
        PageSize: 20,
    }
    operatorContainer2Ref.value.paramData = "ALL"
    timeRange.value = []
}
</script>

<style scoped lang="scss">
.el-link .el-icon--right.el-icon {
  vertical-align: text-bottom;
}
</style>
