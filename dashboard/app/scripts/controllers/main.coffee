'use strict'

class MainCtrl
  constructor: (@scope) ->

angular.module('dashboardApp')
  .controller 'MainCtrl', ($scope) ->
    $scope.model = new MainCtrl($scope)
