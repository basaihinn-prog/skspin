/*
*@Author: 西西米
*@Date: 2022-12-23 14:26:05
*@Description: 管理员列表
*/
<template>
    <div class="adminList">
        <div class="searchView">
            <el-form
                    :model="adminListSearch"
                    style="max-width: 100%"
            >
                <el-space wrap>

                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operatorInfo="operatorListChange" :is-init="true" :hase-all="true"/>
                    <el-form-item :label="$t('用户名称')">
                        <el-input clearable v-model="adminListSearch.UserName" :placeholder="$t('请输入')"></el-input>
                    </el-form-item>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="getList(adminListSearch)">{{ $t('搜索') }}</el-button>
            </el-space>

        </div>

        <div class="page_table_context">
            <customTable
                table-name="adminList"
                :table-header="tableHeader"
                :table-data="tableData"
                :page-size="adminListSearch.PageSize"
                :page="adminListSearch.PageIndex"
                :count="Count"
                @page-change="paginationEmit">

                <template #handleTools>


                    <el-button type="primary" @click="addDialog = true">{{ $t('添加管理员') }}</el-button>
                </template>


                <template #Qrcode="scope">
                    <template v-if="scope.scope.IsOpenGoogle">
                        <el-popover placement="right" trigger="click">
                            <template #reference>
                                <el-button text v-if="scope.scope.IsOpenGoogle" style="margin-top: .5rem" size="small">
                                    {{ $t('查看') }}
                                </el-button>
                            </template>
                            <el-image class="tableImg" ref="myButton" :src="scope.scope.Qrcode"
                                      :preview-src-list="[scope.scope.Qrcode]" :initial-index="0" preview-teleported
                                      fit="cover"/>
                        </el-popover>
                    </template>
                    <template v-else>
                        <span>{{ $t('需要开启谷歌验证码') }}</span>
                    </template>
                </template>


                <template #GoogleCode="scope">
                    <div>
                        <el-switch v-model="scope.scope.IsOpenGoogle"
                                   :before-change="() => changeOpenGoole(scope.scope.IsOpenGoogle, scope.scope.Id)"/>
                    </div>
                </template>


                <template #Status="scope">
                    <div>
                        <el-switch v-model="scope.scope.Status" :active-value="0" :inactive-value="1"
                                   :before-change="() => changeStatus(scope.scope.Status, scope.scope.Id)"/>
                    </div>
                </template>
                <template #Operator="scope">
                    <div style="display: flex">
                            <el-button size="small" type="primary" @click="openEditDialog(scope.scope)" plain :disabled="!scope.scope.OperatorAdmin || (scope.scope.OperatorAdmin && store.AdminInfo.GroupId > 1)">{{
                                    $t('编辑')
                                }}
                            </el-button>





                        <el-button size="small" type="danger" @click="resetPassword(scope.scope)" plain>
                            {{ $t('重置密码') }}
                        </el-button>


                    </div>
                </template>


            </customTable>


        </div>


        <!-- 添加弹框 -->
        <el-dialog v-model="addDialog" @closed="resetForm" :title="$t('添加管理员')"
                   :width="store.viewModel === 2 ? '85%' : '956px'">
            <el-form ref="addFormRef" :model="addForm" :rules="addRules" label-position="top" label-width="80px">
                <el-form-item :label="$t('用户名')" prop="Username">
                    <el-input v-model="addForm.Username" :placeholder="$t('请输入')"/>
                </el-form-item>

                <el-form-item :label="$t('密码')" prop="Password">
                    <el-input v-model="addForm.Password" type="password" :placeholder="$t('请输入')"/>
                </el-form-item>
                <el-form-item :label="$t('权限组')" prop="PermissionId">
                    <el-select v-model="addForm.PermissionId" :placeholder="$t('请选择')">
                        <el-option v-for="item in groupList" :key="item.Id" :label="item.Name" :value="item.Id"/>
                    </el-select>
                </el-form-item>


            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('添加') }}</el-button>
                </span>
            </template>
        </el-dialog>


        <!-- 添加弹框 -->
        <el-dialog v-model="editDialog" @closed="resetForm" :title="$t('编辑管理员')"
                   :width="store.viewModel === 2 ? '85%' : '956px'">
            <el-form ref="editFormRef" :model="addForm" :rules="addRules" label-position="top" label-width="80px">
                <el-form-item :label="$t('用户名')" prop="Username">
                    <el-input disabled v-model="addForm.Username" :placeholder="$t('请输入')"/>
                </el-form-item>
                <el-form-item :label="$t('权限组')" prop="PermissionId">
                    <el-select v-model="addForm.PermissionId" :placeholder="$t('请选择')">
                        <el-option v-for="item in groupList" :key="item.Id" :label="item.Name" :value="item.Id"/>
                    </el-select>
                </el-form-item>


            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="editAdminer">{{ $t('修改') }}</el-button>
                </span>
            </template>
        </el-dialog>


    </div>
