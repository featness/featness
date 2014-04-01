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
        isAuthenticated: true
      .when '/login',
        templateUrl: 'views/login.html'
        controller: 'LoginCtrl'
        isAuthenticated: false
      .otherwise
        redirectTo: '/'

  .run(($rootScope, $location, $route, AuthService) ->
    $rootScope.$on("$locationChangeStart", (event, next, current) ->
      nextPath = next \
        .replace("#{$location.protocol()}://#{$location.host()}:#{$location.port()}", "") \
        .replace("#{$location.protocol()}://#{$location.host()}", "")

      requiresAuthentication = $route.routes[nextPath].isAuthenticated
      if (requiresAuthentication and not AuthService.isAuthenticated())
        $location.url("/login")
    )
  )
