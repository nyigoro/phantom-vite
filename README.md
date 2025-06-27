# 🕴️ Phantom Vite

A hybrid headless browser CLI inspired by PhantomJS — powered by Go, Node.js (Puppeteer & Vite), and Python. Phantom Vite brings together fast builds, headless browsing, and agent-based automation — right from your terminal.

---

## 🚀 Features

- ✅ `open <url>` — Headless Puppeteer screenshot + title
- ✅ `build` — Vite build pipeline for frontend assets
- ✅ `serve <file>` — Vite preview mode for local development
- ✅ `agent <prompt>` — Python-based AI agent handler
- ✅ `gemini <prompt>` — Google Gemini CLI integration
- ✅ `<script.js>` — Run any Node.js script directly
- ✅ Cross-platform (Linux & Windows)
- ✅ Optional config via `phantomvite.config.json`
- ✅ Clear CLI help and Puppeteer dependency checks

---

## 📦 Installation

```bash

### Prerequisites
- Go 1.21+
- Node.js 20+
- Python 3.9+
- [Gemini CLI](https://github.com/google-gemini/gemini-cli) (optional)

```bash
npm install -g @google/gemini-cli
````

### Build from source

```bash
git clone https://github.com/yourname/phantom-vite
cd phantom-vite
npm install
npx vite build
go build -o phantom-vite ./cmd
```

---

## 🛠️ Usage

```bash
phantom-vite open https://example.com
phantom-vite build
phantom-vite agent "summarize this repo"
phantom-vite gemini "generate a blog post on Go concurrency"
phantom-vite serve dist/index.html
phantom-vite myscript.js
```

---

## 🧠 Config (Optional)

```json
{
  "headless": true,
  "viewport": {
    "width": 1280,
    "height": 720
  }
}
```

---

## 🧪 Tests

```bash
go test ./cmd
```

---

## 📂 Project Structure

```
phantom-vite/
├── cmd/                  # Go CLI
│   ├── main.go
│   └── main_test.go
├── python/               # Python agent
│   └── agent.py
├── dist/                 # Vite build output
├── phantomvite.config.json
├── go.mod / go.sum
├── package.json / vite.config.js
├── README.md
└── .github/workflows/
    └── build.yml
```

---

## 📄 License

MIT — PhantomJS-inspired, Gemini-powered, developer-first.

```bash
nyigoro
```

---
