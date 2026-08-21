<template>
    <div class='gameList'>
        <div class="searchView" style="margin-bottom: 1rem">
            <el-date-picker
                v-model="datePickerValue" :editable="false" clearable
                @change="queryList"
                type="year"
                :disabled-date="pickerOptions"
                :prefix-icon="customPrefix"
                format="YYYY"
                value-format="YYYY"
            />
        </div>
        <el-row :gutter="12">
            <el-col :xs="12" :sm="8" :md="8" :xl="6" class="elCard" v-for="item in OpenTime.List" v-loading="loading">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <el-text style="font-size: 1rem;margin-right: 1rem" type="primary">
                                {{ item.OpenTime }}
                            </el-text>
                            <el-tag :type="item.IsOpened?'success':'danger'" class="mx-1" effect="plain" round>
                                {{ item.IsOpened ? $t('已开奖') : $t('未开奖') }}
                            </el-tag>
                            <div style="float: right">
                                <el-popconfirm
                                    width="220"
                                    :confirm-button-text="$t('确定开奖')"
                                    :cancel-button-text="$t('取消')"
                                    icon="InfoFilled"
                                    icon-color="#626AEF"
                                    @confirm="OpenPrizeConfirm(item.OpenTime)"
                                    :title="$t('确定要进行手动兑奖？') + item.OpenTime"
                                >
                                    <template #reference>
                                        <el-button v-if="!item.IsOpened" type="primary" text bg size="small"
                                                   style="margin-bottom: 0">
                                            {{ $t('手动兑奖') }}
                                        </el-button>
                                    </template>
                                </el-popconfirm>

                            </div>
                        </div>
                    </template>
                    <el-radio-group v-model="item.type" size="small">
                        <el-radio-button v-for="(item,index) in WebSideList" :label="index">
                            {{ item }}
                        </el-radio-button>
                    </el-radio-group>
                    <div style="float: right;font-size: 12px;position: relative;">
                        <el-popover :visible="item.popOver" placement="right" :width="800" trigger="click">
                            <template #reference>
                                <el-button type="info" text style="box-shadow: none" @click="toDetail(item)">
                                    {{ $t('查看详情') }}
                                </el-button>
                            </template>
                            <div class="dialog-openList-body fs-7 py-2" v-if="getPrizeByDateList">
                                <div class="" style="display:inline-block;width: 100%">
                                    <el-button size="small" text style="float: right" @click="item.popOver = false">{{$t('关闭')}}</el-button>
                                </div>
                                <el-row :gutter="20">
                                    <el-col :span="6">
                                        <el-card shadow="never">
                                            <template #header>
                                                <div class="card-header">
                                                    {{ $t(dialogTopList[0][0].name) }}
                                                </div>
                                            </template>
                                            <el-space wrap v-for="(item,index) in getPrizeByDateList.slice(0,1)">
                                                <div class="d-flex" v-for="i in item">
                                                    <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                                </div>
                                            </el-space>
                                        </el-card>
                                    </el-col>
                                    <el-col :span="6" wrap v-for="(item,index) in getPrizeByDateList.slice(1,4)"
                                            :key="index">
                                        <el-card shadow="never">
                                            <template #header>
                                                <div class="card-header">
                                                    {{ dialogTopList[1][index].name }}
                                                </div>
                                            </template>
                                            <el-space wrap v-for="i in item">
                                                <div class="d-flex">
                                                    <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                                </div>
                                            </el-space>
                                        </el-card>
                                    </el-col>
                                </el-row>

                                <el-row :gutter="20">
                                    <el-col :span="8">
                                        <el-card shadow="never">
                                            <template #header>
                                                <div class="card-header">
                                                    {{ dialogTopList[2][0].name }}
                                                </div>
                                            </template>
                                            <el-space wrap v-for="(item,index) in getPrizeByDateList.slice(4,5)">
                                                <div class="d-flex" v-for="i in item">
                                                    <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                                </div>
                                            </el-space>
                                        </el-card>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-card shadow="never">
                                            <template #header>
                                                <div class="card-header">
                                                    {{ dialogTopList[3][0].name }}
                                                </div>
                                            </template>
                                            <el-space wrap v-for="(item,index) in getPrizeByDateList.slice(5,6)">
                                                <div class="d-flex" v-for="i in item">
                                                    <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                                </div>
                                            </el-space>
                                        </el-card>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-card shadow="never">
                                            <template #header>
                                                <div class="card-header">
                                                    {{ dialogTopList[4][0].name }}
                                                </div>
                                            </template>
                                            <el-space wrap v-for="(item,index) in getPrizeByDateList.slice(6,7)">
                                                <div class="d-flex" v-for="i in item">
                                                    <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                                </div>
                                            </el-space>
                                        </el-card>
                                    </el-col>
                                </el-row>
                                <el-card shadow="never">
                                    <template #header>
                                        <div class="card-header">
                                            {{ dialogTopList[5][0].name }}
                                        </div>
                                    </template>
                                    <el-space wrap v-for="(item,index) in getPrizeByDateList.slice(7,8)">
                                        <div class="d-flex" v-for="i in item">
                                            <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                        </div>
                                    </el-space>
                                </el-card>
                                <el-card shadow="never">
                                    <template #header>
                                        <div class="card-header">
                                            {{ dialogTopList[6][0].name }}
                                        </div>
                                    </template>
                                    <el-space wrap v-for="(item,index) in getPrizeByDateList.slice(8,9)">
                                        <div class="d-flex" v-for="i in item">
                                            <p class="bg py-2 mx-2 my-2 rounded-2">{{ i.toString() }}</p>
                                        </div>
                                    </el-space>
                                </el-card>
                            </div>
                        </el-popover>
                    </div>
                    <p style="font-size: 12px;margin-top: .5rem" v-if="item.IsOpened">{{ item.OpenTime }}{{$t('共销售彩票张数')}}:
                        {{item.Count}},{{$t('中奖总金额')}}({{item.gold}})</p>

                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, h, shallowRef} from 'vue';
