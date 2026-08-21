<template>
    <el-form v-model="PlayRTP" label-position="right" label-width="140px">
        <el-form-item :label="$t('游戏设置')">

            <el-select v-model="selectGameList" :disabled="props.PlayerBuyRTPDisabled && props.PlayerRTPDisabled" @change="selectGameChange"
                       style="width: 280px;" :placeholder="$t('请选择游戏')"
                       multiple filterable collapse-tags clearable>
                <el-option v-for="(item, index) in games" :key="index" :label="$t(item.Name)"
                           :value="item.ID">

                    <div style="display: flex;align-items: center;justify-content: space-between">
                    <span style="float: left">{{ $t(item.Name)}}</span>
                    <el-tag
                        v-if="item.ID != 'ALL'"
                        type="primary"
                        size="small"
                        style="float: right">

                        {{ $t(item.ManufacturerName) }}
                      </el-tag>
                    </div>
                </el-option>
            </el-select>

            <div class="RTPControlGameList">

                <el-tag v-for="(item, index) in selectGameListData" :closable="props.PlayerRTPDisabled"
                        @close="removeGame(index)" type="info"
                        style="margin-right: 10px">{{ $t(item.Name)}}
                </el-tag>
            </div>

        </el-form-item>


        <el-form-item :label="$t('账号RTP')">
            <el-space>
                <el-select :disabled="props.PlayerRTPDisabled" v-model="PlayRTP.ContrllRTP" style="width: 150px" >

                    <template v-for="(item, index) in RTPList">
                        <el-option
                            :label="Number(item) == 93 ? `93%${$t('(常用)')}` : `${item}%`"
                            :value="Number(item)"/>
                    </template>


                </el-select>

            </el-space>
        </el-form-item>

        <el-form-item :label="$t('自动解除控制')">

            <el-select v-model="PlayRTP.AutoRemoveRTP" :disabled="props.PlayerRTPDisabled" style="width: 150px">
                <el-option v-for="(item, index) in AutoRemoveRTPList" :key="index" :label="item == 0 ? $t('无') : item + '%'"
                           :value="Number(item)"></el-option>
            </el-select>



        </el-form-item>
        <div style="color: red; line-height: 20px;font-size: 13px;text-indent: 10px">
            <span style="display: block;text-indent: 10px" >1.{{ $t('账号RTP＞自动解除控制,为增加用户RTP.(RTP在区间内解除)') }}</span>
            <span style="display: block;text-indent: 10px" >2.{{ $t('账号RTP＜自动解除控制,为降低用户RTP.(RTP在区间内解除)') }}</span>
            <span style="display: block;text-indent: 10px" >3.{{ $t('账号RTP = 自动解除控制.（当RTP = 目标值时才会解除）') }}</span>
        </div>

        <el-form-item :label="$t('当前回报率')" v-if="PlayerInfo && PlayerInfo.huibao" style="visibility: hidden">
            <el-space>
                {{ percentFormatter(0, 0, PlayerInfo.huibao) }}
            </el-space>
        </el-form-item>

        <el-divider />

        <el-form-item :label="$t('购买免费游戏RTP')" v-if="props.PlayerBuyRTPDisabled">
            <el-space>
                <el-select v-model="PlayRTP.BuyRTP" style="width: 150px" >

                    <template v-for="(item, index) in BuyRTPList">
                        <el-option
                            :label="item == 0 ? $t('无') : `${item}%`"
                            :value="Number(item)"/>
                    </template>
                </el-select>
            </el-space>
        </el-form-item>

        <el-form-item :label="$t('个人赢取最高钱数')" style="margin-bottom: 0">
            <el-input v-model.number="PlayRTP.PersonWinMaxScore"
                       style="width: 150px"
                      @blur="MaxWinPointsInput"
                      :placeholder="$t('请输入')" maxlength="7">
            </el-input>
        </el-form-item>

        <el-text style="color:red;text-align: left;width: 90%;display: block;margin-bottom: 20px;text-indent: 10px;font-size: 13px">

            {{ $t('赢取最高钱数可设置区间{Num}', {Num:"1-1,000,000"})}}
        </el-text>

        <el-form-item :label="$t('个人赢取最高倍数')" style="margin-bottom: 0">
            <el-space>
                <el-input v-model.number="PlayRTP.PersonWinMaxMult"
                          style="width: 150px"
                          @blur="maxnultiple"
                          :placeholder="$t('请输入')" maxlength="5">
                </el-input>

            </el-space>
        </el-form-item>
        <el-text style="color:red;text-align: left;width: 90%;display: block;margin-bottom: 20px;text-indent: 10px;font-size: 13px">
            {{ $t('赢取最高倍数可设置{Num}', { Num:"30~10000" }) }}
        </el-text>



    </el-form>
