import type {GameInterface} from "@/interface/gameInterface";
import {axiosComm} from "@/util/requestComm";
import {useLanguageStore} from "@/stores/store/langageStore";
import {useGameStore} from "@/stores/store/game";
import type {ThemeInterface} from "@/interface/themeInterface";
import {moteAxiosComm} from "@/util/moteRequestComm";
import type {manufacturerInterface} from "@/interface/manufacturer";

export const generatorGameData = (gameList: GameInterface[]): Record<string, GameInterface[]> => {

    let manufacturerRecordData: Record<string, GameInterface[]> = {
        "PG": [],
        "JILI": [],
        "PP": [],
    }

    gameList.forEach(item => {
        if (!manufacturerRecordData[item.manufacturerName]) {
            manufacturerRecordData[item.manufacturerName] = []
        }
        manufacturerRecordData[item.manufacturerName].push(item)
    })

    return manufacturerRecordData
}


let LangeStore: any = null
let GameStore: any = null


export const initStore = () => {
    LangeStore = useLanguageStore()
    GameStore = useGameStore()
}


export const GetSystemAllData = async () => {

    await getLanguageList()
    await getLanguageToStore()
    await getLanguageList()
    await getAboutInfo()
    await getManufacturerList()
    await getGames()
    await getThemeList()
}

const getLanguageToStore = async () => {
    // 获取语言系列表
    const Language = localStorage.getItem("systemLanguage")
    const lang = await axiosComm.get(`/mock/language/${Language}.json`)

    LangeStore.setSystemText(lang.data)
    LangeStore.setLanguage(Language)

}

const getLanguageList = async () => {

    const lang = await axiosComm.get("/mock/language.json")

    LangeStore.setLanguageList(lang.data)

}

const getAboutInfo = async () => {
    const about = await axiosComm.get(`/mock/about.json`)

    LangeStore.setSystemAboutInfo(about.data)

}

// 获取所有厂商
const getManufacturerList = async () => {

    const manufacturer = await moteAxiosComm.get(`/AdminInfo/Demo/getDemoManufacturerList`)
    // 为了约束厂商排序
    let manufacturerSortConstraint = ["PG", "JILI", "PP", "TADA", "JDB", "PLAYSTAR", "SKYXSPIN"]
    manufacturer.data.data = manufacturer.data.data.map((item: any, index: number) => ({
        id: Number(index) + 1,
        label: item.ManufacturerName.toUpperCase(),
        localHref: "",
        sort: manufacturerSortConstraint.indexOf(item.ManufacturerName.toUpperCase()) < 0 ? -1 : Number(manufacturerSortConstraint.indexOf(item.ManufacturerName.toUpperCase())) + 1
    })).sort((a: any, b: any): number => a.sort - b.sort)


    GameStore.setManufacturer(manufacturer.data.data)

}

// 获取所有游戏
const getGames = async () => {

    const games = await moteAxiosComm.get(`/AdminInfo/Demo/getGameList`)
    GameStore.setGameList(games.data.data)

}


// 获取所有厂商
const getThemeList = async () => {

    const theme = await moteAxiosComm.get(`/AdminInfo/Demo/getDemoThemeList`)

    theme.data.data = theme.data.data.map((item: any) => ({
        id: Number(item.id),
        label: item.value,
        sort: Number(item.sort)
    }))
    GameStore.setThemeList(theme.data.data)

}