import type {ElInput, FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {WebSideMap} from '@/api/gamepb/customer';
import {AdminGameCenter} from '@/api/gamepb/admin';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import moment from "moment";

const {t} = useI18n()
const store = useStore()
const formRef = ref<FormInstance>()
const WebSideList: JsMap<WebSideMap, string> = {
    [WebSideMap.FROM_SITE_SANOOK]: t('sanook'),
    [WebSideMap.FROM_SITE_THAIRATH]: t('thairath'),
    [WebSideMap.FROM_SITE_MYHORA]: t('myhora'),
    [WebSideMap.FROM_SITE_GLO]: t('glo'),
};
let getPrizeByDateStorage = ref([])

const customPrefix = shallowRef({
    render() {
        return h('p', t('年份'))
    },
})
const pickerOptions = (time: Date) => {
    return time.getFullYear() > new Date().getFullYear()
}
const datePickerValue = ref(moment().format('YYYY'))

let getPrizeByDateList = ref(null)
const toDetailVisible = ref(false)
let OpenTime = reactive({
    List: [],
})
const dialogTopList = reactive(
    [
        [{name: t('一等奖'), content: 'รางวัลละ 6,000,000 บาท'}],
        [{name: t('前三'), content: '2 รางวัลๆละ 4,000 บาท'}, {
            name: t('后三'),
            content: '2 รางวัลๆละ 4,000 บาท'
        }, {name: t('后二'), content: '1 รางวัลๆละ 2,000 บาท'}],
        [{name: t('一等奖边奖'), content: '2 รางวัลๆละ 100,000 บาท'}],
        [{name: t('二等奖'), content: 'รางวัลที่ 2 มี 5 รางวัลๆละ 200,000 บาท'}],
        [{name: t('三等奖'), content: 'รางวัลที่ 3 มี 10 รางวัลๆละ 80,000 บาท'}],
        [{name: t('四等奖'), content: 'รางวัลที่ 4 มี 50 รางวัลๆละ 40,000 บาท'}],
        [{name: t('五等奖'), content: 'รางวัลที่ 5 มี 100 รางวัลๆละ 20,000 บาท'}]
    ]
)
let loading = ref(false)

async function callFunc(args: any) {
    // return Client.send('mq/lotter/admin/getOpenDaysByYear', {Year:datePickerValue.value})
    return Client.send(`mq/lotter/admin/getOpenDaysByYear`, args)
}
import ut from '@/lib/util'
const queryList = async () => {
    loading.value = true
    let [combine, err] = await callFunc({Year: datePickerValue.value})
    if (err) {
        return tip.e(err)
    }
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    OpenTime.List = combine.OpenTimeInfos?.map(v => {
        return {
            ...v,
            gold:ut.toNumberWithComma(v.Gold),
            type: 0,
            popOver:false
        }
    })
    console.log(OpenTime.List);
}
const OpenPrizeConfirm = (e) => {
    OpenPrizeButton(e)
}
const OpenPrizeButton = async (e) => {
    loading.value = true
    let args = {
        OpenDay: e
    }
    let [data,err] = await Client.send(`mq/lotter/admin/manualOpenPrize`, args)

    loading.value = false
    if (err) {
        return tip.e(err)
    }
    console.log(data);
    queryList()
}
const toDetail = async (i) => {
    const exists = getPrizeByDateStorage.value.find((item) => item.date === i.OpenTime);
    i.popOver = true
    if (exists) {
        getPrizeByDateList.value = exists.arr;
    } else {
        if (!i.IsPulled) {
            let args = {
                OpenDay: i.OpenTime,
                FromSite: Number(i.type)
            }
            let [data, err] = await Client.send(`mq/lotter/admin/manualPullPrize`, args)
            if (err) {
                return tip.e(err)
            }
            getPrizeByDateList.value = data.Prizes
        } else {
            let arr = {OpenDay: i.OpenTime,}
            let [data, err] = await Client.send(`mq/lotter/admin/getPrizeByDate`, arr)
            if (err) {
                return tip.e(err)
            }
            getPrizeByDateList.value = data.Prizes
        }
        getPrizeByDateStorage.value.push({
            arr: getPrizeByDateList.value,
            date: i.OpenTime
        })
    }
}
onMounted(queryList);

</script>
<style scoped lang='scss'>
.elCard {
    border: none;
    margin-bottom: 1rem;
    height: 160px;
    .el-card {
        border-radius: .5rem;
        height: 100%;
        :deep(.el-card__header){
            padding: 1rem;
        }
        :deep(.el-card__body){
            padding: 1rem;
        }
    }

    .topCard {
        //padding-left: 1.5rem;
        position: relative;

        span {
            //margin-right: .5rem;
        }

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

.dialog-openList-body {
    .el-card {
        margin-bottom: 1rem;

        .d-flex {
            padding: 0.2rem 0.4rem;
            border: 1px solid #ccc;
            border-radius: 1rem;
            min-width: 83px;
            text-align: center;
            p{
                font-size: 12px;
            }
        }
    }
}
.el-popover.el-popper{

}
</style>
