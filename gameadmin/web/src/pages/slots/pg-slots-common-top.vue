<template>
    <div class="title">
        <el-popover placement="right" :width="400" trigger="hover">
            <template #reference>
                        <span class="table-icon">
                            <el-icon><InfoFilled/></el-icon>
                        </span>
            </template>
            <div class="info">
                <img :src="dyImport(`${props.DiaLogGame}/info_page.png`)"/>
            </div>
        </el-popover>
    </div>
    <template v-if="data.Turns">
        <el-dropdown  @command="handleCommand" class="text-center" :max-height="300">
                    <span class="el-dropdown-link">
                      旋转定位
                      <el-icon class="el-icon--right">
                        <arrow-down/>
                      </el-icon>
                    </span>
            <template #dropdown>
                <el-dropdown-menu style="width: 300px">
                    <el-dropdown-item v-for="(item,index) in datas.Turns" :command="index" class="dropdown-item-style">
                        <span>N0.{{ index + 1 }}</span>
                        <span style="float: right">THB:{{ item.PanG / 10000 }}</span>
                    </el-dropdown-item>
                </el-dropdown-menu>
            </template>
        </el-dropdown>

        <span class="text-center">当前位置: N0.{{ datas.myCarouselSelectCommand }}/{{ datas.Turns.length }}</span>
    </template>
</template>

<script setup lang="ts">
import { dyImport, getImages } from "./allimage";
import {defineProps, reactive,ref} from "vue";
interface Props {
    PlayResp: any  // 盘面
    DiaLogGame: any  // 盘面
}
const props = defineProps(['DiaLogGame','data'])
let datas = reactive({
    Turns: props.data.Turns,
    images: getImages(props.DiaLogGame),
    myCarouselSelectCommand: 1,
})
const myCarousel = ref(null)
const handleCommand = (command: number) => {
    myCarousel.value.setActiveItem(command);

    datas.myCarouselSelectCommand = command+1
}
</script>

<style scoped lang="less">
.title{
    display: flex;
    justify-content: flex-end;
    align-items: center;
    font-size: 16px;
    color: rgb(255, 200, 36);
    .table-icon{
        margin: 0;
        margin-left: .5rem;
    }
}
</style>
