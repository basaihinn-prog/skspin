<template>

    <div>
        <div style="width: 100%; height: 30px;display: flex;font-weight: bolder">
            {{ $t("游戏下注") }}
        </div>


        <template v-if="props.gameInfo.GameManufacturer != 'SPRIBE'">

            <div class="setLimitsContainer" v-if="gameInfo.GameManufacturer == 'PP'">
                <el-form label-width="110px" :inline="true" style="width: 100%">
                    <el-form-item :label="`${$t('设定下限值')}:`">
                        <el-input v-model="onlineUpNum" @blur="inputChange"/>
                    </el-form-item>
                    <el-form-item :label="`${$t('设定上限值')}:`">
                        <el-input v-model="onlineDownNum" @blur="inputChange"/>
                    </el-form-item>

                </el-form>
            </div>

            <custom-table :table-header="tableHeader" :table-data="tableData" height="270">

                <template #handleTools>

                    <template v-if="gameInfo.GameManufacturer == 'PG' || gameInfo.GameManufacturer == 'JDB'|| gameInfo.GameManufacturer == 'HACKSAW'">

                        <el-button type="primary" plain @click="addBettingInfo" v-if="props.gameInfo.ChangeBetOff == 1">{{ $t("添加") }}</el-button>
                    </template>

                </template>
                <template #BetValue="scope">
                    <el-input v-model="scope.scope.BetValue" @input="betChange(scope.index, $event)" style="width: 140px"
                              :placeholder="$t('请输入')"/>
                </template>
                <template #info="scope">
                    <template v-if="!propLinNum">
                        <div>{{ $t('配置值') }}({{ scope.scope.BetValue }})x{{ $t('投注等级') }}(1)x{{
                                $t('线数')
                            }}({{ linNum }})={{ $t('总投注') }}
                            <span>{{ formatInfo(scope.scope.BetValue, linNum) }} </span></div>
                    </template>
                    <template v-else>
                        <div>
                            {{ $t('总投注=线数（X）*投注等级*配置值') }}
                            <br/>
                            {{ $t('常用线数：10线数，20线数') }}
                        </div>
                    </template>
                </template>
                <template #JDBinfo="scope">

                    <div style="display: flex;justify-content: space-between">


                        <template v-if="!propLinNum">
                            <div>{{ $t('配置值') }}({{ scope.scope.BetValue }})x{{ $t('投注等级') }}(1)x{{
                                    $t('线数')
                                }}({{ linNum }})={{ $t('总投注') }}
                                <span>{{ formatInfo(scope.scope.BetValue, linNum) }} </span></div>
                        </template>
                        <template v-else>
                            <div>
                                {{ $t('总投注=线数（X）*投注等级*配置值') }}
                                <br/>
                                {{ $t('常用线数：10线数，20线数') }}
                            </div>
                        </template>
                        <el-button @click="delBettingInfo(scope.index)" size="small" type="danger" plain
                                   v-if="gameInfo.ChangeBetOff == 1">

                            {{ $t("删除") }}
                        </el-button>
                    </div>
                </template>

                <template #HACKSAWinfo="scope">

                    <div style="display: flex;justify-content: center">


                        <el-button @click="delBettingInfo(scope.index)" size="small" type="danger" plain
                                   v-if="gameInfo.ChangeBetOff == 1">

                            {{ $t("删除") }}
                        </el-button>
                    </div>
                </template>
                <template #BetInfo="scope">
                    <el-space wrap style="justify-content: center">
                        <template v-if="!propLinNum">
                            <el-popover placement="right" :width="400" trigger="click">
                                <template #reference>
                                    <el-button size="small" type="primary" plain style="margin-right: 16px">{{
                                            $t("查看详情")
                                        }}
                                    </el-button>
                                </template>

                                <el-table :data="scope.scope.BetInfo">
                                    <el-table-column width="120" property="BetValue" :label="$t('投注额')"/>
                                    <el-table-column width="80" property="multipleValue" :label="$t('面额')"/>
                                    <el-table-column width="80" property="linNum" :label="$t('游戏线数')"/>
                                    <el-table-column width="80" property="result" :label="$t('BET值')">
                                        <template #default="scope">
                                            {{ scope.row.result.toFixed(4) }}
                                        </template>
                                    </el-table-column>
                                </el-table>
                            </el-popover>
                        </template>
                        <template v-else>

                            <div>
                                {{ $t('总投注=线数（X）*投注等级*配置值') }}
                                <br/>
                                {{ $t('常用线数：10线数，20线数') }}
                            </div>



                        </template>

                        <el-button @click="delBettingInfo(scope.index)" size="small" type="danger" plain
                                   v-if="gameInfo.GameManufacturer == 'PG' && gameInfo.ChangeBetOff == 1">

                            {{ $t("删除") }}
                        </el-button>
                    </el-space>
                </template>

            </custom-table>
        </template>

        <template v-else>



            <div class="setLimitsContainer setLimitsContainer2" >


                <el-form ref="addFormRef" :model="SPRIBE_AircraftData" :rules="addFormRules" label-width="120px" :inline="true" class="dialog__form">

                    <el-row :gutter="18">

                        <el-col :span="12">
                            <el-form-item :label="$t('刻度') + ':'">
                                <el-input-number :controls="false"
                                                 style="text-align: left"
                                                 v-model="SPRIBE_AircraftData.Scale"
                                                 :precision="1"
                                                 :max="Number(SPRIBE_AircraftData.OnlineDownNum)"/>

                            </el-form-item>
                            <el-form-item :label="$t('投注下限') + ':'" prop="OnlineUpNum">

                                <el-input-number :controls="false" :min="1" v-model="SPRIBE_AircraftData.OnlineUpNum" />
                            </el-form-item>


                            <el-form-item :label="$t('投注上限') + ':'"  prop="OnlineUpNum">

                                <el-input-number :controls="false" :min="1" v-model="SPRIBE_AircraftData.OnlineDownNum"/>
                            </el-form-item>


                            <el-form-item :label="$t('默认值') + ':'" prop="DefaultCs">
                                <el-input-number :controls="false"
                                                 v-model="SPRIBE_AircraftData.DefaultCs"
                                                 :min="Number(SPRIBE_AircraftData.OnlineUpNum)"
                                                 :max="Number(SPRIBE_AircraftData.OnlineDownNum)"

                                />


                            </el-form-item>


                        </el-col>

                        <el-col :span="12">
                            <el-form-item :label="$t('快捷投注1') + ':'" prop="defaultCS_0">
                                <el-input-number :controls="false"
                                                 :precision="2"
                                                 v-model="SPRIBE_AircraftData.BetBase[0]"
                                                 :min="Number(SPRIBE_AircraftData.OnlineUpNum)"
                                                 :max="Number(SPRIBE_AircraftData.OnlineDownNum)"
                                                 :placeholder="$t('请输入')">
                                </el-input-number>

                            </el-form-item>
                            <el-form-item :label="$t('快捷投注2') + ':'" prop="defaultCS_1">
                                <el-input-number :controls="false"
                                                 :precision="2"
                                                 :min="Number(SPRIBE_AircraftData.OnlineUpNum) + 1"
                                                 :max="Number(SPRIBE_AircraftData.OnlineDownNum)"
                                                 v-model="SPRIBE_AircraftData.BetBase[1]"
                                                 :placeholder="$t('请输入')">
                                </el-input-number>

                            </el-form-item>
                            <el-form-item :label="$t('快捷投注3') + ':'" prop="defaultCS_2">
                                <el-input-number :controls="false"
                                                 :precision="2"
                                                 :min="Number(SPRIBE_AircraftData.OnlineUpNum) + 2"
                                                 :max="Number(SPRIBE_AircraftData.OnlineDownNum)"
                                                 v-model="SPRIBE_AircraftData.BetBase[2]"
                                                 :placeholder="$t('请输入')">
                                </el-input-number>

                            </el-form-item>
                            <el-form-item :label="$t('快捷投注4') + ':'" prop="defaultCS_3">
                                <el-input-number :controls="false"
                                                 :precision="2"
                                                 :min="Number(SPRIBE_AircraftData.OnlineUpNum) + 3"
                                                 :max="Number(SPRIBE_AircraftData.OnlineDownNum)"
                                                 v-model="SPRIBE_AircraftData.BetBase[3]"
                                                 :placeholder="$t('请输入')">
                                </el-input-number>

                            </el-form-item>


                        </el-col>
                    </el-row>
                </el-form>




            </div>



        </template>
    </div>
