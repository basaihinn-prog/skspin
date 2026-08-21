import {useStore} from '@/pinia/index'

// 这些路由无需登录就可以随便进的
let noLimitRouters = [
    // 登陆
    {path: '/', redirect: "/dashboard"},
    {path: '/login', component: () => import('@/pages/Login.vue'),},
    {path: '/404', component: () => import('@/pages/404.vue'),},
    {
        path: "/:pathMatch(.*)",
        redirect: "/404"
    },
    {
        path: "/", component: () => import('@/pages/layout.vue'), redirect: '/dashboard',
        children: [
            {
                path: "/dashboard",
                component: () => import('@/pages/runtime/runtime_operator.vue'),
                meta: {title: "游戏报表", parent: "游戏数据"}
            },
            {path: "/config", component: () => import('@/pages/systemManagement/config/json.vue'), meta: {title: "配置文件"}},
            //接口文档-单一钱包
            {path: "/singleWallet", component: () => import('@/pages/mdPreview/singleWallet.vue'),},
            // 接口文档-转账钱包
            {path: "/transferWallet", component: () => import('@/pages/mdPreview/transferWallet.vue'),},
            //后台管理
            {
                path: '/menuList',
                component: () => import('@/pages/systemManagement/menuList.vue'),
                meta: {title: "功能菜单维护", parent: "后台管理"}
            },
            {
                path: '/scheme',
                component: () => import('@/pages/systemManagement/scheme.vue'),
                meta: {title: "菜单管理", parent: "后台管理"}
            },
            {
                path: '/groupList',
                component: () => import('@/pages/systemManagement/groupList.vue'),
                meta: {title: "功能权限维护", parent: "后台管理"}
            },
            {
                path: '/adminList',
                component: () => import('@/pages/systemManagement/adminList.vue'),
                meta: {title: "用户维护", parent: "后台管理"}
            },
            {
                path: '/mock_appAnalysis',
                component: () => import('@/pages/mockpage/app.vue'),
                meta: {title: "数据统计", parent: "产品数据"}
            },
            {
                path: '/newUserRTP',
                component: () => import('@/pages/systemManagement/merchants/newUserRTP.vue'),
                meta: {title: "新玩家RTP保护", parent: "商户管理"}
            },
            {
                path: '/preAepositAmountLog',
                component: () => import('@/pages/systemManagement/merchants/preAepositAmountLog.vue'),
                meta: {title: "预存额日志", parent: "商户管理"}
            },
            {
                path: '/multiLanguage',
                component: () => import('@/pages/systemManagement/multiLanguage.vue'),
                meta: {title: "多语言管理", parent: "后台管理"}
            },
            // { path: '/notification', meta:{ title:"公告"}, component: () => import('@/pages/systemManagement/notification.vue') },

            {
                path: '/appList',
                component: () => import('@/pages/app/appList.vue'),
                meta: {title: "产品列表", parent: "产品管理"}
            },
            {path: '/IPAccessQuery', component: () => import('@/pages/app/IPAccessQuery.vue')},
            {path: '/serverStatus', component: () => import('@/pages/systemManagement/serverStatus.vue')},
            {
                path: '/notification',
                meta: {title: "公告"},
                component: () => import('@/pages/systemManagement/pushMessage.vue')
            },


            {
                path: '/rtpMonitor',
                component: () => import('@/pages/player/rtpMonitor/rtpMonitor.vue'),
                meta: {title: "游戏动态调控", parent: "玩家信息"}
            },
            {
                path: '/playerInfoOperator',
                component: () => import('@/pages/player/info.vue'),
                meta: {title: "玩家信息", parent: "玩家管理"}
            },
            {
                path: '/playerRTPData',
                component: () => import('@/pages/player/rtpControl.vue'),
                meta: {title: "玩家RTP控制列表", parent: "玩家管理"}
            },
            {
                path: '/playerRestrictions',
                component: () => import('@/pages/player/playerRestrictions/playerRestrictions.vue'),
                meta: {title: "玩家终身限制", parent: "玩家管理"}
            },
            {
                path: '/platformMaintenance',
                component: () => import('@/pages/riskManagement/platformMaintenance.vue'),
                meta: {title: "平台维护", parent: "风控管理"}
            },
            {
                path: '/operatorApproval',
                component: () => import('@/pages/riskManagement/operatorApproval.vue'),
                meta: {title: "开户审批", parent: "商户管理"}
            },
            {
                path: '/riskManagement',
                component: () => import('@/pages/riskManagement/riskManagement.vue'),
                meta: {title: "平台维护", parent: "风控管理"}
            },
            {
                path: '/systemAlarm/balanceAlarm',
                component: () => import('@/pages/systemAlarm/balanceAlarm.vue'),
                meta: {title: "余额预警", parent: "风控管理"}
            },
            {
                path: '/systemAlarm/returnRateAlarm',
                component: () => import('@/pages/systemAlarm/returnRateAlarm.vue'),
                meta: {title: "回报率预警", parent: "风控管理"}
            },
            {
                path: '/systemAlarm/transferRateAlarm',
                component: () => import('@/pages/systemAlarm/transferAlarm.vue'),
                meta: {title: "转账预警", parent: "风控管理"}
            },
            {path: '/playerList', component: () => import('@/pages/player/list.vue')},
            {
                path: '/playerAnalysis',
                component: () => import('@/pages/analysis/player.vue'),
                meta: {title: "玩家数据", parent: "玩家管理"}
            },
            {
                path: '/gameAnalysis',
                component: () => import('@/pages/analysis/game.vue'),
                meta: {title: "每日游戏详情", parent: "数据统计"}
            },
            {
                path: '/appAnalysis',
                component: () => import('@/pages/analysis/app.vue'),
                meta: {title: "产品数据", parent: "数据统计"}
            },
            {
                path: '/gameUserRetained',
                component: () => import('@/pages/analysis/gameUserRetained/gameUserRetained.vue'),
                meta: {title: "游戏留存", parent: "数据统计"}
            },
            {
                path: '/ranking',
                component: () => import('@/pages/analysis/ranking/ranking.vue'),
                meta: {title: "排行榜", parent: "数据统计"}
            },
            {
                path: '/betLog',
                component: () => import('@/pages/analysis/betlist.vue'),
                meta: {title: "下注历史", parent: "数据统计"}
            },
            {path: '/downloadManagement', component: () => import('@/pages/analysis/downloadManagement.vue')},
            {path: '/goldLog', component: () => import('@/pages/analysis/modifylist.vue')},

            // 月末处理
            {
                path: '/statement',
                component: () => import('@/pages/monthEndProcess/statement.vue'),
                meta: {title: "对账单", parent: "月末处理"}
            },
            // 平台汇总 每日平台汇总 每月平台汇总 平台玩家汇总
            {
                path: '/platform_summary',
                component: () => import('@/pages/analysis/platform_summary.vue'),
                meta: {title: "平台汇总", parent: "数据统计"}
            },
            {
                path: '/daily_platform_summary',
                component: () => import('@/pages/analysis/daily_platform_summary.vue'),
                meta: {title: "每日平台汇总", parent: "数据统计"}
            },
            {
                path: '/monthly_platform_summary',
                component: () => import('@/pages/analysis/monthly_platform_summary.vue'),
                meta: {title: "每月平台汇总", parent: "数据统计"}
            },
            {
                path: '/platform_player_summary',
                component: () => import('@/pages/analysis/platform_player_summary.vue'),
                meta: {title: "平台玩家汇总", parent: "数据统计"}
            },
            {
                path: '/betWinStatistics',
                component: () => import('@/pages/analysis/betWinStatistics.vue'),
                meta: {title: "打码量与盈利统计", parent: "数据统计"}
            },

            // 日志管理 LogManagement
            {
                path: '/pendingApprovalLog',
                component: () => import('@/pages/analysis/pendingApprovalLog.vue'),
                meta: {title: "待审批日志-奖池", parent: "个人奖池", parents: "风控管理"}
            },
            {
                path: '/logQuery',
                component: () => import('@/pages/analysis/logQuery.vue'),
                meta: {title: "日志查询-奖池", parent: "个人奖池", parents: "风控管理"}
            },
            {
                path: '/prizePool/pendingApprovalLog',
                component: () => import('@/pages/analysis/prizePool/pendingApprovalLog.vue')
            },
            {path: '/prizePool/logQuery', component: () => import('@/pages/analysis/prizePool/logQuery.vue')},
            // 彩票每期情况统计
            {
                path: '/lotteryStatisticsForEach',
                component: () => import('@/pages/analysis/lotteryStatisticsForEach.vue')
            },
            {
                path: '/operatorMaintenance',
                component: () => import('@/pages/systemManagement/merchants/merchantsManage.vue'),
                meta: {title: "商户列表", parent: "商户管理"}
            },
            {
                path: '/mock_operatorMaintenance',
                component: () => import('@/pages/mockpage/operatorMaintenance.vue'),
                meta: {title: "商户列表", parent: "商户管理"}
            },

            {
                path: '/mock_gameUserRetained',
                component: () => import('@/pages/mockpage/gameUserRetained.vue'),
                meta: {title: "游戏留存", parent: "数据统计"}
            },

            {
                path: '/mock_daily_platform_summary',
                component: () => import('@/pages/mockpage/daily_platform_summary.vue'),
                meta: {title: "游戏留存", parent: "每日平台汇总"}
            },
            {
                path: '/mock_monthly_platform_summary',
                component: () => import('@/pages/mockpage/monthly_platform_summary.vue'),
                meta: {title: "游戏留存", parent: "每月平台汇总"}
            },
            {
                path: '/mock_platform_summary',
                component: () => import('@/pages/mockpage/platform_summary.vue'),
                meta: {title: "游戏留存", parent: "平台汇总"}
            },

            {path: '/lotteryReport', component: () => import('@/pages/analysis/lotteryReport.vue')},
            {path: '/lotteryGuess', component: () => import('@/pages/analysis/lotteryGuess.vue')},
            {path: '/lotterList', component: () => import('@/pages/games/lotterList.vue')},

            // 游戏运行时
            {path: '/runtime_normal', component: () => import('@/pages/runtime/runtime_normal.vue')},
            {
                path: '/runtime_operator',
                component: () => import('@/pages/runtime/runtime_operator.vue'),
                meta: {title: "商户", parent: "数据统计"}
            },
            // 游戏运行时(汇总)
            {
                path: '/runtime_all',
                component: () => import('@/pages/runtime/runtime_all.vue'),
                meta: {title: "汇总", parent: "数据统计"}
            },
            {
                path: '/runtime_dataSummary',
                component: () => import('@/pages/runtime/runtime_dataSummary.vue'),
                meta: {title: "数据汇总", parent: "游戏数据"}
            },
            // 扑克
            {path: '/puke_daily_platform_summary', component: () => import('@/pages/puke/daily_platform_summary.vue')},
            {
                path: '/puke_monthly_platform_summary',
                component: () => import('@/pages/puke/monthly_platform_summary.vue')
            },

            /*{ path: '/runtime', component: () => import('@/pages/runtime/runtime.vue') },
			{ path: '/runtime_niubi', component: () => import('@/pages/runtime/runtime.vue') },
			{ path: '/runtime_baozang', component: () => import('@/pages/runtime/runtime.vue') },
			{ path: '/runtime_zhaocaimao', component: () => import('@/pages/runtime/runtime.vue') },
			{ path: '/runtime_roma', component: () => import('@/pages/runtime/runtime.vue') },
			{ path: '/runtime_romax', component: () => import('@/pages/runtime/runtime.vue') },
			{ path: '/runtime_hilo', component: () => import('@/pages/runtime/runtime_hilo.vue') },
			{ path: '/runtime_xingyunxiang', component: () => import('@/pages/runtime/runtime.vue') },*/
            {path: '/slotsgen_xingyunxiang', component: () => import('@/pages/slotsgen/slotsgen.vue')},
            {path: '/slotsgen_yingcaishen', component: () => import('@/pages/slotsgen/yingcaishen.vue')},
            {path: '/slotsgen_pg_39', component: () => import('@/pages/slotsgen/pg_39.mjs')},
            {path: '/slotsgen_pg_1489936', component: () => import('@/pages/slotsgen/pg_1489936.mjs')},
            {path: '/slotsgen_pg_common', component: () => import('@/pages/slotsgen/pg_common.mjs')},
            {path: '/slotsgen_pg_common_list', component: () => import('@/pages/slotsgen/slotsgen_pg_common.vue')},
            {path: '/slotsgen_jili_common_list', component: () => import('@/pages/slotsgen/slotsgen_jili_common.vue')},
            {path: '/slotsgen_pp_common_list', component: () => import('@/pages/slotsgen/slotsgen_pp_common.vue')},

            {path: '/slotsgen_tuzi', component: () => import('@/pages/slotsgen/tuzi.vue')},
            {path: '/slotsgen_majiang', component: () => import('@/pages/slotsgen/majiang.vue')},
            {path: '/slotsgen_majiang2', component: () => import('@/pages/slotsgen/majiang2.vue')},
            {path: '/slotsgen_jinniu', component: () => import('@/pages/slotsgen/jinniu.vue')},
            {path: '/slotsgen_bonaza', component: () => import('@/pages/slotsgen/bonaza.vue')},
            {path: '/slotsgen_cowBoy', component: () => import('@/pages/slotsgen/cowBoy.vue')},
            {path: '/slotsgen_olympus', component: () => import('@/pages/slotsgen/olympus.vue')},
            {path: '/slotsgen_olympus1000', component: () => import('@/pages/slotsgen/olympus1000.vue')},
            {path: '/slotsgen_sugarRush', component: () => import('@/pages/slotsgen/sugarRush.vue')},
            {path: '/slotsgen_StarlightChristmas', component: () => import('@/pages/slotsgen/starlightChristmas.vue')},
            {path: '/slotsgen_Starlight_Princess', component: () => import('@/pages/slotsgen/starlight.vue')},
            {path: '/slotsgen_Starlight_Princess1000', component: () => import('@/pages/slotsgen/starlight1000.vue')},
            {path: '/game_list', component: () => import('@/pages/games/list.vue'), meta: {title: "游戏列表", parent: "产品管理"}},
        ]
    }
]

