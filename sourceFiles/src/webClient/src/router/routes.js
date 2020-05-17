const routes = [
  {path: '/', component: () => import(`../app/components/home.vue`)},
  {
    path: '/check_user_email',
    component: () => import(`../app/components/auth/email/checkUserEmail.vue`),
  },
  {
    path: '/email_auth_recover_password',
    component: () => import(`../app/components/auth/email/emailAuthRecoverPassword.vue`),
  },
  {path: '/users', component: () => import(`../app/components/users/index.vue`)},
  {path: '/users/:id', component: () => import(`../app/components/users/item.vue`), props: true},
  {path: '/profile', component: () => import(`../app/components/currentUser/profile.vue`)},
  {path: '/task', component: () => import(`../app/components/task/index.vue`)},
  {path: '/task/:id', component: () => import(`../app/components/task/item.vue`), props: true},
  {path: '/taskType', component: () => import(`../app/components/taskType/index.vue`)},
  {path: '/taskType/:id', component: () => import(`../app/components/taskType/item.vue`), props: true},
  // for codeGenerate ##routes_slot1
]

// Always leave this as last one
if (process.env.MODE !== 'ssr') {
  routes.push({
    path: '*',
    component: () => import(`../app/components/error404.vue`)
  })
}

export default routes
