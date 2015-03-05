(function () {
    var app = angular.module('myApp', []);
    app.config(function ($sceDelegateProvider) {
        $sceDelegateProvider.resourceUrlWhitelist([
   // Allow same origin resource loads.
   'self',
   // Allow loading from our assets domain.  Notice the difference between * and **.
   'https://www.google.com/maps/embed/**']);
    });
    
    
    app.filter("getMapUrl", function () {
        return function (input) {
            return "https://www.google.com/maps/embed/v1/place?key=***REMOVED***&q=s" + input;
        }
    });
    app.filter("getStaticMapUrl", function () {
        return function (input) {
            this.output = "http://maps.googleapis.com/maps/api/staticmap?center=";
            this.output = this.output + input;
            this.output = this.output + "&zoom=12&size=300x200&maptype=roadmap&sensor=false&key=***REMOVED***";
            return this.output;
        }
    });
    
    //"http://maps.googleapis.com/maps/api/staticmap?center=malm%C3%B6+ellstorp&zoom=12&size=300x300&maptype=roadmap&markers=color:red%7ccolor:red%7clabel:c%7c55.605422,13.031286&sensor=false&key=***REMOVED***"

    app.controller('EventController', ["$http", function ($http) {
        var ctrl = this;
        ctrl.events = {};
        this.selected = "skane";
        this.url = "http://localhost:3000/";
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