package keychord

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"

	"github.com/ge-editor/gelog"
)

type KeySpec struct {
	Key tcell.Key
	Str string
	Mod tcell.ModMask
}

var (
	ErrInvalidKeyEvent = errors.New("error invalid key event")
)

type DecodedKey struct {
	Key tcell.Key     // 特殊キー
	Str string        // printable 文字
	Mod tcell.ModMask // ctrl / alt / shift
}

// KeySpec を返すユーティリティ
func (d DecodedKey) KeySpec() KeySpec {
	return KeySpec{
		Key: d.Key,
		Str: d.Str,
		Mod: d.Mod,
	}
}

func Decode(s string) (DecodedKey, error) {
	var out DecodedKey
	if s == "" {
		gelog.Info("Blank, invalid", "s", s)
		return out, ErrInvalidKeyEvent
	}

	parts := strings.Split(s, "+")
	var keyPart string

	for _, p := range parts {
		pl := strings.ToLower(p)
		switch pl {
		case "ctrl":
			out.Mod |= tcell.ModCtrl
		case "alt":
			out.Mod |= tcell.ModAlt
		case "shift":
			out.Mod |= tcell.ModShift
		case "meta":
			out.Mod |= tcell.ModMeta
		default:
			keyPart = p
		}
	}

	if keyPart == "" {
		gelog.Info("keyPart is blank, invalid", "keyPart", keyPart, "s", s)
		return out, ErrInvalidKeyEvent
	}

	// Specila Keys
	// "Delete", "Esc", ... tcell/v3 の key.go を見てね
	for k, name := range tcell.KeyNames {
		if strings.EqualFold(keyPart, name) {
			gelog.Info("Special key", "s", s, "KeyName", name)
			out.Key = k
			return out, nil
		}
	}

	// 0 	000 0000 	000 	0x00 	NUL ^@ 	Null
	//
	// 27 	001 1011 	033 	0x1b 	ESC	^[ 	Escape
	// 28 	001 1100 	034 	0x1c 	FS 	^\ 	File Separator
	// 29 	001 1101 	035 	0x1d 	GS 	^] 	Group Separator
	// 30 	001 1110 	036 	0x1e 	RS 	^^ 	Record Separator
	// 31 	001 1111 	037 	0x1f 	US 	^_ 	Unit Separator
	// 32 	010 0000 	040 	0x20 	SP 	^@	Space
	//
	// 33-47	! - /
	// 48-57	0-9		Number
	//
	// 127 	111 1111 	177 	0x7f 	DEL 	^? 	delete character

	// 2) printable: KeyRune
	if len(keyPart) == 1 {
		r := rune(keyPart[0])

		// Ctrl+A - Ctrl+Z, ^@, ^[, ^\, ^], ^^, ^_
		if out.Mod&tcell.ModCtrl != 0 {
			lower := unicode.ToLower(r)
			code := tcell.Key(r) - '@'
			if code == 27 { // ^[
				out = DecodedKey{Key: code}
				return out, nil
			}
			if code == 0 { // ^@ (^Space)
				gelog.Info("^@ (^Space)", "r", r, "s", s, "code", code, "str", out.Str, "key", out.Key, "spec", out.KeySpec())
				out = DecodedKey{Key: 256, Str: " ", Mod: tcell.ModCtrl}
				return out, nil
			}
			if code >= 28 && code <= 31 { // ^\, ^], ^^, ^_ (^/)
				gelog.Info("code: 28 to 31", "r", r, "s", s, "code", code, "str", out.Str, "key", out.Key, "spec", out.KeySpec())
				out = DecodedKey{Key: 256, Str: string(r), Mod: tcell.ModCtrl}
				return out, nil
			}
			if lower >= 'a' && lower <= 'z' {
				out.Key = tcell.Key(lower - 'a' + rune(tcell.KeyCtrlA))
				gelog.Info("Ctrl+A to Ctrl+Z", "r", r, "s", s, "code", code, "str", out.Str, "key", out.Key, "spec", out.KeySpec())
				return out, nil
			}
			// Ctrl + 非ASCII → 無効
			// key: 0
			gelog.Info("Ctrl + Non ASCII: invalid", "r", r, "s", s, "code", code, "str", out.Str, "key", out.Key, "spec", out.KeySpec())
			return out, ErrInvalidKeyEvent
		}

		// Rune 普通の文字
		out.Key = tcell.KeyRune
		out.Str = keyPart
		gelog.Info("Rune (standard string)", "s", s, "str", out.Str)
		return out, nil
	}

	// 3) Unicode 1文字
	if utf8.RuneCountInString(keyPart) == 1 {
		out.Key = tcell.KeyRune
		out.Str = keyPart
		gelog.Info("Unicode (a charactor)", "s", s, "str", out.Str)
		return out, nil
	}

	gelog.Info("Invalid", "s", s)
	return out, ErrInvalidKeyEvent
}

