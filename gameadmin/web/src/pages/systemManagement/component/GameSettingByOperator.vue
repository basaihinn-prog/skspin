<template>
    <el-dialog :title="$t('Slots游戏设置')"
               align-center
               v-model="props.modelValue"
               destroy-on-close
               @open-auto-focus="openDialog"
               @close="emits('update:modelValue')"
               :width="store.viewModel === 2 ? '100%' : '650px'">
        <el-descriptions :column="1" border>
            <el-descriptions-item :label="$t('PP游戏设置')" v-if=" GameManufacturer.indexOf('PP') != -1 && (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf('PP') != -1)">
                <el-form>
                    <el-form-item :label="`${$t('时间和游戏展示')}:`" style="margin-bottom: 0">
                        <el-select
                            v-model="PPSubmitConfig.ShowNameAndTimeOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>

                    </el-form-item>
                    <!--                    <el-form-item :label="`${$t('游戏货币符号显示')}:`">-->
                    <!--                        <el-select-->
                    <!--                            v-model="PPSubmitConfig.CurrencyVisibleOff"-->
                    <!--                            :placeholder="$t('选择状态')"-->
                    <!--                            style="width: 100px;"-->


                    <!--                        >-->
                    <!--                            <el-option-->
                    <!--                                v-for="(item, key) in statusEnmu"-->
                    <!--                                :key="key"-->
                    <!--                                :label="$t(item.label)"-->
                    <!--                                :value="item.value"-->
                    <!--                            />-->
                    <!--                        </el-select>-->

                    <!--                    </el-form-item>-->
                </el-form>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('PG游戏设置')"  v-if=" GameManufacturer.indexOf('PG') != -1 && (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf('PG') != -1)">
                <el-form>

                    <el-form-item :label="`${$t('首页轮播')}:`">
                        <el-select
                            v-model="PGSubmitConfig.CarouselOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>

                    </el-form-item>

                    <el-form-item :label="`${$t('时间和游戏展示')}:`">
                        <el-select
                            v-model="PGSubmitConfig.ShowNameAndTimeOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>

                    </el-form-item>

                    <!--                    <el-form-item :label="`${$t('游戏货币符号显示')}:`">-->
                    <!--                        <el-select-->
                    <!--                            v-model="PGSubmitConfig.CurrencyVisibleOff"-->
                    <!--                            :placeholder="$t('选择状态')"-->
                    <!--                            style="width: 100px;"-->


                    <!--                        >-->
                    <!--                            <el-option-->
                    <!--                                v-for="(item, key) in statusEnmu"-->
                    <!--                                :key="key"-->
                    <!--                                :label="$t(item.label)"-->
                    <!--                                :value="item.value"-->
                    <!--                            />-->
                    <!--                        </el-select>-->

                    <!--                    </el-form-item>-->

                    <el-form-item :label="`${$t('止盈止损是否显示')}:`">
                        <el-select
                            v-model="PGSubmitConfig.StopLoss"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in StopLossStatusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>

                    </el-form-item>

                    <el-form-item :label="`${$t('菜单内退出按钮')}:`">
                        <el-select
                            v-model="PGSubmitConfig.ExitBtnOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>


                    </el-form-item>
                    <el-form-item :label="`${$t('菜单内退出按钮链接')}:`" v-if="PGSubmitConfig.ExitBtnOff">
                        <el-input v-model="PGSubmitConfig.ExitLink" style="width: 150px;margin-left: 10px" :placeholder="$t('请输入退出按钮链接')"/>

                    </el-form-item>
                    <el-form-item :label="`${$t('是否开启官方验证')}:`">
                        <el-select
                            v-model="PGSubmitConfig.OfficialVerify"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>

                    </el-form-item>
                </el-form>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('JILI游戏设置')" v-if=" GameManufacturer.indexOf('JILI') != -1 && (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf('JILI') != -1)">
                <el-form>

                    <el-form-item :label="`${$t('背包按钮显示')}:`">
                        <el-select
                            v-model="JILISubmitConfig.BackPackOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item :label="`${$t('开屏提示')}:`">
                        <el-select
                            v-model="JILISubmitConfig.OpenScreenOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item :label="`${$t('弹窗和侧边栏')}:`">
                        <el-select
                            v-model="JILISubmitConfig.SidebarOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>
                </el-form>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('TADA游戏设置')"  v-if=" GameManufacturer.indexOf('TaDa') != -1 && (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf('TaDa') != -1)">
                <el-form>

                    <el-form-item :label="`${$t('背包按钮显示')}:`">
                        <el-select
                            v-model="TADASubmitConfig.BackPackOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>


                    <el-form-item :label="`${$t('开屏提示')}:`">
                        <el-select
                            v-model="TADASubmitConfig.OpenScreenOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"
                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>

                    <el-form-item :label="`${$t('弹窗和侧边栏')}:`">
                        <el-select
                            v-model="TADASubmitConfig.SidebarOff"
                            :placeholder="$t('选择状态')"
                            style="width: 100px;"


                        >
                            <el-option
                                v-for="(item, key) in statusEnmu"
                                :key="key"
                                :label="$t(item.label)"
                                :value="item.value"
                            />
                        </el-select>
                    </el-form-item>
                </el-form>
            </el-descriptions-item>

        </el-descriptions>
        <div slot="footer" class="dialog-footer flex flex_child_end" style="margin-top: 15px">
            <el-button @click="emits('update:modelValue')">{{ $t('取消') }}</el-button>
            <el-button type="primary" @click="submitData">{{ $t('确定') }}</el-button>
        </div>
    </el-dialog>
