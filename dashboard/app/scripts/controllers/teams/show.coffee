'use strict'

class ShowTeamCtrl
    constructor: (@teamId, @scope, @http) ->
        @loadTeamData()

    loadTeamData: ->
        @http({method: 'GET', url: "http://local.featness.com:8000/team/#{ @teamId }"}).
            success((data, status, headers, config) =>
                @team = data
            ).
            error((data, status, headers, config) =>
                # TODO: error msg
                console.log 'error', arguments
            )

angular.module('dashboardApp')
    .controller 'ShowTeamCtrl', ($scope, $routeParams, $http) ->
        teamId = $routeParams.teamId
        $scope.model = new ShowTeamCtrl(teamId, $scope, $http)
