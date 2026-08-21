<template>

    <div>
  <!-- 添加弹框 -->
    <el-dialog v-model="addDialog"
               :title="$t('基础信息')"
               destroy-on-close
               @open="initEditData"
               :align-center="true"
               :width="store.viewModel === 2 ? '100%' : '950px'" @close="emits('update:modelValue')">
        <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="130px" :inline="true"
                 class="dialog__form">

            <el-row :gutter="18">
                <el-col :span="12">

                    <el-form-item :label="$t('状态') + ':'" v-if="store.AdminInfo.GroupId == 1">
                        <div class="switchContainer">
                            <el-switch
                                    v-model="editForm.Status"
                                    :active-value="1"
                                    :inactive-value="0"
                            />
                        </div>
                    </el-form-item>


                    <el-form-item :label="$t('合作模式') + ':'" v-if="editForm.OperatorType == 2">

                        <div style="font-size: 14px;text-align: right">
                            {{ editForm.CooperationTypeText }}

                        </div>
                    </el-form-item>

                    <el-form-item :label="$t('费率') + ':'"  v-if="editForm.OperatorType == 1 || editForm.CooperationType == 1">

                        <div style="font-size: 14px;text-align: right">

                            <el-input style="width: 150px" v-model.trim="editForm.PlatformPay"
                                      :disabled="store.AdminInfo.GroupId > 1"
                                      onkeyup="value=value.replace(/[^\d\.\d{0,2}]/g, '')" maxlength="4"/>

                        </div>

                    </el-form-item>
                    <el-form-item :label="$t('流水比例') + ':'"  v-if="editForm.OperatorType == 1 || editForm.CooperationType == 2">
                        <div style="font-size: 14px;text-align: right">
                        <el-input style="width: 150px" v-model.trim="editForm.TurnoverPay"
                                  :disabled="store.AdminInfo.GroupId > 1"
                                  onkeyup="value=value.replace(/[^\d\.\d{0,4}]/g, '')" maxlength="6"/>
                        </div>
                    </el-form-item>

                    <el-form-item :label="$t('钱包类型') + ':'" v-if="editForm.OperatorType == 2">

                        <div style="font-size: 14px;text-align: right">
                            {{ editForm.WalletMode == 2 ? $t('单一钱包') : $t('转账钱包') }}
                        </div>

                    </el-form-item>
                    <el-form-item :label="$t('归属国家') + ':'">
                        <el-select
                                v-model="editForm.BelongingCountry"
                                :placeholder="$t('请选择归属国家')"
                                :disabled="store.AdminInfo.GroupId > 1"
                                filterable
                                clearable
                        >
                            <template v-for="item in props.belongingCountryList">

                                <el-option :label="$t(item.label)" v-if="item.value != 'ALL'" :value="item.value"/>
                            </template>

                        </el-select>
                    </el-form-item>

                    <el-form-item :label="$t('备注') + ':'" v-if="store.AdminInfo.GroupId <= 2">
                        <el-input
                            v-model="editForm.Remark"
                            :rows="3"
                            type="textarea"
                            :placeholder="$t('备注')"
                        />
                    </el-form-item>

                    <el-form-item :label="$t('商户币种') + ':'">
                        <el-select
                                v-model.trim="editForm.CurrencyKey"
                                :placeholder="$t('请选择商户币种')"
                                filterable
                                disabled
                                clearable
                        >
                            <template v-for="item in CurrencyConfig">
                                <el-option :label="`【${item.value}】${item.label}`" :value="item.value"
                                           v-if="item.value != 0"/>
                            </template>

                        </el-select>
                    </el-form-item>

                    <el-form-item :label="$t('基础货币单位') + ':'" v-if="editForm.OperatorType == 2">
                        <el-select
                            v-model="editForm.CurrencyCtrlStatus"
                            :placeholder="$t('请选择基础货币单位')"
                            disabled

                        >
                            <template v-for="(item, index) in props.currencyRate">

                                <el-option :label="$t(item)" :value="Number(index)"/>
                            </template>

                        </el-select>
