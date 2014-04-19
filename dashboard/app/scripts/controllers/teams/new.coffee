'use strict'

class NewTeamCtrl
    constructor: (@scope) ->
        @selectedMembers = []
        @nameAvailable = null
        @nameAvailableClass = ''
        @user =
            name: "heynemann"
            picture: 'http://graph.facebook.com/bernardo.heynemann/picture'

        @scope.$watch('model.teamName', (newValue, oldValue) =>
            @validateTeamName(newValue)
        )

        @availableMembers = [
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

        existingTeams = [
            'timehome',
            'appdev'
        ]

        name = name.toLowerCase()

        for team in existingTeams
            if team == name
                @nameAvailable = false
                @nameAvailableClass = 'has-error'
                return

        @nameAvailable = true
        @nameAvailableClass = 'has-success'

    addMember: () =>
        @selectedMembers.push(@selectedMember)
        @selectedMember = null

    removeMember: (member) =>
        index = @selectedMembers.indexOf(member)
        return if index == -1

        @selectedMembers.splice(index, 1);


angular.module('dashboardApp')
    .controller 'NewTeamCtrl', ($scope) ->
        $scope.model = new NewTeamCtrl($scope)
