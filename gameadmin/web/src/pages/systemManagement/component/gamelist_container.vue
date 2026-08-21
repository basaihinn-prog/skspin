<template>

    <el-form-item :label="haseManufacturer ? $t('游戏厂商') : ''" v-if="haseManufacturer">



            <el-select v-model="manufacturer" @clear="clearManufacturer" popper-class="manufacturer_container" @change="changManufacturerList" style="width: 300px">

                <template v-for="item in ManufacturerList">
                    <el-option
                            v-if="!visibleManufacturer || visibleManufacturer.indexOf(item.ManufacturerName) != -1"
                            :label="$t(item.ManufacturerName)"
                            :value="item.ManufacturerCode"
                    >

                        <span :class="item.ManufacturerCode != 'ALL' ? 'manufacturer_options' : 'manufacturer_options_all'"> {{
                            $t(item.ManufacturerName)
                            }} </span>
                    </el-option>
                </template>

            </el-select>

    </el-form-item>
    <el-form-item :label="haseManufacturer ? $t('游戏名称') : ''">

        <el-select v-model="manufacturer" @clear="clearManufacturer" popper-class="manufacturer_container"
                   v-if="!haseManufacturer" @change="changManufacturerList" style="width: 150px">

            <template v-for="item in ManufacturerList">
                <el-option
                        v-if="!visibleManufacturer || visibleManufacturer.indexOf(item.ManufacturerName) != -1"
                        :label="$t(item.ManufacturerName)"
                        :value="item.ManufacturerCode"
                >

                    <span :class="item.ManufacturerCode != 'ALL' ? 'manufacturer_options' : 'manufacturer_options_all'"> {{
                        $t(item.ManufacturerName)
                        }} </span>
                </el-option>
            </template>
        </el-select>


        <el-select
                v-model="gameData"
                @change="selectGameList"
                clearable
                filterable
                @clear="clearGame"
                popper-class="game_popper"
                multiple
                collapse-tags
                :disabled="gameLoading"
                :placeholder="$t('请输入')"
                style="width: 300px"
            >
                <el-option
                    v-for="item in options"
                    :key="item.value"
                    :label="$t(item.label)"
                    :value="item.value"
                />
            </el-select>
            <el-select v-if="showLotteryGame" v-model="lotteryGameType" @change="selectLottery" multiple :disabled="gameLoading"  style="width: 300px">
                <el-option
                    v-for="item in Lottery_Peroid"
                    :key="item.value"
                    :label="$t(item.label)"
                    :value="item.value"
                />
            </el-select>
            <div style="margin-top: 10px">
                <el-tag type="primary" v-for="item in gameList"
                        :key="item.label"
                        @close="closeTag(item.value)"
                        closable style="margin-right: 10px">{{ $t(item.label) }}</el-tag>
            </div>

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
const { t } = useI18n();
const store = useStore();
const { defaultGameEvent, isInit, haseManufacturer, haseAll,visibleManufacturer } = defineProps(['defaultGameEvent', "isInit", "haseManufacturer", "haseAll","visibleManufacturer"]);
const emit  = defineEmits();
let gameData = ref()
let gameList = ref()
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

watch(manufacturer,()=>{
    selectGameList([])

    selectSearch()
})

onMounted(() => {
    selectSearch()
    getManufacturerList()

})

const getManufacturerList = async () => {
    const [ response, err ] = await Client.Do(Manufacturer.GetManufacturerList, {} as any)

    if (err){
        ManufacturerList.value = []
        return
    }

    response.List.unshift({
        Id:"ALL",
        ManufacturerName:"全部",
        ManufacturerCode:"ALL",
    })


    ManufacturerList.value = response.List
    manufacturer.value = response.List[0].ManufacturerCode

}

const clearGame = () => {

    selectGameList([])
}
const clearManufacturer = () => {


    manufacturer.value = haseAll ? "ALL" :  ManufacturerList.value[0].ManufacturerCode
}

const changManufacturerList = () => {

    selectSearch()
    selectGameList([])
}

const closeTag = (value) => {
    gameData.value = gameData.value.filter(item=> item != value)

    selectGameList(gameData.value)
}

const selectGameList = (value?) => {



    if(value[0] != 'ALL' && value.indexOf("ALL") != -1){
        gameData.value = gameData.value.filter(item=> item == "ALL")
    }else if(value.length == 0){
        gameData.value = ["ALL"]
    }else{
        gameData.value = gameData.value.filter(item=> item != "ALL")
    }


    if (gameData.value.length >= 1){

        gameList.value = options.value.filter(item=> gameData.value.indexOf(item.value) > -1)
    }else{
        gameList.value = options.value.filter(item=> item.value == "ALL")

    }


        // 触发查询操作，可以通过 emit 发送事件到父组件
    emit('update:gameData', gameData);
    emit('update:lotteryGameType', lotteryGameType);
    emit('select-operator', {gameData: gameData, manufacturer: manufacturer.value},lotteryGameType);
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
        data.List = response.List.filter(item=> item.ChangeBetOff == 1 || (item.ChangeBetOff == 0 && item.ManufacturerName == 'PP'))
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

            gameData.value = ["ALL"]
        }

        options.value = data.List
            .map(item => ({
                value: item.ID,
                label: item.Name,
                Status: item.Status,
            }));
        gameList.value = options.value.filter(item=> item.value == "ALL")

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
