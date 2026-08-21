<template>

    <div>
        <!-- 添加弹框 -->
        <el-dialog v-model="addDialog" :title="$t('添加商户')"
                   @open-auto-focus="openDialog"
                   destroy-on-close
                   :align-center="true"
                   :width="store.viewModel === 2 ? '100%' : '950px'" @close="emits('update:modelValue')">
            <el-form ref="addFormRef" :model="addForm" label-width="140px" :rules=addRules
                     class="dialog__form">

                <el-row :gutter="18">
                    <el-col :span="12">

                        <el-form-item prop="OperatorType">
                            <template #label>
                                <div class="customLabel"> {{ $t('商户类型') }}:</div>
                            </template>
                            <el-select
                                    v-model="addForm.OperatorType"
                                    :placeholder="$t('请选择商户类型')"
                                    :disabled="store.AdminInfo.GroupId == 2"
                            >
                                <template v-for="item in props.operatorType">
                                    <el-option :label="$t(item.label)" :value="item.value" v-if="item.value != 0"/>
                                </template>

                            </el-select>
                        </el-form-item>


                        <el-form-item prop="AppID">

                            <template #label>
                                <div class="customLabel"> {{ $t('商户AppID') }}:</div>
                            </template>
                            <el-input v-model="addForm.AppID" :placeholder="$t('请输入用户名称')" maxlength="15"/>
                        </el-form-item>

                        <el-form-item prop="UserName">

                            <template #label>
                               <div class="customLabel"> {{ $t('商户账号') }}:</div>
                            </template>


                            <el-input v-model="addForm.UserName" onkeyup="value=value.replace(/[^a-zA-Z0-9_]/g, '')"
                                      maxlength="15" :placeholder="$t('请输入商户账号')"></el-input>

                        </el-form-item>

                        <el-form-item prop="CooperationType" v-if="addForm.OperatorType == 2">
                            <template #label>
                                <div class="customLabel"> {{ $t('合作模式') }}:</div>
                            </template>
                            <el-select
                                v-model="addForm.CooperationType"
                                @change="CooperationTypeChange"
                                :placeholder="$t('请选择合作模式')"

                            >
                                <template v-for="item in props.incomeType">
                                    <el-option :label="$t(item.label)" :value="item.value" />
                                </template>

                            </el-select>
                        </el-form-item>

                        <el-form-item :label="$t('平台费率') + ':'" prop="PlatformPay" class="PlatformPay" v-if="(addForm.OperatorType == 2 && addForm.CooperationType == 1) || addForm.OperatorType == 1">

                            <el-input v-model="addForm.PlatformPay" :placeholder="$t('请输入平台费率')"
                                      @blur="inputNum('PlatformPay', $event)" maxlength="4" :max="100">
                                <template #suffix>
                                    %
                                </template>
                            </el-input>
                        </el-form-item>

                        <el-form-item :label="$t('流水比例') + ':'" prop="TurnoverPay" class="PlatformPay" v-if="(addForm.OperatorType == 2 && addForm.CooperationType == 2) || addForm.OperatorType == 1">

                            <el-input v-model="addForm.TurnoverPay" :placeholder="$t('请输入流水比例')"
                                      @blur="inputNum('TurnoverPay', $event, 8)" maxlength="14" :max="100">
                                <template #suffix>
                                    %
                                </template>
                            </el-input>
                        </el-form-item>



                        <!--                    <el-form-item :label="$t('预付款金额') + ':'" prop="Advance">-->
                        <!--                        <el-input v-model="addForm.Advance" :placeholder="$t('请输入预付款金额')"-->
                        <!--                                  @blur="inputNum('Advance', $event)" maxlength="12"/>-->
                        <!--                    </el-form-item>-->

                        <el-form-item :label="$t('归属国家') + ':'">
                            <el-select
                                    v-model="addForm.BelongingCountry"
                                    :placeholder="$t('请选择归属国家')"
                                    clearable
                                    filterable
                            >
                                <template v-for="item in props.belongingCountryList">

                                    <el-option :label="$t(item.label)" v-if="item.value != 'ALL'" :value="item.value"/>
                                </template>

                            </el-select>
                        </el-form-item>



                        <el-form-item prop="CurrencyKey">
                            <template #label>
                                <div class="customLabel"> {{ $t('商户币种') }}:</div>
                            </template>
                            <el-select
                                    v-model.trim="addForm.CurrencyKey"
                                    :placeholder="$t('请选择商户币种')"
                                    filterable
                                    clearable
                            >
                                <template v-for="item in CurrencyConfig">
                                    <el-option :label="`【${item.value}】${item.label}`" :value="item.value"
                                               v-if="item.value != 0"/>
                                </template>

                            </el-select>
                        </el-form-item>
                        <el-form-item prop="CurrencyRate" v-if="addForm.OperatorType == 2 && (store.AdminInfo.GroupId <= 1)">
                            <template #label>
                                <div class="customLabel"> {{ $t('RTP区间') }}:</div>
                            </template>
                            <el-select
                                v-model="addForm.HighRTPOff"
                                :placeholder="$t('请选择RTP区间')">
                                <template v-for="(item, index) in OperatorRTP" >
                                    <el-option :label="`${item.split('-')[0]}%~${item.split('-')[1]}%`" :value="Number(index)" v-if="index != 0"/>
                                </template>

                            </el-select>
