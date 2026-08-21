<template>
    <el-row :gutter="24">

        <el-col :lg="4" :sm="6" :xs="8" :span="8" v-for="(item, index) in props.gameList" :key="index" class="game_item" :class="`game_item_${props.componentType}_${item.manufacturer}_${item.id}_${index}`" @click="startGame(item)">

            <div class="imageContent" :class="`imageContent_${props.componentType}_${item.manufacturer}_${item.id}_${index}`" :style="{minHeight:imageHeight}" >

                <el-image class="game_img"
                          :src="`${item.icon}?_=${new Date().getDate()}`" fit="fill" :lazy="true">
                    <template #placeholder>
                        <div :style="{minHeight:imageHeight}"></div>
                    </template>
                </el-image>

                <div class="runStatus_model" :style="{minHeight:imageHeight}">
                    <img :src="Number(item.runStatus) == 1 ? `${HOST}${aboutInfo.runStatus?.enabled}` : `${HOST}${aboutInfo.runStatus?.disabled}`"
                         :style="{width: item.runStatus == 1 ? '56%' : '94%'}">
                </div>
            </div>
            <div class="game_name">{{item.gameName}}</div>

        </el-col>
    </el-row>
</template>

<script setup lang="ts">

    import type {GameInterface} from "@/interface/gameInterface";
    import type { manufacturerInterface } from "@/interface/manufacturer";
    import type {CharacterGame} from "@/interface/characterGame";
    import type {AboutInterface, RunStatusInterface} from "@/interface/aboutInterface";
    import {onMounted, onUnmounted, Ref, ref, watch} from "vue";
    import {ElLoading, ElNotification} from "element-plus";
    import failImg from "@/assets/failImg.png"
    import loadding from "@/assets/loadding.webp"
    import {LanguageConfigInterface} from "@/interface/languageConfigInterface";
    import {moteAxiosComm, externalMoteAxiosComm} from "@/util/moteRequestComm";
    import router from "@/router";
    import {SetOperator} from "@/util/util";
    import {useLanguageStore} from "@/stores/store/langageStore";
    import {storeToRefs} from "pinia";

    const imageContentRef = ref("")
    const imageHeight = ref("")
    const LanguageStore = useLanguageStore()
    const {aboutInfo, systemText }: {
        aboutInfo:Ref<AboutInterface>
        systemText:Ref<LanguageConfigInterface>
    } = storeToRefs(LanguageStore)



    const props = defineProps({
        gameList:{
            type: Array<GameInterface>,
            default: []
        },
        componentType: {
            type: String,
            default: ""
        }
    })



    const HOST:string =  window.location.protocol + "//" + window.location.host + import.meta.env.VITE_SOURCE_BASE




    const setImageContentHeight = (className: string) => {

        let dom:HTMLCollectionOf<Element> = document.getElementsByClassName(className)

        if (dom[0] && dom[0].clientWidth && dom[0].clientWidth != undefined && dom[0].clientWidth != 0){

            localStorage.setItem("imageHeight", dom[0].clientWidth.toString())
        }

        return dom[0] ? dom[0].clientWidth : 0
    }


    const isPhone = () => {
        const userAgentInfo = navigator.userAgent;
        const mobileAgents = ["Android", "iPhone", "SymbianOS", "Windows Phone", "iPad", "iPod"];
        const mobileFlag = mobileAgents.some((mobileAgent) => {
            return userAgentInfo.indexOf(mobileAgent) > 0;
        });

        return mobileFlag;
    }



    const initUserId = (callback) => {
        let UserID = localStorage.getItem("userId")

        if (!UserID){
            const randomSource = "QWERTYUIOPASDFGHJKLZXCVBNMwertyuiopasdfghjklzxcvbnm"

            let result = ""
            for (let i = 0; i < 8; i++) {

                result += randomSource.charAt(Math.floor(Math.random() * randomSource.length));

            }

            result += `-${new Date().getTime()}-DemoUser`



            localStorage.setItem("userId",result)
            UserID = result
        }


        moteAxiosComm.request({
            url: "/v1/player/create",
            method: "POST",
            data:{
                UserID
            },
        }).then(res=>{
            if(res.error){
                ElNotification({
                    title: systemText.value.hitTitle,
                    message: systemText.value.hitDescription,
                    type: 'error',
                    showClose: false,
                })
                return
            }

            callback()
        })
    }
    const startGame = (game: GameInterface) => {

        if (game.runStatus == 0){
            ElNotification({
                title: systemText.value.hitTitle,
                message: systemText.value.hitDescription,
                type: 'error',
                duration:2000,
            })
            return;
        }



        let param = SetOperator()
        if (param){





            initUserId(() => {

                const loading = ElLoading.service({
                    lock: true,
                    text: 'Loading',
                    background: 'rgba(0, 0, 0, 0.7)',
                })



                const ap = localStorage.getItem("app_id")




                if (game.manufacturerName.toLocaleLowerCase() == "playstar" || game.manufacturerName.toLocaleLowerCase() == "skyxspin") {
                    externalMoteAxiosComm.request({
                    url: "/v1/game/launch",
                    method: "POST",
                    data: {
                        UserID: localStorage.getItem("userId"),
                        GameID: game.manufacturerName.toLocaleLowerCase() + "_" + game.id.toString(),
                        Platform: "desktop",
                        Language: localStorage.getItem("systemLanguage") || "en"
                    }
                    }).then(async res => {
                        if (res.error) {


                            if (res.error == "Line provider does not allow access") {

                                ElNotification({
                                    title: systemText.value.hitTitle,
                                    message: res.data.error,
                                    duration:2000,
                                    type: 'error',

                                })

                            }


                            ElNotification({
                                title: systemText.value.hitTitle,
                                message: systemText.value.hitDescription,
                                duration:2000,
                                type: 'error',

                            })
                            loading.close()
                            return;
                        }
                        var u = navigator.userAgent;
                        var isiOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端

                        loading.close()


                        localStorage.setItem("launchedScroll", window.pageYOffset.toString())
                        // let Uri = `${res.data.Url}&game_name=${game.gameName}&mu=${game.manufacturerName}&_=${new Date().getTime()}`

                        if (isiOS) {

                            window.location.href = res.data.Url

                        } else {

                            window.open(res.data.Url)
                        }
                        // let Uri = `https://m-pg.ppply.site/launch.html?game_uri=${encodeURIComponent(res.data.data.Url)}&AppId=${AppId}&GameName=${game.gameName}`
                        if (ap == 'faketrans') {

                            let ac = await moteAxiosComm.request({
                                url: "/SetDemoPlayerRTP",
                                method: "POST",
                                data: {
                                    DemoUserId: localStorage.getItem("userId"),
                                    GameID: game.manufacturerName.toLocaleLowerCase() + "_" + game.id.toString(),
                                    ContrllRTP: Number(localStorage.getItem("RTPValue")) || 96
                                },
                            })

                        }

                    })
                } else {
                    moteAxiosComm.request({
                    url: "/v1/game/launch",
                    method: "POST",
                    data: {
                        UserID: localStorage.getItem("userId"),
                        GameID: game.manufacturerName.toLocaleLowerCase() + "_" + game.id.toString(),
                        Platform: "desktop",
                        Language: localStorage.getItem("systemLanguage") || "en"
                    }
                    }).then(async res => {
                        if (res.error) {


                            if (res.error == "Line provider does not allow access") {

                                ElNotification({
                                    title: systemText.value.hitTitle,
                                    message: res.data.error,
                                    duration:2000,
                                    type: 'error',

                                })

                            }


                            ElNotification({
                                title: systemText.value.hitTitle,
                                message: systemText.value.hitDescription,
                                duration:2000,
                                type: 'error',

                            })
                            loading.close()
                            return;
                        }
                        var u = navigator.userAgent;
                        var isiOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端

                        loading.close()


                        localStorage.setItem("launchedScroll", window.pageYOffset.toString())
                        let Uri = `${res.data.Url}&game_name=${game.gameName}&mu=${game.manufacturerName}&_=${new Date().getTime()}`

                        if (isiOS) {

                            window.location.href = Uri

                        } else {

                            window.open(Uri)
                        }
                        // let Uri = `https://m-pg.ppply.site/launch.html?game_uri=${encodeURIComponent(res.data.data.Url)}&AppId=${AppId}&GameName=${game.gameName}`
                        if (ap == 'faketrans') {

                            let ac = await moteAxiosComm.request({
                                url: "/SetDemoPlayerRTP",
                                method: "POST",
                                data: {
                                    DemoUserId: localStorage.getItem("userId"),
                                    GameID: game.manufacturerName.toLocaleLowerCase() + "_" + game.id.toString(),
                                    ContrllRTP: Number(localStorage.getItem("RTPValue")) || 96
                                },
                            })

                        }

                    })
                }
                
               


            })


        }
    }



    onMounted(()=>{
        imageHeight.value = setImageContentHeight("imageContent") + "px"

        window.addEventListener("resize", ()=> {
            let width = setImageContentHeight("imageContent")
            imageHeight.value = width.toString() + "px"
        })
    })

    onUnmounted(()=>{
        window.addEventListener("resize",null)
    })


