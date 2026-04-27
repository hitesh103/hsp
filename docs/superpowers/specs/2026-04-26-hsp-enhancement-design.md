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

Every output element has a defined color. Colors use `fatih/color` package attributes.

| Element | Color | Hex / ANSI | Usage |
|---------|-------|------------|-------|
| **Box Borders** | Cyan | `color.FgCyan` | All ASCII box frame lines (`.`, `-`, `|`, `+`) |
| **Section Title** | Bold White | `color.FgWhite + color.Bold` | "REQUEST", "PREVIEW", "RESPONSE", "BODY", "HEADERS" labels |
| **HTTP Method** | Bold Cyan | `color.FgCyan + color.Bold` | GET, POST, PUT, PATCH, DELETE badges |
| **Status 2xx** | Bold Green | `color.FgGreen + color.Bold` | 200, 201, 202, 204 |
| **Status 3xx** | Bold Yellow | `color.FgYellow + color.Bold` | 301, 302, 304 |
| **Status 4xx** | Bold Red | `color.FgRed + color.Bold` | 400, 401, 403, 404, 422 |
| **Status 5xx** | Bold Magenta | `color.FgMagenta + color.Bold` | 500, 502, 503, 504 |
| **Status Text** | Same as status code | | "OK", "Created", "Bad Request", etc. |
| **Error / FAIL** | Bold Red | `color.FgRed + color.Bold` | Error messages, `[FAIL]` label |
| **Warning** | Bright Yellow | `color.FgYellow` | Missing variable warnings, deprecation notices |
| **Success / PASS** | Bold Green | `color.FgGreen + color.Bold` | `[PASS]`, success checkmarks |
| **URL** | Bright Cyan | `color.FgCyan` | Full resolved URL in output |
| **URL (unresolved)** | Dim Yellow | `color.FgYellow + color.Dim` | `{{BASE_URL}}` with unresolved variable hint |
| **Header Key (important)** | Bright Yellow | `color.FgYellow + color.Bold` | Auth, Content-Type, Accept keys |
| **Header Key (standard)** | Dark Yellow | `color.FgYellow` | Standard header names |
| **Header Value (sensitive)** | Dim White | `color.FgWhite + color.Dim` | Token values masked as `***` |
| **Header Value (standard)** | White | `color.FgWhite` | Non-sensitive header values |
| **Query Param Key** | Yellow | `color.FgYellow` | Parameter names |
| **Query Param Value** | White | `color.FgWhite` | Parameter values |
| **Priority Badge [high]** | Bright Green | `color.FgGreen` | `[priority: high]` label |
| **Priority Badge [low]** | Dim White | `color.FgWhite + color.Dim` | `[priority: low]` label |
| **Body Key (important)** | Bright Magenta | `color.FgMagenta + color.Bold` | name, email, token, etc. |
| **Body Key (standard)** | Magenta | `color.FgMagenta` | Standard body field names |
| **Body Value (string)** | White | `color.FgWhite` | String values `"..."` |
| **Body Value (number)** | Cyan | `color.FgCyan` | Numeric values |
| **Body Value (boolean)** | Yellow | `color.FgYellow` | true / false |
| **Body Value (null)** | Dim White | `color.FgWhite + color.Dim` | null |
| **Response Key (important)** | Bright Cyan | `color.FgCyan + color.Bold` | id, name, email, etc. |
| **Response Key (standard)** | Cyan | `color.FgCyan` | Standard response field names |
| **Response Value (string)** | White | `color.FgWhite` | String values |
| **Response Value (number)** | Bright Green | `color.FgGreen` | Numeric values |
| **Response Value (boolean)** | Yellow | `color.FgYellow` | true / false |
| **Metadata / Timestamps** | Dim Gray | `color.FgBlack` | "Time:", "Saved:", timestamps, hints |
| **Interactive Prompt ?** | Bold White | `color.FgWhite + color.Bold` | `? URL:`, `? Method:` prompts |
| **Interactive Input** | Bright White | `color.FgHiWhite` | User-typed input display |
| **Inline Variable `{{VAR}}`** | Bright Yellow | `color.FgYellow + color.Bold + color.Underline` | Variable tokens in prompts and preview |
| **Resolved Variable** | Dim White | `color.FgWhite + color.Dim` | `(resolved from config)` hints |
| **Test Suite Title** | Bold White | `color.FgWhite + color.Bold` | "TEST SUITE: ..." |
| **Test Name (passing)** | Green | `color.FgGreen` | Test names with OK status |
| **Test Name (failing)** | Red | `color.FgRed` | Test names with FAIL status |
| **Test Name (skipped)** | Dim White | `color.FgWhite + color.Dim` | Test names with SKIP status |
| **Test Progress dots** | Dim Gray | `color.FgBlack` | `............` progress |
| **Command Name** | Cyan | `color.FgCyan` | `hsp`, `request`, `get`, etc. in help |
| **Command Description** | White | `color.FgWhite` | Help text descriptions |
| **Suggestion / Hint** | Dim Cyan | `color.FgCyan + color.Dim` | Default values, inline suggestions |
| **Masked Secret** | Dim White | `color.FgWhite + color.Dim` | `***`, `eyJ...`, etc. |

