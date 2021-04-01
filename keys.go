package main

import (
	"bytes"
)

type KeyType string

const (
	Escape             KeyType = "Escape"
	ControlSpace       KeyType = "ControlSpace"
	ControlA           KeyType = "ControlA"
	ControlB           KeyType = "ControlB"
	ControlC           KeyType = "ControlC"
	ControlD           KeyType = "ControlD"
	ControlE           KeyType = "ControlE"
	ControlF           KeyType = "ControlF"
	ControlG           KeyType = "ControlG"
	ControlH           KeyType = "ControlH"
	ControlK           KeyType = "ControlK"
	ControlL           KeyType = "ControlL"
	ControlM           KeyType = "ControlM"
	ControlN           KeyType = "ControlN"
	ControlO           KeyType = "ControlO"
	ControlP           KeyType = "ControlP"
	ControlQ           KeyType = "ControlQ"
	ControlR           KeyType = "ControlR"
	ControlS           KeyType = "ControlS"
	ControlT           KeyType = "ControlT"
	ControlU           KeyType = "ControlU"
	ControlV           KeyType = "ControlV"
	ControlW           KeyType = "ControlW"
	ControlX           KeyType = "ControlX"
	ControlY           KeyType = "ControlY"
	ControlZ           KeyType = "ControlZ"
	ControlBackslash   KeyType = "ControlBackslash"
	ControlSquareClose KeyType = "ControlSquareClose"
	ControlCircumflex  KeyType = "ControlCircumflex"
	ControlUnderscore  KeyType = "ControlUnderscore"
	Backspace          KeyType = "Backspace"
	Up                 KeyType = "Up"
	Down               KeyType = "Down"
	Right              KeyType = "Right"
	Left               KeyType = "Left"
	Home               KeyType = "Home"
	End                KeyType = "End"
	Enter              KeyType = "Enter"
	Delete             KeyType = "Delete"
	ShiftDelete        KeyType = "ShiftDelete"
	ControlDelete      KeyType = "ControlDelete"
	PageUp             KeyType = "PageUp"
	PageDown           KeyType = "PageDown"
	Tab                KeyType = "Tab"
	BackTab            KeyType = "BackTab"
	Insert             KeyType = "Insert"
	F1                 KeyType = "F1"
	F2                 KeyType = "F2"
	F3                 KeyType = "F3"
	F4                 KeyType = "F4"
	F5                 KeyType = "F5"
	F6                 KeyType = "F6"
	F7                 KeyType = "F7"
	F8                 KeyType = "F8"
	F9                 KeyType = "F9"
	F10                KeyType = "F10"
	F11                KeyType = "F11"
	F12                KeyType = "F12"
	F13                KeyType = "F13"
	F14                KeyType = "F14"
	F15                KeyType = "F15"
	F16                KeyType = "F16"
	F17                KeyType = "F17"
	F18                KeyType = "F18"
	F19                KeyType = "F19"
	F20                KeyType = "F20"
	F21                KeyType = "F21"
	F22                KeyType = "F22"
	F23                KeyType = "F23"
	F24                KeyType = "F24"
	ControlUp          KeyType = "ControlUp"
	ControlDown        KeyType = "ControlDown"
	ControlRight       KeyType = "ControlRight"
	ControlLeft        KeyType = "ControlLeft"
	ShiftUp            KeyType = "ShiftUp"
	ShiftDown          KeyType = "ShiftDown"
	ShiftRight         KeyType = "ShiftRight"
	ShiftLeft          KeyType = "ShiftLeft"
	Ignore             KeyType = "Ignore"
)

type ASCIICode struct {
	Key  KeyType
	Code []byte
}

func (c *ASCIICode) String() string {
	return string(c.Key)
}

