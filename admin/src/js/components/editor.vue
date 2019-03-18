<template>
    <div>
        <div class="field">
            <div class="control">
                <input class="input is-medium" type="text" placeholder="Medium input">
            </div>
        </div>
        <VueTrix v-model="editorContent"
                placeholder="Enter content"
                trix-file-accept="alert('something is beign dragged')"
                localStorage/>
    </div>
</template>

<script>

import VueTrix from "vue-trix";

export default{


    mounted(){
        console.log("Editor component has been mounted");
    },

    data(){
        return{
            editorContent:"none"
        }
    }

}
</script>