### Environment Color Coding

| Environment | Color | Usage |
|------------|-------|-------|
| `default` | Dim White | Default environment, no special color |
| `dev` | Green | Development environment label |
| `staging` | Yellow | Staging environment label |
| `prod` | Red + Bold | Production environment label (warning color) |

### Priority Shading Logic

- **Important keys:** `id`, `name`, `email`, `token`, `key`, `status`, `message`, `error`, `code`, `url`, `type`, `data`, `result`.
- **Less important:** Numeric IDs, timestamps, internal fields, pagination fields.
- Priority detected automatically by key name matching OR user can mark fields with `@important` comment in body.

### TUI Panel Output

Pure ASCII box drawing characters for maximum compatibility. Font weights create visual hierarchy.

```
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
[white bold]|  REQUEST                                                           [cyan bold]POST[cyan bold]        |
[cyan bold]|[cyan bold]  [dim]Time: 143ms                                                                |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  URL       : [cyan]https://{{BASE_URL}}/users[yellow]   (unresolved: {{BASE_URL}})[yellow]                    |
[cyan bold]|[cyan bold]  Resolved  : [cyan underline]https://api.example.com/users[cyan underline]                                               |
[dim]|[dim]             [dim](resolved from config)[dim]                                                    |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  Headers   : [yellow bold]Authorization[yellow bold]: [white dim]Bearer ***[white dim]                                           |
[cyan bold]|[cyan bold]             [yellow bold]Content-Type[yellow bold]: [white]application/json[white]                                             |
[cyan bold]|[cyan bold]             [yellow bold]Accept[yellow bold]: [white]application/json[white]                                                   |
[cyan bold]|[cyan bold]  Params    : [yellow bold]page[yellow bold]=[cyan]1[cyan] | [yellow bold]limit[yellow bold]=[cyan]20[cyan]                                                      |
[cyan bold]|[cyan bold]  Env       : [green]dev[green]                                                                    |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[white bold]|  BODY (Payload)                                                                   |
[cyan bold]|[cyan bold]  +--------------------------------------------------------------------------+ |
[cyan bold]|[cyan bold]  | [magenta bold]name      [white]: [white]"John Doe"[white]              [green][priority: high][green]    |
[cyan bold]|[cyan bold]  | [magenta bold]email     [white]: [white]"john@example.com"[white]  [green][priority: high][green]  |
[cyan bold]|[cyan bold]  | [magenta]userId    [white]: [cyan bold]1[cyan bold]                [dim][priority: low][dim]       |
[cyan bold]|[cyan bold]  +--------------------------------------------------------------------------+ |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  [green bold]201 Created[green bold]                                                                  |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[white bold]|  RESPONSE                                                                          |
[cyan bold]|[cyan bold]  +--------------------------------------------------------------------------+ |
[cyan bold]|[cyan bold]  | [cyan bold]id        [white]: [green bold]123[green bold]                   [green][priority: high][green]  |
[cyan bold]|[cyan bold]  | [cyan bold]name      [white]: [white]"John Doe"[white]             [green][priority: high][green]  |
[cyan bold]|[cyan bold]  | [cyan bold]email     [white]: [white]"john@example.com"[white]       [green][priority: high][green] |
[cyan bold]|[cyan bold]  | [cyan bold]created_at[white]: [dim]"2026-04-26T12:00:00Z"[dim]      [dim][priority: low][dim]   |
[cyan bold]|[cyan bold]  +--------------------------------------------------------------------------+ |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  [green]Request saved[green]: ~/.hsp/history/POST_2026-04-26_12-00-00.json                          |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
```

### Interactive TUI View (hsp request)

Step-by-step prompts with clean ASCII section separators between each step, all prompt text and input in defined colors with proper font weight hierarchy.

