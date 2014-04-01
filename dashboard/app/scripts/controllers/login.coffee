'use strict'

class LoginCtrl
  constructor: (@scope) ->

angular.module('dashboardApp')
  .controller 'LoginCtrl', ($scope) ->
    $scope.model = new LoginCtrl($scope)
