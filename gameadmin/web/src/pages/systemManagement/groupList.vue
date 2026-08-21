/*
*@Author: 西西米
*@Date: 2022-12-24 10:02:31
*@Description: 权限组列表
*/
<template>
    <div class='groupList'>

        <!-- 数据 -->

        <div class="page_table_context">
            <div class="flex flex_child_end">

            </div>
            <customTable
                table-name="groupList_list"
                :table-header="tableHeader"
                :table-data="groupList"
                :count="params.Count"
                :page="params.PageIndex"
                :page-size="params.PageSize"
            >


                <template #handleTools>
                    <el-button type="primary" @click="addDialogClick">{{
                        $t('添加权限组')
                        }}
                    </el-button>
                </template>

                <template #Operator="scope">

                    <el-button type="primary" @click="getPowerDetails(scope.scope)"
                               plain
                               size="small"
                               style="margin-bottom: 3px">
                        {{ $t("设置") }}
                    </el-button>
                </template>
            </customTable>

        </div>
        <!-- 添加弹框 -->
        <el-dialog v-model="addDialog" :title="$t('权限')" :width="store.viewModel === 2 ? '85%' : '550px'">
            <el-form ref="addFormRef" :model="changeForm" :rules="changeRules" label-position="top" label-width="80px">
                <el-form-item :label="$t('权限名称')" prop="Name">
                    <el-input v-model="changeForm.Name" :placeholder="$t('请输入')"/>
                </el-form-item>

                <el-form-item :label="$t('设置权限')" v-if="arrMenuTreeList.length>0" prop="MenuIds">
                    <el-tree ref="treeRef" :data="arrMenuTreeList" show-checkbox accordion node-key="ID" :props="{
                    children: 'Children',
                    label: 'Title',
                }"/>
                </el-form-item>
                <el-form-item :label="$t('备注')">
                    <el-input v-model="changeForm.Remark" :placeholder="$t('请输入')"/>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddMenu">{{ $t('添加') }}</el-button>
                </span>
            </template>
        </el-dialog>
        <!-- 设置权限弹框 -->
        <el-dialog v-model="setMenuDialog" :title="$t('设置权限')" :destroy-on-close="true"
                   :width="store.viewModel === 2 ?  '85%' : '956px'">
            <el-form ref="editFormRef" :model="editForm" :rules="changeRules" label-position="top" label-width="100px">
                <el-form-item :label="$t('权限名称')" prop="Name">
                    <el-input v-model="editForm.Name" :placeholder="$t('请输入')"/>
                </el-form-item>
                <!--            <template v-if="store.AdminInfo.AppID === 'admin'">
                                <el-form-item :label="$t('商户名称')"  prop="OperatorId">
                                    <el-select v-model="editForm.OperatorId" clearable :placeholder="$t('请选择')" @change="selectMenuTreeListr('edit')">
                                        <el-option v-for='item in operatorList' :key="item.Id" :label="item.Name" :value="item.Id" />
                                    </el-select>
                                </el-form-item>
                            </template>-->
                <el-form-item :label="$t('设置权限')">
                    <el-tree ref="editTreeRef" :data="editMenuTreeList" show-checkbox accordion node-key="ID" :props="{
                    children: 'Children',
                    label: 'Title',
                }"/>
                </el-form-item>
                <el-form-item :label="$t('备注')">
                    <el-input v-model="editForm.Remark" :placeholder="$t('请输入')"/>
                </el-form-item>
            </el-form>

            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="editPowerDetails">{{ $t('保存') }}</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick} from 'vue'
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia/index';
import type {FormInstance, FormRules, ElTree} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {AdminGroup} from '@/api/adminpb/group';
import {AdminInfo} from "@/api/adminpb/info";

let {t} = useI18n()
const store = useStore()
let groupList = ref([])
let loading = ref(false)
let addDialog = ref(false)
const addFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()

import customTable from "@/components/customTable/tableComponent.vue"
import {ElMessageBox} from "element-plus";

let changeForm = reactive({
    Name: '',
    Remark: '',
    MenuIds: [],
})
let editForm = reactive({
    Name: '',
    PermissionId: '',
    Remark: '',
    MenuIds: [],
})

let tableHeader = [
    {
        label: "ID",
        value: "Id",
        width: "150"
    },

    {
        label: "权限名称",
        value: "Name",
    },


    {
        label: "操作",
        value: "Operator",
        type: "custom",
        width: "180px",
        hiddenVisible: true
    }
]

