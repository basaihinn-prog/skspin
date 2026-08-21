import {defineStore} from "pinia";
import type {LanguageInterface} from "@/interface/languageInterface";
import type {LanguageConfigInterface} from "@/interface/languageConfigInterface";
import type {AboutInterface} from "@/interface/aboutInterface";


export const useLanguageStore= defineStore( {
    id:"languageStore",
    state: ()=> ({
        language: "",                                        // 当前语言
        systemText: <LanguageConfigInterface>{},             // 系统语言文字
        LanguageList:<LanguageInterface[]>[],                // 语种列表
        aboutInfo: <AboutInterface>{},                       // 试玩站相关信息
    }),
    getters: {
      getLanguage: (state): string => state.language
    },
    actions: {
        // 设置当前语言
        setLanguage(language: string = "en"){
            this.$state.language = language
        },
        // 设置系统语言文字
        setSystemText(systemText: LanguageConfigInterface){
            this.$state.systemText = systemText
        },
        // 设置语种列表
        setLanguageList(languageList: LanguageInterface[]){
            this.$state.LanguageList = languageList
        },
        // 设置试玩站相关信息
        setSystemAboutInfo(aboutInfo: AboutInterface){
            this.$state.aboutInfo = aboutInfo
        },


    },
    persist: {
        key:"languageConfig",
        storage: localStorage,
    }
})