<!--                        <div class="switchContainer">-->
<!--                            <el-switch-->
<!--                                v-model="editForm.CurrencyRate"-->
<!--                                :active-value="1"-->
<!--                                :inactive-value="0"-->
<!--                            />-->
<!--                        </div>-->
                    </el-form-item>


<!--                    <el-form-item :label="$t('预付款金额') + ':'" prop="Advance">-->
<!--                        <el-input v-model="editForm.Advance" :placeholder="$t('请输入预付款金额')"-->
<!--                                  :disabled="store.AdminInfo.GroupId >1" @keydown="RegexpFloat"-->
<!--                                  @blur="RegexpFloat"></el-input>-->
<!--                        <div style="font-size: 14px;text-align: right"> {{ $t('剩余预付款') }}：-->


<!--                            &lt;!&ndash;                            <span v-if="editForm.CurrencyKeyName">{{`${editForm.Advance}(${editForm.CurrencyKeyName})`}}</span>&ndash;&gt;-->
<!--                            <span v-if="editForm.CurrencyKeyName">{{ `${editForm.Advance}(USDT)` }}</span>-->
<!--                        </div>-->
<!--                    </el-form-item>-->



                    <el-form-item prop="CurrencyRate" v-if="editForm.OperatorType == 2">
                        <template #label>
                            <div class="customLabel"> {{ $t('RTP区间') }}:</div>
                        </template>
                        <el-select
                            v-model="editForm.HighRTPOff"
                            :disabled="store.AdminInfo.GroupId > 1"
                            :placeholder="$t('请选择RTP区间')">
                            <template v-for="(item, index) in OperatorRTP" >
                                <el-option :label="`${item.split('-')[0]}%~${item.split('-')[1]}%`" :value="Number(index)" v-if="index != 0"/>
                            </template>

                        </el-select>
                    </el-form-item>

                    <el-form-item :label="$t('联系方式') + ':'">

                        <el-row :gutter="5" v-for="(item, index) in editForm.Contact" :key="index"
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
                                               v-if="index == editForm.Contact.length - 1" @click="addContact"/>

                                </div>
                            </el-col>

                        </el-row>
                    </el-form-item>
                </el-col>


                <el-col :span="12" class="col-right">

                    <el-form-item :label="$t('谷歌密钥和密码') + ':'" style="display: none">
                        <el-input
                                v-model="editForm.Address"
                                :placeholder="$t('谷歌密钥和密码')"
                        />

                    </el-form-item>


                    <el-form-item :label="$t('游戏厂商') + ':'"
                                  v-if="editForm.OperatorType == 2 && store.AdminInfo.GroupId == 1"
                                  style="margin-bottom: 5px">
                        <el-select
                            v-model="editForm.DefaultManufacturerOn"
                            @change="ManufacturerChange"
                            :placeholder="$t('请选择游戏厂商')"
                            multiple
                            collapse-tags
                        >
                            <template v-for="(item, index) in ManufacturerList">

                                <el-option :label="$t(item.ManufacturerName)" :value="item.ManufacturerCode"/>
                            </template>

                        </el-select>

                    </el-form-item>
                    <template v-if="editForm.OperatorType == 2 && store.AdminInfo.GroupId != 2">
                        <!--                        <el-form-item :label="$t('服务器地址') + ':'" >-->
                        <!--                            <el-input-->
                        <!--                                v-model="editForm.WhiteIps"-->
                        <!--                                :placeholder="$t('请输入服务器地址')"-->
                        <!--                            />-->
                        <!--                        </el-form-item>-->

                        <el-form-item v-if="editForm.WalletMode == 2" :label="$t('服务器回调地址') + ':'"
                                      prop="Address" style="margin-bottom: 5px">
                            <el-input
                                    v-model="editForm.Address"
                                    :placeholder="$t('请输入服务器回调地址')"
                            />

                        </el-form-item>

                        <div style="color: var(--el-color-danger);text-align: right;font-size: 12px;margin-top: 15px;margin-bottom: 15px; width: 85%"
                             v-if="editForm.WalletMode == 2">
                            {{ $t('慎重填写回调地址,填写有误将无法正常游戏.') }}
                        </div>
                        <!--                        prop="UserWhite"-->
                        <el-form-item v-if="editForm.WalletMode == 1" :label="$t('服务器IP白名单') + ':'">
                            <el-input
                                    v-model="editForm.UserWhite"
                                    :rows="3"
                                    type="textarea"
                                    :placeholder="$t('请输入服务器IP白名单')"
                            />
                        </el-form-item>


                    </template>


                    <el-form-item :label="$t('开户权限') + ':'" prop="WalletMode" v-if="props.operatorData.OperatorType == 1">

                        <el-radio-group v-model="editForm.ReviewType" style="margin-top: -4px">
                            <el-radio :value="item.value" size="large" v-for="item in CREATE_MERCHANT_RULE">
                                {{ $t(item.label) }}
                            </el-radio>
                        </el-radio-group>
                    </el-form-item>

                    <el-form-item :label="$t('风控语言') + ':'">
                        <el-select
                            v-model="editForm.Lang"
                            :placeholder="$t('请选择风控语言')"

                            clearable
                        >
                            <template v-for="item in languageConfig">

                                <el-option :label="$t(item.label)" :value="item.value"/>
                            </template>

                        </el-select>
                    </el-form-item>

                    <el-form-item :label="$t('机器人识别码') + ':'"
                                  v-if="props.operatorData.OperatorType != 1 && store.AdminInfo.GroupId != 2">
                        <el-input
                                v-model="editForm.Robot"

                                :placeholder="$t('请输入机器人识别码')"
                        >
                            <template #suffix>
                                <el-tooltip
                                        class="box-item"
                                        effect="dark"
                                        :content="$t('Telegram机器人识别码')"
                                        placement="bottom-end"

                                >
                                    <el-icon size="18" color="#000000">
                                        <QuestionFilled/>
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('chatId') + ':'"
                                  v-if="props.operatorData.OperatorType != 1 && store.AdminInfo.GroupId != 2">
                        <el-input
                                v-model="editForm.ChatID"
                                :placeholder="$t('请输入chatId')"
                        >
                            <template #suffix>
                                <el-tooltip
                                        class="box-item"
                                        effect="dark"
                                        :content="$t('Telegram机器人所在频道ID')"
                                        placement="bottom-end"

                                >
                                    <el-icon size="18" color="#000000">
                                        <QuestionFilled/>
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-input>
                    </el-form-item>

                    <el-form-item :label="$t('总控机器人识别码') + ':'"
                                  v-if="props.operatorData.OperatorType != 1 && store.AdminInfo.GroupId == 1">
                        <el-input
                                v-model="editForm.AdminRobot"
                                :placeholder="$t('请输入机器人识别码')"
                        >
                            <template #suffix>
                                <el-tooltip
                                        class="box-item"
                                        effect="dark"
                                        :content="$t('Telegram机器人识别码')"
                                        placement="bottom-end"

                                >
                                    <el-icon size="18" color="#000000">
                                        <QuestionFilled/>
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('余额报警chatId') + ':'"
                                  v-if="props.operatorData.OperatorType != 1 && store.AdminInfo.GroupId == 1">
                        <el-input
                                v-model="editForm.AdminChatID"
                                :placeholder="$t('请输入余额报警chatId')"
                        >
                            <template #suffix>
                                <el-tooltip
                                        class="box-item"
                                        effect="dark"
                                        :content="$t('Telegram机器人所在频道ID')"
                                        placement="bottom-end"

                                >
                                    <el-icon size="18" color="#000000">
                                        <QuestionFilled/>
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-input>
                    </el-form-item>

                    <el-form-item :label="$t('余额') + ':'" v-if="store.AdminInfo.GroupId <= 1  && props.operatorData.OperatorType != 1">

                        <el-input

                            v-model="editForm.Balance"
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
                        <div style="color: var(--el-color-danger);font-size: 12px">{{$t('在商户基础信息界面中点击修改才会生效')}}</div>
                    </el-form-item>


                    <el-form-item v-if="store.AdminInfo.GroupId <= 1 && props.operatorData.OperatorType != 1" style="width:100%">


                        <el-form-item :label="$t('余额不足阈值') + ':'" style="margin-right: 0;width: 100%">
                            <el-input v-model="editForm.BalanceThreshold"
                                      :placeholder="$t('请输入余额不足阈值')"/>
                        </el-form-item>
                        <el-form-item :label="$t('余额不足间隔') + ':'" style="margin-right: 0;width: 100%">
                            <el-select v-model="editForm.BalanceThresholdInterval"
                                       :placeholder="$t('请选择余额不足间隔时间')">

                                <el-option v-for="item in BALANCE_THRESHOLD" :label="$t(item.label)"
                                           :value="item.value"/>

                            </el-select>


                        </el-form-item>
                        <el-form-item :label="$t('欠费提示间隔') + ':'" style="margin-right: 0;width: 100%">
                            <el-select v-model="editForm.ArrearsThresholdInterval"
                                       :placeholder="$t('请选择欠费提示间隔')">

                                <el-option v-for="item in ARREARS_HINT_INTERVAL" :label="$t(item.label)"
                                           :value="item.value"/>

                            </el-select>
                        </el-form-item>



                        <el-row style="width: 90%">
                            <el-col :span="12" class="col-right">
                                <el-form-item :label="$t('退出按钮') + ':'">
                                    <div class="switchContainer">
                                        <el-switch
                                                v-model="editForm.ShowExitBtnOff"
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
                                            v-model="editForm.CurrencyVisibleOff"
                                            :active-value="1"
                                            :inactive-value="0"
                                        />
                                    </div>
                                </el-form-item>
                            </el-col>
                            <el-form-item :label="$t('退出按钮链接') + ':'" v-if="editForm.ShowExitBtnOff">
                                <el-input
                                    v-model="editForm.ExitLink"
                                    :placeholder="$t('请输入退出按钮链接')"
                                />

                            </el-form-item>
                            <el-col :span="12" class="col-right btn-right-content">

                                <el-form-item :label="$t('RTP设置') + ':'">
                                    <div class="switchContainer">
                                        <el-switch
                                            v-model="editForm.RTPOff"
                                            :active-value="1"
                                            :inactive-value="0"
                                        />
                                    </div>
                                </el-form-item>

                                <el-form-item :label="$t('赢取最高倍数') + ':'">
                                    <div class="switchContainer">
                                        <el-switch
                                            v-model="editForm.MaxMultipleOff"
                                            :active-value="1"
                                            :inactive-value="0"
                                        />
                                    </div>
                                </el-form-item>
                                <el-form-item :label="$t('购买RTP') + ':'">
                                    <div class="switchContainer">
                                        <el-switch
                                            v-model="editForm.BuyRTPOff"
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
                                            v-model="editForm.PlayerRTPSettingOff"
                                            :active-value="1"
                                            :inactive-value="0"
                                        />
                                    </div>
                                </el-form-item>
