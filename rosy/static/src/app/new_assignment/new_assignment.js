angular.module( 'rosy.newAssignment', [
  'ui.router'
])

.config(function config($stateProvider) {
  $stateProvider.state( 'newAssignment', {
    url: '/assignments/new',
    views: {
      'main': {
        controller: 'NewAssignmentCtrl',
        templateUrl: 'new_assignment/new_assignment.tpl.html'
      }
    },
    data: { pageTitle: 'New Assignment' }
  });
})

.controller( 'NewAssignmentCtrl', ['$scope', '$http', '$state', function NewAssignmentCtrl($scope, $http, $state) {
  $scope.form = {
    title: '',
    language: 'Python',
    description: '',
    output: '',
    code: '# students will work from the code here\n# Press the test button to execute'
  };

  $scope.editorOptions = {
    theme: 'tomorrow',
    mode: $scope.form.language.toLowerCase
  };

  $scope.output = 'Waiting for submission...';

  $scope.languages = [
    'Ruby',
    'Python',
    'JavaScript',
    'Haskell',
    'C++',
    'Go'
  ];

  $scope.submitForm = function(form) {
    console.log('submitting assignment', form);
    data = JSON.stringify(form);
    $http.post('/assignments/new', data).
      success(function(data) {
        console.log('new assignment: ', data);
        $state.go('assignment', {id: data.id});
    }).
      error(function(data) {
        console.log(data);
    });
  };

  $scope.runCode = function(code, language) {
    var data = JSON.stringify({code: code});
    $http.post('/eval/' + language, data).
      success(function(data) {
        console.log(data);
        $scope.output = data.output;
      }).
      error(function(data) {
        console.log(data);
        $scope.output = 'error running code';
      });
  };
}])

;
