<template>
    <customTable :tableHeader="languageTableConfig" :tableData="tableData">
        <template #handleTools>
            <el-button type="primary" plain>{{ $t('提交') }}</el-button>
        </template>
    </customTable>
</template>
<script setup lang="ts">

import {onMounted, ref} from "vue";
import CustomTable from "@/components/customTable/tableComponent.vue";

const languageTableConfig = ref([
    {label:"币种名称", value:"CurrencyName"},
    {label:"币种编号", value:"CurrencyCode"},
    {label:"币种符号", value:"CurrencySymbol"},
    {label:"操作", value:"LanguageOperator"},
])
const tableData = ref([])

onMounted(()=>{
    alert("币种")
})

</script>
<style scoped lang="scss">

</style>
