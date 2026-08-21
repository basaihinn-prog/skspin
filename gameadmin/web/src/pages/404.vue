/*
 *@Author: 西西米
 *@Date: 2023-02-25 13:58:34
 *@Description: 404
*/
<template>
    <div class='view-404'>
        <el-empty :description="$t('～找不到当前页面了!～')">
            <el-button @click="goHome()" type="primary">{{$t('返回首页')}}</el-button>
        </el-empty>
    </div>
</template>

<script setup lang='ts'>
import { onMounted } from 'vue';
import { useStore } from '@/pinia/index';
import router from '@/router';
const store = useStore()
const piniaStore = useStore()
const {setToken, setAdminInfo,setMenuList} = piniaStore
const goHome = () => {
    setToken('')
    setMenuList('')
    setAdminInfo({})
    localStorage.clear();
    router.push({ path: '/dashboard' })
}
onMounted(() => {

});
</script>

<style lang='scss'>
.view-404 {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
