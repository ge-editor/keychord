# keychord

`keychord` is a flexible key binding management library for CUI / TUI applications,
built on top of `github.com/gdamore/tcell/v3`.

- Supports single keys and hierarchical key sequences (key chords)
- Supports Ctrl / Alt / Shift / Meta modifiers
- Prefix-key style bindings
- Designed to support application mode switching (modal UI)
- Directly processes `tcell.EventKey`

---

## Features

### ✅ Based on tcell/v3

This library uses [`github.com/gdamore/tcell/v3`](https://github.com/gdamore/tcell)
to handle:

- Terminal-dependent key input
- Modifier keys (Ctrl / Alt / Shift / Meta)
- Special keys (Esc, Enter, Delete, etc.)

---

### ✅ Hierarchical Key Bindings (Key Chords)

You can define key sequences such as:

```text
Ctrl+X Ctrl+S
g g
Ctrl+C Esc
```

While a prefix key is active, the internal state is preserved
and the dispatcher waits for the next key input.

---

### ✅ Designed for Modal Applications

`RootNode` maintains internal state (`Current`), making it easy to support
modal application designs such as:

- Normal / Insert / Visual
- Command / Search
- Application-specific operation modes

By preparing a separate `RootNode` for each mode, you can switch key maps
cleanly and explicitly.

```go
normalMode := keychord.NewRootNode()
insertMode := keychord.NewRootNode()
```

---

### ✅ Human-friendly Key Notation

Key bindings are defined using simple string notation:

```text
"Ctrl+X"
"Alt+Enter"
"Esc"
"g"
"Ctrl+A"
```

These strings are decoded internally into `tcell.Key` and `tcell.ModMask`.

---

## Installation

```bash
go get github.com/ge-editor/keychord
```

---

## Usage

### Defining Key Bindings

```go
root := keychord.NewRootNode()

root.Bind("Ctrl+X", "Ctrl+S").Do(func() {
    saveFile()
})

root.Bind("g", "g").Do(func() {
    goToTop()
})
```

---

### Dispatching Events

Pass `tcell.EventKey` directly to the dispatcher.

```go
status, result := root.Dispatch(ev)

switch result {
case keychord.DispatchExecuted:
    // action executed
case keychord.DispatchPrefix:
    // waiting for next key in a sequence
case keychord.DispatchNotFound:
    // no matching binding
}
```

`status` contains the currently entered key sequence (e.g. `C-x`),
which can be used for status bars or on-screen hints.

---

### Listing Candidate Keys

```go
candidates := root.Candidates()
```

This allows you to implement features such as:

- Contextual help
- which-key–style UIs

---

### Supported Key Types

- Special keys (as defined by tcell)
- ASCII characters
- Single Unicode characters
- Ctrl+A through Ctrl+Z
- Control characters such as ^@, ^[, ^\, ^], ^^, ^_

---

## Logging

Internal logging uses `github.com/ge-editor/gelog`.

- Logging is enabled only in debug builds
- Logging calls are no-ops in release builds

---

## License

MIT License

---
