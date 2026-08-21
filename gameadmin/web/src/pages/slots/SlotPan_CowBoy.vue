<template>
    <div class="flex_vbox flex_child_center">
        <div v-for="(rowData, rowIndex) in newPanListData" :key="'row-' + rowIndex" class="flex_hbox row">
            <div v-for="(imgId, colIndex) in rowData" :key="'col-' + colIndex" class="block">
                <template v-for="items in hitLines">
                    <template v-if="items === imgId">
                        <div class="borders"></div>
                    </template>
                </template>
                <img :src="images[imgId.data]"/>
                <img v-if="imgId.star" :src="star" class="star">
                <span v-if="imgId.baiDaMul && imgId.baiDaMul===2" class="baiDaMul blue">x{{ imgId.baiDaMul }}</span>
                <span v-if="imgId.baiDaMul && imgId.baiDaMul===3" class="baiDaMul purple">x{{ imgId.baiDaMul }}</span>
                <span v-if="imgId.baiDaMul && imgId.baiDaMul===5" class="baiDaMul red">x{{ imgId.baiDaMul }}</span>
            </div>
        </div>
    </div>
</template>

<script setup>

import {defineProps} from "vue"
// import { dyImport } from "./allimage";
// import { DocBetLog } from "@/api/gamepb/customer";
import {useCssModule} from 'vue'
import star from '@/assets/CowBoyMap/1.png'
/*interface Props {
    Pan: number[][]// 盘面
    images: JsMap<string, string>  // ID索引图片
    highlight: (x, y, idx) => boolean
}*/
const props = defineProps(['Di', 'Pan', 'images', 'highlight','hitLines'])
const styles = useCssModule()
const newPanListData = props.Pan.Data
const newPanListStar = props.Pan.Star
const newPanListBaiDaMul = props.Pan.BaiDaMul
const init = () => {
    for (let i = 0; i < newPanListData.length; i++) {
        for (let j = 0; j < newPanListStar[0].length; j++) {
            newPanListData[i][j] = {
                data: newPanListData[i][j],
                star: newPanListStar[i][j],
                baiDaMul: newPanListBaiDaMul[i][j],
            };
        }
    }
}

init()
</script>

<style scoped>
.flex {
    display: flex;
    display: -webkit-flex;
}

.flex_vbox {
    display: flex;
    display: -webkit-flex;
    flex-direction: column;
}

.flex_hbox {
    display: flex;
    display: -webkit-flex;
    flex-direction: row;
}

.flex_grow {
    flex-grow: 1;
}

.star {
    max-width: 20px;
    position: absolute;
    right: 0;
}
.borders{
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    border: 2px solid #f9fcff;
}
.flex_child_center {
    align-items: center;
    justify-content: center;
}

.flex_self_center {
    align-self: center;
}

.fill_height {
    height: 100%;
}

.fill_width {
    width: 100%;
}

.block {
    max-width: 60px;
    max-height: 60px;
    min-width: 60px;
    min-height: 60px;
    width: 60px;
    height: 60px;
    border: 3px solid transparent;
    position: relative;
}

.block img {
    width: 100%;
    height: auto;
}

.block span {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 1.7rem;
    font-weight: bold;
    text-shadow: 0 0 0 #ffe351, 0 1px 0 #e99707, 0 2px 0 #e99707, 0 3px 0 #e99707;
    -webkit-background-clip: text;
    /*color: #f8d94d;*/
    color: transparent;
    background-image: linear-gradient(to top, #ffb125 20%, #f7fc00 50%, #e99707 100%);
    -webkit-text-stroke: .5px #242622;
}

.block .baiDaMul {
    text-shadow: none;
    /*-webkit-text-stroke:0*/

    transform: rotate(30deg);
    font-size: 1.2rem;
    left: 60%;
}
.block .blue{
    background-image: linear-gradient(180deg, rgba(237, 237, 237, 1) 0%, rgba(23, 134, 219, 1) 100%);
}
.block .purple{
    background-image: linear-gradient(180deg, rgba(237, 237, 237, 1) 0%, rgb(101, 23, 219) 100%);
}
.block .red{
    background-image: linear-gradient(180deg, rgba(237, 237, 237, 1) 0%, rgb(219, 69, 23) 100%);
}
.highlight {
    border: 3px solid red;
    animation-name: upAnimation;
    transform-origin: center bottom;
    animation-duration: 2s;
    animation-fill-mode: both;
    animation-iteration-count: infinite;
    animation-delay: .5s;
}

@keyframes upAnimation {
    0% {
        transform: rotate(0deg);
        transition-timing-function: cubic-bezier(0.215, .61, .355, 1)
    }

    10% {
        transform: rotate(-12deg);
        transition-timing-function: cubic-bezier(0.215, .61, .355, 1)
    }

    20% {
        transform: rotate(12deg);
        transition-timing-function: cubic-bezier(0.215, .61, .355, 1)
    }

    28% {
        transform: rotate(-10deg);
        transition-timing-function: cubic-bezier(0.215, .61, .355, 1)
    }

    36% {
        transform: rotate(10deg);
        transition-timing-function: cubic-bezier(0.755, .5, .855, .06)
    }

    42% {
        transform: rotate(-8deg);
        transition-timing-function: cubic-bezier(0.755, .5, .855, .06)
    }

    48% {
        transform: rotate(8deg);
        transition-timing-function: cubic-bezier(0.755, .5, .855, .06)
    }

    52% {
        transform: rotate(-4deg);
        transition-timing-function: cubic-bezier(0.755, .5, .855, .06)
    }

    56% {
        transform: rotate(4deg);
        transition-timing-function: cubic-bezier(0.755, .5, .855, .06)
    }

    60% {
        transform: rotate(0deg);
        transition-timing-function: cubic-bezier(0.755, .5, .855, .06)
    }

    100% {
        transform: rotate(0deg);
        transition-timing-function: cubic-bezier(0.215, .61, .355, 1)
    }
}
</style>
