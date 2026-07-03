import { expect, test } from '@playwright/test';
import { requireScenario } from './support/gherkin';

test('SCN-AUTH-001 anonymous visitor sees the login experience', async ({ page, request }) => {
	requireScenario('SCN-AUTH-001', 'e2e');

	const session = await request.get('/api/auth/session');
	await expect(session).toBeOK();
	await expect(await session.json()).toMatchObject({ authenticated: false });

	await page.goto('/');
	await expect(page.getByRole('heading', { name: 'Admin login' })).toBeVisible();
	await expect(page.getByLabel('Username')).toHaveValue('admin');
	await expect(page.getByLabel('Password')).toHaveValue('admin');
});

test('SCN-AUTH-002 admin signs in with valid credentials', async ({ page }) => {
	requireScenario('SCN-AUTH-002', 'e2e');

	await page.goto('/');
	await page.getByRole('button', { name: 'Log in' }).click();

	await expect(page.getByRole('navigation')).toBeVisible();
	await expect(page.getByRole('link', { name: 'Settings' })).toBeVisible();
	await expect(page.getByRole('link', { name: 'System' })).toBeVisible();
});
