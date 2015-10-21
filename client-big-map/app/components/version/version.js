'use strict';

angular.module('bigmapApp.version', [
  'bigmapApp.version.interpolate-filter',
  'bigmapApp.version.version-directive'
])

.value('version', '0.1');