<!--                            <div class="switchContainer">-->
<!--                                <el-switch-->
<!--                                    v-model="addForm.CurrencyRate"-->
<!--                                    :active-value="1"-->
<!--                                    :inactive-value="0"-->
<!--                                />-->
<!--                            </div>-->
                        </el-form-item>
                        <el-form-item prop="CurrencyRate" v-if="addForm.OperatorType == 2">
                            <template #label>
                                <div class="customLabel"> {{ $t('基础货币单位') }}:</div>
                            </template>
                            <el-select
                                v-model="addForm.CurrencyCtrlStatus"
                                :placeholder="$t('请选择基础货币单位')"

                            >
                                <template v-for="(item, index) in props.currencyRate">

                                    <el-option :label="$t(item)" :value="Number(index)"/>
                                </template>

                            </el-select>
<!--                            <div class="switchContainer">-->
<!--                                <el-switch-->
<!--                                    v-model="addForm.CurrencyRate"-->
<!--                                    :active-value="1"-->
<!--                                    :inactive-value="0"-->
<!--                                />-->
<!--                            </div>-->
                        </el-form-item>
                        <el-form-item :label="$t('备注') + ':'" v-if="store.AdminInfo.GroupId != 2">
                            <el-input
                                v-model="addForm.Remark"
                                :rows="3"
                                type="textarea"
                                :placeholder="$t('请输入备注')"
                            />
                        </el-form-item>

                        <el-form-item :label="$t('联系方式') + ':'" class="Contact_context">
                            <el-row :gutter="5" v-for="(item, index) in addForm.Contact" :key="index"
                                    style="margin-right: 0">
                                <el-col :span="8">
                                    <el-input v-model="item.name" :placeholder="$t('联系人')"
                                              @input="inputContact(index, 'name',$event)"></el-input>
                                </el-col>
                                <el-col :span="16" style="position: relative;padding-right: 0">
                                    <el-input v-model="item.value" :placeholder="$t('联系方式')"
                                              @input="inputContact(index, 'value',$event)"></el-input>
                                    <div style="position: absolute;left: calc(100% + 10px); width: 100%;top: 0">
                                        <el-button size="small" type="info" :icon="Minus" circle
                                                   @click="delContact(index)"/>
                                        <el-button size="small" type="primary" :icon="Plus" circle
                                                   v-if="index == addForm.Contact.length - 1" @click="addContact"/>

                                    </div>
                                </el-col>

                            </el-row>
                        </el-form-item>
                    </el-col>


                    <el-col :span="12">

                        <template v-if="addForm.OperatorType == 1">
                            <el-form-item :label="$t('开户权限') + ':'" prop="WalletMode">

                                <el-radio-group v-model="addForm.ReviewType" style="margin-top: -4px">
                                    <el-radio :value="item.value" size="large" v-for="item in CREATE_MERCHANT_RULE">
                                        {{ $t(item.label) }}
                                    </el-radio>
                                </el-radio-group>
                            </el-form-item>
                        </template>
                        <template v-else>


                            <el-form-item :label="$t('钱包类型') + ':'" prop="WalletMode">

                                <el-radio-group v-model="addForm.WalletMode" style="margin-top: -4px">
                                    <el-radio :value="item.value" size="large" v-for="item in props.walletOptions">
                                        {{ $t(item.label) }}
                                    </el-radio>
                                </el-radio-group>
                            </el-form-item>


                            <el-form-item :label="$t('余额') + ':'" v-if="store.AdminInfo.GroupId <= 1">

                                <el-input

                                        v-model="addForm.Balance"
                                        :placeholder="$t('商户剩余平台费用')"
                                        disabled
                                >
                                    <template #suffix>
                                        <div style="margin-right: -10px">
                                            <el-button type="primary" plain @click="operatorBalanceVisible = true">
                                                {{ $t('充值') }}
                                            </el-button>
                                        </div>
                                    </template>
                                </el-input>
                                <div style="color: var(--el-color-danger);font-size: 12px">{{$t('在添加用户界面中点击添加才会生效')}}</div>
                            </el-form-item>

                            <operator_container v-if="store.AdminInfo.GroupId <= 1 && addForm.OperatorType == 2"
                                                :label-name="$t('线路商') + ':'"
                                                :defaultOperatorEvent="defaultOperatorEvent"
                                                @select-operatorInfo="operatorListChange" :is-init="true"
                                                :operator-type="1"></operator_container>

                            <el-form-item :label="$t('游戏厂商') + ':'" prop="Address"

                                          style="margin-bottom: 15px">
                                <el-select
                                    v-model="addForm.DefaultManufacturerOn"
                                    @change="ManufacturerChange"
                                    :placeholder="$t('请选择游戏厂商')"
                                    multiple
                                    collapse-tags
                                    min
                                >
                                    <template v-for="(item, index) in ManufacturerList">

                                        <el-option :label="$t(item.ManufacturerName)" :value="item.ManufacturerCode"/>
                                    </template>

                                </el-select>
                            </el-form-item>

                            <template v-if="store.AdminInfo.GroupId != 2">


                                <el-form-item :label="$t('服务器回调地址') + ':'" prop="Address"
                                              v-if="addForm.WalletMode == 2"
                                              style="margin-bottom: 5px">
                                    <el-input
                                            v-model="addForm.Address"
                                            :placeholder="$t('请输入服务器回调地址')"
                                    />
                                </el-form-item>
                                <div style="color: var(--el-color-danger);text-align: right;font-size: 12px;margin-top: 15px;margin-bottom: 15px; width: 85%"
                                     v-if="addForm.WalletMode == 2">{{ $t('慎重填写回调地址,填写有误将无法正常游戏.') }}
                                </div>

                                <el-form-item :label="$t('服务器IP白名单') + ':'" v-if="addForm.WalletMode == 1">
                                    <el-input
                                            v-model="addForm.UserWhite"
                                            :rows="3"
                                            type="textarea"
                                            :placeholder="$t('请输入服务器IP白名单')"
                                    />
                                </el-form-item>
                            </template>

