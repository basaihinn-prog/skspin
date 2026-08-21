

<template>

    <el-upload :action="actionUrl"
               class="avatar-uploader"
               :auto-upload="true" ref="uploadRef"
               :file-list="fileList"
               :http-request="uploadLangeConfig"
               :limit="1"
               :show-file-list="false">

        <slot></slot>
    </el-upload>
</template>

<script setup lang="ts">

import {Plus} from "@element-plus/icons-vue";
import {ref} from "vue";
import {Client} from "@/lib/client";
import {Upload} from "@/api/comm";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {Language} from "@/api/adminpb/language";


const {t} = useI18n()

let props = defineProps(["uploadFileKeys", "fileName"])

let emits = defineEmits(["uploadFile"])


let fileList = ref([])
let actionUrl:string = "http://action.com"

const uploadLangeConfig = async (file) => {




    let formData = new FormData();

    formData.append('file', file.file);
    formData.append('fileName', props.fileName);


    const [response, err] = await Client.Do(Language.UploadConfig, formData as any)

    if (err){
        return tip.e(t("上传失败"))
    }

    emits("uploadFile", {file: file.file, key: props.uploadFileKeys})

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
