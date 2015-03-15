'use strict';

// Declare app level module which depends on views, and components
angular.module('bigmapApp', [
  'ngRoute',
  'bigmapApp.view1',
  'bigmapApp.view2',
  'bigmapApp.version'
]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.otherwise({redirectTo: '/view1'});
}]);

//config(['$routeProvider', function($routeProvider) {
//  $routeProvider.otherwise({redirectTo: '/view1'});
//}]);
