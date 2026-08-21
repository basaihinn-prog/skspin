import { AdminInfo, AdminInfoResp , DBMenu } from '@/api/adminpb/info';
import { MenuListResp } from '@/api/adminpb/menu';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { defineStore } from 'pinia'
import router from "@/router";
import {AdminGameCenter, GetExcleBetLogAddReq} from "@/api/gamepb/admin";
import {Language} from "@/api/adminpb/language";
// import {useI18n} from "vue-i18n";
// import i18n from "@/language/i18n";
// const { t } = i18n.global;
export function MenuListToTree(data: MenuListResp) {

    let map = {
        0: {Children:[]},
    }

    for (let one of data.MenuList) {
        map[one.ID] = one
        one.Url = one.Url.startsWith('/') ? one.Url : '/' + one.Url
    }

    for (let one of data.MenuList) {
        let parent = map[one.Pid]
        if (parent) {
            if (!parent.Children) {
                parent.Children = []
            }
            parent.Children.push(one)
        }
    }

    return map[0].Children
}

export async function GetAdminInfo() : Promise<AdminInfoResp> {
    let [data, err] = await Client.Do(AdminInfo.GetAdminInfo, {})
    if (err) {
        tip.e(err)
        return
    }
    data.MenuList = MenuListToTree(data)
    // console.log(data,'pina');
    return data
}

export async function GetTimeUtc() : Promise<AdminInfoResp> {
    let [data, err] = await Client.Do(AdminInfo.GetTimeUtc, {})
    if (err) {
        tip.e(err)
        return
    }
    // console.log(data,'pina');
    return data
}

function convertTimeZoneToUTC(data) {
    // Check if the time zone offset is valid
    let  timeZoneOffset= data.Str
    const regex = /^([+-])(\d{2})(\d{2})$/;
    const match = timeZoneOffset.match(regex);

    if (!match) {
        throw new Error("输入的时区格式不正确。");
    }
    // Extract the hours and minutes from the time zone offset
    const hours = parseInt(timeZoneOffset.slice(1, 3));
    const minutes = parseInt(timeZoneOffset.slice(3));

    // Convert the hours and minutes to seconds
    const hoursInSeconds = hours * 60 * 60;
    const minutesInSeconds = minutes * 60;

    // Calculate the total offset in seconds
    const totalOffsetInSeconds = hoursInSeconds + minutesInSeconds;

    // Create a UTC time string
    const utcTimeString = `GMT${timeZoneOffset.charAt(0)}${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;

    return utcTimeString;
}

interface ItemEx extends DBMenu {
    Children : ItemEx[]
}

export const useStore = defineStore({
    id: 'mystore',
    state: () => {
        return {
            Token: '',
            SelectedTimeZone: 0,
            tabMenuList: [],
            SystemConfig:{
                tableBorderConfig: null,
            },
            AdminInfo: {
                MenuList: [] as ItemEx[],
                Username: "",
                AppID: "",
                OperatorId: null,
                OperatorName:'',
                IsDefaultPermission:false,
                PermissionId: 0,
                GroupId: 0,
                BusinessesId: 0,
                Businesses: {}
            },
            UTC:'0',
            viewModel: 1,//2 隐藏左边菜单栏
            sidebarShow: false,
            GameList: [],
            LineMap: {},
            tipsMap:{},
            paginationLayout: 'total, sizes, prev, pager, next, jumper',
            language: 'en',
            lang: [],
            RouterMsg:{},
            // activeTabs:t('首页'),
            activeTabs:'',
            elMenuActive: {
                title:"",
                value:""
            },
            langRadio:[],
            runtimeTable:null
        }
    },
    getters: {},
    actions: {
        setViewModel(value: 1 | 2) {
            this.viewModel = value
            this.paginationLayout = value === 1 ? 'total, sizes, prev, pager, next, jumper' : 'total, prev, pager, next,'
        },
        setSidebarShow(value: boolean) {
            this.sidebarShow = value
        },
        setLanguage(value) {
            this.language = value
        },
        setToken(value: string) {
            this.Token = value
        },
        setActiveTabs(value: string) {
            this.activeTabs = value
        },
        setElMenuActive(value) {
            this.elMenuActive = value
        },
        setMenuList(value: string) {
            this.MenuList = value
            this.AdminInfo.MenuList = value
        },
        setTableSetConfig(value) {
            this.SystemConfig.tableBorderConfig = value
        },
        setSelectedTimeZone(value) {
            this.SelectedTimeZone = value
        },
        setGameList(value) {
            this.GameList = value


            value.forEach(item=>{
                this.LineMap[item.ID] = item.LineNum
            })

        },
        async setTipsMap() {
            const [res, err] = await Client.Do(Language.SelectConfig, {fileName: "tips.json"})

            let tipMapData = {}

            for (let i in res.Context) {
                let moudle = res.Context[i].moudle
                let tableHeader = res.Context[i].tableHeader
                let tips = res.Context[i].tips

                if (!tipMapData[moudle]){
                    tipMapData[moudle] = {}
                }
                tipMapData[moudle][tableHeader] = tips

            }

            this.tipsMap = tipMapData

        },
        setLang(value) {
            this.lang = JSON.parse(JSON.stringify(value))
        },
        setLangRadio(value) {
            this.langRadio = value
        },
        async initAdminInfo() {
            let info = await GetAdminInfo()
            let utc = await GetTimeUtc()
            if (info) {
                this.setAdminInfo(info)
            }
            if (utc) {
                this.setUTC(convertTimeZoneToUTC(utc))
            }
        },
        setAdminInfo(info) {
            this.AdminInfo = info
        },
        setUTC(info) {
            this.UTC = info
        },
        addTab(tab) {
            // 检查选项卡是否已经存在
            const existingTab = this.tabMenuList.find((t) => t.Url === tab.Url);
            if (!existingTab) {
                this.tabMenuList.push(tab);
            }
            // 可选：将此选项卡设置为活动选项卡
            this.activeTabs = tab.Title;
        },
        setRuntimeTable(info) {
            this.runtimeTable = info
        },
    },
    persist: {
        key: 'game_store',
        storage: window.localStorage,
        paths: ['language', 'Token', 'routerList','AdminInfo','MenuList','tabMenuList','activeTabs','lang','langRadio','elMenuActive'],
    },

})
