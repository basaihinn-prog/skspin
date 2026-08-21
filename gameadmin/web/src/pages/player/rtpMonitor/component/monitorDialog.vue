<template>
    <el-dialog
            v-model="dialogVisible"
            :title="props.type == 'MoniterConfig' ? $t('机台RTP自动控制编辑'): $t('个人RTP自动控制编辑')"
            :width="store.viewModel === 2 ? '85%' : '1050px'"
            @close="closeDialog"
            height="200px"
            align-center
            @open-auto-focus="openDialog"
    >
        <div style="height: 700px;overflow-y: auto">
            <!--   RTP基础配置     -->
            <div class="config">
                <div class="config_item">
                    <div class="flex">
                        <div class="config_label">{{ $t("定义新手玩家得游戏局数") }}</div>
                        <div class="config_value">
                            <el-select v-model="MonitorForm.MoniterNewbieNum" style="width: 80%">
                                <el-option v-for="item in GamePlayedNum" :key="item" :label="`${$t('{Num}局', {Num:item})}`"
                                           :value="item"/>
                            </el-select>
                        </div>
                        <div class="config_description">{{
                            $t("玩家生涯得游戏局数，在此局数范围为新手玩家，反之为非新手玩家，每个游戏新手局数单独计数")
                            }}
                        </div>
                    </div>
                </div>
                <div class="config_item">
                    <div class="flex">
                        <div class="config_label">{{ $t("玩家RTP误差范围") }}</div>
                        <div class="config_value">
                            <el-select v-model="MonitorForm.MoniterRTPErrorValue" style="width: 80%"
                                       @change="ConfigChange">
                                <el-option v-for="item in RTPDeviation" :key="item" :label="`${item}${$t('%')}`"
                                           :value="item"/>
                            </el-select>
                        </div>
                        <div class="config_description">
                            {{ $t("玩家在游戏内的RTP偏离游戏设定RTP时，玩家进入系统判断是否控制") }}
                        </div>
                    </div>
                </div>
                <div class="config_item">
                    <div class="flex">
                        <div class="config_label">{{ $t("游戏内统计数据周期") }}</div>
                        <div class="config_value">
                            <el-select v-model="MonitorForm.MoniterNumCycle" style="width: 80%">
                                <el-option v-for="item in StaticGamePlayedNum" :key="item"
                                           :label="`${$t('{Num}局', {Num:item})}`" :value="item"/>
                            </el-select>
                        </div>
                        <div class="config_description">{{ $t("玩家在间隔对应局数进行一次RTP的监测") }}</div>
                    </div>
                </div>
            </div>

            <div class="tableTitle"><el-tooltip
                            class="box-item"
                            effect="dark"

                            placement="top"
                        >

                            <template #content>
                               <div>{{ $t("玩家RTP低于配置值时，有概率使玩家强制中奖。玩家RTP是指玩家在对应的机台的终生RTP")}}</div>
                            </template>
                            <el-icon><QuestionFilled /></el-icon>
                        </el-tooltip>
                {{ $t("增加RTP范围内容") }}</div>
            <customTable
                    table-name="playerRTPData_list"
                    :hideTableHandel="true"
                    v-loading="loading"
                    :table-header="tableHeader"
                    :table-data="MoniterAddRTPRangeValue"
                    height="250px">


                <template #description="scope">
                    {{
                    $t("每{Round}局游戏进行检测，若玩家RTP低于配置{MinRTP}%~{MaxRTP}%", {
                      Round: MonitorForm.MoniterNumCycle,
                      MinRTP: scope.scope.RangeMinValue,
                      MaxRTP: scope.scope.RangeMaxValue
                    })
                    }}
                </template>

                <template #NewbieValue="scope">
                    <el-input-number :controls="false" :mine="0" style="width: 120px" v-model="scope.scope.NewbieValue"
                                     :precision="0">
                        <template #suffix>%
                        </template>
                    </el-input-number>
                </template>

                <template #NotNewbieValue="scope">
                    <el-input-number :controls="false" :mine="0" style="width: 120px"
                                     v-model="scope.scope.NotNewbieValue" :precision="0">
                        <template #suffix>%
                        </template>
                    </el-input-number>
                </template>

                <template #Operator="scope">
                    <el-button size="small" type="info" :icon="Minus" circle
                               @click="delContact(scope.index, 'MoniterAddRTPRangeValue')"
                               v-if="scope.index != MoniterAddRTPRangeValue.length - 1"
                    />
                    <el-button size="small" type="primary" :icon="Plus" circle
                               v-if="scope.index != MoniterAddRTPRangeValue.length - 1 ||
                               (scope.index == MoniterAddRTPRangeValue.length - 1 && MoniterAddRTPRangeValue.length != 0)"
                               @click="addContact(scope.index,'MoniterAddRTPRangeValue')"/>

                </template>
            </customTable>

            <div class="tableTitle">
                <el-tooltip
                            class="box-item"
                            effect="dark"

                            placement="top"
                        >

                            <template #content>
                               <div>{{ $t("玩家RTP高于配置值时，有概率使玩家强制不中奖。玩家RTP是指玩家在对应的机台的终生RTP")}}</div>
                            </template>
                            <el-icon><QuestionFilled /></el-icon>
                        </el-tooltip>
                {{ $t("减少RTP范围内容") }}
            </div>
            <customTable
                    table-name="playerRTPData_list"
                    :hideTableHandel="true"
                    v-loading="loading"
                    :table-header="tableHeader"
                    :table-data="MoniterReduceRTPRangeValue"
                    height="250px">


                <template #description="scope">
                    {{
                    $t("每{Round}局游戏进行检测，若玩家RTP高于配置{MinRTP}%~{MaxRTP}%",
                        {
                          Round: MonitorForm.MoniterNumCycle,
                          MinRTP: scope.scope.RangeMinValue,
                          MaxRTP: scope.scope.RangeMaxValue
                        })
                    }}
                </template>

                <template #NewbieValue="scope">
                    <el-input-number :controls="false" :mine="0" style="width: 120px" v-model="scope.scope.NewbieValue"
                                     :precision="0">
                        <template #suffix>%
                        </template>
                    </el-input-number>
                </template>

                <template #NotNewbieValue="scope">
                    <el-input-number :controls="false" :mine="0" style="width: 120px"
                                     v-model="scope.scope.NotNewbieValue" :precision="0">
                        <template #suffix>%
                        </template>
                    </el-input-number>
                </template>

                <template #Operator="scope">
                    <el-button size="small" type="info" :icon="Minus" circle
                               v-if="scope.index != MoniterReduceRTPRangeValue.length - 1"
                               @click="delContact(scope.index, 'MoniterReduceRTPRangeValue')"/>
                    <el-button size="small" type="primary" :icon="Plus" circle
                               v-if="scope.index != MoniterReduceRTPRangeValue.length - 1 ||
                                (scope.index == MoniterReduceRTPRangeValue.length - 1 && MoniterReduceRTPRangeValue.length != 0)"
                               @click="addContact(scope.index,'MoniterReduceRTPRangeValue')"/>

                </template>
            </customTable>

        </div>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="closeDialog">{{ $t('关闭') }}</el-button>
                <el-button type="primary" @click="commitData">
                    {{ $t('提交') }}
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import {useStore} from "@/pinia";
import {useI18n} from "vue-i18n";
import {GamePlayedNum, RTPDeviation, StaticGamePlayedNum} from "@/pages/player/rtpMonitor/component/enmu";

