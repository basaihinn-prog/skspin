<template>

  <div>

      <div class="searchView">
          <el-form :model="RuleSearchForm" :inline="true" style="max-width: 100%">

              <el-form-item :label="$t('规则名称') + ':'">
                  <el-input v-model="RuleSearchForm.RuleName" style="width: 140px" :placeholder="$t('请输入规则名称')" />
              </el-form-item>

              <el-form-item :label="$t('状态') + ':'">
                  <el-select
                          v-model="RuleSearchForm"
                          :placeholder="$t('请选择状态')"
                          style="width: 140px"
                  >
                      <el-option
                              v-for="item in RuleStatus"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                      />
                  </el-select>
              </el-form-item>

              <el-form-item>
                  <el-button style="margin-right: 15px" type="primary">搜索</el-button>
                  <el-button style="margin-right: 15px">重置</el-button>
              </el-form-item>

          </el-form>
      </div>



      <div class="page_table_context">


          <div class="flex_child_end flex">

              <el-button style="margin-bottom: 15px" type="primary" @click="rulesDialog = true">新增</el-button>
          </div>




          <customTable :tableHeader="tableHeader" :tableData="tableData">
              <template #Status="scope">
                  <el-text :type="scope.scope.Status == 1 ? 'success' : 'danger'">{{ scope.scope.Status == 1 ? $t("启用") : $t("禁用") }}</el-text>
              </template>
              <template #Operator="scope">
                  <el-button style="margin-bottom: 15px" plain type="primary">编辑</el-button>
                  <el-button style="margin-bottom: 15px" :type="scope.scope.Status == 0 ? 'success' : 'danger'">{{ scope.scope.Status == 0 ? $t("启用") : $t("禁用") }}</el-button>
                  <el-button style="margin-bottom: 15px" plain type="primary">删除</el-button>
              </template>
          </customTable>
      </div>




      <!--    维护面板    -->
      <el-dialog v-model="rulesDialog" :title="$t('添加规则集')" :width="store.viewModel === 2 ? '85%' : '60%'" @close="closeDialog">


          <div style="margin: 15px auto">{{ $t('基础配置') }}</div>

          <el-form :model="BasicsForm" :inline="true" :rules="BasicsFormRule" ref="formEl">
              <el-form-item prop="UserIp">
                  <el-input v-model="BasicsForm.RuleAggregateName" style="width: 240px" :placeholder="$t('请输入用户IP')" />
              </el-form-item>
              <el-form-item prop="TimeRange">

                  <el-date-picker
                      v-model="BasicsForm.TimeRange"
                      style="width: 300px"
                      type="daterange"
                      :range-separator="$t('至')"
                      :start-placeholder="$t('开始时间')"
                      :end-placeholder="$t('结束时间')"
                      format="YYYY-MM-DD HH:mm:ss"
                  />
              </el-form-item>

          </el-form>

          <div style="padding-top: 15px;border-top: 1px solid rgb(217,217,217)">

              <div class="flex">

                  <el-button style="margin-bottom: 15px" type="primary" @click="rulesDialog = true">添加规则</el-button>
              </div>

              <div style="width: 98%; height: 300px;overflow-x: hidden;overflow-y: auto">

                  <template v-for="(item, index) in RulesList">
                      <el-form :model="item" :inline="true" :rules="BasicsFormRule">

                          <el-row>
                              <el-form-item prop="RuleName">
                                  <el-input v-model="item.RuleName" style="width: 240px"
                                            :placeholder="$t('请输入用户IP')"/>
                              </el-form-item>

                              <el-form-item prop="RuleType">

                                  <el-radio-group v-model="item.RuleType">
                                      <el-radio value="1" size="large">游戏</el-radio>

                                      <el-radio value="2" size="large">玩家</el-radio>
                                  </el-radio-group>
                              </el-form-item>
                              <el-form-item>
                                  <el-button type="danger" plain>删除</el-button>
                              </el-form-item>
                          </el-row>

                          <!--           如果                   -->

                          <el-row>

                              <el-col :span="6">
                                  <el-tag type="info" size="large">如果</el-tag>
                              </el-col>
                              <el-col :span="12">2</el-col>
                              <el-col :span="6">3</el-col>
                          </el-row>

                          <!--           那么                   -->

                          <el-row style="margin-top: 15px">
                              <el-col :span="6">
                                  <el-tag type="info" size="large">那么</el-tag>
                              </el-col>
                              <el-col :span="12">2</el-col>
                              <el-col :span="6">3</el-col>

                          </el-row>

                      </el-form>
                  </template>


                  </div>

          </div>






          <template #footer>
              <div class="dialog-footer">
                  <el-button @click="closeDialog" style="margin-right: 15px">{{ $t("关闭") }}</el-button>
                  <el-button type="primary" @click="commitWhitList">
                      {{ $t("添加") }}
                  </el-button>
              </div>
          </template>
      </el-dialog>



  </div>

</template>

<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {Ref, ref, watch} from "vue";

import customTable from "@/components/customTable/tableComponent.vue"
import ut from "@/lib/util";
import {Rule, RuleAggregate} from "@/api/adminpb/riskManage";
import {useStore} from "@/pinia";
import {Client} from "@/lib/client";
import {AdminPlayer, PlayerReq} from "@/api/adminpb/adminPlayer";
const store = useStore()

const { t } = useI18n()

const rulesDialog = ref(false)

watch(rulesDialog, (newData) => {
    if (!newData){
        BasicsForm.value = {
            RuleAggregateName: "",
            EffectTime: "",
            InvalidTime: "",
            Status: 0,
            Rules: "",
            TimeRange: []
        }
    }
})

const RuleStatus = ref([
    {label: "全部",value: -1},
    {label: "启用", value: 1},
    {label: "禁用", value: 0},
])

const tableData = ref([

])


const BasicsForm:Ref<RuleAggregate> = ref<RuleAggregate>({
    RuleAggregateName: "",
    EffectTime: "",
    InvalidTime: "",
    Status: 0,
    Rules: "",
    TimeRange: []
})

const BasicsFormRule = ref({
    RuleAggregateName: [{required: true, message: t('规则集名称不能为空'), trigger: 'blur'},],
})



const RulesList:Ref<Rule[]> = ref<Rule[]>([
    {
        // 规则名称
        RuleName: "",
        // 规则类型
        RuleType: 1,
        // 规则条件
        Condition: "",
        // 规则结果
        Result:"",
        // 规则顺序
        Sort:1,
    }
])

const tableHeader = ref([
    {label: "编号", value: "Date", type: "selection"},
    {label: "规则集合名称", value: "RuleAggregateName"},
    {label: "创建时间", value: "CreatedTime", format: (row) => ut.fmtDateSecond(row.CreatedTime) },
    {label: "生效时间", value: "EffectTime", format: (row) => ut.fmtDateSecond(row.EffectTime) },
    {label: "失效时间", value: "InvalidTime", format: (row) => ut.fmtDateSecond(row.InvalidTime) },
    {label: "状态", value: "Status", type: "custom"},
    {label: "操作者", value: "Operator", type: "custom"},
])
const RuleSearchForm = ref({
    RuleName: "",
    Status: -1
})

const request:PlayerReq = <PlayerReq>{
    Pid: 100001,
}

const initRuleAggregate = () => {
    Client.Do(AdminPlayer.GetPlayerInfo, request)
}

initRuleAggregate()
const closeDialog = () => {
    initRuleAggregate()
}
</script>


<style scoped lang="scss">

</style>
