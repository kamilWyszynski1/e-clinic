import Vue from 'vue'
import VueRouter from 'vue-router'
import Profile from '../views/Profile.vue'
import Home from '../views/Home.vue'
import Drugs from '../views/Drugs.vue'
import Drug from '../views/Drug.vue'
import Appointments from '../views/Appointments.vue'
import Appointment from '../views/Appointment.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
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
    name: 'Appointments',
    component: Appointments,
  },
  {
    path: '/appointment/:id',
    name: 'appointment',
    component: Appointment,
  }
]

const router = new VueRouter({
  routes
})

export default router
