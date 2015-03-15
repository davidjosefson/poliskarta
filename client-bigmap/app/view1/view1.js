'use strict';

angular.module('bigmapApp.view1', ['ngRoute','uiGmapgoogle-maps'])

.config(['$routeProvider','uiGmapGoogleMapApiProvider', function($routeProvider,uiGmapGoogleMapApiProvider) {
  $routeProvider.when('/view1', {
    templateUrl: 'view1/view1.html',
    controller: 'View1Ctrl'
  });
        uiGmapGoogleMapApiProvider.configure({
        key: '***REMOVED***',
        v: '3.17',
        libraries: 'weather,geometry,visualization'
    });

}])

.controller('View1Ctrl', ['$scope', 'uiGmapGoogleMapApi', function($scope, uiGmapGoogleMapApi) {
$scope.map = { center: { latitude: 45, longitude: -73 }, zoom: 8 };
    uiGmapGoogleMapApi.then(function(maps) {
        
    });
    
}]);