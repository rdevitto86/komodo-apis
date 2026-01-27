export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.png"]),
	mimeTypes: {".png":"image/png"},
	_: {
		client: {start:"_app/immutable/entry/start._RxN-VGp.js",app:"_app/immutable/entry/app.gAWLgfwJ.js",imports:["_app/immutable/entry/start._RxN-VGp.js","_app/immutable/chunks/LrPuS6WJ.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/Xs3Esnmx.js","_app/immutable/entry/app.gAWLgfwJ.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/CDPdhYCU.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/Xs3Esnmx.js","_app/immutable/chunks/oi9Nr6Ux.js","_app/immutable/chunks/CsoU4Bzd.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/2.js')),
			__memo(() => import('./nodes/3.js')),
			__memo(() => import('./nodes/4.js')),
			__memo(() => import('./nodes/5.js')),
			__memo(() => import('./nodes/6.js')),
			__memo(() => import('./nodes/7.js')),
			__memo(() => import('./nodes/8.js')),
			__memo(() => import('./nodes/9.js')),
			__memo(() => import('./nodes/10.js')),
			__memo(() => import('./nodes/11.js')),
			__memo(() => import('./nodes/12.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/about",
				pattern: /^\/about\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 3 },
				endpoint: null
			},
			{
				id: "/api/health",
				pattern: /^\/api\/health\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/health/_server.ts.js'))
			},
			{
				id: "/api/v1/admin/manage/content/invalidate",
				pattern: /^\/api\/v1\/admin\/manage\/content\/invalidate\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/admin/manage/content/invalidate/_server.ts.js'))
			},
			{
				id: "/api/v1/admin/manage/content/upsert",
				pattern: /^\/api\/v1\/admin\/manage\/content\/upsert\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/admin/manage/content/upsert/_server.ts.js'))
			},
			{
				id: "/api/v1/landing",
				pattern: /^\/api\/v1\/landing\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/landing/_server.ts.js'))
			},
			{
				id: "/api/v1/marketing/content/[id]",
				pattern: /^\/api\/v1\/marketing\/content\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/marketing/content/_id_/_server.ts.js'))
			},
			{
				id: "/api/v1/marketing/user/[id]",
				pattern: /^\/api\/v1\/marketing\/user\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/marketing/user/_id_/_server.ts.js'))
			},
			{
				id: "/api/v1/orders/[id]",
				pattern: /^\/api\/v1\/orders\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/orders/_id_/_server.ts.js'))
			},
			{
				id: "/api/v1/products/[id]",
				pattern: /^\/api\/v1\/products\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/products/_id_/_server.ts.js'))
			},
			{
				id: "/api/v1/services/scheduling/[id]",
				pattern: /^\/api\/v1\/services\/scheduling\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/services/scheduling/_id_/_server.ts.js'))
			},
			{
				id: "/api/v1/services/[id]",
				pattern: /^\/api\/v1\/services\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/v1/services/_id_/_server.ts.js'))
			},
			{
				id: "/contact",
				pattern: /^\/contact\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 4 },
				endpoint: null
			},
			{
				id: "/faq",
				pattern: /^\/faq\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 5 },
				endpoint: null
			},
			{
				id: "/landing",
				pattern: /^\/landing\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 6 },
				endpoint: null
			},
			{
				id: "/marketing/content/[id]",
				pattern: /^\/marketing\/content\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: { layouts: [0,], errors: [1,], leaf: 7 },
				endpoint: null
			},
			{
				id: "/marketing/user/[id]",
				pattern: /^\/marketing\/user\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: { layouts: [0,], errors: [1,], leaf: 8 },
				endpoint: null
			},
			{
				id: "/orders/[id]",
				pattern: /^\/orders\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: { layouts: [0,], errors: [1,], leaf: 9 },
				endpoint: null
			},
			{
				id: "/products/[id]",
				pattern: /^\/products\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: { layouts: [0,], errors: [1,], leaf: 10 },
				endpoint: null
			},
			{
				id: "/services/[id]",
				pattern: /^\/services\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: { layouts: [0,], errors: [1,], leaf: 11 },
				endpoint: null
			},
			{
				id: "/terms",
				pattern: /^\/terms\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 12 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
