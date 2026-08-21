import VNodes, { modelProp, modelPropNumber, r } from '../../lib/VNodes.mjs'
import { onMounted, ref, reactive, h } from 'vue';
import { Client } from '@/lib/client';
import ut from '@/lib/util'
import { tip } from '@/lib/tip';

export default {
    setup
}

// const gameName = 'pg_39'
async function callFunc(api, args) {
    return Client.send(`mq/pg_1489936/AdminSlotsRpc/${api}`, args)
}

const bets = [1, 75]
function setup() {
    const tableData = ref([])
    const combineData = ref([])
    const loading = ref(false)
    let totalCount = [0, 0]
    let totalTimes = [0, 0]

    onMounted(() => {
        console.log("onMounted!!")
        refreshData()
    })

    async function refreshData() {
        let [data, err] = await callFunc("GetPoolStatus", {})
        loading.value = false
        if (err) {
            console.error(err)
            return tip.e(err)
        }
        tableData.value = data.RateStatus

        totalCount = [0, 0]
        totalTimes = [0, 0]
        for (const vv of data.RateStatus) {
            const t = vv[0].Type
            totalCount[t] += vv.reduce((sum, item) => sum + item.DstCount, 0)
            totalTimes[t] += vv.reduce((sum, item) => sum + item.DstTimes, 0)
        }

        let [combine, err2]  = await  callFunc("combine/get", {})
        if (err2) {
            return tip.e(err2)
        }

        combineData.value = combine

    }

    async function ExtractAction(index) {
        let datas = tableData.value[index]
        // console.log(datas)
        let args = []
        for (let i = 0; i < datas.length; i++) {
            let data = datas[i]
            if (ut.isNull(data.DstCount)) {
                data.DstCount = 0
            }

            if (data.CurrentCount <= 0 && data.DstCount > 0) {
                return tip.e(`区间(${data.MinRate},${data.MaxRate}]，没有生成的，不可配置 `)
            }
            args.push({ index: data.Index, count: data.DstCount })
        }
        loading.value = true
        let [data, err] = await callFunc("ExtractAction", args)
        loading.value = false

        console.log(data)
        if (err) {
            return tip.e(err)
        } else {
            await refreshData()
            return tip.s("succss")
        }
    }
    async function  genCombineDataAction() {
        await callFunc("combine/gen", combineData.value)
    }

    return () =>  
        h('div', { style:{'text-align': 'center'}}, [
            r('el-button', {onClick: refreshData, type: 'primary'}, ()=>'刷新'),
            h('br'),
            "Normal RTP: " + fmtfloat( totalTimes[0] / totalCount[0]),
            h('br'),
            "Buy RTP: " + fmtfloat( totalTimes[1] / (totalCount[1] * bets[1])),
            r('el-table', { data: tableData.value, style: { width: '100%' } }, () => [
                r('el-table-column', { label: '类', width: "40" }, ({ row }) => row.map(item => h('div', item.Type))),
                r('el-table-column', { label: '区间' }, ({ row }) => row.map(item => h('div', item.Flag))),
                r('el-table-column', { label: '库存' }, ({ row }) => row.map(item => h('div',  item.CurrentCount))),
                // r('el-table-column', { label: '剩余' }, ({ row }) => row.map(item => h('div', item.NotUseCombineCount))),
                // r('el-table-column', { label: '掉个数' }, ({ row }) => row.map(item => h('div', item.UseCombineCount))),
                r('el-table-column', { label: '配置个数' }, ({ row }) => row.map(item => r('el-input', {
                    placeholder: "配置个数",
                    modelValue: item.DstCount.toString(),
                    "onUpdate:modelValue": val => item.DstCount = parseInt(val),
                }))),
                r('el-table-column', { label: '分布比例' }, ({ row }) => [
                    h('div', { style: { display: 'flex', 'justify-content': 'space-around', 'align-items': 'flex-end', 'line-height': '40px' } }, [
                        h('div', null, row.map(item=>h('div', fmtfloat( item.DstCount / totalCount[item.Type])))),
                        h('div', {style:{'padding-left': '20px'}},  fmtfloat( row.reduce((sum, item) => sum + item.DstCount, 0)/ totalCount[row[0].Type] )),
                    ]),
                ]),
                r('el-table-column', { label: '中奖倍数',  }, ({ row }) => row.map(item => h('div', { style: {'line-height': '40px'}}, item.DstTimes.toFixed(2)))),
                r('el-table-column', { label: '中奖倍数-奖池',  }, ({ row }) => row.map(item => h('div', { style: {'line-height': '40px'}}, (item.DstTimes-item.PoolCost*item.DstCount).toFixed(2)))),

                r('el-table-column', { label: '回报分布' }, ({ row }) => [
                    h('div', { style: { display: 'flex', 'justify-content': 'space-around', 'align-items': 'flex-end', 'line-height': '40px' } }, [
                        h('div', null, row.map(item=>h('div', fmtfloat( item.DstTimes / totalCount[item.Type]/ bets[item.Type])))),
                        h('div', {style:{'padding-left': '20px'}},  fmtfloat( row.reduce((sum, item) => sum + item.DstTimes, 0) / totalCount[row[0].Type] / bets[row[0].Type] )),
                    ]),
                ]),

                r('el-table-column', { label: '操作' }, (scope) => r('el-button', { style: { 'vertical-align': 'bottom' }, type: 'primary', onClick: () => ExtractAction(scope.$index) }, () => '提取')),
                
            ]),
            
            h('h4', {style:{'text-align': 'left'}}, '组合数组设置'),
            r('el-table', {data:  combineData.value}, ()=>[ 
                r('el-table-column', {label: '组合', prop: 'Name'}),
                r('el-table-column', {label: '配置组数', width: '180'}, ({row})=>
                    r('el-input', {
                        placeholder: '配置组数', 
                        modelValue: row.Count.toString(),
                        "onUpdate:modelValue": val => row.Count = parseInt(val),
                    })),
             ]),

             h('div', {style:{'margin': '20px', 'display': 'flex'}}, 
                r('el-button', {type:'primary', onClick: genCombineDataAction}, ()=> '生成组合数据')
            ),

        ])
}

function fmtfloat(fv) {
    fv *= 100;
    return fv.toFixed(3).toString() + '%'
}