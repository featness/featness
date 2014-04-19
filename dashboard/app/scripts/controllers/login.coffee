'use strict'

class LoginCtrl
  constructor: (@scope, @location, @auth) ->
    @auth.initializeFacebook(@handleAuthenticated)

  authenticateWithGoogle: ->
    @auth.authenticateWithGoogle(@handleAuthenticated)

  authenticateWithFacebook: ->
    @auth.authenticateWithFacebook(@handleAuthenticated)

  handleAuthenticated: (email, name, token) =>
    if email? and token?
      @location.path('/')
    else
      @authenticationFailed()

  authenticationFailed: ->
    alert('WOOT?')  # TODO CHANGE THIS

angular.module('dashboardApp')
  .controller 'LoginCtrl', ($scope, $location, AuthService) ->
    $scope.model = new LoginCtrl($scope, $location, AuthService)