</template>

<script setup lang="ts">

import {RTP_VALUE} from "@/lib/RTP_config";
import {computed, ref, Ref, watch, watchEffect} from "vue";
import {AdminPlayer, PlayerRTP} from "@/api/adminpb/adminPlayer";
import {Client} from "@/lib/client";
import {Manufacturer} from "@/api/gamepb/game";
import {AdminGameCenter} from "@/api/gamepb/admin";
import {useStore} from "@/pinia";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";
import {OperatorRTP} from "@/enum";
import {debug} from "util";

const props = defineProps(["SelectPlayInfo", "GameList", "marginBottom", "PlayerRTPDisabled", "PlayerBuyRTPDisabled", "highRTP"])
const selectManu = ref("")
const ManufacturerList: Ref<Manufacturer[]> = ref<Manufacturer[]>([])
const games = ref([])
const {t} = useI18n()
const store = useStore()

const selectGameList = ref([])
const selectGameListData = ref([])


const BuyRTPList = import.meta.env.VITE_BUY_RTP.split(",")
const RTPList = ref([...RTP_VALUE])
const AutoRemoveRTPList = ref([0,...RTP_VALUE])





const PlayerInfo = computed(() => {

    if (props.SelectPlayInfo && props.SelectPlayInfo.gameList && props.SelectPlayInfo.gameListIds) {


        selectGameListData.value = props.SelectPlayInfo.gameList
        selectGameList.value = props.SelectPlayInfo.gameListIds
        PlayRTP.value.ContrllRTP = Number(props.SelectPlayInfo.ContrllRTP)
        PlayRTP.value.AutoRemoveRTP = Number(props.SelectPlayInfo.AutoRemoveRTP)
        PlayRTP.value.BuyRTP = Number(props.SelectPlayInfo.BuyRTP)
        PlayRTP.value.PersonWinMaxScore = Number(props.SelectPlayInfo.PersonWinMaxScore) || 1000000
        PlayRTP.value.PersonWinMaxMult = Number(props.SelectPlayInfo.PersonWinMaxMult) || 100

    }



    return props.SelectPlayInfo
})


const PlayRTP: Ref<PlayerRTP> = ref<PlayerRTP>({

    GameID: "",  //游戏列表
    Pid: 0,  //玩家ID
    AppID: "",  //所属商户
    ContrllRTP: 93, //控制RTP
    BuyRTP: 90, //控制RTP
    AutoRemoveRTP: Number(AutoRemoveRTPList.value[0]), //自动解除RTP
    PersonWinMaxScore: 1000000, //自动解除RTP
    PersonWinMaxMult: 10000 //自动解除RTP

})

if (store.AdminInfo.GroupId == 3 && store.AdminInfo.Businesses.HighRTPOff > 0) {
    let minRTP = OperatorRTP[store.AdminInfo.Businesses.HighRTPOff].split('-')[0]
    let maxRTP = OperatorRTP[store.AdminInfo.Businesses.HighRTPOff].split('-')[1]

    RTPList.value = []
    for (let i in RTP_VALUE) {
        // 去除当前用户可设置的区间值  与  原设定值
        let flag = Number(minRTP) <= Number(RTP_VALUE[i])
            && Number(maxRTP) >= Number(RTP_VALUE[i])
            || PlayRTP.value == Number(RTP_VALUE[i])

        if (flag) {
            RTPList.value.push(Number(RTP_VALUE[i]))
        }
    }
}




const initManufacturerList = async () => {
    const [res, err] = await Client.Do(Manufacturer.GetManufacturerList, {} as any)
    ManufacturerList.value = res.List

}


