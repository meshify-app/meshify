import Vue from 'vue'
import Vuex from 'vuex'
import auth from "./modules/auth";
import host from "./modules/host";
import server from "./modules/server";
import mesh from "./modules/mesh";

Vue.use(Vuex)

export default new Vuex.Store({
  state: {},
  getters : {},
  mutations: {},
  actions:{},
  modules: {
    auth,
    host,
    mesh,
    server
  }
})
