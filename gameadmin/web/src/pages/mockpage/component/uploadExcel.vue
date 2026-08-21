

<template>

    <el-upload :action="actionUrl"
               class="avatar-uploader"
               :auto-upload="true" ref="uploadRef"
               :file-list="fileList"
               :http-request="uploadLangeConfig"
               :limit="100"
               :show-file-list="false">

        <template #trigger>
            <el-button type="primary" plain>上传文件</el-button>
        </template>
    </el-upload>


    <el-button type="primary" plain style="margin-left: 10px" @click="downloadFile">下载模板</el-button>
</template>

<script setup lang="ts">

import {ref} from "vue";
import {useI18n} from "vue-i18n";

import * as XLSX from 'xlsx'
import {useRoute} from "vue-router";


const {t} = useI18n()

let props = defineProps(["uploadFileKeys", "fileName"])

let emits = defineEmits(["uploadFile"])


let fileList = ref([])
let route = useRoute()
let actionUrl:string = "http://action.com"

const uploadLangeConfig = async (file) => {
    const reader = new FileReader();
    reader.onload = (e) => {
        const data = e.target.result;
        const workbook = XLSX.read(data, { type: 'binary' });
        const sheetName = workbook.SheetNames[0];
        const worksheet = workbook.Sheets[sheetName];
        const jsonData = XLSX.utils.sheet_to_json(worksheet);
        emits("uploadFile", {data: jsonData})
    };
    reader.readAsBinaryString(file.file);


}


const downloadFile = () => {
    const baseURL = window.location.protocol + "//" + window.location.host + "/excel/"
    let fileName = route.path.replace("/", "")

    const excelUri = `${baseURL}/${fileName}.xlsx`
    fetch(excelUri,{
        method:"get",
        headers: {
            'Content-Type': 'application/vnd.ms-excel',
        },
    })
        .then(response => response.blob())
        .then(blob => {
        let blobUrl = window.URL.createObjectURL(blob)

        let aDoc = document.createElement("a")
        aDoc.setAttribute("download",`${props.fileName}.xlsx`)
        aDoc.setAttribute("target","_blank")
        aDoc.setAttribute("href", blobUrl)


        document.body.appendChild(aDoc)

        aDoc.click()
        aDoc.remove()
        window.URL.revokeObjectURL(blobUrl); // 释放blob UR
    })

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
