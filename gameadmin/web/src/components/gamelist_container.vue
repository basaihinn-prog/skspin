<template>
    <el-form-item :label="haseManufacturer ? $t('游戏厂商') : ''" v-if="haseManufacturer">
        <el-space>

            <el-select v-model="manufacturer" @clear="clearManufacturer" popper-class="manufacturer_container" @change="changManufacturerList" clearable style="width: 150px">
                <el-option
                    v-for="item in ManufacturerList"
                    :key="item.ManufacturerCode"
                    :label="$t(item.ManufacturerName)"
                    :value="item.ManufacturerCode"
             >

                <span :class="item.ManufacturerCode != 'ALL' ? 'manufacturer_options' : 'manufacturer_options_all'"> {{ $t(item.ManufacturerName) }} </span>
                </el-option>
            </el-select>
        </el-space>
    </el-form-item>
    <el-form-item :label="haseManufacturer ? $t('游戏名称') : ''" >
        <el-space>


            <el-select v-model="manufacturer" @clear="clearManufacturer" popper-class="manufacturer_container" v-if="!haseManufacturer"  @change="changManufacturerList" clearable  style="width: 150px">
                <el-option
                    v-for="item in ManufacturerList"
                    :key="item.ManufacturerCode"
                    :label="$t(item.ManufacturerName)"
                    :value="item.ManufacturerCode"
                >

                    <span :class="item.ManufacturerCode != 'ALL' ? 'manufacturer_options' : 'manufacturer_options_all'"> {{ $t(item.ManufacturerName) }} </span>
                </el-option>
            </el-select>


            <el-select
                v-model="gameData"
                @change="selectGameList"
                clearable
                filterable
                @clear="clearGame"
                popper-class="game_popper"
                :disabled="gameLoading"
                :placeholder="$t('请输入')"
                style="width: 150px"
            >
                <el-option
                    v-for="item in options"
                    :key="item.value"
                    :label="$t(item.label)"
                    :value="item.value"
                />
            </el-select>
            <el-select v-if="showLotteryGame" v-model="lotteryGameType" @change="selectLottery" :disabled="gameLoading"  style="width: 150px">
                <el-option
                    v-for="item in Lottery_Peroid"
                    :key="item.value"
                    :label="$t(item.label)"
                    :value="item.value"
                />
            </el-select>
        </el-space>
    </el-form-item>
</template>

<script setup lang="ts">
import {onMounted, ref, reactive, defineProps, defineEmits, watch, watchEffect} from 'vue';
import { useStore } from '@/pinia/index';
import { useI18n } from 'vue-i18n';
import {Client} from "@/lib/client";
import {AdminInfo} from "@/api/adminpb/info";
import {tip} from "@/lib/tip";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {GameStatus} from "@/api/gamepb/customer";
import {Manufacturer} from "@/api/gamepb/game";
import {storeToRefs} from "pinia";
const { t } = useI18n();
const store = useStore();
const { defaultGameEvent, isInit, haseManufacturer, haseAll } = defineProps(['defaultGameEvent', "isInit", "haseManufacturer", "haseAll"]);
const emit  = defineEmits();
const { AdminInfo } = storeToRefs(store)
let gameData = ref()
let gameLoading = ref()
let options = ref([])
let ManufacturerList = ref([])
let manufacturer = ref()
let Lottery_Peroid =[
    {
        value: 0,
        label: t('一期')
    },{
        value: 1,
        label: t('二期')
    },
]
const StatusMap: JsMap<GameStatus, string> = {
    [GameStatus.Open]: t('正常'),
    [GameStatus.Maintenance]: t('维护'),
    [GameStatus.Hide]: t('隐藏'),
    [GameStatus.Close]: t('关闭'),
};

let showLotteryGame = ref(false)
let lotteryGameType = ref(null)
let UserName = ref('')

watch(gameData,()=>{
    selectGameList()
})

onMounted(() => {
    selectSearch()
    getManufacturerList()
})
let DefaultManufacturerOn = AdminInfo.value.Businesses.DefaultManufacturerOn
const getManufacturerList = async () => {
    const [ response, err ] = await Client.Do(Manufacturer.GetManufacturerList, {} as any)

    let data = response.List


    if (err){
        ManufacturerList.value = []
    }else{
        data.unshift({
            Id:"",
            ManufacturerName:"全部",
            ManufacturerCode:"ALL",
        })
        manufacturer.value = "ALL"
    }



    ManufacturerList.value = data


}

const clearGame = () => {

    if (haseAll){

        gameData.value = "ALL"
    }else{

        gameData.value = options.value[0].value
    }


    let resultId = gameData.value == "ALL" ? null : gameData.value
    let manufacturers = manufacturer.value == "ALL" ? null : manufacturer.value

    emit('update:lotteryGameType', lotteryGameType);
    emit('select-operator', {gameData: resultId, manufacturer: manufacturers},lotteryGameType);


}
const clearManufacturer = () => {


    manufacturer.value = haseAll ? "ALL" :  ManufacturerList.value[0].ManufacturerCode
}

const changManufacturerList = () => {

    selectSearch()
    selectGameList()
}

const selectGameList = () => {
    let id = gameData

    if (gameData){
        if (gameData.value !== 'lottery') {
            lotteryGameType.value = null
            showLotteryGame.value = false
        } else {
            lotteryGameType.value = 0
            showLotteryGame.value = true
        }


        if (typeof gameData != "string"){
            id = gameData.value
        }
    }

    let resultId = id == "ALL" ? null : id
    let manufacturers = manufacturer.value == "ALL" ? null : manufacturer.value


    // 触发查询操作，可以通过 emit 发送事件到父组件
    emit('update:gameData', id);
    emit('update:lotteryGameType', lotteryGameType);
    emit('select-operator', {gameData: resultId, manufacturer: manufacturers},lotteryGameType);
};
const selectLottery = () => {
    emit('update:lotteryGameType', lotteryGameType);
    emit('select-operator',gameData,lotteryGameType);
}
let selectSearch = async () => {

    try {
        gameLoading.value = true


        let data = {
            List: []
        }

        let err = ""

        let [response, error] = await Client.Do(AdminGameCenter.GameListOperator, {Maintenance: manufacturer.value == "ALL" ? '' : manufacturer.value})
        data.List = response.List
        err = error


        gameLoading.value = false

        if (err) {
            return tip.e(err)
        }
        if (haseAll){

            data.List.unshift({
                ID:"ALL",
                Name:"全部",
                Status:"0",
            })

            gameData.value = "ALL"
        }

        options.value = data.List
            .map(item => ({
                value: item.ID,
                label: item.Name,
                Status: item.Status,
            }));



        if (isInit && !haseAll){

            gameData.value = options.value[0].value
        }


    } catch (error) {
        tip.e(error.message || error);
        console.error('Failed to fetch game list:', error);
    }

}
</script>
<style scoped lang='scss'>
:deep(.el-select-v2__placeholder){
    font-size: .8rem;
}
.searchView .el-space{
    margin-bottom: 0 !important;
    margin-right: 15px!important;
}
</style>
