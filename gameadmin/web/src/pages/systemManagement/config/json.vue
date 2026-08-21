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
                        :placeholder="$t('关键字过滤搜索')"
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
                        <span>{{ $t(data.Name) }}</span>
                    </span>
                        </template>
                    </el-tree>
                </div>
            </el-col>
            <el-col :lg="20">
                <div class="mainview" v-if="fileEvent">
                    <el-descriptions :column="2" :border="true" size="small" style="width: 98%;margin-left: 1%;margin-bottom: 12px;">
                        <el-descriptions-item :label="$t('名称')">{{ fileEvent.Name }}</el-descriptions-item>
                        <el-descriptions-item :label="$t('文件名')">{{ fileEvent.fileName }}</el-descriptions-item>


                    </el-descriptions>

                    <component :is="fileEvent.tableComponent" :fileName="fileEvent.fileName"/>
                </div>
            </el-col>
        </el-row>


    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref } from 'vue';
import languageTable from "./component/languageTable.vue"
import tipsTable from "./component/tips.vue"
import demoUserTable from "./component/demoGame.vue"
import themeTable from "./component/themeTable.vue"
import { useI18n } from 'vue-i18n';
const { t } = useI18n()
const dNone = ref(false)
onMounted(() => {
    getList()
    window.screen.width<500?dNone.value = true:dNone.value = false
});
let loading = ref(false)



let list = ref([])
const getList = async () => {
    list.value = [
        {Name:"多语言", tableConfig:"language",fileName:"lang.json", tableComponent: languageTable},
        {Name:"tips提示", tableConfig:"tips",fileName:"tips.json", tableComponent: tipsTable},
        {Name:"试玩站", tableConfig:"demoGame",fileName:"demoGame.json", tableComponent: demoUserTable},
        {Name:"试玩站主题配置", tableConfig:"theme",fileName:"theme.json", tableComponent: themeTable},
    ]
}


let fileEvent = ref(null)

const treeNodeClick = async (event) => {

    fileEvent.value = event

}
const filterText = ref('')
import { ElTree } from 'element-plus'
const tree = ref<InstanceType<typeof ElTree>>()

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


    .codeView:hover {
        ::-webkit-scrollbar {
            width: 4px !important;
            height: 4px !important;
            background: var(--el-color-primary);
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
