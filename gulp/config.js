var client = "./client/";
var server = "./server/";

var stage = "./staging/";
var dist = "./dist/";

var assets = {
		styles: "./src/styles/",
		images: "./src/images/",
		scripts: "./src/scripts/",
		svg: "./src/svg/",
		cache: "./src/cache/"
};

var	staging = {
		styles: stage + "css/",
		images: stage + "img/",
		scripts: stage + "js/",
};

module.exports = {
	clean: {
		staging: stage,
		dist: dist
	},

	templates: {
		src: client + "app/**/*.html",
		dst: assets.cache
	},

	scripts: {
		vendors: client + "vendor/**/*.js",
		src: [
			client + "app/**/*module*.js",
			client + "app/**/*.js"
		],
		dst: staging.scripts
	},

	styles: {
		vendors: client + "vendor/**/*.css",
		src: assets.images + "styles.scss",
		dst: staging.styles
	},

	images: {
		cache: assets.cache,
		src: assets.images + "*",
		dst: staging.images
	},

	svg: {
		src: assets.svg + "*.svg",
		dst: staging.images
	},

	fingerprint: {
		minFilter: "**/*.{css,js,jpg,png,svg}",
		index: "index.html",

		src: [
			stage + "**/*.{css,js,jpg,png,svg}",
			client + "index.html"
		],
		dst: dist
	},

	reference: {
		ext: [
			".html",
			".js"
		],
		src: dist + "**/*.{html,js}",
		dst: dist
	}

	// "templates": "./client/app/**/*.html",
	// "scripts": [
	// 	"./client/app/**/*module*.js",
	// 	"./client/app/**/*.js"
	// ],
	// "styles": [
	//     "./src/styles/styles.scss"	
	// ],
	// "images": [
	//     "./src/img/*"
	// ],
	// "svg": [
	//     "./src/svg/*"
	// ],

	// "vendorjs":	"./client/vendor/**/*.js",
	// "vendorcss": "./client/vendor/**/*.css",

	// "dist": "./dist"
}