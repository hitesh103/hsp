# HSP Enhancement Design Spec

## Overview

Enhance HSP (HTTP Superpowers) with 4 major capability improvements:
speed/productivity via variables and session memory, a proper TUI with warm color scheme, profiles/templates, and test suites.

**Date:** 2026-04-26
**Author:** hitesh103
**Status:** Draft

---

## 1. Variable System

### Syntax

Use `{{VARIABLE_NAME}}` syntax (Postman-style) throughout URL, headers, query params, and body.

Example:
```
URL: {{BASE_URL}}/users
Header: Authorization: Bearer {{API_KEY}}
Body: { "name": "{{TEST_USER}}", "email": "{{TEST_EMAIL}}" }
```

### Variable Storage

`~/.hsp/config.yaml` stores all variables organized by environment.

```yaml
environments:
  default:
    BASE_URL: "https://api.example.com"
    API_KEY: ""

  dev:
    BASE_URL: "http://localhost:3000"
    API_KEY: "dev-token-123"

  staging:
    BASE_URL: "https://staging-api.example.com"
    API_KEY: "staging-token-456"

  prod:
    BASE_URL: "https://api.example.com"
    API_KEY: "prod-token-789"
```

### Active Environment

- Default: `default` environment always active.
- Switching: `hsp env [name]` or `--env [name]` flag.
- Interactive prompt in `hsp request` to select environment.

### Variable Substitution

1. Load config from `~/.hsp/config.yaml`.
2. Load active environment variables.
3. Replace all `{{VAR}}` tokens in URL, headers, params, body.
4. Unknown variables: error with message listing missing variables.
5. Display resolved values (with masking for secrets like `***`).

### Phase 1 Feature: Variable Commands

```bash
hsp var list                          # List all variables for active env
hsp var set BASE_URL https://api.com  # Set a variable
hsp var delete API_KEY               # Delete a variable
hsp var export                       # Export config to stdout
hsp env                              # Show current env
hsp env dev                          # Switch to dev env
```

---

## 2. Session Memory

### Last Request Auto-fill

- Every request saves state to `~/.hsp/.last_request.json`.
- Fields saved: url, method, headers, params, body_format, body.
- On `hsp request`, if last URL starts with same base, auto-fill method/headers.
- User presses Enter to accept, or types new value to override.

### Session Resume

```bash
hsp request --resume        # Load last request and start from step 1
hsp request --last         # Jump to confirmation (just re-send last request)
```

### Named Profiles

Save and name requests as reusable profiles.

```bash
hsp profile save my-api       # Save current request as "my-api"
hsp profile list             # List all saved profiles
hsp profile run my-api       # Run profile "my-api" (with variable substitution)
hsp profile delete my-api    # Delete profile
hsp profile edit my-api      # Edit profile in $EDITOR
```

Profile stored in `~/.hsp/profiles/my-api.json`:

```json
{
  "name": "my-api",
  "url": "{{BASE_URL}}/users",
  "method": "POST",
  "headers": {
    "Authorization": "Bearer {{API_KEY}}",
    "Content-Type": "application/json"
  },
  "params": {},
  "body_format": "json",
  "body": "{ \"name\": \"{{TEST_USER}}\" }",
  "created_at": "2026-04-26T12:00:00Z"
}
```

---

## 3. TUI Output System

### Color Palette (Warm Scheme)

| Component     | Color              | Notes                                  |
|--------------|--------------------|----------------------------------------|
| Headers (important) | Bright Yellow (FgHiYellow)  | Most visible, light        |
| Headers (less)     | Dark Yellow (FgYellow)    | Standard yellow              |
| Body (important)   | Light Magenta (FgMagenta)  | Key fields, light          |
| Body (less)        | Deep Magenta (FgHiMagenta) | Less important, rich      |
| Response (important)| Bright Cyan (FgHiCyan)    | Key response fields        |
| Response (less)     | Dark Cyan (FgCyan)        | Standard cyan              |
| Error/4xx/5xx      | Bright Red (FgRed)         | Errors, critical          |
| Warning/3xx         | Bright Yellow (FgYellow)   | Redirects, warnings       |
| Success/2xx         | Bright Green (FgGreen)      | Success, positive          |
| Section labels      | Bold White (FgWhite+Bold)   | Headers like "REQUEST"    |
| Metadata           | Dim Gray (FgBlack)         | Timestamps, hints        |
| Borders            | Cyan (FgCyan)              | ASCII box borders          |