<!--                            <el-form-item :label="$t('风控语言') + ':'">-->
<!--                                <el-select-->
<!--                                    v-model="addForm.Lang"-->
<!--                                    :placeholder="$t('请选择风控语言')"-->

<!--                                    clearable-->
<!--                                >-->
<!--                                    <template v-for="item in languageConfig">-->

<!--                                        <el-option :label="$t(item.label)" :value="item.value"/>-->
<!--                                    </template>-->

<!--                                </el-select>-->
<!--                            </el-form-item>-->
                            <template v-if="store.AdminInfo.GroupId != 2">

<!--                                <el-form-item :label="$t('机器人识别码') + ':'">-->
<!--                                    <el-input-->
<!--                                            v-model="addForm.Robot"-->
<!--                                            :placeholder="$t('请输入机器人识别码')"-->
<!--                                    >-->

<!--                                        <template #suffix>-->
<!--                                            <el-tooltip-->
<!--                                                    class="box-item"-->
<!--                                                    effect="dark"-->
<!--                                                    :content="$t('Telegram机器人识别码')"-->
<!--                                                    placement="bottom-end"-->

<!--                                            >-->
<!--                                                <el-icon size="18" color="#000000">-->
<!--                                                    <QuestionFilled/>-->
<!--                                                </el-icon>-->
<!--                                            </el-tooltip>-->
<!--                                        </template>-->
<!--                                    </el-input>-->

<!--                                </el-form-item>-->
<!--                                <el-form-item :label="$t('chatId') + ':'">-->
<!--                                    <el-input-->
<!--                                            v-model="addForm.ChatID"-->
<!--                                            :placeholder="$t('请输入chatId')"-->
<!--                                    >-->
<!--                                        <template #suffix>-->
<!--                                            <el-tooltip-->
<!--                                                    class="box-item"-->
<!--                                                    effect="dark"-->
<!--                                                    :content="$t('Telegram机器人所在频道ID')"-->
<!--                                                    placement="bottom-end"-->

<!--                                            >-->
<!--                                                <el-icon size="18" color="#000000">-->
<!--                                                    <QuestionFilled/>-->
<!--                                                </el-icon>-->
<!--                                            </el-tooltip>-->
<!--                                        </template>-->
<!--                                    </el-input>-->