</template>

<script lang='ts' setup>
import {onMounted, ref, reactive} from 'vue'
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia/index';
import type {FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {AddOperatorAdminReq, AdminInfo} from '@/api/adminpb/info';
import {AdminGroup} from '@/api/adminpb/group';
import customTable from "@/components/customTable/tableComponent.vue"
import {ElMessageBox} from "element-plus";
import Operator_container from "@/components/operator_container.vue";

const {t} = useI18n()
const store = useStore()
let tableHeader = [
    {
        label: "ID",
        value: "Id"
    },
    {
        label: "权限名称",
        value: "PermissionName"
    },
    {
        label: "商户AppID",
        value: "AppID"
    },
    {
        label: "用户名",
        value: "UserName"
    },
    {
        label: "谷歌二维码",
        value: "Qrcode",
        type: "custom",
        hiddenVisible: true

    },
    {
        label: "谷歌验证码",
        value: "GoogleCode",
        hiddenVisible: true,
        type: "custom"
    },
    {
        label: "有效标识",
        value: "Status",
        type: "custom",
        hiddenVisible: true
    },
    {
        label: "操作",
        value: "Operator",
        type: "custom",
        fixed:"right",

        width:"200px",
        hiddenVisible: true
    }
]
let tableData = ref([])
let param = reactive({
    PageIndex: 1,
    PageSize: 20
})
let adminListSearch = reactive({
    PageIndex: 1,
    PageSize: 20,
    AppID: '',
    UserName: '',
})

let groupList = ref([])
let groupListAll = reactive([])
let Count = ref(0)
let loading = ref(false)
let operatorList = ref([])
onMounted(() => {
    getGroupAllList()
    getGroupList()
    if (store.AdminInfo.AppID === 'admin') {
        getOperator()
    } else {
        adminListSearch.AppID = store.AdminInfo.OperatorName
    }
});


const getList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.AdminerList, adminListSearch)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    Count.value = data.Count
    tableData.value = data.Count === 0 ? [] : data.List
    const groupMap = {};
    groupListAll.forEach((t) => {
        groupMap[t.Id] = t.Name;
    });

// 遍历 tableData 并更新 PermissionName
    tableData.value.forEach((i) => {
        const permissionId = i.PermissionId;
        if (groupMap.hasOwnProperty(permissionId)) {
            i['PermissionName'] = groupMap[permissionId];
        }
    });
}
const getGroupAllList = async () => {
    let [data, err] = await Client.Do(AdminGroup.GroupListAll, {
        PageIndex: 1,
        PageSize: 10000
    })


    groupListAll = data.Count === 0 ? [] : data.List
    getList()
}
const getGroupList = async () => {
    let [data, err] = await Client.Do(AdminGroup.GroupList, {
        PageIndex: 1,
        PageSize: 10000
    })


    groupList.value = data.Count === 0 ? [] : data.List

}
const getOperator = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, {
        PageIndex: 1,
        PageSize: 10000,
        OperatorType: 2,
        Status: 1
    } as any)
    if (err) {
        return tip.e(err)
    }
    operatorList.value = data.AllCount === 0 ? [] : data.List
}
const changeOpenGoole = (IsOpenGoogle, Id) => {
    return new Promise(async (resolve) => {
        let [data, err] = await Client.Do(AdminInfo.UpdateOpenGoole, {
            Open: !IsOpenGoogle, Id
        })
        if (err) {
            tip.e(err)
            return resolve(false)
        }
        getList()
        return resolve(true)

    })
}

