<template xmlns="http://www.w3.org/1999/html">
    <div>


        <div class="searchView">
            <el-form :model="systemStatus" style="max-width: 100%">

                <el-form-item :label="$t('半闭维护')" style="margin-bottom: 15px">
                    <el-tooltip
                            class="box-item"
                            effect="dark"
                            :content="$t('半闭维护是白名单用户可以进入游戏')"
                            placement="top-start"

                    >
                        <el-icon size="20">
                            <QuestionFilled/>
                        </el-icon>
                    </el-tooltip>


                    <el-switch v-model="systemStatus.HalfMaintenance" :active-value="1" :inactive-value="0"  @change="switchChange('halfMaintenance', $event)"></el-switch>
                </el-form-item>

                <el-form-item :label="$t('全闭维护')" style="margin-bottom: 15px">
                    <el-tooltip
                            class="box-item"
                            effect="dark"
                            :content="$t('全闭维护是所有用户无法进入场馆')"
                            placement="top-start"

                    >
                        <el-icon size="20">
                            <QuestionFilled/>
                        </el-icon>
                    </el-tooltip>

                    <el-switch v-model="systemStatus.EntireMaintenance" :active-value="1" :inactive-value="0"  @change="switchChange('entireMaintenance', $event)"></el-switch>
                </el-form-item>

            </el-form>

        </div>
        <div class="page_table_context">


            <el-tabs v-model="activeName" class="demo-tabs" @tabClick="handleClick">
            <el-tab-pane :label="$t('操作日志')" name="operatorList">



                    <customTable
                        table-name="platformMaintenance_log_list"
                        :table-header="OperatorTableHeader"
                        :table-data="maintenanceLogRes.List"
                        :count="maintenanceLogRes.Count"
                        :page-size="maintenanceLogReq.PageSize"
                        :page="maintenanceLogReq.Page"
                        @pageChange="maintenancePageChange"
                        @refresh-table="refreshTableData"
                        >

                        <template #handleTools>
                            <el-button type="primary" plain @click="whitListDialog = true">{{ $t('添加白名单') }}</el-button>
                        </template>

                </customTable>


            </el-tab-pane>
            <el-tab-pane :label="$t('半闭白名单')" name="WhitList">


                    <customTable
                        table-name="platformMaintenance_whit_list"
                        :table-header="WhitListTableHeader"
                        :table-data="WhitListRes.List"
                        :checkData="checkWhitUser"
                        :count="WhitListRes.Count"
                        :page-size="whitListReq.PageSize"
                        :page="whitListReq.Page"
                        @pageChange="whitListPageChange"
                        @selectionChange="whitSelect"
                    >
                        <template #handleTools>

                            <el-button type="primary" plain @click="whitListDialog = true">{{ $t('添加白名单') }}</el-button>
                            <el-button type="warning" @click="deleteWhiteUsers" plain style="margin-right: 15px">
                                {{ $t('批量删除白名单') }}</el-button>
                        </template>
                    </customtable>

            </el-tab-pane>
        </el-tabs>
        </div>



        <!--    维护面板    -->
        <el-dialog v-model="maintenanceDialog" :title="$t('维护')" width="550px" @close="closeDialog" destroy-on-close>
            <el-form :model="maintenanceTime">
                <el-form-item>
                    <el-date-picker
                        v-model="maintenanceTime.timeRang"
                        type="datetimerange"
                        :range-separator="$t('至')"
                        :start-placeholder="$t('开始时间')"
                        :end-placeholder="$t('结束时间')"
                        format="YYYY-MM-DD HH:mm:ss"

                    />
                </el-form-item>
                <el-form-item>
                    <div>
                    <span style="color: red">{{ $t('注') }}</span>：{{ $t('全闭维护是所有用户无法进入场馆') }}
                        <br/>

                        {{ $t('半闭维护是在系统当中的IP才能访问场馆') }}</div>
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="closeDialog" style="margin-right: 15px">{{ $t("关闭") }}</el-button>
                    <el-button type="primary" @click="confirmMaintenance">
                        {{ $t("添加") }}
                    </el-button>
                </div>
            </template>
        </el-dialog>



        <!--    维护面板    -->
        <el-dialog v-model="whitListDialog"
                   :title="$t('添加白名单')"
                   width="550px"
                   @close="closeDialog"
                   @open-auto-focus="openWhitDialog"
               >
            <el-form ref="formEl" :model="whitForm" :rules="whitFormRule" >
                <el-form-item :label="$t('用户IP') + ':'" label-width="120px" prop="UserIp">
                    <el-input v-model="whitForm.UserIp" style="width: 240px" :placeholder="$t('请输入用户IP')" />
                </el-form-item>
                <el-form-item :label="$t('备注') + ':'" label-width="120px">
                    <el-input
                        v-model="whitForm.Remark"
                        style="width: 240px"
                        :rows="2"
                        type="textarea"
                        :placeholder="$t('请输入备注')"
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="closeDialog" style="margin-right: 15px">{{ $t("关闭") }}</el-button>
                    <el-button type="primary" @click="commitWhitList">
                        {{ $t("添加") }}
                    </el-button>
                </div>
            </template>
        </el-dialog>


    </div>
