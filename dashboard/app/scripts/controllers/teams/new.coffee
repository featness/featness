'use strict'

errorClass = 'has-error'
successClass = 'has-success'

class NewTeamCtrl
    constructor: (@scope, @http, @location, @auth) ->
        @selectedMembers = []
        @nameAvailable = null
        @nameAvailableClass = ''
        @user =
            name: "heynemann"
            picture: 'http://graph.facebook.com/bernardo.heynemann/picture'

        @scope.$watch('model.teamName', (newValue, oldValue) =>
            @validateTeamName(newValue)
        )

        @teamOwner = @auth.getUser()

    getUsersThatMatch: (name) ->
        promise = @http({method: 'GET', url: "http://local.featness.com:8000/users/find?name=#{ name }"}).then((response) ->
            return response.data
        )
        return promise

    validateTeamName: (name) ->
        if not name? or name == ''
            @nameAvailable = null
            @nameAvailableClass = ''
            return

        @http({method: 'GET', url: "http://local.featness.com:8000/teams/available?name=#{ name }"}).
            success((data, status, headers, config) =>
                if data? and data.toLowerCase() == "true"
                    @nameAvailable = true
                    @nameAvailableClass = successClass
                else
                    @nameAvailable = false
                    @nameAvailableClass = errorClass
            ).
            error((data, status, headers, config) =>
                @nameAvailable = false
                @nameAvailableClass = errorClass
            )

    addMember: () =>
        member = @selectedMember
        memberAlreadyAdded = @selectedMembers.indexOf(member) != -1 or member.UserId == @teamOwner.account
        @selectedMember = null
        return if memberAlreadyAdded

        @selectedMembers.push(member)

    removeMember: (member) =>
        index = @selectedMembers.indexOf(member)
        return if index == -1

        @selectedMembers.splice(index, 1)

    createTeam: =>
        formData =
            name: @teamName
            owner: @teamOwner.account
            users: (user.UserId for user in @selectedMembers)

        @http(
            url: "http://local.featness.com:8000/teams/new",
            method: "POST",
            headers: {'Content-Type': 'application/x-www-form-urlencoded'}
            data: $.param(formData)
        ).success((data, status, headers, config) =>
            @location.url("/team/#{ data[0] }")
        ).error((data, status, headers, config) =>
            # TODO: Show error message
        )

        return false

angular.module('dashboardApp')
    .controller 'NewTeamCtrl', ($scope, $http, $location, AuthService) ->
        $scope.model = new NewTeamCtrl($scope, $http, $location, AuthService)
