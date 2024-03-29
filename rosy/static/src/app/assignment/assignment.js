angular.module( 'rosy.assignment', [
  'ui.router'
])

.config(function config($stateProvider) {
  $stateProvider.state( 'assignment', {
    url: '/assignment/:id',
    views: {
      'main': {
        controller: 'AssignmentCtrl',
        templateUrl: 'assignment/assignment.tpl.html'
      }
    },
    data: { pageTitle: 'Assignment' }
  });
})

.controller( 'AssignmentCtrl', ['$scope', '$http', '$stateParams', function AssignmentCtrl($scope, $http, $stateParams) {
  var assignmentId = $stateParams.id;
  $scope.output = 'Waiting for submission...';
  $scope.editorOptions = {
    theme: 'tomorrow',
    mode: $scope.language
  };
  $http.get('/assignment/' + assignmentId).
    success(function(data) {
      data.safeDescription = data.description.
        replace(/</g, '&lt;').
        replace(/>/g, '&gt;').
        replace(/&/g, '&amp;').
        replace(/\n/g, '<br/>');
      $scope.language = data.language || 'python';
      $scope.assignment = data;
  });

  $scope.runCode = function(code, language) {
    var data = JSON.stringify({code: code});
    $http.post('/assignment/' + assignmentId + '/submit', data).
      success(function(data) {
        $scope.submitted = true;
        console.log(data);
        $scope.correct = data.correct;
        $scope.output = data.output;
        $scope.assignment.attempts = data.attempts;
      }).
      error(function(data) {
        console.log(data);
        $scope.correct = false;
        $scope.output = 'error running code';
      });
  };
}])

;
