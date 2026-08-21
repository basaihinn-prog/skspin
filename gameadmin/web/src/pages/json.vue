/*
 *@Author: 西西米
 *@Date: 2023-01-30 14:12:34
 *@Description: json
*/
<template>
    <div v-loading="loading" class='jsonView'>
        <el-row>
            <el-col :lg="4" :sm="24">
                <el-affix target=".treeView">
                    <el-input
                        v-model="filterText"
                        style="margin-bottom: 1rem"
                        placeholder="关键字过滤搜索"
                    />
                </el-affix>
                <div class="treeView" :class="dNone?'d-none':''">
                    <el-tree :data="list" :default-expand-all="false"
                             :filter-node-method="filterNode"
                             :highlight-current="true"
                             accordion
                             @node-click="treeNodeClick"
                             node-key="ID" ref="tree">
                        <template #default="{ node, data }">
                    <span>
                        <span>{{ data.Name }}</span>
                    </span>
                        </template>
                    </el-tree>
                </div>
            </el-col>
            <el-col :lg="20">
                <div class="mainview" v-if="fileEvent">
                    <el-descriptions :column="2" :border="true" size="small" style="margin-bottom: 12px;">
                        <el-descriptions-item :label="$t('名称')">{{ fileEvent.Name }}</el-descriptions-item>
                        <el-descriptions-item :label="$t('JSON文件名')">{{ fileEvent.File }}</el-descriptions-item>
                        <el-descriptions-item :label="$t('高亮')">
                            <el-switch v-model="isDark" :before-change="changeDark" /></el-descriptions-item>
                        <el-descriptions-item :label="$t('操作')">
                            <el-button type="primary" @click="SaveConfig">{{ $t('提交') }}</el-button>
                            <el-button type="primary" @click="isFileHistory = true, getFileHistory()">{{ $t('查看历史') }}</el-button>
                        </el-descriptions-item>
                        <el-descriptions-item :label="$t('历史记录')" v-if="isFileHistory">
                            <div class="fileHistory" v-for="(item, index) in FileHistory" :key="index">
                                <p :class="FileHistoryIndex === index ? 'act' : ''"
                                   @click="FileHistoryIndex = index, fileHistoryClick(item)">{{ item.FileName }} {{
                                        dateFormater(0, 0, item.CreateAt)
                                    }}</p>
                            </div>
                        </el-descriptions-item>
                    </el-descriptions>

                    <div class="codeView">
                        <el-input :rows="30" type="textarea" v-if="!isDark" v-model="json"></el-input>
                        <codemirror v-else-if="fileExt !== 'csv' && fileExt !== ''" v-model="json" :placeholder="$t('请输入')"
                                    :indent-with-tab="true" :extensions="[javascript(), oneDark]" :tabSize="4" style="min-height: 400px;" />
                        <hot-table v-else-if="fileExt === 'csv'" :data="handsontableData" :settings="settings"></hot-table>
                        <div v-else></div>
                    </div>
                </div>
            </el-col>
        </el-row>


</div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive, watch, nextTick } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { AdminConfigFile } from '@/api/adminpb/json';
import { Codemirror } from "vue-codemirror";
import { javascript } from "@codemirror/lang-javascript";
import { oneDark } from "@codemirror/theme-one-dark";
import Papa from "papaparse";

import { HotTable } from '@handsontable/vue3';
import 'handsontable/dist/handsontable.full.css';//表格样式
import 'handsontable/languages/zh-CN'; //汉语包
import beautify from "js-beautify";
import { useI18n } from 'vue-i18n';
const { t } = useI18n()
const dNone = ref(false)
onMounted(() => {
    getList()
    console.log(window.screen.width);
    window.screen.width<500?dNone.value = true:dNone.value = false
});
let loading = ref(false)
let list = ref([])
const getList = async () => {
    let listFile = 'admin_configFiles.csv'
    let data = await getLoadConfig(listFile)

    let arr: any[] = Papa.parse(data, { header: true }).data

    let idx = arr.findIndex((val: any) => val.File == listFile)


    let newList = []
    const filterFn = (name) => list => list.Name.slice(0, 2) === name;
    const pgList = arr.filter(filterFn('pg'));
    const jiliList = arr.filter(filterFn('ji'));
    const newArr = arr.filter(list => !filterFn('pg')(list) && !filterFn('ji')(list));

    newArr.forEach(item => {
        if (item.Name) {
            let newItem = newList.find((i) => i.Name == item.Group)
            if (item.Group === '') {
                newList.push(item)
            } else if (!newItem) {
                newList.push({ Name: item.Group, children: [item] })
            } else {
                newItem.children.push(item)
            }
        }
    })
    console.log(newArr);
    list.value = newList
}

