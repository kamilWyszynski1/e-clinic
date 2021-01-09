import Vue from 'vue'
import App from './App.vue'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import VueFormulate from '@braid/vue-formulate'

Vue.use(VueFormulate)

import VueFormGenerator from 'vue-form-generator'
import 'vue-form-generator/dist/vfg.css'

Vue.use(VueFormGenerator)
// Install BootstrapVue
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)

Vue.config.productionTip = false

import VueRouter from 'vue-router'
import NavBar from "./components/layout/NavBar.vue";

Vue.use(VueRouter)

// const NotFound = { template: '<p>Page not found</p>' }
const Home = App

// const routes = {
//   '*': Home,
//   '/profile': Home,
//   '/appointments': Home,
//   '/drugs/*': Home,
//   '/drugs/:id': Home,
// }

const Foo = { template: '<div>foo</div>' }

const routes2 = [
  // { path: "*", component: Home },
  { path: "/", component: Home },
  { path: '/foo', component: Foo },
  { path: '/nav', component: NavBar },

]

const router = new VueRouter({mode: 'history', routes: routes2 })


new Vue({
  router,
  el: '#app',
  template: `<div id="app">
  <div class="hello">
    <nav
      class="navbar navbar-expand-lg navbar-light bg-light justify-content-center"
    >
      <a class="navbar-brand" href="/"
        ><router-link to="/">Home</router-link></a
      >
      <a class="navbar-brand" ><router-link to="/">Profile</router-link></a>
      <a class="navbar-brand" ><router-link to="/">Appointments</router-link></a>
      <a class="navbar-brand" ><router-link to="/">Drug</router-link></a>
    </nav>
  </div>
</div>`
}).$mount('#app')