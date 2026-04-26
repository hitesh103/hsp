# HSP - HTTP Superpowers

**The ultimate HTTP client in the terminal** - Postman-like experience with variables, profiles, test suites, and beautiful TUI output.

---

## What is HSP?

HSP is an interactive CLI HTTP client that makes API testing incredibly fast. No complex flags, no GUI overhead - just run commands and see beautiful output.

```
$ hsp request
+------------------------------------------------------------------------------+
|  REQUEST                                                           [POST]  |
+------------------------------------------------------------------------------+
|  URL: https://api.example.com/users                                       |
|  Headers: Authorization: Bearer ***                                      |
+------------------------------------------------------------------------------+
|  BODY (Payload)                                                       |
|  +--------------------------------------------------------------------------+
|  | name     : "John Doe"                                                |
|  | email    : "john@example.com"                                       |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
|  Time: 143ms                              [201 Created]                    |
+------------------------------------------------------------------------------+
|  RESPONSE                                                            |
|  +--------------------------------------------------------------------------+
|  | id       : 123                                                    |
|  | name     : "John Doe"                                             |
|  | email    : "john@example.com"                                     |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
```

---

## Quick Comparison

### vs. cURL
```bash
# cURL - Remember all the flags
curl -X POST https://api.example.com/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token123" \
  -d '{"name": "John", "email": "john@example.com"}'

# HSP - Just run and answer prompts
$ hsp request
? URL: https://api.example.com/users
? Method: POST
...
```
**Result:** Same API call, but zero memorization needed.

### vs. Postman
- **500MB+ download** vs **15MB binary**
- **Mouse required** vs **100% keyboard**
- **Slow launch** vs **instant start**
- **No history by default** vs **auto-saved**

---

## Installation

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

### Verify
```bash
hsp --help
```

---

## All Commands

```
Available Commands:
  completion   Generate autocompletion script
  env          Manage environments (dev/staging/prod)
  get, g       Send GET request
  post, p      Send POST request
  profile      Manage request profiles
  request, r   Interactive request builder
  test         Run API test suites
  var          Manage variables
```

---

## Features

### 1. Variables & Environments

Set once, use everywhere with `{{VAR}}` syntax.

```bash
# Set variables
$ hsp var set BASE_URL https://api.example.com
Set BASE_URL = https://api.example.com in environment 'default'

$ hsp var set API_KEY secret123 --env dev
Set API_KEY = secret123 in environment 'dev'

# Create environments
$ hsp env create staging
Created environment 'staging'

$ hsp env staging
Switched to 'staging'

# List environments
$ hsp env --list
* default
  dev
  staging
  prod (current)

# Use in requests
$ hsp request
? URL: {{BASE_URL}}/users
    Resolved: https://api.example.com/users
```

**Output:**
```
+------------------------------------------------------------------------------+
|  REQUEST                                                           [GET]  |
+------------------------------------------------------------------------------+
|  URL       : {{BASE_URL}}/users                                       |
|  Resolved  : https://api.example.com/users                           |
+------------------------------------------------------------------------------+
```

### 2. Session Memory

Never lose your last request.

```bash
# Last request auto-saved
$ hsp request --last          # Just re-send last request
$ hsp request --resume        # Load and modify before sending
```

**Output (--last):**
```
+------------------------------------------------------------------------------+
|  GET https://api.example.com/users                                  |
+------------------------------------------------------------------------------+
|  Headers: Accept: application/json                               |
+------------------------------------------------------------------------------+
|  Time: 89ms                                   [200 OK]     |
+------------------------------------------------------------------------------+
```

### 3. Named Profiles

Save and reuse request templates.

```bash
# Save current request
$ hsp profile save github-api
Saved profile 'github-api'

# List profiles
$ hsp profile list
  github-api  - GET to https://api.github.com (updated: 2026-04-27)
  my-api    - POST to {{BASE_URL}}/users

# Run a profile
$ hsp profile run github-api
+------------------------------------------------------------------------------+
|  TEST SUITE: github-api                                     [dev]  |
+------------------------------------------------------------------------------+
```

### 4. Command Aliases

Shortcuts for power users.

```bash
hsp r          # = hsp request
hsp g          # = hsp get
hsp p          # = hsp post
hsp d          # = hsp delete
```

### 5. Test Suites

JSON-based test runner with assertions.

**Example test suite:**
```json
{
  "name": "User API Tests",
  "tests": [
    {
      "name": "Create user",
      "request": {
        "method": "POST",
        "url": "{{BASE_URL}}/users",
        "body": { "name": "test" }
      },
      "assertions": [
        { "type": "status", "expected": 201 },
        { "type": "body_contains", "path": "$.name", "value": "test" }
      ]
    }
  ]
}
```

**Run tests:**
```bash
$ hsp test run user-api.json
+------------------------------------------------------------------------------+
|  TEST SUITE: User API Tests                                      [dev]  |
+------------------------------------------------------------------------------+
|  PASS  Create user ....................................... OK   45ms |
|  PASS  Get user ........................................ OK   12ms |
|  PASS  Update user .................................... OK   89ms |
|  FAIL  Delete user ................................. FAIL             |
+------------------------------------------------------------------------------+
|  FAILURES:                                                            |
|  [4] Delete user                                                     |
|    Expected: 200                                                   |
|    Actual:   500 Internal Server Error                              |
+------------------------------------------------------------------------------+
|  Summary: 3/4 passed                                         [FAIL]  |
+------------------------------------------------------------------------------+
```

