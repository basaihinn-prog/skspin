<template>

    <el-dialog :width="store.viewModel === 2 ? '85%' : '950px'" v-model="props.modelValue" :title="$t('列表显示栏目')"
               @open="openDialog" @close="emits('update:modelValue')" destroy-on-close>

        <el-checkbox-group v-model="checkedCities" @change="changeHeaderVisible" :min="1">
            <el-row>
                <template v-for="item in VisibleHeader">
                    <template v-if="!item.hiddenVisible">
                        <el-checkbox :key="item.value" :label="item.label" :value="item.value" style="min-width: 20%;">
                            {{ $t(item.label) }}
                        </el-checkbox>
                    </template>
                </template>

            </el-row>

        </el-checkbox-group>

        <template #footer>
            <div class="dialog-footer">
                <el-button @click="emits('update:modelValue')">{{ $t('关闭') }}</el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">

import {ref} from "vue";
import {useStore} from "@/pinia";
import {useI18n} from "vue-i18n";

const store = useStore()
const props = defineProps(["modelValue", "tableHeader"])
const emits = defineEmits(["update:modelValue", "controlTableHeader"])

const checkedCities = ref([])
const VisibleHeader = ref()

const {t} = useI18n()
const openDialog = () => {

    VisibleHeader.value = props.tableHeader
    checkedCities.value = props.tableHeader.filter(item => item.visible).map(item => item.value)

}

const changeHeaderVisible = (value) => {
    emits("controlTableHeader", value)
}


</script>

<style scoped lang="scss">

</style>
