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
    props:["menus"],

    mounted(){
        var vm = this
        vm.Menus = JSON.parse(vm.menus)
        console.log(vm.Menus)
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
            axios.post('/admin/theme-settings/set-menu',{
                'menus':vm.Menus
            }).then((result) => {
                console.log(result)
            }).catch((err) => {
                console.log(err)
            });
        }
    }

}
</script>