</template>
<script setup lang="ts">


import customTable from "@/components/customTable/tableComponent.vue"
import {ElMessageBox, FormInstance, FormRules} from "element-plus";
import {useI18n} from "vue-i18n";
import {computed, onMounted, reactive, ref, watchEffect} from "vue";
import {GeneratorDefaultCs} from "@/pages/systemManagement/component/gameSettingInfoComponent/generatorDefaultCS";
import {useStore} from "@/pinia";
import {tip} from "@/lib/tip";

let MIN_BET = 0.005
const props = defineProps(["merchantInfo", "gameInfo", "tableData", "propLinNum"])

const {t} = useI18n()

let tableData = ref([])
let SPRIBE_AircraftData = ref({
    Scale:0.1,
    OnlineDownNum:100,
    OnlineUpNum:0.1,
    DefaultCs:1,
    BetBase:[1,2,5,10],
})


const defaultCSValid = (rule: any, value: any, callback: any) => {

    const currentNum = SPRIBE_AircraftData.value.BetBase[rule.field.split("_")[1]]

    if(SPRIBE_AircraftData.value.OnlineDownNum < currentNum || currentNum >SPRIBE_AircraftData.value.OnlineUpNum){
        callback("投注下限不能大于投注上线")
        return
    }
    callback()
}
const OnlineValid = (rule: any, value: any, callback: any) => {

    if(SPRIBE_AircraftData.value.OnlineDownNum < SPRIBE_AircraftData.value.OnlineUpNum){
        callback("投注下限不能大于投注上线")
        return
    }
    callback()
}
const addFormRules = reactive<FormRules>({
    OnlineUpNum: [{trigger: 'blur', validator:OnlineValid}],
    OnlineDownNum:  [{trigger: 'blur', validator:OnlineValid}],
    defaultCS_0: [{trigger: 'blur', validator:defaultCSValid}],
    defaultCS_1: [{trigger: 'blur', validator:defaultCSValid}],
    defaultCS_2: [{trigger: 'blur', validator:defaultCSValid}],
    defaultCS_3: [{trigger: 'blur', validator:defaultCSValid}],
})

