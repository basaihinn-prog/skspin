<template>
    <el-form-item :label="$t('币种')">
        <el-select v-model.number="paramData" @clear="clearCurrency" filterable clearable style="width: 150px" :placeholder="$t('请选择')"
                   @change="selectOperator()">
            <el-option v-for='item in operatorData' :key="item.CurrencyCode" :label="item.CurrencyCode != 'All' ? `【${item.CurrencyCode}】${item.CurrencyName}` : item.CurrencyName" :value="item.CurrencyCode">{{  item.CurrencyCode != 'All' ? `【${item.CurrencyCode}】${item.CurrencyName}` : item.CurrencyName }}</el-option>

        </el-select>
    </el-form-item>


</template>

<script setup lang="ts">
import {onMounted, ref, reactive, defineProps, defineEmits} from 'vue';
import {useStore} from '@/pinia/index';
import {useI18n} from 'vue-i18n';
import {Client} from "@/lib/client";
import {AdminInfo} from "@/api/adminpb/info";
import {tip} from "@/lib/tip";

const {t} = useI18n();
const store = useStore();
const {defaultOperatorEvent, haseAll} = defineProps(['defaultOperatorEvent', "haseAll"]);
const emit = defineEmits();
let operatorData = ref([])
let defaultOperatorObj = reactive({
    name: '',
    id: null,
})
let paramData = ref()
let operatorParam = reactive({
    PageIndex: 1,
    PageSize: 10000,
    OperatorType: 2
})
import {useRoute} from 'vue-router'
import {Comm} from "@/api/comm";
import {debug} from "util";

const route = useRoute()
onMounted(async () => {
        await operatorList()
        if (operatorData.value.length && route.query.OperatorId) {
            let a = operatorData.value.find(list => list.CurrencyCode === route.query.CurrencyCode)

            paramData.value = a.Id
        }




})
const selectOperator = () => {
    if (!paramData.value) {
        paramData.value = null
    }
    // 触发查询操作，可以通过 emit 发送事件到父组件
    emit('update:paramData', paramData);
    const selectItem = operatorData.value.find(list => list.CurrencyCode === paramData.value) || {
        CurrencyCode: "All",
        CurrencyName: t("全部"),

    }

    if(selectItem.CurrencyCode == "All"){
        selectItem.CurrencyCode = null
    }
    emit('select-operator', paramData, selectItem);
    emit('select-operatorInfo', selectItem);
};
const operatorList = async () => {
    let [data, err] = await Client.Do(Comm.GetCurrency, {} as any)
    if (err) {
        return tip.e(err)
    }

    if (data.List.length > 0){

        if (haseAll){
            data.List.unshift({
                CurrencyCode:"All",
                CurrencyName:t("全部"),
            })
            paramData.value = "All"
        }

    }

    operatorData.value = data.List || []
}

const clearCurrency = () => {
    paramData.value = "All"
    selectOperator()
}
</script>
<style scoped lang='scss'>


</style>
