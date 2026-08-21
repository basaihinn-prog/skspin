
<template>
    <div class="flex_vbox flex_child_center">
        <div v-for="(rowData, rowIndex) in updatedData" :key="'row-' + rowIndex" class="flex_hbox row">
            <div v-for="(imgId, colIndex) in rowData" :key="'col-' + colIndex" class="block">
                <template v-if="hitLines === imgId">
                    <div class="borders"></div>
                </template>
                <template v-if="typeof imgId === 'object'">
                    <template v-if="imgId.nums === 11">
                        <template v-if="imgId.findNum">
                            <img v-if="imgId.urls === 0" :src="img1"  alt=""/>
                            <img v-if="imgId.urls === 1" :src="img2"  alt=""/>
                            <img v-if="imgId.urls === 2" :src="img3"  alt=""/>
                            <img v-if="imgId.urls === 3" :src="img4"  alt=""/>
                            <span>{{imgId.findNum}}x</span>
                        </template>
                        <template v-else>
                            <img :src="images[imgId.nums]" />
                        </template>
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
import img1 from '@/assets/Starlight1000/Tumble0.png'
import img2 from '@/assets/Starlight1000/Tumble1.png'
import img3 from '@/assets/Starlight1000/Tumble2.png'
import img4 from '@/assets/Starlight1000/Tumble3.png'
/*interface Props {
    Pan: number[][]// 盘面
    images: JsMap<string, string>  // ID索引图片
    highlight: (x, y, idx) => boolean
}*/
const arrs = {
    0:[100,250,500],
    1:[25,50],
    2:[10,12,15,20],
    3:[2,3,4,5,6,8]
}
const props = defineProps(['Di','Pan','images','highlight','hitLines'])
const styles = useCssModule()
const newPanListData = props.Pan.Data
const newPanListType = props.Pan.Type
let updatedData = null
const init = () => {
    updatedData = newPanListData.map((dataRow, i) => {
        if (dataRow.length === newPanListType[i].length) {
            return dataRow.map((data, j) => {
                if (data > 0 && newPanListType[i][j] > 0) {
                    if (data === 11) {
                        for (let arrsKey in arrs) {
                            let nums = arrs[arrsKey].findIndex(u => u === newPanListType[i][j]);
                            if (nums !== -1) {
                                let findNum = arrs[arrsKey][nums];
                                return {
                                    nums: data,
                                    type: newPanListType[i][j],
                                    findNum: findNum,
                                    urls: Number(arrsKey)
                                };
                            }
                        }
                    } else {
                        return {
                            nums: data,
                            type: newPanListType[i][j],
                            findNum: null,
                            urls: null
                        };
                    }
                }
                return data;
            });
        }
        return dataRow;
    });
    return updatedData;
};


init()
</script>

<style scoped>
.flex {
    display: flex;
    display: -webkit-flex;
}
.borders{
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    border: 2px solid #f9fcff;
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
    height: 100%;
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
