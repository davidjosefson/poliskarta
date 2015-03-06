(function () {
    var app = angular.module('myApp', []);
    app.config(function ($sceDelegateProvider) {
        $sceDelegateProvider.resourceUrlWhitelist([
   // Allow same origin resource loads.
   'self',
    'http://192.168.1.78:3000/',        
   // Allow loading from our assets domain. 
   'https://www.google.com/maps/embed/**']);
    });
    
    
    app.filter("getMapUrl", function () {
        return function (input) {
            return "https://www.google.com/maps/embed/v1/place?key=***REMOVED***&q=s" + input;
        }
    });
    app.filter("getStaticMapUrl", function () {
        return function (input) {
//        return function (input, lat, lng) {
            this.output = "http://maps.googleapis.com/maps/api/staticmap?center=";
            this.output += input;
            this.output += "&zoom=12&&scale=2&size=300x200&maptype=roadmap&sensor=false&key=***REMOVED***";
            
            //Detta ska läggas till när vi har koordinater
            //this.output += "&markers=color:red%7ccolor:red%7clabel:c%7c"
            //this.output += lat + "," + lng
            
            return this.output;
        }
    });
    
    //"http://maps.googleapis.com/maps/api/staticmap?center=malm%C3%B6+ellstorp&zoom=12&size=300x300&maptype=roadmap&markers=color:red%7ccolor:red%7clabel:c%7c55.605422,13.031286&sensor=false&key=***REMOVED***"

    app.controller('EventController', ["$http", function ($http) {
        var ctrl = this;
        ctrl.events = {};
        this.selected = "skane";
        this.url = "http://192.168.1.78:3000/";
        this.mapURL = "https://www.google.com/maps/embed/v1/place?key=***REMOVED***&q=s";

        $http.get(this.url + this.selected).success(function (data) {
            ctrl.events = data;
        });

        this.getData = function (input) {
            this.selected = input;
            $http.get(this.url + input).success(function (data) {
                ctrl.events = {};
                ctrl.events = data;
            });
        };

    }]);

    app.directive("event", function () {
        return {
            restrict: 'E',
            templateUrl: 'event.html'
        };
    });

})();