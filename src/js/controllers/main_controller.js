angular.module('MobyOSAdmin.controllers.Main', [])

.directive('back', ['$window', function($window) {
  return {
    restrict: 'A',
    link: function (scope, elem, attrs) {
      elem.bind('click', function () {
        $window.history.back();
      });
    }
  };
}])

.controller('MainController', ['$scope', '$window','$http', '$rootScope', '$location', function($scope, $window, $http, $rootScope, $location){
  var self = this;
  $scope.app = null;


  this.runApp = function(app) {
    $http.post('/apps/' + app.id + '/start').then(function(response) {
      app.is_running = true;
    }, function(response) {
    });
  };

  this.stopApp = function(app) {
    $http.post('/apps/' + app.id + '/stop').then(function(response) {
      app.is_running = false;
    }, function(response) {
    });
  };

  this.openRemote = function(app) {
    $window.open(app.remote_url);
  };


  this.uninstallApp = function(app) {
    $http.delete('/apps/' + app.id).then(function(response){
      self.getApps().then(function(response) {
        $location.path('/');
      });
    }, function(response){
    });
  };

  $scope.actions = function(app) {
    $scope.app = app;
    $rootScope.Ui.turnOn('modal1');
  }

  this.getApps = function() {
    $scope.installedApps = [];
    return $http.get('/apps').then(function(response) {
      response.data.forEach(function(app) {
        $scope.installedApps.push(app);
      });
    }, function(response) {
    });
  };

  this.getApps();


}])

.directive('ngRightClick', ['$parse', function($parse) {
    return function(scope, element, attrs) {
        var fn = $parse(attrs.ngRightClick);
        element.bind('contextmenu', function(event) {
            scope.$apply(function() {
                event.preventDefault();
                fn(scope, {$event:event});
            });
        });
    };
}])


.controller('AppController', ['$scope', '$window','$http','$routeParams', function($scope, $window, $http, $routeParams){
  $scope.app = null;

  var appId = $routeParams.id;

  $http.get('/apps/' + appId).then(function(response) {
    $scope.app = response.data;
  }, function(response) {
  });



}])

.controller('StoreController', ['$scope', function($scope) {

}])

.controller('PrefsController', ['$scope', function($scope) {

}]);
