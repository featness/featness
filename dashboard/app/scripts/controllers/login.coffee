'use strict'

class LoginCtrl
  constructor: (@scope, @auth) ->

  authenticateWithGoogle: ->
    @auth.authenticateWithGoogle(@handleAuthenticated)

  handleAuthenticated: (authResult, profile) ->
    console.log(authResult, profile)

angular.module('dashboardApp')
  .controller 'LoginCtrl', ($scope, AuthService) ->
    $scope.model = new LoginCtrl($scope, AuthService)
