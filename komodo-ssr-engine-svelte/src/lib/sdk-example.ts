import type { Product, Service, MarketingContent } from '@komodo-forge-sdk/typescript/types';

export function exampleProduct(): Product {
	return {
		id: '1',
		slug: 'example-product',
		name: 'Example Product',
		description: 'This is an example product',
		status: 'active',
		trackInventory: false,
		variants: [
			{
				id: 'v1',
				name: 'Default Variant',
				price: 99.99
			}
		]
	};
}

export function exampleService(): Service {
	return {
		id: 's1',
		slug: 'example-service',
		name: 'Example Service',
		description: 'This is an example service',
		category: 'installation',
		status: 'active',
		price: 199.99,
		locationTypes: ['residential']
	};
}

export function exampleMarketingContent(): MarketingContent {
	return {
		id: 'm1',
		code: 'HERO_BANNER_001',
		name: 'Hero Banner',
		type: 'hero-banner',
		status: 'active',
		format: 'structured',
		priority: 1,
		content: {
			title: 'Welcome to Komodo',
			subtitle: 'Your trusted partner',
			ctaText: 'Get Started',
			ctaUrl: '/contact'
		}
	};
}
