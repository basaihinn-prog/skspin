<template>
    <div class="adminList">
        <div class="searchView">
            <el-space wrap>
                <el-form :inline="true" :model="searchParams" class="demo-form-inline">
                    <el-form-item :label="$t('是否存在游戏')">
                        <el-select
                            v-model="searchParams.isExits"
                            :placeholder="$t('请选择是否存在游戏')"
                            style="width: 150px"
                            clearable
                        >
                            <el-option v-for="(item, key) in isExits" :label="item" :value="key"/>
                        </el-select>
                    </el-form-item>


                    <el-form-item>
                        <el-button type="primary" @click="onSearchNotify">{{ $t('查询') }}</el-button>
                    </el-form-item>
                </el-form>
            </el-space>
        </div>


        <div class="page_table_context">
            <div class="flex flex_child_end">
                <el-button type="primary" @click="addDialog = true" style="margin-bottom: 15px">{{ $t('新增') }}</el-button>
            </div>

            <customTable :table-header="tableHeader" :table-data="tableData">
                <template #gameType="scope">

                    <el-input v-model="scope.scope.GameTypeName" style="width: 240px" placeholder="Please input" />

                    <el-button type="primary" @click="EditGameClass">{{ $t('修改') }}</el-button>
                </template>
            </customTable>
        </div>

    </div>
</template>

<script setup lang='ts'>
import {useStore} from "@/pinia";

import customTable from "@/components/customTable/tableComponent.vue"
import {useI18n} from "vue-i18n";

const store = useStore()

let { t } = useI18n()
import {ref} from "vue";
import {GameType} from "@/api/adminpb/game";

let isExits = {
    "2": t("全部"),
    "1": t("是"),
    "0": t("否"),
}
const addForm = ref({})

const searchParams = ref({
    isExits: ""
})

const tableHeader = [
    {
        label: "序号",
        value: "index",
        type: "index",
        width: "220px"
    },
    {
        label: "游戏类型",
        value: "gameType",
        width: "530px"
    },
    {
        label: "游戏数量",
        value: "gameType",

    },
    {
        label: "操作",
        value: "Operator",
        type: "custom",
        width: "180px"
    }
]
const tableData = ref([])

const addDialog = ref(false)

const GetGameTypeList = () => {
    // GameType.GameTypeList()
}

const EditGameClass = () => {

}
const onSearchNotify = () => {

}
const addGameType = () => {

}

</script>

<style scoped lang="scss">
.dialog__form{
    width: 80%;
    margin: 0 auto;
}
</style>
