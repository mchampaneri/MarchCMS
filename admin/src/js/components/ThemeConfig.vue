<template>
  <div class="card">
                        <div class="card-header">
                            <p class="card-header-title">Theme Configuration</p>
                        </div>
                        <div class="card-content">
                            <div class="card">
                                <div class="card-header">
                                    <p class="card-header-title">Menus</p>
                                </div>
                                <div class="card-content">
                                    <div class="field" v-for="(menu,index) in Menus" v-bind:key="index">
                                        <div class="control">
                                            <label for="" class="label">{{ menu.Place }}</label>
                                            <input type="text" class="input" v-model="menu.Menu" :name="menu.Place">
                                        </div>
                                    </div>

                                    <div class="field">
                                            <button type="button" class="button is-primary is-pulled-right"
                                             v-bind:class="[isSaving ? 'is-loading':'']"
                                              @click="SaveSettings()"
                                            >
                                            <i class="fa fa-save"></i>&nbsp Save
                                            </button>
                                    </div>
                                    <div style="clear: both"></div>
                                </div>


                            </div>

                        </div>
                    </div>

</template>

<script>
export default {

    mounted(){
        var vm = this
          axios.get('/admin/theme/settings')
          .then((result) => {
              console.log(result.data)
              vm.Menus = result.data.config.Menus
            }).catch((err) => {
                console.log(err)
            });
    },

    data(){
        return {
            Menus:[],
            isSaving:false,
        }
    },

    methods:{
        SaveSettings(){
            let vm = this;
            vm.isSaving = true;
            axios.post('/admin/theme/settings',{
                'menus':vm.Menus
            }).then((result) => {
                vm.Menus  = result.data.menus
                vm.isSaving = false;
            }).catch((err) => {
                console.log(err)
                vm.isSaving = false;
            });
        }
    }

}
</script>

