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

.controller( 'TryCtrl', ['$http', function TryCtrl($scope) {
  $scope.languages = [
    {name: 'Ruby'},
    {name: 'Python'},
    {name: 'JavaScript'}
  ];
  $scope.language = $scope.languages[1];

  //$scope.code = 'console.log("hello");';
  $scope.code = 'print \'hello\'';

  $scope.editorOptions = {
    theme: 'tomorrow',
    mode: $scope.language.name.toLowerCase()
  };

  $scope.runCode = function(code, language) {
    $http.post('http://50.23.95.61/' + language, {code: code}).
      success(function(data) {
        console.log(data);
      });
  };
}])

;