<!--                                </el-form-item>-->
                                <el-form-item :label="$t('余额不足阈值') + ':'">


                                    <el-input v-model="addForm.BalanceThreshold"
                                              :placeholder="$t('请输入余额不足阈值')"/>
                                </el-form-item>
                                <el-form-item :label="$t('余额不足间隔') + ':'">
                                    <el-select v-model="addForm.BalanceThresholdInterval"
                                               :placeholder="$t('请选择余额不足间隔时间')">

                                        <el-option v-for="item in BALANCE_THRESHOLD" :label="$t(item.label)"
                                                   :value="item.value"/>

                                    </el-select>


                                </el-form-item>
                                <el-form-item :label="$t('欠费提示间隔') + ':'">
                                    <el-select v-model="addForm.ArrearsThresholdInterval"
                                               :placeholder="$t('请选择欠费提示间隔')">

                                        <el-option v-for="item in ARREARS_HINT_INTERVAL" :label="$t(item.label)"
                                                   :value="item.value"/>

                                    </el-select>
                                </el-form-item>

                                <el-row style="width: 90%" v-if="store.AdminInfo.GroupId == 1">
                                    <el-col :span="12" class="col-right">
                                        <el-form-item :label="$t('退出按钮') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.ShowExitBtnOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>
                                    </el-col>
                                    <el-col :span="12" class="col-right btn-right-content">
                                        <el-form-item :label="$t('是否显示币种') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.CurrencyVisibleOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>
                                    </el-col>
                                    <el-form-item :label="$t('退出按钮链接') + ':'" v-if="addForm.ShowExitBtnOff">
                                        <el-input
                                            v-model="addForm.ExitLink"
                                            :placeholder="$t('请输入退出按钮链接')"
                                        />

                                    </el-form-item>
                                    <el-col :span="12" class="col-right btn-right-content">

                                        <el-form-item :label="$t('RTP设置') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.RTPOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>

                                        <el-form-item :label="$t('赢取最高倍数') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.MaxMultipleOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>
                                        <el-form-item :label="$t('购买RTP') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.BuyRTPOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>

                                    </el-col>
                                    <el-col :span="12" class="col-right btn-right-content">
                                        <el-form-item :label="$t('个人RTP设置') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.PlayerRTPSettingOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>

                                        <el-form-item :label="$t('赢取最高钱数') + ':'">
                                            <div class="switchContainer">
                                                <el-switch
                                                    v-model="addForm.MaxWinPointsOff"
                                                    :active-value="1"
                                                    :inactive-value="0"
                                                />
                                            </div>
                                        </el-form-item>


                                    </el-col>

                                </el-row>

                            </template>
                        </template>
                    </el-col>


                </el-row>

            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('添加') }}</el-button>
                </span>
            </template>
        </el-dialog>

        <operatorBanlance v-model="operatorBalanceVisible" handel-type="add" :merchant-balance="0" @commitMerchantBalance="calcBalance"/>
    </div>

</template>


<script setup lang="ts">

import {computed, reactive, Ref, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import {Minus, Plus} from "@element-plus/icons-vue";
import {tip} from "@/lib/tip";
import {Client} from "@/lib/client";
import {AdminInfo, EditMaintenance} from "@/api/adminpb/info";
import {AdminGroup} from "@/api/adminpb/group";
import ut, {generatePassword, PasswordRSAEncryption} from "@/lib/util";

import {ElLoading, ElMessageBox, FormRules} from "element-plus";
import Operator_container from "./operator_container.vue";
import OperatorBanlance from "./operatorBanlance.vue";
import {ARREARS_HINT_INTERVAL, BALANCE_THRESHOLD, CREATE_MERCHANT_RULE, LANG_LIST} from "@/lib/config";
import {AddLog} from "@/lib/commRequest";
import {LOG_OPERATOR_TYPE, LOG_PRIMARY_MODULE} from "@/api/adminpb/log";
import {OperatorRTP} from "@/enum";
import {Manufacturer} from "@/api/gamepb/game";


let {t} = useI18n()

const addFormRef = ref(null)
const lineMerchantPlatformPay = ref(0)
const lineMerchantTurnoverPay = ref(0)
const PlatformPayValid = (rule: any, value: any, callback: any) => {

    if (value == 0){
        callback(new Error(t("平台费不可为0")))
    }
    if (addForm.value.OperatorType == 2){

        if (lineMerchantPlatformPay.value == 0){
            callback()
        }


        if (value >= lineMerchantPlatformPay.value){
            callback()
        }else{
            callback(new Error(t(`商户平台费不可小于线路商平台费({Num}%)`, {Num: lineMerchantPlatformPay.value})))
        }
    }else{
        callback()
    }
}
const TurnoverValid = (rule: any, value: any, callback: any) => {

    if (value == 0){
        callback(new Error(t("流水比例不可为0")))
    }
    if (addForm.value.OperatorType == 2){

        if (lineMerchantTurnoverPay.value == 0){
            callback()
        }


        if (value >= lineMerchantTurnoverPay.value){
            callback()
        }else{
            callback(new Error(t(`商户流水比例不可小于线路商流水比例({Num}%)`, {Num: lineMerchantTurnoverPay.value})))
        }
    }else{
        callback()
    }
}
const isIPv4validator = (rule: any, value: any, callback: any) => {

    if (!value) {
        callback()
        return
    }

    let whitelistString = value.split(',').filter(item => item !== "")
    const isIPv4 = new RegExp(/^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})$/);
    let whiteBools = whitelistString.filter(t => isIPv4.test(t));
    if (whiteBools.length) {
        callback()
    } else {
        callback(new Error(t('IP不合法')))
    }
}


