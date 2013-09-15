angular.module( 'rosy.try', [
  'ui.state'
])

.config(function config($stateProvider) {
  $stateProvider.state( 'try', {
    url: '/try',
    views: {
      'main': {
        controller: 'TryCtrl',
        templateUrl: 'try/try.tpl.html'
      }
    },
    data: { pageTitle: 'Try It!' }
  });
})

.controller( 'TryCtrl', ['$scope', '$http', function TryCtrl($scope, $http) {
  $scope.languages = [
    {name: 'Ruby'},
    {name: 'Python'},
    {name: 'JavaScript'}
  ];
  $scope.language = $scope.languages[1];

  //$scope.code = 'console.log("hello");';
  $scope.code = 'print \'hello\'';
  $scope.output = 'Waiting for code...';

  $scope.editorOptions = {
    theme: 'tomorrow',
    mode: $scope.language.name.toLowerCase()
  };

  $scope.runCode = function(code, language) {
    var data = JSON.stringify({code: code});
    $http.post('/eval/' + language.name.toLowerCase(), data).
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
