<template>
    <div class="searchView ">
        <el-form
            :model="uiData"
            style="max-width: 100%"
        >
            <el-space wrap >
                <operator_container :defaultOperatorEvent ="defaultOperatorEvent" @select-operatorInfo="operatorListChange" :is-init="true" :hase-all="haseAll"></operator_container>
                <el-input v-if="props.currencyVisible" v-model="CurrencyCode"></el-input>
                <!--                <el-form-item label="货币">
                                    <el-autocomplete
                                        v-model="uiData.Pid"
                                        :fetch-suggestions="querySearch"
                                        clearable
                                        class="inline-input w-50"
                                        :placeholder="$t('货币')"
                                    />
                                </el-form-item>-->
                <el-form-item :label="$t('日期查询')" v-if="!hiddenTime">
                    <el-date-picker locale="zh-cn" v-model="uiData.times" type="daterange" unlink-panels :range-separator="$t('至')"
                                    :disabled-date="option"
                                    :start-placeholder="$t('开始时间')"
                                    :end-placeholder="$t('结束时间')"
                                    :shortcuts="shortcuts"/>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="queryList">{{ $t('搜索') }}</el-button>
                </el-form-item>
            </el-space>
        </el-form>

    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, defineProps, defineEmits } from 'vue';
import { useStore } from '@/pinia/index';
import { useI18n } from 'vue-i18n';
import Operator_container from "@/components/operator_container.vue";
import moment from 'moment';
import Currency_container from "@/components/currency_container.vue";
import {debug} from "util";
// 获取当前日期
const defaultOperatorEvent = ref({})
const defaultCurrencyEvent = ref({})
const { t } = useI18n();
const store = useStore();
let loading = ref(false);
const props = defineProps({
    uiData: Object as () => Record<string, any>,// 传递表单数据的属性
    currencyVisible: Boolean as () => Boolean,// 传递表单数据的属性
    hiddenTime: Boolean as () => Boolean,// 传递表单数据的属性
    haseAll: Boolean as () => Boolean
});
const emit  = defineEmits();
const uiData = ref(props.uiData);

const option = (value) => {


    return value > new Date()
}

interface RestaurantItem {
    name: string
    value: string
}
const restaurants = ref<RestaurantItem[]>([])

