<template>

    <div>
        <!-- 查询条件 -->
        <div class="searchView">
            <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
                <el-form-item :label="$t('公告标题')">
                    <el-input v-model="searchInfo.notifyTitle" placeholder="请输入公告标题" clearable/>
                </el-form-item>
                <el-form-item :label="$t('公告类型')">
                    <el-select
                        v-model="searchInfo.notifyType"
                        :placeholder="$t('请选择查询类型')"
                        style="width: 150px"
                        clearable
                    >
                        <el-option v-for="key in notifyType" :label="notifyType[key]" :value="key"/>
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('状态')">

                    <el-select
                        v-model="searchInfo.notifyStatus"
                        :placeholder="$t('请选择公告状态')"
                        style="width: 150px"
                        clearable
                    >

                        <el-option v-for="key in notifyStatus" :label="notifyStatus[key]" :value="key"/>
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('创建时间')">
                    <el-date-picker
                        v-model="searchInfo.notifyCreateTime"
                        type="date"
                        :placeholder="$t('请选择创建时间')"
                        clearable
                    />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="onSearchNotify">{{ $t('查询') }}</el-button>
                </el-form-item>
            </el-form>

        </div>

        <!-- 表格数据 -->
        <div class="page_table_context">

            <el-tabs v-model="activeName" class="demo-tabs" @tabChange="languageChange">
                <el-tab-pane :label="language.lanNameNew" :name="language.abbreviation"
                             v-for="(language, index) in languageList" :key="index"/>
            </el-tabs>


            <div class="flex flex_child_end">
                <el-button type="primary" @click="addDialog = true">新增</el-button>
            </div>

            <customTable :table-header="tableHeader" :table-data="tableData" :page="searchInfo.page"
                         :pageSize="searchInfo.pageSize" :count="searchInfo.count">

                <template #operation="scope">
                    <el-button type="primary" plain @click="editNotify(scope)">编辑</el-button>
                    <el-button type="danger" plain @click="deleteNotify(scope)">删除</el-button>
                </template>
                <template #notifyStatus="scope">
                    <el-button type="primary" plain v-if="scope.scope.notifyStatus == 1" size="small" style=" pointer-events: none;">启用</el-button>
                    <el-button type="danger" plain v-if="scope.scope.notifyStatus == 2" size="small" style=" pointer-events: none;">停用</el-button>
                </template>

            </customTable>
        </div>


        <!-- 新增修改弹窗       -->

        <el-dialog v-model="addDialog" :title="isEdit ? '修改公告' : '添加公告'" width="60%" center>

            <el-form :model="dataForm" label-width="auto" style="max-width: 600px">
                <el-form-item :label="$t('公告标题')">
                    <el-input v-model="dataForm.notifyTitle" placeholder="请输入公告标题" clearable/>
                </el-form-item>
                <el-form-item>
                    <el-form-item :label="$t('公告类型')">
                        <el-select
                            v-model="dataForm.notifyType"
                            :placeholder="$t('请选择查询类型')"
                            style="width: 150px"
                            clearable
                        >
                            <el-option v-for="(item, key) in notifyType" :value="key" :label="item">{{ item }}</el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('状态')">
                        <el-select
                            v-model="dataForm.notifyStatus"
                            :placeholder="$t('请选择公告状态')"
                            style="width: 150px"
                            clearable
                        >
                            <el-option v-for="(item, key) in notifyStatus" :value="key" :label="item">{{ item }}</el-option>
                        </el-select>
                    </el-form-item>
                </el-form-item>

                <el-form-item :label="$t('上下架时间')">
                    <el-date-picker
                        v-model="dataForm.limitTime"
                        type="datetimerange"
                        start-placeholder="Start date"
                        end-placeholder="End date"
                        format="YYYY-MM-DD HH:mm:ss"
                        date-format="YYYY/MM/DD ddd"
                        time-format="A hh:mm:ss"
                    />
                </el-form-item>
                <el-form-item>
                    <el-tabs v-model="formActiveName" class="demo-tabs" @tabChange="formLanguageChange">
                        <el-tab-pane :label="language.lanNameNew" :name="language.abbreviation"
                                     v-for="(language, index) in languageList" :key="index"/>
                    </el-tabs>



                    <el-input
                        v-model="notifyContext"
                        style="width: 240px"
                        :autosize="{ minRows: 2, maxRows: 4 }"
                        type="textarea"
                        placeholder="Please input"
                        @change="notifyContextChange"
                    />

                </el-form-item>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="addDialog = false">关闭</el-button>
                    <el-button type="primary" @click="commitNotify">
                        {{ $t('提交') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>


<script setup lang="ts">

import {useI18n} from "vue-i18n";

import {Notify} from "@/api/adminpb/notify";
import {Client} from "@/lib/client";
import {AdminGroup} from "@/api/adminpb/group";
import ut from "@/lib/util";
import {ref, watch} from "vue";
import customTable from "@/components/customTable/tableComponent.vue"
import {ElMessage} from "element-plus";

// 国际化语言
let {t} = useI18n()

let tableHeader = [
    {
        label: "序号",
        value: "",
        type: "index",
    },
    {
        label: "公告类型",
        value: "notifyType",
        format: (row) => notifyType[row.notifyType]
    },
    {
        label: "公告标题",
        value: "notifyTitle",
    },
    {
        label: "上架时间",
        value: "",
        format:(row) => `${ut.fmtDateSecond(row.onlineLimitTime)}`,
        width: "180px"
    },
    {
        label: "下架时间",
        value: "",
        format:(row) => `${ut.fmtDateSecond(row.lowerLimitTime)}`,
        width: "180px"
    },
    {
        label: "创建时间",
        value: "created_time",
    },
    {
        label: "排序",
        value: "sort",

    },
    {
        label: "状态",
        value: "notifyStatus",
        type: "custom",
        format: (row) => notifyStatus[row.notifyStatus]
    },
    {
        label: "操作者",
        value: "operator",
    },
    {
        label: "操作",
        value: "operation",
        type: "custom",
    }
]
let tableData = ref([])


// 公告类型
const notifyType = {
    "1": "系统",
    "2": "维护",
    "3": "优惠",
    "4": "推荐"
}
// 状态
const notifyStatus = {
    "1": "启用",
    "2": "停用",
}

// 当前系统中存在的语言类别
const languageList = ref([])

// 语言类别选中的值
const activeName = ref("")          // 查询条件中
const formActiveName = ref("")      // 语言类别中

// 查询条件
const searchInfo = ref({
    notifyTitle: "",
    notifyType: "",
    notifyStatus: "",
    notifyCreateTime: "",
    page: 1,
    pageSize: 20,
    count: 0
})

// 数据实体
const dataForm = ref({
    notifyTitle: "",
    notifyType: "",
    notifyStatus: "",
    notifyCreatedTime: "",
    limitTime:[],
    sort: 0,
    onlineLimitTime: "",
    lowerLimitTime: "",
    languageContext: {}
})
const notifyLanguageConfig = ref({

})

// 弹窗弹出的状态
const addDialog = ref(false)
const isEdit = ref(false)


const notifyContext = ref("")


watch(addDialog, (newData) => {
    if (!newData){
        dataForm.value = {
            notifyTitle: "",
            notifyType: "",
            notifyStatus: "",
            notifyCreatedTime: "",
            limitTime:[],
            sort: 0,
            onlineLimitTime: "",
            lowerLimitTime: "",
            languageContext: {}
        }
        formActiveName.value = languageList.value[0].abbreviation
        dataForm.value.languageContext = {}
        notifyLanguageConfig.value = {}
        notifyContext.value = ""
        isEdit.value = false
    }

})


// 获取通知信息列表
const getNotifyData = async () => {

    let response = await Client.Do(Notify.GetNotify, {languageCode: activeName.value})
    tableData.value = response[0].Notifications
    searchInfo.value.count = response[0].Count
}

// 获取系统现有语言类别
const getLanguageData = async () => {
    let response = await Client.Do(AdminGroup.GetLang, {})
    let langCodeList = Object.keys(response[0].List[0] || {}).filter(list => list !== 'Permission' && list !== '_id');
    languageList.value = ut.LangList.filter(item => langCodeList.indexOf(item.abbreviation.toLowerCase()) > -1)
    activeName.value = languageList.value[0].abbreviation
    formActiveName.value = languageList.value[0].abbreviation
}

// 页面初始化
const init = async () => {

    await getLanguageData()
    await getNotifyData()
}
init()

// 查询通告
const onSearchNotify = () => {
    getNotifyData()
}
const languageChange = () => {
    getNotifyData()
}

const formLanguageChange = () => {
    notifyContext.value = notifyLanguageConfig.value[formActiveName.value] || ""
}

const notifyContextChange = () => {
    notifyLanguageConfig.value[formActiveName.value] = notifyContext.value
}


const commitNotify = async () => {

    dataForm.value.onlineLimitTime  =  dataForm.value.limitTime[0]
    dataForm.value.lowerLimitTime   =  dataForm.value.limitTime[1]

    let languageContext = {}
    for (const i in languageList.value) {
        let context = ""
        if(notifyLanguageConfig.value[languageList.value[i].abbreviation]){
            context = notifyLanguageConfig.value[languageList.value[i].abbreviation]
        }
        languageContext[languageList.value[i].abbreviation] = context
    }

    dataForm.value.languageContext   =  JSON.stringify(languageContext)


    let response
    if (isEdit.value){

        response = await Client.Do(Notify.EditNotify, dataForm.value as any)
          if (!response[1]){
              ElMessage.success("修改成功")
          }else {
              ElMessage.success("修改失败")
          }
    }else{

        response = await Client.Do(Notify.AddNotify, dataForm.value as any)
          if (!response[1]){
              ElMessage.success("添加成功")
          }else {
              ElMessage.success("添加失败")
          }
    }


    addDialog.value = false
    getNotifyData()
}

const deleteNotify = async (scope) => {

    const response = await Client.Do(Notify.DelNotify, {_id: scope.scope._id} as any)

    if (!response[1]){
        ElMessage.success("删除成功")
        getNotifyData()
    }else{
        ElMessage.success("删除失败")
    }
}

const editNotify = (data) => {

    isEdit.value = true
    addDialog.value = true
    dataForm.value = data.scope
    dataForm.value.limitTime = [dataForm.value.onlineLimitTime, dataForm.value.lowerLimitTime]
    languageList.value.forEach(item =>{
        notifyLanguageConfig.value[item.abbreviation] = data.scope[item.abbreviation] || ""
    })

    notifyContext.value = data.scope[languageList.value[0].abbreviation]


}

</script>


<style scoped lang="scss">
.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

.flex_child_end {
  margin-bottom: 10px;
}
</style>
