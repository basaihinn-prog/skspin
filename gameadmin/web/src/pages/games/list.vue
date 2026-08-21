<template>
    <div class='gameList'>
        <div class="searchView">
            <el-space wrap>
                <el-button type="primary" @click="addApp">{{ $t('添加游戏') }}</el-button>
                <Gamelist_container
                    :hase-manufacturer="true"
                    :hase-all="true"
                    :defaultGameEvent="defaultGameEvent"
                    @select-operator="selectGameList"
                />
                <uploadExcel request-uri="AdminInfo/uploadGame">
                    <el-button>上传</el-button>
                </uploadExcel>

            </el-space>
        </div>

        <el-row :gutter="12" style="margin-top: 20px">
            <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="6" class="elCard" v-for="item in gameData"
                    v-loading="loading">
                <el-card shadow="hover">
                    <el-row>
                        <el-col :span="4" style="display: flex;align-items: center;">
                            <el-avatar :size="50" :src="item.IconUrl"/>
                        </el-col>

                        <el-col :span="20">
                            <div slot="header" class="clearfix topCard">
                                <div>
                                    <span>{{ $t(item.Name) }}</span>
                                    <p class="topCardId">( {{ item.ID }} )</p>
                                </div>
                                <div>
                                    <el-space style="margin-right: -8px">
                                        <el-link type="info" @click="copyUrl(item.IconUrl)">
                                            <el-icon>
                                                <CopyDocument/>
                                            </el-icon>
                                        </el-link>

                                        <!--                                       IconUrl-->
                                        <el-button circle class="small-font" @click="editPlat(item)"
                                                   style="box-shadow: none">
                                            <el-icon @click="editPlat(item)">
                                                <Edit/>
                                            </el-icon>
                                        </el-button>
                                        <!--                                       IconUrl-->

                                    </el-space>
                                </div>

                            </div>
                            <div class="bottom clearfix">
                                <el-tag effect="plain">
                                    {{ $t(GameClass[item.Type]) }}
                                </el-tag>
                                <el-tag effect="plain" v-for="(i,key) in item.ExtraData">
                                    <template v-if="key === 'boyibo_balance'">
                                        {{ $t('余额') }}：{{ ut.toNumberWithComma(i) }}
                                    </template>
                                </el-tag>

                                <div>
                                    <el-tag
                                            :type="item.Status?item.Status===1?'warning':'info': 'success'"
                                            effect="light" round>
                                        {{ StatusMap[item.Status] }}
                                    </el-tag>
                                    <el-button type="danger" size="small"
                                               @click="gameDel({'ID':item.ID})"
                                               style="margin-bottom: 0;margin-left: .5rem"
                                               :icon="Delete" v-if="item.Status===3" circle/>
                                </div>
                            </div>
                        </el-col>
                    </el-row>

                </el-card>
            </el-col>
        </el-row>
        <el-dialog v-model="rhs.dialogVisible" :title="rhs.dialogAction == 'add' ? $t('添加') : $t('修改')" class="file-container"
                   destroy-on-close
                   @close="initGameForm"
                   :align-center="true"
                   :width="store.viewModel === 2 ? '85%' : '950px'">

            <el-form ref="formRef" :model="forms" :rules="formRule" label-position="top">

                <el-form-item :label="`${$t('Bet是否可修改')}:`">

                    <el-switch v-model="forms.ChangeBetOff"
                               size="default"
                               :active-value="1"
                               :inactive-value="0"
                    />
                </el-form-item>

                <el-row :gutter="18">
                    <el-col :span="12">

                        <el-form-item :label="`${$t('ID')}:`" prop="GameId">
                            <el-input v-model.trim="forms.GameId" clearable
                                      :placeholder="$t('请输入')" size="default"/>
                        </el-form-item>
                        <el-form-item :label="`${$t('游戏ID')}:`" prop="ID">
                            <el-input v-model.trim="forms.ID" clearable
                                      :disabled="rhs.dialogAction != 'add' "
                                      :placeholder="$t('请输入')" size="default"/>
                        </el-form-item>


                        <el-form-item :label="`${$t('厂商名称')}:`" prop="ManufacturerName">
                            <el-input v-model.trim.number="forms.ManufacturerName" clearable :placeholder="$t('请输入')"
                                      size="default"/>
                        </el-form-item>

                        <template v-for="(game, index) in gameConfig" :key="index">


                            <el-form-item :label="`${$t('游戏名称')}(${ index }):`">
                                <el-input v-model="game.GameName"></el-input>
                            </el-form-item>
                        </template>


                    </el-col>
                    <el-col :span="12">

                        <el-form-item :label="`${$t('类型')}:`" size="default" prop="Type">
                            <el-select v-model="forms.Type" :placeholder="$t('请选择')">
                                <el-option v-for='(value, key) in GameClass' :label="value"
                                           :value="parseInt(key as unknown as string)"/>
                            </el-select>
                        </el-form-item>
                        <el-form-item :label="`${$t('状态')}:`" prop="Status">
                            <el-select v-model="forms.Status" :placeholder="$t('请选择')">
                                <el-option v-for='(value, key) in StatusMap' :label="value"
                                           :value="parseInt(key as unknown as string)"/>
                            </el-select>
                        </el-form-item>

                        <el-row :gutter="18">

                            <template v-for="(game, index) in gameConfig" :key="index">

                                <el-col :span="8">
                                    <el-form-item :label="`${$t('游戏图标')}(${ index }):`" class="gamePage">

                                        <div v-if="game.Icon" style="width: 81px;height: 81px;position: relative">
                                            <el-icon style="position:absolute;right: -5px;top: -5px;z-index: 999" @click="reupload(index)"><CircleCloseFilled /></el-icon>
                                            <img :src="`${protocol}//${baseUri}${game.Icon}`" class="avatar"  style="width: 81px;height: 81px"/>
                                        </div>
                                        <uploadComponent @uploadFile="uploadFileRequest($event, index)" :upload-file-keys="index" v-else/>

                                    </el-form-item>
                                </el-col>
                            </template>
                        </el-row>
                    </el-col>
                </el-row>

            </el-form>

            <template #footer>
                <span class="dialog-footer">
                    <el-button v-if="rhs.dialogAction == 'add'" type="primary" @click="onAddPlat">{{
                        $t('添加')
                        }}</el-button>
                    <el-button v-if="rhs.dialogAction == 'modify'" type="primary" @click="onModifyPlat">{{
                        $t('修改')
                        }}</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick, watch, Ref} from 'vue';
