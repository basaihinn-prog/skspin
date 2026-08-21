/*
 *@Author: 西西米
 *@Date: 2022-12-23 17:14:13
 *@Description: 菜单列表
*/
<template>
    <div class='menuList'>

        <!-- 数据 -->


        <div class="page_table_context">

            <customChildTable
                table-name="menuList_list"
                v-loading="loading"
                rowKey="ID"
                :treeProps="{ children: 'Children' }"
                :table-data="menuList"
                :table-header="tableHeader"
                @refreshTable="getMenuList">
            <template #handleTools>
                <el-button type="primary" @click="addDialog = true">{{ $t('添加菜单') }}</el-button>
                <el-button type="success" @click="copyJson">{{ $t('复制为JSON') }}</el-button>

            </template>
            <template #Title="scope">
                <div style="display: flex;align-items: center;">
                    <el-icon>
                        <component :is="scope.scope.MenuIconFont" />
                    </el-icon>&nbsp;
                    {{ scope.scope.Title }}
                </div>
            </template>

            <template #Operator="scope">

                <el-button type="primary" size="small" plain @click="addForm.Pid = scope.scope.ID, addDialog = true">{{ $t('新增')
                    }}</el-button>
                <el-button type="primary" size="small" plain @click="revise(scope.scope)">{{ $t('编辑') }}</el-button>

                <el-button type="danger" size="small" @click="DelMenu(scope.scope.ID)" plain>{{ $t('删除') }}</el-button>

            </template>
        </customChildTable>
        </div>
        <!-- 添加弹框 -->
        <el-dialog v-model="addDialog" @closed="resetForm" :title="$t('菜单')" :width="store.viewModel === 2 ? '85%' : '1200px'">
            <el-form ref="addFormRef" :model="addForm" :rules="addRules" label-position="top" label-width="80px">

                <el-form-item :label="$t('上级菜单')" v-if="reviseID">
                    <el-select v-model="addForm.Pid">
                        <el-option label="无" :value="0"/>
                        <el-option v-for="menu in menuList" :label="menu.Title" :value="menu.ID"/>
                    </el-select>

                </el-form-item>
                <el-form-item :label="$t('中文名称')" prop="Title">
                    <el-input v-model="addForm.Title" :placeholder="$t('请输入')" />
                </el-form-item>
                <el-form-item :label="$t('前端路由')" prop="Url">
                    <el-input v-model="addForm.Url" :placeholder="$t('请输入')" />
                </el-form-item>
                <el-form-item :label="$t('图标')">
                    <div class="iconView">
                        <div class="icon" v-for="(item, index) in icon" :key="index"
                            @click="addForm.Icon = item.name"
                            :style="{ 'background': addForm.Icon === item.name ? '#67c23a' : '#4478f9' }">
                            <el-icon color="#ffffff" size="20">
                                <component :is="item" />
                            </el-icon>
                        </div>
                    </div>
                </el-form-item>
                <el-form-item :label="$t('排序')">
                    <el-input-number v-model="addForm.Sort" :precision="0" :step="1" :min="0" />
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddMenu">{{ $t('保存') }}</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive, nextTick, watch } from 'vue'
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { MenuListToTree, useStore } from '@/pinia/index';
import type { FormInstance, FormRules } from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n';
import { AdminMenu } from '@/api/adminpb/menu';

import customChildTable from "@/components/customTable/tableChildrenComponent.vue"

import ut from '@/lib/util';
import {Delete, Edit, Plus, RefreshRight, Setting} from "@element-plus/icons-vue";
import lashenIcon from "@/assets/login/lashen.png";
import {ElMessageBox} from "element-plus";
import TableChildrenComponent from "@/components/customTable/tableChildrenComponent.vue";
const { t } = useI18n()
const store = useStore()

const tableBorder = ref(true)
const visibleTableHeader = ref(true)
const lineStatus = ref("large")

