angular.module('MobyOSAdmin.controllers.Main', [])

.controller('MainController', ['$scope', '$window', function($scope, $window){
    $scope.installedApps = [];

    function getInstalledApps() {
        $scope.installedApps = [
            {
                name: 'VLC media player',
                description: 'VLC is a free and open source cross-platform multimedia player and framework that plays most multimedia files as well as DVDs, Audio CDs, VCDs, and various streaming protocols.',
                running: true,
                icon: 'http://i.utdstc.com/icons/256/vlc-media-player-1-0-5.png',
                admin_location: 'http://127.0.0.1:30001'
            },
            {
                name: 'Netflix',
                description: 'Netflix is the world’s leading subscription service for watching TV episodes and movies on your phone. This Netflix mobile application delivers the best experience anywhere, anytime.',
                running: false,
                icon: 'http://icons.iconarchive.com/icons/chrisbanks2/cold-fusion-hd/128/netflix-icon.png',
                admin_location: 'http://127.0.0.1:30002'
            },
            {
                name: 'Spotify',
                description: 'Spotify is a digital music service that gives you access to millions of songs.',
                running: false,
                icon: 'http://www.iconarchive.com/download/i98446/dakirby309/simply-styled/Spotify.ico',
                admin_location: 'http://127.0.0.1:30003'
            }
        ];
    }

    $scope.runApp = function(app) {
        app.running = true;
        $window.location.href = app.admin_location;
    }

    $scope.stopApp = function(app) {
        app.running = false;
    }

    getInstalledApps();
}]);
