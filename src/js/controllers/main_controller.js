angular.module('MobyOSAdmin.controllers.Main', [])

.controller('MainController', ['$scope', '$window','$http','$location', function($scope, $window, $http, $location){
    $scope.installedApps = [];

    $http.get('/apps/installed').then(function(response) {
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


    $scope.runApp = function(app) {
        app.running = true;
        $window.location.href = app.admin_location;
    }

    $scope.stopApp = function(app) {
        app.running = false;
    }

}]);
