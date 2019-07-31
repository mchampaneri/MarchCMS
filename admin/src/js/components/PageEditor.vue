<template>
    <div class="container box">

        <div>
            <div class="field is-grouped is-pulled-right" >
                <div class="control">
                    <button class="button is-primary is-pulled-right" v-bind:class="[isSaving ? 'is-loading':'']" @click="SavePage()"> <i class="fa fas fa-save"></i> &nbsp Save</button>
                </div>
            </div>
            <div style="clear:both"></div>
        </div>

     <div class="tabs is-boxed">
  <ul style="border-bottom-color: #dbdbdb;
    border-bottom-style: solid;
    border-bottom-width: 1px;">
    <li @click="function(){fillMeta = true}" v-bind:class="[fillMeta ? 'is-active'  : '']">
      <a>
        <span class="icon is-small"><i class="fas fa-image" aria-hidden="true"></i></span>
        <span>Meta</span>
      </a>
    </li>
    <li  @click="function(){fillMeta = false}" v-bind:class="[!fillMeta ? 'is-active' : '']">
      <a>
        <span class="icon is-small"><i class="fas fa-music" aria-hidden="true"></i></span>
        <span>Content</span>
      </a>
    </li>
  </ul>
</div>

        <div  v-if="fillMeta">

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
                <div class="control">
                    <textarea class="textarea"
                    v-model="HTML"
                    placeholder="Textarea"></textarea>
                </div>
        </div>

        <br>
        <div>
            <div class="field is-grouped is-pulled-right" >
                <div class="control">
                    <button class="button is-primary is-pulled-right" v-bind:class="[isSaving ? 'is-loading':'']" @click="SavePage()"> <i class="fa fas fa-save"></i> &nbsp Save</button>
                </div>
            </div>
            <div style="clear:both"></div>
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
            isSaving:false,
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
            vm.isSaving = true;
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
                    vm.isSaving = false;
                })
                .catch(function (error) {
                    console.log(error);
                    vm.isSaving = false;
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
                     vm.isSaving = false;
                 })
                 .catch(function (error) {
                     console.log(error);
                     vm.isSaving = false;
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