<template>
  <aside id="aside" class="w-20 text-base fixed top-0 z-40 h-screen transition-position lg:left-0 overflow-y-auto dark:border-r dark:border-gray-800 dark:bg-gray-900/70 lg:dark:bg-gray-900 xl:dark:bg-gray-900/70 bg-white border-r border-gray-100 aside-scrollbars-light -left-60 lg:hidden xl:block">
    <div class="flex flex-row w-full flex-1 h-14 items-center dark:bg-transparent">
      <div class="hidden lg:flex xl:hidden items-center grow-0 shrink-0 relative cursor-pointer py-2 px-3">
        <span class="inline-flex justify-center items-center w-6 h-6 cursor-pointer">
          <svg viewBox="0 0 24 24" width="24" height="24" class="inline-block"><path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z"></path></svg>
        </span>
      </div><div class="flex-1 px-3"><span>软件工程训练营</span> </div>
    </div>
    <div>
      <ul>
        <li>
          <a class="router-link-active router-link-exact-active flex cursor-pointer dark:hover:bg-gray-700/50 text-blue-600 hover:text-black dark:text-white py-2" aria-current="page" @click='onToHome'>
          <span class="inline-flex justify-center items-center w-4 h-6 flex-none font-bold text-black dark:text-white">
            <svg viewBox="0 0 24 24" width="16" height="16" class="inline-block"><path fill="currentColor" d="M21,14H3V4H21M21,2H3C1.89,2 1,2.89 1,4V16A2,2 0 0,0 3,18H10L8,21V22H16V21L14,18H21A2,2 0 0,0 23,16V4C23,2.89 22.1,2 21,2Z"></path></svg>
          </span>
          <span class="inline-flex justify-center items-center w-10 h-6 flex-none font-bold text-black dark:text-white">所有文件</span>
          </a><!---->
        </li>
      </ul>

    </div>
  </aside>

</template>

<script setup>
import router from '@/router'

const onToHome = () => {
// 配置跳转方式
router.push('/')
}
</script>