</template>

<script setup lang="ts">

import {reactive, Reactive, Ref, ref} from "vue";
import {QuestionFilled} from "@element-plus/icons-vue";
import customTable from "@/components/customTable/tableComponent.vue"
import ut from "@/lib/util";
import {Client} from "@/lib/client";
import {AdminGameCenter, MaintenanceLog, MaintenanceLogResponse, SystemConfig} from "@/api/gamepb/admin";
import {useI18n} from "vue-i18n";
import {CommPageRequest} from "@/api/comm";
import {ElMessage, ElMessageBox, FormRules} from "element-plus";
import {WhitData, White, WhitResponse} from "@/api/adminpb/WhitList";
import {tip} from "@/lib/tip";
const { t } = useI18n()

const formEl = ref(null)
const activeName = ref("operatorList")
const maintenanceDialog = ref(false)
const whitListDialog = ref(false)

const checkWhitUser = ref([])

const maintenanceTime = ref({
    timeRang:[]
})
const whitForm:Ref<WhitData> = ref<WhitData>({
    UserIp: "",
    Remark: ""
})

const isIPv4validator = (rule: any, value: any, callback: any) => {
    let whitelistString = value.split(',').filter(item => item !== "")
    const isIPv4 = new RegExp(/^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})$/);
    let whiteBools = isIPv4.test(whitelistString)
    if (whiteBools) {
        callback()
    } else {
        callback(new Error(t('白名单不合法')))
    }
}

const whitFormRule = reactive<FormRules<WhitData>>({
    UserIp: [
        {required: true, message: t('用户IP不能为空'), trigger: 'blur'},
        {trigger: 'blur',validator: isIPv4validator}
    ],

})


const maintenanceLogRes:Ref<MaintenanceLogResponse> = ref<MaintenanceLogResponse>({
    List:[],
    Count: 0
})

const WhitListRes:Ref<WhitResponse> = ref<WhitResponse>({
    List:[],
    Count: 0
})

const maintenanceLogReq:Ref<CommPageRequest> = ref<CommPageRequest>({
    Page: 1,
    PageSize: 20
})

const whitListReq:Ref<CommPageRequest> = ref<CommPageRequest>({
    Page: 1,
    PageSize: 20
})

const systemStatus:Ref<SystemConfig> = ref<SystemConfig>({
    EntireMaintenance: 0,
    HalfMaintenance: 0,
    StartTime: 0,
    EndTime: 0,
})



const OperatorTableHeader = [
    {label: "序号", value: "Date", type: "index", width:"60px"},
    {label: "维护模式", value: "MaintenanceMod", format: (row) => {
        if (row.MaintenanceMod == 0){
            return t("关闭");
        } if (row.MaintenanceMod == 1){
            return t("半闭维护");
        } if (row.MaintenanceMod == 2){
            return t("全闭维护");
        }
    }},
    {label: "维护状态", value: "MaintenanceStatus", format: (row) => row.MaintenanceStatus == 0 ? t("维护未完成"): t("维护完成") },
    {label: "开始时间", value: "StartTime", format: (row) => ut.fmtSelectedUTCDateFormat(row.StartTime * 1000)},
    {label: "结束时间", value: "RealEndTime", format: (row) => ut.fmtSelectedUTCDateFormat(row.RealEndTime * 1000)},
    {label: "维护操作人", value: "OperatorId", hiddenVisible:true},
]

const WhitListTableHeader = [
    {label: "", value: "Date", type: "selection", width: "40px", hiddenVisible:true},
    {label: "IP地址", value: "UserIp"},
    {label: "添加时间", value: "OperatorTime", format: (row) => ut.fmtSelectedUTCDateFormat(row.OperatorTime) },
    {label: "备注", value: "Remark"},
    {label: "添加操作人", value: "OperatorName", hiddenVisible:true},
]