import customTable from "@/components/customTable/tableComponent.vue";
import {computed, ref} from "vue";
import {Minus, Plus} from "@element-plus/icons-vue";
import { Monitor} from "@/api/adminpb/monitor";
import {Client} from "@/lib/client";
import {tip} from "@/lib/tip";

let {t} = useI18n()
let store = useStore()
const emits = defineEmits(["update:modelValue"])
const props = defineProps(["modelValue", "dialogTitle", "MonitorData", "type"])

const MonitorForm = ref({
    MoniterNewbieNum: GamePlayedNum[0],
    MoniterRTPErrorValue: RTPDeviation[0],
    MoniterNumCycle: StaticGamePlayedNum[0],
})
const dialogVisible = computed(() => props.modelValue)
const loading = ref(false)

let tableHeader = ref([
    {label: "序号", value: "AppID", type: "index", width: "80px"},
    {label: "玩家RTP区间（根据商户设定的玩家RTP误差范围）", value: "description", type: "custom"},
    {label: "新手", value: "NewbieValue", type: "custom", width: "150px"},
    {label: "非新手", value: "NotNewbieValue", type: "custom", width: "150px"},
    {label: "操作", value: "Operator", type: "custom", width: "180px"},
])

let MoniterAddRTPRangeValue = ref([])

let MoniterReduceRTPRangeValue = ref([])