### 6. Beautiful TUI Output

Warm color scheme with priority highlighting.

```
+------------------------------------------------------------------------------+
|  REQUEST                                                           [POST]  |
+------------------------------------------------------------------------------+
|  URL       : https://api.example.com/users                              |
|  Headers  : Authorization: Bearer ***                               |
|             Content-Type: application/json                       |
+------------------------------------------------------------------------------+
|  BODY (Payload)                                                   |
|  +--------------------------------------------------------------------------+
|  | name     : "John Doe"        [priority: high]                  |
|  | email    : "john@example.com" [priority: high]             |
|  | userId  : 1               [priority: low]                   |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
|  Time: 143ms                              [201 Created]               |
+------------------------------------------------------------------------------+
|  RESPONSE                                                          |
|  +--------------------------------------------------------------------------+
|  | id       : 123                 [priority: high]           |
|  | name     : "John Doe"                                 |
|  | created_at: "2026-04-27T00:00:00Z" [priority: low]  |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
```

### 7. Auto History

Every request saved automatically.

```bash
$ ls ~/.hsp/history/
GET_2026-04-27_10-30-00.json
POST_2026-04-27_10-35-15.json
```

---

## Usage Examples

### Example 1: Interactive Request
```bash
$ hsp request

--- Step 1: URL ---
? URL: {{BASE_URL}}/users
    Resolved: https://api.example.com/users

--- Step 2: Method ---
? Method: (default: GET)
    1) GET       [retrieve data]
    2) POST      [create new]
    3) PUT       [update]
    4) PATCH     [partial]
    5) DELETE    [remove]
Choose (1-7) or type method: POST

--- Step 3: Headers ---
? Add headers? (y/n): n
    [Auto-set: Accept: application/json]

--- Step 4: Query Parameters ---
? Add query parameters? (y/n): n

--- Step 5: Body ---
? Add request body? (y/n): y
    Body format:
    1) JSON
    2) Form data
    3) Raw text
Choose (1-3): 1
Enter JSON body (press Enter twice when done):
{
  "name": "John Doe",
  "email": "john@example.com"
}
    Valid JSON
    [Auto-set: Content-Type: application/json]

--- Step 6: Preview ---
+------------------------------------------------------------------------------+
|  PREVIEW                                                            |
+------------------------------------------------------------------------------+
|  POST https://api.example.com/users                               |
+------------------------------------------------------------------------------+
|  Headers: Accept: application/json                               |
|            Content-Type: application/json                        |
+------------------------------------------------------------------------------+
|  Body: name: "John Doe"                                         |
|        email: "john@example.com"                               |
+------------------------------------------------------------------------------+
? Send request? (y/n): y

+------------------------------------------------------------------------------+
|  Time: 286ms                              [201 Created]        |
+------------------------------------------------------------------------------+
|  RESPONSE                                                          |
|  +--------------------------------------------------------------------------+
|  | id        : 101                                             |
|  | name      : "John Doe"                                      |
|  | email     : "john@example.com"                             |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
Request saved: ~/.hsp/history/POST_2026-04-27_10-40-00.json
```

### Example 2: Quick GET
```bash
$ hsp g https://jsonplaceholder.typicode.com/posts/1

Status: 200 (48ms)

{
  "body": "quia et suscipit...",
  "id": 1,
  "title": "sunt aut facere...",
  "userId": 1
}
```

### Example 3: Variables + Profile
```bash
# Set up variables
$ hsp env create prod
$ hsp var set BASE_URL https://api.production.com
$ hsp var set API_KEY prod-secret-token

# Save a profile
$ hsp profile save prod-api
Saved profile 'prod-api'

# Run later
$ hsp profile run prod-api
+------------------------------------------------------------------------------+
|  GET https://api.production.com/users              |
+------------------------------------------------------------------------------+
```

---

## Configuration

### Environment Variables
```bash
export HSP_HISTORY_DIR="$HOME/Documents/api-requests"
export HSP_PRETTY=true
export HSP_TIMEOUT=30s
export HSP_COLOR=true    # Enable/disable colors
```

### Config File
`~/.hsp/config.yaml` stores all variables and environments:
```yaml
environments:
  default:
    BASE_URL: ""
    API_KEY: ""
  dev:
    BASE_URL: "http://localhost:3000"
    API_KEY: "dev-token"
  prod:
    BASE_URL: "https://api.example.com"
    API_KEY: "prod-token"
```

---

## File Locations

| Path | Purpose |
|------|---------|
| `~/.hsp/config.yaml` | Variables & environments |
| `~/.hsp/history/` | Request history |
| `~/.hsp/profiles/` | Saved profiles |
| `~/.hsp/suites/` | Test suites |
| `~/.hsp/.last_request.json` | Session memory |

---

## Troubleshooting

### "URL must start with http:// or https://"
Check your URL: needs `http://` or `https://` prefix.

### "Invalid JSON"
Your JSON body has a syntax error. Check brackets and quotes.

### "Variable not found: {{VAR}}"
Set the variable first: `hsp var set VAR_NAME value`

### "Environment not found"
Create it: `hsp env create myenv`

---

## License

MIT License

---

**Built for developers who love the terminal.**