const MenuIdsCheck = (rule: any, value: any, callback: any) => {
    let arr = treeRef.value!.getCheckedKeys(false)
    if (arr.length === 0) {
        callback(new Error(t('请给商户管理员配置权限')))
    } else {
        callback()
    }
}

const editMenuIdsCheck = (rule: any, value: any, callback: any) => {
    let arr = editTreeRef.value!.getCheckedKeys(false)
    if (arr.length === 0) {
        callback(new Error(t('请给商户管理员配置权限')))
    } else {
        callback()
    }
}
const changeRules = reactive<FormRules>({
    Name: [{required: true, message: t('请填写权限名称'), trigger: 'blur'}],
    OperatorId: [{required: true, message: t('请选择商户'), trigger: 'blur'}],
    MenuIds: [{required: true, validator: MenuIdsCheck, trigger: 'change'}],
})
const editRules = reactive<FormRules>({
    Name: [{required: true, message: t('请填写权限名称'), trigger: 'blur'}],
    OperatorId: [{required: true, message: t('请选择商户'), trigger: 'blur'}],
    MenuIds: [{required: true, validator: editMenuIdsCheck, trigger: 'change'}],
})

let setMenuDialog = ref(false)
const treeRef = ref<InstanceType<typeof ElTree>>()
const editTreeRef = ref<InstanceType<typeof ElTree>>()
let menuTreeList = reactive([])
let arrMenuTreeList = ref([])
let editMenuTreeList = ref([])
let editPowerID = 0
let checkedArr: number[] = []
let operatorList = ref([])
let params = ref({
    PageIndex: 1,
    PageSize: 10,
    Count: 0
})


const init = () => {
    getGroupList()
    menuTreeList = JSON.parse(JSON.stringify(store.AdminInfo.MenuList))
    arrMenuTreeList.value = JSON.parse(JSON.stringify(store.AdminInfo.MenuList))
    if (store.AdminInfo.AppID === 'admin') {
        getOperator()
    } else {
        const getAllIds = (arr) => {
            const ids = [];
            for (const item of arr) {
                ids.push(item.ID);
                if (item.Children) {
                    ids.push(...getAllIds(item.Children));
                }
            }
            return ids;
        };

        changeForm.MenuIds = getAllIds(store.AdminInfo.MenuList)
    }
}

onMounted(() => {
    init()
});

const getGroupList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminGroup.GroupList, {
        PageIndex: params.value.PageIndex,
        PageSize: params.value.PageSize
    })
    loading.value = false
    if (err) {
        return tip.e(err)
    }

    params.value.Count = data.AllCount

    groupList.value = data.AllCount === 0 ? [] : data.List
}

function filterTree(data, ids) {
    return data.map(item => {
        const clonedItem = {...item};
        const children = item.Children;
        if (children) {
            clonedItem.Children = filterTree(children, ids);
        }
        return ids.includes(item.ID) || (clonedItem.Children && clonedItem.Children.length > 0) ? clonedItem : null;
    }).filter(Boolean);
}

const addDialogClick = () => {
    addDialog.value = true

    init()

    changeForm.Name = ""
    changeForm.Remark = ""

    if (store.AdminInfo.AppID != 'admin') {
        const arr = changeForm.MenuIds
        arrMenuTreeList.value = filterTree(menuTreeList, arr)

    }
}

// 切换商户用的
const selectMenuTreeListr = (bools) => {
    let arr = null
    let arreditForm = null
    if (bools === 'add') {
        arr = operatorList.value.find((t) => t.Id === changeForm.OperatorId)?.MenuIds || []
        arreditForm = operatorList.value.find((t) => t.Id === editForm.OperatorId)?.MenuIds || []
        nextTick(() => {
            arrMenuTreeList.value = filterTree(menuTreeList, arr)
        })
    } else {
        arr = editForm.MenuIds || []
        arreditForm = editForm.MenuIds || []
        nextTick(() => {
            editMenuTreeList.value = filterTree(menuTreeList, arreditForm)
        })
    }
}
const AddMenu = async () => {
    if (!addFormRef) return


    changeForm.MenuIds = [
        ...treeRef.value!.getCheckedKeys(false),
        ...treeRef.value!.getHalfCheckedKeys()
    ]

    await addFormRef.value.validate(async (valid, fields) => {
        if (valid) {


            ElMessageBox.confirm(
                t('确认添加该权限'),
                t('是否确认'),
                {
                    confirmButtonText: t('确定'),
                    cancelButtonText: t('关闭'),
                    type: 'warning',
                }
            )
                .then(async () => {
                    let [data, err] = await Client.Do(AdminGroup.AddGroup, changeForm)
                    if (err) {
                        return tip.e(err)
                    }
                    tip.s(t('添加成功'))
                    addDialog.value = false
                    getGroupList()
                })
        }
    })
}

