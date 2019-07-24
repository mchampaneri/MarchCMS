<template>

<nav class="navbar is-fixed-top has-background-white" role="navigation" aria-label="main navigation">
    <div class="container">
        <div class="navbar-brand">
            <a class="navbar-item" href="/">
                <b>AdminPanel</b>
            </a>
            <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false">
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
            </a>
        </div>

        <div class="navbar-menu navbar-end">

            <a class="navbar-item" href="/admin">
                Home
            </a>

            <a class="navbar-item" href="/admin/pages/list">
                Pages
            </a>

            <a class="navbar-item" href="/admin/posts/list">
                Posts
            </a>

            <a class="navbar-item" href="/admin/menus/list">
                Menus
            </a>

            <a class="navbar-item" href="/admin/assets/images">
                Assests
            </a>

            <a class="navbar-item" href="/admin/users/list">
                Users
            </a>

            <a class="navbar-item" href="/admin/settings">
                Settings
            </a>

            <!-- Adminitstaration Menu -->
            <div class="navbar-item has-dropdown "
                :class="[ profileDropDown ? 'is-active':'']"
                @click="toggleProfileMenu">

                <a class="navbar-link">
                    <div class="image is-32X32">
                        <img class="is-rounded" src="https://bulma.io/images/placeholders/128x128.png">
                    </div>
                </a>

                <div class="navbar-dropdown">
                    <a class="navbar-item">
                       Hi, &nbsp <b>{{ name }}</b>
                    </a>
                    <hr class="navbar-divider">
                    <!-- User & Settings  -->
                    <a class="navbar-item">
                        Profile
                    </a>
                    <a class="navbar-item" href="/signout">
                        SignOut
                    </a>
                </div>
            </div>

        </div>
    </div>
</nav>
</template>

<script>
    export default{

        mounted(){
            console.log("Navbar moutned")
        },

        props:["name","id","picture"],

        data(){
            return{
                profileDropDown: false
            }
        },
        methods:{
            toggleProfileMenu(){
                let vm = this;
                vm.profileDropDown = (vm.profileDropDown) ? false : true;
            }
        }
    }
</script>