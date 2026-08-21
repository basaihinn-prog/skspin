<template>

    <div style="text-align: center; margin-top: 1rem">
        <el-collapse>
            <el-collapse-item title="pp游戏列表" name="1">
                <div style="padding: 1rem">
                    <el-radio-group v-model="checkedModel" size="small" @change="checkedChange" border="false">
                        <el-radio-button v-for="(item,index) in checkedList" :key="item.value +index" :label="item">
                            <el-space>
                                <el-tag type="success" effect="light">
                                    {{ item.value }}
                                </el-tag>
                                <el-tag type="info" effect="light">
                                    {{ item.label }}
                                </el-tag>
                            </el-space>
                        </el-radio-button>
                    </el-radio-group>
                </div>
            </el-collapse-item>
        </el-collapse>


        <div style="margin-top: 1rem">
            <el-space>
                <el-autocomplete
                    v-model="autocompleteState"
                    :fetch-suggestions="querySearch"
                    clearable
                    placeholder="搜索相应游戏即可查询"
                    @select="handleSelect"
                >
                    <template #suffix>
                        <el-icon class="el-input__icon">
                            <edit/>
                        </el-icon>
                    </template>
                    <template #default="{ item }">
                        <span class="value">
                            <el-tag type="success" effect="light">
                                {{ item.value }}
                            </el-tag>
                        </span>
                        <span class="link" style="float: right">
                            <el-tag type="info" effect="light">
                                {{ item.label }}
                            </el-tag>
                        </span>
                    </template>
                </el-autocomplete>
                <el-button @click="refreshData" type="primay" style="margin:0" :disabled="!checkedModel.ID">
                    <el-space>

                        <template v-if="checkedModel.ID">
                            刷新
                            <el-tag type="success" effect="light">
                                {{ checkedModel.ID }}
                            </el-tag>
                            <el-tag type="info" effect="light">
                                {{ checkedModel.label }}
                            </el-tag>
                        </template>
                        <template v-else>
                            <span style="color: #999999">请选择相应游戏</span>
                        </template>
                    </el-space>
                </el-button>
            </el-space>
            <el-row :gutter="20" style="margin-top: 1rem">
                <el-col :span="4">
                    <el-space>
                        <el-statistic title="Normal RTP"
                                      :value="ut.fmtPercent(uiData.totalTimes[0] / uiData.totalCount[0])"/>
                        <el-divider direction="vertical"/>
                        <el-statistic title="Count" :value="uiData.totalCount[0]"/>
                    </el-space>
                </el-col>
                <el-col :span="20">
                    <el-space  v-for="(item,index) in typeFor" style="margin-right: 2rem">
                        <el-statistic title="下注倍数" :value="typeMulti[index].toString()" />
                        <el-divider direction="vertical"/>
                        <el-statistic title="Buy RTP"
                                      :value="(uiData.totalCount[item] * typeMulti[index])?
                                          ut.fmtPercent(uiData.totalTimes[item] / (uiData.totalCount[item] * typeMulti[index]))
                                          : $t('暂无')"/>
                        <el-divider direction="vertical"/>

                        <el-statistic title="Count" :value="uiData.totalCount[item]"/>
                    </el-space>
                </el-col>
            </el-row>

        </div>
        <el-table :data="tableData" style="width: 100%" v-loading="loading" fit stripe highlight-current-row
                  @header-click="tableDataHeader">
            <el-table-column label="类" width="50">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.Type }}
                    </div>
                </template>
            </el-table-column>

            <el-table-column label="区间" width="200" fixed>
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.Flag }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="库存">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.CurrentCount }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="组合占用">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.UseCombineCount }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="配置个数">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        <el-input placeholder="配置个数" v-model.number="item.DstCount"></el-input>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="分布比例" class-name="flex-150">
                <template #default="scope">
                    <div class="flex-cell">
                        <div>
                            <div :key="index" v-for="(item, index) in scope.row">
                                {{
                                    uiData.totalCount[item.Type]?
                                        ut.fmtPercent(item.DstCount / uiData.totalCount[item.Type]):
                                        '--'
                                }}
                            </div>
                        </div>
                        <div style="padding-left: 1rem">{{
                                uiData.totalCount[scope.row[0].Type]?
                                    ut.fmtPercent(scope.row.reduce((sum, item) => sum + item.DstCount, 0) / uiData.totalCount[scope.row[0].Type]):
                                    '--'
                            }}
                        </div>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="中奖倍数">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.DstTimes.toFixed(2) }}
                    </div>
                </template>
            </el-table-column>

            <el-table-column label="奖池">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ item.PoolCost }}
                    </div>
                </template>
            </el-table-column>

            <el-table-column label="中奖倍数-奖池">
                <template #default="scope">
                    <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                        {{ (item.DstTimes - item.PoolCost * item.DstCount).toFixed(2) }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="回报分布">
                <template #default="scope">
                    <div style="display: flex;justify-content: flex-start;align-items: flex-end;line-height: 40px;">
                        <div>
                            <div :key="index" style="line-height: 40px;" v-for="(item, index) in scope.row">
                                {{
                                    (uiData.totalCount[item.Type] && bets[item.Type])?
                                        ut.fmtPercent((item.DstTimes- item.PoolCost * item.DstCount) / uiData.totalCount[item.Type] / bets[item.Type]):
                                        '--'
                                }}
                            </div>
                        </div>
                        <div style="padding-left: 20px">
                            {{
                                (uiData.totalCount[scope.row[0].Type] && bets[scope.row[0].Type])?
                                    ut.fmtPercent(scope.row.reduce((sum, item) => sum + (item.DstTimes- item.PoolCost * item.DstCount), 0) / uiData.totalCount[scope.row[0].Type] / bets[scope.row[0].Type]): '--'
                            }}
                        </div>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button @click="ExtractAction(scope.$index)" style="vertical-align: bottom;margin-bottom: 0"
                               type="primary">提取
                    </el-button>
                </template>
            </el-table-column>

        </el-table>
        <div style="text-align:left;margin-top: 1rem">组合数组设置</div>
        <el-table :data="combineData" v-loading="loading">
            <!-- <el-table-column label="组合" prop="name"></el-table-column> -->
            <el-table-column label="组合" prop="Name"></el-table-column>
            <el-table-column label="配置组数" width="180">
                <template #default="scope">
                    <el-input placeholder="配置组数" v-model.number="scope.row.Count"></el-input>
                </template>
            </el-table-column>
        </el-table>
        <el-button style="margin-top: 1rem" @click="genCombineDataAction" type="primary">生成组合数据</el-button>
    </div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
// import { AdminZhaoCaiMaoRpc } from '@/api/zhaocaimao/admin_rpc';
// import { AdminBaoZangRpc } from '@/api/baozang/admin_rpc';
// import { AdminNiubiRpc } from '@/api/niubi/admin_rpc';
// import { AdminRomaRpc } from '@/api/roma/admin_rpc';
// import { AdminRomaXRpc } from '@/api/romax/admin_rpc';
import ut from '@/lib/util'
import {useRoute} from 'vue-router'
import {AdminGameCenter} from "@/api/gamepb/admin";

const route = useRoute()
let uiData = reactive({
    times: 0,
    value: '',
    combineData: [],
    totalCount: new Map(),
    totalTimes: new Map(),
    gameName: ""
})
const tableData = ref([])
const checkedList = ref(null)
const checkedModel = ref({})
const combineData = ref([])
let gid = ref("pp_000")
let buy = ref(999)
let bets = ref([1, 999])
let loading = ref(false)
const autocompleteState = ref('')
const restaurants = ref([])
const querySearch = (queryString: string, cb: any) => {
    const results = queryString
        ? restaurants.value.filter(createFilter(queryString))
        : restaurants.value
    // call callback function to return suggestions
    cb(results)
}
const createFilter = (queryString: string) => {
    return (restaurant) => {
        return (
            restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
        )
    }
}
const handleSelect = async (item) => {
    console.log(item);
    isDisType.value = false
    typeFor.value = []
    checkedModel.value = item
    uiData.gameName = gid.value = item.ID
    buy.value = item.BuyBetMulti
    bets.value = [1, buy.value]
    loading.value = true
    await refreshData()
}
onMounted(async () => {
    try {
        let [datas, errs] = await Client.Do(AdminGameCenter.GameList, {})
        if (errs) {
            return tip.e(errs)
        }
        // 检查datas.List的存在性和类型
        if (!Array.isArray(datas.List)) {
            tip.e('Invalid data format');
        }
        // 合并map和filter操作
        const filteredList = datas.List.reduce((acc, l) => {
            if (typeof l.ID === 'string' && l.ID.startsWith('pp')) {
                acc.push({
                    ...l,
                    label: l.Name,
                    value: l.ID
                });
            }
            return acc;
        }, []);
        checkedList.value = filteredList;
        restaurants.value = filteredList;
    } catch (err) {
        tip.e(err.message || 'An unknown error occurred');
    }
})

const checkedChange = async () => {
    isDisType.value = false
    typeFor.value = []
    uiData.gameName = gid.value = checkedModel.value.ID
    buy.value = checkedModel.value.BuyBetMulti
    bets.value = [1, buy.value]
    loading.value = true
    await refreshData()
}

async function callFunc(api: string, args: any) {
    // return Client.send(`mq/${uiData.gameName}/AdminSlotsRpc/${api}`, { Data: JSON.stringify(args) })
    return Client.send(`mq/${uiData.gameName}/AdminSlotsRpc/${api}`, args)
}

async function refreshCombine() {
    let [combine, err] = await callFunc("combine/get", {})
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    combineData.value = combine
}

const isDisType = ref(false)
const typeFor = ref([])
const typeMulti = ref([])
const refreshData = async () => {
    let [data, err] = await callFunc("GetPoolStatus", {})
    loading.value = false
    if (err) return tip.e(err)
    const uniqueTypes = new Set();
    const totalCountByType = {};
    const totalTimesByType = {};
    tableData.value = data.RateStatus
    isDisType.value = data.DisType && data.DisType === 1;

    typeFor.value = []
    typeMulti.value = []


    data.RateStatus.forEach(list => {
        list.forEach(op => {

            if (op.Type !== 0) uniqueTypes.add(op.Type);
            if (!totalCountByType[op.Type]) {
                totalCountByType[op.Type] = 0;
                totalTimesByType[op.Type] = 0;
            }
            totalCountByType[op.Type] += op.DstCount;
            totalTimesByType[op.Type] += op.DstTimes - op.PoolCost*op.DstCount;
        });
    });
    typeFor.value = [...uniqueTypes];
    console.log(buy.value);
    for (const opType of typeFor.value) {
        if (data?.BuyGame) typeMulti.value.push(buy.value);

    }
    uiData.totalCount = totalCountByType;
    uiData.totalTimes = totalTimesByType;
    await refreshCombine()

}
const tableDataHeader = async (column) => {
    const columnMap = {
        2: {name:'库存',data:'CurrentCount'},
        3: {name:'组合占用',data:'UseCombineCount'},
        4: {name:'配置个数',data:'DstCount'},
        6: {name:'中奖倍数',data:'DstTimes'},
        7: {name:'奖池',data:'PoolCost'}
    };
    if (columnMap.hasOwnProperty(column.no)) {
        try {
            const summaries = await Promise.all(Object.keys(columnMap).map(async (key) => {
                const label = columnMap[key];
                return await getSummaries(label.data, label.name);
            }));
            console.log(summaries);
            // 生成表格样式的内容
            const columnData = transformToColumns(summaries);
            const csvContent = columnData.join('\n');
            await navigator.clipboard.writeText(csvContent);
            tip.s('复制成功，可以粘贴到Excel中');
        } catch (err) {
            console.error('Failed to get summaries or copy data:', err);
            tip.e('Copy Failed');
        }
    } else {
        console.warn(`不支持的列号: ${column.no} ${column.label}`);
        return;
    }
}
const transformToColumns = (summaries) => {
    const rows = summaries.map(summary => summary.split('\n'));
    const maxRows = Math.max(...rows.map(row => row.length));

    return Array.from({length: maxRows}, (_, i) =>
        rows.map(row => row[i] || '').join(',')
    );
};
const getSummaries = async (sumValue,label) => {
    const validValues = tableData.value.flatMap(item =>
        Object.values(item).map(subItem => Number(subItem[sumValue]))
    ).filter(num => !Number.isNaN(num));

    if (validValues.length === 0) {
        return 'N/A';
    }
    const csvRows = [label, ...validValues];
    return csvRows.join('\n');
}
async function ExtractAction(index: number) {
    let datas = tableData.value[index]
    let args = []
    for (let i = 0; i < datas.length; i++) {
        let data = datas[i]
        if (ut.isNull(data.DstCount)) {
            data.DstCount = 0
        }

        if (data.CurrentCount <= 0 && data.DstCount > 0) {
            return tip.e(`区间(${data.MinRate},${data.MaxRate}]，没有生成的，不可配置 `)
        }
        args.push({index: data.Index, count: data.DstCount})
    }
    loading.value = true
    let [data, err] = await callFunc("ExtractAction", args)
    loading.value = false

    if (err) {
        return tip.e(err)
    } else {
        await refreshData()
        return tip.s("succss")
    }
}

async function genCombineDataAction() {
    await callFunc("combine/gen", combineData.value)
}

// const generatAction = async () => {
//     await callFunc("GenerateData", {})
// }
</script>
<style lang="scss">
.el-collapse-item__content{
    padding-bottom: 0;
    .el-radio-button__inner{
        border: none;
    }
    .el-radio-button:first-child .el-radio-button__inner{
        border-left: 0;
    }
}
.el-collapse-item__header{
    padding-left: 1rem;
}
.flex-150{
    min-width: 250px;
    .flex-cell{
        display: flex;justify-content: flex-start;align-items: flex-end;line-height: 40px;
        div{
            min-width: 60px;
        }
    }
}
</style>
