/*
 *@Author: 西西米
 *@Date: 2023-03-02 09:34:40
 *@Description:
*/
<template>
    <el-pagination :current-page="param.Page" :page-size="param.PageSize" :total="Count" :page-sizes="[20, 30, 60, 100]"
        :pager-count="5" :layout="store.paginationLayout" @size-change="handleSizeChange"
        @current-change="handleCurrentChange" class="elPagination" />
</template>

<script setup lang='ts'>
import { ref, reactive, onMounted } from 'vue';
import { useStore } from '@/pinia/index';
const store = useStore()

interface Props {
    Count: number
}
const props = withDefaults(defineProps<Props>(), {
    Count: 0
})

let param = reactive({
    Page: 1,
    PageSize: 20,
})

// onMounted( () => {
//     emit('paginationEmit', param.Page, param.PageSize)
// })

const emit = defineEmits<{
    (e: 'paginationEmit', Page: number, PageSize: number): void
}>()

const handleCurrentChange = async (page: number) => {
    param.Page = page
    emit('paginationEmit', param.Page, param.PageSize)
}

const handleSizeChange = (size: number) => {
    param.PageSize = size
    emit('paginationEmit', param.Page, param.PageSize)
}
</script>
<style scoped lang='scss'></style>
