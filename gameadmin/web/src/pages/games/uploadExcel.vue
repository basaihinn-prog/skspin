

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
import {Client} from "@/lib/client";
import {Upload} from "@/api/comm";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {ElLoading} from "element-plus";


const {t} = useI18n()

let props = defineProps(["uploadFileKeys", "requestUri"])

let emit = defineEmits(["uploadFile"])


let fileList = ref([])
let actionUrl:string = "http://action.com"

const uploadGameIcon = async (file) => {

    const loading = ElLoading.service({
        lock: true,
        text: '上传更新游戏中,请勿操作...',
        background: 'rgba(0, 0, 0, 0.7)',
    })
    let formData = new FormData();

    formData.append('file', file.file);


    const [response, err] = await Client.Do(AdminGameCenter.UploadGameFile, formData as any)

    if (err && err.indexOf('上传成功') == -1){
        loading.close()
        return tip.e(t("上传失败"))

    }

    loading.close()
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
