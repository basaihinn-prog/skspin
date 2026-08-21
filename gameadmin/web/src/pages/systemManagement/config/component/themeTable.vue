<template>
    <div>
        <customTable :tableHeader="languageTableConfig" :tableData="tableData">
            <template #handleTools>


                <uploadExcel :fileName="fileName">
                    <el-button type="primary" plain>{{ $t('上传') }}</el-button>
                </uploadExcel>

            </template>

            <template #gameName="scope">
                <div style="display: flex;align-items: center;">
                    <el-tag size="small" type="primary">{{ scope.scope.GameManufacturerName }}</el-tag>
                    <div style="margin: 0 10px">
                        <el-avatar :src="scope.scope.GameIcon" size="small" />
                    </div>
                    {{ scope.scope.gameName }}
                </div>
            </template>
            <template #tips="scope">

                <div v-html="scope.scope.tips">

                </div>
            </template>
        </customTable>

    </div>
</template>
<script setup lang="ts">

import {onMounted, ref} from "vue";
import CustomTable from "@/components/customTable/tableComponent.vue";
import {Client} from "@/lib/client";
import UploadExcel from "@/pages/systemManagement/config/component/component/uploadExcel.vue";
import Gamelist_container from "@/components/gamelist_container.vue";
import {Language} from "@/api/adminpb/language";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";

const {t} = useI18n()
const store = useStore()


const {fileName} = defineProps(["fileName"])

const languageTableConfig = ref([
    {label:"主题名称", value:"label"},
    {label:"排序", value:"sort"},
])
const tableData = ref([])
const filterStatus = ref(false)
const dialogVisible = ref(false)
const GameId = ref("")
const Manufacturer = ref("")
const GameMap = ref({})


onMounted(() => {
    store.GameList.forEach(item=>{
        GameMap.value[item.ID] = item
    })

    initTips()
})

const openDialog = () => {

}

const defaultGameEvent = ref({})
const selectGameList = (value) => {

    GameId.value = value.gameData || null
    Manufacturer.value = value.manufacturer || null


}





const initTips = async () => {
    const [response, err] = await Client.Do(Language.SelectConfig, {fileName: fileName} as any)
    if (err) return tip.e(t(err))

    tableData.value = response.Context



}

</script>
<style scoped lang="scss">

</style>
