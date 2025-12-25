# keychord

`keychord` は、`github.com/gdamore/tcell/v3` を用いた CUI / TUI アプリケーション向けの柔軟なキーバインド管理ライブラリです。

- 単一キー／階層キー（キーシーケンス）に対応
- Ctrl / Alt / Shift / Meta 修飾キー対応
- prefix key（前置キー）方式のキーバインド
- アプリケーションのモード切替（modal UI）に対応可能な設計
- `tcell.EventKey` を直接処理

---

## 特徴

### ✅ tcell/v3 ベース

本ライブラリは [`github.com/gdamore/tcell/v3`](https://github.com/gdamore/tcell) を使用し、下記を扱います。

- 端末依存のキー入力
- 修飾キー（Ctrl / Alt / Shift / Meta）
- 特殊キー（Esc, Enter, Delete など）

---

### ✅ 階層的なキーバインド（Key Chord）

以下のようなキーシーケンスを定義できます。

```text
Ctrl+X Ctrl+S
g g
Ctrl+C Esc
```

プレフィックスキー入力中は状態を保持し、
次の入力を待つ動作が可能です。

### ✅ モード切替に対応した設計

RootNode は状態（Current）を持つため、

- Normal / Insert / Visual
- Command / Search
- アプリ固有の操作モード

など、モードごとに別の RootNode を用意することで
エディタや CUI アプリ特有のモード切替設計に対応できます。

```go
normalMode := keychord.NewRootNode()
insertMode := keychord.NewRootNode()
```

### ✅ 人間が定義しやすいキー表記

キー定義は文字列で行います。

```
"Ctrl+X"
"Alt+Enter"
"Esc"
"g"
"Ctrl+A"
```

内部で tcell.Key / tcell.ModMask にデコードされます。

---

## インストール

```bash
go get github.com/ge-editor/keychord
```

## 使い方

### キーバインド定義

```go
root := keychord.NewRootNode()

root.Bind("Ctrl+X", "Ctrl+S").Do(func() {
    saveFile()
})

root.Bind("g", "g").Do(func() {
    goToTop()
})
```

### イベントのディスパッチ

tcell.EventKey をそのまま渡します。

```go
status, result := root.Dispatch(ev)

switch result {
case keychord.DispatchExecuted:
    // アクション実行済み
case keychord.DispatchPrefix:
    // プレフィックスキー入力中
case keychord.DispatchNotFound:
    // 該当なし
}
```

status には現在入力中のキー列（例: C-x）が入るため、
ステータス表示などに利用できます。

### 候補キーの取得

```go
candidates := root.Candidates()
```

入力可能なキー一覧を取得できるため、

- ヘルプ表示
- which-key 風 UI

などの実装が可能です。

### 対応キー種別

- 特殊キー（tcell 定義の Key）
- ASCII 文字
- Unicode 1 文字
- Ctrl+A ～ Ctrl+Z
- ^@, ^[, ^, ^], ^^, ^_ など制御文字

---

## ログについて

内部ログには `github.com/ge-editor/gelog` を使用しています。

- debug ビルド時のみログ出力
- release ビルドでは no-op

---

## ライセンス

MIT License

---