import type {ElInput, FormInstance, FormRules, UploadFile,} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {Game, GameStatus, GameType} from '@/api/gamepb/customer';
import {AdminGameCenter} from '@/api/gamepb/admin';
import {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import {Delete, Download, Plus, ZoomIn} from '@element-plus/icons-vue'
import {excel} from "@/lib/excel";
import { Upload } from "@/api/comm"
import Papa from "papaparse";
let uploadImageUrl = ref(null)
const {t} = useI18n()
const store = useStore()
const formRef = ref<FormInstance>()
const actionUrl = "http://11.com"
const StatusMap: JsMap<GameStatus, string> = {
    [GameStatus.Open]: t('正常'),
    [GameStatus.Maintenance]: t('维护'),
    [GameStatus.Hide]: t('隐藏'),
    [GameStatus.Close]: t('关闭'),
};

const formRule = reactive<FormRules>({
    ID: [{required: true, message: t('请填写ID'), trigger: 'blur'}],
    GameId: [{required: true, message: t('请填写ID'), trigger: 'blur'}],
    LineNum: [{required: true, message: t('请填写游戏线数'), trigger: 'blur'}],
    Bet: [{required: true, message: t('请填写默认Bet'), trigger: 'blur'}],
    Name: [{required: true, message: t('请填写名称'), trigger: 'blur'}],
    Type: [{required: true, message: t('请选择类型'), trigger: 'change'}],
    Status: [{required: true, message: t('请选择状态'), trigger: 'change'}],
})


let forms = reactive(<Game>{})
let gameConfig = ref({})

const queryForm = ref({
    GameId: "",
    Maintenance: "",
})
let rhs = ref({
    List: [],
    dialogVisible: false,
    dialogAction: "modify",
})

const gameData:Ref<Game[]> = ref(<Game[]>[])

let loading = ref(false)
let newList = ref(null)

const {lang, language} = store

const defaultGameEvent = ref({})
const selectGameList = (value) => {

    if (value.gameData){

    }
    queryForm.value.GameId = value.gameData

    if (value.manufacturer || value.manufacturer == null){

        queryForm.value.Maintenance = value.manufacturer
    }

    queryList()

}
const initLanguageList = async () => {

    let [data, err] = await Client.Do(AdminConfigFile.LoadConfig, {FileName: "lang.csv"})

    let arr: any[] = Papa.parse(data.Content).data

    for (let key in arr[0]) {

        let item = arr[0][key]
        if (item == "id"){
            item = "zh"
        }

        gameConfig.value[item] = {
            GameName: "",
            Icon: ""
        }


    }


}
initLanguageList()


const queryList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminGameCenter.GameList, queryForm.value)
    loading.value = false
    if (err) {
        return tip.e(err)
    }




    gameData.value = data.List



}

