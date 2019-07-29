<template>
 <div>
                    <div class="card">
                        <header class="card-header">
                            <p class="card-header-title">
                                Website details
                            </p>
                        </header>
                        <div class="card-content">
                            <div class="field">
                                <div class="control">
                                    <label for="" class="label">Title</label>
                                    <input type="text" class="input" v-model="Name" >
                                </div>
                            </div>
                            <div class="field">
                                <div class="control">
                                    <label for="" class="label">Logo URL</label>
                                    <input type="text" class="input"  v-model="LogoURL"  >
                                </div>
                            </div>
                            <div class="field">
                                <div class="control">
                                    <label for="" class="label">Favicon URL</label>
                                    <input type="text" class="input"  v-model="FaviconURL" >
                                </div>
                            </div>
                            <div class="field is-pulled-right">
                                <div class="control">
                                    <button class="button is-primary" @click="SaveConfig()">
                                        <i class="fa fa-save"></i> &nbsp Save
                                    </button>
                                </div>
                            </div>
                            <div style="clear:both"></div>
                        </div>
                    </div>
                    <!-- </div> -->
                </div>
</template>

<script>

export default {
    mounted(){
        let vm = this;
        console.log("site-settings component mounted")
        axios.get('/admin/site/settings')
            .then(function(response){
                console.log(response)
                vm.Name = response.data.config.Name;
                vm.FaviconURL = response.data.config.FaviconURL;
                vm.LogoURL = response.data.config.LogoURL;
            }).catch(function(err){
                console.log(err)
            })
    },
    data(){
        return{
            'Name':'',
            'FaviconURL':'',
            'LogoURL':'',
        }
    },
    methods:{
        SaveConfig(){
            let vm = this;
        axios.post('/admin/site/settings',{
            "Name": vm.Name,
            "LogoURL":vm.LogoURL,
            "FaviconURL": vm.FaviconURL,
            }).then(function(response){
                      console.log(response)

                  }).catch(function(err){
                      console.log(err)
                  })
              }
    }
}
</script>


