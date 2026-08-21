

<template>
    <el-form ref="addFormRef" :model="commit" label-width="120px" :inline="true" class="dialog__form">

        <el-row :gutter="18">

            <el-col :span="12">
                <el-form-item :label="$t('代理商户') + ':'">
                    {{ gameInfoData.userName }}
                </el-form-item>
                <el-form-item :label="$t('游戏编号') + ':'">
                    {{ gameInfoData.GameId }}
                </el-form-item>


                <el-form-item :label="$t('购买免费游戏RTP') + ':'" v-if="props.gameInfo.BuyType && ((merchantInfoData.BuyRTPOff == 0 && store.AdminInfo.GroupId <= 1) || merchantInfoData.BuyRTPOff == 1)">
                    <el-select
                        v-model="gameInfoData.BuyRTP"
                        :placeholder="$t('请选择')"
                    >
                        <template v-for="item in BuyRTP">
                            <el-option
                                :label="item == 0 ? $t('无') : `${item}%`"
                                :value="Number(item)"
                            />
                        </template>

                    </el-select>
                </el-form-item>

                <el-form-item :label="$t('赢取最高钱数') + ':'" style="margin-bottom: 0">
                    <el-input v-model.number="gameInfoData.MaxWinPoints"
                              :disabled="!((merchantInfoData.MaxWinPointsOff == 0 && store.AdminInfo.GroupId <= 1) || merchantInfoData.MaxWinPointsOff == 1)"
                              :placeholder="$t('请输入')" @blur="MaxWinPointsInput" maxlength="7">
                    </el-input>

                </el-form-item>

                <el-text type="danger" style="text-align: right;width: 90%;display: block;margin-top: 10px">
                    {{ $t('赢取最高钱数可设置区间{Num}', {Num:"1-1,000,000"})}}

                </el-text>

<!--                <el-form-item :label="$t('背包显示') + ':'" style="margin-bottom: 0" v-if='gameInfoData.GameManufacturer == "JILI"'>-->
<!--                    <div class="switchContainer">-->
<!--                        <el-switch-->
<!--                            v-model="gameInfoData.ShowBag"-->
<!--                            :active-value="1"-->
<!--                            :inactive-value="0"-->
<!--                        />-->
<!--                    </div>-->
<!--                </el-form-item>-->

            </el-col>

            <el-col :span="12">

                <el-form-item v-if="props.gameInfo.GameManufacturer != 'SPRIBE'">

                    <template #label>
                        <el-tooltip
                            class="box-item"
                            effect="dark"

                            placement="top"
                        >

                            <template #content>
                               <div>{{ $t("心跳型：中奖率约5%-15%，中奖和小奖概率低，有概率出现超大奖，波动很大。")}}</div>
                               <div>{{ $t("波动型：中奖率约15%-20%，中奖和小奖概率低，有概率出现超大奖，波动大。")}}</div>
                               <div>{{ $t("仿正版（默认）：中奖率在20%-30%之间，与官方原版体验相同。")}}</div>
                               <div>{{ $t("混合型：中奖率约25%-35%，大奖概率略低，波动较小。")}}</div>
                               <div>{{ $t("稳定型：中奖率约35%-45%，大奖概率低，数值体验相对稳定。")}}</div>
                               <div>{{ $t("高中奖率：中奖率约45%-60%，大奖概率很低，中奖和小奖的中奖率高，数值稳定性高。")}}</div>
                               <div>{{ $t("超高中奖率：中奖率约40%-70%左右，大奖概率极低，中奖和小奖的中奖率较高，数值稳定性较高。")}}</div>
                            </template>
                            <el-icon><QuestionFilled /></el-icon>
                        </el-tooltip>



                        {{ $t('游戏类型') + ':' }}
                    </template>


                    <el-select
                            v-model="gameInfoData.GamePattern"
                            :placeholder="$t('请选择')"

                    >
                        <el-option
                                v-for="(item, index) in GameType"
                                :key="index"
                                :label="item.value == 5 ? $t(item.label) + '(' + $t('默认') + ')': $t(item.label) "
                                :value="item.value"
                        />
                    </el-select>

                </el-form-item>

                <el-form-item v-if="props.gameInfo.GameId == 'spribe_01'">

                    <template #label>
                        <el-tooltip
                            class="box-item"
                            effect="dark"

                            placement="top"
                        >

                            <template #content>
                                <div>{{ $t("假设配置为1%， 即通杀概率为1%。配置越高，杀的越厉害")}}</div>
                            </template>
                            <el-icon><QuestionFilled /></el-icon>
                        </el-tooltip>



                        {{ $t('飞机1.00倍坠毁率') + ':' }}
                    </template>
                    <el-input
                        v-model="gameInfoData.CrashRate"
                        :placeholder="$t('请输入')"
                    >

                    </el-input>
                </el-form-item>

                <el-form-item style="margin-bottom: 20px"  v-if="props.gameInfo.GameId == 'spribe_01'">
                    <template #label>
                        <el-tooltip
                            class="box-item"
                            effect="dark"

                            placement="top"
                        >

                            <template #content>
                                <div>{{ $t("假设配置为5%， 就类似平台抽水5%， 平台只赢这5%（略有波动）")}}</div>
                            </template>
                            <el-icon><QuestionFilled /></el-icon>
                        </el-tooltip>



                        {{ $t('利润率') + ':' }}
                    </template>
                    <el-input v-model.number="gameInfoData.ProfitMargin"
                              :placeholder="$t('请输入')" @blur="ProfitMarginChange" maxlength="7">
                    </el-input>

                </el-form-item>

                <el-form-item :label="$t('RTP') + ':'" v-if="props.gameInfo.GameId != 'spribe_01'">
                    <el-select
                        v-model="gameInfoData.RTP"
                        :placeholder="$t('请选择')"
                        :disabled="!((merchantInfoData.RTPOff == 0 && store.AdminInfo.GroupId <= 1) ||  merchantInfoData.RTPOff == 1)"
                    >
                        <template v-for="item in RTP_VALUE">
                            <el-option
                                :key="Number(item)"
                                :label="item == 93 ? '93%' + $t('(常用)') : item + '%' "
                                :value="Number(item)"
                                v-if="store.AdminInfo.GroupId <= 1 ||
                                 (merchantInfoData.HighRTPOff > 0
                                 && (Number(item) >= OperatorRTP[merchantInfoData.HighRTPOff].split('-')[0]
                                 && Number(item) <= OperatorRTP[merchantInfoData.HighRTPOff].split('-')[1]) || Number(item) == gameInfoData.RTP)"
                            />
                        </template>
                    </el-select>
                </el-form-item>


                <el-form-item :label="$t('止损止盈开关') + ':'" v-if="gameInfoData.GameManufacturer.toUpperCase() == 'PG'">

                    <div class="switchContainer">
                        <el-switch
                            v-model="gameInfoData.StopLoss"
                            :active-value="0"
                            :inactive-value="1"
                        />
                    </div>
                </el-form-item>


                <!--                <el-form-item :label="$t('显示名字和时间') + ':'" v-if="gameInfoData.GameManufacturer.toUpperCase() != 'JILI'">-->
