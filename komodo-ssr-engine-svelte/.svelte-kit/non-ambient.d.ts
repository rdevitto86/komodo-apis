
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/about" | "/api" | "/api/health" | "/api/v1" | "/api/v1/admin" | "/api/v1/admin/manage" | "/api/v1/admin/manage/content" | "/api/v1/admin/manage/content/invalidate" | "/api/v1/admin/manage/content/upsert" | "/api/v1/landing" | "/api/v1/marketing" | "/api/v1/marketing/content" | "/api/v1/marketing/content/[id]" | "/api/v1/marketing/user" | "/api/v1/marketing/user/[id]" | "/api/v1/orders" | "/api/v1/orders/[id]" | "/api/v1/products" | "/api/v1/products/[id]" | "/api/v1/services" | "/api/v1/services/scheduling" | "/api/v1/services/scheduling/[id]" | "/api/v1/services/[id]" | "/contact" | "/faq" | "/landing" | "/marketing" | "/marketing/content" | "/marketing/content/[id]" | "/marketing/user" | "/marketing/user/[id]" | "/orders" | "/orders/[id]" | "/products" | "/products/[id]" | "/services" | "/services/[id]" | "/terms";
		RouteParams(): {
			"/api/v1/marketing/content/[id]": { id: string };
			"/api/v1/marketing/user/[id]": { id: string };
			"/api/v1/orders/[id]": { id: string };
			"/api/v1/products/[id]": { id: string };
			"/api/v1/services/scheduling/[id]": { id: string };
			"/api/v1/services/[id]": { id: string };
			"/marketing/content/[id]": { id: string };
			"/marketing/user/[id]": { id: string };
			"/orders/[id]": { id: string };
			"/products/[id]": { id: string };
			"/services/[id]": { id: string }
		};
		LayoutParams(): {
			"/": { id?: string };
			"/about": Record<string, never>;
			"/api": { id?: string };
			"/api/health": Record<string, never>;
			"/api/v1": { id?: string };
			"/api/v1/admin": Record<string, never>;
			"/api/v1/admin/manage": Record<string, never>;
			"/api/v1/admin/manage/content": Record<string, never>;
			"/api/v1/admin/manage/content/invalidate": Record<string, never>;
			"/api/v1/admin/manage/content/upsert": Record<string, never>;
			"/api/v1/landing": Record<string, never>;
			"/api/v1/marketing": { id?: string };
			"/api/v1/marketing/content": { id?: string };
			"/api/v1/marketing/content/[id]": { id: string };
			"/api/v1/marketing/user": { id?: string };
			"/api/v1/marketing/user/[id]": { id: string };
			"/api/v1/orders": { id?: string };
			"/api/v1/orders/[id]": { id: string };
			"/api/v1/products": { id?: string };
			"/api/v1/products/[id]": { id: string };
			"/api/v1/services": { id?: string };
			"/api/v1/services/scheduling": { id?: string };
			"/api/v1/services/scheduling/[id]": { id: string };
			"/api/v1/services/[id]": { id: string };
			"/contact": Record<string, never>;
			"/faq": Record<string, never>;
			"/landing": Record<string, never>;
			"/marketing": { id?: string };
			"/marketing/content": { id?: string };
			"/marketing/content/[id]": { id: string };
			"/marketing/user": { id?: string };
			"/marketing/user/[id]": { id: string };
			"/orders": { id?: string };
			"/orders/[id]": { id: string };
			"/products": { id?: string };
			"/products/[id]": { id: string };
			"/services": { id?: string };
			"/services/[id]": { id: string };
			"/terms": Record<string, never>
		};
		Pathname(): "/" | "/about" | "/about/" | "/api" | "/api/" | "/api/health" | "/api/health/" | "/api/v1" | "/api/v1/" | "/api/v1/admin" | "/api/v1/admin/" | "/api/v1/admin/manage" | "/api/v1/admin/manage/" | "/api/v1/admin/manage/content" | "/api/v1/admin/manage/content/" | "/api/v1/admin/manage/content/invalidate" | "/api/v1/admin/manage/content/invalidate/" | "/api/v1/admin/manage/content/upsert" | "/api/v1/admin/manage/content/upsert/" | "/api/v1/landing" | "/api/v1/landing/" | "/api/v1/marketing" | "/api/v1/marketing/" | "/api/v1/marketing/content" | "/api/v1/marketing/content/" | `/api/v1/marketing/content/${string}` & {} | `/api/v1/marketing/content/${string}/` & {} | "/api/v1/marketing/user" | "/api/v1/marketing/user/" | `/api/v1/marketing/user/${string}` & {} | `/api/v1/marketing/user/${string}/` & {} | "/api/v1/orders" | "/api/v1/orders/" | `/api/v1/orders/${string}` & {} | `/api/v1/orders/${string}/` & {} | "/api/v1/products" | "/api/v1/products/" | `/api/v1/products/${string}` & {} | `/api/v1/products/${string}/` & {} | "/api/v1/services" | "/api/v1/services/" | "/api/v1/services/scheduling" | "/api/v1/services/scheduling/" | `/api/v1/services/scheduling/${string}` & {} | `/api/v1/services/scheduling/${string}/` & {} | `/api/v1/services/${string}` & {} | `/api/v1/services/${string}/` & {} | "/contact" | "/contact/" | "/faq" | "/faq/" | "/landing" | "/landing/" | "/marketing" | "/marketing/" | "/marketing/content" | "/marketing/content/" | `/marketing/content/${string}` & {} | `/marketing/content/${string}/` & {} | "/marketing/user" | "/marketing/user/" | `/marketing/user/${string}` & {} | `/marketing/user/${string}/` & {} | "/orders" | "/orders/" | `/orders/${string}` & {} | `/orders/${string}/` & {} | "/products" | "/products/" | `/products/${string}` & {} | `/products/${string}/` & {} | "/services" | "/services/" | `/services/${string}` & {} | `/services/${string}/` & {} | "/terms" | "/terms/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): "/favicon.png" | string & {};
	}
}