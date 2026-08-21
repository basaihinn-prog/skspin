<template>
    <div class='groupList'>
        <div class="searchView">
            <el-form
                label-position="top" label-width="100px" :model="param" style="max-width: 100%">
                <el-space wrap>
                    <el-input v-model="param.Lang" clearable></el-input>
                    <el-button type="primary" @click="getLangList()">{{ $t('搜索') }}</el-button>
                    <el-upload
                        v-model:file-list="fileList"
                        class="upload-excel"
                        :action="doUpload"
                        :http-request="ossUpload"
                        auto-upload
                        :limit="3"
                    >
                        <el-button type="primary">导入</el-button>

                    </el-upload>
                </el-space>
            </el-form>

        </div>


        <!--  表格数据      -->

        <div class="page_table_context">
            <div class="flex flex_child_end">
                <el-button type="primary" @click="addDialog = true">新增</el-button>
            </div>

            <customTable :table-header="tableHeader" :table-data="langList" :page="param.PageIndex"
                     :pageSize="param.PageSize" :count="param.Count">

            <template #Permission="scope">
                <el-space wrap>
                    <el-tag v-if="scope.scope.Permission == '通用'" type="success" disable-transitions>通用</el-tag>
                    <el-tag v-else-if="scope.scope.Permission == '游戏'" type="info" disable-transitions>游戏</el-tag>
                    <el-tag v-else type="danger" disable-transitions>{{ scope.scope.Permission }}</el-tag>
                </el-space>
            </template>
            <template #operation="scope">
                <el-button type="primary" plain @click="openDialogEvent(scope.scope,'edit')">编辑</el-button>
                <el-button type="danger" plain @click="deleteEvent(scope.scope)">删除</el-button>
            </template>


        </customTable>
        </div>

        <!--   添加     -->
        <el-dialog v-model="addDialog" destroy-on-close title="添加多语言" @close="resetForm()"
                   width="680">
           <div>
               <el-radio-group v-model="addDialogList.radio" class="ml-4">
                   <el-radio :label="0" size="large">通用</el-radio>
                   <el-radio :label="-1" size="large">其他</el-radio>
               </el-radio-group>
               <el-checkbox-group v-if="addDialogList.radio !== 0" v-model="addDialogList.checkList">
                   <el-checkbox v-for="item in addDialogList.list" :label="item.label">{{ item.name }}</el-checkbox>
               </el-checkbox-group>
           </div>
            <div :class="showAddLangInput?'div_border':'div_margin_b'">
                <el-space>
                    <el-button type="primary" plain  @click="addLang()" style="margin-bottom: 0" size="small">
                        {{ showAddLangInput?'取消新增':'新增语言' }}</el-button>
                    <template v-if="showAddLangInput">
                        <el-select
                            v-model="autocompleteState"
                            filterable size="small"
                            placeholder="下拉选择语言"
                            @change="handleSelect"
                        >
                            <el-option
                                v-for="item in restaurants"
                                :label="item.lanName"
                                :value="item.abbreviation"
                                :disabled="!item.bol"
                            >
                                <el-text style="float: left" :class="item.bol?'':'itemBol'" size="small">{{ item.lanName }}</el-text>
                                <el-text style="float: right" :class="item.bol?'':'itemBol'" type="info" size="small">{{ item.abbreviation }}</el-text>
                            </el-option>
                        </el-select>
                    </template>
                </el-space>
                <div>
                    <el-text class="" size="small">每次添加只可添加一种语言</el-text>
                </div>
            </div>
            <div class="div_border">
                <div>
                    <el-text class="" size="small">粘贴格式为:英文，泰文，中文，（所添加语言）...</el-text>
                </div>
                <div>
                    <el-text class="" size="small">泰文，中文为必填项</el-text>
                </div>
                <div>
                    <el-text class="" size="small">tips: 使用英文逗号将不同语言隔开，最后一列语言数据右侧没有逗号</el-text>
                </div>
                <el-input
                    v-model="addFormText"
                    placeholder="添加多语言"
                    show-word-limit
                    type="textarea"
                />
            </div>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="addLangList">添加</el-button>
                </span>
            </template>
        </el-dialog>
        <!-- 添加弹框 -->
        <el-dialog v-model="openDialog" destroy-on-close :title="whitchButton?$t('修改多语言1'):$t('添加多语言1')" @close="resetForm()" width="550px">
            <el-form ref="langFormRef" :model="changeForm" :rules="changeRules" label-position="top" label-width="80px">
                <el-form-item>
                    <el-space wrap>
                        <el-radio-group v-model="changeForm.radio" class="ml-4" size="large"
                                        @change="changeFormRadioChange">
                            <el-radio-button :label="0" >通用</el-radio-button>
                            <el-radio-button :label="-1">其他</el-radio-button>
                        </el-radio-group>
                        <div :class="changeForm.radio !== 0?'div_border':''" style="margin-bottom: 0">
                            <el-checkbox-group v-if="changeForm.radio !== 0" v-model="changeForm.checkList">
                                <el-checkbox v-for="item in addDialogList.list" :label="item.label">{{
                                        item.name
                                    }}
                                </el-checkbox>
                            </el-checkbox-group>
                        </div>
                    </el-space>
                </el-form-item>
                <template v-for="(items,index) in Object.keys(changeForm)" >
                    <el-form-item :label="tableListLabelFun(items)" :prop="items"
                                  v-if="items !== 'Permission' &&
                                  items !== 'checkList' &&
                                  items !== 'radio'