### Priority Shading Logic

- **Important keys:** `id`, `name`, `email`, `token`, `key`, `status`, `message`, `error`, `code`, `url`, `type`, `data`, `result`.
- **Less important:** Numeric IDs, timestamps, internal fields, pagination fields.
- Priority detected automatically by key name matching OR user can mark fields with `@important` comment in body.

### TUI Panel Output

Pure ASCII box drawing characters for maximum compatibility.

```
.------------------------------------------------------------------------------.
|  REQUEST                                                       [POST]        |
+------------------------------------------------------------------------------+
|  URL       : https://{{BASE_URL}}/users                                       |
|  Resolved  : https://api.example.com/users                                   |
+------------------------------------------------------------------------------+
|  Headers   : Authorization: Bearer ***                                        |
|             Content-Type: application/json                                     |
|             Accept: application/json                                          |
+------------------------------------------------------------------------------+
|  Params    : page=1 | limit=20                                               |
+------------------------------------------------------------------------------+
|  BODY                                                                      |
|  +--------------------------------------------------------------------------+
|  | name      : "John Doe"           [priority: high]                           |
|  | email     : "john@example.com"   [priority: high]                           |
|  | userId    : 1                 [priority: low]                             |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
|  Time: 143ms                              [201 Created]                        |
+------------------------------------------------------------------------------+
|  RESPONSE                                                                   |
|  +--------------------------------------------------------------------------+
|  | id        : 123                              [priority: high]              |
|  | name      : "John Doe"                        [priority: high]              |
|  | email     : "john@example.com"                  [priority: high]              |
|  | created_at: "2026-04-26T12:00:00Z"           [priority: low]               |
|  +--------------------------------------------------------------------------+
+------------------------------------------------------------------------------+
|  [PASS] Request saved: ~/.hsp/history/POST_2026-04-26_12-00-00.json           |
'------------------------------------------------------------------------------'
```

### Interactive TUI View (hsp request)

Step-by-step prompts with clean ASCII section separators between each step.

```
--- Step 1: URL ---
? URL: https://{{BASE_URL}}/users
    Resolved: https://api.example.com/users
    (resolved from config)

--- Step 2: Method ---
? Method: (default: GET)
    1) GET       [retrieve data]
    2) POST      [create new resource]
    3) PUT       [update entire resource]
    4) PATCH     [partial update]
    5) DELETE    [remove resource]
    6) HEAD      [like GET, no body]
    7) OPTIONS   [allowed methods]
Choose (1-7) or type method: 2

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
  "name": "{{TEST_USER}}",
  "email": "test@example.com"
}
    Valid JSON
    [Auto-set: Content-Type: application/json]

--- Step 6: Preview ---
```

Request preview with full resolved URL and color-coded sections.

```
.------------------------------------------------------------------------------.
|  PREVIEW                                                                    |
+------------------------------------------------------------------------------+
|  POST https://api.example.com/users                                          |
+------------------------------------------------------------------------------+
|  Headers   : Authorization: Bearer ***                                       |
|             Content-Type: application/json                                     |
|             Accept: application/json                                          |
+------------------------------------------------------------------------------+
|  Body      : name: "John Doe" [priority: high]                               |
|               email: "test@example.com" [priority: high]                      |
+------------------------------------------------------------------------------+
|  ? Send request? (y/n): y                                                   |
'------------------------------------------------------------------------------'
```

---

## 4. Command Aliases

Short aliases for power users.

```bash
# Phase 3 commands (aliases)
hsp r          # = hsp request
hsp r --last   # = hsp request --last
hsp r --resume # = hsp request --resume
hsp g          # = hsp get
hsp p          # = hsp post
hsp pu         # = hsp put
hsp pa         # = hsp patch
hsp d          # = hsp delete
hsp t          # = hsp test
hsp var        # = hsp var
hsp profile    # = hsp profile
hsp env        # = hsp env
```

Register aliases in `~/.hsp/config.yaml`:

```yaml
aliases:
  r: request
  g: get
  p: post
  pu: put
  pa: patch
  d: delete
  t: test
```

---

## 5. Test Suites (Phase 4)

### Test Suite Format

JSON-based (Postman-compatible), stored in `~/.hsp/suites/` or project root.

