"use strict";

// Declare app level module which depends on views, and components
var bigmapApp = angular.module("bigmapApp", ["uiGmapgoogle-maps"]);

bigmapApp.config(["uiGmapGoogleMapApiProvider", function(uiGmapGoogleMapApiProvider) {
    uiGmapGoogleMapApiProvider.configure({
        key: "***REMOVED***",
        v: "3.17",
        libraries: "weather,geometry,visualization"
    });
}]);

bigmapApp.controller("mainController", ["uiGmapGoogleMapApi", "$scope", "$http", function(uiGmapGoogleMapApi, $scope, $http) {
    var vm = this;
    vm.markerlist = [];
    vm.events = {};

    function onClickTest(title, id) {
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

    function closeAllWindows() {
        return function() {
            console.log("Clicked on the map!");
            //            vm.markerlist[1].show = !vm.markerlist[1].show;
            console.log("[1] = " + vm.markerlist[1].title);
            for (var i = 0; i < vm.markerlist.length; i++)  {
                if (vm.markerlist[i].show) {
                    vm.markerlist[i].show = !vm.markerlist[i].show;
                }
            }
        };
    }
 
    vm.closeAll = closeAllWindows();
    vm.test = "hej";

    


    vm.singleCallToAllEvents = function() {
        for (var i = 0; i < vm.markerlist.length; i++) {
          vm.singleEventCall(i);
        }
    };

    vm.singleEventCall = function(index, marker) {
        $http.get(marker.eventURI).success(function(data) {
            var event = data;
            marker.coords.longitude = event.Longitude;
            marker.coords.latitude = event.Latitude;
            vm.markerlist.push(marker);
        });
    };


    $http.get("http://localhost:3000/areas/skane/?limit=50").success(function(data) {
        var events = data;
        for (var i = 0; i < events.Events.length; i++) {
            var marker = {
                id: i + 100,
                eventURI: events.Events[i].EventURI,
                title: events.Events[i].Title,
                description: events.Events[i].DescriptionShort,
                coords: {
                },
                show: false,
                onClick: onClickTest(events.Events[i].Title, i + 100)
            };
            marker.onClick = onClickTest(marker.title, marker.id);
            vm.singleEventCall(i, marker);
        }

        // vm.singleCallToAllEvents();
    });

    //    vm.markerlist = [
    //        {
    //            id: 50,
    //            title: "2015-03-07 19:42, Inbrott, Eslöv",
    //            description: "Inbrott i bostad, Löberöd.",
    //            coords: {
    //                latitude: 55.63755,
    //                longitude: 13.063351
    //            },
    //            show: true,
    //            onClick: onClickTest("2015-03-07 19:42, Inbrott, Eslöv", 50)
    //        }, {
    //            id: 70,
    //            title: "2015-03-07 19:42, Rån, Örkelljunga",
    //            description: "Bla i bla i bla.",
    //            coords: {
    //                latitude: 55.710692,
    //                longitude: 14.293691
    //            },
    //            show: true,
    //            onClick: onClickTest("2015-03-07 19:42, Rån, Örkelljunga", 70)
    //    }, {
    //            id: 90,
    //            title: "2015-03-07 19:42, Hej, Hejsingborg",
    //            description: "MMMM mmmmm oooooo.",
    //            coords: {
    //                latitude: 55.563947,
    //                longitude: 13.552114
    //            },
    //            show: true,
    //            onClick: onClickTest("MMMM mmmmm oooooo.", 90)
    //    }
    //    ];

    vm.testtest = onClickTest();
    //    vm.testtest();



    vm.map = {
        center: {
            latitude: 55.60714,
            longitude: 13.004377
        },
        zoom: 8
            //        events: {
            //            click: //vm.closeAll()
            //                closeAllWindows()
            //        }
    };


    vm.marker2 = {
        id: 2,
        latitude: 45.240727,
        longitude: -73.537372,
        title: "marker2 xxxxxxxxxxxxxx",
        show: true
    };

    vm.marker2.onClick = function() {
        vm.marker2.show = vm.marker2.show;
    };



    vm.windowOptions = {
        visible: false
    };

    vm.onClick = function() {
        vm.windowOptions.visible = !vm.windowOptions.visible;
        vm.test = "mjeee";
    };

    vm.closeClick = function() {
        vm.windowOptions.visible = false;
    };

    vm.title = "Window Title!";

    uiGmapGoogleMapApi.then(function(maps) {

    });
}]);