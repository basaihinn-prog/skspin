
<template>
    <App></App>
</template>

<script setup lang="tsx">

import { VNode } from "vue"
import { dyImport } from "./allimage";
import { DocBetLog } from "@/api/gamepb/customer";
import { useCssModule } from 'vue'

interface Props {
    Pan: number[][]// 盘面
    images: JsMap<string, string>  // ID索引图片
    highlight: (x, y, idx) => boolean
}

const props = withDefaults(defineProps<Props>(), {})

const styles = useCssModule()

let App = () => {
    let list: VNode[] = []
    let n = 0

    for (let i = 0; i < props.Pan.length; i++) {
        let rowData = props.Pan[i]
        let imgs: VNode[] = []
        for (let j = 0; j < rowData.length; j++) {
            // let src = dyImport(props.GameID +'/'+ rowData[i] + '.png')
            let src = props.images[rowData[j]]
            if (src) {
                imgs.push(<img class={[styles.block, props.highlight(i, j, n) ? styles.highlight : ""]} src={src}></img>)
            } else {
                imgs.push(<div class={styles.block}></div>)
            }
            n++
        }
        list.push(<div class={[styles.row, "flex_hbox"]}>{imgs}</div>)
    }
    return <div class="flex_vbox flex_child_center">
        {list}
    </div>
}

</script>

<style lang="less" module>
.block {
    max-width: 60px;
    max-height: 60px;
    min-width: 60px;
    min-height: 60px;
    width: 60px;
    height: 60px;
    border: 3px solid transparent;
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

// https://juejin.cn/post/6844904000307855374
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