const resetForm = () => {
    // changeForm.Name = ''
    // changeForm.Remark = ''
    // if (!changeFormRef) return
    // changeFormRef.value.resetFields()
}

function removeParentIds(list, arr) {
    let newArr = arr.slice(); // 创建一个新数组，避免修改原始数组

    for (let item of list) {
        if (newArr.includes(item.ID)) {
            // 当前项的ID在newArr中
            if (item.Pid !== 0 && newArr.includes(item.Pid)) {
                // 当前项有父级ID且父级ID在newArr中，则移除父级ID
                newArr = newArr.filter(id => id !== item.Pid);
            }
        }
        if (item.Children && item.Children.length > 0) {
            // 递归处理子集
            newArr = removeParentIds(item.Children, newArr);
        }
    }

    return newArr;
}

const getPowerDetails = (params) => {
    menuTreeList = store.AdminInfo.MenuList
    setMenuDialog.value = true
    Object.assign(editForm, params);
    console.log(editForm);
    editForm.OperatorId = store.AdminInfo.AppID === 'admin' ? params.OperatorId : store.AdminInfo.OperatorId
    editForm.OperatorName = store.AdminInfo.AppID === 'admin' ? operatorList.value.find((t) => t.Id === params.OperatorId)?.Name : store.AdminInfo.OperatorName
    editForm.PermissionId = params.Id
    editMenuTreeList.value = menuTreeList
    // let idList = [];
    let checkedKeys = removeParentIds(editMenuTreeList.value, params.MenuIds)
    // console.log(a);
    nextTick(() => {
        editTreeRef.value!.setCheckedKeys(checkedKeys, false)
    })
}

function findParentIds(list, arr, idList) {
    for (let item of list) {
        if (arr.includes(item.ID)) {
            // 当前项的ID在arr中
            if (item.Pid !== 0 && !idList.includes(item.Pid)) {
                // 当前项有父级ID且父级ID不在idList中，则添加到idList中
                idList.push(item.Pid);
                // 递归查找父级ID的父级ID
                findParentIds(list, [item.Pid], idList);
            }
        }
        if (item.Children && item.Children.length > 0) {
            // 递归处理子集
            findParentIds(list, arr, idList);
        }
    }
}

const editPowerDetails = async () => {


    editForm.MenuIds = [
        ...editTreeRef.value!.getCheckedKeys(false),
        ...editTreeRef.value!.getHalfCheckedKeys()
    ]
    await editFormRef.value.validate(async (valid, fields) => {
        if (valid) {
            ElMessageBox.confirm(
                t('确认修改该权限'),
                t('是否确认'),
                {
                    confirmButtonText: t('确定'),
                    cancelButtonText: t('关闭'),
                    type: 'warning',
                }
            )
                .then(async () => {
                    let [data, err] = await Client.Do(AdminGroup.GetPowerDetailsById, editForm)
                    if (err) {
                        return tip.e(err)
                    }
                    tip.s(t('保存成功'))
                    setMenuDialog.value = false
                    getGroupList()
                })
        }
    })

    /* let checkedKeys = treeRef.value.getCheckedKeys()
     let hafCheckedKeys = treeRef.value.getHalfCheckedKeys()
     let PowerList = checkedKeys.concat(hafCheckedKeys);
     let [data, err] = await Client.Do(AdminGroup.EditPowerDetails, {
         PowerList: <number[]>PowerList,
         Gid: editPowerID
     })
     if (err) {
         return tip.e(err)
     }
     tip.s(t('保存成功'))
     setMenuDialog.value = false*/
}
const getOperator = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, {
        PageIndex: 1,
        PageSize: 10000
    })
    if (err) {
        return tip.e(err)
    }
    operatorList.value = data.AllCount === 0 ? [] : data.List
}

</script>
<style scoped lang='scss'></style>
