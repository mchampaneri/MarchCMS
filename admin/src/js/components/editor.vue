<template>
    <div>
        <div class="field">
            <div class="control">
                <input class="input is-medium" type="text"
                v-model="PageTitle"
                placeholder="Page Title">
            </div>
             <div class="control">
                <input class="input is-medium" type="text"
                v-model="PageURL"
                placeholder="Page URL">
            </div>
             <div class="control">
                <input class="input is-medium" type="text"
                v-model="Desc"
                placeholder="Page Description">
            </div>
             <div class="control">
                <input class="input is-medium" type="text"
                v-model="Keywords"
                placeholder="Page Keywords">
            </div>
        </div>
        <VueTrix v-model="HTML"
                placeholder="Enter content"
                trix-file-accept="alert('something is beign dragged')"
                localStorage/>

        <div class="control">
            <button class="btn btn-default" @click="SavePage()">Save</button>
        </div>
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
            HTML:"none",
            PageTitle:"",
            PageContent:"",
            PageURL:"",
            Desc:"",
            Keywords:"",
        }
    },

    methods:{
        SavePage:function(){
            var vm = this;
            axios.post('/admin/page/create', {
                'PageTitle':vm.PageTitle,
                'PageURL':vm.PageURL,
                'Desc':vm.Desc,
                'Keywords':vm.Keywords,
                'HTML':vm.HTML,
            })
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                console.log(error);
            });
        }
    },

}
</script>