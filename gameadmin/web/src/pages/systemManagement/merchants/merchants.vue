
<template>
    <div class="page_table_context">

      <el-row :gutter="20">



          <el-col :span="12">
              <el-descriptions
                      direction="horizontal"
                      size="large"
                      :column="4"
                      border
              >

                  <el-descriptions-item
                      :span="4"
                      :label="$t('状态')">
                      <el-switch
                          v-model="editForm.Status"
                          :active-value="1"
                          :inactive-value="0"
                      />
                  </el-descriptions-item>

                  <el-descriptions-item
                      :span="4"
                      :label="$t('合作模式')">

                      <div v-if="editForm.CooperationType" style="font-size: 14px;text-align: right">{{ incomeType.find(item => item.value == editForm.CooperationType ).label}}，
                          {{ $t('本月费率') }}{{(editForm.PlatformPay * 100).toFixed(2)}}%，{{ $t('下个月费率') }}{{(editForm.NextRate * 100).toFixed(2)}}
                      </div>
                  </el-descriptions-item>

                  <el-descriptions-item :label="$t('钱包类型')" :span="4">
                      <div style="font-size: 14px;text-align: right">  {{ editForm.WalletMode == 2 ? $t('单一钱包') : $t('转账钱包') }}</div>

                  </el-descriptions-item>

                  <el-descriptions-item :label="$t('商户币种')" :span="4">
                      <el-select
                          v-model.trim="editForm.CurrencyKey"
                          :placeholder="$t('请选择商户币种')"
                          disabled
                          clearable
                      >
                          <template v-for="item in CurrencyConfig">
                              <el-option :label="item.label" :value="item.value" v-if="item.value != 0"/>
                          </template>

                      </el-select>
                  </el-descriptions-item>

                  <el-descriptions-item :label="$t('客户端默认语言')" :span="4">
                      <el-select
                          v-model="editForm.Lang"
                          :placeholder="$t('请选择商户类型')"

                          clearable
                      >
                          <template v-for="item in languageConfig">

                              <el-option  :label="$t(item.lanName)" :value="item.abbreviation"/>
                          </template>

                      </el-select>
                  </el-descriptions-item>


                  <el-descriptions-item :label="$t('服务器地址') + ':'" :span="4">
                      <el-input
                          v-model="editForm.WhiteIps"
                          :placeholder="$t('请输入服务器地址')"
                      />
                  </el-descriptions-item>

                  <el-descriptions-item :label="$t('服务器回调地址') + ':'" :span="4">
                      <el-input
                          v-model="editForm.Address"
                          :placeholder="$t('请输入服务器回调地址')"
                      />
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('RTP设置') + ':'" :span="4">
                      <el-switch
                          v-model="editForm.RTPOff"
                          :active-value="0"
                          :inactive-value="1"
                      />
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('止损止赢开关') + ':'" :span="4">
                      <el-switch
                          v-model="editForm.StopLoss"
                          :active-value="1"
                          :inactive-value="0"
                      />

                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('赢取最高押注倍数') + ':'" :span="4">
                      <el-switch
                          v-model="editForm.MaxMultipleOff"
                          :active-value="1"
                          :inactive-value="0"
                      />
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('联系方式') + ':'" :span="4">
                          <el-row :gutter="24" v-for="(item, index) in editForm.Contact" :key="index">
                              <el-col :span="9">
                                  <el-input v-model="item.name" :placeholder="$t('请输入联系人的姓名')" @input="inputContact(index, 'name',$event)"></el-input>
                              </el-col>
                              <el-col :span="9">
                                  <el-input v-model="item.value" :placeholder="$t('请输入联系方式')" @input="inputContact(index, 'value',$event)"></el-input>
                              </el-col>
                              <el-col :span="6">
                                  <el-button size="small" type="info" :icon="Minus" circle @click="delContact(index)"/>
                                  <el-button size="small" type="primary" :icon="Plus" circle v-if="index == editForm.Contact.length - 1" @click="addContact"/>

                              </el-col>
                          </el-row>
                  </el-descriptions-item>

              </el-descriptions>
          </el-col>

          <el-col :span="12">

              <el-form ref="addFormRef" :model="gameSearch" label-width="120px" :inline="true" class="dialog__form">

                  <el-form-item :label="$t('游戏分类') + ':'">
                      <el-radio-group v-model="activeName">
                          <el-radio-button  size="large" value="SLOT">SLOT</el-radio-button>
                      </el-radio-group>

                  </el-form-item>
                  <el-form-item :label="$t('游戏搜索') + ':'">

                      <el-input v-model="gameSearch.GameName" style="width: 240px" @focus="getGameList" placeholder="请输入游戏名称" />
                      <el-select
                          v-model="gameSearch.GameOn"
                          placeholder="选择游戏状态"
                          style="width: 240px;margin-left: 15px"
                          @change="getGameList"
                      >
                          <el-option
                              v-for="(item, key) in gameState"
                              :key="key"
                              :label="item"
                              :value="parseInt(key)"
                          />
                      </el-select>



                  </el-form-item>
              </el-form>


              <customTable :table-header="tableHeader" :table-data="tableData" :page="gameSearch.PageIndex" :page-size="gameSearch.PageSize" :count="Count" height="500px">

                  <template #StopLoss="scope">
                      <el-button size="small" type="success" v-if="scope.scope.StopLoss == 1">开启</el-button>
                      <el-button size="small" type="danger"  v-if="scope.scope.StopLoss == 0">关闭</el-button>

                  </template>
                  <template #operator="scope">
                      <el-button type="primary" size="small" plain @click="openGameSetting(scope.scope)">游戏设置</el-button>
                      <el-button type="success"  v-if="scope.scope.StopLoss == 1" size="small" plain @click="editGameStatus(scope.scope)">开启</el-button>
                      <el-button type="danger"  v-if="scope.scope.StopLoss == 0" size="small" plain @click="editGameStatus(scope.scope)">关闭</el-button>
                      <!--                <el-button type="warning" size="small" plain>试玩</el-button>-->

                  </template>

              </customTable>



          </el-col>

      </el-row>

    </div>

