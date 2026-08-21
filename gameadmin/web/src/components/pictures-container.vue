/*
 *@Author: 西西米
 *@Date: 2023-01-10 11:15:54
 *@Description: 选择图片
*/
<template>
    <div class='upload-pictures' @click="isPreservation = false, dialogVisible = true, getFileList()">
        <div class="imgList" v-if='imgList.length - 1'>
            <template v-for="(item, index) in imgList">
                <div class="imgView" :key="index" v-if="item !== 'draggable'">
                    <el-image :src="imageUrl(item)" class="imgView" fit="scale-down" />
                </div>
            </template>
        </div>
        <el-button v-else type="primary" plain>{{$t('选择图片')}}</el-button>
    </div>
    <!-- <el-button type="primary" @click="isPreservation = false, dialogVisible = true, getFileList()" plain>
                选择图片
                <span v-if='imgList.length - 1'>/当前选择{{ imgList.length - 1 }}张</span>
            </el-button> -->
    <el-dialog v-model="dialogVisible" :title="$t('选择图片')" class="file-container" @close="resetForm"
        :width="store.viewModel === 2 ? '85%' : '50%'">
        <div class='upload-pictures'>
            <draggable :list="imgList" :move="onMove" filter=".forbid" ghost-class="imgList" :force-fallback="true"
                class="imgList" chosen-class="chosenClass" animation="300" item-key="item">
                <template #item="{ element, index }">
                    <div class="imgView" v-if="element !== 'draggable'">
                        <el-image :src="imageUrl(element)" class="imgView" fit="scale-down" />
                        <div class="delView">
                            <el-icon color="#ff0000" size="20" @click="delImg(index)">
                                <Delete />
                            </el-icon>
                        </div>
                    </div>
                    <div class="forbid upload" v-else-if="props.multiple || (!props.multiple && imgList.length === 1)">
                        <el-icon>
                            <Plus />
                        </el-icon>
                        <input type="file" :multiple="props.multiple" accept="image/png, image/jpeg" @change="changeFile">
                    </div>
                </template>
            </draggable>
        </div>

        <p class="existence">{{$t('当前路径')}}: </p>
        <el-breadcrumb separator-icon="ArrowRight">
            <el-breadcrumb-item @click="currPath = [], getFileList()">/</el-breadcrumb-item>
            <el-breadcrumb-item v-for="(item, index) in currPath" @click="filePathTag(index)">
                {{ item }}
            </el-breadcrumb-item>
        </el-breadcrumb>
        <div class="filePath">
            <p v-for="item in filePathList" @click="filePathClick(item.Path)">{{ item.Path }}</p>
        </div>

        <p class="existence">{{$t('已存在文件')}}</p>
        <div class="upload-pictures">
            <div class="imgList">
                <div class="imgView" v-for="(item, index) in existenceList" :key="index">
                    <el-image @click="pushImage(item.Path)" :src="imageUrl(item.Path)" class="imgView" loading="lazy"
                        fit="scale-down" />
                    <div class="delView">
                        <el-icon color="#ff0000" size="20" @click="removeFile(index)">
                            <Delete />
                        </el-icon>
                    </div>
                </div>
            </div>
        </div>

        <template #footer>
            <span class="dialog-footer">
                <el-button @click="dialogVisible = false">{{$t('取消')}}</el-button>
                <el-button type="primary" @click="confirm">{{$t('确定')}}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive, watch } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { AdminUploadFile } from '@/api/admin_upload_file'
import { useStore } from '@/pinia/index';
import draggable from "vuedraggable";
const dialogVisible = ref(false)
const store = useStore()

interface Props {
    multiple?: boolean,
    filelist?: string,
    name?: string,
    uploadType?: number
}
const props = withDefaults(defineProps<Props>(), {
    multiple: true,
    filelist: '',
    name: '',
    uploadType: 0
})

