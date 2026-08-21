
<template>
    <div>
        <div class="searchView">
            <el-form
                label-position="top"
                label-width="100px"
                :model="searchList"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operator="operatorListChange"></operator_container>

                    <el-form-item :label="$t('投注时间')">
                        <el-date-picker
                            v-model="searchList.times"
                            type="datetimerange"
                            value-format="x"
                            format="YYYY-MM-DD HH:mm:ss"
                            range-separator="To"
                            :shortcuts="shortcuts"
                            start-placeholder="start time"
                            end-placeholder="end time"
                        />
                    </el-form-item>
                </el-space>
            </el-form>
            <el-space>
                <el-button type="primary" @click="queryListSearch" >
                    {{ $t('搜索') }}
                    <el-icon style="margin-left: .5rem"><Search /></el-icon>
                </el-button>
                <el-tag type="info">{{$t('下载记录的有效期为7天，超过有效期下载记录自动删除')}}</el-tag>
            </el-space>
        </div>
        <!-- 数据 -->
        <el-table :data="tableData" class="elTable" v-loading="loading"
                  :header-cell-style="{ background: '#F5F7FA', color: '#333333' }">
            <el-table-column type="index" :label="$t('序号')" width="80px" />
            <el-table-column prop="Username" :label="$t('用户名称')"   />
            <el-table-column prop="OperatorName" v-if="store.AdminInfo.AppID === 'admin'" :label="$t('商户')" />
            <el-table-column prop="FileName" :label="$t('文档名称')" />
            <el-table-column prop="StartTime" :label="$t('开始日期')"  :formatter="dateSecondFormater" />
            <el-table-column prop="EndTime" :label="$t('结束日期')" :formatter="dateSecondFormater" />
            <el-table-column prop="ProcessStatus" :label="$t('状态')">
                <template #default="scope">
                    {{mapStatus[scope.row.ProcessStatus]}}
                </template>
            </el-table-column>
            <el-table-column prop="CreateAt" :label="$t('创建日期')" :formatter="dateSecondFormater" />
            <el-table-column prop="InsertTime" :label="$t('操作')" :formatter="dateSecondFormater" >
                <template #default="scope">
                    <el-button circle type="" @click="queryListSearch" >
                        <el-icon><Refresh /></el-icon>
                    </el-button>
                    <el-button circle type="" :disabled="!mapStatus[scope.row.ProcessStatus]" @click="GetBHdownload(scope.row)" >
                        <el-icon><Download /></el-icon>
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
        <!-- 分页 -->
        <Pagination :Count='searchList.Count' @paginationEmit='paginationEmit'></Pagination>
    </div>
</template>
<script setup lang='ts'>
import {onMounted, reactive, ref} from 'vue';
import {useStore} from '@/pinia/index';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useI18n} from 'vue-i18n';
import {AdminGameCenter} from '@/api/gamepb/admin';
import Operator_container from "@/components/operator_container.vue";
import {Download, Refresh} from "@element-plus/icons-vue";
import {AdminInfo} from "@/api/adminpb/info";
import {useRoute} from 'vue-router'
import {storeToRefs} from 'pinia'

const store = useStore()
const { activeTabs } = storeToRefs(store)
const route = useRoute()
const defaultOperatorEvent = ref({})
let loading = ref(false)
let tableData =  ref([])
let pageList = reactive({
    row: null,
    dialogVisible: false,
    isPg:false,
    whitchPg:'',
})
let searchList = reactive({
    OperatorId:null,
    PageSize:20,
    PageNumber:1,
    Count:0,
    StartTime:null,
    EndTime:null,
    times:null,
})
let {t} = useI18n()
const mapStatus = [
    /** 已生成下载任务 */
    t('已生成下载任务'),
    /** 下载任务处理中 */
    t('下载任务处理中'),
    /** 下载文件已生成 */
    t('下载文件已生成'),
]
const shortcuts = [
    {
        text: t('过去7天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            return [Date.parse(start.toString()), Date.parse(end.toString())]
        },
    },
    {
        text: t('今天'),
        value: () => {
            const today = new Date();
            const year = today.getFullYear();
            const month = today.getMonth();
            const date = today.getDate();

            const startTime = new Date(year, month, date, 0, 0, 0, 0o00).getTime();
            const endTime = new Date(year, month, date, 23, 59, 59, 0o00).getTime();

            return [startTime, endTime];
        },
    },
]


