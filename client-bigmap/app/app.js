'use strict';

// Declare app level module which depends on views, and components
var bigmapApp = angular.module('bigmapApp', [
    'uiGmapgoogle-maps',
    'ngRoute',
    'bigmapControllers'
]);

bigmapApp.config(['uiGmapGoogleMapApiProvider', function(uiGmapGoogleMapApiProvider) {
    uiGmapGoogleMapApiProvider.configure({
        key: '***REMOVED***',
        v: '3.17',
        libraries: 'weather,geometry,visualization'
    });
}]);

bigmapApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/', {
        templateUrl: 'areas.html',
        controller: 'mainController'
      }).
      when('/:area', {
        templateUrl: 'areas.html',
        controller: 'mainController'
      }).
      otherwise({
        redirectTo: '/'
      });
}]);