const uploadType = {
    0: 'gameIcon', //游戏图标类图片
    1: 'notice', //公告类图片
    2: 'goods', //积分商城类图片
    3: 'rule', //规则文本
}

let imgList = ref<string[]>([])
let isPreservation = ref(true)

watch(() => props.filelist, (newVal) => {
    imgList.value = newVal === '' ? ['draggable'] : newVal.split(',').concat(['draggable'])
}, { immediate: true, deep: true })

const emit = defineEmits<{
    (e: 'pictureCallback', name: string, url: string): void
}>()

const changeFile = async (file: any) => {
    for (let i = 0; i < file.target.files.length; i++) {
        let formData = new FormData()
        formData.append('file', file.target.files[i])
        formData.append('path', uploadType[props.uploadType])
        let [data, err] = await Client.send("AdminUpload", formData)
        if (err) {
            return tip.e(err)
        }
        imgList.value.splice(imgList.value.length - 1, 0, data)
    }
}

const onMove = (e: any) => {
    if (e.relatedContext.element === 'draggable') return false;
    return true;
}

const delImg = (index: number) => {
    imgList.value.splice(index, 1)
}

const confirm = () => {
    let text1 = imgList.value.join(',')
    let text2 = props.filelist.split(',').concat(['draggable']).join(',')
    if (text1 !== text2) { //改变
        imgList.value.pop()
        emit('pictureCallback', props.name, imgList.value.join(','))
    }
    isPreservation.value = true
    dialogVisible.value = false
}

const resetForm = () => {
    if (!isPreservation.value) {
        imgList.value = props.filelist === '' ? ['draggable'] : props.filelist.split(',').concat(['draggable'])
        isPreservation.value = true
    }
}

let existenceList = ref([])
let filePathList = ref([])
let currPath = ref([])
const getFileList = async () => {
    let Path = currPath.value.length === 0 ? '' : currPath.value[currPath.value.length - 1]
    let [data, err] = await Client.Do(AdminUploadFile.FileList, { Path })
    if (err || !data.List) {
        tip.e(err)
        data.List = []
    }
    let file = []
    let filePath = []
    data.List.forEach(item => {
        if (item.IsDir) {
            filePath.push(item)
        } else {
            file.push(item)
        }
    })
    existenceList.value = file
    filePathList.value = filePath
}

const filePathTag = async (index) => {
    currPath.value.splice(index + 1, currPath.value.length)
    getFileList()
}

const filePathClick = async (path) => {
    currPath.value.push(path.split(':')[1])
    getFileList()
}

const removeFile = async (index) => {
    let isRemove = await tip.ask("确定删除此图片吗？")
    if (isRemove === 'ok') {
        let Path = existenceList.value[index].Path
        let [data, err] = await Client.Do(AdminUploadFile.RemoveFile, { Path })
        if (err) {
            return tip.e(err)
        }
        existenceList.value.splice(index, 1)
        tip.s('删除成功')
    }
}

const pushImage = (path) => {
    if (props.multiple) {
        imgList.value.splice(imgList.value.length - 1, 0, path)
    } else {
        imgList.value = [path, 'draggable']
    }
}

</script>
<style scoped lang='scss'>
.upload-pictures .imgList {
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    column-gap: 10px;
    row-gap: 10px;
    max-height: 440px;
    overflow-y: scroll;

    .imgView,
    .upload {
        width: 90px;
        height: 90px;
        position: relative;
        overflow: hidden;
        box-sizing: border-box;

        &:hover .delView {
            display: flex;
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
    }

    .upload {
        border: 1px dashed #dcdfe6;
        background: #fafafa;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 20px;
        color: #dcdfe6;

        input {
            cursor: pointer;
            width: 100%;
            height: 100%;
            position: absolute;
            top: 0;
            left: 0;
            opacity: 0;
        }

        .progress {
            position: absolute;
            background: #fafafa;
            top: 0;
            left: 0;
        }
    }
}

.existence {
    padding: 10px 0;
    font-size: 16px;
}
</style>