</template>

<script setup lang="ts">

import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {Ref, ref} from "vue";
import {AdminInfo} from "@/api/adminpb/info";
import type { JILIConfig, PGConfig, PPConfig, UpdateGameUniversalConfigParam} from "@/api/adminpb/info";
import {Client} from "@/lib/client";
import {tip} from "@/lib/tip";
import {ElLoading} from "element-plus";
import {Manufacturer} from "@/api/gamepb/game";

const {t} = useI18n()
const store = useStore()


const props = defineProps(['modelValue', 'gameSettingInfoData'])
const emits = defineEmits(["update:modelValue"])

const GameManufacturer = ref([])
const loading = ref(false)
const statusEnmu = [{
    label: '开启',
    value: 1
},{
    label: '关闭',
    value: 0
}]

const StopLossStatusEnmu = [{
    label: '开启',
    value: 0
},{
    label: '关闭',
    value: 1
}]


const PPSubmitConfig:Ref<PPConfig> = ref<PPConfig>(<PPConfig> {
    ShowNameAndTimeOff: 1,
    CurrencyVisibleOff: 1,
})

const PGSubmitConfig:Ref<PGConfig> = ref<PGConfig>(<PGConfig>{
    CarouselOff: 1,
    ShowNameAndTimeOff: 1,
    CurrencyVisibleOff: 1,
    StopLoss: 0,
    OfficialVerify: 1,
    ExitLink: "http://",
})
const JILISubmitConfig:Ref<JILIConfig> = ref<JILIConfig>(<JILIConfig>{
    BackPackOff: 1,
    CurrencyVisibleOff: 1,
    SidebarOff: 0
})
const TADASubmitConfig:Ref<JILIConfig> = ref<JILIConfig>(<JILIConfig>{
    BackPackOff: 1,
    CurrencyVisibleOff: 1,
    SidebarOff: 0,
})

const openDialog = async () => {

    let [res, err ] = await Client.Do(AdminInfo.GetOperatorInfo, {AppID:props.gameSettingInfoData.AppID})


    PPSubmitConfig.value =  {
        ShowNameAndTimeOff: res.OperatorInfo.PPConfig.ShowNameAndTimeOff,
        CurrencyVisibleOff: res.OperatorInfo.PPConfig.CurrencyVisibleOff,
    }

    PGSubmitConfig.value = {
        CarouselOff: res.OperatorInfo.PGConfig.CarouselOff,
        ShowNameAndTimeOff: res.OperatorInfo.PGConfig.ShowNameAndTimeOff,
        CurrencyVisibleOff: res.OperatorInfo.PGConfig.CurrencyVisibleOff,
        StopLoss: res.OperatorInfo.PGConfig.StopLoss,
        ExitBtnOff: res.OperatorInfo.PGConfig.ExitBtnOff,
        ExitLink: res.OperatorInfo.PGConfig.ExitBtnOff ? res.OperatorInfo.PGConfig.ExitLink : "",
        OfficialVerify: res.OperatorInfo.PGConfig.OfficialVerify,
    }
    JILISubmitConfig.value = {
        BackPackOff: res.OperatorInfo.JILIConfig.BackPackOff,
        CurrencyVisibleOff: res.OperatorInfo.JILIConfig.CurrencyVisibleOff,
        OpenScreenOff:res.OperatorInfo.JILIConfig.OpenScreenOff,
        SidebarOff:res.OperatorInfo.JILIConfig.SidebarOff,
    }
    TADASubmitConfig.value = {
        BackPackOff: res.OperatorInfo.TADAConfig.BackPackOff,
        CurrencyVisibleOff: res.OperatorInfo.TADAConfig.CurrencyVisibleOff,
        OpenScreenOff:res.OperatorInfo.TADAConfig.OpenScreenOff,
        SidebarOff:res.OperatorInfo.TADAConfig.SidebarOff,
    }
}

const initGameManufacturer = async () => {
    const [res, err] = await Client.Do(Manufacturer.GetManufacturerList, {} as any)

    GameManufacturer.value = res.List.map(item => item.ManufacturerCode)
}
initGameManufacturer()
const submitData = async () => {

    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    let params = <UpdateGameUniversalConfigParam>{
        AppID: props.gameSettingInfoData.AppID,


    }
    if (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf("PP") != -1){
        params.PPConfig =  PPSubmitConfig.value
    }


    if (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf("PG") != -1){
        params.PGConfig =  PGSubmitConfig.value
    }


    if (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf("JILI") != -1){
        params.JILIConfig =  JILISubmitConfig.value
    }


    if (!props.gameSettingInfoData.DefaultManufacturerOn || props.gameSettingInfoData.DefaultManufacturerOn.indexOf("TaDa") != -1){
        params.TADAConfig =  TADASubmitConfig.value
    }

    const [res, err] = await Client.Do(AdminInfo.UpdateGameUniversalConfig, params)

    loading.close()
    emits('update:modelValue')
    if (err){
        tip.e(err)
        return
    }
    tip.s(t('修改成功'))
    emits("update:modelValue")

}

</script>
<style scoped lang="scss">

</style>
