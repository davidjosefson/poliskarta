"use strict";


/*
TODO:
1. Repetera alla areas och skapa länkar för var och en
2. Fixa routing så att ett klick på en länk tar en till en ny karta och laddar in pins för just det länet
*/

var bigmapControllers = angular.module("bigmapControllers", []);

bigmapControllers.controller("mainController", ["uiGmapGoogleMapApi", "$scope", "$http", "$routeParams", function(uiGmapGoogleMapApi, $scope, $http, $routeParams) {
    var vm = this;
    vm.params = $routeParams;
    vm.markerlist = [];
    vm.events = {};
    vm.areas = {};
    vm.selectedArea = {};

    vm.map = {
        center: {
            latitude: 62.1983366,
            longitude: 17.567198
        },
        zoom: 5
    };

    //Get the list of areas
    $http.get("http://localhost:3000/api/v1/areas/").success(function(data) {
        vm.areas = data;

        for (var i = 0; i < vm.areas.areas.length; i++) {
            //If the route parameter matches one of the fetched areas:
            if (vm.areas.areas[i].value === $routeParams.area) {
                vm.selectedArea = vm.areas.areas[i];
                vm.map.center.latitude = vm.selectedArea.latitude;
                vm.map.center.longitude = vm.selectedArea.longitude;
                vm.map.zoom = 8;
                
                //Get events for the chosen area
                getEvents();
            }
        }
    });

    //Get all events for the chosen area and for each event: get the coordinates in a separate get-request
    function getEvents() {
        $http.get("http://localhost:3000/api/v1/areas/" + vm.selectedArea.value + "/?limit=10").success(function(data) {
            var events = data;
            for (var i = 0; i < events.events.length; i++) {
                var marker = {
                    id: i + 100,
                    eventURI: events.events[i].links[0].href,
                    title: events.events[i].title,
                    description: events.events[i].descriptionShort,
                    coords: {
                    },
                    show: false,
                    onClick: showOrHideMarker(events.events[i].Title, i + 100)
                };
                marker.onClick = showOrHideMarker(marker.title, marker.id);
                getCoordinatesForEvent(i, marker);
            }
        });
    }

    function getCoordinatesForEvent(index, marker) {
        $http.get(marker.eventURI).success(function(data) {
            var event = data;

            if (typeof event.location === "undefined") {
                console.log("Couldn't add marker: " + marker.title + ", no coordinates");
            } else {
                marker.coords.latitude = event.location.latitude;
                marker.coords.longitude = event.location.longitude;
                vm.markerlist.push(marker);

                console.log("Added marker: " + marker.title);
            }
        });
    }

    function showOrHideMarker(title, id) {
        return function() {
            for (var i = 0; i < vm.markerlist.length; i++)  {
                if (id == vm.markerlist[i].id) {
                    vm.markerlist[i].show = !vm.markerlist[i].show;
                } else if (vm.markerlist[i].show) {
                    vm.markerlist[i].show = !vm.markerlist[i].show;
                }
            }
        };
    }

    function closeAllMarkerWindows() {
        return function() {
            console.log("Clicked on the map!");
            for (var i = 0; i < vm.markerlist.length; i++)  {
                if (vm.markerlist[i].show) {
                    vm.markerlist[i].show = !vm.markerlist[i].show;
                }
            }
        };
    }

    //Have to register this as a "vm"-function separately, 
    //otherwise it won't execute when called from areas.html
    vm.closeAllMarkerWindows = closeAllMarkerWindows();

    uiGmapGoogleMapApi.then(function(maps) {

    });
}]);