const queryListSearch = async () => {
    await queryList()
}
function getEnrichedList(data) {
    const operatorMap = new Map(operatorList.value.map(op => [op.Id, op]));
    return data.map(item => ({
        ...item,
        OperatorName: operatorMap.get(Number(item.OperatorId))?.Name || null
    }));
}
const queryList = async () => {
    let dataList  = {
        ...searchList,
        PageSize:Number(searchList.PageSize),
        StartTime:searchList.times?searchList.times[0]/1000:null,
        EndTime:searchList.times?searchList.times[1]/1000:null,
    }
    let [data, err] = await Client.Do(AdminGameCenter.GetExcleBetLogList, dataList)
    loading.value=false
    if (err) {
        return tip.e(err)
    }
    searchList.Count = data.Count
    tableData.value = data.Count === 0 ? [] : data.List
    console.log(tableData.value);
    console.log(data.List);
    if (store.AdminInfo.AppID === 'admin') {
        await getOperator()
        console.log(operatorList.value);
        if (operatorList.value.length && data.List) {
            tableData.value = getEnrichedList(data.List)
        }
    }
}
const paginationEmit = (PageIndex, PageSize) => {
    searchList.PageNumber = PageIndex
    searchList.PageSize = PageSize
    queryList()
}

let operatorList = ref([])
const getOperator = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, {
        PageIndex: 1,
        PageSize: 10000
    })
    if (err) {
        return tip.e(err)
    }
    operatorList.value = data.AllCount === 0 ? [] : data.List
}
onMounted(() => {

    console.log(route.params);
    if (route.query.OperatorId) {
        let times = route.query.Times;
        if (Array.isArray(times)) {
            times = times.map(list => {
                return Number(list)
            });
        } else if (typeof times === 'string') {
            times = [new Date(times).getTime()];
        }
        searchList.times = times
        searchList.OperatorId = Number(route.query.OperatorId)

        queryListSearch()

    }

});
const operatorListChange = (value) =>{
    searchList.OperatorId = value.value
}
const GetBHdownload = async (item) => {
    try {
        const response = await fetch('/api/AdminInfo/BHdownload', {
            method: 'POST',  // 默认情况下 fetch 使用 GET 请求，你可能需要使用 POST
            headers: {
                'Content-Type': 'application/json',
                'Authorization':store.Token
            },
            body: JSON.stringify({ "ID": item.ID })
        });

        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }

        const data = await response.blob();  // 使用 response.blob() 而不是 response.text()
        const blob = new Blob([data]);
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = item.FileName || 'downloaded_file';  // 设置下载文件的默认名称
        document.body.appendChild(a);
        a.click();
        setTimeout(() => {
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        }, 0);
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
        return tip.e(t('下载文件失败'));
    }
};


function isWithinSevenDays(startTime, endTime) {
    const diffTime = endTime - startTime;
    const sevenDaysInMilliseconds = 7 * 24 * 60 * 60 * 1000;
    return diffTime <= sevenDaysInMilliseconds;
}
</script>
<style lang="scss" scoped>
.el-table .cell {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.dialog-top{
    background: #282834;
    font-size: 14px;
    color: hsla(0,0%,100%,.6);
    .el-col-6{
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        text-align: center;
        .top{
            display: flex;
            justify-content: center;
            flex-direction: column;
            height: 100%;
            word-wrap: break-word;
            color: rgba(255, 255, 255, 0.6);
            p{
                text-align: center;
            }
        }
        .color-text{
            text-align: center;
            word-wrap: break-word;
            color: rgb(88, 245, 109);
        }
    }

}
:deep .el-carousel__container{
    min-height: 600px;
    max-height: 900px;
    height: auto !important;
}
.py-10{
    padding: .5rem;
}

</style>
