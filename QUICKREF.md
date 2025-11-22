# HSP Quick Reference

## Commands

```bash
hsp request          # Start interactive request builder (RECOMMENDED)
hsp get <url>       # Quick GET request
hsp post <url>      # Quick POST request
hsp --help          # Show help
```

## Interactive Flow (hsp request)

```
1. Enter URL
   ✓ https://api.example.com/endpoint

2. Select Method
   ✓ GET / POST / PUT / PATCH / DELETE / HEAD / OPTIONS

3. Add Headers (optional)
   ✓ Authorization: Bearer token
   ✓ X-API-Key: secret-key

4. Add Query Parameters (optional)
   ✓ page: 1
   ✓ limit: 20

5. Add Body (if POST/PUT/PATCH)
   ✓ Format: JSON / Form Data / Raw Text

6. Pretty-print Response?
   ✓ Yes / No

7. Review Preview & Confirm
   ✓ Send request
```

## Examples

### GET Request with Headers
```bash
hsp request
? URL: https://api.github.com/user
? Method: GET
? Add headers? y
  Header name: Authorization
  Value: token gh_xxxx
  Add another? n
? Add query params? n
? Pretty? y
? Send? y
```

### POST with JSON
```bash
hsp request
? URL: https://jsonplaceholder.typicode.com/posts
? Method: POST
? Add headers? n
? Add query params? n
? Add body? y
? Format: 1 (JSON)
Enter JSON:
{
  "title": "Test",
  "body": "Test body",
  "userId": 1
}

? Pretty? y
? Send? y
```

### PUT with Form Data
```bash
hsp request
? URL: https://api.example.com/users/1
? Method: PUT
? Add headers? y
  Header name: Authorization
  Value: Bearer token
  Add another? n
? Add query params? n
? Add body? y
? Format: 2 (Form data)
  Field: name
  Value: John Doe
  Add another? n
? Pretty? y
? Send? y
```

## History

All requests auto-saved:
```bash
ls ~/.hsp/history/
cat ~/.hsp/history/POST_*.json
```

## Status Colors

- 🟢 **Green**: 2xx (Success)
- 🔴 **Red**: 4xx/5xx (Error)

## Tips & Tricks

✨ Type 'done' to finish adding headers/params
✨ Press Ctrl+C to cancel
✨ JSON is auto-validated before sending
✨ Content-Type auto-set for JSON bodies
✨ Query params URL-encoded automatically
✨ Response headers always shown
✨ Requests stored with timestamps

## Environment Variables

```bash
export HSP_HISTORY_DIR="$HOME/Documents/requests"  # Change history location
export HSP_PRETTY=true                             # Default pretty-print
export HSP_TIMEOUT=30s                             # Request timeout
```

## What Makes HSP Different?

| Feature | HSP | curl | Postman |
|---------|-----|------|---------|
| Interactive | ✅ | ❌ | ✅ |
| Terminal-only | ✅ | ✅ | ❌ |
| Easy headers | ✅ | ⚠️ | ✅ |
| Auto history | ✅ | ❌ | ✅ |
| Zero config | ✅ | ✅ | ❌ |
| Fast | ✅ | ✅ | ⚠️ |
| Beautiful output | ✅ | ⚠️ | ✅ |
| Lightweight | ✅ | ✅ | ❌ |

---
**Stay productive. Stay in the terminal. 🚀**
