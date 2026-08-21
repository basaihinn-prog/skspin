
<template>
    <router-view/>
</template>
<script setup lang="ts">

import {onBeforeMount} from "vue";
import { initStore} from "@/util/generatorData";

onBeforeMount(()=>{
    initStore()
})
</script>
