

<template>

    <el-upload :action="actionUrl"
               class="avatar-uploader"
               :auto-upload="true" ref="uploadRef"
               :file-list="fileList"
               :http-request="uploadGameIcon"
               :limit="1"
               :show-file-list="false">
       <slot></slot>
    </el-upload>
</template>

<script setup lang="ts">

import {Plus} from "@element-plus/icons-vue";
import {ref} from "vue";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

let props = defineProps(["uploadFileKeys"])

let emit = defineEmits(["uploadFile"])


let fileList = ref([])
let actionUrl:string = "http://action.com"

const uploadGameIcon = (file) => {

    fileList.value = []
    let fileMaxSize = 5120;//1M
    let fileSize = (file.file.size) / 5120
    let fileSizeBool = fileSize > fileMaxSize || fileSize <= 0
    if (file.file.type !== "image/png" && file.file.type !== "image/webp" && fileSizeBool) {
        tip.e(t("请选择小于100M的png图片"))
        return
    }

    emit("uploadFile", {file: file.file, key: props.uploadFileKeys})
}

</script>

<style scoped lang="scss">

</style>

<style>

.gamePage .avatar-uploader .el-upload {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);
}

.gamePage .avatar-uploader .el-upload:hover {
    border-color: var(--el-color-primary);
}

.gamePage .el-icon.avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 80px;
    height: 80px;
    text-align: center;
}
</style>
