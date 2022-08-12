import {
  getFiles,
  downloadUrl,
  uploadFile,
  deleteFile,
  updateFile
} from '@/api/files'

export default {
  namespaced: true,
  state: () => ({
    // file list 展示的数据源
    files: []
  }),
  mutations: {
    /**
     * 为 categorys 赋值
     */
    setFiles(state, files) {
      console.log(files)
      if (files) {
        state.files = [...files]
      } else {
        state.files = files
      }
      console.log(state.files)
    },
    deleteFile(state, filehash) {
      let deleteIndex
      state.files.forEach((item, index) => {
        if (item.file_hash === filehash) {
          deleteIndex = index
        }
      })

      state.files.splice(deleteIndex, 1)
    },
    updateFile(state, fileParams) {
      let updateIndex
      state.files.forEach((item, index) => {
        if (item.file_hash === fileParams.file_hash) {
          updateIndex = index
        }
      })
      state.files[updateIndex].file_name = fileParams.file_name
      console.log(fileParams)
    }
  },
  actions: {
    /**
     * 获取 category 数据，并自动保存到 vuex 中
     */
    async useFilesData(context, params) {
      const files = await getFiles(params)
      context.commit('setFiles', files)
    },

    downloadUrl(context, params) {
      return new Promise(resolve => {
        downloadUrl(params, (response) => {
          resolve(response)
        })
      })
    },

    upload(context, params) {
      return new Promise(resolve => {
        uploadFile(params, (response) => {
          console.log("store then")
          console.log(response)
          resolve(response)
        })
      })
    },

    delete(context, params) {
      return new Promise(resolve => {
        deleteFile(params, (response) => {
          context.commit('deleteFile', params.filehash)
          resolve(response)
        })
      })
    },

    update(context, params) {
      return new Promise(resolve => {
        updateFile(params, (response) => {
          context.commit('updateFile', response)
          resolve(response)
        })
      })
    }
  }
}