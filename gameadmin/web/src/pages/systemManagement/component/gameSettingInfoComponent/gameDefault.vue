
<template>

  <div>
      <div style="width: 100%; height: 30px;display: flex;font-weight: bolder;line-height: 30px">
          {{ $t("默认投注金额") }}
      </div>
      <div style="width: 100%; height: auto;display: flex;font-weight: bolder;line-height: 30px">


          <template v-if="gameInfo.GameManufacturer == 'HACKSAW'">
              {{ $t("设置玩家进入游戏时的初始默认投注金额") }}
          </template>

          <template v-else>
              {{ $t("设置玩家进入游戏时的初始默认投注金额：投注金额=投注额*投注倍数*游戏线数") }}
          </template>

      </div>

      <div>

          <el-form label-width="120px">
              <el-row :gutter="18">
                  <el-col :span="12" v-if="gameInfo.GameManufacturer == 'HACKSAW'">
                      <el-form-item :label="`${$t('线注额度')}:`">
                          <el-select style="width: 150px" v-model="gameInfo.DefaultCs">
                              <el-option v-for="item in BetList" :value="item.BetValue" :key="item.BetValue" :label="item.BetValue"></el-option>
                          </el-select>
                      </el-form-item>
                  </el-col>
                  <el-col :span="12">
                      <el-form-item :label="`${$t('投注倍数')}:`" v-if="gameInfo.GameManufacturer == 'PG'">
                          <el-select style="width: 150px" v-model="gameInfo.DefaultBetLevel">
                              <el-option v-for="item in BetList[0].BetInfo" :value="item.multipleValue" :key="item.multipleValue" :label="item.multipleValue"></el-option>
                          </el-select>
                      </el-form-item>
                  </el-col>
                  <el-form-item>
                      <div style="width: 100%; height: auto;display: flex">
                          <template v-if="gameInfo.GameManufacturer == 'PG'">
                              {{ `${$t('线注额度')}(${gameInfo.DefaultCs})x${$t("投注倍数")}(${gameInfo.DefaultBetLevel})x${$t('游戏线数')}(${linNum})=${$t('总投注')}${gameInfo.totalBet}` }}

                          </template>

                      </div>
                  </el-form-item>


              </el-row>
          </el-form>

      </div>
  </div>
</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {computed, onMounted, ref, watch} from "vue";
import {GeneratorDefaultCs} from "@/pages/systemManagement/component/gameSettingInfoComponent/generatorDefaultCS";

const props = defineProps(["merchantInfo", "gameInfo", "tableData"])
const store = useStore()
const {t}=useI18n()
const linNum = ref()

let BetList = ref(props.tableData)

watch(props.tableData, (newData)=> {
    BetList.value = newData
})




const onlineupNum = ref(0)
const onlinedownNum = ref(0)

onMounted(()=>{

    linNum.value = store.LineMap[props.gameInfo.GameId]
    GeneratorDefaultCs(onlineupNum.value, onlinedownNum.value)
})

</script>


<style scoped lang="scss">

</style>
