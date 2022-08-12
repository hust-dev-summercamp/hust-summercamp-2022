<template>
  <div class="w-full">
    <div class="flex-1 px-3 text-base">
      <div v-if="userName">
        <span>用户名：</span>
        <span>{{userName}}</span>
      </div>
      <div v-else>
        <span>未登录</span>
      </div>
    </div>
  </div>
</template>


<script setup>
import { useStore } from 'vuex'
import { computed } from 'vue'

const store = useStore()

const userName = computed(() => {
  if(store.getters.username) {
    return store.getters.username
  }
  return
})
</script>