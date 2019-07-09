<template>
    <div>

        <div>
            <div class="field is-grouped is-pulled-right" >
                <div class="control">
                    <div class="select">
                        <select v-model="stage" @click="checkStatus" >
                            <option value="meta">Page Meta</option>
                            <option value="content" :disabled="blockEditor">Page Content</option>
                        </select>
                    </div>
                </div>
                <div class="control">
                    <button class="button is-primary is-pulled-right" @click="SavePage()"> <i class="fa fas fa-save"></i> &nbsp Save</button>
                </div>
            </div>
            <div style="clear:both"></div>
        </div>

        <div v-if="blockEditor" class="notification is-warning">
          <button class="delete"></button>
          Please fill all page meta info before writing content
        </div>

        <div  v-if="fillMeta">
            <h4 class="title is-4"> Page Meta Information </h4>
            <div class="columns is-multiline">
                <div class="column is-6">
                <div class="field">
                    <div class="control">
                         <label class="label">Page Title</label>
                        <input class="input is-medium" type="text"
                        v-model="PageTitle"
                        placeholder="Page Title">
                    </div>
                </div>
                </div>
                <div class="column is-6">
                <div class="field">
                    <div class="control">
                        <label class="label">Page URL</label>
                        <input class="input is-medium" type="text"
                        v-model="PageURL"
                        placeholder="Page URL">
                    </div>
                </div>
                </div>
                <div class="column is-6">
                <div class="field">
                     <div class="control">
                        <label class="label">Page Description</label>
                        <input class="input is-medium" type="text"
                        v-model="Desc"
                        placeholder="Page Description">
                    </div>
                </div>
                </div>
                <div class="column is-6">
                <div class="field">
                     <div class="control">
                        <label class="label">Page keywords </label>
                        <input class="input is-medium" type="text"
                        v-model="Keywords"
                        placeholder="Page Keywords">
                    </div>
                </div>
                </div>
                <div class="column is-6">
                <label  class="label">Page Template</label>
                <div class="field">
                    <div class="control">
                        <div class="select">
                            <select  v-model="PageTemplate" >
                                <option v-bind:value="template" v-bind:key="i" v-for="(template,i) in PageTemplates">
                                        {{ template }}
                                </option>
                            </select>
                        </div>
                    </div>
                </div>
                </div>
            </div>
            </div>
            <div v-if="!fillMeta" class="field">
            <h4 class="title is-4"> Page Content ( Markdown ) </h4>
                <div class="control">
                    <textarea class="textarea"
                    v-model="HTML"
                    placeholder="Textarea"></textarea>
                </div>
        </div>
    </div>
</template>

<script>


export default{

    props:['isedit','opagetitle','opageurl','odesc','okeywords','ohtml','opagenumber', 'pagetemplates', 'pagetemplate'],


    mounted(){
        console.log("Editor component has been mounted");
        var vm = this
        vm.PageTemplates = vm.pagetemplates
        if(vm.isedit == "true"){
            vm.PageTitle = vm.opagetitle
            vm.PageURL = vm.opageurl
            vm.HTML = vm.ohtml
            vm.Desc = vm.odesc
            vm.Keywords = vm.okeywords
            vm.PageNumber = vm.opagenumber
            vm.PageTemplate = vm.pagetemplate
        }
    },


data(){
        return{
            blockEditor:true,
            fillMeta: true,
            stage:'meta',

            HTML:"",
            PageTitle:"",
            PageURL:"",
            Desc:"",
            Keywords:"",
            PageNumber:"-",
            PageTemplate:"-",
            PageTemplates:"-",
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
                    'PageTemplate':vm.PageTemplate,
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
                       'PageTemplate':vm.PageTemplate,
                 })
                 .then(function (response) {
                     console.log(response);
                     window.location = "/admin/page/"+response.data.PageNumber+"/edit"
                 })
                 .catch(function (error) {
                     console.log(error);
                 });
            }
        },

        checkStatus(){
            let vm =this;
            if (vm.PageTitle != "" &&
                vm.PageURL != "" &&
                vm.PageTemplate != ""){
                    vm.blockEditor = false;
                }
        }
    },

}
</script>