const MaxWinPointsInput = (value) => {


    PlayRTP.value.PersonWinMaxScore = Number(PlayRTP.value.PersonWinMaxScore.toString().replace(/^(0+)|[^\d]+/g,''))

    if (PlayRTP.value.PersonWinMaxScore < 1) {
        PlayRTP.value.PersonWinMaxScore = 1
    }
    if (PlayRTP.value.PersonWinMaxScore > 1000000) {
        PlayRTP.value.PersonWinMaxScore = 1000000
    }

}



const maxnultiple = (value) => {
    PlayRTP.value.PersonWinMaxMult = PlayRTP.value.PersonWinMaxMult.toString().replace(/^(0+)|[^\d]+/g,'')
    if (PlayRTP.value.PersonWinMaxMult < 30) {
        PlayRTP.value.PersonWinMaxMult = 30
    }
    if (PlayRTP.value.PersonWinMaxMult > 10000) {
        PlayRTP.value.PersonWinMaxMult = 10000
    }

}


const initGameList = async () => {
    const [gameList, err] = await Client.Do(AdminGameCenter.GameList, {} as any)


    let currentManuData = []
    let gameManuData  = []


    // console.log(PlayerInfo.value.Operator.DefaultManufacturerOn)


    if(PlayerInfo.value){
        gameList.List = gameList.List.filter(item=>  !PlayerInfo.value.Operator.DefaultManufacturerOn || PlayerInfo.value.Operator.DefaultManufacturerOn.indexOf(item.ManufacturerName) > -1)
    }


    games.value =  gameList.List.sort((a,b)=> a.Id - b.Id)


    if (!selectManu.value) {
        games.value.unshift({
            Name: "全部",
            ID: "ALL"
        })
    }

}

const init = async () => {
    selectManu.value = ""
    selectGameList.value = []
    selectGameListData.value = []
    PlayRTP.value = {
        GameID: "",  //游戏列表
        Pid: 0,  //玩家ID
        AppID: "",  //所属商户
        ContrllRTP: 93, //控制RTP
        BuyRTP: 90, //控制RTP
        AutoRemoveRTP: Number(AutoRemoveRTPList.value[0]), //自动解除RTP
        PersonWinMaxScore: 1000000, //自动解除RTP
        PersonWinMaxMult: 10000
         //自动解除RTP
    }
    await initManufacturerList()
    await initGameList()
}


init()


const selectGameChange = (item) => {

    if (item.indexOf("ALL") != -1) {
        selectGameListData.value = [{
            Name: "全部",
            ID: "ALL",

        }]
        selectGameList.value = ["ALL"]

    } else {
        selectGameListData.value = games.value.filter(game => item.indexOf(game.ID) > -1)

    }


}
const manufactureChange = () => {
    initGameList()
}
const removeGame = (index) => {
    selectGameListData.value.splice(index, 1)
    selectGameList.value.splice(index, 1)
}

const RTPConfig = async (RTPData) => {

    let appId = ""

    let commit = RTPData.map(item => {
        appId = item.AppID

        return `${item.AppID}:${item.Pid}`
    })
    const ids = selectGameList.value.map(item => item)


    const PlayerRTPData: PlayerRTP = <PlayerRTP>{
        GameID: ids.join(","),  //游戏列表
        AppIDandPlayerID: commit.join(","),
        AppID:appId,
        ContrllRTP: Number(PlayRTP.value.ContrllRTP), //控制RTP
        BuyRTP: Number(PlayRTP.value.BuyRTP), //控制RTP
        AutoRemoveRTP: Number(PlayRTP.value.AutoRemoveRTP), //自动解除RTP
        PersonWinMaxScore: Number(PlayRTP.value.PersonWinMaxScore), //自动解除RTP
        PersonWinMaxMult: Number(PlayRTP.value.PersonWinMaxMult) //自动解除RTP
    }


    if (PlayerRTPData.GameID == "") {
        tip.e(t("游戏不能为空"))
        return false
    }

    if (PlayerRTPData.AppIDandPlayerID == "") {
        tip.e(t("玩家不能为空"))
        return false
    }


    const [res, err] = await Client.Do(AdminPlayer.EditPlayerRTP, PlayerRTPData)


    if (err) {

        tip.e(t(err))
        return false

    }

    return true

}

defineExpose({
    RTPConfig,
    init,
    selectManu,
    selectGameList,
})

</script>


<style scoped lang="scss">
.RTPControlGameList {
    width:100%;
  max-height: 150px;
  overflow-y: auto;
}
</style>
