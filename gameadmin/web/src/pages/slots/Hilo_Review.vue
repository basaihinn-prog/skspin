<template>
    <div class="boxcolor flex_hbox">
        <div class="pan flex_grow">
            <el-table :data="data.Turns" style="width: 100%" max-height="550">
                <el-table-column prop="Area" :label="$t('投注区域')" />
                <el-table-column prop="Bet" :label="$t('投注金额')" width="100">
                    <template #default="scope">
                        {{goldFormater(0,0,scope.row.Bet) }}
                    </template>
                </el-table-column>
                <el-table-column prop="Rate" :label="$t('中奖倍率')" width="100"/>
            </el-table>
            <p>{{ $t("中奖结果") }}：
                <template v-for="item in props.PlayResp.T">
                    <img :src="data.images[item]" class="block_img" />
                </template>
            </p>
            <p>{{ $t('玩家投注') }}：{{goldFormater(0,0,props.PlayResp.BetTotal)}}</p>
            <p>{{ $t('玩家赢分') }}：{{goldFormater(0,0,props.PlayResp.WinTotal)}}</p>
        </div>

    </div>
</template>

<script lang="ts" setup>
// 1号区域 1号区域点数相加小
// 2号区域 2号区域点数相加大
// 3号区域 3号区域点数相加11
// 4号区域 4号区域骰子1点
// 5号区域 5号区域骰子2点
// 6号区域 6号区域骰子3点
// 7号区域 7号区域骰子4点
// 8号区域 8号区域骰子5点
// 9号区域 9号区域骰子6点
// 10号区域 10号区域点数相加小且骰子包含5点
// 11号区域 11号区域点数相加小且骰子包含6点
// 12号区域 12号区域骰子包含1点和2点
// 13号区域 13号区域骰子包含1点和4点
// 14号区域 14号区域骰子包含1点和5点
// 15号区域 15号区域骰子包含1点和6点
// 16号区域 16号区域骰子包含2点和6点
// 17号区域 17号区域骰子包含2点和4点
// 18号区域 18号区域骰子包含2点和5点
// 19号区域 19号区域骰子包含2点和3点
// 20号区域 20号区域骰子包含3点和5点
// 21号区域 21号区域骰子包含3点和6点
// 22号区域 22号区域骰子包含3点和4点
// 23号区域 23号区域骰子包含4点和5点
// 24号区域 24号区域骰子包含5点和6点
// 25号区域 25号区域点数相加小且骰子包含1点
// 26号区域 26号区域点数相加小且骰子包含3点
// 27号区域 27号区域点数相加大且骰子包含4点
// 28号区域 28号区域点数相加大且骰子包含6点
// 29号区域 29号区域骰子包含1、2、3点中的2个
import SlotPan from "./SlotPan.vue";
import { onMounted, onUpdated, reactive, ref } from "vue";
import { dyImport, getImages } from "./allimage";
interface Props {
    PlayResp: any  // 盘面
    DiaLogGame: any  // 盘面
}
const list = [
    '',
    '1号区域点数相加小',
    '2号区域点数相加大',
    '3号区域点数相加11',
    '4号区域骰子1点',
    '5号区域骰子2点',
    '6号区域骰子3点',
    '7号区域骰子4点',
    '8号区域骰子5点',
    '9号区域骰子6点',
    '10号区域点数相加小且骰子包含5点',
    '11号区域点数相加小且骰子包含6点',
    '12号区域骰子包含1点和2点',
    '13号区域骰子包含1点和4点',
    '14号区域骰子包含1点和5点',
    '15号区域骰子包含1点和6点',
    '16号区域骰子包含2点和6点',
    '17号区域骰子包含2点和4点',
    '18号区域骰子包含2点和5点',
    '19号区域骰子包含2点和3点',
    '20号区域骰子包含3点和5点',
    '21号区域骰子包含3点和6点',
    '22号区域骰子包含3点和4点',
    '23号区域骰子包含4点和5点',
    '24号区域骰子包含5点和6点',
    '25号区域点数相加小且骰子包含1点',
    '26号区域点数相加小且骰子包含3点',
    '27号区域点数相加大且骰子包含4点',
    '28号区域点数相加大且骰子包含6点',
    '29号区域骰子包含1',
    '30号区域骰子包含4、5、6点中的2个',
]
const props = withDefaults(defineProps<Props>(), {})
let data = reactive({
    Turns: null,
    images: getImages(props.DiaLogGame),
})

let hightLight = (pan, x, y) => {
    return false
}

let updateData = () => {
    let turns = []
    // let isStartGame = false
    console.log(props,'props');
    console.log(props.PlayResp.T,'props.PlayResp');
    if (props.PlayResp && props.PlayResp.WinInfos) {

        let arrs = props.PlayResp.WinInfos
        for (const turnsKey in arrs) {

            for (let a = 0;a<list.length;a++) {
                if (a === Number(turnsKey)) {
                    arrs[turnsKey].Area = list[a]
                    arrs[turnsKey].Rate = !arrs[turnsKey].Rate?'':'X'+arrs[turnsKey].Rate
                }
            }
            turns.push(arrs[turnsKey])
        }
        /*props.PlayResp.WinInfos.forEach(v => {
            /!*if (isStartGame) {
                v.free = true
            }*!/
            console.log(v,'---------v');
            turns.push(v)
            /!*if (v.IsStartGame) {
                isStartGame = true
            }*!/
        })*/
    }
    data.Turns = turns
    console.log(data.Turns,'data.Turns');
}
// onUpdated(updateData)
onMounted(updateData)
</script>

<style lang="less" scoped>
.pan {
    padding: .5rem;
    background: #30303c;
}
.color{
    color: rgb(255, 200, 36);
}
.text-center{
    text-align: center;
    color: hsla(0, 0%, 100%, 0.6);
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
    max-width: 400px;
    height: auto;

    img {
        max-width: 100%;
    }
}

.line {
    cursor: pointer;
}

.carousel_wh {
    height: 600px !important;
}

.boxcolor {
    flex-wrap: wrap;
}
</style>