const querySearch = (queryString: string, cb: any) => {
    console.log(restaurants,'restaurants');
    const results = queryString
        ? restaurants.value.filter(createFilter(queryString))
        : restaurants.value
    // call callback function to return suggestions
    cb(results)
}
const createFilter = (queryString: string) => {
    return (restaurant: RestaurantItem) => {
        console.log(restaurant.value,'restaurant.abbreviation');
        return (
            restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
        )
    }
}
const currency =()=> {
    return [
        {value:'CNY',name:'人民币'},
        {value:'USD',name:'美元'},
        {value:'AED',name:'阿联酋迪拉姆'},
        {value:'AFN',name:'阿富汗尼'},
        {value:'ALL',name:'列克'},
        {value:'AMD',name:'亚美尼亚德拉姆'},
        {value:'ANG',name:'荷属安的列斯盾'},
        {value:'AOA',name:'宽扎'},
        {value:'ARS',name:'阿根廷披索'},
        {value:'AUD',name:'澳元'},
        {value:'AWG',name:'阿鲁巴弗罗林'},
        {value:'AZN',name:'阿塞拜疆马纳特'},
        {value:'BAM',name:'自由兑换马克'},
        {value:'BBD',name:'巴巴多斯元'},
        {value:'BDT',name:'塔卡'},
        {value:'BGN',name:'保加利亚列弗'},
        {value:'BHD',name:'巴林第纳尔'},
        {value:'BIF',name:'布隆迪法郎'},
        {value:'BMD',name:'百慕大元'},
        {value:'BND',name:'文莱元'},
        {value:'BOB',name:'玻利维亚诺'},
        {value:'BRL',name:'巴西雷亚尔'},
        {value:'BS',vname:'强势玻利瓦尔'},
        {value:'BSD',name:'巴哈马元'},
        {value:'BSS',name:'主权玻利瓦尔'},
        {value:'BTN',name:'努扎姆'},
        {value:'BWP',name:'普拉'},
        {value:'BYN',name:'白俄罗斯卢布'},
        {value:'BZD',name:'伯利兹元'},
        {value:'CAD',name:'加拿大元'},
        {value:'CDF',name:'刚果法郎'},
        {value:'CHE',name:'WIR欧元'},
        {value:'CHF',name:'瑞士法郎'},
        {value:'CHW',name:'WIR法郎'},
        {value:'CLP',name:'智利比索'},
        {value:'COP',name:'哥伦比亚比索'},
        {value:'CRC',name:'哥斯达黎加科朗'},
        {value:'CUC',name:'可兑换比索'},
        {value:'CUP',name:'古巴比索'},
        {value:'CVE',name:'佛得角埃斯库多'},
        {value:'CZK',name:'捷克克朗'},
        {value:'DJF',name:'吉布提法郎'},
        {value:'DKK',name:'丹麦克朗'},
        {value:'DOP',name:'多明尼加比索'},
        {value:'DZD',name:'阿尔及利亚第纳尔'},
        {value:'EGP',name:'埃及镑'},
        {value:'ERN',name:'纳克法'},
        {value:'ETB',name:'埃塞俄比亚比尔'},
        {value:'EUR',name:'欧元'},
        {value:'FJD',name:'斐济元'},
        {value:'FKP',name:'福克兰镑'},
        {value:'GBP',name:'英镑'},
        {value:'GEL',name:'拉里'},
        {value:'GHS',name:'加纳塞地'},
        {value:'GIP',name:'直布罗陀镑'},
        {value:'GMD',name:'冈比亚达拉西'},
        {value:'GNF',name:'几内亚法郎'},
        {value:'GTQ',name:'格查尔'},
        {value:'GYD',name:'圭亚那元'},
        {value:'HKD',name:'港币'},
        {value:'HNL',name:'伦皮拉'},
        {value:'HRK',name:'库纳'},
        {value:'HTG',name:'古德'},
        {value:'HUF',name:'福林'},
        {value:'IDR',name:'印尼卢比'},
        {value:'ILS',name:'以色列新锡克尔'},
        {value:'INR',name:'印度卢比'},
        {value:'IQD',name:'伊拉克第纳尔'},
        {value:'IRR',name:'伊朗里亚尔'},
        {value:'ISK',name:'冰岛克郎'},
        {value:'JMD',name:'牙买加元'},
        {value:'JOD',name:'约旦第纳尔'},
        {value:'JPY',name:'日元'},
        {value:'KES',name:'肯尼亚先令'},
        {value:'KGS',name:'吉尔吉斯斯坦索姆'},
        {value:'KMF',name:'科摩罗法郎'},
        {value:'KPW',name:'朝鲜元'},
        {value:'KRW',name:'韩元'},
        {value:'KWD',name:'科威特第纳尔'},
        {value:'KYD',name:'开曼群岛元'},
        {value:'KZT',name:'哈萨克坚戈'},
        {value:'LAK',name:'老挝基普'},
        {value:'LBP',name:'黎巴嫩镑'},
        {value:'LKR',name:'斯里兰卡卢比'},
        {value:'LRD',name:'利比里亚元'},
        {value:'LSL',name:'莱索托洛蒂'},
        {value:'LYD',name:'利比亚第纳尔'},
        {value:'MAD',name:'摩洛哥道拉姆'},
        {value:'MDL',name:'摩尔多瓦列伊'},
        {value:'MGA',name:'马达加斯加阿里亚里'},
        {value:'MKD',name:'马其顿第纳尔'},
        {value:'MMK',name:'缅甸元'},
        {value:'MNT',name:'图格里克'},
        {value:'MOP',name:'澳门元'},
        {value:'MRU',name:'毛里塔尼亚乌吉亚'},
        {value:'MUR',name:'毛里求斯卢比'},
        {value:'MVR',name:'马尔代夫拉菲亚'},
        {value:'MWK',name:'马拉维克瓦查'},
        {value:'MXN',name:'墨西哥比索'},
        {value:'MXV',name:'墨西哥(资金)'},
        {value:'MYR',name:'马来西亚林吉特'},
        {value:'MZN',name:'梅蒂卡尔'},
        {value:'NAD',name:'纳米比亚元'},
        {value:'NGN',name:'奈拉'},
        {value:'NIO',name:'尼加拉瓜科多巴'},
        {value:'NOK',name:'挪威克朗'},
        {value:'NPR',name:'尼泊尔卢比'},
        {value:'NZD',name:'新西兰元'},
        {value:'OMR',name:'阿曼里亚尔'},
        {value:'PAB',name:'巴拿马巴波亚'},
        {value:'PEN',name:'秘鲁索尔'},
        {value:'PGK',name:'巴布亚新几内亚基那'},
        {value:'PHP',name:'菲律宾比索'},
        {value:'PKR',name:'巴基斯坦卢比'},
        {value:'PLN',name:'波兰兹罗提'},
        {value:'PYG',name:'巴拉圭瓜拉尼'},
        {value:'QAR',name:'卡塔尔里亚尔'},
        {value:'RON',name:'罗马尼亚列伊'},
        {value:'RSD',name:'塞尔维亚第纳尔'},
        {value:'RUB',name:'俄罗斯卢布'},
        {value:'RWF',name:'卢旺达法郎'},
        {value:'SAR',name:'沙特里亚尔'},
        {value:'SBD',name:'所罗门群岛元'},
        {value:'SCR',name:'塞舌尔卢比'},
        {value:'SDG',name:'苏丹镑'},
        {value:'SEK',name:'瑞典克朗'},
        {value:'SGD',name:'新加坡元'},
        {value:'SHP',name:'圣赫勒拿镑'},
        {value:'SLL',name:'塞拉利昂利昂'},
        {value:'SOS',name:'索马里先令'},
        {value:'SRD',name:'苏里南元'},
        {value:'SSP',name:'南苏丹镑'},
        {value:'STN',name:'圣多美多布拉'},
        {value:'SVC',name:'萨尔瓦多科朗'},
        {value:'SYP',name:'叙利亚镑'},
        {value:'SZL',name:'斯威士兰里兰吉尼'},
        {value:'THB',name:'泰铢'},
        {value:'TJS',name:'塔吉克斯坦索莫尼'},
        {value:'TMT',name:'土库曼斯坦马纳特'},
        {value:'TND',name:'突尼斯第纳尔'},
        {value:'TOP',name:'汤加潘加'},
        {value:'TRY',name:'新土耳其里拉'},
        {value:'TTD',name:'特立尼达和多巴哥元'},
        {value:'TWD',name:'新台币'},
        {value:'TZS',name:'坦桑尼亚先令'},
        {value:'UAH',name:'乌克兰赫里纳'},
        {value:'UGX',name:'乌干达先令'},
        {value:'UYI',name:'乌拉圭比索（乌拉圭）'},
        {value:'UYU',name:'乌拉圭比索'},
        {value:'UZS',name:'乌兹别克斯坦苏姆'},
        {value:'VES',name:'委内瑞拉玻利瓦尔'},
        {value:'VND',name:'越南盾'},
        {value:'VUV',name:'瓦努阿图瓦图'},
        {value:'WST',name:'萨摩亚塔拉'},
        {value:'XAF',name:'中洲金融共同体法郎'},
        {value:'XAG',name:'银'},
        {value:'XAU',name:'黄金'},
        {value:'XBA',name:'债券市场单位欧洲综合单位（EURCO）'},
        {value:'XBB',name:'债券市场单位欧洲货币单位（E.M.U.-6）'},
        {value:'XBC',name:'债券市场单位欧洲账户单位9（E.U.A.-9）'},
        {value:'XBD',name:'债券市场单位欧洲账户单位17（E.U.A.-17）'},
        {value:'XCD',name:'东加勒比元'},
        {value:'XDR',name:'国际货币基金组织特别提款权'},
        {value:'XOF',name:'西非金融共同体法郎'},
        {value:'XPD',name:'钯金'},
        {value:'XPF',name:'太平洋法郎'},
        {value:'XPT',name:'铂金'},
        {value:'XUA',name:'亚行账户单位'},
        {value:'YER',name:'也门里亚尔'},
        {value:'ZAR',name:'南非兰特'},
        {value:'ZMW',name:'赞比亚克瓦查'},
        {value:'ZWD',name:'津巴布韦元'},
        {value:'ZWL',name:'津巴布韦元'},
    ]
}
// 获取当前日期
const currentDate = moment();
// 获取上一个月的月初
const lastMonthStart = moment(currentDate).startOf('month').subtract(1, 'month');
// 获取上一个月的月末
const lastMonthEnd = moment(lastMonthStart).endOf('month');
const shortcuts = [
    {
        text: t('今天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime())
            return [start, end]
        },
    },
    {
        text: t('昨天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24)
            return [start, end]
        },
    },
    {
        text: t('近7天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            return [start, end]
        },
    },
    {
        text: t('近30天'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            return [start, end]
        },
    },
    {
        text: t('本月'),
        value: () => {
            const end = new Date()
            const start = new Date()
            start.setDate(1)
            return [start, end]
        },
    },
    {
        text: t('上月'),
        value: () => {
            return [lastMonthStart, lastMonthEnd]
        },
    },
]
const queryList = () => {
    // 触发查询操作，可以通过 emit 发送事件到父组件
    emit('query-list');
};
const CurrencyCode = ref("")
const operatorListChange = (value) =>{

    if (value){

        uiData.value.rate = value.PlantRate
        uiData.value.Pid = value.Id
        CurrencyCode.value = value.CurrencyKey
    }else{
        uiData.value.rate = 1
        uiData.value.Pid = 0
        CurrencyCode.value = ""
    }

}
const currencyListChange = (value) =>{
    uiData.value.CurrencyCode = value.CurrencyCode
}
onMounted(() => {
    // restaurants.value = currency()
    console.log(uiData.value,'----312312----');

})
</script>
<style scoped lang='scss'>


</style>
