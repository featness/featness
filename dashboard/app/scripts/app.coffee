'use strict'

angular
  .module('dashboardApp', [
    'ngCookies',
    'ngResource',
    'ngSanitize',
    'ngRoute',
    'ui.bootstrap'
  ])
  .config ($routeProvider, $locationProvider, $httpProvider) ->
    $locationProvider.html5Mode(true)
    $routeProvider
      .when '/',
        templateUrl: '/views/main.html'
        controller: 'MainCtrl'
        isAuthenticated: true
      .when '/login',
        templateUrl: '/views/login.html'
        controller: 'LoginCtrl'
        isAuthenticated: false
      .when '/teams/join',
        templateUrl: '/views/teams/join.html'
        controller: 'JoinTeamCtrl'
        isAuthenticated: true
      .when '/teams/new',
        templateUrl: '/views/teams/new.html'
        controller: 'NewTeamCtrl'
        isAuthenticated: true
      .otherwise
        redirectTo: '/'

    $httpProvider.interceptors.push('httpRequestInterceptor')

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

  .factory('httpRequestInterceptor', ->
    request: (config) ->
      config['cache'] = false

      storage = window.sessionStorage
      token = storage.getItem("featness-token")

      if token?
        config.headers['X-Auth-Token'] = token

      return config
  )
