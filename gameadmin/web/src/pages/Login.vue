<template>
    <div class="flex flex_grow flex_child_center " v-loading="loading">
        <div class="login-container">
            <div class="img"></div>
            <div class="login-box">
                <div class="login-form">
                    <!-- <avatar class="avatar" /> -->
                    <div class="title">{{ $t('游戏管理后台') }}</div>
                    <el-form @keyup.enter="loginAction">
                        <el-form-item class="languageChange">
                            <el-select style="min-width: 240px"
                                v-model="selectLanguages"
                                :placeholder="$t('游戏管理后台')"
                                @change="changeLanguage"
                            >
                                <el-option
                                    v-for="(item,index) in langRadio"
                                    :key="index" :label="item?.lanNameNew"
                                    :value="item?.abbreviation"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item>
                            <el-input prefix-icon="Message" v-model.trim="loginForm.Username"
                                      :placeholder="$t('用户名')"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-input prefix-icon="Lock" v-model.trim="loginForm.Password"
                                      :placeholder="$t('密码')" show-password
                                      type="password"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-input prefix-icon="Iphone" v-model="loginForm.GooleCode"
                                      :placeholder="$t('谷歌验证码')"></el-input>
                        </el-form-item>
                        <el-button type="primary" style="width: 100%; margin-bottom: 30px"
                                   @click="loginAction">{{
                                $t('登录')
                            }}
                        </el-button>
                    </el-form>
                </div>

            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import router from '@/router';
