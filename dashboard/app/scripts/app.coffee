'use strict'

angular
  .module('dashboardApp', [
    'ngCookies',
    'ngResource',
    'ngSanitize',
    'ngRoute'
  ])
  .config ($routeProvider, $locationProvider) ->
    $locationProvider.html5Mode(true)
    $routeProvider
      .when '/',
        templateUrl: 'views/main.html'
        controller: 'MainCtrl'
      .when '/login',
        templateUrl: 'views/login.html'
        controller: 'LoginCtrl'
      .otherwise
        redirectTo: '/'

