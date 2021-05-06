import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    user: JSON.parse(localStorage.getItem('user') || '{"id":0, "name":null, "uuid":null}'),
  },
  mutations: {
    updateUser: (state, user) => {
      state.user = user;
      localStorage.setItem('user', JSON.stringify(user));
    },
    removeUser: (state) => {
      state.user = null;
      localStorage.setItem('user', '{"id":0, "name":null, "uuid":null}');
    },
  },
  getters: {
    user: (state) => state.user,
  },
  actions: {}
})
