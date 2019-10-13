<template>
    <div class="card">
        <div class="card-header">
            <p class="card-header-title">Your profile</p>
        </div>
        <div class="card-content">
            <div class="field">
                <div class="control">
                    <label for="" class="label">Name</label>
                    <input type="text" class="input" v-model="Name" />
                </div>
            </div>
            
            <div class="field">
                <div class="control">
                    <label for="" class="label">Picture URL</label>
                    <input type="text" class="input" v-model="Picture" />
                </div>
            </div>
            <div class="field">
                <div class="control">
                    <label for="" class="label">About Self</label>
                    <textarea type="text" class="textarea" v-model="SmallDesc" />
                </div>
            </div>

            <br>
            <div>
                <div class="field is-grouped is-pulled-right" >
                    <div class="control">
                        <button class="button is-primary is-pulled-right" v-bind:class="[isSaving ? 'is-loading':'']" @click="SaveBasicProfile()"> <i class="fa fas fa-save"></i> &nbsp Save</button>
                    </div>
                </div>
                <div style="clear:both"></div>
            </div>

        </div>
    </div>
</template>

<script>
export default {
    mounted(){
        let vm = this;
        axios.get("/admin/user/profile")
            .then((result) => {
                console.log(result)
                vm.Name = result.data.user.Name
                vm.Picture  = result.data.user.Picture
                vm.SmallDesc = result.data.user.SmallDesc
            }).catch((err) => {
                console.log(err)
            });
    },

    data(){
        return{
            "isSaving":false,
            "Name":'',
            "SmallDesc":'',
            "Picture":'',
        }
    },

    methods:{

            SaveBasicProfile(){
                let vm=this;
                vm.isSaving = true;
                // Making save request for menu
                axios.post("/admin/user/profile",{
                    'Name':vm.Name,
                    'SmallDesc':vm.SmallDesc,
                    'Picture':vm.Picture,
                })
                .then((result) => {
                    console.log(result)
                    vm.Name = result.data.user.Name
                    vm.Picture  = result.data.user.Picture
                    vm.SmallDesc = result.data.user.SmallDesc
                    vm.isSaving = false;
                }).catch((err) => {
                    console.log(err)
                    vm.isSaving = false;
                });
            }
    }
}
</script>

