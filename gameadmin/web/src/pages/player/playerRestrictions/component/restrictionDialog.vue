<template>
    <el-dialog
            v-model="dialogVisible"
            :title="props.RestrictionData == 'MoniterConfig' ? $t('限值玩家'): $t('修改限值')"
            :width="store.viewModel === 2 ? '85%' : '650px'"
            @close="closeDialog"
            height="200px"
            destroy-on-close
            align-center
            @open-auto-focus="openDialog"
    >
        <div>


            <el-form
                    ref="commitFormRef"
                    :model="commitForm"
                    :rules="commitRules"
                    label-width="160px"
                    style="max-width: 100%"
            >


                <el-form-item :label="$t('玩家ID')" prop="Pids">
                    <template #label>

                        {{ $t('玩家ID') }}
                    </template>
                    <el-input type="textarea" clearable v-model.trim="commitForm.Pids"
                              :placeholder="$t('请输入')"
                              resize="none"
                              :disabled="props.RestrictionData && props.RestrictionData.Pid"
                              onkeyup="value=value.replace(/[^\d,]/g, '').replace(/^,/, '');"/>

                </el-form-item>

                <el-form-item prop="RestrictionsMaxWin">
                    <template #label>

                        <table-tips :tips="$t('注意：购买小游戏或freegame时不会受到限制，只计算slots游戏的数值')"/>
                        <div class="customLabel">{{ $t('最多累积赢钱总数') }}</div>
                    </template>
                    <el-input clearable v-model.trim.number="commitForm.RestrictionsMaxWin"
                              :placeholder="$t('请输入')"
                              onkeyup="Number(value=value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1'))"/>
                </el-form-item>
                <el-form-item prop="RestrictionsMaxMulti">
                    <template #label>
                        <table-tips :tips="$t('注意：购买小游戏或freegame时不会受到限制，只计算slots游戏的数值')"/>
                        <div class="customLabel">{{ $t('最多累积赢取总倍数') }}
                        </div>
                    </template>
                    <el-input clearable v-model.trim.number="commitForm.RestrictionsMaxMulti"
                              :placeholder="$t('请输入')"
                              onkeyup="Number(value=value.replace(/^\D*(\d*(?:\.\d{0,3})?).*$/g,'$1'))"/>
                </el-form-item>

            </el-form>


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
import tableTips from "@/components/customTable/tableTips.vue"
import {computed, reactive, ref} from "vue";
import {useStore} from "@/pinia";
import {AdminGameCenter, RestrictionParams} from "@/api/gamepb/admin";
import {Client} from "@/lib/client";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {ElMessageBox, FormRules} from "element-plus";

const props = defineProps(["modelValue", "RestrictionData"])
const emits = defineEmits(["update:modelValue", "updateFrom"])
const store = useStore()
const {t} = useI18n()


const dialogVisible = computed(() => {
    return props.modelValue
})

const commitRules = reactive<FormRules>({
    Pids: [
        {required: true, message: t("玩家ID不能为空"), trigger: "blur"},
    ],
    RestrictionsMaxWin: [
        {
            message: t("最多累积赢钱总数不能为空"), trigger: "blur", validator: (rule: any, value: any, callback: any) => {
                if (value == "") {
                    callback(new Error(t("最多累积赢钱总数不能为空")))
                }


                callback()
            }
        },
    ],
    RestrictionsMaxMulti: [
        {
            message: t("最多累积赢取总倍数不能为空"),
            trigger: "blur",
            validator: (rule: any, value: any, callback: any) => {
                if (value == "") {
                    callback(new Error(t("最多累积赢取总倍数不能为空")))
                }
                callback()
            }
        },
    ],
})

const commitFormRef = ref(null)
const commitForm = ref({
    Pids: "",
    RestrictionsMaxWin: 0,
    RestrictionsMaxMulti: 0
})

const openDialog = () => {
    commitForm.value = {
        Pids: "",
        RestrictionsMaxWin: 0,
        RestrictionsMaxMulti: 0
    }
    if (props.RestrictionData) {
        commitForm.value = {
            Pids: props.RestrictionData.Pid.toString(),
            RestrictionsMaxWin: props.RestrictionData.RestrictionsMaxWin,
            RestrictionsMaxMulti: props.RestrictionData.RestrictionsMaxMulti,
        }

    }

}
const closeDialog = () => {
    emits("update:modelValue", false)
}
const commitData = async () => {

    commitFormRef.value.validate((A, ABA, C) => {
        if (A) {
            let setParam = <RestrictionParams>{
                ...commitForm.value
            }

            if (typeof setParam.Pids === "string") {
                setParam["Pids"] = setParam.Pids.split(",").map(item => parseInt(item)).filter(item => !isNaN(item) || item > 0)
            }


            if (setParam["Pids"].length <= 0) {
                tip.e(t("请填写玩家ID"))
                return
            }


            ElMessageBox.confirm(
                t('确认调整限制'),
                t('是否确认'),
                {
                    confirmButtonText: t('确定'),
                    cancelButtonText: t('关闭'),
                    type: 'warning',
                }).then(async () => {


                const [response, err] = await Client.Do(AdminGameCenter.SetPlayerRestrictions, setParam)

                if (err != null) {
                    tip.e(t("设置失败"))
                    return
                }
                if (response.NoExitPId != "") {
                    tip.w(t("无效的玩家ID:{Ids}", {Ids: response.NoExitPId}))
                }


                closeDialog()
                emits("updateFrom")
            })
        }
    })

}
</script>

<style scoped lang="scss">

.customLabel:before {
  color: var(--el-color-danger);
  content: "*";
  margin-right: 4px;
}
</style>
