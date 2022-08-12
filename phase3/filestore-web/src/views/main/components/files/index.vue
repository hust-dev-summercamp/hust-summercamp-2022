<template>
  <div>
    <div class="md:rounded text-sm p-5 dark:bg-gray-700 bg-white border border-gray-100 dark:border-gray-800 mb-6 dark:text-white"><!---->
      <div class="">
        <table>
          <thead>
          <tr><th></th><th></th><th>文件hash</th><th>文件名</th><th>文件大小</th><th>上传时间</th><th>最近更新</th><th>操作</th><th></th></tr>
          </thead>
          <tbody>
          <tr v-for='(item, index) in $store.getters.files'>
            <td></td>
            <td class='hidden'>{{ index }}</td>
            <td class="checkbox-cell">
              <label class="checkbox"><input type="checkbox"><span class="check"></span></label>
            </td>
            <td class="text-center lg:w-32">{{ item.file_hash }}</td>
            <td class="text-center lg:w-32">{{ item.file_name }}</td>
            <td class="text-center lg:w-20 whitespace-nowrap">{{ item.file_size }}</td>
            <td class="lg:w-32">{{ item.upload_at }}</td>
            <td class="lg:w-32">{{ item.last_updated }}</td>
            <td class="before:hidden lg:w-1 whitespace-nowrap">
              <div class="flex items-center justify-start lg:justify-end flex-nowrap -mb-3">
                <m-button
                  class="inline-flex justify-center bg-blue-500 text-white border-blue-600 hover:bg-blue-600 m-1"
                  icon="download"
                  iconColor="#fff"
                  @click="onDownloadClick(item.file_hash, item.file_name)"
                ></m-button>
                <m-button
                  class="inline-flex justify-center bg-green-600 text-white border-green-700 hover:bg-green-700 m-1"
                  icon="setting"
                  iconColor="#fff"
                  @click="onShowEditClick(item.file_hash, item.file_name)"
                ></m-button>

                <m-button
                  class="inline-flex justify-center bg-blue-500 text-white border-blue-500 hover:bg-blue-500 m-1"
                  icon="unfold"
                  iconColor="#fff"
                  @click="onShowCdlClick(item.file_hash, item.file_name)"
                ></m-button>
                <m-button
                  class="inline-flex justify-center bg-red-600 text-white border-red-700 hover:bg-red-700 m-1"
                  icon="delete"
                  iconColor="#fff"
                  @click="onDeleteClick(item.file_hash)"
                ></m-button>
              </div>
            </td>
          </tr>

          </tbody>
        </table>
      </div>
    </div>
    <m-popup v-model="isOpenEditPopup">
      <div class="h-[30px] p-5">
        <m-input
          v-model="editingFileName"
          class="w-full"
          type="text"
          max="100"
        ></m-input>
      </div>
      <div class='flex flex-grow justify-end p-5'>
        <m-button
          class="inline-flex justify-center bg-green-600 text-white border-green-700 hover:bg-green-700 m-1"
          iconColor="#fff"
          @click="onFileUpdateClick()"
        >确定</m-button>
      </div>
    </m-popup>
    <m-popup v-model="isOpenCdlPopup">
      <div id="progressing">[下载进度] {{cdlProgress}}%</div>
      <m-button
        class="inline-flex justify-center bg-red-600 text-white border-red-700 hover:bg-red-700 m-1"
        iconColor="#fff"
        @click="abortDownload"
      >暂停/继续</m-button>
    </m-popup>
  </div>
</template>

<script setup>
import { useStore } from 'vuex'
import { saveAs } from 'file-saver'
import { ref, watch } from 'vue'
import { useVModel } from '@vueuse/core'

const store = useStore()
store.dispatch('file/useFilesData', {limit: 15})
// popup 展示
const isOpenEditPopup = ref(false)
const isOpenCdlPopup = ref(false)

