<template>

  <!-- 添加弹框 -->
    <el-dialog @open="opendialog" v-model="props.modelValue" :title="$t('游戏设置')"
               destroy-on-close
               :width="store.viewModel === 2 ? '100%' : '650px'" @close="emits('update:modelValue') + ':'">
        <el-form ref="addFormRef" :model="gameInfo" label-width="120px" :inline="true" class="dialog__form">

            <gamelist_container
                    :hase-all="true"
                    :hase-manufacturer="true"
                    :visible-manufacturer="props.gameSettingInfoData.DefaultManufacturerOn"
                    :defaultGameEvent="defaultGameEvent"
                    @select-operator="selectGameList"/>




        </el-form>

            <el-form-item label-width="120px">


                <template #label>
                    <el-tooltip
                            class="box-item"
                            effect="dark"

                            placement="top"
                    >

                        <template #content>
                            <div>{{ $t("根据当前投注金额（Bet），将其乘以预设的倍数，生成新的投注金额（Bet）。由于机台较多，批量处理时间大概需要半分钟左右。")}}</div>
                        </template>
                        <el-icon><QuestionFilled /></el-icon>
                    </el-tooltip>



                    {{ $t('bet调整倍数') + ':' }}
                </template>



                <el-select v-model.number="gameInfo.BetMultiples" style="width:300px">
                    <template v-for="a in BetMultiples">
                        <el-option :label="a" :value="a"/>
                    </template>
                </el-select>
            </el-form-item>


        <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('确认') }}</el-button>
                </span>
        </template>
    </el-dialog>


</template>


<script setup lang="ts">

import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {nextTick, Ref, ref} from "vue";
import {RTP_VALUE} from "@/lib/RTP_config";
import {Client} from "@/lib/client";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {BatchEditGameRTPParam, Merchant} from "@/api/adminpb/merchant";
import {tip} from "@/lib/tip";
import {ElLoading, ElMessageBox} from "element-plus";
import Gamelist_container from "./gamelist_container.vue";
import GameSet from "@/pages/systemManagement/component/gameSettingInfoComponent/gameSet.vue";
import GameBet from "@/pages/systemManagement/component/gameSettingInfoComponent/gameBet.vue";
import {BetMultiples} from "@/enum";


let {t} = useI18n()

const store = useStore()
const props = defineProps(["modelValue", "gameSettingInfoData"])
const emits = defineEmits(['update:modelValue', 'update:gameSettingInfoData', "commit"])
const defaultGameEvent = ref({})

const gameInfo: Ref<BatchEditGameRTPParam> = ref(<BatchEditGameRTPParam>{
    GameList: [],
    AppID: "",
    RTP: Number(RTP_VALUE[0]),
    GamePattern: 1,
    MaxWinPoints: 1000000,
    MaxMultiple: 10000,
    BetMultiples:1
})
const gameBetVisible = ref(false)
const componentInit = ref(false)
const gameBetRef = ref(null)
const merchantInfo = ref({
    GameManufacturer:"",
    GameIds:[],
    OnlineUpNum:0.05,
    OnlineDownNum:6000,
    ChangeBetOff:1,
    tableData:[],
    lineNum:1
})

const selectGameList = (value) => {

    merchantInfo.value.GameManufacturer = value.manufacturer
    merchantInfo.value.GameIds = value.gameData

}

const opendialog = () => {
    gameInfo.value.BetMultiples = 1
}

const AddAdminer = async () => {
    ElMessageBox.confirm(
        t('确认对选中的游戏进行设置Bet'),
        t('是否确认'),
        {
            confirmButtonText: t('确定'),
            cancelButtonText: t('关闭'),
            type: 'warning',
        }
    )
        .then(async () => {
            const loading = ElLoading.service({
                lock: true,
                text: 'Loading',
                background: 'rgba(0, 0, 0, 0.7)',
            })

            const params = {
                AppID: props.gameSettingInfoData.AppID,
                Manufacturer: merchantInfo.value.GameManufacturer == 'ALL' ? '' : merchantInfo.value.GameManufacturer,
                GameIds: merchantInfo.value.GameIds.join(","),
                BetMult:gameInfo.value.BetMultiples,
            }




            const [response, err] = await Client.Do(Merchant.BatchSetGameBetMultiples, params)
            loading.close()
            if (err) {

                tip.e(t("修改失败"))
                return
            }
            tip.s(t("修改成功"))

            emits("update:modelValue")
            emits("commit")
        })


}

</script>

<style scoped>


</style>
