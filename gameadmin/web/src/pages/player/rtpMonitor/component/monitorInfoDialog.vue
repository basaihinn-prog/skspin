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
            <div class="tableTitle">
                <el-tooltip
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
            <div class="config">
               <div class="config_item" v-for="(item, index) in MoniterAddRTPRangeValue">

                   <div class="config_index">{{index + 1}}</div>
                   <div class="config_description">
                       {{
                           $t("每经过{Round}局游戏进行监测，玩家RTP低于机台{MinRTP}%~{MaxRTP}%间触发判断，新手玩家将有{NewbieValue}%概率被控制，非新手玩家将有{NotNewbieValue}.00%概率控制",
                               {
                                   Round: MonitorForm.MoniterNumCycle,
                                   MinRTP: item.RangeMinValue,
                                   MaxRTP: item.RangeMaxValue,
                                   NewbieValue: item.NewbieValue,
                                   NotNewbieValue: item.NotNewbieValue,
                                   Step: MonitorForm.MoniterRTPErrorValue,
                               }
                           )
                       }}
                   </div>
               </div>
            </div>
            <div class="tableTitle">
                <el-tooltip
                    class="box-item"
                    effect="dark"

                    placement="top"
                >

                    <template #content>
                        <div>{{ $t("玩家RTP高于配置值时，有概率使玩家强制不中奖。玩家RTP是指玩家在对应的机台的终生RTP") }}
                        </div>
                    </template>
                    <el-icon>
                        <QuestionFilled/>
                    </el-icon>
                </el-tooltip>
                {{ $t("减少RTP范围内容") }}
            </div>
            <div class="config">
                <div class="config_item" v-for="(item, index) in MoniterReduceRTPRangeValue">

                    <div class="config_index">{{index + 1}}</div>
                    <div class="config_description">
                        {{
                            $t("每经过{Round}局游戏进行监测，玩家RTP高于机台{MinRTP}%~{MaxRTP}%间触发判断，新手玩家将有{NewbieValue}%概率被控制，非新手玩家将有{NotNewbieValue}.00%概率控制",
                                {
                                    Round: MonitorForm.MoniterNumCycle,
                                    MinRTP: item.RangeMinValue,
                                    MaxRTP: item.RangeMaxValue,
                                    NewbieValue: item.NewbieValue,
                                    NotNewbieValue: item.NotNewbieValue,
                                    Step: MonitorForm.MoniterRTPErrorValue,
                                }
                            )
                        }}
                    </div>
                </div>
            </div>
        </div>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="closeDialog">{{ $t('关闭') }}</el-button>

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
const closeDialog = () => {
    emits("update:modelValue", false)
}

const ConfigChange = () => {
    let MoniterRTPErrorValue = MonitorForm.value.MoniterRTPErrorValue

    let ErrorValue = 0
    if (MoniterRTPErrorValue == 1){
        ErrorValue = 2
    }else{
        ErrorValue = MoniterRTPErrorValue
    }

    let model = [{
        RangeMinValue: 1,
        RangeMaxValue: ErrorValue,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }, {
        RangeMinValue: ErrorValue,
        RangeMaxValue: ErrorValue * 2,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }, {
        RangeMinValue: ErrorValue * 2,
        RangeMaxValue: 999999,
        NewbieValue: 0,
        NotNewbieValue: 0,
    }]
    MoniterAddRTPRangeValue.value = JSON.parse(JSON.stringify(model))
    MoniterReduceRTPRangeValue.value = JSON.parse(JSON.stringify(model))


}


</script>

<style scoped lang="scss">
.config{
    width: 98%;
    margin: 0 auto ;
    border: 1px solid #e5e5e5;
}
.config_item:last-child{

    border-bottom: none;
}
.config_item{
    display: flex;
    align-items: center;
    padding: 10px 0;
    border-bottom: 1px solid #e5e5e5;
}
.config_index{
    width: 50px;
    text-indent: 10px;
}
.config_description{
   flex: 1;
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


