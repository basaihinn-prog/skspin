<template>

   <span style="font-size: 18px;color: #EB334D; font-weight: bolder">RTP</span>



  <!--  系统语言  -->
    <Dialog v-model="dialogVisible" :title="systemText.RTPSetting && systemText.RTPSetting['title']" width="50%" max-width="450px" @mask-close="maskClose">
        <div class="dialog-list">

            <el-form ref="form" :model="RTPForm" label-width="70px" label-position="top" style="width: 100%;display: flex;justify-content: center">
                <el-form-item style="width: 100%; max-width: 300px">

                    <template #label>
                        <div class="setRtp__label">{{systemText.RTPSetting && systemText.RTPSetting['label']}}:</div>
                    </template>

                    <el-select style="width: 100%; max-width: 300px" v-model="RTPForm.RTP" popper-class="RTPSelect">
                        <template v-for="item in RTPList">
                            <el-option :label="item + '%' + `${Number(item) == 96 ? `(${systemText.RTPSetting && systemText.RTPSetting['comm']})` : ''}`" :value="Number(item)" v-if="appId == 'faketrans' || (appId != 'faketrans' && item <= 120)"/>
                        </template>
                    </el-select>
                </el-form-item>
            </el-form>
        </div>

        <div class="dialog--foot">
            <el-button type="default" @click="closeDialog">{{systemText.RTPSetting && systemText.RTPSetting['cancel']}}</el-button>
            <el-button type="danger" @click="saveRTP">{{systemText.RTPSetting && systemText.RTPSetting['sure']}}</el-button>
        </div>


    </Dialog>


</template>


<script setup lang="ts">

import {computed, ref} from "vue";
import Dialog from "@/components/dialog.vue";
import {useLanguageStore} from "@/stores/store/langageStore";
import {storeToRefs} from "pinia";



const languageStore = useLanguageStore()
const { systemText } = storeToRefs(languageStore)
const RTPList = import.meta.env.VITE_RTP.split(",")
const appId = localStorage.getItem("app_id")
const randomStr = localStorage.getItem("randomStr")


const props = defineProps(['RTPDialogVisible'])
const emits = defineEmits(['closeRTPDialogVisible'])

const RTPForm = ref({
    RTP: 96
})

const dialogVisible = computed(()=>{


    if (props.RTPDialogVisible){
        document.body.style.overflow = "hidden"
        const RTPValue = localStorage.getItem("RTPValue")
        RTPForm.value.RTP = Number(RTPValue) || 96
    }else{

        document.body.style.overflow = "auto"
    }
    return props.RTPDialogVisible
})



const saveRTP = () => {
    localStorage.setItem("RTPValue", RTPForm.value.RTP)
    closeDialog()
}
const closeDialog = () => {

    emits("closeRTPDialogVisible")
}
const maskClose = () => {
    setTimeout(()=>{

        closeDialog()
    },50)
}

</script>

<style scoped>
.dialog-list{
    width: 80%;
    margin: -30px auto 0 auto;
}

.setRtp__label{
    font-size: 16px;
    font-weight: bolder;
    color: #ffffff;
}

.dialog--foot{
    width: 95%;
    display: flex;
    justify-content: flex-end;
    margin: 60px auto 0 auto;
}
</style>
