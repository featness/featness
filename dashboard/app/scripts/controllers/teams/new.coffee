'use strict'

errorClass = 'has-error'
successClass = 'has-success'

class NewTeamCtrl
    constructor: (@scope, @http, @auth) ->
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
        return [
            {name: 'guilhermef', picture: 'http://graph.facebook.com/guilherme.souza/picture'},
            {name: 'rfloriano', picture: 'http://graph.facebook.com/rafael.floriano/picture'},
            {name: 'scorphus', picture: 'http://graph.facebook.com/pablo.aguiar/picture'},
            {name: 'metal', picture: 'http://graph.facebook.com/marcelo.vieira/picture'}
        ]

    validateTeamName: (name) ->
        if not name? or name == ''
            @nameAvailable = null
            @nameAvailableClass = ''
            return

        @http({method: 'GET', url: "http://local.featness.com:8000/teams/available?name=#{ name }"}).
            success((data, status, headers, config) =>
                if data? and data
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
        memberAlreadyAdded = @selectedMembers.indexOf(member) != -1
        @selectedMember = null
        return if memberAlreadyAdded

        @selectedMembers.push(member)

    removeMember: (member) =>
        index = @selectedMembers.indexOf(member)
        return if index == -1

        @selectedMembers.splice(index, 1)


angular.module('dashboardApp')
    .controller 'NewTeamCtrl', ($scope, $http, AuthService) ->
        $scope.model = new NewTeamCtrl($scope, $http, AuthService)
