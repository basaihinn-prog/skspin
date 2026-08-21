<template>
    <div class="boxcolor flex_hbox">
        <div class="pan flex_grow">
            <div class="title">
                <el-popover placement="right" :width="400" trigger="hover">
                    <template #reference>
                            <span class="table-icon">
                                <el-icon><InfoFilled /></el-icon>
                            </span>
                    </template>
                    <div class="info">
                        <img :src="dyImport(`${props.DiaLogGame}/info_page.png`)" />
                    </div>
                </el-popover>
            </div>
            <template v-if="data.Turns">
                <el-dropdown  @command="handleCommand" class="text-center" :max-height="300">
                    <span class="el-dropdown-link">
                      旋转定位
                      <el-icon class="el-icon--right">
                        <arrow-down />
                      </el-icon>
                    </span>
                    <template #dropdown>
                        <el-dropdown-menu style="width: 300px;" >
                            <el-dropdown-item v-for="(item,index) in data.Turns" :command="index" class="dropdown-item-style">
                                <span>N0.{{ index+1 }}</span>
                                <span style="float: right">THB:{{ item.G/10000 }}</span>
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
                <span class="text-center">当前位置: N0.{{data.myCarouselSelectCommand}}</span>
            </template>
            <el-carousel :autoplay="false" max-height="900px" trigger="click" v-if="data.Turns" ref="myCarousel" indicator-position="none" @change="carouselChange">
                <el-carousel-item :key="'p'+panid" v-for="(pan,panid) in data.Turns">
                    <div class="flex_vbox">
                        <SlotPan :images="data.images" :Pan="pan.Pan" :highlight="(x,y) => hightLight(panid, x, y)" />
                        <div style="height:20px" />
                        <div class="flex_vbox">
                            <div v-if="pan.free">免费游戏</div>
                            <template v-if="pan.HitLines && pan.HitLines.length>0">
                                <!--                                <div>【中奖元素】总共{{ pan.HitLines.length }}个元素</div>-->
                                <p class="text-center">中奖情况</p>
                                <div class="text-center">
                                    <div :key="'line' + i" class="line" v-for="(p2, i) in pan.HitLines">
                                        <el-space wrap>
                                            <img :src="data.images[p2.Color]" class="block_img" />
                                            <div>
                                                {{ $t('金额')}}: <span class="color">{{ p2.G/10000 }}</span>
                                            </div>
                                            <div>
                                                <el-popover placement="right" :width="400" trigger="hover">
                                                    <template #reference>
                                                        {{ $t('详情') }}
                                                        <span class="table-icon">
                                                                <el-icon><InfoFilled /></el-icon>
                                                            </span>
                                                    </template>
                                                    <div class="info">
                                                        <div>
                                                            <p>
                                                                {{ $t('金额')}}
                                                                {{ p2.Formula }}
                                                            </p>
                                                            <p>
                                                                【详情】
                                                                各列数量：{{ p2.Counts?.join("*") }}，{{ p2.Counts?.length }}列({{ p2.Rate }}分)，倍率({{
                                                                    pan.Multi
                                                                }})。
                                                            </p>
                                                        </div>
                                                    </div>
                                                </el-popover>

                                            </div>
                                        </el-space>
                                    </div>
                                </div>
                            </template>
                        </div>
                    </div>
                </el-carousel-item>
            </el-carousel>
        </div>
    </div>
</template>
<script lang="ts" setup>

import SlotPan from "./SlotPan.vue";
import { onMounted, onUpdated, reactive, ref } from "vue";
import { dyImport, getImages } from "./allimage";
interface Props {
    PlayResp: any  // 盘面
    DiaLogGame: any  // 盘面
    Bet: any  // 投注金额
}
const props = withDefaults(defineProps<Props>(), {})
let data = reactive({
    Turns: null,
    images: getImages(props.DiaLogGame),
    myCarouselSelectCommand: 1,
})
const myCarousel = ref(null)
const handleCommand = (command: number) => {
    myCarousel.value.setActiveItem(command);
    data.myCarouselSelectCommand = command+1
}
let carouselChange = (index) => {
    console.log(index,'carouselChange');
    data.myCarouselSelectCommand = index+1
}
let hightLight = (pan, x,y)=> {
    return false
}
let updateData = () => {
    let turns = []
    let isStartGame = false

    if (props.PlayResp && props.PlayResp.DropPan) {
        props.PlayResp.DropPan.forEach(v => {
            if (isStartGame) {
                v.free = true
            }

            turns.push(v)

            if (v.IsStartGame) {
                isStartGame = true
            }
        })
    }
    data.Turns = turns
}
// onUpdated(updateData)
onMounted(updateData)
</script>
<style lang="less" scoped>
.pan {
    padding: .5rem;
    background: #30303c;
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
}
.color{
    color: rgb(255, 200, 36);
}
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
    height: 600px !important;
}
.boxcolor {
    background-color: #30303c;
}
.text-center{
    text-align: center;
    color: hsla(0, 0%, 100%, 0.6);
}
</style>