let menuList = ref([])
let tableHeader = [
    {
        label:"ID",
        value:"ID",
        width:"150"
    },
    {

        label:"菜单名称",
        value:"Title",
    },
    {
        label:"前端路由",
        value:"Url",
    },
    {
        label:"排序",
        value:"Sort",

    },
    {
        label:"备注",
        value:"Remark",
    },
    {
        label:"操作",
        value:"Operator",
        type: "custom",
    }
]


let loading = ref(false)
const icon = ElementPlusIconsVue
onMounted(() => {
    getMenuList()
    let tableSetConfig = store.SystemConfig.tableBorderConfig

    if (tableSetConfig){
        tableBorder.value = tableSetConfig["menuList_list"]

    }
});


watch(tableBorder, (newData)=>{


    let tableSetConfig = store.SystemConfig.tableBorderConfig
    tableSetConfig["menuList_list"] = newData
    localStorage.setItem("tableSetConfig", JSON.stringify(tableSetConfig))


})


const copyJson = async () => {
    let [data, err] = await Client.Do(AdminMenu.MenuList, {})
    if (err) {
        return tip.e(err)
    }

    ut.copyTextToClipBoard(JSON.stringify(data.MenuList, null, "\t"))
}

const getMenuList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminMenu.MenuList, {})
    loading.value = false
    if (err) {
        return tip.e(err)
    }


    menuList.value = MenuListToTree(data)
}

let addDialog = ref(false)
const addFormRef = ref<FormInstance>()
let addForm = reactive({
    Pid: 0,
    Title: '',
    Url: '',
    Icon: 'Setting',
    Sort: 0,
})
const addRules = reactive<FormRules>({
    Title: [{ required: true, message: t('请填写菜单名'), trigger: 'blur' }],
    Url: [{ required: true, message: t('请填写前端路由'), trigger: 'blur' }],
})

const changeTableLine = (changStatus) => {
    lineStatus.value = changStatus
}
const AddMenu = async () => {

    if (!addFormRef) return
    await addFormRef.value.validate(async (valid, fields) => {
        if (valid) {
        ElMessageBox.confirm(
            t(`是否确定${reviseID.value ? '修改': '新增'}菜单`),
            t('是否确认'),
            {
                confirmButtonText: t('确定'),
                cancelButtonText: t('关闭'),
                type: 'warning',
            }
        )
            .then(async () => {

                    let url = reviseID.value ? AdminMenu.UpdateMenu : AdminMenu.AddMenu
                    let Form = addForm as any
                    if (reviseID.value) {
                        Form.ID = reviseID.value
                    }
                    let [data, err] = await Client.Do(url, Form)
                    if (err) {
                        return tip.e(err)
                    }
                    tip.s(t('操作成功'))
                    addDialog.value = false
                    getMenuList()

            })
        }
    })

}

const DelMenu = async (ID) => {


    ElMessageBox.confirm(
        t(`是否确定删除该菜单`),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {

                let [data, err] = await Client.Do(AdminMenu.DelMenu, {ID})
                if (err) {
                    return tip.e(err)
                }
                tip.s(t('操作成功'))
                getMenuList()
            }
        )
}

let reviseID = ref(null)
const revise = (data) => {
    for (const key in addForm) {

        if (!addForm[key]){
            addForm[key] = {}
        }
        addForm[key] = data[key]
    }
    reviseID.value = data.ID
    addDialog.value = true
}

const resetForm = () => {
    nextTick(() => {
        addForm.Pid = 0
        addForm.Title = ''
        addForm.Sort = 0
        reviseID.value = null
    })
    if (!addFormRef) return
    addFormRef.value.resetFields()
}
</script>
<style scoped lang='scss'>
.iconView {
    display: flex;
    flex-wrap: wrap;
    row-gap: 10px;
    column-gap: 10px;

    .icon {
        width: 25px;
        height: 25px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
    }
}

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

.tableHandleSwitch{
    height: 30px;
    display: flex;
    align-items: center;

}
.tableHandleSwitch>span{
    margin-right: 5px;
}

</style>
