<template>
    <div>
        <div class="field">
             <div class="control">
                 <label for="">Menu Name</label>
                 <input type="text" v-model="Name" placeholder="Menu name">
            </div>
        </div>
        <div>
            <input type="text" v-model="NewItem.Title" placeholder="Title">
            <input type="text" v-model="NewItem.URL" placeholder="URL">
            <button type="button" @click="AddItem(NewItem)">Add Item</button>
        </div>
        <div>
            <li v-for="(Item,Index) in Items" :key="Index">
                {{Item.Title}} | {{Item.URL}}
            </li>
        </div>
        <draggable v-model="myArray" draggable=".item">
            <div v-for="element in myArray" :key="element.id" class="item">
                    {{element.name}}
                </div>
            <button slot="footer" @click="addPeople">Add</button>
        </draggable>
    </div>
</template>

<script>
import draggable from 'vuedraggable';

export default {

    Components:{
        draggable:draggable,
    },

    data(){
        return{
            Name:'',
            NewItem:{
                Title:"",
                URL:"",
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
            vm.Items.push(Item);
            vm.NewItem = {
                Title:"",
                URL:"",
            }
        }
    }

}
</script>
