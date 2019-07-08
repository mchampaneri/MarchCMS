<template>
    <div style="max-height:150px;min-height:auto;">
        <div v-if="showDropZone">
          <button @click="showDropZoneTrigger" class="button is-danger  is-pulled-right">
            <i class="fa fas fa-times"></i> Cancle
          </button>
          <vue-dropzone
          ref="myVueDropzone"
          id="dropzone"
          :options="dropzoneOptions"
          v-on:vdropzone-success="success">
          </vue-dropzone>
        </div>
        <div v-if="!showDropZone">
          <button @click="showDropZoneTrigger" class="button is-primary  is-pulled-right"> <i class="fa fas fa-plus"></i> Upload New Item</button>
        </div>
    </div>
</template>

<script>

import vue2Dropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'

export default {
  components: {
    vueDropzone: vue2Dropzone
  },
  data() {
    return {
      showDropZone: false,
      dropzoneOptions: {
          url: '/asset/upload/file',
          thumbnailWidth: 150,
          maxFilesize: 0.5,
          headers: { "My-Awesome-Header": "header value" }
      }
    }
  },
  methods:{
      success(){
          window.location.reload();
      },
      showDropZoneTrigger(){
          let vm = this;
          vm.showDropZone = (vm.showDropZone  == true)?false : true;
      }
  }
}
</script>
