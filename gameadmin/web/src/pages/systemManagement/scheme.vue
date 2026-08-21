<template>
    <div class='gameList'>
        <el-row :gutter="12">
            <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="4"  class="elCard" v-for="item in rhs.List" v-loading="loading">
                <el-card shadow="hover">
                    <div slot="header" class="clearfix topCard">
                        <el-tag type="info" v-if="!item['editPlat']">
                            {{item.Name}}
                        </el-tag>
                        <el-input clearable v-if="item['editPlat']" style="width: 200px" v-model.trim="item.Name" :placeholder="$t('请输入')"></el-input>
                        <el-space>
                            <el-button typeof="primary" circle v-if="!item['editPlat']" @click="item['editPlat'] = true" style="margin-bottom: 0">
                                <el-icon >
                                    <Edit/>
                                </el-icon>
                            </el-button>
                            <el-button typeof="primary" circle v-if="item['editPlat']" @click="clickExcludeGameUp(item)" style="margin-bottom: 0">
                                <el-icon >
                                    <Check/>
                                </el-icon>
                            </el-button>
                            <el-button typeof="danger" circle class="" @click="deleteExcludeGame(item)" style="margin-bottom: 0">
                                <el-icon>
                                    <Delete/>
                                </el-icon>
                            </el-button>
                        </el-space>
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang='ts'>
import {onMounted, ref, reactive, nextTick, watch} from 'vue';
import type {ElInput, FormInstance, FormRules,UploadFile,} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {Game, GameStatus, GameType} from '@/api/gamepb/customer';
import ReturnErrorHandle, {Client} from '@/lib/client';
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia';
import {Delete, Download, Plus, ZoomIn} from '@element-plus/icons-vue'
let uploadImageUrl = ref(null)
const {t} = useI18n()
const store = useStore()
const clickExcludeGameUp = async (item) => {
    let res = await tip.ask(t("确定修改参数吗？"))
    if (res != "ok") {
        item['editPlat'] = false
        return false
    }
    console.log(res);
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.GetExcludeGameUp, {ID:item.ID,Name:item.Name})
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    await queryList()
}

const deleteExcludeGame = async (item) => {
    let res = await tip.ask(t("确定删除？"))
    console.log(res);
    if (res === 'cancel') return
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.GetExcludeGameDel, {ID:item.ID})
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    await queryList()
}
let forms = reactive(<Game>{})
const gameListCollapse = ref('1')
let rhs = reactive({
    List: [],
    dialogVisible: false,
    dialogAction: "modify",
})

let loading = ref(false)
let editPlat = ref(false)

const queryList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.GetexcludeGameAll, {})
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    rhs.List = ReturnErrorHandle(data,'Name')
    console.log(rhs.List);
}

onMounted(queryList);

const addApp = () => {
    rhs.dialogAction = "add"
    rhs.dialogVisible = true
    uploadFileList.value = []
    uploadImageUrl.value = null
    for (let k in forms) {
        forms[k] = ""
    }
}
let gameName = t("游戏列表")

import {AdminInfo} from "@/api/adminpb/info";


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
        align-items: center;
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
.avatar{
    max-width: 200px;
}
.small-font{
    margin: 0;
}
:deep(.el-card__body){
    padding: 15px;
}
</style>
