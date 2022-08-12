import { createStore } from 'vuex'
// import category from './modules/category'
import file from './modules/file'
import theme from './modules/theme'
import app from './modules/app'
// import search from './modules/search'
import user from './modules/user'
import getters from './getters'
import createPersistedState from 'vuex-persistedstate'

const store = createStore({
  getters,
  modules: {
    file,
    // category,
    theme,
    app,
    // search,
    user
  },
  plugins: [
    createPersistedState({
      // 保存到 localStorage 中的 key
      key: 'filestore-key',
      // 需要保存的模块
      // paths: ['file', 'category', 'theme', 'search', 'user']
      paths: ['file', 'theme', 'user']
    })
  ]
})

export default store