const store = useStore()

const onlineUpNum = ref(1)
const onlineDownNum = ref(10)
const level = ref(10)
const linNum = ref(1)

const inputChange = () => {
    tableData.value = []

    onlineUpNum.value = parseFloat(onlineUpNum.value.toString().replace(/(\.\d{0,4})\d*$/,'$1'))
    onlineDownNum.value = parseFloat(onlineDownNum.value.toString().replace(/(\.\d{0,4})\d*$/,'$1'))



    if (onlineUpNum.value < MIN_BET){
        onlineUpNum.value = MIN_BET
    }


    if (onlineDownNum.value / onlineUpNum.value < 10) {
        onlineDownNum.value = onlineUpNum.value * 10
    }


    const betList = GeneratorDefaultCs(onlineUpNum.value * 10000, onlineDownNum.value * 10000)

    for (const index in betList) {
        let BetInfoList = []
        for (let i = 0; i < level.value; i++) {
            let currentBet = i + 1
            BetInfoList.push({
                BetValue: betList[index],
                multipleValue: currentBet,
                result: parseFloat(betList[index]) * linNum.value * currentBet,
                linNum: linNum.value
            })
        }
        tableData.value.push({
            BetValue: betList[index],
            BetInfo: [...BetInfoList]
        })
    }

}
const tableHeader = ref([
    {label: "序号", value: "index", type: "index", width: "80px"},
    {label: "投注额", value: "BetValue", type: `${props.gameInfo.ChangeBetOff == 1 ? 'custom' : ''}`, width: props.gameInfo.GameManufacturer == 'HACKSAW' ? "370px" : "170px"},
])

