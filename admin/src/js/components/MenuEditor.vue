<template>
    <div>
        <div class="field">
             <div class="control">
                 <label for="">Menu Name</label>
                 <input type="text" v-model="Name" placeholder="Menu name">
            </div>
        </div>
        <div>
            <input type="text" v-model="NewItem.nav" placeholder="Title">
            <button type="button" @click="AddItem(NewItem)">Add Item</button>
        </div>
        <modal v-if="deleteItem" modal_title="Are you sure?" @clicked="clicked(deleteItemId)">
            <p> Delete item from menu list.</p>
            <p> Action is irreversible.</p>
        </modal>
        <nav class="panel">
            <p class="panel-heading">
                Sample Menu
            </p>
            <draggable v-model="Items" >
                <div v-for="element in Items" :key="element.id" class="draggable-item">
                        <a class="panel-block">
                            <span class="panel-icon">
                                <i class="fa fas fa-times" @click="DeleteItem(element.id)"></i>
                                <i class="fa fas fa-pen" @click="EditItem(element)"></i>
                            </span>
                            {{element.nav}}
                        </a>
                </div>
                <button slot="footer" @click="AddItem(NewItem)">Add</button>
            </draggable>
        </nav>

    </div>
</template>

<style lang="scss" scoped>
.draggable-item{
    padding:10px 20px;
    border:2px solid #cfcfcf;
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
            deleteItem:false,
            deleteItemId:-1,
            Name:'',
            NewItem:{
                nav:"",
            },
            Items:[],
        }
    },

    mounted(){
        console.log("menu editor has been mounted")
    },

    methods:{
        AddItem:function(Item){
            let vm=this;
            Item.id = vm.Items.length
            vm.Items.push(Item);
            vm.NewItem = {
                nav:"",
            }
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
