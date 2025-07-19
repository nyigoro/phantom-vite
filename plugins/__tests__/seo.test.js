import { onStart } from '../seo';

describe('SEO plugin', () => {
  let consoleSpy;

  beforeEach(() => {
    consoleSpy = jest.spyOn(console, 'log').mockImplementation(() => {});
    jest.spyOn(console, 'warn').mockImplementation(() => {});
  });

  afterEach(() => {
    consoleSpy.mockRestore();
  });

  test('should log correct messages when context is provided', () => {
    const mockContext = {
      engine: 'puppeteer',
      meta: {
        command: 'open',
        url: 'https://example.com',
      },
    };

    onStart(mockContext);

    expect(consoleSpy).toHaveBeenCalledWith('[SEO Plugin] onStart triggered by command: open');
    expect(consoleSpy).toHaveBeenCalledWith('[SEO Plugin] Engine: puppeteer');
    expect(consoleSpy).toHaveBeenCalledWith('[SEO Plugin] Target URL: https://example.com');
  });

  test('should warn when no context is provided', () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation(() => {});
    onStart(null);
    expect(consoleWarnSpy).toHaveBeenCalledWith('[SEO Plugin] No context');
    consoleWarnSpy.mockRestore();
  });

  test('should use default values if meta is missing', () => {
    const mockContext = {
      engine: 'playwright',
      meta: {},
    };

    onStart(mockContext);

    expect(consoleSpy).toHaveBeenCalledWith('[SEO Plugin] onStart triggered by command: unknown');
    expect(consoleSpy).toHaveBeenCalledWith('[SEO Plugin] Engine: playwright');
    expect(consoleSpy).toHaveBeenCalledWith('[SEO Plugin] Target URL: N/A');
  });
});