const validContact = (rule: any, value: any, callback: any) => {

    let Contact = addForm.value["Contact"].filter(item => item.name != "" && item.value != "")

    if (Contact.length == 0) {
        callback(new Error(t('联系方式至少有一条填写联系人与联系方式')))
    } else {
        callback()
    }

}

const isIPv4validator2 = (rule: any, value: any, callback: any) => {

    if (!value) {
        callback()
        return
    }
    let whiteBools = false
    const isIPv4 = new RegExp(/^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})$/);
    //url= 协议://(ftp的登录信息)[IP|域名](:端口号)(/或?请求参数)
    var strRegex = '^((https|http|ftp)://)?'//(https或http或ftp):// 可有可无
        + '(([\\w_!~*\'()\\.&=+$%-]+: )?[\\w_!~*\'()\\.&=+$%-]+@)?' //ftp的user@  可有可无
        + '(([0-9]{1,3}\\.){3}[0-9]{1,3}' // IP形式的URL- 3位数字.3位数字.3位数字.3位数字
        + '|' // 允许IP和DOMAIN（域名）
        + '(localhost)|'	//匹配localhost
        + '([\\w_!~*\'()-]+\\.)*' // 域名- 至少一个[英文或数字_!~*\'()-]加上.
        + '\\w+\\.' // 一级域名 -英文或数字  加上.
        + '[a-zA-Z]{1,6})' // 顶级域名- 1-6位英文
        + '(:[0-9]{1,5})?' // 端口- :80 ,1-5位数字
        + '((/?)|' // url无参数结尾 - 斜杆或这没有
        + '(/[\\w_!~*\'()\\.;?:@&=+$,%#-]+)+/?)$';//请求参数结尾- 英文或数字和[]内的各种字符

    var strRegex1 = '^(?=^.{3,255}$)((http|https|ftp)?:\/\/)?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/)?(?:\/(.+)\/?$)?(\/\w+\.\w+)*([\?&]\w+=\w*|[\u4e00-\u9fa5]+)*$';
    var re = new RegExp(strRegex, 'i');//i不区分大小写

    whiteBools = re.test(value) || isIPv4.test(value)
    if (whiteBools) {
        callback()
    } else {
        callback(new Error(t('服务器回调不合法')))
    }
}

const store = useStore()


const props = defineProps([
    "modelValue",
    "operatorType",
    "operatorStatus",
    "incomeType",
    "walletOptions",
    "belongingCountryList",
    "currencyRate"
])
const emits = defineEmits(['update:modelValue', "addOperator"])

const addDialog = computed(() => {

    return Boolean(props.modelValue)
})

const defaultOperatorEvent = ref({})
const operatorListChange = (value) => {
    if (value) {

        lineMerchantPlatformPay.value = value.PlatformPay
        lineMerchantTurnoverPay.value = value.TurnoverPay


        addFormRef.value.validateField("PlatformPay")
        addFormRef.value.validateField("TurnoverPay")
        addForm.value.LineMerchantID = value.AppID
    } else {

        addForm.value.LineMerchantID = ""
    }
}

