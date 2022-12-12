import { route } from 'quasar/wrappers'
import { createRouter, createMemoryHistory, createWebHistory, createWebHashHistory } from 'vue-router'
import routes from './routes'

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

export default route(function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : (process.env.VUE_ROUTER_MODE === 'history' ? createWebHistory : createWebHashHistory)

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,

    // Leave this as is and make changes in quasar.conf.js instead!
    // quasar.conf.js -> build -> vueRouterMode
    // quasar.conf.js -> build -> publicPath
    history: createHistory(process.env.MODE === 'ssr' ? void 0 : process.env.VUE_ROUTER_BASE)
  })

  // в случаях когда список на экране index.vue загружается до конца, то не срабатывает кнопка back в браузере
  // Ошибка DOMException: A history state object with URL 'http://localhost:8080undefined/'
  // для этого добавляем данную функцию
  Router.beforeEach((to, from, next) => {
    if (!window.history.state.current) window.history.state.current = to.fullPath;
    if (!window.history.state.back) window.history.state.back = from.fullPath;
    return next();
  });

  return Router
})
