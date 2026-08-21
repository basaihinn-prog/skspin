/*
 *@Author: 西西米
 *@Date: 2023-01-04 11:16:28
 *@Description: 文件列表
*/
<template>
    <div class='fileList'>
        
        <div class="file" v-for="(item, index) in imageList" :key="index">
            <el-image :src="item" :preview-src-list="imageList" :initial-index="index" loading="lazy" preview-teleported
                fit="cover" />
            <div class="delView">
                <el-icon color="#ff0000" size="20" @click="removeFile(index)">
                    <Delete />
                </el-icon>
            </div>
        </div>
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { AdminUploadFile } from '@/api/admin_upload_file'
import { useStore } from '@/pinia/index';
import uploadPictures from '@/components/upload-pictures.vue';
const store = useStore()

onMounted(() => {
    getFileList()
});

let imageList = ref([])

const getFileList = async () => {
    let [data, err] = await Client.Do(AdminUploadFile.FileList, { Path: 'avatar' })
    if (err) {
        return tip.e(err)
    }
    let list = []
    data.List.forEach((item) => {
        list.push(store.fileURL + item.Path)
    })
    imageList.value = list
}

const removeFile = async (index) => {
    let Path = imageList.value[index].split(store.fileURL)[1]
    let [data, err] = await Client.Do(AdminUploadFile.RemoveFile, { Path })
    if (err) {
        return tip.e(err)
    }
    imageList.value.splice(index, 1)
    tip.s('删除成功')
}
</script>
<style scoped lang='scss'>
.fileList {
    display: flex;
    row-gap: 10px;
    column-gap: 10px;
    flex-wrap: wrap;

    .file {
        width: 80px;
        height: 80px;
        cursor: pointer;
        position: relative;

        .el-image {
            width: 100%;
            height: 100%;
        }

        .delView {
            width: 30px;
            height: 30px;
            justify-content: center;
            align-items: center;
            background: rgba($color: #000000, $alpha: .3);
            position: absolute;
            top: 0;
            right: 0;
            display: none;
            column-gap: 15px;
        }

        &:hover {
            .delView {
                display: flex;
            }
        }
    }
}
</style>