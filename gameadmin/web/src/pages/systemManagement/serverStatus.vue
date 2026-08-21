<template>
    <div class='gameList'>
        <div class="searchView gameList" >
            <el-space>
<!--                <el-autocomplete
                    v-model="autocompleteState"
                    :fetch-suggestions="querySearch"
                    clearable
                    placeholder="搜索相应游戏即可查询"
                    @select="handleSelect"
                >
                    <template #default="{ item,index }">
                        <span class="link">
                            <el-tag type="info" effect="light">
                                {{ item.value }}
                            </el-tag>
                        </span>
                    </template>
                </el-autocomplete>-->
                <el-switch
                    v-model="switchValue" @change="switchValueChange"
                    size="small"
                    active-text="正常"
                    inactive-text="其他"
                />
            </el-space>
        </div>

        <el-row :gutter="12">
            <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="4" class="elCard"
                    v-for="(item,name) in getServicesStatusList" v-loading="loading">
                <el-badge is-dot class="item" :type="item.length?'primary':'red'">
                    <el-card shadow="hover">
                        <div slot="header" class="clearfix">
                            <el-space>
                                <div class="el-space-class el-space-class-left">
                                    <el-text type="info">
                                        {{ $t('游戏名称') }}
                                    </el-text>
                                    <el-tag type="info">{{ name }}</el-tag>
                                </div>
                                <div class="el-space-class">
                                    <template v-if="item.length === 1">
                                        <div v-for="(i,index) in item">
                                            <el-tag type="info">
                                                {{ $t('服务器') }}：{{ i.ServiceName }}
                                            </el-tag>
                                            <el-tag type="info">
                                                {{ $t('主机名称') }}：{{ i.Hostname }}
                                            </el-tag>
                                            <el-tag type="info">
                                                {{ $t('进程id') }}：{{ i.Pid }}
                                            </el-tag>
                                            <el-tag type="info">
                                                {{ $t('在线时长') }}：{{ i.Uptime }}
                                            </el-tag>
                                        </div>
                                    </template>
                                    <template v-else>
                                        <div v-for="(i,index) in item">
                                            <template v-if="index === 0">
                                                <el-tag type="info">
                                                    {{ $t('服务器') }}：{{ i.ServiceName }}
                                                </el-tag>
                                                <el-tag type="info">
                                                    {{ $t('主机名称') }}：{{ i.Hostname }}
                                                </el-tag>
                                                <el-tag type="info">
                                                    {{ $t('进程id') }}：{{ i.Pid }}
                                                </el-tag>

                                                <el-popover placement="right" :width="400" trigger="click">
                                                    <template #reference>
                                                        <el-button text size="small">
                                                            {{ $t('查看更多') }}
                                                        </el-button>
                                                    </template>
                                                    <el-space>
                                                        <div v-for="(i,index) in item" class="el-space-class">
                                                            <el-tag type="info">
                                                                {{ $t('服务器') }}：{{ i.ServiceName }}
                                                            </el-tag>
                                                            <el-tag type="info">
                                                                {{ $t('主机名称') }}：{{ i.Hostname }}
                                                            </el-tag>
                                                            <el-tag type="info">
                                                                {{ $t('进程id') }}：{{ i.Pid }}
                                                            </el-tag>
                                                            <el-tag type="info">
                                                                {{ $t('在线时长') }}：{{ i.Uptime }}
                                                            </el-tag>
                                                        </div>
                                                    </el-space>
                                                </el-popover>
                                            </template>
                                        </div>
                                    </template>
                                </div>
                            </el-space>
                        </div>
                    </el-card>
                </el-badge>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick, watch} from 'vue';
import type {ElInput, FormInstance, FormRules, UploadFile,} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {Game, GameStatus, GameType} from '@/api/gamepb/customer';
import ReturnErrorHandle, {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import {Delete, Download, Plus, ZoomIn} from '@element-plus/icons-vue'

let uploadImageUrl = ref(null)
const {t} = useI18n()
const store = useStore()
import ut from "@/lib/util";

let forms = reactive(<Game>{})
const gameListCollapse = ref('1')
let getServicesStatusList = ref(null)
let getServicesStatusListNew = ref(null)

let loading = ref(false)
let editPlat = ref(false)

const queryList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.getServicesStatus, {})
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    getServicesStatusList.value = data
    getServicesStatusListNew.value = data
    let arr = Object.keys(data)
    console.log(arr);
    restaurants.value = arr.map(l => {
        return {
            value: l,
        }
    }).filter(Boolean)
}
const autocompleteState = ref('')
const switchValue = ref(true)
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

const switchValueChange = async () => {
    const servicesArray = Object.keys(getServicesStatusList.value)
        .filter(key => Array.isArray(getServicesStatusList.value[key]) && getServicesStatusList.value[key].length === 0)
        .reduce((acc, key) => ({ ...acc, [key]: getServicesStatusList.value[key] }), {});

    console.log(servicesArray);
    getServicesStatusList.value = switchValue.value ? getServicesStatusListNew.value : servicesArray;
}

const handleSelect = async (item) => {
    if (item.value) {
        const key = item.value;
        const value = getServicesStatusListNew.value[key];

        // 检查 servicesStatus 是否为数组
        if (Array.isArray(value)) {
            getServicesStatusList.value[item.value] = value;
        } else {
            console.error('The data corresponding to item.value is not an array');
        }
    } else {
        getServicesStatusList.value = getServicesStatusListNew.value
    }
}

onMounted(queryList);

let gameName = t("游戏列表")

import {AdminInfo} from "@/api/adminpb/info";


</script>
<style scoped lang='scss'>

.elCard {
    border: none;
    margin-bottom: 1rem;

    .el-card {
        border-radius: .5rem;
    }

    .topCard {
        //padding-left: 1.5rem;
        display: flex;
        align-items: flex-start;
        flex-direction: column;
        flex-wrap: wrap;
        align-content: flex-start;

        .topCardId {
            font-size: 10px;
            color: #999999;
            margin-bottom: .2rem;
        }

    }

    .bottom {
        display: flex;
        justify-content: space-between;
        width: 100%;
        position: relative;
        //padding-left: 1.5rem;
        font-size: 13px;

        div {
            padding: 0 3px;
            line-height: 24px;
        }

    }
}

.avatar {
    max-width: 200px;
}

.small-font {
    margin: 0;
}

:deep(.el-card__body) {
    padding: 15px;
}

.el-space-class {
    display: flex;
    flex-direction: column;
    align-items: center;
    span {
        margin-bottom: .5em;
    }
}

.el-space-class-left {
    width: 100px;
}

:deep(.el-badge__content.is-dot) {
    width: 12px;
    height: 12px;
}

:deep(.el-badge__content--primary) {
    background: #176c0a;
}
</style>
