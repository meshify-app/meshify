import Vue from 'vue'
import Vuex from 'vuex'
import auth from "./modules/auth";
import host from "./modules/host";
import server from "./modules/server";
import mesh from "./modules/mesh";
import user from "./modules/users"
import account from "./modules/account"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {},
  getters : {},
  mutations: {},
  actions:{},
  modules: {
    account,
    auth,
    host,
    mesh,
    user,
    server
  }
})
