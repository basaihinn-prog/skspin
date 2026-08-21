<template>
    <table class="table">
        <thead>
        <tr>
            <th scope="col">{{$t('指标名称')}}</th>
            <th scope="col">{{$t('总数据')}}</th>
            <th scope="col">{{$t('七日数据')}}</th>
            <th scope="col">{{$t('本月数据')}}</th>
            <th scope="col">{{$t('购买小游戏总数据')}}</th>
            <th scope="col">{{$t('购买小游戏七日')}}</th>
            <th scope="col">{{$t('购买小游戏本月')}}</th>
            <template v-if="gameMsgProps">
                <th scope="col">{{$t('额外模式总数据')}}</th>
                <th scope="col">{{$t('额外模式七日')}}</th>
                <th scope="col">{{$t('额外模式本月')}}</th>
            </template>
        </tr>
        </thead>
        <tbody>
        <tr>
            <th scope="row">{{$t('流水')}}</th>
            <td>{{runtimeTable?.TotalData.BetAmount}}</td>
            <td>{{ runtimeTable?.SevenDays.BetAmount }}</td>
            <td>{{ runtimeTable?.MonthData.BetAmount }}</td>
            <td>{{ runtimeTable?.TotalData.BuyBetAmount }}</td>
            <td>{{ runtimeTable?.SevenDays.BuyBetAmount }}</td>
            <td>{{ runtimeTable?.MonthData.BuyBetAmount }}</td>
            <template v-if="gameMsgProps">
                <td>{{ runtimeTable?.TotalData.ExtraBetAmount }}</td>
                <td>{{ runtimeTable?.SevenDays.ExtraBetAmount }}</td>
                <td>{{ runtimeTable?.MonthData.ExtraBetAmount }}</td>
            </template>
        </tr>
        <tr>
            <th scope="row">{{$t('产出')}}</th>
            <td>{{ runtimeTable?.TotalData.WinAmount }}</td>
            <td>{{ runtimeTable?.SevenDays.WinAmount }}</td>
            <td>{{ runtimeTable?.MonthData.WinAmount }}</td>
            <td>{{ runtimeTable?.TotalData.BuyWinAmount }}</td>
            <td>{{ runtimeTable?.SevenDays.BuyWinAmount }}</td>
            <td>{{ runtimeTable?.MonthData.BuyWinAmount }}</td>
            <template v-if="gameMsgProps">
                <td>{{ runtimeTable?.TotalData.ExtraWinAmount }}</td>
                <td>{{ runtimeTable?.SevenDays.ExtraWinAmount }}</td>
                <td>{{ runtimeTable?.MonthData.ExtraWinAmount }}</td>
            </template>

        </tr>
        <tr>
            <th scope="row">{{$t('回报率')}}</th>
            <td>{{ runtimeTable?.TotalData.TotalReturnRate }}</td>
            <td>{{ runtimeTable?.SevenDays.TotalReturnRate }}</td>
            <td>{{ runtimeTable?.MonthData.TotalReturnRate }}</td>
            <td>{{ runtimeTable?.TotalData.buyRateTotalReturnRate }}</td>
            <td>{{ runtimeTable?.SevenDays.buyRateTotalReturnRate }}</td>
            <td>{{ runtimeTable?.MonthData.buyRateTotalReturnRate }}</td>
            <template v-if="gameMsgProps">
                <td>{{ runtimeTable?.TotalData.extra }}</td>
                <td>{{ runtimeTable?.SevenDays.extra }}</td>
                <td>{{ runtimeTable?.MonthData.extra }}</td>
            </template>
        </tr>

        </tbody>
    </table>
</template>

<script setup lang="ts">
import {computed, reactive, ref, onUpdated} from "vue";
import ut from "@/lib/util";
import {useStore} from '@/pinia/index';
import {storeToRefs} from "pinia";
const store = useStore()
const { runtimeTable } = storeToRefs(store)
interface Props {
    SevenDays: any
    MonthData: any
    TotalData: any
    tableData: any
    allTableData: any
    gameMsg:any
}
const props = withDefaults(defineProps<Props>(), {})
let gameMsgProps = computed(() => props.gameMsg?.startsWith('jili'))
console.log(gameMsgProps);
// let isJili = gameMsgProps.value.startsWith('jili')
/*let runtimeTable = reactive({
    SevenDays: props.SevenDays,
    MonthData: props.MonthData,
    TotalData: props.TotalData,
    tableData: props.tableData,
    allTableData: props.allTableData,
})*/
// console.log(isJili.value.startsWith('jili'));
onUpdated(() => {
    init()
});
let init = () => {
    // console.log(runtimeTable.value.TotalData.BetAmount);
    // console.log('---');
    /*if (!runtimeTable.allTableData || Object.keys(runtimeTable.allTableData).length === 0) {
        return
    }
    let LastWeek = JSON.parse(JSON.stringify(runtimeTable.allTableData.WeekData.LastWeek))
    let NowWeek = JSON.parse(JSON.stringify(runtimeTable.allTableData.WeekData.NowWeek))
    runtimeTable.tableData = NowWeek.concat(Object.values(LastWeek))
    console.log(runtimeTable.tableData);*/
    /*runtimeTable.SevenDays = returnData(runtimeTable.allTableData.SevenData)
    runtimeTable.MonthData = returnData(runtimeTable.allTableData.MonthData)
    runtimeTable.TotalData = returnData(runtimeTable.allTableData.TotalData)*/

}

</script>

<style scoped>
.table{
    border-collapse: collapse;
    width:100%;
    border:1px solid #c6c6c6 !important;
    margin:20px 0;
}
.table th{
    border-collapse: collapse;
    border-right:1px solid #c6c6c6;
    border-bottom:1px solid #c6c6c6;
    background-color:#ddeeff;
    padding:5px 9px;
    font-size:14px;
    font-weight:normal;
    text-align:center;
}
.table td{
    border-collapse: collapse;
    border-right:1px solid #c6c6c6;
    border-bottom:1px solid #c6c6c6;
    padding:5px 9px;
    font-size:12px;
    font-weight:normal;
    text-align:center;
    word-break: break-all;
}
.table tr:nth-child(odd){
    background-color:#fff !important;
}
.table tr:nth-child(even){
    background-color: #f8f8f8 !important;
}
</style>