</script>

<style scoped>

    .game_name{
        width: 100%;
        height: 30px;
        line-height: 20px;
        margin-bottom: 15px;
        text-align: center;
        font-size: 15px;
        color: #ffffff;
        position: relative;
        margin-top: 10px;
    }
    .game_img{
        max-width: 228px;
        height: auto;
        overflow: hidden;
        position: relative;
        margin: 0 auto;
    }
    .game_img_load{
        width: 58%;
    }
    .imageContent{
        width: 100%;
        height: auto;
        min-height: 20px;
        border-radius: 10px;
        display: flex;
        align-items: center;
    }
    .imageContent:hover .runStatus_model{
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .image-slot{
        max-width: 228px;
        position: absolute;
        width: 100%;
        height: 100%;
        display: flex;
        background: rgba(0,0,0,0);
        align-items: center;
        justify-items: center;
    }



    .runTimeModal{
        width: 100%;
        height: 100%;
        position: absolute;
        left: 0;
        top: 0;
        background: rgba(0,0,0, .8);
    }
    .game_item{
        width: 100%;
        height: auto;
        display: block;
        position: relative;
        margin-bottom: 22px;
    }
    .runStatus_model{
        background-color:rgba(25,25,25,0.7);
        width: 100%;
        height: calc(100% - 50px);
        position: absolute;
        left: 0;
        top: 0;
        display: none;
    }
    .runStatus_model>img{
        height: auto;
    }

</style>

<style>
.game_item:hover .el-image__inner{
    transform: scale(1.1);
    transition: 0.5s ease-in-out;
}

.game_item .el-image__inner{
    transform: scale(1);
    transition: 0.5s ease-in-out;
}
</style>
