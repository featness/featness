'use strict'

class JoinTeamCtrl
  constructor: (@scope, @auth) ->
    @loadTeams()

  loadTeams: ->
    @teams = [
      { name: "TimeHome", description: "Time responsável pela manutenção da home da globo.com, do thumbor, cocoon e outros projetos.", projects: 10 }
      { name: "AppDev", description: "Time responsável por projetos que aumentem a produtividade dos desenvolvedores, como o Holmes e o Featness.", projects: 7 }
      { name: "G1", description: "Time responsável pelos projetos necessários ao bom funcionamento do G1.", projects: 10 }
      { name: "GloboEsporte", description: "Time responsável pelos projetos necessários ao bom funcionamento do globoesporte.globo.com.", projects: 10 }
      { name: "Entretenimento", description: "Time responsável pelos projetos necessários ao bom funcionamento do etc.globo.com.", projects: 20 }
    ]

angular.module('dashboardApp')
  .controller 'JoinTeamCtrl', ($scope) ->
    $scope.model = new JoinTeamCtrl($scope)
