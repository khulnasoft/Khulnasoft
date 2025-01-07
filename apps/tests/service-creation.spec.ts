import { expect, PlaywrightTestArgs, test } from '@playwright/test';

import { authenticate, catchNewPage, deleteKhulnasoftResources } from './test-utils';

async function cleanup({ page }: PlaywrightTestArgs) {
  test.setTimeout(5 * 60 * 1000);
  await authenticate(page);
  await deleteKhulnasoftResources(page);
}

test.beforeEach(cleanup);
test.afterEach(cleanup);

test('github service creation', async ({ context, page }) => {
  test.setTimeout(15 * 60 * 1000);
  await authenticate(page);

  await page.getByRole('link', { name: 'Create service' }).click();
  await page.getByRole('button', { name: 'GitHub' }).click();
  await page.getByPlaceholder('https://github.com/khulnasoft/example-go').fill('khulnasoft/khulnasoft/tree/main/examples/go-gin');
  await page.getByRole('button', { name: 'Import' }).click();
  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('button', { name: 'Deploy' }).click();

  await expect(page.getByText('Your service is ready')).toBeVisible({ timeout: 15 * 60 * 1_000 });

  const servicePage = await catchNewPage(context, async () => {
    await page.getByRole('link', { name: 'Visit your service public domain' }).click();
  });

  await expect(async () => {
    await servicePage.reload();
    await expect(servicePage.getByText('Hello, world!')).toBeVisible();
  }).toPass({ timeout: 60 * 1000 });
});

test('docker service creation', async ({ context, page }) => {
  test.setTimeout(15 * 60 * 1000);

  await page.goto('/');

  await page.getByRole('link', { name: 'Create service' }).click();
  await page.getByRole('button', { name: 'Docker' }).click();
  await page.getByPlaceholder('docker.io/khulnasoft/demo:latest').fill('khulnasoft/demo');
  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('button', { name: 'Next' }).click();
  await page.getByRole('button', { name: 'Deploy' }).click();

  await expect(page.getByText('Your service is ready')).toBeVisible({ timeout: 15 * 60 * 1_000 });

  const servicePage = await catchNewPage(context, async () => {
    await page.getByRole('link', { name: 'Visit your service public domain' }).click();
  });

  await expect(async () => {
    await servicePage.reload();
    await expect(servicePage.getByText('Welcome to Khulnasoft')).toBeVisible();
  }).toPass({ timeout: 60 * 1000 });
});
