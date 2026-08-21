<template>
    <div class="boxcolor">
        <div>
            <img class="info" :src="dyImport(`${props.DiaLogGame}/page_info.jpg`)" />
        </div>
        <el-carousel :autoplay="false" height="400px" trigger="click" v-if="data.Turns">
            <el-carousel-item :key="'p'+panid" v-for="(pan,panid) in data.Turns">
                <div class="flex_vbox">
                    <SlotPan :images="data.images" :Pan="pan" :highlight="(x,y) => hightLight(panid, x, y)" />
                    <div style="height:20px" />
                    <div class="flex_vbox"></div>
                </div>
            </el-carousel-item>
        </el-carousel>
    </div>
</template>

<style lang="less" scoped>
.block_img {
    max-width: 30px;
    max-height: 30px;
    min-width: 30px;
    min-height: 30px;
}

.weapon {
    max-width: 20px;
    height: auto;
}
.info {
    max-width: 100%;
    height: auto;
}
.line {
    cursor: pointer;
}
.carousel_wh {
    height: 300px !important;
}
.boxcolor {
    background-color: #30303c;
}
</style>

<script lang="ts" setup>
import img1 from "@/assets/niubi/1.png"
import img2 from "@/assets/niubi/2.png"
import img3 from "@/assets/niubi/3.png"
import img4 from "@/assets/niubi/4.png"

import SlotPan from "./SlotPan.vue";
import { onMounted, onUpdated, reactive } from "vue";
import { dyImport, getImages } from "./allimage";
interface Props {
    PlayResp: any  // 盘面
    DiaLogGame: any  // 盘面
}
const props = withDefaults(defineProps<Props>(), {})
let data = reactive({
    Turns: null,
    images: getImages(props.DiaLogGame),
})

let hightLight = (pan, x,y)=> {
    return false
}

let updateData = () => {
    let turns = []
    turns.push(props.PlayResp.Pan)

    data.Turns = turns
}
// onUpdated(updateData)
onMounted(updateData)
</script>
