import {defineStore} from "pinia";
export const useGameStore= defineStore( {
    id:"gameStore",
    state: ()=> ({
        Games:[],                       // 游戏列表
        Manufacturer:[],                // 厂商
        ThemeList:[],                   // 厂商
    }),
    actions: {
        setGameList(games:[]){
            this.$state.Games = games
        },
        setManufacturer(Manufacturer: []){
            this.$state.Manufacturer = Manufacturer
        },
        setThemeList(ThemeList: []){
            this.$state.ThemeList = ThemeList
        }

    },
    persist: {
        key:"GameConfig",
        storage: localStorage,
    }
})