```json
{
  "name": "User API Test Suite",
  "version": "1.0",
  "description": "Complete user management API tests",
  "env": "dev",
  "variables": {
    "TEST_USER": "AutoTestUser",
    "TEST_EMAIL": "autotest@example.com"
  },
  "tests": [
    {
      "name": "Create user",
      "description": "Create a new user with valid data",
      "request": {
        "method": "POST",
        "url": "{{BASE_URL}}/users",
        "headers": {
          "Authorization": "Bearer {{API_KEY}}",
          "Content-Type": "application/json"
        },
        "body": {
          "name": "{{TEST_USER}}",
          "email": "{{TEST_EMAIL}}"
        }
      },
      "assertions": [
        { "type": "status", "expected": 201 },
        { "type": "body_contains", "path": "$.name", "value": "{{TEST_USER}}" },
        { "type": "header", "name": "Content-Type", "contains": "json" }
      ],
      "save": { "var": "USER_ID", "path": "$.id" }
    },
    {
      "name": "Get created user",
      "description": "Retrieve the user created above",
      "request": {
        "method": "GET",
        "url": "{{BASE_URL}}/users/{{USER_ID}}",
        "headers": {
          "Authorization": "Bearer {{API_KEY}}"
        }
      },
      "assertions": [
        { "type": "status", "expected": 200 },
        { "type": "body_contains", "path": "$.id", "value": "{{USER_ID}}" },
        { "type": "response_time_ms", "max": 500 }
      ]
    },
    {
      "name": "Update user",
      "description": "Update the user's name",
      "request": {
        "method": "PUT",
        "url": "{{BASE_URL}}/users/{{USER_ID}}",
        "headers": {
          "Authorization": "Bearer {{API_KEY}}",
          "Content-Type": "application/json"
        },
        "body": {
          "name": "{{TEST_USER}} Updated"
        }
      },
      "assertions": [
        { "type": "status", "expected": 200 }
      ]
    },
    {
      "name": "Delete user",
      "description": "Delete the test user",
      "request": {
        "method": "DELETE",
        "url": "{{BASE_URL}}/users/{{USER_ID}}",
        "headers": {
          "Authorization": "Bearer {{API_KEY}}"
        }
      },
      "assertions": [
        { "type": "status", "expected": 200 }
      ]
    }
  ]
}
```

### Assertion Types

| Type | Parameters | Description |
|------|------------|-------------|
| `status` | `expected: int` | Assert response status code |
| `body_contains` | `path: string`, `value: string` | JSON path contains value |
| `body_equals` | `path: string`, `value: string` | JSON path exactly equals value |
| `header` | `name: string`, `contains: string` | Header contains substring |
| `header_exists` | `name: string` | Header exists |
| `response_time_ms` | `max: int` | Response time under threshold |
| `body_size` | `max: int` | Response body size under threshold |
| `regex` | `path: string`, `pattern: string` | JSON path matches regex |

### Test Commands

```bash
hsp test run suite.json              # Run all tests in suite
hsp test run suite.json --env dev   # Run with specific environment
hsp test run suite.json --verbose   # Show all output
hsp test run suite.json --stop-on-fail  # Stop at first failure
hsp test list                       # List all test suites
hsp test create                     # Interactive test creator
hsp test export suite.json         # Export test suite template
```

### Test Output

Clean ASCII output showing pass/fail per test.

```
.------------------------------------------------------------------------------.
|  TEST SUITE: User API Test Suite                                    [dev]   |
+------------------------------------------------------------------------------+
|  Test: Create user                                       ................ OK |
|  Test: Get created user                                   .............. OK |
|  Test: Update user                                         .......... OK |
|  Test: Delete user                                         .......... OK |
+------------------------------------------------------------------------------+
|  Summary: 4/4 passed (143ms total)                              [PASS]      |
'------------------------------------------------------------------------------'
```

### Failed Test Output

```
.------------------------------------------------------------------------------.
|  TEST SUITE: User API Test Suite                                    [dev]   |
+------------------------------------------------------------------------------+
|  Test: Create user                                       ................ OK |
|  Test: Get created user                                   .............. OK |
|  Test: Update user                                         .......... FAIL |
|  Test: Delete user                                         .......... SKIP |
+------------------------------------------------------------------------------+
|  FAILURES:                                                                    |
|  [2] Update user                                                             |
|    Assertion: status                                                         |
|    Expected: 200                                                             |
|    Actual:   500 Internal Server Error                                        |
|    Body:     {"error": "Database connection failed"}                        |
+------------------------------------------------------------------------------+
|  Summary: 1/4 passed (143ms total)                              [FAIL]      |
'------------------------------------------------------------------------------'
```

