angular.module('MobyOSAdmin', [
  'ngRoute',
  'mobile-angular-ui',
  'mobile-angular-ui.gestures',
  'MobyOSAdmin.controllers.Main'
])

.config(function($locationProvider) {
  $locationProvider.hashPrefix('');
})

.config(function($routeProvider) {
  $routeProvider.when('/', {templateUrl:'apps.html', controller: 'MainController',  reloadOnSearch: false});
  $routeProvider.when('/apps/:id', {templateUrl:'app.html',  controller: 'AppController', reloadOnSearch: false});
  $routeProvider.when('/store', {templateUrl:'store.html',  controller: 'StoreController', reloadOnSearch: false});
  $routeProvider.when('/profiles', {templateUrl:'profiles.html',  controller: 'ProfileController', reloadOnSearch: false});
  $routeProvider.when('/prefs', {templateUrl:'prefs.html',  controller: 'PrefsController', reloadOnSearch: false});
  $routeProvider.when('/about', {templateUrl:'about.html',  controller: 'AboutController', reloadOnSearch: false});
});
