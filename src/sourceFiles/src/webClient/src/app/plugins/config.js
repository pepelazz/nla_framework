export default {
  appName: '[[appName]]',
  uiAppName: '[[uiAppName]]',
  apiUrl: () => process.env.NODE_ENV === 'development' ? 'http://localhost:[[webPort]]' : 'https://[[url]]',
  wsUrl: () => process.env.NODE_ENV === 'development' ? 'ws://localhost:[[webPort]]' : 'wss://[[url]]',
  isEmailAuth: {
    firstName: true,
    lastName: true,
  },
  logoSrc: '[[logoSrc]]',
  dadataToken: '1cf3a086e3dbe1306ed142d2b5fbc1b324b8eb60',
  // yandexMetrikaId: 54433825,
  breadcrumbIcons: {
    [[breadcrumbIcons]]
  }
}
