
<template>
    <router-view/>
</template>
<script setup lang="ts">

import {GetSystemAllData, initStore} from "@/util/generatorData";
import {onBeforeMount} from "vue";


onBeforeMount(()=>{

    initStore()
    GetSystemAllData()
})



</script>
