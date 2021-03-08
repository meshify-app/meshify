import ApiService from "../../services/api.service";

const state = {
  error: null,
  accounts: [],
}

const getters = {
  error(state) {
    return state.error;
  },
  accounts(state) {
    return state.accounts;
  },
}

const actions = {
  error({ commit }, error){
    commit('error', error)
  },

  readAll({ commit, dispatch }, email){
    ApiService.get(`/accounts/${email}`)
      .then(resp => {
        commit('accounts', resp)
      })
      .catch(err => {
        commit('error', err)
      })
  },
  
  create({ commit, dispatch }, account){
    ApiService.post(`/accounts/${account.id}`)
      .then(resp => {
        commit('create', resp)
      })
      .catch(err => {
        commit('error', err)
      })
  },

  update({ commit, dispatch }, account){
    ApiService.patch(`/accounts/${account.id}`)
      .then(resp => {
        commit('update', resp)
      })
      .catch(err => {
        commit('error', err)
      })
  },

  delete({ commit }, account){
    ApiService.delete(`/accounts/${account.id}`)
      .then(() => {
        commit('delete', account)
      })
      .catch(err => {
        commit('error', err)
      })
  },

}

const mutations = {
  error(state, error) {
    state.error = error;
  },
  accounts(state, accounts){
    state.accounts = accounts
  },
  create(state, account){
    state.accounts.push(account)
  },
  update(state, account){
    let index = state.accounts.findIndex(x => x.id === account.id);
    if (index !== -1) {
      state.accounts.splice(index, 1);
      state.accounts.push(account);
    } else {
      state.error = "update account failed, not in list"
    }
  },
  delete(state, account){
    let index = state.accounts.findIndex(x => x.id === account.id);
    if (index !== -1) {
      state.accounts.splice(index, 1);
    } else {
      state.error = "delete mesh failed, not in list"
    }
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
