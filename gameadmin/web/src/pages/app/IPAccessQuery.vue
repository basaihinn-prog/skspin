<template>
    <div>
        <div class="searchView gameList" >
            <el-form
                label-position="top"
                label-width="100px"
                style="max-width: 100%"
                @keyup.enter="initSearch"
            >
                <el-space wrap>

                    <operator_container :defaultOperatorEvent="defaultOperatorEvent"
                                        @select-operator="operatorListChange"></operator_container>
                    <gamelist_container :defaultGameEvent="defaultGameEvent"
                                        @select-operator="selectGameList"></gamelist_container>
                    <el-form-item :label="$t('访问时间')">
                        <el-date-picker
                            v-model="datePicker" value-format="x"
                            locale="zh-cn" type="datetimerange" format="YYYY-MM-DD HH:mm:ss"
                        />
                    </el-form-item>
                </el-space>
            </el-form>
            <el-space wrap>
                <el-button type="primary" @click="initSearch">{{ $t('搜索') }}</el-button>
            </el-space>
        </div>
        <el-table :data="tableData" style="width: 100%">
            <el-table-column :label="$t('访问时间')" prop="LoginTime" :formatter="dateSecondFormater"></el-table-column>
            <el-table-column :label="$t('商户')" prop="Pid"></el-table-column>
            <el-table-column :label="$t('游戏名称')" prop="GameID" >
                <template #default="scope">
                    <!--                    {{uiData.GameList.find(t=>t.ID === scope.row.GameID)?.Name}}-->
                    {{ findGameName(scope.row.GameID?.startsWith('lottery') ? 'lottery' : scope.row.GameID) }}
                </template>
            </el-table-column>
            <el-table-column :label="$t('IP')" prop="Ip" ></el-table-column>
            <el-table-column :label="$t('IP地址')" prop="Loc" ></el-table-column>
        </el-table>
        <Pagination :Count='searchList.Count' @paginationEmit='paginationEmit'></Pagination>
    </div>
</template>

<script lang="ts">
export default {
    name: "IPAccessQuery"
}
</script>
<script setup lang="ts">
import {onMounted, ref, reactive, nextTick} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {AdminStatsRpc} from '@/api/stats/stats';
import {AdminGameCenter} from "@/api/gamepb/admin";
import Operator_container from "@/components/operator_container.vue";
import {useI18n} from 'vue-i18n';
import Gamelist_container from "@/components/gamelist_container.vue";
const {t} = useI18n()
import ut from '@/lib/util'
const defaultOperatorEvent = ref({})
const defaultGameEvent = ref({})
const datePicker = ref([])
let searchList = reactive({
    OperatorId:null,
    GameID:'',
    StartTime:null,
    EndTime:null,
    PageSize:null,
    PageNumber:null,
    Count:0,
})
const tableData = ref([])
const operatorListChange = (value) => {
    searchList.OperatorId = value.value
}
const selectGameList = (value) =>{
    searchList.GameID = value
}

const initSearch = async () => {
    console.log(datePicker.value);
    if (datePicker.value && datePicker.value.length) {
        searchList.StartTime = datePicker.value[0]/1000
        searchList.EndTime = datePicker.value[1]/1000
    } else {
        searchList.StartTime = null
        searchList.EndTime = null
    }
    let [data, err] = await Client.Do(AdminStatsRpc.GetGamesLoginList, searchList)
    if (err) {
        return tip.e(err)
    }
    searchList.Count = data.Count
    searchList.PageNumber = data.PageNumber
    searchList.PageSize = data.PageSize
    tableData.value = data.Count === 0 ? [] : data.List
}

const paginationEmit = (PageIndex, PageSize) => {
    searchList.PageNumber = PageIndex
    searchList.PageSize = PageSize
    initSearch()
}

let initGameList = ref(null)
let GameListSearch = async () => {
    let [datas, errs] = await Client.Do(AdminGameCenter.GameList, {})
    if (errs) {
        return tip.e(errs)
    }
    initGameList.value = datas.List
}
const findGameName = (gameID) => {
    let game = initGameList.value?.find(t => t.ID === gameID);
    console.log(game);
    return game ? game.Name : ''
}
onMounted(() => {
    initSearch()
    GameListSearch()
});
</script>
<style scoped>

</style>
