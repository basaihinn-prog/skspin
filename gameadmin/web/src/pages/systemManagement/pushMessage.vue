<template>
    <div class='groupList' v-loading="loadingUpdate">
        <div class="group">
            <el-input
                v-model="textarea"
                :disabled="cantEdit"
                style="width: 600px"
                :autosize="{ minRows: 6, maxRows: 20 }"
                type="textarea"
                :placeholder="$t('请粘贴您的文本')"
            />
            <div class="flex">
                <el-button style="margin-top: 12px" @click="cantEdit = false">修改</el-button>
                <el-button style="margin-top: 12px" :disabled="!textarea" @click="getSaveMsg4GameUpdateMsg()">保存
                </el-button>
            </div>
            <div class="gap">
                <el-checkbox-group v-model="checkList">
                    <el-checkbox v-for="item in operatorData" :label="item.Id">
                        {{ item.Name }}
                    </el-checkbox>
                </el-checkbox-group>
            </div>
            <div class="flex">
                <el-button @click="getSaveOperator4GameUpdateMsg" :disabled="!checkList.length">完成</el-button>
                <el-button @click="getSendMsg4GameUpdate" :disabled="!checkList.length || !textarea">发送</el-button>
                <el-button @click="getResult4GameUpdateMsg" :disabled="!checkList.length || !textarea">查看结果</el-button>
            </div>
           <div class="gap gap-more" v-if="FailOperatorID && FailOperatorID.length">
               <p v-for="item in FailOperatorID">
                   {{item}}
               </p>
           </div>
        </div>
    </div>

</template>

<script setup lang="ts">
import {ref, onMounted, reactive} from "vue";
import {Client} from "@/lib/client";
import {AdminInfo} from "@/api/adminpb/info";
import {tip} from "@/lib/tip";
import {useI18n} from "vue-i18n";

const textarea = ref('')
const cantEdit = ref(true)
let operatorData = ref([])
let checkList = ref([])
let loadingUpdate = ref(false)
let FailOperatorID = ref([])
let getConfigData = ref(null)
let {t} = useI18n()
onMounted(() => {
    getStatus4GameUpdateMsg()
    getConfig4GameUpdateMsg()
    operatorList()
});
let operatorParam = reactive({
    PageIndex: 1,
    PageSize: 10000
})

const getConfig4GameUpdateMsg = async () => {
    let [data, err] = await Client.Do(AdminInfo.getConfig4GameUpdateMsg, {})
    if (err) {
        return tip.e(err)
    }
    console.log(data);
    getConfigData.value = data
    textarea.value = getConfigData.value.Msg
    checkList.value = getConfigData.value.DestOperatorID
    console.log(checkList.value);
}

const getSaveOperator4GameUpdateMsg = async () => {
    let [data, err] = await Client.Do(AdminInfo.getSaveOperator4GameUpdateMsg, {DestOperatorID: checkList.value})
    if (err) {
        return tip.e(err)
    }
    tip.s(t('商户方案保存成功'))
}
const getSaveMsg4GameUpdateMsg = async () => {
    cantEdit.value = true
    let [data, err] = await Client.Do(AdminInfo.getSaveMsg4GameUpdateMsg, {Msg: textarea.value})
    if (err) {
        return tip.e(err)
    }
    tip.s(t('消息保存成功'))
}

const getSendMsg4GameUpdate = async () => {
    loadingUpdate.value = true
    if (!textarea.value.trim()) {
        tip.e(t('消息不能为空'))
        loadingUpdate.value = false
        return
    }
    if (!checkList.value.length) {
        tip.e(t('请选择商户'))
        loadingUpdate.value = false
        return
    }
    let [data, err] = await Client.Do(AdminInfo.getSendMsg4GameUpdate,
        {
            Msg: textarea.value,
            DestOperatorID: checkList.value
        }
    )
    if (err) {
        return tip.e(err)
    }
    loadingUpdate.value = false
    await getStatus4GameUpdateMsg()
}

const getResult4GameUpdateMsg = async () => {
    loadingUpdate.value = true
    let [data, err] = await Client.Do(AdminInfo.getResult4GameUpdateMsg, {})
    loadingUpdate.value = false
    if (err) {
        return tip.e(err)
    }
    if (data.Completed === false){
        tip.w(t('消息正在处理中'))
    } else {
        FailOperatorID.value = []
        operatorData.value.forEach(i=>{
            data.FailOperatorID.forEach(u=>{
                if (i.Id === u) {
                    FailOperatorID.value.push(i.Name)
                }
            })
        })
    }
}

const getStatus4GameUpdateMsg = async () => {
    loadingUpdate.value = true
    let intervalId;
    let status = null
    const fetchStatus = async () => {
        let [data, err] = await Client.Do(AdminInfo.getStatus4GameUpdateMsg, {});
        if (err) {
            tip.e(err);
            clearInterval(intervalId); // 请求失败，停止定时任务
            return;
        }

        if (!data.Status) {
            clearInterval(intervalId); // 状态为0，停止定时任务
            loadingUpdate.value = false;
            status = data.Status;
            tip.s(t("消息推送成功，可以点击'查看结果'查询细节"));
        } else {
            tip.w(t("消息推送中"));
        }
    };
    // 首次执行 fetchStatus
    let [data, err] = await Client.Do(AdminInfo.getStatus4GameUpdateMsg, {});
    if (err) {
        tip.e(err);
        loadingUpdate.value = false;
        return;
    }

    if (!data.Status) {
        loadingUpdate.value = false;
        tip.s(t("消息推送成功，可以点击'查看结果'查询细节"));
    } else {
        // 如果首次 Status 不为 0，设置定时器
        intervalId = setInterval(fetchStatus, 10000); // 每10秒执行一次
    }
}

const operatorList = async () => {
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, operatorParam)
    if (err) {
        return tip.e(err)
    }
    operatorData.value = data.AllCount === 0 ? [] : data.List.filter(list => !list.Status)
}
</script>
<style scoped>
.group {
    display: flex;
    width: auto;
    flex-direction: column;
    align-items: center;
}

.flex {
    justify-content: space-between;
}

:deep(.el-step__title) {
    font-size: 14px;
    font-weight: normal;
}

.gap {

    display: grid;
    width: 50%;
    max-height: 200px;
    overflow-y: auto;


}

.gap-more{
    gap: 1rem;
    grid-template-rows: repeat(4, minmax(2px, 2fr));
    grid-auto-flow: column;
    justify-items: start;
}
:deep(.gap .el-checkbox) {
    min-width: 125px;
}
</style>
