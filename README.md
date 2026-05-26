# keychord

`keychord` is a flexible key binding management library for CUI / TUI applications,
built on top of `github.com/gdamore/tcell/v3`.

* Supports single keys and hierarchical key sequences (key chords)
* Supports Ctrl / Alt / Shift / Meta modifiers
* Prefix-key style bindings
* Explicit state transitions (DFA-based dispatcher)
* Designed for modal UI architectures
* Directly processes `tcell.EventKey`

---

## Features

### ✅ Based on tcell/v3

This library uses **gdamore/tcell** (`github.com/gdamore/tcell/v3`) to handle:

* Terminal-dependent key input
* Modifier keys (Ctrl / Alt / Shift / Meta)
* Special keys (Esc, Enter, Delete, etc.)

---

### ✅ Deterministic Key Dispatcher (DFA)

The dispatcher is implemented as a **deterministic finite automaton (DFA)**.

It maintains one of two internal states:

* `Root`
* `Prefix(n)`

State transitions are explicit and predictable.

Returned transitions:

* `DispatchNotFound`
* `DispatchPrefix`
* `DispatchExecuted`
* `DispatchInvalidAfterPrefix`

This allows applications to precisely distinguish:

* No binding exists
* Waiting for next key
* Action executed
* Invalid continuation after a prefix

---

### ✅ Hierarchical Key Bindings (Key Chords)

You can define key sequences such as:

```text
Ctrl+X Ctrl+S
g g
Ctrl+C Esc
```

When a prefix key is entered:

* Internal state transitions to `Prefix(n)`
* Dispatcher waits for the next key
* Candidate keys can be queried
* Invalid continuation returns `DispatchInvalidAfterPrefix`

Example:

```
C-x   → DispatchPrefix
C-x z → DispatchInvalidAfterPrefix
```

---

### ✅ Designed for Modal Applications

`RootNode` maintains internal state (`current`), making it easy to support
modal application designs such as:

* Normal / Insert / Visual
* Command / Search
* Application-specific operation modes

Prepare a separate `RootNode` per mode:

```go
normalMode := keychord.NewRootNode()
insertMode := keychord.NewRootNode()
```

Switching modes simply means switching the active `RootNode`.

No global state required.

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

These strings are decoded internally into:

* `tcell.Key`
* `tcell.ModMask`

Supported:

* Special keys (from tcell)
* ASCII characters
* Single Unicode characters
* Ctrl+A through Ctrl+Z
* Control characters (^@, ^[, ^, ^], ^^, ^_)

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

Non-blocking bindings:

```go
root.Bind("Ctrl+L").DoAlso(func() {
    refreshScreen()
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

case keychord.DispatchInvalidAfterPrefix:
    // invalid continuation after prefix

case keychord.DispatchNotFound:
    // no matching binding
}
```

`status` contains the currently entered key sequence
(e.g. `C-x`), suitable for:

* Status bars
* Key hints
* Debug display

---

## State Transition Diagram

```
            +-------------------+
            |       Root        |
            +-------------------+
             |   |         |
     NotFound|   |Action   |Prefix
             |   |         v
             |   |    +----------------+
             |   |    |   Prefix(n)    |
             |   |    +----------------+
             |   |      |   |        |
             |   |      |   |Action  |Prefix
             |   |      |   |        v
             |   |      |   |   Prefix(n')
             |   |      |   |
             |   |      |Invalid
             |   |      v
             |   |     Root
             v   v
            Root Root
```

Formal definition:

```
δ : S × Σ → S × O
```

Where:

```
S = { Root, Prefix(n) }
O = {
  DispatchNotFound,
  DispatchPrefix,
  DispatchExecuted,
  DispatchInvalidAfterPrefix
}
```

The dispatcher is deterministic.

---

### Listing Candidate Keys

```go
candidates := root.Candidates()
```

Useful for:

* Contextual help
* which-key–style UI
* Interactive key discovery

---

## Logging

Internal logging uses `github.com/ge-editor/gelog`.

* Logging enabled in debug builds
* No-op in release builds
* Each dispatch cycle has a unique ID for traceability

---

## Design Principles

* Deterministic state machine
* Explicit reset semantics
* No hidden global state
* Prefix failure is distinguishable
* Safe for modal editor architectures

---

## License

MIT License

---
