from selenium import webdriver
from selenium.webdriver.chrome.options import Options

options = Options()
options.add_argument("--headless")
driver = webdriver.Chrome(options=options)

driver.get("https://example.com")
print("[Selenium] Page title:", driver.title)
driver.quit()
