'use strict'

class LoginCtrl
  constructor: (@scope, @auth) ->

  authenticateWithGoogle: ->
    @auth.authenticateWithGoogle(@handleAuthenticated)

  handleAuthenticated: (email, token) ->
    console.log(email, token)

angular.module('dashboardApp')
  .controller 'LoginCtrl', ($scope, AuthService) ->
    $scope.model = new LoginCtrl($scope, AuthService)