</template>

<script setup lang="ts">
import {ref} from "vue";
import type {Ref} from "vue";
import {AdminInfo, EditMaintenance} from "@/api/adminpb/info";
import {Client} from "@/lib/client";
import {AdminGroup} from "@/api/adminpb/group";
import customTable from "@/components/customTable/tableComponent.vue"
import {tip} from "@/lib/tip";
import ut from "@/lib/util";
import {useI18n} from "vue-i18n";
import {Merchant, MerchantGameReq} from "@/api/adminpb/merchant";
import {useStore} from "@/pinia";
import {Minus, Plus} from "@element-plus/icons-vue";

let {t} = useI18n()
const tableHeader = [
    { label: "游戏编号", value: "GameId", width: "180px"},
    { label: "游戏名称", value: "GameName", width: "180px"},
    { label: "RTP", value: "RTP", width: "80px"},
    { label: "止盈止损开关", value: "StopLoss", type: "custom", width: "120px"},
    { label: "赢取最高押注倍数", value: "MaxMultiple", width: "120px"},
    { label: "投注金额", value: "BetBase",width: "180px"},
    { label: "操作", value: "operator", type: "custom", fixed:"right",  width: "180px", hiddenVisible: true},

]
const { AdminInfo } = useStore()
const Count = ref(0)
let tableData = ref([])
let CurrencyConfig = ref([])
let languageConfig = ref([])
let incomeType = [
    {label: t("收益分成"), value: 1},
]

const gameState = ref({
    "-1": "全部",
    "1": "开启",
    "0": "关闭",
})
let activeName = "SLOT"

