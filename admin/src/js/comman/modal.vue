<template>
      <div class="modal is-active">
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title" >{{modal_title}}</p>
            <button class="delete" aria-label="close"></button>
          </header>
          <section class="modal-card-body">
            <slot></slot>
          </section>
          <footer class="modal-card-foot">
            <button class="button is-success" @click="callOkFunc">Save changes</button>
            <button class="button">Cancel</button>
          </footer>
        </div>
    </div>
</template>

<script>
export default {
    props:['modal_title'],
    methods:{
        callOkFunc:function(e){
            this.$emit('clicked');
        }
    }
}
</script>
