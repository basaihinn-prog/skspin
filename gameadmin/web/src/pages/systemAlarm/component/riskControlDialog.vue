<template>

    <el-dialog :width="store.viewModel === 2 ? '85%' : '650px'" v-model="props.modelValue" destroy-on-close
               @close="emits('update:modelValue', false)" @open="openDialog"
               :title="$t(`${props.title}（${store.AdminInfo.GroupId > 1 ? '商户' : '总控' }）`)">

        <div>
            <el-form :model="form">

                <el-space>
                <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                    @select-operatorInfo="operatorListChange"
                                    :hase-all="false" :is-init="true"
                                    v-if="store.AdminInfo.GroupId <= 1"></operator_container>

                </el-space>
                <el-form-item>
                    <el-button type="primary" plain  @click="addRule" v-if="(props.rickType == 'transfer' && ruleList.length == 0) || props.rickType != 'transfer'">
                        {{ $t('添加规则') }}
                    </el-button>
                </el-form-item>
                <div class="ruleListContainer">
                    <el-row class="rule_item" v-for="(item, ruleIndex) in ruleList" :key="ruleIndex">

                        <el-col :span="24" style="margin-bottom: 10px">
                            <el-button type="danger" plain
                                       @click="removeRule(ruleIndex)">{{ $t('删除') }}
                            </el-button>
                        </el-col>
                        <template v-for="(conditionItem , conditionIndex) in item.condition" :key="conditionIndex">
                            <el-col :span="3">
                                <div style="line-height: 32px">
                                    {{ conditionIndex == 0 ? $t('如果') : $t('并且') }}
                                </div>
                            </el-col>
                            <el-col :span="18">

                                <el-space>
                                    <el-select style="width: 120px" v-model="conditionItem.filed"
                                               @change="changeFiled($event,ruleIndex ,conditionIndex)">
                                        <template v-for="item in condition">

                                            <el-option :label="$t(item.label)" :value="item.value"/>
                                        </template>
                                    </el-select>
                                    <el-select style="width: 120px" v-model="conditionItem.contrast">
                                        <el-option v-for="item in contrast" :label="$t(item.label)" :value="item.value"/>
                                    </el-select>
                                    <el-input style="width: 120px" v-model.number.trim="conditionItem.value"
                                              maxlength="12"
                                              @input="conditionItemValueChange(conditionItem.filed, $event,ruleIndex ,conditionIndex)"
                                              :placeholder="$t('请输入触发值')"></el-input>
                                </el-space>

                            </el-col>
                            <el-col :span="3">

                                <el-button type="danger" circle plain size="small"
                                           @click="removeCondition(ruleIndex, conditionIndex)"
                                           v-if="conditionIndex != 0">
                                    <el-icon><Minus /></el-icon>
<!--                                           v-if="conditionIndex != 0">{{ $t('删除') }}-->
                                </el-button>
                                <template v-else>


                                    <template v-if="props.rickType != 'transfer'">
                                        <el-button type="primary" circle plain size="small"
                                                   @click="addCondition(ruleIndex)"
                                                   v-if="conditionItem.filed != 'TransferOutMoney'">   <el-icon><Plus /></el-icon>
<!--                                            {{$t('添加如果') }}-->
                                        </el-button>
                                        <el-tooltip
                                                class="box-item"
                                                effect="dark"
                                                :content="$t('转账值不能添加多个条件')"
                                                placement="top-start"
                                                v-else
                                        >
                                            <el-button type="primary" circle plain size="small"
                                                       @click="addCondition(ruleIndex)" disabled>
                                         <el-icon><Plus /></el-icon>
<!--                                                {{ $t('添加如果') }}-->
                                            </el-button>
                                        </el-tooltip>
                                    </template>
                                </template>
                            </el-col>

                            <el-col :span="24">
                                <div style="height: 15px;width: 100%"></div>
                            </el-col>

                        </template>
                        <el-col :span="3">
                            <div style="line-height: 32px">
                                {{ $t("那么") }}
                            </div>
                        </el-col>
                        <el-col :span="15">

                            <el-space>
                                <el-select style="width: 120px" v-model="item.execute.value">
                                    <el-option v-for="item in execute" :label="$t(item.label)" :value="item.value"/>
                                </el-select>
                            </el-space>

                        </el-col>


                    </el-row>
                </div>
            </el-form>

        </div>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="emits('update:modelValue', false)">{{ $t('关闭') }}</el-button>
                <el-button type="primary" @click="submitRule">
                    {{ $t('设置') }}
                </el-button>
            </div>
        </template>
    </el-dialog>


</template>

<script setup lang="ts">
import {reactive, ref, watchEffect} from "vue";
import Operator_container from "@/components/operator_container.vue";
import {useStore} from "@/pinia";
import {Client} from "@/lib/client";
import {Rule} from "@/api/adminpb/riskManage";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {ElMessage, ElMessageBox} from "element-plus";

const props = defineProps(["modelValue", "isAll", "rickType", "title"])
const emits = defineEmits(["update:modelValue"])

const defaultOperatorEvent = ref({})

const store = useStore()
const {t} = useI18n()

const default_transfer = [{
    condition: [
        {filed: "TransferOutMoney", contrast: "gt", value: 1000},

    ],
    execute: {
        value: "1"
    }
}]
const default_returnRate = [{
    condition: [
        {filed: "Single", contrast: "gt", value: 10000},
        {filed: "Total", contrast: "gt", value: 100000},
        {filed: "RTP", contrast: "gt", value: 106},

    ],
    execute: {
        value: "1"
    }
}]

