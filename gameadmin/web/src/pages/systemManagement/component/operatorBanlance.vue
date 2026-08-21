<template>

    <el-dialog v-model="addDialog" :title="$t('预付款余额')"
               @open-auto-focus="openDialog"
               destroy-on-close
               align-center
               :width="store.viewModel === 2 ? '100%' : '650px'" @close="emits('update:modelValue')">
        <el-form ref="operatorBalanceFormRef" :model="operatorBalanceForm" label-width="140px" :rules=BalanceRules :inline="true"
                 class="dialog__form">



            <el-form-item :label="$t('余额操作类型') + ':'" prop="OperatorType">
                <el-select
                    v-model="operatorBalanceForm.OperatorType"
                    :placeholder="$t('请选择余额操作类型')"
                    @change="merchantBalanceChange"
                >
                    <template v-for="item in BALANCE_OPERATOR_TYPE">
                        <el-option :label="$t(item.label)" :value="item.value"/>
                    </template>

                </el-select>
            </el-form-item>

            <el-form-item :label="$t('金额') + ':'" prop="Balance">
                <el-input
                    v-model="operatorBalanceForm.Balance"
                    @input="merchantBalanceChange"

                    :placeholder="$t('请输入金额')"

                />

<!--                <div style="text-align: right">{{ $t('目前金额') }}：{{ props.merchantBalance }}</div>-->
<!--                <div style="text-align: right">{{ $t('提交后的金额')}}：{{ handledBalance }}</div>-->
            </el-form-item>



            <el-form-item :label="$t('备注') + ':'">
                <el-input
                    v-model="operatorBalanceForm.Remark"
                    :rows="3"
                    type="textarea"
                    :placeholder="$t('请输入备注')"
                />
            </el-form-item>
        </el-form>
        <template #footer>
                <span class="dialog-footer">
                    <el-button type="default" @click="emits('update:modelValue')">{{ $t('取消') }}</el-button>
                    <el-button type="primary" @click="setOperatorBalance">{{ $t('确定') }}</el-button>
                </span>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">

import {computed, reactive, Ref, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {BALANCE_OPERATOR_TYPE} from "@/lib/config";
import {FormRules} from "element-plus";

const props = defineProps(["modelValue", "merchantBalance", "handelType"])
const emits = defineEmits(["update:modelValue", "commitMerchantBalance"])
const addDialog = computed(()=>props.modelValue)
const {t} = useI18n()
const store = useStore()

interface operatorBalance {
    OperatorType: number
    Balance: number
    Remark: string
}

const handledBalance = ref(0)

const operatorBalanceFormRef = ref(null)
const balanceValid = (rule: any, value: any, callback: any) => {

    if (value == 0 || value == ''){
        callback(t("充值或扣款的金额不能为0"))
    }else{
        callback()
    }
}
const operatorValid = (rule: any, value: any, callback: any) => {




        if(value == 2 && props.handelType == 'add'){
            callback(t("新增商户不可以选择扣款"))
        }else{
            callback()
        }

}


const operatorBalanceForm:Ref<operatorBalance> = ref<operatorBalance>({
    OperatorType: 1,
    Balance:0,
    Remark:""
})
const BalanceRules = reactive<FormRules<operatorBalance>>({
    OperatorType: [
        {required: true, message: t('请选择余额操作类型'), trigger: 'blur'},
        {validator: operatorValid, trigger: 'change'}],
    Balance:[{validator: balanceValid, trigger: 'blur'}],

})



const openDialog = (value) => {
    handledBalance.value = 0
    operatorBalanceForm.value = {
        OperatorType: 1,
        Balance:0,
        Remark:""
    }
}

const setOperatorBalance = (value) => {

    operatorBalanceFormRef.value.validate(async valid => {

        if (valid) {
            let emitData = {
                data:operatorBalanceForm.value,
                handledBalance: handledBalance.value
            }
            emits("update:modelValue")
            emits("commitMerchantBalance", emitData)
        }
    })
}

const merchantBalanceChange = () => {
    operatorBalanceForm.value.Balance=operatorBalanceForm.value.Balance.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1')

    if (operatorBalanceForm.value.OperatorType == 1){

        handledBalance.value = parseFloat(props.merchantBalance) + parseFloat(operatorBalanceForm.value.Balance)
    }else{

        handledBalance.value = parseFloat(props.merchantBalance) - parseFloat(operatorBalanceForm.value.Balance)
    }
    handledBalance.value = parseFloat(handledBalance.value).toFixed(2)
}


</script>

<style scoped lang="scss">

</style>
