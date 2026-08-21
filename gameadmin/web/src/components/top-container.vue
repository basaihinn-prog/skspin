/*
*@Author: 西西米
*@Date: 2022-12-06 17:21:20
*@Description: 顶部
*/
<template>
    <div class='top-container' :style="store.viewModel === 1 ? 'justify-content: flex-end;' :''">

        <div style="width: 100%;display: flex;justify-content: space-between;line-height: 46px">
            <div style="color: #000;display: flex;align-items: center">




                <el-breadcrumb separator="/">
                    <el-breadcrumb-item>
                        <el-button class="pageHeader_btn" text :icon="Expand" @click="sidBarShow" size="small"></el-button>
                        <el-button class="pageHeader_btn" text :icon="House" @click="goHome" size="small"></el-button>
                    </el-breadcrumb-item>
                    <el-breadcrumb-item v-if="router.currentRoute.value.meta.parents"><div class="pageHeader_btn" style="line-height: 30px">{{ $t(router.currentRoute.value.meta.parents) }}</div></el-breadcrumb-item>
                    <el-breadcrumb-item v-if="router.currentRoute.value.meta.parent"><div class="pageHeader_btn" style="line-height: 30px">{{ $t(router.currentRoute.value.meta.parent) }}</div></el-breadcrumb-item>
                    <el-breadcrumb-item>
                        <div class="pageHeader_btn" style="line-height: 30px">{{ $t(router.currentRoute.value.meta.title) }}</div>
                    </el-breadcrumb-item>
                </el-breadcrumb>

            </div>
            <el-space>
                <span class="font_size16"  style="display: flex;align-items: center" v-if="store.AdminInfo.GroupId > 2" >

                        <el-icon size="18" color="#9F9F9F" v-if="BalanceLoad">
                             <loading/>
                         </el-icon>

                    {{ $t('商户余额') }}：{{ store.AdminInfo.Businesses.Balance.toFixed(2) }}


                         <el-icon size="18" color="#9F9F9F" @click="refreshBalance">
                             <refresh/>
                         </el-icon>

                </span>
<!--                <span class="font_size16"  style="font-weight: bolder">{{ $t('服务器时区') }}：{{ UTC }}</span>-->
                <span class="font_size16"  style="font-weight: bolder">
                    {{ $t('服务器时区') }}：
                    <el-select v-model="TimeZone" placeholder="请选择" style="width: 100px" size="small" @change="TimeZoneChange">
                        <el-option
                            v-for="item in TimeZoneList"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value">
                        </el-option>
                      </el-select>
                </span>

