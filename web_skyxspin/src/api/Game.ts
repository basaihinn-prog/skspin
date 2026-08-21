import {axiosComm} from "@/util/requestComm";
import type {GameInterface} from "@/interface/gameInterface";
import type {FilterInterface} from "@/interface/filterInterface";
import type {CharacterGame} from "@/interface/characterGame";
import {useGameStore} from "@/stores/store/game";
import {generatorGameData} from "@/util/generatorData";


// 获取游戏列表
// @params gameLanguage        语言配置
// @params searchName          名称查询条件
// @params themeFilterList     通过主题与厂商筛选
// @params classification      通过类型
// @params sort                排序

export const GetGameList = async (
    gameLanguage:Record<string, string>,
    searchName: string = "",
    themeFilterList:FilterInterface[],
    classification: string = "ALL",
    sort: number = 1,
) => {


    const { Games } = useGameStore()

    let gameData: GameInterface[] = Games

    // 排序
    if (sort <= 1){

        // 通过名称查询
        if (searchName){
            gameData = gameData.filter(item => item.gameName && item.gameName.toLowerCase().indexOf(searchName.toLowerCase()) != -1)
        }

        // 对厂商类型进行筛选
        if (classification){

            if (classification != "MINI"){
                gameData = gameData.filter(item => item.manufacturerName != "SPRIBE")
            }else{
                gameData = gameData.filter(item => item.manufacturerName == "SPRIBE")

            }
        }



        if (themeFilterList && themeFilterList.length){


            let themeList:string[] = []
            let manufacturerList:string[] = []

            themeFilterList.forEach(item => {

                if (item.type == "manufacturer"){
                    manufacturerList.push(item.label)
                }
                if (item.type == "theme"){
                    themeList.push(item.label)
                }
            })

            if (manufacturerList.length){

                gameData = gameFilterByManufacturer(gameData, manufacturerList)
            }
            if (themeList.length){

                gameData = await gameFilterByTheme(gameData, themeList)
            }

        }


        gameData = gameData.sort((a,b):number=> {
            if (sort == -1) {
                return b.sort - a.sort
            } else if (sort == 1) {
                return a.newSort - b.newSort
            }else{
                return a.sort - b.sort
            }
        })


        gameData = generatorGameData(gameData) as any

    }



    if (sort > 1){

        if (classification != "MINI"){
            gameData = await gameFilterByHotGame(gameData, "hot")
        }else{
            gameData = await gameFilterByHotGame(gameData, "miniHot")
        }


    }

    return gameData
}



// 通过厂商筛选游戏
const gameFilterByHotGame = async (gameData: GameInterface[], filterField: string) => {

    gameData = gameData.filter((item:any) => {
        return item[filterField] != ""
    }).map((item:any)=>({
        ...item,
        sort: Number(item[filterField])
    })).sort((a:any,b:any)=>a.sort - b.sort)


    return gameData

}
// 通过厂商筛选游戏
const gameFilterByManufacturer = (gameData: GameInterface[],manufacturerFilterList:string[]) => {
    return gameData.filter(item => manufacturerFilterList.indexOf(item.manufacturerName) != -1)

}

// 通过主题筛选游戏
const gameFilterByTheme = async (gameData: GameInterface[],themeFilterList:string[]) => {



    let themeGameList = gameData.filter(item=> item[themeFilterList[0].toLowerCase()] != "").map(item => ({
        ...item,
        sort: item[themeFilterList[0].toLowerCase()]
    }))

    // ======= 数据格式化 ======

    let characterMap: Record<string, number> = {}
    themeGameList.forEach(item => {
        characterMap[`${item.manufacturerName}-${item.id}`] = item.sort
    })

    gameData = gameData.filter(item => {
        item.sort = characterMap[`${item.manufacturerName}-${item.id}`]
        return characterMap[`${item.manufacturerName}-${item.id}`]
    })


    return gameData

}
