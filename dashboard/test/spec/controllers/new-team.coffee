'use strict'

describe 'Controller: NewTeamCtrl', ->

  # load the controller's module
  beforeEach module 'dashboardApp'

  NewTeamCtrl = {}
  scope = {}

  # Initialize the controller and a mock scope
  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()
    NewTeamCtrl = $controller 'NewTeamCtrl', {
      $scope: scope
    }

  it 'should attach a list of awesomeThings to the scope', ->
    expect(scope.awesomeThings.length).toBe 3