<!--                <span class="font_size16"  style="font-weight: bolder">-->
<!--                      {{ $t('服务器区域') }}：-->
<!--                    <el-select v-model="serviceArea" placeholder="请选择" style="width: 100px" size="small" @change="serviceAreaChange">-->
<!--                        <el-option-->
<!--                            v-for="item in serviceAreaList"-->
<!--                            :key="item.value"-->
<!--                            :label="item.label"-->
<!--                            :value="item.value">-->
<!--                        </el-option>-->
<!--                      </el-select>-->
<!--                </span>-->
                <div class="language">
                    <el-dropdown size="small">

                        <el-button class="pageHeader_btn" text size="small">
                            <el-image :src="languageIcon" style="width: 18px"/>
                        </el-button>

                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item v-for="item in langRadio" @click="changeLanguage(item)">
                                    {{ item.lanNameNew }}
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
                <el-dropdown>
                    <div style="display:flex; align-items: center;">
                        <el-avatar size="small" :icon="UserFilled" />
                        <span style="margin-left: .2rem">{{ store.AdminInfo.Username }}</span>
                    </div>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <!-- <el-dropdown-item @click="refresh()">{{$t("刷新页面")}}</el-dropdown-item> -->
                            <el-dropdown-item @click="dialog = true">{{ $t("修改密码") }}</el-dropdown-item>
                            <el-dropdown-item @click="signOut()">{{ $t("退出登录") }}</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </el-space>
        </div>
    </div>

    <el-dialog v-model="dialog" @closed="" :title="$t('修改密码')" :width="store.viewModel === 2 ? '85%' : '650px'">
        <el-form ref="formRef" :model="param" :rules="rules" label-position="right" label-width="100px">
            <el-form-item prop="OldPasswd" :label="$t('旧密码')">
                <el-input v-model.trim="param.OldPasswd" clearable :placeholder="$t('请输入')"/>
            </el-form-item>
            <el-form-item prop="NewPasswd" :label="$t('新密码')">
                <el-input v-model.trim="param.NewPasswd" clearable :placeholder="$t('请输入')"/>
            </el-form-item>
            <el-form-item prop="ConfirmPasswd" :label="$t('确认密码')">
                <el-input v-model.trim="param.ConfirmPasswd" clearable :placeholder="$t('请输入')"/>
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="updatePasswdBtn">{{ $t('保存') }}</el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script setup lang='ts'>
import router from '@/router';
import {GetAdminInfo, useStore} from '@/pinia/index';
import {getCurrentInstance, ref, reactive, shallowRef, nextTick, onMounted} from "vue";
import {useI18n} from 'vue-i18n';
import {AdminInfo} from '@/api/adminpb/info';
import type {FormInstance, FormRules} from 'element-plus'
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import ut from "@/lib/util";
import {ElMessage, ElMessageBox} from 'element-plus'
import {Expand, House, RefreshRight, UserFilled} from '@element-plus/icons-vue'
import languageIcon from "@/assets/login/languageIcon.png"
let {t, locale} = useI18n();
const {proxy} = getCurrentInstance();
const store = useStore()
let dialog = ref(false)


let lang = shallowRef()
const formRef = ref<FormInstance>()
import {storeToRefs} from "pinia";
import {ServiceAreaList, TimeZone as TimeZoneList} from "@/enum";
import {AdminAuth} from "@/api/adminpb/auth";
const {language, UTC, langRadio} = storeToRefs(store)
const piniaStore = useStore()
const BalanceLoad = ref(false)
const {setLanguage, setToken, setAdminInfo,setMenuList, setActiveTabs,addTab,setElMenuActive} = piniaStore

let serviceArea = ref("AS")
let TimeZone = ref(0)
let serviceAreaList = ref(ServiceAreaList)


let param = reactive({
    OldPasswd: '',
    NewPasswd: '',
    ConfirmPasswd: ''
})
const rules = reactive<FormRules>({
    OldPasswd: [{required: true, message: t('请填写旧密码'), trigger: 'blur'}],
    NewPasswd: [{required: true, message: t('请填写新密码'), trigger: 'blur'}],
    ConfirmPasswd: [{required: true, message: t('请再次确认密码'), trigger: 'blur'}],
})
let initVal = {
    ID: 0,
    Icon: "Home",
    Pid: 0,
    Sort: 0,
    Title: '游戏报表',
    Url: "/dashboard"
}


let langFilter = langRadio.value?.filter(item => item.abbreviation === language.value)
const lanTitle = ref(langFilter.length > 0 ? langFilter[0].lanNameNew : langRadio.value[0].lanNameNew)


const sidBarShow = () => {
    store.setSidebarShow(!store.sidebarShow)
}

const updatePasswdBtn = async () => {
    if (!formRef) return
    await formRef.value.validate(async (valid, fields) => {
        if (valid) {
            try {
                let [data, err] = await Client.Do(AdminInfo.UpdatePasswd, param);
                if (err) {
                    tip.e(err);
                    return;
                }
                tip.s(t('修改成功'));
                dialog.value = false;
            } catch (error) {
                console.error(error); // 记录错误信息
                tip.e('修改失败'); // 显示错误提示
            }
        }
    })
}


onMounted(async () => {
    locale.value = store.language;
    serviceArea.value = localStorage.getItem("serviceArea")
    // getServiceAreaUri(serviceArea.value)
    TimeZone.value = store.SelectedTimeZone
});


