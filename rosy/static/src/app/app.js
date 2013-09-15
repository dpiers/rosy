angular.module( 'rosy', [
  'templates-app',
  'templates-common',
  'rosy.home',
  'rosy.about',
  'rosy.try',
  'ui.state',
  'ui.route',
  'ui.ace'
])

.config( function myAppConfig ( $stateProvider, $urlRouterProvider ) {
  $urlRouterProvider.otherwise( '/home' );
})

.run( function run () {
})

.controller( 'AppCtrl', function AppCtrl ( $scope, $location ) {
  $scope.$on('$stateChangeSuccess', function(event, toState, toParams, fromState, fromParams){
    if ( angular.isDefined( toState.data.pageTitle ) ) {
      $scope.pageTitle = toState.data.pageTitle + ' | rosy' ;
    }
  });
})

;

