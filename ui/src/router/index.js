import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '../store';
//----------------------------------
import Main from '../views/Home.vue';
import ServerConfig from '../views/ServerConfig.vue';
import ClientConfig from '../views/ClientConfig.vue';
import Dashboard from '../views/Dashboard.vue';
import Login from '../views/Login.vue';



Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'LoginPage',
    component: Login,
  },
  {
    path: '/home',
    name: 'Main',
    component: Main,
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: '/home',
        name: 'Dashboard',
        index: 1,
        component: Dashboard,
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: '/server',
        name: 'ServerConfig',
        index: 2,
        component: ServerConfig,
        meta: {
          requiresAuth: true
        },
      },
      {
        path: '/client',
        name: 'ClientConfig',
        index: 3,
        component: ClientConfig,
        meta: {
          requiresAuth: true
        },
      },
    ],
  }
]



const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  store.commit("closeSidebar");
  if (to.matched.some(record => record.meta.requiresAuth)) {
    console.log(store.getters.getUserLogin.isLogiedIn)
    console.log("CHECK USER", store.getters.getUserLogin.isLogiedIn)
    if (store.getters.getUserLogin.isLogiedIn) {
      console.log("USER IS LOGIN")
      next()
      return
    }
    next('/')
  } else {
    next()
  }
})


export default router;