type RootNode struct {
	Root      *KeyNode
	Current   *KeyNode
	KeyStatus string
}

// NewRootNode
func NewRootNode() *RootNode {
	root := NewKeyNode()
	return &RootNode{
		Root:    root,
		Current: root,
	}
}

// Bind 単一キーまたは階層キー列にアクションを登録
func (r *RootNode) bind(keyStrs []string, action func()) error {
	if len(keyStrs) == 0 {
		return errors.New("empty key sequence")
	}
	node := r.Root
	for i, ks := range keyStrs {
		k, err := Decode(ks)
		if err != nil {
			return fmt.Errorf("invalid key %q: %w", ks, err)
		}
		if node.Next == nil {
			node.Next = make(map[KeySpec]*KeyNode)
		}
		next, ok := node.Next[k.KeySpec()]
		if !ok {
			next = NewKeyNode()
			node.Next[k.KeySpec()] = next
		}

		// KeyStr をセット（作成時のみ）
		if next.KeyStr == "" {
			next.KeyStr = ks
		}
		if i == len(keyStrs)-1 {
			next.Action = action
		}
		node = next
	}
	return nil
}

// Reset 現在のノードを初期化
func (r *RootNode) Reset() {
	r.Current = r.Root
	r.KeyStatus = ""
}

type DispatchResult int

const (
	DispatchNotFound DispatchResult = iota // 該当するキーがない
	DispatchPrefix                         // プレフィックスキーで次の入力を待つ
	DispatchExecuted                       // アクション実行済み
)

// Dispatch 入力イベントに応じてキーバインドを実行
func (r *RootNode) Dispatch(ev *tcell.EventKey) (string, DispatchResult) {
	k := KeySpec{
		Key: ev.Key(),
		Str: ev.Str(),
		Mod: ev.Modifiers(),
	}

	if r.Current == nil {
		r.Current = r.Root
	}

	node, ok := r.Current.Next[k]
	if !ok {
		// 存在しないキー → 状態リセット
		r.KeyStatus = ""
		r.Current = r.Root
		return r.KeyStatus, DispatchNotFound
	}

	if node.Action != nil {
		node.Action()
		r.Reset()
		return r.KeyStatus, DispatchExecuted
	}

	if len(node.Next) > 0 {
		// プレフィックスキー → 次の入力待ち
		r.Current = node
		r.appendKeyStatus(node)
		return r.KeyStatus, DispatchPrefix
	}

	// 到達不能だが安全策
	r.KeyStatus = ""
	r.Current = r.Root
	return r.KeyStatus, DispatchNotFound
}

func (r *RootNode) appendKeyStatus(n *KeyNode) {
	if r.KeyStatus != "" {
		r.KeyStatus += " "
	}
	r.KeyStatus += formatMods(n.KeyStr)
}

func (r *RootNode) Candidates() []string {
	if r.Current == nil || r.Current.Next == nil {
		return nil
	}

	out := make([]string, 0, len(r.Current.Next))
	for _, node := range r.Current.Next {
		out = append(out, formatMods(node.KeyStr))
	}
	sort.Strings(out)
	return out
}

type binder struct {
	root *RootNode
	keys []string
}

func (r *RootNode) Bind(keys ...string) *binder {
	return &binder{root: r, keys: keys}
}

func (b *binder) Do(action func()) error {
	return b.root.bind(b.keys, action)
}

// KeyNode はキーバインドのノード（階層構造）
type KeyNode struct {
	Next   map[KeySpec]*KeyNode
	Action func()
	KeyStr string
}

// NewKeyNode 作成
func NewKeyNode() *KeyNode {
	return &KeyNode{
		Next: make(map[KeySpec]*KeyNode),
	}
}

func formatMods(s string) string {
	repls := []struct {
		from string
		to   string
	}{
		{"ctrl+", "C-"},
		{"alt+", "M-"},
		{"shift+", "S-"},
		{"meta+", "Meta-"}, // or "M-" にするならここを調整
	}

	ls := strings.ToLower(s)
	for _, r := range repls {
		if strings.HasPrefix(ls, r.from) {
			return r.to + s[len(r.from):]
			// return r.to + ls[len(r.from):] // lower case
		}
	}
	return s
}
