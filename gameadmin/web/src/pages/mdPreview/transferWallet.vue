<template>
    <div style="width: 100%;position: relative" class="flex">

        <MdPreview :editorId="id" :modelValue="inputmarkdown" class="md-editor"/>
        <MdCatalog :editorId="id" :scrollElement="scrollElement" />
    </div>
</template>

<script setup lang='ts'>
import { onMounted, ref, reactive } from 'vue';
import { Client } from '@/lib/client';
import { tip } from '@/lib/tip';
import { MdPreview, MdCatalog } from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
const inputmarkdown = ref('')
const scrollElement = document.documentElement;
const init = async () => {
    let url = 'https://dl.kafa010.com/api-transfer.md';
    const [resp, err] = await Client.send(url,{});
    if (err) {
        tip.e(err);
        return;
    }
}
const id = 'preview-only';
onMounted(async () => {
    fetch('https://dl.kafa010.com/api-transfer.md')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.text();
        })
        .then(data => {
            inputmarkdown.value = data
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
});
</script>

<style scoped>
.md-editor-catalog{
    position: -webkit-sticky;
    position: sticky;
    top: 1rem;
    right: 0;
    background: rgb(143 143 143 / 11%);
    padding: 20px;
    backdrop-filter: blur(8px);
    border-radius: 1rem;
    max-height: 800px;
    overflow-y: scroll;
    margin-top: 20px;
}
.md-editor{
    width: calc(100% - 330px);
}
</style>
