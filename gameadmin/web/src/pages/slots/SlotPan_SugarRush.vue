
<template>
    <div class="flex_vbox flex_child_center">
        <div v-for="(rowData, rowIndex) in newPanList" :key="'row-' + rowIndex" class="flex_hbox row">
            <div v-for="(imgId, colIndex) in rowData" :key="'col-' + colIndex" class="block" :class="{ 'highlight': highlight(rowIndex, colIndex, colIndex) }">
                <template v-if="images[imgId.nums]">
                    <template v-if="arrs[imgId.multi] === undefined">
                        <img :src="images[imgId.nums]" />
                    </template>
                    <template v-if="arrs[imgId.multi] === null">
                        <div class="opacity"></div>
                        <img :src="images[imgId.nums]" />
                    </template>
                    <template v-if="arrs[imgId.multi] > 1">
                        <div class="bg">
                            <img :src="img" alt="">
                        </div>
                        <img :src="images[imgId.nums]" />
                        <span>x{{imgId.multi}}</span>
                    </template>
                </template>

                <template v-else>
                    <div class="block"></div>
                </template>
            </div>
        </div>
    </div>
</template>

<script setup>

import {defineProps} from "vue"
// import { dyImport } from "./allimage";
// import { DocBetLog } from "@/api/gamepb/customer";
import { useCssModule } from 'vue'
import img from '@/assets/SugarRush/0.png'
/*interface Props {
    Pan: number[][]// 盘面
    images: JsMap<string, string>  // ID索引图片
    highlight: (x, y, idx) => boolean
}*/

const props = defineProps(['Di','Pan','images','highlight'])
const styles = useCssModule()
const newPanListData = props.Pan.Data
const newPanListMulti = props.Pan.Multi
let newPanList = []
let arrs = [
    undefined,
    null,
    2,
    4,
    8,
    16,
    32,
    64,
    128,
]
const init = () => {
    newPanList = newPanListData.map((data, i) => {
        if (data.length === newPanListMulti[i].length) {
            return data.map((num, j) => ({
                nums: num,
                multi: newPanListMulti[i][j],
            }));
        }
        return [];
    });

}

init()
</script>

<style scoped>
.opacity{
    opacity: .5;
    width: 100%;
    height: 100%;
    position: absolute;
    background: #f0f8ff6b;
}
.bg{
    width: 100%;
    height: 100%;
    position: absolute;
    z-index: -1;
}
.bg img{
    width: 100%;
    height: 100%;
    max-width: 100% !important;
    max-height: 100% !important;
}
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
    display: flex;
    justify-content: center;
    align-items: center;
}
.block img{
    max-width: 80%;
    max-height: 80%;
}
.block span{
    position: absolute;
    right: 0px;
    top: -10px;
    font-size: 1.7rem;
    font-weight: bold;
    -webkit-background-clip: text;
    /*color: #f8d94d;*/
    color: white;
    background-image: linear-gradient(to top, #ffb125 20%, #f7fc00 50%, #e99707 100%);
    -webkit-text-stroke: .5px #242622;
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
