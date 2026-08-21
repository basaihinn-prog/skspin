<template>
    <customTable v-loading="loading" :tableHeader="languageTableConfig" :tableData="tableData" height="600">

        <template #handleTools>

            <el-space wrap>


                <template v-if="filterStatus">

                    <el-button type="primary" plain @click="commitEdit">{{ $t('提交') }}</el-button>


                    <uploadExcel @uploadFile="uploadSuccess" fileName="lang.json">
                        <el-button type="primary" plain>{{ $t('上传') }}</el-button>
                    </uploadExcel>
                </template>

                <template v-else>
                    <el-input v-model="filter" @input="filterLange"/>
                </template>

                <el-button type="primary" plain @click="filterVisible">{{
                    $t(!filterStatus ? '关闭' : '过滤语言')
                    }}
                </el-button>

            </el-space>


        </template>


        <template #LanguageOperator="scope">
            <el-button type="danger" plain size="small" @click="removeLange(scope.index)">{{ $t('删除') }}</el-button>
        </template>

    </customTable>


    <el-dialog :title="$t('编辑语言')" v-model="dialogFormVisible">
        <el-form :model="EditLang">
            <el-form-item v-for="(key, item) in EditLang" :label="item" label-width="60px">
                <el-input v-model="EditLang[item]"></el-input>
            </el-form-item>

        </el-form>
        <div slot="footer" class="dialog-footer flex flex_child_end">
            <el-button @click="dialogFormVisible = false">取 消</el-button>
            <el-button type="primary" @click="commitDialog">确 定</el-button>
        </div>
    </el-dialog>


</template>
<script setup lang="ts">

import {onMounted, ref} from "vue";
import CustomTable from "@/components/customTable/tableComponent.vue";
import {useStore} from "@/pinia";
import {Client} from "@/lib/client";
import {AdminGroup} from "@/api/adminpb/group";
import {tip} from "@/lib/tip";
import {ElMessageBox} from "element-plus";
import {useI18n} from "vue-i18n";
import {Language, LanguageReq} from "@/api/adminpb/language";
import UploadExcel from "./component/uploadExcel.vue";


const dialogFormVisible = ref(false)

const {t} = useI18n()
const store = useStore()
const loading = ref(true)
const languageTableConfig = ref([])
const baseData = ref([])
const tableData = ref([])
const filterStatus = ref(true)


const EditLang = ref({})

const filter = ref("")

onMounted(() => {
    initTable()
})

const filterLange = () => {

    tableData.value = baseData.value.filter(item => {
        let flag = false
        for (let i = 0; i < languageTableConfig.value.length; i++) {
            if (filter.value) {
                if (item[languageTableConfig.value[i].label] && item[languageTableConfig.value[i].label].indexOf(filter.value) > -1) {
                    flag = true
                    break
                }
            } else {

                flag = true
            }
        }

        return flag
    })

    tableData.value = tableData.value.split(0, 100)

}

const filterVisible = () => {
    filterStatus.value = !filterStatus.value;
    filter.value = ''
}

const initTable = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminGroup.GetLang, {
        PageIndex: 1,
        PageSize: 10000
    })
    loading.value = false
    if (err) return tip.e(err)

    languageTableConfig.value = []
    tableData.value = []
    for (const langekey in data.List[0]) {

        if (langekey == "zh") {
            languageTableConfig.value.unshift({label: langekey, value: langekey, width: "300px", fixed:"left",hiddenVisible: true,})
        } else {
            languageTableConfig.value.push({label: langekey, value: langekey, width: "300px"})
        }
    }
    languageTableConfig.value.push({
        label: "操作",
        value: "LanguageOperator",
        type: "custom",
        fixed: "right",
        hiddenVisible: true,
        width: "200px"
    })


    tableData.value = data.List.slice(0, 200)
    baseData.value = data.List

}


const editLang = (data) => {
    EditLang.value = data

    dialogFormVisible.value = true
}
const removeLange = (index) => {

    ElMessageBox.confirm(
        t('需要点击提交后才能生效'),
        t('确认将要删除该语言'),

        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            let delData = tableData.value[index]
            tableData.value.splice(index, 1)


            let delIndex = baseData.value.findIndex(item => item["zh"] == delData["zh"])
            baseData.value.splice(delIndex, 1)
        })


}

const commitEdit = async () => {

    let EditLanguage: LanguageReq = {
        FileName: "lang.json",
        Context: [...baseData.value],
    }

    const [response, err] = await Client.Do(Language.EditLanguageConfig, EditLanguage)

    if (err) return tip.e(t(err))

    tip.s(t("修改完成"))

    initTable()
}

const uploadSuccess = () => {
    initTable()
}

const commitDialog = () => {

}

</script>
<style scoped lang="scss">

</style>