let fileExt = ref('')
let settings = ref({
    language: 'zh-CN', // 官方汉化
    licenseKey: 'non-commercial-and-evaluation', //去除底部非商用声明
    currentRowClassName: 'currentRow', // 突出显示行
    currentColClassName: 'currentCol', // 突出显示列
    // colHeaders: ["Label", "File", "Type", "Group"],
    colHeaders: true,
    trimWhitespace: false, //去除空格
    rowHeaderWidth: 50, //单元格宽度
    stretchH: 'all',
    rowHeaders: true, // 行标题   布尔值 / 数组/  函数
    contextMenu: true, //右键菜单
    manualColumnResize: true,
    autoWrapRow: true, //自动换行
    width: "100%",
    height: "auto",
})
let handsontableData = ref([])

const getLoadConfig = async (FileName) => {
    loading.value = true
    let [data, err] = await Client.Do(AdminConfigFile.LoadConfig, { FileName })
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    return data.Content
}

let fileEvent = ref(null)
const treeNodeClick = async (event) => {
    fileExt.value = ''
    fileEvent.value = event
    FileHistory.value = []
    FileHistoryIndex.value = null
    isFileHistory.value = false
    console.log(event);
    if (event.File) {
        let data = await getLoadConfig(event.File)
        showData(data)
    }
}
const filterText = ref('')
import { ElTree } from 'element-plus'
const tree = ref<InstanceType<typeof ElTree>>()
watch(filterText, (val) => {
    tree.value!.filter(val)
})

const filterNode = (value, data) => {
    if (!value) return true
    return data.Name.includes(value)
}

let FileHistory = ref([])
let isFileHistory = ref(false)
let FileHistoryIndex = ref(null)
const getFileHistory = async () => {
    let [data, err] = await Client.Do(AdminConfigFile.FileHistory, {
        Count: 5,
        FileName: fileEvent.value.File
    })
    if (err) {
        return tip.e(err)
    }
    FileHistory.value = data.List
}

const fileHistoryClick = (data) => {
    fileExt.value = ''
    showData(data.Content)
}

const showData = (data) => {
    nextTick(() => {
        if (fileEvent.value.File.endsWith(".csv")) {
            handsontableData.value = Papa.parse(data).data
        }
        if (fileEvent.value.File.endsWith(".json")) {
            json.value = beautify(data, {
                indent_size: 2,
                space_in_empty_paren: true,
            });
        } else {
            json.value = data
        }
        fileExt.value = fileEvent.value.File.split('.')[1]
        console.log(fileExt.value);
    })
}

let json = ref()
let isDark = ref(true)

const SaveConfig = async () => {
    synchronization()
    let param = {
        FileName: fileEvent.value.File,
        Content: fileExt.value === 'csv' ? Papa.unparse(handsontableData.value) : json.value
    }
    let [data, err] = await Client.Do(AdminConfigFile.SaveConfig, param)
    if (err) {
        return tip.e(err)
    }
    FileHistory.value = []
    FileHistoryIndex.value = null
    isFileHistory.value = false
    tip.s(t('保存成功'))
}

// watch(() => json.value, (newVal) => {
// if (fileEvent.value && fileEvent.value.File.endsWith(".csv")) {
//     handsontableData.value = Papa.parse(newVal).data
// }
// }, { immediate: true, deep: true })
// watch(() => handsontableData.value, (newVal) => {
// if (fileEvent.value && fileEvent.value.File.endsWith(".csv")) {
//     json.value = Papa.unparse(newVal)
// }
// }, { immediate: true, deep: true })
const changeDark = () => {
    return new Promise(async (resolve) => {
        synchronization()
        return resolve(true)
    })
}

const synchronization = () => {
    if (fileExt.value === 'csv' && isDark.value) {
        json.value = Papa.unparse(handsontableData.value)
    }
    if (fileExt.value === 'csv' && !isDark.value) {
        handsontableData.value = Papa.parse(json.value).data
    }
}
</script>
<style scoped lang='scss'>
.jsonView {
    width: 100%;
    overflow: scroll;
    display: flex;
}

.treeView {
    width: 200px;
    max-height: 600px;
    overflow: scroll;
}
.d-none{
    height: 100px;
    width: 100%;
    margin-bottom: 1rem;
}
.el-row{
    width: 100%;
}
.mainview {
    width: 100%;
    padding-left: 10px;

    .codeView {
        height: calc(100vh - 230px);
        overflow: scroll;

        ::-webkit-scrollbar {
            width: 4px !important;
            height: 4px !important;
        }
    }
}



.fileHistory {
    width: 100%;
    display: flex;
    justify-content: space-between;

    p {
        cursor: pointer;
    }

    .act {
        color: #409eff
    }
}
</style>
