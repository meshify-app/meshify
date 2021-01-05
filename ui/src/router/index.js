import Vue from 'vue'
import VueRouter from 'vue-router'
import store from "../store";

Vue.use(VueRouter);

const routes = [
  {
    path: '/hosts',
    name: 'hosts',
    component: function () {
      return import(/* webpackChunkName: "Hosts" */ '../views/Hosts.vue')
    },
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/mesh',
    name: 'mesh',
    component: function () {
      return import(/* webpackChunkName: "Mesh" */ '../views/Mesh.vue')
    },
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/users',
    name: 'users',
    component: function () {
      return import(/* webpackChunkName: "Users" */ '../views/Users.vue')
    },
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/server',
    name: 'server',
    component: function () {
      return import(/* webpackChunkName: "Server" */ '../views/Server.vue')
    },
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/',
    name: 'root',
    meta: {
      requiresAuth: false
    }
  },

  {
      // catch all 404 - define at the very end
      path: "*",
      meta: {
        requiresAuth: false
      },
      component: () => import("../views/NotFound.vue")

  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
});

router.beforeEach((to, from, next) => {
  if(to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters["auth/isAuthenticated"]) {
      next()
      return
    }
    next('/')
  } else {
    next()
  }
})

export default router
