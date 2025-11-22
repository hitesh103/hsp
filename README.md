# HSP - HTTP Superpowers 🚀

**The easiest HTTP client in the terminal - Postman-like experience without the complexity**

HSP is an interactive CLI tool that makes HTTP requests as simple as answering prompts. No need to remember curl flags or compose complex commands!

## ✨ Why HSP?

### vs. cURL
```bash
# cURL - Need to remember flags and syntax
curl -X POST https://api.example.com/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token123" \
  -d '{"name": "John", "email": "john@example.com"}'

# HSP - Just answer prompts
hsp request
# ? URL: https://api.example.com/users
# ? Method: (GET, POST, PUT, PATCH, DELETE) POST
# ? Add headers? y
#   Header name: Authorization
#   Value: Bearer token123
#   Add another? y
#   Header name: Custom-Header
#   Value: custom-value
# ... (clear preview and confirmation)
```

### vs. Postman
- **No UI overhead** - pure terminal speed
- **Lightweight** - single binary, ~15MB
- **Scriptable** - pipe input for automation
- **History** - all requests saved automatically
- **Keyboard-driven** - never reach for mouse

## 🎯 Features

- 🎨 **Interactive Request Builder** - Step-by-step guided workflow
- 🔄 **Auto-formatting** - JSON body formatting & Content-Type auto-detection
- 💾 **Request History** - All requests stored in `~/.hsp/history/`
- 📋 **Request Preview** - See exactly what will be sent before confirming
- 🌈 **Colored Output** - Beautiful response display with syntax highlighting
- ⚡ **Quick Commands** - `hsp get <url>`, `hsp post <url>` for fast requests
- ✅ **Input Validation** - Prevents malformed URLs and invalid JSON

## 📦 Installation

### From Source
```bash
git clone https://github.com/hitesh103/hsp.git
cd hsp
go build -o hsp
sudo mv hsp /usr/local/bin/
```

### macOS (Homebrew)
```bash
brew install hitesh103/hsp/hsp
```

### Verify Installation
```bash
hsp --version
hsp --help
```

## 🚀 Quick Start

### Interactive Mode (Recommended)
```bash
hsp request
```

You'll be guided through:
1. **URL** - Where to send the request
2. **Method** - GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
3. **Headers** - Add custom headers easily
4. **Query Parameters** - Key-value pairs appended to URL
5. **Body** (if POST/PUT/PATCH) - JSON, form data, or raw text
6. **Pretty-print** - Format response nicely
7. **Preview** - Review before sending
8. **Confirmation** - Send or cancel

### Quick GET
```bash
hsp get https://api.github.com/users/golang
```

### Quick POST
```bash
hsp post https://api.example.com/data --json '{"key": "value"}'
```

## 📖 Usage Examples

### Example 1: GET Request with Headers
```bash
$ hsp request
? URL: https://api.github.com/repos/golang/go/issues
? Method: GET
? Add headers? (y/n): y
  Header name: Authorization
  Value: token ghp_XXXXXXXXXXXX
  Add another? (y/n): n
? Add query parameters? (y/n): y
  Parameter name: state
  Value: open
  Add another? (y/n): y
  Parameter name: labels
  Value: bug
  Add another? (y/n): n
? Pretty response? (y/n, default: y): y

PREVIEW
======================================================================
GET https://api.github.com/repos/golang/go/issues?state=open&labels=bug
Headers:
  Accept: application/json
  Authorization: token ghp_XXXXXXXXXXXX
======================================================================

? Send request? (y/n): y

✔ 200 OK (143ms)

Response Headers:
  Content-Type: application/json; charset=utf-8
  X-RateLimit-Limit: 60
  X-RateLimit-Remaining: 59

Response Body:
[
  {
    "id": 12345,
    "title": "Example issue",
    "state": "open",
    ...
  }
]
✓ Request saved to history: /Users/dev/.hsp/history/GET_2025-11-22_12-49-07.json
```

### Example 2: POST with JSON Body
```bash
$ hsp request
? URL: https://jsonplaceholder.typicode.com/posts
? Method: POST
? Add headers? (y/n): n
? Add query parameters? (y/n): n
? Add request body? (y/n): y
  Body format:
  1) JSON
  2) Form data
  3) Raw text
  Choose (1-3): 1
  Enter JSON body (press Enter twice when done):
  {
    "title": "HSP is awesome",
    "body": "Making HTTP requests easy",
    "userId": 1
  }
  
  ✓ JSON body set
? Pretty response? (y/n, default: y): y

PREVIEW
======================================================================
POST https://jsonplaceholder.typicode.com/posts
Headers:
  Accept: application/json
  Content-Type: application/json

Body:
  {
    "body": "Making HTTP requests easy",
    "title": "HSP is awesome",
    "userId": 1
  }
======================================================================

? Send request? (y/n): y

✔ 201 Created (286ms)

Response Headers:
  Content-Type: application/json; charset=utf-8
  Location: https://jsonplaceholder.typicode.com/posts/101

Response Body:
{
  "body": "Making HTTP requests easy",
  "id": 101,
  "title": "HSP is awesome",
  "userId": 1
}
✓ Request saved to history: /Users/dev/.hsp/history/POST_2025-11-22_12-49-49.json
```

