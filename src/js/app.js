angular.module('MobyOSAdmin', [
  'ngRoute',
  'mobile-angular-ui',
  'mobile-angular-ui.gestures',
  'MobyOSAdmin.controllers.Main'
])

.config(function($routeProvider) {
  $routeProvider.when('/', {templateUrl:'apps.html',  reloadOnSearch: false});
  $routeProvider.when('/apps/:id', {templateUrl:'app.html',  controller: 'AppController', reloadOnSearch: false});
  $routeProvider.when('/store', {templateUrl:'store.html',  controller: 'StoreController', reloadOnSearch: false});
  $routeProvider.when('/prefs', {templateUrl:'prefs.html',  controller: 'PrefsController', reloadOnSearch: false});
});