---

## 6. Architecture

### File Structure

```
~/.hsp/
  config.yaml              # Variables, environments, aliases
  .last_request.json      # Session memory
  profiles/               # Named request profiles
    my-api.json
    github.json
  history/                # Request history (existing)
  suites/                 # Test suites
    user-api.json
    auth-tests.json
```

### Key Modules

| Module | File | Responsibility |
|--------|------|---------------|
| Config | `cmd/config.go` | Load/save `config.yaml`, environment management |
| Variables | `cmd/variables.go` | `{{VAR}}` substitution engine |
| Session | `cmd/session.go` | Last request memory, resume |
| Profiles | `cmd/profiles.go` | Named profile CRUD |
| TUI Output | `cmd/output.go` | ASCII panel rendering, color system |
| Priority | `cmd/priority.go` | Field importance detection |
| Test Runner | `cmd/test.go` | Test suite execution, assertions |

### Config Module API

```go
type Config struct {
    Environments map[string]map[string]string  // env name -> variable name -> value
    Aliases      map[string]string              // short alias -> full command
    ActiveEnv    string
}

func LoadConfig() (*Config, error)
func SaveConfig(*Config) error
func GetEnv(name string) map[string]string
func ResolveVariables(input string, env map[string]string) (string, []string)
```

### Variable Substitution Algorithm

```
1. Load config.yaml
2. Select active environment
3. For each `{{VAR}}` in input:
   a. Look up in environment map
   b. If found: replace with value
   c. If not found: add to missing list
4. If missing list non-empty: print error with list of missing vars, exit
5. Return resolved string
```

### Priority Detection Algorithm

```
1. Check if key matches high-priority pattern list:
   ["id", "name", "email", "token", "key", "status", "message", "error",
    "code", "url", "type", "data", "result", "user", "userId", "token",
    "access_token", "refresh_token", "secret", "password", "auth"]
2. Check if key matches low-priority pattern list:
   ["id", "uuid", "_id", "created_at", "updated_at", "timestamp", "v",
    "__v", "index", "seq", "offset", "cursor", "page", "per_page", "total"]
3. Default: medium priority
```

---

## 7. Implementation Phases

### Phase 1: Variables + Session (Core Speed)
- [ ] `cmd/config.go` — config loading/saving, `~/.hsp/config.yaml`
- [ ] `cmd/variables.go` — `{{VAR}}` substitution, environment groups
- [ ] `cmd/env.go` — `hsp env` command, env switching
- [ ] `cmd/var.go` — `hsp var` commands (list/set/delete)
- [ ] Session memory in `cmd/session.go`
- [ ] `--resume` and `--last` flags on `hsp request`
- [ ] Auto-fill from last request in prompts
- [ ] Variable highlighting in prompts (show `{{VAR}}` in distinct color)
- [ ] Error on missing variables with helpful message

### Phase 2: TUI + Colors
- [ ] `cmd/output.go` — ASCII panel rendering
- [ ] Warm color palette integration (fatih/color)
- [ ] Priority detection in `cmd/priority.go`
- [ ] Priority-based shading in output
- [ ] Updated response display with TUI panels
- [ ] Updated request preview with TUI panels
- [ ] Clean section separators between interactive steps
- [ ] `HSP_PRETTY` and `HSP_COLOR` env vars

### Phase 3: Profiles + Templates + Aliases
- [ ] `cmd/profiles.go` — profile CRUD
- [ ] `hsp profile save|list|run|delete|edit`
- [ ] Command aliases in `cmd/aliases.go`
- [ ] Short aliases registered via Cobra
- [ ] Profile template creation from `hsp request`
- [ ] Profile editing with `$EDITOR`

### Phase 4: Test Suites
- [ ] `cmd/test.go` — test runner core
- [ ] JSON test suite parsing
- [ ] All assertion types implemented
- [ ] `hsp test run|list|create|export`
- [ ] `save.var` — extract values from response for chaining
- [ ] Test suite execution with pass/fail output
- [ ] `hsp test --stop-on-fail` and `--verbose`

---

## 8. Open Questions

- [ ] Should variable values be masked in config file display? (Yes for keys with names like `token`, `secret`, `password`, `key`)
- [ ] Should `--last` show a preview or just resend immediately? (Preview first for safety)
- [ ] Should test suites auto-detect environment or require explicit `--env`? (Require explicit for clarity)
- [ ] Should profiles support variable slots that prompt user if not provided? (Yes, like Postman)