
<template>
    <el-dialog v-model="dialogVisible" :title="$t('调控玩家')" width="40%" @close="emits('update:modelValue', false)">
        <el-form :model="submitUsers">
            <el-form-item :label="$t('调控玩家')" label-width="100px">


                <div>

                    <span v-for="(item, index) in props.PreSubmitPlayer">
                        <el-tag type="primary" style="margin-right: 10px">

                            【{{ item.AppID }}】{{item.Uid}}
                        </el-tag>
                    </span>
                </div>

            </el-form-item>
        </el-form>

        <PlayerControl ref="PlayerControlRef" marginBottom="0" :highRTP="props.highRTP"></PlayerControl>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="emits('update:modelValue', false)">{{ $t('关闭') }}</el-button>
                <el-button type="primary" style="margin-left: 10px" @click="commitRTPControl">
                    {{ $t('提交') }}
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">

import {computed, ref} from "vue";
import PlayerControl from "@/pages/analysis/component/PlayerControl.vue";
import {ElLoading, ElMessageBox} from "element-plus";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const props = defineProps(["modelValue", "PreSubmitPlayer", "commitController", "highRTP"])
const emits = defineEmits(["update:modelValue", "commitController"])
const PlayerControlRef = ref(null)
const dialogVisible = computed(()=>{

    if (!props.modelValue){
        if (PlayerControlRef.value && PlayerControlRef.value.init){
            PlayerControlRef.value.init()
        }
    }

    return props.modelValue
})

const submitUsers = ref({
    GameID: "",
    ContrllRTP: 0,
    AutoRemoveRTP: 0,
})

const commitRTPControl = async () => {
    ElMessageBox.confirm(
        t('确认对用户进行RTP设置'),
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
            const tag = await PlayerControlRef.value.RTPConfig(props.PreSubmitPlayer)
            loading.close()
            if (tag) {


                emits("commitController")
            }
        })
}

</script>

<style scoped lang="scss">

</style>