const openDialog = () => {

    let propsMonitorData = {...props.MonitorData[props.type]}

    MonitorForm.value = propsMonitorData
    MonitorForm.value.MoniterNewbieNum = propsMonitorData.MoniterNewbieNum || GamePlayedNum[0]
    MonitorForm.value.MoniterRTPErrorValue = propsMonitorData.MoniterRTPErrorValue || RTPDeviation[0]
    MonitorForm.value.MoniterNumCycle = propsMonitorData.MoniterNumCycle || StaticGamePlayedNum[0]

    if (propsMonitorData.MoniterAddRTPRangeValue && propsMonitorData.MoniterReduceRTPRangeValue) {

        MoniterAddRTPRangeValue.value = propsMonitorData.MoniterAddRTPRangeValue
        MoniterReduceRTPRangeValue.value = propsMonitorData.MoniterReduceRTPRangeValue
    } else {

        ConfigChange()
    }


}
const commitData = async () => {

    loading.value = true

    let commitData = {
        ...MonitorForm.value,
        AppID: props.MonitorData.AppID,
        MoniterType: props.type == 'MoniterConfig' ? 0 : 1,
        IsMoniter: 1,
        MoniterAddRTPRangeValue: MoniterAddRTPRangeValue.value,
        MoniterReduceRTPRangeValue: MoniterReduceRTPRangeValue.value,
    }


    let [resp, err] = await Client.Do(Monitor.SetGameMonitor, commitData)
    loading.value = false
    if (err){
        tip.e(t(err))
        return
    }
    tip.s(t('修改成功'))
    closeDialog()
}
const closeDialog = () => {
    emits("update:modelValue", false)
}

const ConfigChange = () => {
    let MoniterRTPErrorValue = MonitorForm.value.MoniterRTPErrorValue

    let ErrorValue = 0


    let model = [{
        RangeMinValue: 1,
        RangeMaxValue: MoniterRTPErrorValue + 1,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }, {
        RangeMinValue: MoniterRTPErrorValue + 1,
        RangeMaxValue: (MoniterRTPErrorValue * 2) + 1,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }, {
        RangeMinValue: (MoniterRTPErrorValue * 2) + 1,
        RangeMaxValue: 999999,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }]
    MoniterAddRTPRangeValue.value = JSON.parse(JSON.stringify(model))
    MoniterReduceRTPRangeValue.value = JSON.parse(JSON.stringify(model))


}

const delContact = (index, key) => {

    if (key == "MoniterAddRTPRangeValue") {
        MoniterAddRTPRangeValue.value.splice(index, 1)
        MoniterAddRTPRangeValue.value = generatorRTPRange([...MoniterAddRTPRangeValue.value])
    }
    if (key == "MoniterReduceRTPRangeValue") {
        MoniterReduceRTPRangeValue.value.splice(index, 1)
        MoniterReduceRTPRangeValue.value = generatorRTPRange([...MoniterReduceRTPRangeValue.value])
    }
}
const addContact = (index, key) => {
    let ac = {
        RangeMinValue: 0,
        RangeMaxValue: 0,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }
    if (key == "MoniterAddRTPRangeValue") {
        MoniterAddRTPRangeValue.value.splice(index + 1, 0, ac)
        MoniterAddRTPRangeValue.value = generatorRTPRange([...MoniterAddRTPRangeValue.value])
    }
    if (key == "MoniterReduceRTPRangeValue") {
        MoniterReduceRTPRangeValue.value.splice(index + 1, 0, ac)
        MoniterReduceRTPRangeValue.value = generatorRTPRange([...MoniterReduceRTPRangeValue.value])
    }

}


const generatorRTPRange = (data) => {
    let dataLen = data.length

    let MoniterRTPErrorValue = MonitorForm.value.MoniterRTPErrorValue




    data[0] = {
        RangeMinValue: 1,
        RangeMaxValue: MoniterRTPErrorValue + 1,
        NewbieValue: data[0].NewbieValue,
        NotNewbieValue: data[0].NotNewbieValue,
    }

    for (let i = 1; i < dataLen - 1; i++) {

        data[i] = {
            RangeMinValue: data[i - 1].RangeMaxValue,
            RangeMaxValue: MoniterRTPErrorValue + data[i - 1].RangeMaxValue,
            NewbieValue: data[i].NewbieValue,
            NotNewbieValue: data[i].NotNewbieValue,
        }
    }

    data[dataLen - 1] = {
        RangeMinValue: data[dataLen - 2] ? data[dataLen - 2].RangeMaxValue : 1,
        RangeMaxValue: 99999,
        NewbieValue: data[dataLen - 1].NewbieValue,
        NotNewbieValue: data[dataLen - 1].NotNewbieValue,
    }

    return data

}


</script>

<style scoped lang="scss">

.config {
  width: 98%;
  border: 1px solid #e5e5e5;
  border-radius: 5px;
  margin: 0 auto;
}

.config_label {
  flex: 2;
  text-indent: 10px;
}

.config_value {
  flex: 2;
}

.config_description {
  flex: 6;
  padding-right: 10px;
}

.config_item:nth-child(2) {
  border-width: 1px 0 1px;
  border-color: #e5e5e5;
  border-style: solid;
}

.config_item {
  padding: 10px 0 10px;
  align-items: center;
}

.tableTitle {
  width: 98%;
  height: 60px;
  margin: 15px auto 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-color-primary);
  color: #ffffff;
  border-radius: 5px 5px 0 0;
}
</style>