const addForm: Ref<EditMaintenance> = ref(<EditMaintenance>{
    // 商户编码
    UserName: "",                                               //  商户类型
    OperatorType: 2,                                            //  商户类型
    Remark: "",                                                 //  商户类型
    BelongingCountry: "",                                       //  商户类型
    LineMerchantID: "",                                         //  商户类型
    Password: null,                                             //  商户类型
    TurnoverPay: 0,                                             //  商户类型
    Name: "",                                                   //  商户主账号
    AppID: "",                                                  //  商户名称
    PlatformPay: 1,                                             //  平台费
    BuyRTPOff:0,
    CooperationType: props.incomeType[0].value,                 //  合作模式
    Advance: 1,                                                 //  预付款金额
    CurrencyKey: "",                                            //  商户币种
    CurrencyCtrlStatus: 0,
    CurrencyKeyName: "",                                        //  商户币种
    Contact: [{name: "", value: ""}],                           //  联系方式
    WalletMode: 2,                                              //  钱包类型
    DefaultManufacturerOn:null,
    Surname: "",                                                //  会员前缀
    Lang: "en",                                                   //  客户端默认语言
    ServiceIp: "",                                              //  服务器地址
    WhiteIps: "",                                               //  服务器IP白名单
    Address: "",                                                //  服务器回调地址
    UserWhite: "",                                              //  用户白名单
    LoginOff: 1,                                                //  服务器IP白名单
    FreeOff: 1,                                                 //  服务器回调地址
    DormancyOff: 1,                                             //  用户白名单
    RestoreOff: 1,                                              //  记录开启
    ManualFullScreenOff: 1,                                     //  免游游戏开关板
    NewGameDefaulOff: 1,                                        //  防休眠开启
    MassageOff: 1,                                              //  断复开关
    MassageIp: "",                                              //  手动全屏开启
    RTPOff: 1,                                                  //  RTP设置
    StopLoss: 1,                                                //  消息推送
    PlayerRTPSettingOff: 1,                                     //  消息推送
    MaxMultipleOff: 1,                                          //  消息推送地址
    Balance: 0,                                                 //  消息推送地址
    ReviewType: CREATE_MERCHANT_RULE[0].value,                  //  消息推送地址
    BalanceThreshold: 1500,
    BalanceThresholdInterval: BALANCE_THRESHOLD[0].value,
    ArrearsThresholdInterval: ARREARS_HINT_INTERVAL[0].value,
    _LineMerchant: 1,                                           //  游戏新型预设
    Robot: "",                                                  //  商户类型
    ChatID: "",                                                 //  商户类型
    ShowNameAndTimeOff:0,
    ShowExitBtnOff:1,
    ExitLink:"http://",
    MaxWinPointsOff:1,
    HighRTPOff:1,
    PGConfig:{
        CarouselOff: 1,
        ShowNameAndTimeOff: 1,
        CurrencyVisibleOff: 1,
        StopLoss: 0,
        ExitBtnOff: 0,
        ExitLink: "",
        OfficialVerify: 1,
    },
    CurrencyVisibleOff:1,
    JILIConfig:{
        BackPackOff: 1,
        CurrencyVisibleOff: 1,
        OpenScreenOff:1,
        SidebarOff:0,
    },
    TADAConfig:{
        BackPackOff: 1,
        CurrencyVisibleOff: 1,
        OpenScreenOff:1,
        SidebarOff:0,
    },
    PPConfig: {
        ShowNameAndTimeOff: 1,
        CurrencyVisibleOff: 1,
    }
})

const operatorBalanceData = ref(null)

const languageConfig = ref(LANG_LIST)
const CurrencyConfig = ref([])
const tableData = ref([])
const ManufacturerList = ref([])
const Count = ref(0)


const operatorBalanceVisible = ref(false)


const openDialog = () => {
    addForm.value = {
        // 商户编码
        UserName: "",                       //  商户类型
        BelongingCountry: "",               //  归属国家
        Password: null,               //  归属国家
        OperatorType: 2,                    //  商户类型
        LineMerchantID: "",                    //  商户类型
        Remark: "",                         //  商户类型
        Name: "",                            //  商户主账号
        AppID: "",                           //  商户名称
        PlatformPay: 0,                      //  平台费
        TurnoverPay: 0,                      //  平台费
        BuyRTPOff: 0,                      //  平台费
        CooperationType: props.incomeType[0].value,                 //  合作模式
        DefaultManufacturerOn:null,
        Advance: 0,                          //  预付款金额
        CurrencyKey: "",                     //  商户币种
        CurrencyKeyName: "",                 //  商户币种
        Contact: [{name: "", value: ""}],     //  联系方式
        WalletMode: 2,                       //  钱包类型
        Surname: "",                         //  会员前缀
        Lang: "en",                            //  客户端默认语言
        ServiceIp: "",                       //  服务器地址
        WhiteIps: "",                        //  服务器IP白名单
        Address: "",                         //  服务器回调地址
        UserWhite: "",                       //  用户白名单
        LoginOff: 1,                        //  服务器IP白名单
        FreeOff: 1,                         //  服务器回调地址
        DormancyOff: 1,                     //  用户白名单
        RestoreOff: 1,                      //  记录开启
        ManualFullScreenOff: 1,             //  免游游戏开关板
        NewGameDefaulOff: 1,                //  防休眠开启
        MassageOff: 1,                      //  断复开关
        MassageIp: "",                      //  手动全屏开启
        RTPOff: 1,                          //  RTP设置
        StopLoss: 1,                        //  消息推送
        PlayerRTPSettingOff: 1,              //  消息推送
        MaxMultipleOff: 1,                  //  消息推送地址
        Robot: "",                         //  商户类型
        ChatID: "",                         //  商户类型
        ReviewStatus: 0,                         //  商户类型
        Balance: 0,
        ReviewType: CREATE_MERCHANT_RULE[0].value,
        BalanceThreshold: 1500,
        BalanceThresholdInterval: BALANCE_THRESHOLD[0].value,
        ArrearsThresholdInterval: ARREARS_HINT_INTERVAL[0].value,
        _LineMerchant: 1,                   //  游戏新型预设
        ShowNameAndTimeOff:0,
        ShowExitBtnOff:0,
        ExitLink:"http://",
        MaxWinPointsOff:1,
        HighRTPOff:1,
        CurrencyVisibleOff:1,
        CurrencyCtrlStatus: 0,
        NextRate: 0,
        Status: 0,
        PGConfig:{
            CarouselOff: 1,
            ShowNameAndTimeOff: 1,
            CurrencyVisibleOff: 1,
            StopLoss: 0,
            ExitBtnOff: 1,
            ExitLink: "",
            OfficialVerify: 1,
        },
        JILIConfig:{
            BackPackOff: 1,
            CurrencyVisibleOff: 1,
            OpenScreenOff:1,
            SidebarOff:1,
        },
        TADAConfig:{
            BackPackOff: 1,
            CurrencyVisibleOff: 1,
            OpenScreenOff:1,
            SidebarOff:1,
        },
        PPConfig: {
            ShowNameAndTimeOff: 1,
            CurrencyVisibleOff: 1,
        },
    }
    if (store.AdminInfo.GroupId == 2) {
        addForm.value.OperatorType = 2
        addForm.value.ReviewType = store.AdminInfo.Businesses.ReviewType
    }
    lineMerchantPlatformPay.value = store.AdminInfo.Businesses.PlatformPay
    lineMerchantTurnoverPay.value = store.AdminInfo.Businesses.TurnoverPay

    if (store.language != 'zh'){

        let langList = []
        for (const i in languageConfig.value) {
            langList.push({label:languageConfig.value[i].value, value:languageConfig.value[i].value})
        }

        languageConfig.value = langList
    }

    addForm.value.DefaultManufacturerOn = ["ALL"]
}

