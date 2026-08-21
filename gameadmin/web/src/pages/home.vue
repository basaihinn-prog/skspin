/*
 *@Author: 西西米
 *@Date: 2023-02-25 09:59:41
 *@Description: 首页
*/
<template>
    <div class='home'>
        {{ $t('游戏管理后台') }}
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
onMounted(() => {

});
</script>
<style scoped lang='scss'>
.home{
    text-align: center;
    margin: 0 auto;
    font-size: 12rem;
}
</style>
