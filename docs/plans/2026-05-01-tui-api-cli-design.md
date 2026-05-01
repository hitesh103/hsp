# HSP Interactive TUI Design

## Overview
Transform the current sequential prompt-based `hsp request` interactive mode into a rich, full-screen Terminal User Interface (TUI) to provide a best-in-class "Postman-in-the-terminal" experience.

## Architecture & UI Components
The core interactive flow in `cmd/request.go` will be rewritten using `charmbracelet/bubbletea` and `charmbracelet/bubbles`.

**Layout:**
- **Top Section:** Dynamic input fields for URL (with history/variable auto-complete) and a method selector.
- **Middle Section:** A tabbed interface allowing navigation between 'Headers', 'Query Params', and 'Body'. The Body tab provides a preview of the payload.
- **Bottom Section:** A persistent status bar indicating available hotkeys: `Ctrl+E` (Open Editor), `Ctrl+I` (Import cURL), `Tab` (Next Field), and `Ctrl+S` (Send Request).

**Data Flow:**
A central `RequestBuilder` model acts as the state. User navigation and typing update this state. Pressing `Ctrl+S` triggers a loading spinner, dispatches the HTTP request via the existing engine, and exits the TUI to render the output.

## Advanced Features Integration
- **Editor Integration:** `Ctrl+E` suspends the Bubble Tea program, writes the current body to a temporary file, opens `$EDITOR`, and upon closing, reads the changes back into the TUI state before resuming.
- **cURL Import:** `Ctrl+I` opens a multi-line input overlay. Pasted cURL commands are parsed to extract Method, URL, Headers, and Body, immediately updating the TUI fields.
- **Rich Auto-completion:** Inputs will provide real-time suggestions (e.g., standard HTTP headers, base URLs from active environment variables).

## Error Handling & Testing
- **Validation:** Attempting to send (`Ctrl+S`) with an invalid URL or malformed JSON triggers an inline error flash in the status bar instead of crashing the program.
- **Testing:** Unit tests will verify the `Update` loop (ensuring correct state transitions on keystrokes) and specifically test the cURL parser string extraction logic.
