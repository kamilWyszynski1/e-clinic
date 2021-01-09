import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Patient from '../views/Patient.vue'
import Drugs from '../views/Drugs.vue'
import Drug from '../views/Drug.vue'
import Appointment from '../views/Appointment.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/patient/',
    name: 'Patient',
    component: Patient,
  },
  {
    path: '/drugs',
    name: 'Drugs',
    component: Drugs,
  },
  {
    path: '/drug/:id',
    name: 'drug',
    component: Drug,
  },
  {
    path: '/appointments',
    name: 'Appointment',
    component: Appointment,
  }
]

const router = new VueRouter({
  routes
})

export default router
