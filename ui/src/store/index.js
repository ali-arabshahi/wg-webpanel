import Vue from 'vue';
import Vuex from 'vuex';
import sideBar from './sideBar'
import server from './server'
import client from './client'


Vue.use(Vuex);


export default new Vuex.Store({
  modules: {
    sideBar,
    server,
    client,
  }
})