const CooperationTypeChange = (value) => {
        addForm.value.HighRTPOff = value
}

const addRules = reactive<FormRules>({
    OperatorType: [
        {
            validator: (rule: any, value: any, callback: any) => {
                if (value == ""){
                    callback(new Error(t('商户类型不能为空')))
                }else{
                    callback()
                }
            },
            trigger: 'blur'
        }
    ],
    UserName: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('商户账号不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    AppID: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('商户AppID不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    PlatformPay: [{required: true, message: t('平台费不能为空'), trigger: 'blur'}, {
        validator: PlatformPayValid,
        trigger: 'blur'
    }],
    TurnoverPay: [{required: true, message: t('流水比例不能为空'), trigger: 'blur'},{
        validator: TurnoverValid,
        trigger: 'blur'
    }],
    CooperationType: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('请选择合作模式')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],

    Advance: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('预付金额不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    CurrencyKey: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('币种不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    WalletMode: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('钱包类型不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    Lang: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('语言不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    WhiteIps: [{
        validator: (rule: any, value: any, callback: any) => {
            if (value == ""){
                callback(new Error(t('服务器地址不能为空')))
            }else{
                callback()
            }
        },
        trigger: 'blur'
    }],
    Address: [{validator: isIPv4validator2, trigger: 'blur'}],
    Contact: [{required: true, trigger: 'blur', message: t('联系方式至少有一条')}, {
        validator: validContact,
        trigger: 'blur'
    }],
    UserWhite: [{required: true, trigger: 'blur', message: t('白名单地址不能为空')}, {
        validator: isIPv4validator,
        trigger: 'blur'
    }],

})


const getCurrency = async () => {
    const [res, err] = await Client.Do(AdminGroup.GetCurrency, {} as any)
    if (err) {
        tip.e(t("未找到币种配置"))
        return
    }

    if (res.List instanceof Array && res.List.length > 0) {
        CurrencyConfig.value = res.List.map(item => ({
            label: item.CurrencyName,
            value: item.CurrencyCode,
            CurrencySymbol: item.CurrencySymbol,
        }))

    }


}
const getManufacturerList = async () => {
    const [res, err] = await Client.Do(Manufacturer.GetManufacturerList, {} as any)
    res.List.unshift({
        "Id": "1",
        "ManufacturerName": "全部",
        "ManufacturerCode": "ALL"
    })
    ManufacturerList.value = res.List
    addForm.value.DefaultManufacturerOn = ["ALL"]


}


getCurrency()
getManufacturerList()


const AddAdminer = async () => {

    addFormRef.value.validate(async valid => {
        if (valid) {
            //
            let Contact = addForm.value["Contact"].filter(item => item.name != "" && item.value != "")
            //
            // if (Contact.length == 0) {
            //     tip.e(t("请填写至少一个联系方式"))
            //     return
            // }


            const params = {
                ...addForm.value,
                Contact
            }


            if (typeof params["TurnoverPay"] == "string") {

                params["TurnoverPay"] = parseFloat(params["TurnoverPay"])
            }
            if (typeof params["PlatformPay"] == "string") {

                params["PlatformPay"] = parseFloat(params["PlatformPay"])
            }
            if (typeof params["Advance"] == "string") {

                params["Advance"] = parseFloat(params["Advance"])

            }
            params["WhiteIps"] = params["WhiteIps"].split(",")


            params["Contact"] = JSON.stringify(Contact)
            //
            // if (params["Robot"] == "" && params ["ChatID"] != "") {
            //
            //     tip.e(t("chatId填写机器码不能为空"))
            //     return
            // }
            // if (params["Robot"] != "" && params ["ChatID"] == "") {
            //     tip.e(t("机器码填写chatId不能为空"))
            //     return
            // }


            // 线路商的钱包类型改为0

            if (params.OperatorType == 1) {
                params.WalletMode = 0
                params.BalanceThresholdInterval = 0
                params.ArrearsThresholdInterval = 0
            }
            params.Password = generatePassword(8)
            params.Balance = Number(params.Balance)
            params.BalanceThreshold = Number(params.BalanceThreshold)
            params.BalanceThresholdInterval = Number(params.BalanceThresholdInterval)
            params.ArrearsThresholdInterval = Number(params.ArrearsThresholdInterval)

            params.PGConfig = {
                CarouselOff: 1,
                ShowNameAndTimeOff: 1,
                CurrencyVisibleOff: 1,
                StopLoss: 1,
                ExitBtnOff: 0,
                ExitLink: "http://",
                OfficialVerify: 1,
            }
            params.JILIConfig = {
                BackPackOff: 1,
                CurrencyVisibleOff: 1,
                OpenScreenOff: 1,
                SidebarOff: 1,
            }
            params.TADAConfig = {
               BackPackOff: 1,
               CurrencyVisibleOff: 1,
               OpenScreenOff: 1,
                SidebarOff: 1,
            }
            params.PPConfig =  {
                ShowNameAndTimeOff: 1,
                CurrencyVisibleOff: 1,
            }

            if (params.DefaultManufacturerOn.indexOf("ALL") != -1){
                params.DefaultManufacturerOn = null
            }

            ElMessageBox.confirm(
                t('确认添加商户'),
                t('是否确认'),
                {
                    confirmButtonText: t('确定'),
                    cancelButtonText: t('关闭'),
                    type: 'warning',
                }
            )
                .then(async () => {
                    const loading = ElLoading.service({
                        lock: true,
                        text: 'Loading',
                        background: 'rgba(0, 0, 0, 0.7)',
                    })
                    const [res, err] = await Client.Do(AdminInfo.AddMaintenance, params)
                    loading.close()
                    if (err != null) {
                        tip.e(t(err))
                        return
                    }



                    if (operatorBalanceData.value){

                        AddLog(LOG_PRIMARY_MODULE.OPERATOR_LIST, LOG_OPERATOR_TYPE.BALANCE_EDIT, null, params, params.AppID, operatorBalanceData.value)
                    }




                    emits('update:modelValue')
                    emits('addOperator', {data:params})

                })
        }
    })


}

const ManufacturerChange = (value) => {
    if (value.length == 0){
        addForm.value.DefaultManufacturerOn = ["ALL"]
        return
    }
    if (value[value.length - 1] == "ALL"){
        addForm.value.DefaultManufacturerOn = ["ALL"]
    }else{
        addForm.value.DefaultManufacturerOn = value.filter(item => item != "ALL")
    }

}

const inputNum = (key, value, num = 2) => {
    value = addForm.value[key]

    let regexp = /[^\d\.\d{0, 2}]/g

    if (key == 'TurnoverPay'){
        regexp = /[^\d\.\d{0, 8}]/g
    }

    value = value.replace(regexp, '')
    if (value <= 0) {
        value = 1
    }
    addForm.value[key] = value
}

const inputContact = (index, key, event) => {

    addForm.value.Contact[index][key] = event
}

const addContact = () => {
    let findNull = addForm.value.Contact.find(item => item.name == "" || item.value == "")
    if (findNull) {

        tip.e(t("有未填写的联系方式,请补全后添加"))
        return
    }
    addForm.value.Contact.push({
        name: "",
        value: ""
    })

}

const delContact = (index) => {

    if (addForm.value.Contact.length == 1) {
        return
    }

    addForm.value.Contact.splice(index, 1)
}

const calcBalance = (value) => {
    operatorBalanceData.value = value.data
    addForm.value.Balance = Number(value.handledBalance)
}
</script>


<style lang="scss">
.dialog__form {

  .el-form-item {


    width: 100%;
  }

  .el-select {
    width: 100%;

  }

  .el-form-item__content > * {
    width: 80%;

  }

  //.Contact_context .el-form-item__content > *{
  //    width: 100%;
  //}
}

.customLabel:before{
    color: var(--el-color-danger); content: "*"; margin-right: 4px;
}
.PlatformPay.is-error{
    margin-bottom: 40px!important;
}
</style>
