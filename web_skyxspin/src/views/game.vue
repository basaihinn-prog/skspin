<template>
    <div class="frameContext">
        <iframe id="frame" class="frameId" :src="gameUri" frameborder="0"></iframe>
    </div>
</template>


<script setup lang="ts">


import {useRoute} from "vue-router";

const route = useRoute()

const gameUri = route.params.GameId

const frameId = document.getElementById("frame")
frameId.onload = () => {

    console.log("--------------------------frameId['contentWindow']----------------------------------------")
    console.log(frameId['contentWindow'])
    console.log("-----------------------------frameId['contentWindow'].getElementById('splash').offsetWidth---------------------------------")
    console.log(frameId['contentWindow'].getElementById("splash").offsetWidth)
    frameId['contentWindow'].getElementById("splash").offsetWidth

}


</script>


<style scoped>
.frameId{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    padding: 0;
    margin: 0;
}

.frameContext{
    width: 100%;
    height: 100vh;
    overflow: hidden;
}
</style>
