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
  dadataToken: '[[dadataToken]]',
  // yandexMetrikaId: 54433825,
  breadcrumbIcons: {
    [[breadcrumbIcons]]
  },
  [[telegramConfig]]
  tablesForTask: [[codoGeneratedTablesForTask]],
}
