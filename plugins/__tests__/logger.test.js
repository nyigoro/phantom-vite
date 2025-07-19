import { onStart } from '../logger';

describe('logger plugin', () => {
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

    expect(consoleSpy).toHaveBeenCalledWith('[logger Plugin] onStart triggered by command: open');
    expect(consoleSpy).toHaveBeenCalledWith('[logger Plugin] Engine: puppeteer');
    expect(consoleSpy).toHaveBeenCalledWith('[logger Plugin] Target URL: https://example.com');
  });

  test('should warn when no context is provided', () => {
    const consoleWarnSpy = jest.spyOn(console, 'warn').mockImplementation(() => {});
    onStart(null);
    expect(consoleWarnSpy).toHaveBeenCalledWith('[logger Plugin] No context');
    consoleWarnSpy.mockRestore();
  });

  test('should use default values if meta is missing', () => {
    const mockContext = {
      engine: 'playwright',
      meta: {},
    };

    onStart(mockContext);

    expect(consoleSpy).toHaveBeenCalledWith('[logger Plugin] onStart triggered by command: unknown');
    expect(consoleSpy).toHaveBeenCalledWith('[logger Plugin] Engine: playwright');
    expect(consoleSpy).toHaveBeenCalledWith('[logger Plugin] Target URL: N/A');
  });
});
