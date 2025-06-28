// tests/example.phantom.js
import { phantom } from 'phantom-vite';

describe('Website Tests', () => {
  test('should load homepage', async () => {
    const page = await phantom.goto('https://example.com');
    expect(await page.title()).toBe('Example Domain');
    await page.screenshot('homepage.png');
  });
});
