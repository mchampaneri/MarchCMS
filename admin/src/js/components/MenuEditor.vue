<template>
    <div>
        <div class="field">
            <div class="field">
              <div class="control">
                <input class="input" type="text"
                v-model="MenuName" placeholder="Menu Name">
              </div>
            </div>
        </div>
        <div>
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
        <button type="button" class="button is-normal" @click="AddItem(NewItem)">Add Item</button>
        </div>
        <modal v-if="deleteItem" modal_title="Are you sure?" @clicked="clicked(deleteItemId)">
            <p>Delete item from menu list.</p>
            <p>Action is irreversible.</p>
        </modal>
        <nav class="panel">
            <p class="panel-heading">{{MenuName}}
            </p>
            <draggable v-model="Items" >
                <div v-for="element in Items" :key="element.id" class="draggable-item">
                    <div>
                        <div class="field">
                            <div class="control">
                                <input class="input" type="text"
                                v-model="element.Item.Name" placeholder="Menu Title">
                            </div>
                        </div>
                         <div class="field">
                            <div class="control">
                                <input class="input" type="text"
                                v-model="element.Item.URL" placeholder="Menu URL">
                            </div>
                        </div>
                    </div>
                </div>
            </draggable>
        </nav>
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

    data(){
        return{
            MenuName:"",
            deleteItem:false,
            deleteItemId:-1,
            NewItem:{Name:"",URL:""},
            Items:[],
        }
    },

    mounted(){
        console.log("menu editor has been mounted")
    },

    methods:{
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

        DeleteItem:function(id){
            let vm = this;
            alert("delete modal should be visible.")
            vm.deleteItem = true;
            vm.deleteItemId = id;
        },

        clicked:function(id){
            alert('you got it')
            vm.Items.splice(id,1)
        }
    }

}
</script>
