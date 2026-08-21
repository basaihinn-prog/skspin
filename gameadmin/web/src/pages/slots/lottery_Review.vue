<template>
    <div class="boxcolor flex_hbox">
        <div class="pan flex_grow">
            <el-descriptions>
                <el-descriptions-item :label="$t('时间')+'：'">{{ data.date }}</el-descriptions-item>
                <el-descriptions-item>
                    <el-row :gutter="10">
                        <template v-for="isItem in data.list">
                            <template v-if="isItem.GuessNumber">
                                <el-col>
                                    <p class="mb-2">
                                        <el-tag type="success" effect="light">
                                            {{strArrTitle[(guessFirstPrizeType[isItem.GuessType].type-1)] }}
                                        </el-tag>
                                    </p>
                                    <p>
                                        <el-space wrap>
                                            <span>{{$t('类型')}}:</span>
                                            <el-tag type="info" effect="plain">{{guessFirstPrizeType[isItem.GuessType].name}}</el-tag>
                                        </el-space>
                                    </p>
                                    <p>
                                        <el-space wrap><span>{{$t('描述')}}:</span>
                                            {{isItem.GuessTypeDesc}}
                                        </el-space>
                                    </p>
                                    <p>
                                        <el-space wrap>
                                            <span>{{$t('彩票')}}:</span>
                                            <span>{{isItem.GuessNumber}}</span>
                                        </el-space>
                                    </p>
                                    <p>
                                        <el-space wrap>
                                            <span>{{$t('下注金额')}}:</span>
                                            <span>{{ut.toNumberWithCommaNormal(isItem.GuessGold)}}</span>
                                        </el-space>
                                    </p>
                                </el-col>

                            </template>
                            <template v-else>
                                <el-col :span="6">
                                    <div style="margin-bottom: 1rem">
                                        <p>
                                            <el-space wrap>
                                                <span>{{$t('彩票')}}:</span>
                                                <el-tag type="info" effect="plain">
                                                    {{isItem.LotterNumber}}
                                                </el-tag>
                                            </el-space>
                                        </p>
                                        <p>
                                            <el-space wrap>
                                                <span>{{$t('数量')}}:</span>
                                                <span>{{isItem.Count}}</span>
                                            </el-space>
                                        </p>
                                    </div>
                                </el-col>

                            </template>

                        </template>
                    </el-row>



                </el-descriptions-item>
            </el-descriptions>
        </div>

    </div>
</template>

<script lang="ts" setup>

import SlotPan from "./SlotPan_Olympus.vue";
import { onMounted, onUpdated, reactive, ref } from "vue";
import { dyImport, getImages } from "./allimage";
import {useI18n} from "vue-i18n";
import ut, {initGolbal} from '@/lib/util'
interface Props {
    PlayResp: any  // 盘面
    DiaLogGame: any  // 盘面
}
const props = withDefaults(defineProps<Props>(), {})
let dateJson = JSON.parse(props.PlayResp)
const {t} = useI18n()
console.log(dateJson);
let data = reactive({
    date: dateJson.OpenDay,
    list: dateJson.Numbers,
})
let guessFirstPrizeType = reactive([
    {
        name: t('后三位'),
        len: 3,
        start: 0,
        end: 2,
        type: 1,
        sliceStart: -3,
        sliceEnd: undefined,
        multiple: 1000,
        disordered: true,
        alreadyEnteredList: null
    },
    {
        name: t('后二位'),
        len: 2,
        start: 4,
        end: 5,
        type: 1,
        sliceStart: -2,
        sliceEnd: undefined,
        multiple: 100,
        disordered: true,
        alreadyEnteredList: null
    },
    {
        name: t('后一位'),
        len: 1,
        start: 5,
        end: undefined,
        type: 1,
        sliceStart: -1,
        sliceEnd: undefined,
        multiple: 10,
        disordered: true,
        alreadyEnteredList: null
    },
    {
        name: t('后三位无序'),
        len: 3,
        start: 3,
        end: 5,
        type: 1,
        sliceStart: -3,
        sliceEnd: undefined,
        multiple: 1000,
        disordered: false, // 无序
        alreadyEnteredList: null
    },
    // ----
    {
        name: t('前三位'),
        len: 3,
        start: 0,
        end: 2,
        type: 2,
        sliceStart: 0,
        sliceEnd: 3,
        multiple: 450,
        disordered: true,
        alreadyEnteredList: null
    },
    {
        name: t('后三位'),
        len: 3,
        start: 3,
        end: 5,
        type: 2,
        sliceStart: -3,
        sliceEnd: undefined,
        multiple: 450,
        disordered: true,
        alreadyEnteredList: null
    },
    {
        name: t('后二位'),
        len: 2,
        start: 4,
        end: 5,
        type: 3,
        sliceStart: -2,
        sliceEnd: undefined,
        multiple: 450,
        disordered: true,
        alreadyEnteredList: null
    },
    {
        name: t('后一位'),
        len: 1,
        start: 5,
        end: undefined,
        type: 3,
        sliceStart: -1,
        sliceEnd: undefined,
        multiple: 10,
        disordered: true,
        alreadyEnteredList: null
    },
])
const strArrTitle = [t('猜一等奖号码'), t('猜前三位'), t('猜两位数')]
// onUpdated(updateData)
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
.text-center{
    text-align: center;
    color: hsla(0, 0%, 100%, 0.6);
}
:deep(.el-descriptions__body){
    background: transparent;
}
:deep(.el-descriptions__header){
    margin: 0;
}
:deep(.el-descriptions__label:not(.is-bordered-label)){
    color: #606266;
}
.mb-2{
    margin-bottom: 1rem;
    :deep(.el-tag__content) {
        font-size: .9rem;
    }
}

</style>
