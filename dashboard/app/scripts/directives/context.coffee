'use strict'

class ContextCtrl
    constructor: (@scope, @http, @auth) ->
        @authenticated = @auth.isAuthenticated()

        return unless @authenticated

        @storage = window.sessionStorage
        @loadUserTeams()
        @selectedTeam = null
        @teamText = "Select a Team"

    loadUserTeams: ->
        @http({method: 'GET', url: "http://local.featness.com:8000/teams"}).
            success((data, status, headers, config) =>
                @allTeams = data
                @loadSelectedTeam()
            ).
            error((data, status, headers, config) =>
                @allTeams = []
            )

    setCurrentTeam: (teamId) ->
        @storage.setItem("currentTeam", teamId)
        @loadSelectedTeam()

    loadSelectedTeam: ->
        teamId = @storage.getItem("currentTeam")
        return unless teamId?

        @selectedTeam = teamId
        team = @getTeam(teamId)
        @teamText = "Team: #{ team.name }"

    getTeam: (teamId) ->
        for team in @allTeams
            if team.name == teamId
                return team

        return null

angular.module('dashboardApp')
    .directive('context', ($http, AuthService) ->
        templateUrl: '/views/directives/context.html'
        restrict: 'E'
        link: (scope, element, attrs) ->
            scope.model = new ContextCtrl(scope, $http, AuthService)
    )