### Example 3: PUT Request with Form Data
```bash
$ hsp request
? URL: https://api.example.com/users/123
? Method: PUT
? Add headers? (y/n): y
  Header name: Authorization
  Value: Bearer eyJhbGc...
  Add another? (y/n): n
? Add query parameters? (y/n): n
? Add request body? (y/n): y
  Body format:
  1) JSON
  2) Form data
  3) Raw text
  Choose (1-3): 2
  Form field name: first_name
  Value: John
  Add another? (y/n): y
  Form field name: last_name
  Value: Doe
  Add another? (y/n): n
  ✓ Form body set
? Pretty response? (y/n, default: y): y

PREVIEW
======================================================================
PUT https://api.example.com/users/123
Headers:
  Accept: application/json
  Authorization: Bearer eyJhbGc...
  Content-Type: application/x-www-form-urlencoded

Body:
  first_name=John&last_name=Doe
======================================================================

? Send request? (y/n): y

✔ 200 OK (95ms)
```

## 📚 Request History

All requests are automatically saved to `~/.hsp/history/` with timestamps:

```bash
ls -la ~/.hsp/history/
# GET_2025-11-22_12-49-07.json
# POST_2025-11-22_12-49-49.json
# PUT_2025-11-22_13-10-15.json

cat ~/.hsp/history/POST_2025-11-22_12-49-49.json
# {
#   "timestamp": "2025-11-22_12-49-49",
#   "method": "POST",
#   "url": "https://jsonplaceholder.typicode.com/posts",
#   "headers": {
#     "Accept": "application/json",
#     "Content-Type": "application/json"
#   },
#   "params": {},
#   "body": "{\"title\": \"HSP Test Post\", ...}"
# }
```

## 🎨 Features in Detail

### Auto-Header Management
- **Auto-set Content-Type** for JSON bodies
- **Auto-set Accept** header to `application/json`
- **Easy multiple headers** - add as many as needed
- **Common header templates** (Authorization, X-API-Key, etc.)

### JSON Body Auto-Formatting
- **Validates JSON** before sending
- **Pretty-prints** in preview
- **Auto-detects** objects vs arrays
- **Handles** Unicode and special characters

### Query Parameters
- **Key-value pairs** with interactive prompts
- **URL encoding** handled automatically
- **Multiple parameters** supported
- **Easy to modify** before sending

### Response Display
- **Color-coded status** (green=2xx, red=4xx/5xx)
- **Response headers** displayed
- **JSON pretty-printing** with colors
- **Timing information** for performance analysis

## ⌨️ Keyboard Shortcuts

| Command | Description |
|---------|-------------|
| `Ctrl+C` | Cancel current operation |
| `n` | Skip optional steps |
| `y` | Confirm and proceed |
| `done` | Finish adding headers/params |

## 🔧 Configuration

### Custom History Location
Set environment variable:
```bash
export HSP_HISTORY_DIR="$HOME/Documents/api-requests"
```

### Default Pretty-Print
```bash
export HSP_PRETTY=true
```

### Timeout
```bash
export HSP_TIMEOUT=30s
```

## 📊 Supported HTTP Methods

- ✅ GET
- ✅ POST
- ✅ PUT
- ✅ PATCH
- ✅ DELETE
- ✅ HEAD
- ✅ OPTIONS

## 🐛 Troubleshooting

### "URL required" error
Make sure you enter a URL starting with `http://` or `https://`

### "Invalid JSON" error
Check your JSON syntax. HSP validates before sending.

### Connection timeout
- Check your internet connection
- Increase timeout: `export HSP_TIMEOUT=60s`
- Verify the URL is correct

### Request not saved
History is auto-saved to `~/.hsp/history/`. Check permissions with:
```bash
ls -la ~/.hsp/history/
```

## 🤝 Contributing

Found a bug or have a feature request? Open an issue!

## 📄 License

MIT License - See LICENSE file

---

**Made with ❤️ for developers who love the terminal**

Questions? Create an issue or check documentation at https://github.com/hitesh103/hsp