let serviceAreaChange = async (value) => {

    let serviceAreaStorage = localStorage.getItem("serviceArea")

    if (!serviceAreaStorage){
        localStorage.setItem("serviceArea", value)
    }

    getServiceAreaUri(value)

    const data = {

        ID: store.AdminInfo.ID,
        AppID: store.AdminInfo.AppID,
        CreateAt: store.AdminInfo.CreateAt,
        GoogleCode: store.AdminInfo.GoogleCode,
        IsOpenGoogle: store.AdminInfo.IsOpenGoogle,
        LoginAt: store.AdminInfo.LoginAt,
        OperatorAdmin: store.AdminInfo.OperatorAdmin,
        PermissionId: store.AdminInfo.PermissionId,
        Qrcode: store.AdminInfo.Qrcode,
        Status: store.AdminInfo.Status,
        Token: store.Token,
        TokenExpireAt: store.AdminInfo.TokenExpireAt,
        Username: store.AdminInfo.Username,
        Password: store.AdminInfo.Password,
        GroupId: store.AdminInfo.GroupId,
    }
    const response = await AdminAuth.ChangeService(window.queryLocal, data)

    if (response.error){
        tip.e("切换失败")
        serviceArea.value = localStorage.getItem("serviceArea")
        return false
    }
    localStorage.setItem("serviceArea", value)
    serviceArea.value = value
    window.location.reload()
}

let TimeZoneChange = (value) => {

    localStorage.setItem("timeZone", value)
    store.setSelectedTimeZone(value)
}


let getServiceAreaUri = (val) => {
    const selectArea = serviceAreaList.value.find(item=> item.value == val)
    // window.queryLocal = selectArea.queryLocal

}

const signOut = async () => {
    try {
        let [data, err] = await Client.Do(AdminInfo.LoginOut, {})
        if (err) return tip.e(err)
        setToken('')
        setMenuList('')
        setAdminInfo({})
        store.tabMenuList = [initVal]
        setElMenuActive({
            title:'',
            value:''
        })
        localStorage.removeItem('game_store');
        await router.push({path: '/login'})
    } catch (e) {
        console.error(e); // 记录错误信息
        tip.e(t('退出登录失败')); // 显示错误提示
    }
}

const reloadRouterView = () => {

}

const goHome = () => {
    setElMenuActive({
        title: "游戏报表",
        value: "/dashboard"
    })
    router.push("/dashboard")
}

const changeLanguage = (language) => {
    ElMessageBox.confirm(
        t('切换多语言会导致标签页失效，确认切换吗?'),
        t('提示'),
        {
            confirmButtonText: t('切换'),
            cancelButtonText: t('不切换'),
            type: 'warning',
        }
    )
        .then(() => {
            lanTitle.value = language.lanNameNew
            setLanguage(language.abbreviation)
            Client.setLanguage(language.abbreviation)
            locale.value = language.abbreviation;
            // window.location.reload()
            setElMenuActive({
                title:'',
                value:''
            })
            nextTick(() => {
                router.push({path: '/dashboard'}).then(() => {
                    store.tabMenuList = [initVal]
                    setActiveTabs(t('游戏报表'))
                    location.reload()
                    console.info('s');
                }).catch((error) => {
                    console.info('e:', error);
                });
            });
        })
        .catch(() => {
        })
}
const refreshBalance = async () => {
    BalanceLoad.value = true
    let admininfo = await GetAdminInfo()
    setAdminInfo(admininfo)
    BalanceLoad.value = false
}
</script>

<style lang='scss'>
.top-container {
    width: 100%;
    height: 46px;
    background: #ffffff;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: $size16;
    color: #9F9F9F;
    padding: 0 15px;
    box-sizing: border-box;
    border-bottom: 1px solid #F5F6FA;
    box-shadow: 0 2px 11px 0 rgba(18, 31, 62, .1);

    .el-icon {
        font-size: $size16;
    }


    .language {
        display: flex;

        font-size: $size16;
        i {
            font-style: normal;
            cursor: pointer;
        }
    }
}
.el-dropdown {
    cursor: pointer;
    display: flex;
    align-items: center;
    font-size: $size16;
}
.pageHeader_btn{
    height: 30px;
    margin-left: 0!important;

}
</style>
