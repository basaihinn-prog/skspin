

<template>
    <div>
      <OperatorMaintenance/>
<!--      <Merchant v-else/>-->
    </div>
</template>

<script setup lang="ts">
  import {useStore} from "@/pinia";

  let { AdminInfo } = useStore()

  import OperatorMaintenance from "./operatorMaintenance.vue"
  import Merchant from "./merchants.vue"

</script>

<style scoped lang="scss">

</style>
