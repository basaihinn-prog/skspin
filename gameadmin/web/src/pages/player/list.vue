<template>
    <div class='gameList'>
        <div class="searchView">
            <el-form
                label-position="top"
                label-width="100px"
                :model="param"
                style="max-width: 100%"
            >
                <el-space wrap>
                    <el-form-item :label="$t('所属产品')">
                        <el-select v-model="param.AppID" clearable :placeholder="$t('请选择')">
                            <el-option v-for='value in rhs.AppList' :label="value.AppID" :value="value.AppID"/>
                        </el-select>
                    </el-form-item>
                </el-space>
            </el-form>
            <el-button type="primary" class="" @click="getGameList()">{{ $t('搜索') }}</el-button>
        </div>
        <!-- 数据 -->
        <el-table :data="rhs.AppList" class="elTable" v-loading="loading"
                  :header-cell-style="{ background: '#F5F7FA', color: '#333333' }">
            <el-table-column prop="Pid" label="Pid" width="80px"/>
            <el-table-column prop="Uid" label="Uid"/>
            <el-table-column prop="AppID" label="AppID"/>
            <el-table-column prop="LoginAt" :label="$t('最后登录')" min-width="200px" :formatter="dateFormater"/>
            <el-table-column prop="CreateAt" :label="$t('注册时间')" :formatter="dateFormater"/>
            <el-table-column :label="$t('操作')" width="105px">
                <template #default="scope">
                    <router-link class="noneA" :to="{ path: '/playerInfo', query: { Pid: scope.row.Pid } }">
                        <el-button type="primary">{{ $t('详情') }}</el-button>
                    </router-link>
                </template>
            </el-table-column>
        </el-table>
        <!-- 分页 -->
        <Pagination :Count='rhs.Count' @paginationEmit='paginationEmit'></Pagination>

    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick, watch} from 'vue';
import type {ElInput, FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {AdminGameCenter} from '@/api/gamepb/admin';
import {useStore} from '@/pinia';
import {AdminInfo} from "@/api/adminpb/info";

const {t} = useI18n()
const store = useStore()
let operatorData = reactive([])
let Count = ref(0)
let param = reactive({
    PageIndex: 1,
    PageSize: 20,
    AppID: ''
})
let rhs = reactive({
    AppList: [],
    Count: 0,
    PlayerList: []
})

let loading = ref(false)


onMounted(async () => {
    let [data, err] = await Client.Do(AdminGameCenter.NewAppList, param)
    if (err) {
        return tip.e(err)
    }
    rhs.AppList = data.List

    getGameList()

});
const operatorList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, param)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    console.log(data,'data');
    Count.value = data?.AllCount
    operatorData = data.AllCount === 0 ? [] : data.List.filter(list=> !list.Status)
}
const getGameList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminGameCenter.PlayerList, param)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    rhs.PlayerList = data.List
    rhs.Count = data.Count
}

const paginationEmit = (PageIndex, PageSize) => {
    param.PageIndex = PageIndex
    param.PageSize = PageSize
    getGameList()
}
</script>
<style scoped lang='scss'></style>