const changeStatus = (Status, Id) => {
    // /AdminInfo/UpdateUserStatus
    return new Promise(async (resolve) => {
        let [data, err] = await Client.Do(AdminInfo.AdminerStatus, {
            Status: Status === 1 ? 0 : 1, Id
        })
        if (err) {
            tip.e(err)
            return resolve(false)
        }
        tip.s(t('修改成功'))
        return resolve(true)
    })
}

let addDialog = ref(false)
let editDialog = ref(false)
const addFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()
let addForm = ref({
    Id: "",
    Username: "",
    Password: "",
    AppID: "",
    isRoot: store.AdminInfo.AppID === 'admin',
    PermissionId: "",
})
const addRules = reactive<FormRules>({
    Username: [{required: true, message: t('请填写用户名'), trigger: 'blur'}],
    Password: [{required: true, message: t('请填写密码'), trigger: 'blur'}],
    PermissionId: [{required: true, message: t('请选择权限组'), trigger: 'change'}],
    AppID: [{required: true, message: t('请选择商户'), trigger: 'change'}],
})
const defaultOperatorEvent = ref({})
const operatorListChange = (value) => {
    if (value){

        adminListSearch.AppID = value.AppID
    }else{

        adminListSearch.AppID = ""
    }
}
const AddAdminer = async () => {
    let showApi = AdminInfo.AddAdmin

    if (!addFormRef) return
    await addFormRef.value.validate(async (valid, fields) => {
        if (valid) {
            ElMessageBox.confirm(
                t('确认添加用户'),
                t('是否确认'),
                {
                    confirmButtonText: t('确定'),
                    cancelButtonText: t('关闭'),
                    type: 'warning',
                }
            )
                .then(async () => {
                    let [data, err] = await Client.Do(showApi, addForm.value)
                    if (err) {
                        return tip.e(err)
                    }
                    tip.s(t('添加成功'))
                    addDialog.value = false
                    getList()
                })
        }
    })
}


const editAdminer = async (Pid) => {
    ElMessageBox.confirm(
        t('确认修改用户'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            let [data, err] = await Client.Do(AdminInfo.EditAdminer, {
                Id: addForm.value.Id,
                PermissionId: addForm.value.PermissionId
            } as any)
            if (err) {
                return tip.e(err)
            }
            tip.s(t('修改成功'))
            editDialog.value = false
            getList()
        })
}

const resetForm = () => {
    addForm.value.Username = ''
    addForm.value.Password = ""
    addForm.value.AppID = ""
    addForm.value.PermissionId = ""
    if (!addFormRef) {

        addFormRef.value.resetFields()
    }
    if (!editFormRef) {


        editFormRef.value.resetFields()
    }
}
const openEditDialog = (row) => {
    editDialog.value = true
    addForm.value.Id = row.Id
    addForm.value.Username = row.UserName
    addForm.value.AppID = row.AppID
    addForm.value.PermissionId = row.PermissionId
}
const resetPassword = async (row) => {
    ElMessageBox.confirm(
        t('确认重置该用户密码'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            let [data, err] = await Client.Do(AdminInfo.ResetPassword, {ID: row.Id, AppID: row.AppID} as any)
            if (err) {
                return tip.e(t(err))
            }
            tip.s(t("重置成功"))
        })
}


const paginationEmit = (pageConfig) => {

    adminListSearch.PageIndex = pageConfig.currentPage
    adminListSearch.PageSize = pageConfig.dataSize
    getList()
}
</script>
<style scoped lang='scss'>
:deep(.el-popover) {
  max-width: 100px !important;
  min-width: auto;
}

.flex_child_end {
  margin-bottom: 15px;
}
</style>