var Codes = []*ASCIICode{
	{Key: Escape, Code: []byte{0x1b}},
	{Key: ControlSpace, Code: []byte{0x00}},
	{Key: ControlA, Code: []byte{0x1}},
	{Key: ControlB, Code: []byte{0x2}},
	{Key: ControlC, Code: []byte{0x3}},
	{Key: ControlD, Code: []byte{0x4}},
	{Key: ControlE, Code: []byte{0x5}},
	{Key: ControlF, Code: []byte{0x6}},
	{Key: ControlG, Code: []byte{0x7}},
	{Key: ControlH, Code: []byte{0x8}},
	// {Key: "ControlI", Code: []byte{0x9}},
	// {Key: "ControlJ", Code: []byte{0xa}},
	{Key: ControlK, Code: []byte{0xb}},
	{Key: ControlL, Code: []byte{0xc}},
	{Key: ControlM, Code: []byte{0xd}},
	{Key: ControlN, Code: []byte{0xe}},
	{Key: ControlO, Code: []byte{0xf}},
	{Key: ControlP, Code: []byte{0x10}},
	{Key: ControlQ, Code: []byte{0x11}},
	{Key: ControlR, Code: []byte{0x12}},
	{Key: ControlS, Code: []byte{0x13}},
	{Key: ControlT, Code: []byte{0x14}},
	{Key: ControlU, Code: []byte{0x15}},
	{Key: ControlV, Code: []byte{0x16}},
	{Key: ControlW, Code: []byte{0x17}},
	{Key: ControlX, Code: []byte{0x18}},
	{Key: ControlY, Code: []byte{0x19}},
	{Key: ControlZ, Code: []byte{0x1a}},
	{Key: ControlBackslash, Code: []byte{0x1c}},
	{Key: ControlSquareClose, Code: []byte{0x1d}},
	{Key: ControlCircumflex, Code: []byte{0x1e}},
	{Key: ControlUnderscore, Code: []byte{0x1f}},
	{Key: Backspace, Code: []byte{0x7f}},
	{Key: Up, Code: []byte{0x1b, 0x5b, 0x41}},
	{Key: Down, Code: []byte{0x1b, 0x5b, 0x42}},
	{Key: Right, Code: []byte{0x1b, 0x5b, 0x43}},
	{Key: Left, Code: []byte{0x1b, 0x5b, 0x44}},
	{Key: Home, Code: []byte{0x1b, 0x5b, 0x48}},
	{Key: Home, Code: []byte{0x1b, 0x30, 0x48}},
	{Key: End, Code: []byte{0x1b, 0x5b, 0x46}},
	{Key: End, Code: []byte{0x1b, 0x30, 0x46}},
	{Key: Enter, Code: []byte{0xa}},
	{Key: Delete, Code: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: ShiftDelete, Code: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x32, 0x7e}},
	{Key: ControlDelete, Code: []byte{0x1b, 0x5b, 0x33, 0x3b, 0x35, 0x7e}},
	{Key: Home, Code: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: End, Code: []byte{0x1b, 0x5b, 0x34, 0x7e}},
	{Key: PageUp, Code: []byte{0x1b, 0x5b, 0x35, 0x7e}},
	{Key: PageDown, Code: []byte{0x1b, 0x5b, 0x36, 0x7e}},
	{Key: Home, Code: []byte{0x1b, 0x5b, 0x37, 0x7e}},
	{Key: End, Code: []byte{0x1b, 0x5b, 0x38, 0x7e}},
	{Key: Tab, Code: []byte{0x9}},
	{Key: BackTab, Code: []byte{0x1b, 0x5b, 0x5a}},
	{Key: Insert, Code: []byte{0x1b, 0x5b, 0x32, 0x7e}},
	{Key: F1, Code: []byte{0x1b, 0x4f, 0x50}},
	{Key: F2, Code: []byte{0x1b, 0x4f, 0x51}},
	{Key: F3, Code: []byte{0x1b, 0x4f, 0x52}},
	{Key: F4, Code: []byte{0x1b, 0x4f, 0x53}},
	{Key: F1, Code: []byte{0x1b, 0x4f, 0x50, 0x41}}, // Linux console
	{Key: F2, Code: []byte{0x1b, 0x5b, 0x5b, 0x42}}, // Linux console
	{Key: F3, Code: []byte{0x1b, 0x5b, 0x5b, 0x43}}, // Linux console
	{Key: F4, Code: []byte{0x1b, 0x5b, 0x5b, 0x44}}, // Linux console
	{Key: F5, Code: []byte{0x1b, 0x5b, 0x5b, 0x45}}, // Linux console
	{Key: F1, Code: []byte{0x1b, 0x5b, 0x11, 0x7e}}, // rxvt-unicode
	{Key: F2, Code: []byte{0x1b, 0x5b, 0x12, 0x7e}}, // rxvt-unicode
	{Key: F3, Code: []byte{0x1b, 0x5b, 0x13, 0x7e}}, // rxvt-unicode
	{Key: F4, Code: []byte{0x1b, 0x5b, 0x14, 0x7e}}, // rxvt-unicode
	{Key: F5, Code: []byte{0x1b, 0x5b, 0x31, 0x35, 0x7e}},
	{Key: F6, Code: []byte{0x1b, 0x5b, 0x31, 0x37, 0x7e}},
	{Key: F7, Code: []byte{0x1b, 0x5b, 0x31, 0x38, 0x7e}},
	{Key: F8, Code: []byte{0x1b, 0x5b, 0x31, 0x39, 0x7e}},
	{Key: F9, Code: []byte{0x1b, 0x5b, 0x32, 0x30, 0x7e}},
	{Key: F10, Code: []byte{0x1b, 0x5b, 0x32, 0x31, 0x7e}},
	{Key: F11, Code: []byte{0x1b, 0x5b, 0x32, 0x32, 0x7e}},
	{Key: F12, Code: []byte{0x1b, 0x5b, 0x32, 0x34, 0x7e, 0x8}},
	// Xterm
	{Key: F13, Code: []byte{0x1b, 0x5b, 0x25, 0x7e}},
	{Key: F14, Code: []byte{0x1b, 0x5b, 0x26, 0x7e}},
	// {Key: "F15", Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},  // Conflicts with CPR response
	{Key: F15, Code: []byte{0x1b, 0x5b, 0x28, 0x7e}},
	{Key: F16, Code: []byte{0x1b, 0x5b, 0x29, 0x7e}},
	{Key: F17, Code: []byte{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: F18, Code: []byte{0x1b, 0x5b, 0x32, 0x7e}},
	{Key: F19, Code: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: F20, Code: []byte{0x1b, 0x5b, 0x34, 0x7e}},
	{Key: F13, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x50}},
	{Key: F14, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x51}},
	{Key: F16, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},
	{Key: F17, Code: []byte{0x1b, 0x5b, 0x15, 0x3b, 0x32, 0x7e}},
	{Key: F18, Code: []byte{0x1b, 0x5b, 0x17, 0x3b, 0x32, 0x7e}},
	{Key: F19, Code: []byte{0x1b, 0x5b, 0x18, 0x3b, 0x32, 0x7e}},
	{Key: F20, Code: []byte{0x1b, 0x5b, 0x19, 0x3b, 0x32, 0x7e}},
	{Key: F21, Code: []byte{0x1b, 0x5b, 0x20, 0x3b, 0x32, 0x7e}},
	{Key: F22, Code: []byte{0x1b, 0x5b, 0x21, 0x3b, 0x32, 0x7e}},
	{Key: F23, Code: []byte{0x1b, 0x5b, 0x23, 0x3b, 0x32, 0x7e}},
	{Key: F24, Code: []byte{0x1b, 0x5b, 0x24, 0x3b, 0x32, 0x7e}},
	{Key: ControlUp, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x41}},
	{Key: ControlDown, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x42}},
	{Key: ControlRight, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x43}},
	{Key: ControlLeft, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x44}},
	{Key: ShiftUp, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x41}},
	{Key: ShiftDown, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x42}},
	{Key: ShiftRight, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x43}},
	{Key: ShiftLeft, Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x44}},
	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	{Key: Up, Code: []byte{0x1b, 0x4f, 0x41}},
	{Key: Down, Code: []byte{0x1b, 0x4f, 0x42}},
	{Key: Right, Code: []byte{0x1b, 0x4f, 0x43}},
	{Key: Left, Code: []byte{0x1b, 0x4f, 0x44}},
	{Key: ControlUp, Code: []byte{0x1b, 0x5b, 0x35, 0x41}},
	{Key: ControlDown, Code: []byte{0x1b, 0x5b, 0x35, 0x42}},
	{Key: ControlRight, Code: []byte{0x1b, 0x5b, 0x35, 0x43}},
	{Key: ControlLeft, Code: []byte{0x1b, 0x5b, 0x35, 0x44}},
	{Key: ControlRight, Code: []byte{0x1b, 0x5b, 0x4f, 0x63}}, // rxvt
	{Key: ControlLeft, Code: []byte{0x1b, 0x5b, 0x4f, 0x64}},  // rxvt
	{Key: Ignore, Code: []byte{0x1b, 0x5b, 0x45}},             // Xterm
	{Key: Ignore, Code: []byte{0x1b, 0x5b, 0x46}},             // Linux console
}

func GetKey(b []byte) KeyType {
	for i, j := 0, len(Codes); i < j; i++ {
		if bytes.Equal(b, Codes[i].Code) {
			return Codes[i].Key
		}
	}
	return ""
}
