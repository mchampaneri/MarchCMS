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

    props:['isedit','opagetitle','opageurl','odesc','okeywords','ohtml','opagenumber'],

 components: {
    VueTrix
  },

    mounted(){
        console.log("Editor component has been mounted");
        var vm = this
        if(vm.isedit == "true"){
            vm.PageTitle = vm.opagetitle
            vm.PageURL = vm.opageurl
            vm.HTML = vm.ohtml
            vm.Desc = vm.odesc
            vm.Keywords = vm.okeywords
            vm.PageNumber = vm.opagenumber
        }
    },


data(){
        return{
            HTML:"",
            PageTitle:"",
            PageURL:"",
            Desc:"",
            Keywords:"",
            PageNumber:"-",
        }
    },

    methods:{
        SavePage:function(){
            var vm = this;

            if (vm.isedit == 'true'){
                axios.post('/admin/page/'+vm.PageNumber+'/edit', {
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
            }else{
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
        }
    },

}
</script>