<!--                                <el-form-item label="">-->
<!--                                    <div class="switchContainer" style="visibility: hidden">-->
<!--                                        <el-switch-->
<!--                                            v-model="editForm.PlayerRTPSettingOff"-->
<!--                                            :active-value="1"-->
<!--                                            :inactive-value="0"-->
<!--                                        />-->
<!--                                    </div>-->
<!--                                </el-form-item>-->
                                <el-form-item :label="$t('赢取最高钱数') + ':'">
                                    <div class="switchContainer">
                                        <el-switch
                                            v-model="editForm.MaxWinPointsOff"
                                            :active-value="1"
                                            :inactive-value="0"
                                        />
                                    </div>
                                </el-form-item>

                            </el-col>

                        </el-row>
                    </el-form-item>
                </el-col>


            </el-row>

        </el-form>
        <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="EditAdminer">{{ $t('修改') }}</el-button>
                </span>
        </template>
    </el-dialog>

        <operatorBanlance v-model="operatorBalanceVisible" handel-type="edit" :merchant-balance="LogOriginData.Balance" @commitMerchantBalance="calcBalance"/>
    </div>
</template>


<script setup lang="ts">
import { ElMessageBox, FormRules} from "element-plus";
import {computed, reactive, Ref, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "@/pinia";
import { Minus, Plus} from "@element-plus/icons-vue";
import {tip} from "@/lib/tip";
import {Client} from "@/lib/client";
import {AdminInfo, EditMaintenance} from "@/api/adminpb/info";
import {AdminGroup} from "@/api/adminpb/group";
import {ARREARS_HINT_INTERVAL, BALANCE_THRESHOLD, CREATE_MERCHANT_RULE, LANG_LIST} from "@/lib/config";
import OperatorBanlance from "@/pages/systemManagement/component/operatorBanlance.vue";
import {AddLog} from "@/lib/commRequest";
import {LOG_OPERATOR_TYPE, LOG_PRIMARY_MODULE} from "@/api/adminpb/log";
import {OperatorRTP} from "@/enum";
import {Manufacturer} from "@/api/gamepb/game";

let {t} = useI18n()
const editFormRef = ref(null)

const isIPv4validator = (rule: any, value: any, callback: any) => {
    let whitelistString = value.split(',').filter(item => item !== "")
    const isIPv4 = new RegExp(/^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})$/);
    let whiteBools = whitelistString.filter(t => isIPv4.test(t));
    if (whiteBools.length) {
        callback()
    } else {
        callback(new Error(t('白名单不合法')))
    }
}
const validContact = (rule: any, value: any, callback: any) => {

    let Contact = editForm.value["Contact"].filter(item => item.name != "" && item.value != "")

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
// const addRules = reactive<FormRules>({
//     Name: [{required: true, message: t('商户名称不能为空'), trigger: 'blur',}],
//     Address: [{required: true, message: t('商户回调地址不能为空'), trigger: 'blur'}],
//     WhitelistString: [{required: true, trigger: 'blur', validator: isIPv4validator,}],
//     MenuIds: [{required: true, message: t('请给商户管理员配置权限'), trigger: 'change'}],
// })
const editRules = reactive<FormRules>({
    UserWhite: [{required: true, message: t('用户白名单不能为空'), trigger: 'blur'}, {
        trigger: 'blur',
        validator: isIPv4validator
    }],
    Address: [{
        trigger: 'blur',
        validator: isIPv4validator2
    }],
    Contact: [{required: true, message: t('联系方式至少填写一条'), trigger: 'blur'}, {
        trigger: 'blur',
        validator: validContact
    }],
    Advance: [{required: true, message: t('预付金额不能为空'), trigger: 'blur'}],
})

const store = useStore()

const props = defineProps(["modelValue", "operatorType", "operatorStatus", "incomeType", "walletOptions", "operatorData", "belongingCountryList", "currencyRate"])
const emits = defineEmits(['update:modelValue', "editData"])

const editStatus = ref(false)
const operatorBalanceVisible = ref(false)

// LogOriginData

let LogOriginData = {}

const addDialog = computed(() => {
    if (!props.modelValue) {
        editStatus.value = false
        editForm.value = {
            Robot: "",
            ChatID: "",
            Name: "",
            NextRate: 0,
            Status: 0,
            RTPOff: 0,
            PlayerRTPSettingOff: 0,
            DefaultManufacturerOn: null,
            ReviewStatus: 0,
            ShowExitBtnOff: 0,
            BuyRTPOff:0,
            CurrencyCtrlStatus: 0,
            // 商户编码
            OperatorType: 0,                   // 商户类型
            UserName: "",                        // 商户主账号
            AppID: "",                           // 商户名称
            PlatformPay: 0,                      // 平台费
            CooperationType: "",                 // 合作模式
            Advance: 0,                          // 预付款金额
            CurrencyKey: "",                     // 商户币种
            CurrencyKeyName: "",                     // 商户币种
            Contact: [],     // 联系方式
            WalletMode: 2,                       // 钱包类型
            Surname: "",                         // 会员前缀
            Lang: "en",                            // 客户端默认语言
            ServiceIp: "",                       // 服务器地址
            WhiteIps: "",                        // 服务器IP白名单
            Address: "",                         // 服务器回调地址
            UserWhite: "",                       // 用户白名单
            LoginOff: 0,                        // 服务器IP白名单
            FreeOff: 0,                         // 服务器回调地址
            DormancyOff: 0,                     // 用户白名单
            RestoreOff: 0,                      // 记录开启
            ManualFullScreenOff: 0,             // 免游游戏开关板
            NewGameDefaulOff: 0,                // 防休眠开启
            MassageOff: 0,                      // 断复开关
            MassageIp: "",                      // 手动全屏开启
            Remark: "",                          // 手动全屏开启
            StopLoss: 0,                        // 消息推送
            MaxMultipleOff: 0,                  // 消息推送地址
            _LineMerchant: 1,                   // 游戏新型预设
            MaxWinPointsOff: 1,                   // 游戏新型预设
            CurrencyVisibleOff: 1,                   // 是否显示币种
            AdminRobot: "",                             //总控机器人的chatid
            AdminChatID: "",                            //余额告警阈值
            BalanceAlert: 0,                            //余额告警时间 暂时固定为24每天0点进行检测并触发余额告警
            BalanceAlertTimeInterval: 1,
            Balance: 0,
            HighRTPOff: 0,
            ReviewType: CREATE_MERCHANT_RULE[0].value,
            BalanceThreshold: 0,
            BalanceThresholdInterval: BALANCE_THRESHOLD[0].value,
            ArrearsThresholdInterval: ARREARS_HINT_INTERVAL[0].value,
        }
    }
    return Boolean(props.modelValue)
})

const operatorBalanceForm = ref(null)
const initEditData = async () => {


    let [response, err] = await Client.Do(AdminInfo.GetOperatorInfo, {AppID: props.operatorData.AppID})

    if (err) {
        return tip.e(t(err))
    }

    editForm.value = {...response.OperatorInfo}

    let Contact = []

    if (!editForm.value.DefaultManufacturerOn){
        editForm.value.DefaultManufacturerOn = ["ALL"]
    }

    if (editForm.value.Contact.length && typeof editForm.value.Contact == "string") {
        let parseContact = JSON.parse(editForm.value.Contact)
        for (const i in parseContact) {
            if (parseContact[i].name && parseContact[i].value) {

                Contact.push(parseContact[i])
            }
        }
        if (Contact.length <= 0){
            Contact.push({name:"", value:""})
        }

        editForm.value.Contact = Contact
    }

    if (editForm.value.WhiteIps && editForm.value.WhiteIps instanceof Array) {

        editForm.value.WhiteIps = editForm.value.WhiteIps.join(",")
    }

    let currencyItem = CurrencyConfig.value.find(item => item.value == editForm.value.CurrencyKey)

    if (currencyItem) {

        editForm.value.CurrencyKeyName = currencyItem.CurrencySymbol

    } else {
        editForm.value.CurrencyKeyName = ""
    }

    let incomeType = props.incomeType.find(item => item.value == editForm.value.CooperationType)

    if (incomeType && incomeType.label) {
        editForm.value.CooperationTypeText = t(incomeType.label)
    }


    if (!editForm.value.ReviewType && editForm.value.ReviewType != 0){
        editForm.value.ReviewType = CREATE_MERCHANT_RULE[0].value
    }
    operatorBalanceForm.value = null

    if (store.language != 'zh'){

        let langList = []
        for (const i in languageConfig.value) {
            langList.push({label:languageConfig.value[i].value, value:languageConfig.value[i].value})
        }

        languageConfig.value = langList
    }

    editForm.value.Balance =  parseFloat(parseFloat(editForm.value.Balance).toFixed(2))

    LogOriginData = {...editForm.value}
}


const editForm: Ref<EditMaintenance> = ref(<EditMaintenance>{
    NextRate: 0,
    Status: 0,
    RTPOff: 0,
    CurrencyCtrlStatus: 0,
    DefaultManufacturerOn: null,
    Robot: "",
    ChatID: "",
    ReviewStatus: 0,
    BuyRTPOff: 0,
    TurnoverPay: 0,
    // 商户编码
    UserName: "",                   //  商户类型
    OperatorType: 2,                   //  商户类型
    PlayerRTPSettingOff: 0,                   //  商户类型
    ShowExitBtnOff: 0,
    Name: "",                            //  商户主账号
    AppID: "",                           //  商户名称
    PlatformPay: 0,                      //  平台费
    CooperationType: "",                 //  合作模式
    Advance: 0,                          //  预付款金额
    CurrencyKey: "",                     //  商户币种
    Contact: [{name: "", value: ""}],     //  联系方式
    WalletMode: 2,                       //  钱包类型
    Surname: "",                         //  会员前缀
    Lang: "en",                            //  客户端默认语言
    ServiceIp: "",                       //  服务器地址
    WhiteIps: "",                        //  服务器IP白名单
    Address: "",                         //  服务器回调地址
    UserWhite: "",                       //  用户白名单
    LoginOff: 0,                        //  服务器IP白名单
    FreeOff: 0,                         //  服务器回调地址
    DormancyOff: 0,                     //  用户白名单
    HighRTPOff:1,
    RestoreOff: 0,                      //  记录开启
    ManualFullScreenOff: 0,             //  免游游戏开关板
    NewGameDefaulOff: 0,                //  防休眠开启
    MassageOff: 0,                      //  断复开关
    MassageIp: "",                      //  手动全屏开启
    StopLoss: 0,                        //  消息推送
    MaxMultipleOff: 0,                  //  消息推送地址
    _LineMerchant: 1,                   //  游戏新型预设
    Remark: "",                   //  游戏新型预设
    ExitLink:"http://",
    CurrencyRate: 0,
    CurrencyKeyName: "",
    MaxWinPointsOff: 0,
    CurrencyVisibleOff: 1,                   // 是否显示币种
    AdminRobot: "",                             //总控机器人的chatid
    AdminChatID: "",                            //余额告警阈值
    BalanceAlert: 0,                            //余额告警时间 暂时固定为24每天0点进行检测并触发余额告警
    BalanceAlertTimeInterval: 1,
    Balance: 0,
    HighRTPOff: 1,
    ReviewType: 2,
    BalanceThreshold: 1500,
    BalanceThresholdInterval: BALANCE_THRESHOLD[0].value,
    ArrearsThresholdInterval: ARREARS_HINT_INTERVAL[0].value,
})

const languageConfig = ref(LANG_LIST)
const CurrencyConfig = ref([])
const ManufacturerList = ref([])


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
    editForm.value.DefaultManufacturerOn = ["ALL"]


}

