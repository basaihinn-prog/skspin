
<template>
    <div class="flex_vbox flex_child_center">
        <div v-for="(rowData, rowIndex) in Pan.Data" :key="'row-' + rowIndex" class="flex_hbox row">
            <div v-for="(imgId, colIndex) in rowData" :key="'col-' + colIndex" class="block" :class="{ 'highlight': highlight(rowIndex, colIndex, colIndex) }">
                <template v-if="typeof imgId === 'object'">
                    <template v-if="images[imgId.nums]">
                        <img :src="images[imgId.nums]" />
                        <span>{{imgId.type}}x</span>
                        <!--                        // [0-5) 黄 // [5-50) 红色 // [50-无穷) 粉色-->
                    </template>
                    <template v-else>
                        <div class="block"></div>
                    </template>
                </template>
                <template v-else>
                    <img v-if="images[imgId]" :src="images[imgId]" />
                    <div v-else class="block"></div>
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

/*interface Props {
    Pan: number[][]// 盘面
    images: JsMap<string, string>  // ID索引图片
    highlight: (x, y, idx) => boolean
}*/

const props = defineProps(['Di','Pan','images','highlight'])
const styles = useCssModule()
const newPanListData = props.Pan.Data
const newPanListType = props.Pan.Type
const init = () => {
    console.log(props.Pan,'newPanListData.length');
    console.log(newPanListData,'newPanListData.length');
    for (let i = 0; i < newPanListData.length; i++) {
        if (newPanListData[i].length === newPanListType[i].length) {
            for (let j = 0; j < newPanListData[i].length; j++) {
                if (newPanListData[i][j] > 0 && newPanListType[i][j] > 0) {
                    newPanListData[i][j] = {
                        nums: newPanListData[i][j],
                        type: newPanListType[i][j],
                    };
                }
            }
        }
    }
    console.log(newPanListData,'newPanListData');
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
.block img{
    width: 100%;
    height: auto;
}
.block span{
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