const default_total_value = -99999999999999

let ruleList = ref([])
const init = ref(false)

const form = ref({
    AppID: ""
})


let condition = reactive([])
const contrast = ref([
    {label: ">", value: "gt"},
])
const execute = ref([
    {label: "发出预警", value: "1"},
])


const operatorListChange = (value) => {

    if (value) {
        form.value.AppID = value.AppID
    } else {
        form.value.AppID = ""
    }


        getOperatorRickRule()

}

const openDialog = async () => {
    if (!props.isAll) {
        if (props.rickType == 'returnRate') {
            ruleList.value = JSON.parse(JSON.stringify(default_returnRate))
            condition = [
                {label: "单笔盈利", value: "Single"},
                {label: "累计盈利", value: "Total"},
                {label: "回报率", value: "RTP"}
            ]
        }
        if (props.rickType == 'transfer') {

            ruleList.value = JSON.parse(JSON.stringify(default_transfer))
            condition = [{label: "转账值", value: "TransferOutMoney"}]
        }
    }
   await getOperatorRickRule()
}

const getOperatorRickRule = async () => {
    const [response, err] = await Client.Do(Rule.GetRiskRule, {
        AppID: form.value.AppID
    } as any)

    ruleList.value = []
    if (props.rickType == "returnRate"){
        if(response.OriginReturnRate){


            ruleList.value = JSON.parse(response.OriginReturnRate)


            if (ruleList.value.length) {
                // 懒逼的程序员没爱写递归
                for (const i in ruleList.value) {
                    for (const child_i in ruleList.value[i]) {


                        if (child_i == "condition") {

                            ruleList.value[i].condition = ruleList.value[i].condition.filter(item => item.filed != "Total" || (item.filed == "Total" && item.value != default_total_value))
                        }



                    }
                }
            }
        }else{
            ruleList.value = JSON.parse(JSON.stringify(default_returnRate))
        }
    }

    if (props.rickType == "transfer" ){
        if(response.OriginTransfer){

            ruleList.value = JSON.parse(response.OriginTransfer)

        }else{

            ruleList.value = JSON.parse(JSON.stringify(default_transfer))
        }
    }


}

const conditionItemValueChange = (field, value, ruleIndex, conditionIndex) => {
    let upValue = value.replace(/[^-\d\.]/g,'')
    if (field == 'RTP'){

        upValue = value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1')
    }else{

        upValue = value.replace(/[^-\d\.]/g,'')
    }

    if (value.indexOf(".") != -1 && value.split(".")[1].length > 3){
        upValue = `${value.split(".")[0]}.${value.split(".")[1].slice(0,3)}`
        upValue = parseFloat(upValue)
    }

    ruleList.value[ruleIndex].condition[conditionIndex].value = Number(upValue)

}


const changeFiled = (value, ruleIndex, conditionIndex) => {
    for (const valueKey in ruleList.value[ruleIndex].condition) {
        let item = ruleList.value[ruleIndex].condition[valueKey]

        if(item.filed == value && valueKey != conditionIndex){

            tip.e(t("该类型已存在"))
            ruleList.value[ruleIndex].condition[conditionIndex].filed = "Single"
            return
        }
    }
}

// 添加规则
const addRule = (ruleIndex) => {
    if (props.rickType == 'transfer'){
        ruleList.value.push(...default_transfer)

    }else{
        ruleList.value.push(...default_returnRate)
    }

}

// 添加规则
const removeRule = (ruleIndex) => {
    ruleList.value.splice(ruleIndex, 1);
}

// 添加条件
const addCondition = (ruleIndex) => {
    if (ruleList.value[ruleIndex].condition.length == 3) {
        return
    }

    ruleList.value[ruleIndex].condition.push({filed: "Single", contrast: "gt", value: 0})
}

// 添加条件
const removeCondition = (ruleIndex, conditionIndex) => {
    ruleList.value[ruleIndex].condition.splice(conditionIndex, 1);
}


const submitRule = async () => {


    ElMessageBox.confirm(
        t('确认设置预警值么'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {

            let commitRuleList = JSON.parse(JSON.stringify(ruleList.value));

            if (commitRuleList.length && props.rickType != "transfer") {
                // 懒逼的程序员没爱写递归
                for (const i in commitRuleList) {
                    for (const child_i in commitRuleList[i]) {
                        let exitField = []
                        let flag = false        // 判断是否存在Total
                        for (const last_i in commitRuleList[i].condition) {
                            let item = commitRuleList[i].condition[last_i]

                            if (!flag && item.filed == 'Total'){

                                flag = true
                            }

                            if (exitField.includes(item.filed)) {
                                return tip.e(t("预警类型存在重复的请检查后提交"))
                            }


                            exitField.push(item.filed)
                        }


                        if (!flag){
                            commitRuleList[i].condition.push({filed: "Total", contrast: "gt", value: default_total_value})
                        }

                    }
                }
            }



            let param = {
                AppID: form.value.AppID || store.AdminInfo.AppID,
                RickRule: JSON.stringify(commitRuleList),
                Type:props.rickType.toLowerCase()
            }

            const [res, err] = await Client.Do(Rule.AddRiskRule, param)
            if (err) {
                tip.e(t(err))
                return
            }

            emits("update:modelValue", false)
        })

}
</script>

<style scoped lang="scss">
.ruleListContainer {
  width: 100%;
  height: 300px;
  overflow: auto;

}

.rule_item:last-child {
  border-bottom: 1px solid #dedede;
}

.rule_item {
  padding: 10px 0;
  border-top: 1px solid #dedede;
}
</style>
