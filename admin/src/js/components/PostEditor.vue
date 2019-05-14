<template>
    <div>
        <div class="field">
            <div class="control">
                 <label class="label">Post Title</label>
                <input class="input is-medium" type="text"
                v-model="PageTitle"
                placeholder="Page Title">
            </div>
        </div>
        <div class="field">
             <div class="control">
                  <label class="label">Post URL</label>
                <input class="input is-medium" type="text"
                v-model="PageURL"
                placeholder="Page URL">
            </div>
        </div>
        <div class="field">
             <div class="control">
                <label class="label">Post Description</label>
                <input class="input is-medium" type="text"
                v-model="Desc"
                placeholder="Page Description">
            </div>
        </div>
        <div class="field">
             <div class="control">
                <label class="label">Post keywords </label>
                <input class="input is-medium" type="text"
                v-model="Keywords"
                placeholder="Page Keywords">
            </div>
        </div>
        <div class="field">
             <div class="control">
                <label class="label">Post Template</label>
                <select  v-model="PageTemplate">
                    <option v-bind:value="template" v-bind:key="i" v-for="(template,i) in PageTemplates">
                            {{ template }}
                    </option>
                </select>
            </div>
        </div>

       <div class="field">
                <label class="label">Content ( in markdown )</label>
                <div class="control">
                <textarea class="textarea"
                v-model="HTML"
                placeholder="Textarea"></textarea>
            </div>
        </div>

        <div class="control">
            <button class="button is-link" @click="SavePage()">Save</button>
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
                axios.post('/admin/post/'+vm.PageNumber+'/edit', {
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
                 axios.post('/admin/post/create', {
                     'PageTitle':vm.PageTitle,
                     'PageURL':vm.PageURL,
                     'Desc':vm.Desc,
                     'Keywords':vm.Keywords,
                     'HTML':vm.HTML,
                       'PageTemplate':vm.PageTemplate,
                 })
                 .then(function (response) {
                     console.log(response);
                     window.location = "/admin/post/"+response.data.PageNumber+"/edit"
                 })
                 .catch(function (error) {
                     console.log(error);
                 });
            }
        }
    },

}
</script>