const protocol = window.location.protocol
const baseUri = window.location.host


onMounted(queryList);


const copyUrl = (url: string) => {
    copy(url)
    tip.s(t("复制成功"))
}
const sortGamesByName = (list = []) => {
    return list.sort((a, b) => {
        const nameA = a.Name.toUpperCase(); // Convert names to uppercase for comparison
        const nameB = b.Name.toUpperCase();
        return nameA.localeCompare(nameB); // Use localeCompare for proper sorting
    });
};
const gameListAutocomplete = ref('')
const restaurants = ref([])
const querySearch = (queryString: string, cb: any) => {
    const results = queryString
        ? restaurants.value.filter(createFilter(queryString))
        : restaurants.value
    // call callback function to return suggestions
    cb(results)

}
const createFilter = (queryString: string) => {
    return (restaurant) => {
        return (
            restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
        )
    }
}
const handleSelect = (item) => {
    gameListAutocomplete.value = item.Name
    item ? rhs.value.List = newList.value.filter(list => list.ID === item.ID) : rhs.value.List = newList.value

}
const handleChange = (item) => {
    queryList()
    // item ? rhs.value.List = newList.value.filter(list => list.ID === item.ID) : rhs.value.List = newList.value
}
const addApp = () => {
    rhs.value.dialogAction = "add"
    rhs.value.dialogVisible = true
    uploadFileList.value = []
    uploadImageUrl.value = null
    for (let k in forms) {
        forms[k] = ""
    }
}
let gameName = t("游戏列表")
const downLoad = () => {
    let arr = rhs.value.List.map(list => {
        return {
            ...list,
            Type: TypeMap[list.Type]
        }
    })
    excel.dump(arr, gameName, [
        {key: "Name", name: t("游戏名称")},
        {key: "Name", name: t("英文名称")},
        {key: "ID", name: t("游戏ID")},
        {key: "Type", name: t("游戏类型")},
        {key: "", name: t("Free Spin")},
        {key: "", name: t("Feature Buy")},
        {key: "IconUrl", name: t("icons")}
    ])
}
const editPlat = async (row) => {
    rhs.value.dialogAction = 'modify'
    rhs.value.dialogVisible = true
    uploadFileList.value = []
    await initLanguageList()
    uploadImageUrl.value = null
    for (let k in row) {
        forms[k] = row[k]
    }


    let gameConfigAbout = {}
    for (let rowKey in gameConfig.value) {
        if (!row.GameNameConfig){
            row.GameNameConfig = {}
        }
        if(!row.GameNameConfig[rowKey.toLowerCase()]){

            row.GameNameConfig[rowKey.toLowerCase()] = {
                GameName:"",
                Icon:""
            }
        }
        gameConfigAbout[rowKey.toLowerCase()] = {
            GameName: row.GameNameConfig[rowKey.toLowerCase()].GameName,
            Icon: row.GameNameConfig[rowKey.toLowerCase()].Icon,
        }

    }


    gameConfig.value = gameConfigAbout
}
const uploadRef = ref()
const uploadFileList = ref([])