import {GetAdminInfo, useStore} from '@/pinia/index';
import {AdminAuth} from '@/api/adminpb/auth';
import {getCurrentInstance, onBeforeMount, reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import ut from "@/lib/util";
import {AdminGroup} from "@/api/adminpb/group";
import {storeToRefs} from "pinia";
import {ElMessageBox} from "element-plus";
import {Merchant} from "@/api/adminpb/merchant";

const {proxy} = getCurrentInstance();
const { t,locale } = useI18n();
const store = useStore()
const {language,lang} = storeToRefs(store)
const piniaStore = useStore()
const {setLanguage,setLangRadio,setLang, setToken, setAdminInfo,setMenuList, setActiveTabs,addTab,setElMenuActive, setTipsMap} = piniaStore
const selectLanguages = ref(language.value || 'en')
let langRadio = ref([])
let storedLang = ref(null)
let loading = ref(false)
let initVal = {
    ID: 0,
    Icon: "Home",
    Pid: 0,
    Sort: 0,
    Title: '首页',
    Url: "/dashboard"
}
const loginForm = reactive({
    Username: "",
    Password: "",
    GooleCode: "",
})
onBeforeMount(async () => {
    await getLangList()
    await changeLanguage(selectLanguages.value)
})

const getLangList = async () => {
    try {
        let [data, err] = await Client.Do(AdminGroup.GetLang, {
            PageIndex: 1,
            PageSize: 10000
        })
        if (err) return tip.e(t(err))
        storedLang.value = data.List
        setLang(storedLang.value);
        if (!storedLang.value || !storedLang.value.length) {
            return tip.w(t('语言内容缺失'))
        }
        langRadio.value = storedLang.value[0]
            ? Object.keys(storedLang.value[0])
                .filter(list => list !== 'Permission')
                .map(list => {
                    const match = ut.LangList.find(alist => list.toLowerCase() === alist.abbreviation);
                    return match ? {lanNameNew: match.lanNameNew, abbreviation: match.abbreviation} : null;
                }).filter(item => item !== null)
            : []
        setLangRadio(langRadio.value);
    } catch (e) {
        tip.e(t(e) || t('语言获取失败'))
    }
}

const getTips = () => {
    setTipsMap()
}

const validLoginTime = (time) => {



    let validTime = ut.fmtDate(new Date(time).getTime() / 1000)
    let now = ut.fmtDate(new Date().getTime() / 1000)
    const validTimeFormat = ut.fmtDate(validTime, "YYYY-MM-DD")
    const nowFormat = ut.fmtDate(now, "YYYY-MM-DD")

    return validTimeFormat == nowFormat

}

const loginAction = async () => {
    setToken('')
    setMenuList('')
    setAdminInfo({})
    store.tabMenuList = [initVal]
    setElMenuActive(
        {
            title:'',
            value:''
        }
    )
    loading.value = true

    let userTime = {}



    try {
        let [data, err] = await Client.Do(AdminAuth.Login, {
            Username: loginForm.Username,
            Password: loginForm.Password,
            GooleCode: loginForm.GooleCode,
        })

        if (err) {
            loading.value = false
            return tip.e(t(err))
        }
        Client.setToken(data.Token)
        store.setToken(data.Token)
        let admininfo = await GetAdminInfo()
        store.setAdminInfo(admininfo)
        loading.value = false


        if (admininfo.Businesses.OperatorType == 2 && (admininfo.Businesses.Balance < admininfo.Businesses.BalanceThreshold)){





            const loginTime = JSON.parse(localStorage.getItem("loginTime"))
            let timeFlag = true
            if(loginTime){
                timeFlag = !validLoginTime(loginTime[loginForm.Username])
            }



            if (timeFlag && admininfo.Businesses.BalanceThresholdInterval == 2){
                userTime[loginForm.Username] = new Date()
                localStorage.setItem("loginTime",JSON.stringify(userTime))
                ElMessageBox.confirm(
                    `<div style="display: flex;align-items: center;justify-content: center;flex-wrap: wrap">
                    <div style="width: 100%; text-align: center">${t('商户余额较低，请联系客服及时充值！')}</div>
                    <div style="width:100%;color: red;text-align: center">${t('当前余额:{Num}', {Num: ut.toNumberWithCommaNormal(admininfo.Businesses.Balance.toFixed(2))})}</div>
                </div>`,
                    t('商户余额提示'),
                    {
                        dangerouslyUseHTMLString: true,
                        cancelButtonText: t('关闭'),
                    })

            }
            if (admininfo.Businesses.BalanceThresholdInterval == 1){
                ElMessageBox.confirm(
                    `<div style="display: flex;align-items: center;justify-content: center;flex-wrap: wrap">
                    <div style="width: 100%; text-align: center">${t('商户余额较低，请联系客服及时充值！')}</div>
                    <div style="width:100%;color: red;text-align: center">${t('当前余额:{Num}', {Num: ut.toNumberWithCommaNormal(admininfo.Businesses.Balance.toFixed(2))})}</div>
                </div>`,
                    t('商户余额提示'),
                    {
                        dangerouslyUseHTMLString: true,
                        cancelButtonText: t('关闭'),
                    })

            }


        }


        tip.s(t("登录成功"))

        await router.push({path: '/dashboard'})

    } catch (e) {
        console.error(e); // 记录错误信息
        tip.e(t("登录失败")); // 显示错误提示
    }

}
const changeLanguage = async (language) => {
    try {
        setLanguage(language)
        setActiveTabs(t('游戏报表'))
        locale.value = language
        Client.setLanguage(language)
        document.title = t('游戏管理后台')
    } catch (e) {
        console.error(e); // 记录错误信息
        tip.e(t('切换语言失败')); // 显示错误提示
    }
}
</script>

<style lang="less" scoped>

.login-container {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: url("../assets/login/login-background-light.jpg") no-repeat 100%;
}

.login-box {
    display: flex;
    align-items: center;
    text-align: center;
    background: #ffffff;
    padding: 15px 24px;
    border-radius: 15px;
    box-shadow: 0 2px 15px #bababa ;
    -moz-box-shadow: 0 2px 15px #bababa ;
    -webkit-box-shadow: 0 2px 15px #bababa ;
}

.login-form {
    width: 360px;
}

.avatar {
    width: 350px;
    height: 80px;
}

.login-form .title {
    text-transform: uppercase;
    margin: 15px 0;
    color: #999;
    font: bold 200% Consolas, Monaco, monospace;
}

@media screen and (max-width: 1180px) {
    .login-container {
        grid-gap: 9rem;
    }

    .login-form {
        width: 290px;
    }

    .login-form .title {
        font-size: 2.4rem;
        margin: 8px 0;
    }
}

@media screen and (max-width: 968px) {


    .login-container {
        grid-template-columns: 1fr;
    }

    .login-box {
        justify-content: center;
    }

    .img {
        display: none;
    }
}

.tabs {
    padding: 20px;
}

.login_box {
    max-width: 450px;
    width: 100%;
    margin: 6% 30px 12%;
    background-color: #fff;
    border-radius: 3px;
}

.login_form {
    padding-top: 20px;
}

.languageChange {
    :deep .el-form-item__content {
        justify-content: center;
    }
}
</style>
