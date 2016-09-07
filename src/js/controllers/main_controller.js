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

.controller('MainController', ['$scope', '$window','$http', '$rootScope', function($scope, $window, $http, $rootScope){
  $scope.installedApps = [];
  $scope.app = null;

  $http.get('/apps').then(function(response) {
    response.data.forEach(function(app) {
      $scope.installedApps.push(app);
      //{
      //name: 'VLC media player',
      //description: 'VLC is a free and open source cross-platform multimedia player and framework that plays most multimedia files as well as DVDs, Audio CDs, VCDs, and various streaming protocols.',
      //running: true,
      //icon: 'http://i.utdstc.com/icons/256/vlc-media-player-1-0-5.png',
      //admin_location: 'http://127.0.0.1:30001'
      //},
    });
  }, function(response) {
  });

  $scope.actions = function(app) {
    $scope.app = app;
    $rootScope.Ui.turnOn('modal1');
  }
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

  $scope.runApp = function(app) {
    $http.post('/apps/' + app.id + '/start').then(function(response) {
      $scope.app.is_running = true;
    }, function(response) {
    });
  }

  $scope.stopApp = function(app) {
    $http.post('/apps/' + app.id + '/stop').then(function(response) {
      $scope.app.is_running = false;
    }, function(response) {
    });
  }

  $scope.openRemote = function(app) {
    $window.open(app.remote_url);
  }


}])

.controller('StoreController', ['$scope', function($scope) {

}])

.controller('PrefsController', ['$scope', function($scope) {

}]);
