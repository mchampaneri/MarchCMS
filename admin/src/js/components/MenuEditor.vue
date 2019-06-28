<template>
    <div>
        <modal v-if="deleteItem" modal_title="Are you sure?" @clicked="clicked(deleteItemId)">
            <p>Delete item from menu list.</p>
            <p>Action is irreversible.</p>
        </modal>

        <div class="columns">
            <div class="column is-one-third">
            <div class="field">
              <div class="control">
                <input class="input" type="text"
                v-model="Name" placeholder="Menu Name">
              </div>
            </div>
            <div class="field">
                <label for="">Add New Menu Item</label>
                    <div class="field">
                    <div class="control">
                        <input class="input" type="text"
                        v-model="NewItem.Name" placeholder="Menu Title">
                    </div>
                    </div>
                    <div class="field">
                       <div class="control">
                           <input class="input" type="text"
                           v-model="NewItem.URL" placeholder="Menu URL">
                       </div>
                    </div>
                    <button type="button" class="button is-normal"
                        @click="AddItem(NewItem)">
                        Add Item
                    </button>
                </div>
            </div>
            <div class="column">
          <div class="panel">
            <p class="panel-heading">{{Name || "Unnamed Menu"}}
            </p>
            <draggable v-model="Items" >
                    <div v-for="element in Items" :key="element.id" class="draggable-item columns">
                        <div class="column is-2">
                             <button type="button" class="button is-danger" @click="DeletItem(element.id)">
                                    <i class="fa fa-times"></i>
                                </button>
                        </div>
                        <div class="column is-5">
                            <div class="control">

                                <input class="input" type="text"
                                v-model="element.Item.Name" placeholder="Unnamed Menu">
                            </div>
                        </div>
                        <div class="column is-5">
                            <div class="control">
                                <input class="input" type="text"
                                v-model="element.Item.URL" placeholder="Menu URL">
                            </div>
                        </div>
                    </div>
                    <div>
                        <button class="button" @click="SaveMenu">Save Menu</button>
                    </div>
            </draggable>
        </div>
            </div>
        </div>
</div>

</template>

<style lang="scss" scoped>
    .draggable-item{
        padding:10px 20px;
        border:2px solid #cfcfcf;
        margin: 5px;
    }
</style>

<script>
import draggable from 'vuedraggable';
import modal from '../comman/modal';

export default {

    Components:{
        draggable:draggable,
    },

    props:['id'],

    data(){
        return{
            ID: 0,
            Name:"",
            deleteItem:false,
            deleteItemId:-1,
            NewItem:{Name:"",URL:""},
            Items:[],
        }
    },

    mounted(){
        console.log("menu editor has been mounted")
        let vm = this;
        if (!vm.id){
            vm.ID = 0
        }else{
            vm.ID = parseInt(vm.id)
            axios.get('/admin/menu/'+vm.ID+'/edit')
            .then(function(response){
                console.log(response)
                var data = response.data
                vm.ID =  data.menu.ID
                vm.Name =  data.menu.Name
                vm.Items =  data.menu.Items
            }).catch(function(err){
                console.log(err)
            })
        }
    },

    methods:{

        SaveMenu: function(){
            let vm=this;
            // Making save request for menu
            axios.post('/admin/menu/save', {
                  Name: vm.Name,
                  Items: vm.Items,
                  ID: vm.ID,
                })
                .then(function (response) {
                  console.log(response);
                })
                .catch(function (error) {
                  console.log(error);
                });
        },

        AddItem:function(Item){
            console.log(Item)
            let vm=this;
            vm.Items.push(
                {
                    Item:{"Name":Item.Name,"URL":Item.URL},
                    id: vm.Items.length
                }
            );
            vm.NewItem = {Name:"",URL:""};
        },

        DeletItem:function(id){
            let vm = this;
            vm.Items.splice(id,1)
        }
    }

}
</script>
