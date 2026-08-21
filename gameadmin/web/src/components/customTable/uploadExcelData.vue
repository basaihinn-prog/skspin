

<template>

    <el-dialog v-model="props.modelValue" title="提交数据" width="500">


        <el-form :model="dataForm">
            <el-form-item label="数据列表" label-width="150px">
                <el-link @click="downloadFile">获取模板</el-link>
            </el-form-item>
            <el-form-item label="模版下载" label-width="150px">
                <el-input v-model="dataForm.excelData"
                          :rows="6"
                          type="textarea"/>
            </el-form-item>


        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="closeDialog">关闭</el-button>
                <el-button type="primary" @click="uploadData">
                    提交
                </el-button>
            </div>
        </template>
    </el-dialog>

</template>

<script setup lang="ts">

import {ref} from "vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

let props = defineProps(["modelValue", "uploadFileKeys"])

let emit = defineEmits(["uploadData", "update:modelValue"])

const dataForm = ref({
    excelData:""
})

let fileList = ref([])
const downloadFile = () => {
    // 创建一个<a></a>标签
    const a = document.createElement('a');
    a.href = '/src/assets/index.xlsx';
    a.download = '号码池导入';
    a.style.display = 'none';
// 将a标签追加到文档对象中
    document.body.appendChild(a);
// 模拟点击了<a>标签,会触发<a>标签的href的读取,浏览器就会自动下载了
    a.click();
// 一次性的,用完就删除a标签
    a.remove();
}
const closeDialog = () =>{
    emit("update:modelValue")
}
const uploadData  = () => {
    emit("uploadData", dataForm.value.excelData)
    closeDialog()
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