">
                        <el-input v-model="changeForm[items]" :placeholder="$t('请输入')" :disabled="items==='ZH'?whitchButton:false"/>
                    </el-form-item>
                </template>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddMenu">{{ whitchButton ? $t('修改') : $t('添加') }}</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, reactive, ref, watch} from 'vue'
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia/index';
import type {FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {AdminGroup} from '@/api/adminpb/group';
import ut from "@/lib/util";

import customTable from "@/components/customTable/tableComponent.vue"


// 表格头
const tableHeader = ref([])

let {t} = useI18n()
const store = useStore()
let langList = ref([])
let loading = ref(false)
let openDialog = ref(false)
let addDialog = ref(false)
const langFormRef = ref<FormInstance>()
let whitchButton = ref(false)
let changeForm = ref({
    radio: 0,
    checkList: [],
})



let doUpload = "http://127.0.0.1"
let fileList = ref([])
const ossUpload = () => {

    let formData = new FormData();
    formData.append('file', fileList.value[0].raw);
    Client.Do(AdminGroup.UploadLang, formData)
    fileList.value = []
}

const addDialogList = reactive({
    radio: 0,
    checkList: [],
    list: [
        {
            name: t('游戏'),
            label: 1 << 0,
        },
        {
            name: t('彩票'),
            label: 1 << 1
        }
    ]
})
let param = ref({
    Lang: '',
    Count:"",
    PageIndex: 1,
    PageSize: 20
})

const filterTag = (value, row) => {
    return row.Permission === value
}
let Count = ref(0)
const changeRules = reactive<FormRules>({
    ZH: [{required: true, message: t('请填写中文名称'), trigger: 'blur'}],
})
const autocompleteState = ref('')
const restaurants = ref([])

watch(addDialog, (newData)=>{
    if (newData){

        let exitsLange = restaurants.value.filter(item => !item.bol).map(item => item.lanName)

        addFormText.value = exitsLange.join(",")
    }
})

const handleSelect = () => {
    let arr = []
    arr.push(autocompleteState.value)
    if (arr.length>1 && arr.includes(autocompleteState.value)) {
        tip.e('每次添加只可添加一种语言')
        return
    }
    changeForm.value.newArrs.push({
        [arr.toString().toUpperCase()]:''
    })
}
const showAddLangInput = ref(false)
const addLangvalue = ref('')

const newArrsLen = ref(0)

const addLang = ()=> {
    showAddLangInput.value = !showAddLangInput.value
    if (!showAddLangInput.value && changeForm.value.newArrs.length>newArrsLen.value) {
        changeForm.value.newArrs.pop()
    }
}

onMounted(async () => {
    await getLangList()
});
const tableListLabel = ref('')
let firstLangItemKeys = ref([])
const tableListLabelFun = (i) => {

    return  ut.LangList.find(list => list.abbreviation.toUpperCase() === i.toString())?.lanName || i
}
const getLangList = async () => {
    tableHeader.value = []
    loading.value = true
    let [data, err] = await Client.Do(AdminGroup.GetLang, param.value)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    param.value.Count = data?.Count
    langList.value = data.Count === 0 ? [] : data.List

    for (const key in langList.value[0]) {
        if (key == "_id"){
            break
        }

        if (tableListLabelFun(key) == "Permission"){
            tableHeader.value.push({
                label: "类型",
                value: key,
                type: "custom"
            })
        }else {
            tableHeader.value.unshift({
                label: tableListLabelFun(key),
                value: key,
            })
        }

    }

    const newArrs = firstLangItemKeys.value.map(key => ({ [key.toUpperCase()]: '', }));

    tableHeader.value.push({
        label: "操作",
        value: "operation",
        type: "custom",
        width:"180px"
    })

    changeForm.value = {
        ...changeForm.value,
        newArrs,
        Permission: 0,
    }
    newArrsLen.value = newArrs.length
    let str = newArrs.map(list => Object.keys(list).toString())
    restaurants.value = ut.LangList.map(m=>{
        return {
            ...m,
            bol:!str.find(op => op === m.abbreviation.toUpperCase())
        }
    })

}

const changeFormRadioChange = () => {
    console.log(changeForm.value.radio);
}
const openDialogEvent = (i, str) => {
    const copiedI = {...i};
    changeForm.value = {...i};
    let arr = []
    let permission = i.Permission;
    for (let a = 0; permission !== 0; permission = permission >> 1, a++) {
        if ((permission & 1) === 1) {
            arr.push(a)
        }
    }
    if (changeForm.value['Permission'] === 0) {
        changeForm.value = {
            ...copiedI,
            radio:0
        };
    } else {
        changeForm.value.radio = -1
        let values = arr.map(v => {
            return addDialogList.list.find((i, index) => v === index)
        })
        changeForm.value.checkList = values.map(x => x.label);
    }
    console.log(changeForm.value);
    whitchButton.value = true
    openDialog.value = true
}

const addForm = reactive([])
const addFormText = ref('')

const addList = () => {
    let arr = changeForm.value.newArrs
    arr.splice(0,1)
    let result = arr.map(obj => {
        let label = ""
        try {

            label = ut.LangList.find(list => list.abbreviation.toUpperCase() === Object.keys(obj).toString().toLowerCase()).lanName;
        }catch (e){
            console.log(obj)
        }
        return { label: label, key: Object.keys(obj).toString() };
    })
        .reduce((acc, obj, index) => {
            acc[index] = obj;
            return acc;
        }, {});
    console.log(result);
    const parsedData = addFormText.value.split("\n").map(line => {
        const values = line.split(',');
        return values.map((value, index) => ({
            value,
            ...result[index]
        }));
    });
    addForm.splice(0, addForm.length, ...parsedData);
}
const deleteEvent = async (i) => {
    // deleteEvent
    loading.value = true
    let [data, err] = await Client.Do(AdminGroup.LangDelete, {Id: i._id} as any)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    tip.s(t('删除成功'))
    await getLangList()
}
const paginationEmit = (PageIndex, PageSize) => {
    param.value.PageIndex = PageIndex
    param.value.PageSize = PageSize
    getLangList()
}

const AddMenu = async () => {
    let Permission = 0
    if (changeForm.value.radio !== 0) {
        changeForm.value.checkList.forEach(t => {
            Permission = Permission | t
        })
    }

    changeForm.value.Permission = Permission
    if (whitchButton.value) {
        loading.value = true
        let [data, err] = await Client.Do(AdminGroup.LangUpdate, changeForm.value)
        loading.value = false
        if (err) {
            return tip.e(err)
        }
        tip.s(t('修改成功'))
        openDialog.value = false
        await getLangList()
    }
}

function transformToMap(data) {
    const Langs = {};
    let Permission = 0;
    addDialogList.checkList.forEach(t => {
        Permission |= t;
    });
    console.log(data);
    for (const item of data) {
        const zhLang = item.find(list => list.key === 'ZH')?.value;
        const thLang = item.find(list => list.key === 'TH')?.value;


        if (!zhLang || !thLang) {
            console.log(item);
            tip.e(t('请添加中文和泰文'));
            return null;
        }
        if (!Langs[zhLang]) {
            Langs[zhLang] = { Permission };
        }


        for (const arrs of item) {
            if (!arrs.key) {
                tip.e(t('格式出错'));
                return null;
            }
            Langs[zhLang][arrs.key] = arrs.value;
        }
    }

    return Langs;
}

const addLangList = async () => {
    try {
        addList();
        const transformedLangs = transformToMap(addForm);
        if (!transformedLangs) return; // 出现错误时直接返回

        const datas = { Langs: transformedLangs };

        const [data, err] = await Client.Do(AdminGroup.AddLang, datas);

        if (err) throw err; // 将错误抛出，在 catch 块中统一处理

        if (data.ErrLangs.length) {
            const local = localStorage.getItem('defaultLocal')?.toUpperCase() || 'ZH';
            const errorMessage = data.ErrLangs[local] || data.ErrLangs['ZH'] || t('添加语言时出错');
            throw new Error(errorMessage); // 抛出错误信息
        }

        tip.s(t('添加成功'));
        addDialog.value = false;
        await getLangList();
    } catch (err) {
        tip.e(err.message); // 在 catch 块中处理错误信息
    }
}

const resetForm = () => {
    if (!langFormRef) return
    langFormRef.value?.resetFields()
}

const editNotify = (data) => {

}
const deleteNotify = (data) => {

}

</script>
<style scoped lang='scss'>
.div_border{
    border: 1px solid #e1dede;
    padding: .5rem;
    border-radius: 0.4rem;
    margin-bottom: 1rem;
}
.div_margin_b{
    margin-bottom: 1rem;
}
.itemBol{
    color: #cccccc;
}
</style>

<style>
.upload-excel .el-upload-list{
    display: none;
}
.flex_child_end{

    margin-bottom: 15px;
}
</style>