const handleClick = (tab) => {
    activeName.value = tab
}

const deleteWhiteUsers = async () => {
    if (checkWhitUser.value.length == 0){
        return
    }
    ElMessageBox.confirm(
        t('是否确定删除选中的白名单用户'),
        t('提示'),
        {
            confirmButtonText: t("确认"),
            cancelButtonText: t("关闭"),
            type: 'warning',
        }
    ).then(async () => {

            const [res, err] = await Client.Do(White.DeleteWhitUser, {Ids: checkWhitUser.value.join(",")} as any)
            if (err){
                tip.e(t(err))
                return
            }



            await initWhitUserList()
        })


}

const switchChange = (key, value) => {

    if (value == 1){

        maintenanceDialog.value = true
    }

    if (key == 'halfMaintenance' && value == 1){
        systemStatus.value.EntireMaintenance = 0
    }
    if (key == 'entireMaintenance' && value == 1){
        systemStatus.value.HalfMaintenance = 0
    }

    if (systemStatus.value.EntireMaintenance == 0 && systemStatus.value.HalfMaintenance == 0){

        ElMessageBox.confirm(
            t('是否确定关闭当前维护状态'),
            t('提示'),
            {
                confirmButtonText: t("确认"),
                cancelButtonText: t("关闭"),
                type: 'warning',
            }
        )
            .then(() => {

                confirmMaintenance()
            })
            .catch(() => {
                closeDialog()
            })

    }

}

// 获取场馆是否关闭数据
const initConfig = async () => {
    const [res, err] = await Client.Do(AdminGameCenter.SelectMerchant, null)


    let EntireMaintenance = 0, HalfMaintenance = 0


    if (!err && res.Maintenance){
        switch (res.Maintenance) {
            case 1:
                HalfMaintenance = 1
                break;
            case 2:
                EntireMaintenance = 1
                break;
            default:
                EntireMaintenance = 0
                HalfMaintenance = 0
                break;
        }



    }
    systemStatus.value.EntireMaintenance = EntireMaintenance
    systemStatus.value.HalfMaintenance = HalfMaintenance

}

// 获取日志
const initMaintenanceLog = async () => {
    const [res, err] = await Client.Do(AdminGameCenter.GetMaintenanceLog, maintenanceLogReq.value)

    maintenanceLogRes.value = res
}

// 获取白名单
const initWhitUserList = async () => {

    const [res, err] = await Client.Do(White.GetWhitUserList, whitListReq.value)
    WhitListRes.value = res
}

// 关闭
const closeDialog = () => {
    maintenanceDialog.value = false
    whitListDialog.value = false
    init()
}

// 提交调整维护状态
const confirmMaintenance = async () => {
    if ( maintenanceTime.value.timeRang[0]){

        systemStatus.value.StartTime = maintenanceTime.value.timeRang[0].getTime() / 1000
    }
    if ( maintenanceTime.value.timeRang[1]){

        systemStatus.value.EndTime = maintenanceTime.value.timeRang[1].getTime() / 1000
    }


    await Client.Do(AdminGameCenter.ModifyConfig, systemStatus.value)

    maintenanceDialog.value = false

    await init()
}

//提交白名单信息
const commitWhitList = () => {
    formEl.value.validate(async (valid) => {
        if (valid) {



               const [response, err] = await Client.Do(White.AddWhitUser, whitForm.value)

            if (err){
                tip.e(t(err))
                return
            }



            whitListDialog.value = false
            init()
        } else {
            ElMessage({
                message: t(valid.message),
                type: 'warning',
            })
        }
    })

}

const openWhitDialog = () => {
    whitForm.value = {
        UserIp: "",
        Remark: ""
    }

}



const refreshTableData = () =>{
    init()
}

// 操作日志分页
const maintenancePageChange = async ({ currentPage, dataSize}) => {

    maintenanceLogReq.value.Page = currentPage
    maintenanceLogReq.value.PageSize = dataSize

    await initMaintenanceLog()
}

// 白名单分页
const whitListPageChange = async ({ currentPage, dataSize}) => {

    whitListReq.value.Page = currentPage
    whitListReq.value.PageSize = dataSize

    await initWhitUserList()
}

const whitSelect = (val) => {

    if (!val.length){
        checkWhitUser.value = []
    }

    for (const index in val) {
        checkWhitUser.value.push(val[index].Id)
    }

}
const init = async () => {

    await initConfig()
    await initMaintenanceLog()
    await initWhitUserList()
}
init()


</script>

<style scoped lang="scss">

</style>
