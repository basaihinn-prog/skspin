
<template>
    <div class="flex_vbox flex_child_center">
        <div v-for="(rowData, rowIndex) in Pan" :key="'row-' + rowIndex" class="flex_hbox row">
            <template v-if="rowIndex>0">
                <div v-for="(imgId, colIndex) in rowData" :key="'col-' + colIndex" class="block" :class="{ 'highlight': highlight(rowIndex, colIndex, colIndex) }"
                     :style="
                     ((rowIndex===1 && colIndex===0) ||
                     (rowIndex===2 && colIndex===0) ||
                     (rowIndex===3 && colIndex===0) ||
                     (rowIndex===4 && colIndex===0) ||
                     (rowIndex===5 && colIndex===0) ||
                     (rowIndex===1 && colIndex===4) ||
                     (rowIndex===2 && colIndex===4) ||
                     (rowIndex===3 && colIndex===4) ||
                     (rowIndex===4 && colIndex===4) ||
                     (rowIndex===5 && colIndex===4)
                     )
                     ?{'top':'2rem'}:''"
                >
                    <template v-if="(colIndex === 0 && rowIndex === 5) || (colIndex === 4 && rowIndex === 5)">
                        <div class="block"></div>
                    </template>
                    <template v-else>
                        <img :src="images[imgId]" class="block" />
                    </template>
                </div>
            </template>

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

const props = defineProps(['Pan','images','highlight'])
console.log(props,'props----');
const styles = useCssModule()
console.log(props.Pan,'props.Pan');


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
