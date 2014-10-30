$(document).ready(function() {

    // Smooth scrolling
    $('a[href*=#]:not([href=#])').click(function() {
        if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
            var target = $(this.hash);
            target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
            if (target.length) {
                $('html,body').animate({
                    scrollTop: target.offset().top
                }, 1000);
                return false;
            }
        }
    });

});

var app =angular.module('hashDbApp', ['ngResource']);

app.config(['$interpolateProvider', function($interpolateProvider) {
  $interpolateProvider.startSymbol('[[');
  $interpolateProvider.endSymbol(']]');
}]);

app.controller('HomeController', ['$scope', 'Links', function($scope, Links) {
        $scope.q = "";

        $scope.shortlink = "";
        $scope.links = []

        $scope.submit = function(form) {
            if (form.$invalid) {
                return false;
            }

            Links.save({ url: this.url, message: this.message }, function(data) {
                $scope.link.url = "";
                $scope.link.message = "";

                $scope.links.push(data);
            });
        }
    }]);

app.factory('Links',  ['$resource', function($resource) {
        return $resource('/', {}, {
            save: { method: 'POST', params: { } },
        });
    }]);
