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
    $http.post('http://tryrosy.com:9000/' + language.name.toLowerCase(), {code: escape(code)}).
      success(function(data) {
        console.log(data);
        $scope.output = data;
      }).
      error(function(data) {
        console.log(data);
        $scope.output = 'error running code';
      });
  };
}])

;