const onDownloadClick = (filehash, filename) => {
  console.log(filehash)
  store.dispatch('file/downloadUrl', {filehash}).then((resp)=>{
    console.log('dispatch then')
    console.log(resp)
    setTimeout(() => {
      /**
       * 接收两个参数：
       * 1. 下载的图片链接
       * 2. 下载的文件名称
       */
      saveAs(
        resp,
        filename
      )
    }, 100)

  })
}

const editingFileName = ref('')
const editingFilehash = ref('')
const onShowEditClick = (filehash, filename) => {
  editingFileName.value = filename
  editingFilehash.value = filehash
  isOpenEditPopup.value = true
}

const onFileUpdateClick = () => {
  store.dispatch('file/update', {
    username: store.getters.username,
    filehash: editingFilehash.value,
    op: '0',
    filename: editingFileName.value,
  }).then((response) => {
    isOpenEditPopup.value = false
  })
}

let req = new XMLHttpRequest();
let fileParts = [];
let filesize = -1;
let paused = false;
let lastByteLoaded = 0;
let byteLoaded = 0;
let url = localStorage.getItem('curDownloadUrl');
let filename = "unknown";
const cdlProgress = ref(0)
const onShowCdlClick = (filehash, name) => {
  isOpenCdlPopup.value = true
  filename = name
  url = `http://localhost:8080/api/file/download/range?filehash=${filehash}&username=${store.getters.username}&token=${store.getters.token}`
  req.open("GET", url, true);
  addReqListeners(req);
  req.setRequestHeader('Range', 'bytes=' + lastByteLoaded + '-');
  req.send();
}

function loadedSize() {
  var size = 0;
  for (let i = 0; i < fileParts.length; i++) {
    size += fileParts[i].byteLength;
  }
  return size;
}

const addReqListeners  = (req) => {
  req.onprogress = function (evt) {
    if (evt.lengthComputable) {
      // 第一次获取文件总大小
      if (filesize < 0) {
        filesize = evt.total;
        req.abort();
        // 真正开始下载文件
        resumeDownload();
      } else {
        var percentComplete = (lastByteLoaded + evt.loaded) / filesize;
        // $("#progressing").html('[下载进度] ' + (percentComplete * 100) + "%");
        cdlProgress.value = percentComplete*100
        byteLoaded = lastByteLoaded + evt.loaded;
        // if (byteLoaded >= filesize) {
        //   $("#ctlbtn").attr("disabled", "true");
        // }
      }
    }
  }
  req.responseType = "arraybuffer";
  req.onreadystatechange = function () {
    if (req.status !== 200 && req.status !== 206) {
      return;
    }
    var name = req.getResponseHeader("Content-disposition");
    filename = name.match(/attachment; filename="(.*)"/)[1];

    if (req.readyState === 4) {}
  }
  req.onload = function (oEvent) {
    var arrayBuffer = req.response; // Note: not oReq.responseText
    if (arrayBuffer) {
      fileParts.push(new Uint8Array(arrayBuffer));
      lastByteLoaded = loadedSize();
      // alert(filesize + '  ' + lastByteLoaded + '   ' + arrayBuffer.byteLength);
      if (lastByteLoaded >= filesize) {
        cdlProgress.value = 100
        // $("#progressing").html('[下载进度] 100%');
        // $("#ctlbtn").attr("disabled", "true");

        var blob = new Blob(fileParts);
        saveAs(blob, filename);
      } else {
        resumeDownload();
      }
    }
  }
}

const abortDownload = () => {
  if (!paused) {
    // 暂停下载，结束当前partial download
    req.abort();
  } else {
    // 继续下载，开始下一个指定起始位置的partial download
    resumeDownload();
  }
  paused = !paused;
}

const resumeDownload = () => {
  req.open("GET", url, true);
  addReqListeners(req);
  lastByteLoaded = loadedSize();
  req.setRequestHeader('Range', 'bytes=' + lastByteLoaded + '-' + (lastByteLoaded + 1024 * 1024 * 20 - 1));
  req.send();
}

const onDeleteClick = (filehash, filename) => {
  store.dispatch('file/delete', {filehash, username: store.getters.username}).then(() =>{
    console.log('delete successfully')
  })
}
</script>