```
[yellow bold]---[yellow bold] [yellow]Step 1: URL[yellow] [yellow bold]---[yellow bold]
[white bold]? URL: [bright-white]https://{{BASE_URL}}/users[bright-white]
[dim]    Resolved: [cyan underline]https://api.example.com/users[cyan underline]
[dim]    (resolved from config)

[yellow bold]---[yellow bold] [yellow]Step 2: Method[yellow] [yellow bold]---[yellow bold]
[white bold]? Method: [dim](default: GET)[dim]
[white]    1) [cyan bold]GET[cyan bold]       [dim][retrieve data][dim]
[white]    2) [cyan bold]POST[cyan bold]      [dim][create new resource][dim]
[white]    3) [cyan bold]PUT[cyan bold]       [dim][update entire resource][dim]
[white]    4) [cyan bold]PATCH[cyan bold]     [dim][partial update][dim]
[white]    5) [cyan bold]DELETE[cyan bold]    [dim][remove resource][dim]
[white]    6) [cyan bold]HEAD[cyan bold]      [dim][like GET, no body][dim]
[white]    7) [cyan bold]OPTIONS[cyan bold]   [dim][allowed methods][dim]
[white bold]Choose (1-7) or type method: [bright-white]2[bright-white]

[green bold]    Method: POST[green bold]

[yellow bold]---[yellow bold] [yellow]Step 3: Headers[yellow] [yellow bold]---[yellow bold]
[white bold]? Add headers? (y/n): [bright-white]n[bright-white]
[dim]    [Auto-set: Accept: application/json][dim]

[yellow bold]---[yellow bold] [yellow]Step 4: Query Parameters[yellow] [yellow bold]---[yellow bold]
[white bold]? Add query parameters? (y/n): [bright-white]n[bright-white]

[yellow bold]---[yellow bold] [yellow]Step 5: Body[yellow] [yellow bold]---[yellow bold]
[white bold]? Add request body? (y/n): [bright-white]y[bright-white]
[white]    Body format:
[white]    1) [magenta bold]JSON[magenta bold]
[white]    2) [yellow]Form data[yellow]
[white]    3) [dim]Raw text[dim]
[white bold]Choose (1-3): [bright-white]1[bright-white]
[white bold]Enter JSON body (press Enter twice when done):[white bold]
[bright-white]{
  [yellow bold]"name"[yellow bold]: [yellow]"{{TEST_USER}}"[yellow],
  [yellow bold]"email"[yellow bold]: [white]"test@example.com"[white]
}[bright-white]
[green bold]    Valid JSON[green bold]
[dim]    [Auto-set: Content-Type: application/json][dim]

[yellow bold]---[yellow bold] [yellow]Step 6: Preview[yellow] [yellow bold]---[yellow bold]
```

Request preview with full resolved URL and color-coded sections.

```
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
[white bold]|  PREVIEW                                                                          |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[white bold]|  [cyan bold]POST[cyan bold] [cyan underline]https://api.example.com/users[cyan underline]                                          |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  Headers   : [yellow bold]Authorization[yellow bold]: [white dim]Bearer ***[white dim]                                          |
[cyan bold]|[cyan bold]             [yellow bold]Content-Type[yellow bold]: [white]application/json[white]                                           |
[cyan bold]|[cyan bold]             [yellow bold]Accept[yellow bold]: [white]application/json[white]                                               |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  Body      : [magenta bold]name[magenta bold]: [white]"John Doe"[white]              [green][priority: high][green]       |
[cyan bold]|[cyan bold]               [magenta bold]email[magenta bold]: [white]"test@example.com"[white]    [green][priority: high][green]      |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[white bold]|  ? Send request? (y/n): [bright-white]y[bright-white]                                                         |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
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

Clean ASCII output showing pass/fail per test with color-coded results.

```
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
[white bold]|  TEST SUITE: User API Test Suite                                        [green]dev[green]  |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[green]  [green bold]PASS[green bold]  [green]Create user                      .......................... OK[green]              |
[green]  [green bold]PASS[green bold]  [green]Get created user                ........................ OK[green]              |
[green]  [green bold]PASS[green bold]  [green]Update user                     ......................... OK[green]              |
[green]  [green bold]PASS[green bold]  [green]Delete user                     ......................... OK[green]              |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  Summary: [green bold]4/4 passed[green bold] (143ms total)                                                  |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
```

### Failed Test Output

```
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
[white bold]|  TEST SUITE: User API Test Suite                                        [green]dev[green]  |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[green]  [green bold]PASS[green bold]  [green]Create user                      .......................... OK[green]              |
[green]  [green bold]PASS[green bold]  [green]Get created user                ........................ OK[green]              |
[red bold]  FAIL[red bold]   [red]Update user                     ......................... FAIL[red]              |
[dim]          [dim]Delete user                     ........................ SKIP[dim]              |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[yellow bold]  FAILURES:[yellow bold]                                                                           |
[red bold]  [2] [red bold]Update user[red bold]                                                                          |
[white]    Assertion: [yellow bold]status[yellow bold]                                                                          |
[white]    Expected: [green bold]200[green bold]                                                                          |
[white]    Actual:   [red bold]500[red bold] [red bold]Internal Server Error[red bold]                                                  |
[white]    Body:     {"error": "Database connection failed"}                                           |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]+
[cyan bold]|[cyan bold]  Summary: [red bold]1/4 passed[red bold] (143ms total)                          [red bold]FAIL[red bold]                       |
[cyan bold]+[cyan bold]------------------------------------------------------------------------------[cyan bold]+[cyan bold]
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