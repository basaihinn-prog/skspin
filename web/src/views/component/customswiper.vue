
<template>
    <div class="hidden-sm-and-down" style="width: 100%;position: relative">
        <swiper
                @realIndexChange="slideChange"
                @swiper="setSwiperRef"
                :effect="'coverflow'"
                :centeredSlides="true"
                :initialSlide="1"
                :autoplay="true"
                :loop="true"
                :slidesPerView="1.4"
                :coverflowEffect="{rotate: 5,stretch: 600, depth: 105,modifier: 1,slideShadows: false,}"
                :modules="modules"
                :pagination="true"
                class="mySwiper">
            <swiper-slide>

                <el-image :src="swiper1" fit="fill">
                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>


            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper2" fit="fill">

                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>


            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper3" fit="fill">

                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>
            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper4" fit="fill">
                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>

            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper5" fit="fill">

                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>
            </swiper-slide>
        </swiper>

    </div>

    <div class="hidden-md-and-up" style="height: 65px"></div>
    <div class="hidden-md-and-up" style="width: 100%;position: relative">
        <swiper
                @realIndexChange="slideChange"
                @swiper="setSwiperRef"
                :effect="'coverflow'"
                :centeredSlides="true"
                :initialSlide="1"
                :autoplay="true"
                :loop="true"
                :slidesPerView="1.1"
                :coverflowEffect="{rotate: 5,stretch: 10, depth: 105,modifier: 1,slideShadows: false}"
                :modules="modules"
                :pagination="true"
                class="mySwiper">
            <swiper-slide>

                <el-image :src="swiper1" fit="fill">
                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>

            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper2" fit="fill">
                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>

            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper3" fit="fill">

                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>
            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper4" fit="fill">
                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>

            </swiper-slide>
            <swiper-slide>

                <el-image :src="swiper5" fit="fill">
                    <template #error>
                        <div class="image-slot">

                            <el-icon>
                                <Picture/>
                            </el-icon>

                        </div>
                    </template>
                    <template #placeholder>
                        <div class="placeholder-slot">

                        </div>
                    </template>
                </el-image>

            </swiper-slide>
        </swiper>

    </div>

</template>

<script setup lang="ts">

import {Swiper, SwiperSlide, useSwiper} from 'swiper/vue';

import {EffectCoverflow, Pagination} from "swiper/modules";
import swiper1 from "@/assets/swiper1.webp";
import swiper2 from "@/assets/swiper2.webp";
import swiper3 from "@/assets/swiper3.webp";
import swiper4 from "@/assets/swiper4.webp";
import swiper5 from "@/assets/swiper5.webp";
import {ref} from "vue";

const uSwiper = useSwiper();

import 'swiper/css';
import 'swiper/css/effect-coverflow';
import 'swiper/css/pagination';
const modules = [EffectCoverflow, Pagination]
const swiperHeight = ref(null)
const swiperRef = ref(null)
const indexNum = ref(0)


const slideChange = (swiper) => {

    indexNum.value = swiper.realIndex
}
const setSwiperRef = (swiper) => {
    swiperHeight.value = swiper.height
    swiperRef.value = swiper
}
</script>

<style scoped>

</style>
