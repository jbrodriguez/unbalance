var client = "./client/";
var server = "./server/";
var dist = "./dist/";

//var stage = "./staging/";

var sources = {
		styles: "./src/styles/",
		images: "./src/images/",
		scripts: "./src/scripts/",
		svg: "./src/svg/",
		cache: "./src/cache/",
		tools: "./src/tools/"
};

var	staging = {
		root: "./staging/",
		styles: "./staging/css/",
		images: "./staging/img/",
		scripts: "./staging/js/",
};

module.exports = {
	clean: {
		staging: staging.root,
		dist: dist
	},

	tools: {
		src: sources.tools + '*',
		dst: dist
	},

	build: {
		server: server,
		dist: dist
	},

	templates: {
		src: client + "app/**/*.html",
		dst: sources.cache
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
		src: sources.styles + "styles.scss",
		dst: staging.styles
	},

	images: {
		cache: sources.cache,
		src: sources.images + "*",
		dst: staging.images
	},

	svg: {
		src: sources.svg + "*.svg",
		dst: staging.images
	},

	fingerprint: {
		revFilter: "**/*.{css,js,jpg,png,svg}",
		index: "index.html",

		src: [
			staging.root + "**/*.{css,js,jpg,png,svg}",
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
	},

	publish: {
		src: dist,
		dst: "/boot/custom/unbalance"
	}
}