onMounted(() => {


    if (props.propLinNum) {
        linNum.value = props.propLinNum
    } else {
        linNum.value = store.LineMap[props.gameInfo.GameId]
    }


    tableData.value = props.tableData
    if (props.gameInfo.GameManufacturer == "PG") {
        tableHeader.value.push({
            label: "投注详情",
            value: "BetInfo",
            type: "custom",
            width: "auto",
            hiddenVisible: true
        })
    }
    if (props.gameInfo.GameManufacturer == "JILI" || props.gameInfo.GameManufacturer == "TaDa") {
        tableHeader.value.push({label: "投注详情", value: "info", type: 'custom', width: "auto"})
    }
    if (props.gameInfo.GameManufacturer == "JDB") {
        tableHeader.value.push({label: "投注详情", value: "JDBinfo", type: 'custom', width: "auto"})
    }

    if (props.gameInfo.GameManufacturer == "HACKSAW") {
        tableHeader.value.push({label: "操作", value: "HACKSAWinfo", type: 'custom', width: "auto"})
    }


    if (props.gameInfo.GameManufacturer == "PP" ) {

        tableHeader.value.push({
            label: "投注详情",
            value: "BetInfo",
            type: "custom",
            width: "auto",
            hiddenVisible: true
        })

        onlineUpNum.value = props.gameInfo.OnlineUpNum
        onlineDownNum.value = props.gameInfo.OnlineDownNum


        inputChange()

    }

    if (props.gameInfo.GameManufacturer == "SPRIBE") {

        SPRIBE_AircraftData.value.Scale = props.gameInfo.Scale
        SPRIBE_AircraftData.value.OnlineDownNum = props.gameInfo.OnlineDownNum
        SPRIBE_AircraftData.value.OnlineUpNum = props.gameInfo.OnlineUpNum
        SPRIBE_AircraftData.value.DefaultCs = props.gameInfo.DefaultCs
        SPRIBE_AircraftData.value.BetBase = props.gameInfo.BetBase.split(",")
    }
})


const addBettingInfo = () => {

    let BetInfoList = []


    for (let i = 1; i <= level.value; i++) {

        BetInfoList.push({
            BetValue: MIN_BET,
            multipleValue: i,
            result: MIN_BET * linNum.value * i,
            linNum: linNum.value
        })
    }

    tableData.value.push({
        BetValue: MIN_BET,
        BetInfo: [...BetInfoList]
    })
}

const betChange = (index, val) => {

     val = val.replace(/(\.\d{0,4})\d*$/,'$1')

    if (props.gameInfo.GameManufacturer.toUpperCase() == 'HACKSAW'){
        val = val.replace(/(\.\d{0,2})\d*$/,'$1')
    }

    if (props.gameInfo.GameManufacturer.toUpperCase() != 'JILI'){


        if (val < MIN_BET && val != 0){
            val = MIN_BET
        }

    }
    let BetInfoList = []
    for (let i = 1; i <= level.value; i++) {
        BetInfoList.push({
            BetValue: val,
            multipleValue: i,
            result: parseFloat(val) * linNum.value * i,
            linNum: linNum.value
        })
    }

    tableData.value[index].BetInfo = [...BetInfoList]
    tableData.value[index].BetValue = val

}

const delBettingInfo = (index) => {
    ElMessageBox.confirm(
        t('确认删除下注配置'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {

            if (tableData.value.length > 1){


                tableData.value.splice(index, 1)
            }else{
                tip.e(t("不能清空所有bet项"))
            }
        })
}
// JILI 游戏特有方法
// 通过 配置值x投注等级x线数 得到 总押注
// 再将总押注存入数据库
const formatInfo = (BetValue, linNum, BetNum = 1) => {
    return (BetValue * linNum * BetNum).toFixed(2)
}


defineExpose({
    tableData,
    onlineUpNum,
    onlineDownNum,
    SPRIBE_AircraftData
})

</script>

<style scoped lang="scss">
.setLimitsContainer2 {
    background: #ffffff!important;
    border-radius: 6px;
}
.setLimitsContainer {
  width: 98%;
  height: auto;
  margin: 0 auto;
  background: #f7f7f7;
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-color: #c5c5c5c5;
  border-width: 1px;
  border-style: solid;
}


.setLimitsContainer .el-form .el-form-item {
  margin-bottom: 15px;
  margin-right: 5px;
}

.dialog__form .el-form-item {
    margin-bottom: 15px!important;
}
</style>