import {createRouter, createWebHashHistory, createWebHistory} from "vue-router";
import {Client} from "@/lib/client";
import {AdminGameCenter} from "@/api/gamepb/admin";

let router = createRouter({
    history: createWebHashHistory(),
    // @ts-ignore
    routes: noLimitRouters
})
let update = true
window.addEventListener('online', () => {
    update = true
})  // 网络由异常到正常时触发
window.addEventListener('offline', () => {
    update = false
})   //网络由正常到异常时触发

router.beforeEach(async (to, from, next) => {
    const store = useStore();
    if (store.Token && update) {
        Client.Do(AdminGameCenter.GameList, {} as any).then(res => {
            store.setGameList(res[0].List)
        })
        localStorage.setItem("language", store.language)
        if (to.path != "/login"){
            await store.setTipsMap()
        }

        next() // 已经登录。随意进入
    } else {
        let one = noLimitRouters.find((item) => {
            return to.path == item.path
        })
        if (one) {
            next()
        } else {
            store.setToken('')
            next({path: '/login'})
        }
    }
});

/*router.beforeResolve((to, from, next) => {
    const store = useStore();
    if (store.Token) {
        const duration = performance.navigation.duration;
        console.log(duration);
        if (duration > 1000) {
            // 跳转缓慢
            alert('----')
        }
        next() // 已经登录。随意进入
    }

});*/
export default router
