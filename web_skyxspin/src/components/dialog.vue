<template>




  <div class="dialog--mask animate" v-show="visible">
      <div class="dialog--mask dialog--mask--hh" v-show="visible" @click="closeDilalog">

      </div>
      <div class="dialog--container animate" v-show="visible" :style="{width: props.width, maxWidth: props.maxWidth}">

          <div class="dialog--header">{{ props.title }}</div>


          <div class="dialog--content">
              <slot></slot>
          </div>


      </div>

  </div>

</template>

<script setup lang="ts">

const visible = defineModel()
const props = defineProps(["title", "width", "maxWidth"])
const emits = defineEmits(["maskClose"])

const closeDilalog = () => {
    visible.value = false
    emits("maskClose")

}
</script>

<style scoped>
  .dialog--mask--hh{
      background-color: rgba(0,0,0,0.85)!important;
  }
  .dialog--mask{
      width: 100%;
      height: 100vh;
      position: fixed;
      top: 0;
      left: 0;
      background-color: rgba(0,0,0,0);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 9999999999999999;
  }


  .dialog--container{
      border: 1px solid rgba(130, 136, 151, 0.6);
      border-radius: 50px;
      box-shadow: rgba(238, 119, 0, 0.3) 0px 0px 50px 0px;
      background-color: rgb(20, 22, 32);
      width: 60%;

      min-height: auto;
      position: relative;
      z-index: 9999999999999;

  }


  .dialog--header{
      width: 100%;
      font-size: 30px;
      padding-top: 8%;
      margin-bottom: 46px;
      text-align: center;
      position: relative;
      color: rgb(255, 255, 255);


  }
  .animate{
      --animate-duration: 0.2s;
  }
  .dialog--content{

      padding: 0 20px;
      margin-bottom: 40px;
  }
</style>
