angular.module('MobyOSAdmin', [
  'ngRoute',
  'mobile-angular-ui',
  'MobyOSAdmin.controllers.Main'
])

.config(function($routeProvider) {
  $routeProvider.when('/', {templateUrl:'apps.html',  reloadOnSearch: false});
});