<!--                    <div class="switchContainer">-->
<!--                        <el-switch-->
<!--                            v-model="gameInfoData.ShowNameAndTimeOff"-->
<!--                            :active-value="1"-->
<!--                            :inactive-value="0"-->
<!--                        />-->
<!--                    </div>-->

<!--                </el-form-item>-->


                <el-form-item :label="$t('赢取最高倍数') + ':'" style="margin-bottom: 0">
                    <el-input v-model.number="gameInfoData.MaxMultiple"
                              :disabled="!((merchantInfoData.MaxMultipleOff == 0 && store.AdminInfo.GroupId <= 1) || merchantInfoData.MaxMultipleOff == 1)"
                              :placeholder="$t('请输入')" maxlength="5" @blur="maxnultiple">
                    </el-input>
                </el-form-item>

                <el-text type="danger" style="text-align: right;width: 90%;display: block;margin-top: 10px">
                    {{ $t('赢取最高倍数可设置{Num}', { Num:"30~10000" }) }}
                </el-text>
                <template>
                    <el-form-item>
                            <span style="color: red; font-size: 13px;white-space: nowrap">{{
                                $t('设置玩家当局在游戏内能赢取的最大倍数上限')
                                }}</span>
                    </el-form-item>
                </template>
            </el-col>
        </el-row>
    </el-form>

</template>

<script setup lang="ts">

import {RTP_VALUE} from "@/lib/RTP_config";
import {computed, ref, watchEffect} from "vue";
import {useStore} from "@/pinia";
import {useI18n} from "vue-i18n";
import {OperatorRTP, GameType} from "@/enum";

const store = useStore()
const {t} = useI18n()


const props = defineProps(["merchantInfo", "gameInfo"])


const BuyRTP = import.meta.env.VITE_BUY_RTP.split(",")

let commit = ref({
    ...props.gameInfo,
})


const gameInfoData = computed(()=> props.gameInfo)


const merchantInfoData = computed({
    get(){
        return props.merchantInfo
    },
    set(value){
        commit.value = {
            ...gameInfoData.value
        }
    }
})
watchEffect(()=>{
    commit.value = gameInfoData.value
})


const ProfitMarginChange = (value) => {

    gameInfoData.value.ProfitMargin = gameInfoData.value.ProfitMargin.toString().replace(/^(0+)|[^\d]+/g,'')
    if (gameInfoData.value.ProfitMargin <= 0) {
        gameInfoData.value.ProfitMargin = 1
    }
    gameInfoData.value = {
        ...gameInfoData.value
    }
}
const maxnultiple = (value) => {
    gameInfoData.value.MaxMultiple = gameInfoData.value.MaxMultiple.toString().replace(/^(0+)|[^\d]+/g,'')
    if (gameInfoData.value.MaxMultiple < 30) {
        gameInfoData.value.MaxMultiple = 30
    }
    if (gameInfoData.value.MaxMultiple > 10000) {
        gameInfoData.value.MaxMultiple = 10000
    }
    gameInfoData.value = {
        ...gameInfoData.value
    }
}

const MaxWinPointsInput = (value) => {
    gameInfoData.value.MaxWinPoints = Number(gameInfoData.value.MaxWinPoints.toString().replace(/^(0+)|[^\d]+/g,''))

    if (gameInfoData.value.MaxWinPoints < 1) {
        gameInfoData.value.MaxWinPoints = 1
    }
    if (gameInfoData.value.MaxWinPoints > 1000000) {
        gameInfoData.value.MaxWinPoints = 1000000
    }
    gameInfoData.value = {
        ...gameInfoData.value
    }
}


defineExpose({
    commit
})

</script>


<style scoped lang="scss">

</style>
