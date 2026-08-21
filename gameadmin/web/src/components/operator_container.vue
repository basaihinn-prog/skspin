<template>
    <el-form-item :label="$t(operatorType == 1 ? '线路商' : '商户')">
        <el-select v-model.number="paramData" style="width: 150px" filterable clearable :placeholder="$t('请选择')"
                   :disabled="store.AdminInfo.GroupId == 3"
                   popper-class="operator_proper"
                   @change="selectOperator()"
                   @clear="operatorClear">

            <el-option v-for='item in operatorData' :key="item.Id" :label="$t(item.AppID)" :value="item.Id"/>
        </el-select>


    </el-form-item>



</template>

<script setup lang="ts">
import {onMounted, ref, reactive, defineProps, defineEmits, watch} from 'vue';
import {useStore} from '@/pinia/index';
import {useI18n} from 'vue-i18n';
import {Client} from "@/lib/client";
import {AdminInfo} from "@/api/adminpb/info";
import {tip} from "@/lib/tip";

const {t} = useI18n();
const store = useStore();
const {defaultOperatorEvent, isInit, haseAll, operatorType} = defineProps(['defaultOperatorEvent', "isInit", "haseAll", "operatorType"]);
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
    OperatorType: 2,
    Status: -1
})
import {useRoute} from 'vue-router'
import {storeToRefs} from "pinia";

const storeRef = storeToRefs(store)
const watchStatus = ref(false)

const route = useRoute()
onMounted(async () => {
    await operatorList()
    if (operatorData.value.length && route.query.OperatorId) {

        let a = operatorData.value.find(list => list.Id === Number(route.query.OperatorId))
        paramData.value = a.Id
    }
})
watch(paramData, (newData) => {
    if (watchStatus.value){

        selectOperator()
        watchStatus.value = false
    }
})
const selectOperator = () => {
    if (!paramData.value) {
        paramData.value = null


        operatorClear()
        return
    }
    // 触发查询操作，可以通过 emit 发送事件到父组件
    emit('update:paramData', paramData);
    const selectItem = operatorData.value.find(list => list.Id === Number(paramData.value))
    emit('select-operator', paramData);
    emit('select-operatorInfo', selectItem);
};
const operatorList = async () => {

    if (store.AdminInfo.GroupId != 3) {

        operatorParam.OperatorType = operatorType || 2
        let [data, err] = await Client.Do(AdminInfo.GetOperatorList, operatorParam)
        if (err) {
            return tip.e(err)
        }

        if (data.AllCount !== 0){

            if (haseAll){
                data.List.unshift({
                    Id:"ALL",
                    AppID:"全部",
                })

                paramData.value = "ALL"
            }

        }else{
            data.List = []
        }

        operatorData.value = data.List

        if (isInit) {
            paramData.value = operatorData.value[0].Id
        }

    } else {
        operatorData.value = [{
            Id: store.AdminInfo.BusinessesId,
            AppID: store.AdminInfo.AppID,
        }]
        paramData.value = store.AdminInfo.BusinessesId
    }

    watchStatus.value = true


}

const operatorClear = () => {
    if (haseAll){
        paramData.value = "ALL"
    }else{
        paramData.value = operatorData.value[0].Id
    }

    const selectItem = operatorData.value.find(list => list.Id === Number(paramData.value))
    emit('select-operator', paramData);
    emit('select-operatorInfo', selectItem);
}

defineExpose({
    paramData
})
</script>
<style scoped lang='scss'>


</style>
