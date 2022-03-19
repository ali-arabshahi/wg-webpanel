// import Vue from 'vue'
import Vue from 'vue';
import Vuesax from 'vuesax';
import App from './App.vue';
import router from './router';
import store from './store';
import 'vuesax/dist/vuesax.css';
import 'material-icons/iconfont/material-icons.css';
import axios from 'axios'
//-------------------------------
Vue.use(Vuesax, {
  // options here
});

Vue.config.productionTip = false;

const base = axios.create({
  baseURL: process.env.VUE_APP_BASE_URL
});

Vue.prototype.$http = base;

const token = sessionStorage.getItem('token')

base.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    console.log("midelwar error")
    if (!error.response) {
      Vue.prototype.$vs.notify({
        title: "save",
        text: "server not response",
        color: "danger",
        position: "top-right",
      });
      Vue.prototype.$vs.loading.close();
      console.log("server side error : ", error);
      return Promise.reject(error);
    }
    if (error.response) {
      if (error.response.status === 401) {
        Vue.prototype.$vs.notify({
          title: "error",
          text: error.response.data.Message,
          color: "danger",
          position: "top-right",
        });
        router.push({ name: "LoginPage" }).catch(() => { });
        Vue.prototype.$vs.loading.close();
        return Promise.reject(error);
      }
      Vue.prototype.$vs.notify({
        time: 10000,
        title: "error",
        text: error.response.data.Message,
        color: "danger",
        position: "top-right",
      });
    }
    Vue.prototype.$vs.loading.close();
    return Promise.reject(error);
  });

if (token) {
  Vue.prototype.$http.defaults.headers.common['token'] = token
}
//-------------------------------

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
