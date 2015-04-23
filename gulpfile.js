var requireDir = require('require-dir');

// Require all tasks in gulp, including subfolders
requireDir('./gulp', { recurse: true });