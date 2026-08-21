<template>
    <div class="boxcolor">
        <div>
            <img class="info" :src="dyImport(`${props.DiaLogGame}/Info_Page_1.jpg`)" />
        </div>
        <el-carousel :autoplay="false" height="600px" trigger="click" v-if="data.Turns">
            <el-carousel-item :key="'p'+panid" v-for="(pan,panid) in data.Turns">
                <div class="flex_vbox">
                    <SlotPan :images="data.images" :Pan="pan.Pan" :highlight="(x,y, idx) => hightLight(panid, x,y,idx)" />
                    <div style="height:20px" />
                    <div class="flex_vbox">
                        <div v-if="pan.free">{{ $t("免费游戏") }}</div>
                        <div>{{ $t("中奖线") }}</div>
                        <div :key="'line'+i" class="line" v-for="(p2, i) in pan.HitLines">
                            <div
                            @click="setHightLight(panid, p2.Line.BaseLine)"
                            >{{ $t("元素") }}：{{p2.Line?.BaseLine}} {{ $t('金额') }}：{{p2.G/10000}}</div>
                        </div>
                    </div>
                    <div style="height:20px" />
                    <div class="flex_vbox" v-if="pan.LionSmallGame">
                        <div>{{ $t('砍狮子') }}</div>
                        <div
                            :key="'kan'+i"
                            class="line"
                            v-for="(kan, i) in pan.LionSmallGame.Resps"
                        >
                            <div class="flex_hbox">
                                <img :src="data.attack[kan.Weapon]" class="weapon" />
                                <div>砍狮子：{{goldFormater(0,0,kan.AttackGold)}} 狮子砍我：{{goldFormater(0,0,kan.LionAttackGold)}}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </el-carousel-item>
        </el-carousel>
    </div>
</template>

<style lang="less" scoped>
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
    height: 600px !important;
}
.boxcolor {
    background-color: #30303c;
}
</style>

<script lang="ts" setup>

import BigFight from "@/assets/slots/Roma/BigFight.png"
import SmallFight from "@/assets/slots/Roma/SmallFight.png"
import Defense from "@/assets/slots/Roma/Defense.png"

import SlotPan from "./SlotPan.vue";
import { onMounted, onUpdated, reactive } from "vue";
import { dyImport, getImages } from "./allimage";
import { watch } from "fs";
interface Props {
    PlayResp: any  // 盘面
    DiaLogGame: any  // 盘面
}
const props = withDefaults(defineProps<Props>(), {})
let data = reactive({
    Turns: null,
    images: getImages(props.DiaLogGame),
    attack: [SmallFight, BigFight, Defense],
    hightLightIdx: {},
})

const setHightLight = (pan, line) => {
    data.hightLightIdx[pan] = {}
    line.forEach((v) => data.hightLightIdx[pan][v] = true)
}
const hightLight = (pan, x, y, idx) => {
    if (!data.hightLightIdx[pan]) {
        return false
    }
    return data.hightLightIdx[pan][idx]
}
let updateData = () => {
    data.hightLightIdx = {}

    let turns = []

    props.PlayResp.DropPan.forEach(v => {
        turns.push(v)
    })
    props.PlayResp.FreeSmallGame?.Resps.forEach(v => {
        v.DropPan.forEach(v2 => {
            v2.free = true
            turns.push(v2)
        })
    })

    data.Turns = turns
}
// onUpdated(updateData)
onMounted(updateData)

</script>