getCurrency()
getManufacturerList()


const EditAdminer = async () => {


    editFormRef.value.validate(async valid => {
        if (valid) {

            const params = {...editForm.value}


            if (typeof params["PlatformPay"] == "string") {

                params["PlatformPay"] = parseFloat(params["PlatformPay"])
            }
            if (typeof params["TurnoverPay"] == "string") {

                params["TurnoverPay"] = parseFloat(params["TurnoverPay"])
            }
            if (typeof params["Advance"] == "string") {

                params["Advance"] = parseFloat(params["Advance"])

            }
            if (typeof params["NextRate"] == "string") {

                params["NextRate"] = parseFloat(params["NextRate"])

            }
            let Contact = params["Contact"].filter(item => item.name != "" && item.value != "")


            params["WhiteIps"] = params["WhiteIps"].split(",")
            params["Contact"] = JSON.stringify(Contact)

            if (params["Robot"] == "" && params ["ChatID"] != "") {

                tip.e(t("chatId填写机器码不能为空"))
                return
            }
            if (params["Robot"] != "" && params ["ChatID"] == "") {
                tip.e(t("机器码填写chatId不能为空"))
                return
            }

            params["BalanceThreshold"] = Number(params["BalanceThreshold"])
            params["Balance"] = Number(params["Balance"])

            if (params["DefaultManufacturerOn"].indexOf("ALL") != -1){
                params["DefaultManufacturerOn"] = null
            }

            ElMessageBox.confirm(
                t('确认取消修改当前商户'),
                t('是否确认'),
                {
                    confirmButtonText: t('确定'),
                    cancelButtonText: t('关闭'),
                    type: 'warning',
                }
            )
                .then(async () => {


                    const [res, err] = await Client.Do(AdminInfo.EditMaintenance, params)

                    if (err != null) {
                        tip.e(t("修改失败"))
                        return
                    }


                    if (operatorBalanceForm.value){

                        AddLog(LOG_PRIMARY_MODULE.OPERATOR_LIST, LOG_OPERATOR_TYPE.BALANCE_EDIT, LogOriginData, params, params.AppID, operatorBalanceForm.value)
                    }


                    AddLog(LOG_PRIMARY_MODULE.OPERATOR_LIST, LOG_OPERATOR_TYPE.ADD, LogOriginData, params, store.AdminInfo.AppID, `修改{AppID}`)



                    emits('update:modelValue')
                    emits('editData', {key:params.Name, value:params})


                })
        }
    })

}

