
/**
 * First we will load all of this project's JavaScript dependencies which
 * includes Vue and other libraries. It is a great starting point when
 * building robust, powerful web applications using Vue and Laravel.
 */

require('./bootstrap');

window.Vue = require('vue');

/**
 * Next, we will create a fresh Vue application instance and attach it to
 * the page. Then, you may begin adding components to this application
 * or customize the JavaScript scaffolding to fit your unique needs.
 */

Vue.component("Modal",require("./comman/modal.vue"))

Vue.component("PageEditor",require("./components/PageEditor.vue"))
Vue.component("PostEditor",require("./components/PostEditor.vue"))
Vue.component("MenuEditor",require("./components/MenuEditor.vue"))
Vue.component("AssetsEditor",require("./components/AssetsEditor.vue"))

Vue.component("AssetsImage",require("./components/AssetsImage.vue"))
Vue.component("AssetsDocument",require("./components/AssetsDocument.vue"))

Vue.component("ThemeConfig",require("./components/ThemeConfig.vue"))
Vue.component("SiteConfig",require("./components/SiteConfig.vue"))

Vue.component("Navbar", require("./components/Navbar.vue"))

Vue.component("Login", require("./components/Login.vue"))
const app = new Vue({
    el: '#app'
});
