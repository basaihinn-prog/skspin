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
                                        v-if="store.AdminInfo.GroupId <= 1"/>

                </el-space>

                <div class="ruleListContainer" v-loading="loadingStatus">
                    <el-row class="rule_item" v-for="(item, ruleIndex) in ruleList" :key="ruleIndex">


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
                                        <el-option v-for="item in contrast" :label="$t(item.label)"
                                                   :value="item.value"/>
                                    </el-select>
                                    <el-input style="width: 120px" v-model.number.trim="conditionItem.value"
                                              maxlength="12"
                                              @input="conditionItemValueChange(conditionItem.filed, $event,ruleIndex ,conditionIndex)"
                                              :placeholder="$t('请输入触发值')"></el-input>
                                </el-space>

                            </el-col>


                            <el-col :span="24">
                                <div style="height: 15px;width: 100%"></div>
                            </el-col>

                        </template>

                        <template v-for="(executeItem , executeIndex) in item.execute" :key="executeIndex">
                            <el-col :span="3">
                                <div style="line-height: 32px">
                                    {{ executeIndex == 0 ? $t('那么') : $t('并且') }}

                                </div>
                            </el-col>
                            <el-col :span="15">


                                <el-space v-if="executeItem.type == 1">

                                    <el-select style="width: 120px" v-model="executeItem.value">
                                        <el-option v-for="conditionItem in execute" :label="$t(conditionItem.label)" :value="conditionItem.value"/>
                                    </el-select>

                                </el-space>

                                <el-space v-if="executeItem.type == 2">

                                    <el-select style="width: 120px" v-model="executeItem.value">
                                        <el-option v-for="conditionItem in rate" :label="$t(conditionItem.label)" :value="conditionItem.value"/>
                                    </el-select>
                                    <el-select style="width: 120px" v-model="executeItem.value">
                                        <el-option label="=" value="1"/>
                                    </el-select>
                                    <el-select style="width: 120px" v-model="executeItem.rateTime">
                                        <el-option v-for="conditionItem in rateTime" :label="$t(conditionItem.label)" :value="conditionItem.value"/>
                                    </el-select>

                                </el-space>

                            </el-col>
                            <el-col :span="24">
                                <div style="height: 15px;width: 100%"></div>
                            </el-col>
                        </template>


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
import {nextTick, reactive, ref, watchEffect} from "vue";
import Operator_container from "@/components/operator_container.vue";
import {useStore} from "@/pinia";
import {Client} from "@/lib/client";
import {Rule} from "@/api/adminpb/riskManage";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {ElMessage, ElMessageBox} from "element-plus";
import {AdminInfo} from "@/api/adminpb/info";

const props = defineProps(["modelValue", "isAll", "rickType", "title"])
const emits = defineEmits(["update:modelValue"])

const defaultOperatorEvent = ref({})

const store = useStore()
const {t} = useI18n()
const loadingStatus = ref(false)


const default_returnRate = [{
    condition: [
        {filed: "balance", contrast: "gt", value: 10000},
    ],
    execute: [
        {
            type:1,
            value: "1"
        },
        {
            rateTime: 24,
            type:2,
            value: "1"
        }
    ]
}]

const default_total_value = -99999999999999

let ruleList = ref([])
const init = ref(false)

const form = ref({
    AppID: ""
})


let condition = reactive([
    {label: "余额", value: "balance"},
])
const contrast = ref([
    {label: "<", value: "gt"},
])
const execute = ref([
    {label: "发出预警", value: "1"},
])

const rate = ref([
    {label: "发出频率", value: "1"},
])


const rateTime = ref([
    {label: "1小时", value: 1},
    {label: "3小时", value: 3},
    {label: "5小时", value: 5},
    {label: "12小时", value: 12},
    {label: "一天", value: 24},
])


const operatorInfo = ref(null)

const operatorListChange = (value) => {

    if (value) {
        form.value.AppID = value.AppID
    } else {
        form.value.AppID = ""
    }

    getOperatorRickRule()

}

const openDialog = async () => {

    ruleList.value = [...default_returnRate]

    await getOperatorRickRule()
}

const getOperatorRickRule = async () => {
    loadingStatus.value = true
    if (form.value.AppID){
        const [response, err] = await Client.Do(AdminInfo.getOperatorByAppId, {
            AppID: form.value.AppID
        } as any)
        loadingStatus.value = false
        operatorInfo.value = response

        let a = [{
            condition: [
                {filed: "balance", contrast: "gt", value: response.BalanceAlert},
            ],
            execute: [
                {
                    type:1,
                    value: "1"
                },
                {
                    rateTime: response.BalanceAlertTimeInterval == 0 ? null : response.BalanceAlertTimeInterval,
                    type:2,
                    value: "1"
                }
            ]
        }]

        ruleList.value = a
    }



}

const conditionItemValueChange = (field, value, ruleIndex, conditionIndex) => {
    let upValue = value.replace(/[^-\d\.]/g, '')
    if (field == 'RTP') {

        upValue = value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g, '$1')
    } else {

        upValue = value.replace(/[^-\d\.]/g, '')
    }

    if (value.indexOf(".") != -1 && value.split(".")[1].length > 3) {
        upValue = `${value.split(".")[0]}.${value.split(".")[1].slice(0, 3)}`
        upValue = parseFloat(upValue)
    }

    ruleList.value[ruleIndex].condition[conditionIndex].value = Number(upValue)

}


const changeFiled = (value, ruleIndex, conditionIndex) => {
    for (const valueKey in ruleList.value[ruleIndex].condition) {
        let item = ruleList.value[ruleIndex].condition[valueKey]

        if (item.filed == value && valueKey != conditionIndex) {

            tip.e(t("该类型已存在"))
            ruleList.value[ruleIndex].condition[conditionIndex].filed = "balance"
            return
        }
    }
}

// 添加规则
const addRule = (ruleIndex) => {
    if (props.rickType == 'transfer') {
        ruleList.value.push(...default_transfer)

    } else {
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

    ruleList.value[ruleIndex].condition.push({filed: "balance", contrast: "gt", value: 0})
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




            // operatorInfo.value.BalanceAlert = ruleList.value[0].condition[0].value
            // operatorInfo.value.BalanceAlertTimeInterval = ruleList.value[0].execute[1].rateTime


            let param = {
                AppID: form.value.AppID,
                IncrementKey:["BalanceAlert","BalanceAlertTimeInterval"],
                IncrementValue:{
                    "BalanceAlert":ruleList.value[0].condition[0].value,
                    "BalanceAlertTimeInterval":ruleList.value[0].execute[1].rateTime,
                },
            }

            const [res, err] = await Client.Do(AdminInfo.IncrementUpdataOperator, param)
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