const gameSearch:Ref<MerchantGameReq> = ref<MerchantGameReq>({

    Gametype    : 0,
    GameName    : "",
    GameOn      : -1,
    UserName    : "",
    PageIndex   : 1,
    PageSize    : 20,
})
const editForm:Ref<EditMaintenance> = ref(<EditMaintenance>{
    NextRate:0,
    Status:0,
    RTPOff:0,

    // 商户编码
    OperatorType: 0,                   //  商户类型
    UserName:"",                            //  商户主账号
    Name:"",                            //  商户主账号
    AppID:"",                           //  商户名称
    PlatformPay:0,                     //  平台费
    CooperationType:"",                 //  合作模式
    Advance:0,                         //  预付款金额
    CurrencyKey:"",                     //  商户币种
    CurrencyKeyName:"",                     //  商户币种
    Contact: [{name:"", value:""}],     //  联系方式
    WalletMode:2,                       //  钱包类型
    Surname:"",                         //  会员前缀
    Lang:"",                            //  客户端默认语言
    ServiceIp:"",                       //  服务器地址
    WhiteIps:"",                        //  服务器IP白名单
    Address:"",                         //  服务器回调地址
    UserWhite:"",                       //  用户白名单
    LoginOff: 1,                        //  服务器IP白名单
    FreeOff: 1,                         //  服务器回调地址
    DormancyOff: 1,                     //  用户白名单
    RestoreOff: 1,                      //  记录开启
    ManualFullScreenOff: 1,             //  免游游戏开关板
    NewGameDefaulOff: 1,                //  防休眠开启
    MassageOff: 1,                      //  断复开关
    MassageIp: "",                      //  手动全屏开启
    StopLoss: 1,                        //  消息推送
    MaxMultipleOff: 1,                  //  消息推送地址
    _LineMerchant: 1,                   //  游戏新型预设
})

const initData = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, {})

    if (err) {
        return tip.e(err)
    }

    if (data.List){

        editForm.value = data.List[0]
    }
}

const getCurrency = async () =>{
    const [res, err] = await Client.Do(AdminGroup.GetCurrency, {} as any)
    if (err) {
        tip.e(t("未找到币种配置"))
        return
    }

    if(res.List instanceof Array && res.List.length > 0){
        CurrencyConfig.value = res.List.map(item => ({
            label:item.CurrencyName,
            value:item.CurrencyCode,
            CurrencySymbol:item.CurrencySymbol,
        }))

    }



}

const getLanguage = async () =>{
    const [res, err] = await Client.Do(AdminGroup.GetLang, {} as any)
    if (err) {
        tip.e(t("未找到语言配置"))
        return
    }

    if (res.List[0]){
        for (const listElementKey in res.List[0]) {

            let langItem = ut.LangList.find(list => list.abbreviation.toUpperCase() === listElementKey.toString())

            if (langItem){
                languageConfig.value.push(langItem)
            }

        }
    }


}

const getGameList = async () => {

    // 商户的Id
    gameSearch.value.UserName = AdminInfo.Username

    const [data, err] = await Client.Do(Merchant.GetMerchantGames, gameSearch.value)

    tableData.value = data.List
    Count.value = data.CountAll
}
initData()
getCurrency()
getLanguage()


const editGameStatus = (scope) => {

}



const RegexpFloat = (e) => {
    e.target.value = (e.target.value.match(/^\d*(\.?\d{0,2})/g)[0]) || null

}

const openGameSetting = (data) => {

}

const inputContact = (index, key, value) => {

}
const delContact = () => {

}
const addContact = () => {

}
const addForm = () => {

}

</script>

<style scoped lang="scss">
.dialog__form {
    .el-form-item{
        width: 100%;
        margin-bottom: 20px;
    }

    .el-select{
        width: 100%;

    }
    .el-form-item__content > *{
        width: 80%;

    }

}
.gameType__container{
    width: 100%;

}
.gameType__item{
    width: auto;
    height: auto;
    padding: 8px 10px;
    border: 1px solid #cccccc59;
}

</style>
