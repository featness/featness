'use strict'

class LoginCtrl
  constructor: (@scope, @location, @auth) ->

  authenticateWithGoogle: ->
    @auth.authenticateWithGoogle(@handleAuthenticated)

  handleAuthenticated: (email, token) =>
    if email? and token?
      @location.path('/')
    else
      alert("WOOT?")

angular.module('dashboardApp')
  .controller 'LoginCtrl', ($scope, $location, AuthService) ->
    $scope.model = new LoginCtrl($scope, $location, AuthService)
