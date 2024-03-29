angular.module( 'rosy', [
  'templates-app',
  'templates-common',
  'rosy.home',
  'rosy.about',
  'rosy.lecture',
  'rosy.try',
  'rosy.user',
  'rosy.assignment',
  'rosy.newAssignment',
  'ui.router',
  'ui.ace',
  'firebase'
])

.config( function myAppConfig ( $stateProvider, $urlRouterProvider, $httpProvider ) {
  $urlRouterProvider.otherwise( '/home' );
})

.run( function run () {
})

.controller( 'AppCtrl', function AppCtrl ( $scope, $location, $http, $state ) {
  $scope.$on('$stateChangeSuccess', function(event, toState, toParams, fromState, fromParams){
    if ( angular.isDefined( toState.data.pageTitle ) ) {
      $scope.pageTitle = toState.data.pageTitle + ' | rosy' ;
    }
  });

  $scope.$state = $state;

  $scope.user = $http.get('/user');
  $scope.isTeacher = function(type) {
    return (type === 'teacher');
  };
})

;