const ManufacturerChange = (value) => {

    if (value.length == 0){
        editForm.value.DefaultManufacturerOn = ["ALL"]
        return
    }

    if (value[value.length - 1] == "ALL"){
        editForm.value.DefaultManufacturerOn = ["ALL"]
    }else{
        editForm.value.DefaultManufacturerOn = value.filter(item => item != "ALL")
    }


}

const addContact = () => {
    let findNull = editForm.value.Contact.find(item => item.name == "" || item.value == "")
    if (findNull) {

        tip.e(t("有未填写的联系方式,请补全后添加"))
        return
    }
    editForm.value.Contact.push({
        name: "",
        value: ""
    })

}

const delContact = (index) => {

    if (editForm.value.Contact.length == 1) {
        return
    }

    editForm.value.Contact.splice(index, 1)
}

const inputContact = (index, key, event) => {

    editForm.value.Contact[index][key] = event
}


const RegexpFloat = (e) => {
    e.target.value = (e.target.value.match(/^\d*(\.?\d{0,2})/g)[0]) || null

}


const editNextRate = () => {
    // EditAdminer()
    editForm.value.NextRate = editForm.value.NextRate == "" ? 0 : editForm.value.NextRate


    editForm.value.NextRate = parseFloat(editForm.value.NextRate)
    editStatus.value = false
}

const closeDialog = () => {

}


const calcBalance = (value) => {
    operatorBalanceForm.value = value.data
    editForm.value.Balance = value.handledBalance
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

  .col-right > * {
    width: 80%;

  }

}
.btn-right-content>.el-form--inline .el-form-item{
    width: 100%;
    margin-right: 0;
}

</style>
