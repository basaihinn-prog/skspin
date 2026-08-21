<template>
    <component :is="ReviewComponent" :PlayResp="data.PlayResp" :DiaLogGame="data.DiaLogGame" v-if="data.PlayResp !== null"></component>
    <div v-else-if="data.ID !== props.oneLog.ID">Loading...</div>
    <div v-else>{{ props.oneLog.Comment }}</div>
</template>

<script setup lang="tsx" >
import { tip } from "@/lib/tip";
import { onMounted, watch, ref ,reactive} from "vue";
import { DocBetLog } from '@/api/gamepb/customer';
import { Client } from "@/lib/client";

import BaoZang_Review from "./BaoZang_Review.vue";
import NiuBi_ReviewVue from "./NiuBi_Review.vue";
import Roma_ReviewVue from "./Roma_Review.vue";
import RomaX_ReviewVue from "./RomaX_Review.vue";
import ZhaoCaiMao_ReviewVue from "./ZhaoCaiMao_Review.vue";
import XingYunXiang_ReviewVue from "./XingYunXiang_Review.vue";
import YingCaiShen_ReviewVue from "./YingCaiShen_ReviewVue.vue";
import MaJiang_Review from "./MaJiang_Review.vue";
import MaJiang2_Review from "./MaJiang_Review2.vue";
import TuZi_Review from "./TuZi_Review.vue";
import Hilo_Review from "./Hilo_Review.vue";
import JinNiu_Review from "./JinNiu_Review.vue";
import Bonaza_Review from "./Bonaza_Review.vue";
import Olympus_Review from "./Olympus_Review.vue";
import Olympus1000_Review from "./Olympus1000_Review.vue";
import Starlight1000_Review from "./Starlight1000_Review.vue";
import Starlight_Review from "./Starlight_Review.vue";
import StarlightChristmas_Review from "./StarlightChristmas_Review.vue";
import lottery_Review from "./lottery_Review.vue";
import CowBoy_Review from "./CowBoy_Review.vue";
import SugarRush_Review from "./SugarRush_Review.vue";
import UltimateStriker_Review from "./UltimateStriker_Review.vue";

interface Prop {
    oneLog: DocBetLog;
}

const props = withDefaults(defineProps<Prop>(), {});
const data = reactive({
    PlayResp: null,
    DiaLogGame:null,
    ID: "",
});

const valMap = {
    NiuBi: NiuBi_ReviewVue,
    ZhaoCaiMao: ZhaoCaiMao_ReviewVue,
    XingYunXiang: XingYunXiang_ReviewVue,
    BaoZang: BaoZang_Review,
    Roma: Roma_ReviewVue,
    RomaX: RomaX_ReviewVue,
    YingCaiShen: YingCaiShen_ReviewVue,
    MaJiang: MaJiang_Review,
    MaJiang2: MaJiang2_Review,
    TuZi: TuZi_Review,
    Hilo: Hilo_Review,
    JinNiu: JinNiu_Review,
    Bonaza: Bonaza_Review,
    Olympus: Olympus_Review,
    Olympus1000: Olympus1000_Review,
    Starlight1000: Starlight1000_Review,
    Starlight: Starlight_Review,
    StarlightChristmas: StarlightChristmas_Review,
    lottery: lottery_Review,
    CowBoy: CowBoy_Review,
    SugarRush: SugarRush_Review,
    pg_1489936: UltimateStriker_Review,
};

const ReviewComponent = ref(null);

onMounted(() => {
    loadReviewComponent();
    refreshData();
});

watch(() => props.oneLog.ID, (newID, oldID) => {
    const gameIDParts = props.oneLog.GameID.split('_');

    if (gameIDParts.length > 1 && gameIDParts[0] === 'pg') return false;
    if (newID !== oldID) {
        data.PlayResp = null;
        refreshData();
    }
    loadReviewComponent()
});

const loadReviewComponent = () => {
    ReviewComponent.value = valMap[props.oneLog.GameID] || (props.oneLog.GameID && props.oneLog.GameID.startsWith('lottery')?valMap['lottery']:'');
};
async function refreshData() {
    data.PlayResp = null
    let url = `AdminInfo/GetSpinDetails`;
    const [resp, err] = await Client.send(url, { BetID: props.oneLog.ID,AppID:props.oneLog.AppID });
    if (err) {
        tip.e(err);
        return;
    }
    data.DiaLogGame = props.oneLog.GameID
    data.PlayResp =  resp
}
</script>
<style scoped>
.el-dropdown-menu{
    background-color: #30303c;
}
.el-dropdown-menu__item:not(.is-disabled):focus{
    background-color: #000000;
}
</style>