const initGameForm = () => {
    formRef.value.resetFields()
    initLanguageList()
}


const onAddPlat = async () => {


    let res = await tip.ask(t("确定修改参数吗？"))

    if (res != "ok") return false
    if (!formRef) return




    forms.GameNameConfig = gameConfig.value
    await formRef.value.validate(async (valid, fields) => {
        if (valid) {
            let [_, err] = await Client.Do(AdminGameCenter.AddGame, forms)
            if (err) {
                return tip.e(err)
            }

            tip.s(t("成功"))
            rhs.value.dialogVisible = false
            queryList()
        }
    })
}

const gameDel = async (item) => {
    // /AdminInfo/DeleteGame
    let [_, err] = await Client.Do(AdminInfo.DeleteGame, item)
    if (err) {
        return tip.e(err)
    }

    tip.s(t("修改成功"))
    await queryList()

}

import ut from '@/lib/util'
import {AdminInfo} from "@/api/adminpb/info";
import UploadComponent from "@/components/customTable/uploadComponent.vue";
import {AdminConfigFile} from "@/api/adminpb/json";
import copy from "copy-to-clipboard";
import UploadExcel from "@/pages/games/uploadExcel.vue";
import Gamelist_container from "@/components/gamelist_container.vue";
import tr from "element-plus/es/locale/lang/TR";
import {GameClass} from "@/enum";

const onModifyPlat = async () => {
    let res = await tip.ask(t("确定修改参数吗？"))
    if (res != "ok") return false


    forms.GameNameConfig = gameConfig.value

    forms.DefaultBet = Number(forms.DefaultBet)
    let [data, err] = await Client.Do(AdminGameCenter.ModifyGame, forms)
    if (err) {
        return tip.e(err)
    }
    tip.s(t("修改成功"))
    rhs.value.dialogVisible = false
    await queryList()
}
const toBase64 = async (url) => {
    const response = await fetch(url);
    const arrayBuffer = await response.arrayBuffer();
    const base64 = btoa(
        String.fromCharCode.apply(null, new Uint8Array(arrayBuffer))
    );
    return `data:image/png;base64,${base64}`;
}

const reupload = (key) => {
    gameConfig.value[key].Icon = ""
}

// Client.Do(AdminGameCenter.Clear, {} as any)
const uploadFileRequest = async (uploadMap, lange) => {
    let formData = new FormData();
    formData.append('file', uploadMap.file);
    formData.append("game_id", forms.ID);
    formData.append("language", lange);


    const [response, err] = await Client.Do(Upload.UploadFile, formData )

    if (err){
        return tip.e(t("上传失败"))
    }

    gameConfig.value[lange].Icon = response.url

}

</script>
<style scoped lang='scss'>
.elCard {
  border: none;
  margin-bottom: 1rem;

  .el-card {
    border-radius: .5rem;
  }

  .topCard {
    //padding-left: 1.5rem;
    display: flex;
    justify-content: space-between;

    span {
      margin-right: .5rem;

    }

    .topCardId {
      font-size: 10px;
      color: #999999;
      margin-bottom: .2rem;
    }

  }

  .bottom {
    display: flex;
    justify-content: space-between;
    width: 100%;
    position: relative;
    //padding-left: 1.5rem;
    font-size: 13px;

    div {
      padding: 0 3px;
      line-height: 24px;
    }

  }
}

.avatar {
  max-width: 200px;
}

.small-font {
  margin: 0;
}

.avatar-uploader .avatar {
  width: 80px;
  height: 80px;
  display: block;

}

</style>
