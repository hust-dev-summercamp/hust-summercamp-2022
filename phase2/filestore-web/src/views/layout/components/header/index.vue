<template>
  <div
    class="w-full duration-500 bg-white dark:bg-zinc-800 border-b border-b-zinc-200 dark:border-b-zinc-700 px-2 py-1"
  >
    <div class="flex items-center">
      <m-button
        class="inline-flex justify-center bg-red-600 text-white border-red-700 hover:bg-red-700 m-1 w-5"
        icon="fold"
        iconColor="#fff"
        @click="onShowUploadPopup"
      ></m-button>
      <!-- 隐藏域 -->
      <input
        v-show="false"
        ref="inputFileTarget"
        type="file"
        @change="onSelectFileHandler"
      />
      <header-account class="mr-1"></header-account>
      <header-theme-vue class="mr-1"></header-theme-vue>
      <header-my-vue></header-my-vue>
    </div>
  </div>
</template>

<script setup>
// import headerSearchVue from './header-search/index.vue'
import headerAccount from './header-account.vue'
import headerThemeVue from './header-theme.vue'
import headerMyVue from './header-my.vue'
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import { useStore } from 'vuex'


const router = useRouter()
const store = useStore()
const onToHome = () => {
  router.push('/')
}

// 隐藏域
const inputFileTarget = ref(null)

const onShowUploadPopup = () => {
  console.log('show upload popup')
  inputFileTarget.value.click()
}

const onSelectFileHandler = () => {
  const file = inputFileTarget.value.files[0]
  const bodyFormData = new FormData()
  bodyFormData.append('file', file)
  store.dispatch('file/upload', bodyFormData).then((resp) => {
    store.dispatch("file/useFilesData", {limit: 15})
  })
  inputFileTarget.value.value = null
}
</script>
