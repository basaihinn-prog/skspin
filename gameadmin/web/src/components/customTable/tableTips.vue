
<template>
    <el-tooltip
            class="box-item"
            effect="dark"
            :open-delay="openDelay || 300"
            :hide-after="hideAfter || 300"
            placement="top-start"
    >

        <template #content>

            <div v-html="$t(tips)" style="white-space:pre-line;" ></div>

        </template>
        <el-icon size="20" v-if="tips">
            <QuestionFilled/>
        </el-icon>
    </el-tooltip>
</template>
<script setup lang="ts">

import {QuestionFilled} from "@element-plus/icons-vue";
import {useI18n} from "vue-i18n";

const { tips, tipsCustom, tipsSlotName, openDelay, hideAfter } = defineProps(["tips", "tipsCustom", "tipsSlotName", "openDelay", "hideAfter"])

const { t } = useI18n()
</script>

<style scoped lang="scss">

</style>
