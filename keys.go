package main

import (
	"bytes"
	"runtime"
)

type (
	KeyType  string
	CodeType []byte
)

const (
	KeyEscape             KeyType = "Escape"
	KeyControlSpace       KeyType = "ControlSpace"
	KeyControlA           KeyType = "ControlA"
	KeyControlB           KeyType = "ControlB"
	KeyControlC           KeyType = "ControlC"
	KeyControlD           KeyType = "ControlD"
	KeyControlE           KeyType = "ControlE"
	KeyControlF           KeyType = "ControlF"
	KeyControlG           KeyType = "ControlG"
	KeyControlH           KeyType = "ControlH"
	KeyControlK           KeyType = "ControlK"
	KeyControlL           KeyType = "ControlL"
	KeyControlM           KeyType = "ControlM"
	KeyControlN           KeyType = "ControlN"
	KeyControlO           KeyType = "ControlO"
	KeyControlP           KeyType = "ControlP"
	KeyControlQ           KeyType = "ControlQ"
	KeyControlR           KeyType = "ControlR"
	KeyControlS           KeyType = "ControlS"
	KeyControlT           KeyType = "ControlT"
	KeyControlU           KeyType = "ControlU"
	KeyControlV           KeyType = "ControlV"
	KeyControlW           KeyType = "ControlW"
	KeyControlX           KeyType = "ControlX"
	KeyControlY           KeyType = "ControlY"
	KeyControlZ           KeyType = "ControlZ"
	KeyControlBackslash   KeyType = "ControlBackslash"
	KeyControlSquareClose KeyType = "ControlSquareClose"
	KeyControlCircumflex  KeyType = "ControlCircumflex"
	KeyControlUnderscore  KeyType = "ControlUnderscore"
	KeyBackspace          KeyType = "Backspace"
	KeyUp                 KeyType = "Up"
	KeyDown               KeyType = "Down"
	KeyRight              KeyType = "Right"
	KeyLeft               KeyType = "Left"
	KeyHome               KeyType = "Home"
	KeyEnd                KeyType = "End"
	KeyEnter              KeyType = "Enter"
	KeyDelete             KeyType = "Delete"
	KeyShiftDelete        KeyType = "ShiftDelete"
	KeyControlDelete      KeyType = "ControlDelete"
	KeyPageUp             KeyType = "PageUp"
	KeyPageDown           KeyType = "PageDown"
	KeyTab                KeyType = "Tab"
	KeyBackTab            KeyType = "BackTab"
	KeyInsert             KeyType = "Insert"
	KeyF1                 KeyType = "F1"
	KeyF2                 KeyType = "F2"
	KeyF3                 KeyType = "F3"
	KeyF4                 KeyType = "F4"
	KeyF5                 KeyType = "F5"
	KeyF6                 KeyType = "F6"
	KeyF7                 KeyType = "F7"
	KeyF8                 KeyType = "F8"
	KeyF9                 KeyType = "F9"
	KeyF10                KeyType = "F10"
	KeyF11                KeyType = "F11"
	KeyF12                KeyType = "F12"
	KeyF13                KeyType = "F13"
	KeyF14                KeyType = "F14"
	KeyF15                KeyType = "F15"
	KeyF16                KeyType = "F16"
	KeyF17                KeyType = "F17"
	KeyF18                KeyType = "F18"
	KeyF19                KeyType = "F19"
	KeyF20                KeyType = "F20"
	KeyF21                KeyType = "F21"
	KeyF22                KeyType = "F22"
	KeyF23                KeyType = "F23"
	KeyF24                KeyType = "F24"
	KeyControlUp          KeyType = "ControlUp"
	KeyControlDown        KeyType = "ControlDown"
	KeyControlRight       KeyType = "ControlRight"
	KeyControlLeft        KeyType = "ControlLeft"
	KeyShiftUp            KeyType = "ShiftUp"
	KeyShiftDown          KeyType = "ShiftDown"
	KeyShiftRight         KeyType = "ShiftRight"
	KeyShiftLeft          KeyType = "ShiftLeft"
	KeyIgnore             KeyType = "Ignore"
)

var (
	CodeEscape = CodeType{0x1b}

	CodeControlSpace = CodeType{0x00}
	CodeControlA     = CodeType{0x1}
	CodeControlB     = CodeType{0x2}
	CodeControlC     = CodeType{0x3}
	CodeControlD     = CodeType{0x4}
	CodeControlE     = CodeType{0x5}
	CodeControlF     = CodeType{0x6}
	CodeControlG     = CodeType{0x7}
	CodeControlH     = CodeType{0x8}
	CodeControlK     = CodeType{0xb}
	CodeControlL     = CodeType{0xc}
	CodeControlM     = CodeType{0xd}
	CodeControlN     = CodeType{0xe}
	CodeControlO     = CodeType{0xf}
	CodeControlP     = CodeType{0x10}
	CodeControlQ     = CodeType{0x11}
	CodeControlR     = CodeType{0x12}
	CodeControlS     = CodeType{0x13}
	CodeControlT     = CodeType{0x14}
	CodeControlU     = CodeType{0x15}
	CodeControlV     = CodeType{0x16}
	CodeControlW     = CodeType{0x17}
	CodeControlX     = CodeType{0x18}
	CodeControlY     = CodeType{0x19}
	CodeControlZ     = CodeType{0x1a}

	CodeControlBackslash   = CodeType{0x1c}
	CodeControlSquareClose = CodeType{0x1d}
	CodeControlCircumflex  = CodeType{0x1e}
	CodeControlUnderscore  = CodeType{0x1f}
	CodeBackspace          = CodeType{0x7f}

	CodeUp    = CodeType{0x1b, 0x5b, 0x41}
	CodeDown  = CodeType{0x1b, 0x5b, 0x42}
	CodeRight = CodeType{0x1b, 0x5b, 0x43}
	CodeLeft  = CodeType{0x1b, 0x5b, 0x44}
	CodeHome1 = CodeType{0x1b, 0x5b, 0x48}
	CodeHome2 = CodeType{0x1b, 0x30, 0x48}
	CodeEnd1  = CodeType{0x1b, 0x5b, 0x46}
	CodeEnd2  = CodeType{0x1b, 0x30, 0x46}

	CodeEnter  = CodeType{0xa}
	CodeDelete = CodeType{0x1b, 0x5b, 0x33, 0x7e}
	CodeShiftDelete = CodeType{0x1b, 0x5b, 0x33, 0x3b, 0x32, 0x7e}
	CodeControlDelete = CodeType{0x1b, 0x5b, 0x33, 0x3b, 0x35, 0x7e}
)

type ASCIICode struct {
	Key  KeyType
	Code CodeType
}

func (c *ASCIICode) String() string {
	return string(c.Key)
}

var Codes = []*ASCIICode{
	{Key: KeyEscape, Code: CodeEscape},

	{Key: KeyControlSpace, Code: CodeControlSpace},
	{Key: KeyControlA, Code: CodeControlA},
	{Key: KeyControlB, Code: CodeControlB},
	{Key: KeyControlC, Code: CodeControlC},
	{Key: KeyControlD, Code: CodeControlD},
	{Key: KeyControlE, Code: CodeControlE},
	{Key: KeyControlF, Code: CodeControlF},
	{Key: KeyControlG, Code: CodeControlG},
	{Key: KeyControlH, Code: CodeControlH},
	// {Key: "ControlI", Code: []byte{0x9}},
	// {Key: "ControlJ", Code: []byte{0xa}},
	{Key: KeyControlK, Code: CodeControlK},
	{Key: KeyControlL, Code: CodeControlL},
	{Key: KeyControlM, Code: CodeControlM},
	{Key: KeyControlN, Code: CodeControlN},
	{Key: KeyControlO, Code: CodeControlO},
	{Key: KeyControlP, Code: CodeControlP},
	{Key: KeyControlQ, Code: CodeControlQ},
	{Key: KeyControlR, Code: CodeControlR},
	{Key: KeyControlS, Code: CodeControlS},
	{Key: KeyControlT, Code: CodeControlT},
	{Key: KeyControlU, Code: CodeControlU},
	{Key: KeyControlV, Code: CodeControlV},
	{Key: KeyControlW, Code: CodeControlW},
	{Key: KeyControlX, Code: CodeControlX},
	{Key: KeyControlY, Code: CodeControlY},
	{Key: KeyControlZ, Code: CodeControlZ},

	{Key: KeyControlBackslash, Code: CodeControlBackslash},
	{Key: KeyControlSquareClose, Code: CodeControlSquareClose},
	{Key: KeyControlCircumflex, Code: CodeControlCircumflex},
	{Key: KeyControlUnderscore, Code: CodeControlUnderscore},
	{Key: KeyBackspace, Code: CodeBackspace},

	{Key: KeyUp, Code: CodeUp},
	{Key: KeyDown, Code: CodeDown},
	{Key: KeyRight, Code: CodeRight},
	{Key: KeyLeft, Code: CodeLeft},
	{Key: KeyHome, Code: CodeHome1},
	{Key: KeyHome, Code: CodeHome2},
	{Key: KeyEnd, Code: CodeEnd1},
	{Key: KeyEnd, Code: CodeEnd2},

	{Key: KeyEnter, Code: CodeEnter},
	{Key: KeyDelete, Code: CodeDelete},
	{Key: KeyShiftDelete, Code: CodeShiftDelete},
	{Key: KeyControlDelete, Code: CodeControlDelete},
	{Key: KeyHome, Code: CodeType{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: KeyEnd, Code: CodeType{0x1b, 0x5b, 0x34, 0x7e}},
	{Key: KeyPageUp, Code: CodeType{0x1b, 0x5b, 0x35, 0x7e}},
	{Key: KeyPageDown, Code: CodeType{0x1b, 0x5b, 0x36, 0x7e}},
	{Key: KeyHome, Code: CodeType{0x1b, 0x5b, 0x37, 0x7e}},
	{Key: KeyEnd, Code: CodeType{0x1b, 0x5b, 0x38, 0x7e}},
	{Key: KeyTab, Code: CodeType{0x9}},
	{Key: KeyBackTab, Code: CodeType{0x1b, 0x5b, 0x5a}},
	{Key: KeyInsert, Code: CodeType{0x1b, 0x5b, 0x32, 0x7e}},

	{Key: KeyF1, Code: CodeType{0x1b, 0x4f, 0x50}},
	{Key: KeyF2, Code: CodeType{0x1b, 0x4f, 0x51}},
	{Key: KeyF3, Code: CodeType{0x1b, 0x4f, 0x52}},
	{Key: KeyF4, Code: CodeType{0x1b, 0x4f, 0x53}},

	{Key: KeyF1, Code: CodeType{0x1b, 0x4f, 0x50, 0x41}}, // Linux console
	{Key: KeyF2, Code: CodeType{0x1b, 0x5b, 0x5b, 0x42}}, // Linux console
	{Key: KeyF3, Code: CodeType{0x1b, 0x5b, 0x5b, 0x43}}, // Linux console
	{Key: KeyF4, Code: CodeType{0x1b, 0x5b, 0x5b, 0x44}}, // Linux console
	{Key: KeyF5, Code: CodeType{0x1b, 0x5b, 0x5b, 0x45}}, // Linux console

	{Key: KeyF1, Code: CodeType{0x1b, 0x5b, 0x11, 0x7e}}, // rxvt-unicode
	{Key: KeyF2, Code: CodeType{0x1b, 0x5b, 0x12, 0x7e}}, // rxvt-unicode
	{Key: KeyF3, Code: CodeType{0x1b, 0x5b, 0x13, 0x7e}}, // rxvt-unicode
	{Key: KeyF4, Code: CodeType{0x1b, 0x5b, 0x14, 0x7e}}, // rxvt-unicode

	{Key: KeyF5, Code: CodeType{0x1b, 0x5b, 0x31, 0x35, 0x7e}},
	{Key: KeyF6, Code: CodeType{0x1b, 0x5b, 0x31, 0x37, 0x7e}},
	{Key: KeyF7, Code: CodeType{0x1b, 0x5b, 0x31, 0x38, 0x7e}},
	{Key: KeyF8, Code: CodeType{0x1b, 0x5b, 0x31, 0x39, 0x7e}},
	{Key: KeyF9, Code: CodeType{0x1b, 0x5b, 0x32, 0x30, 0x7e}},
	{Key: KeyF10, Code: CodeType{0x1b, 0x5b, 0x32, 0x31, 0x7e}},
	{Key: KeyF11, Code: CodeType{0x1b, 0x5b, 0x32, 0x32, 0x7e}},
	{Key: KeyF12, Code: CodeType{0x1b, 0x5b, 0x32, 0x34, 0x7e, 0x8}},
	// Xterm
	{Key: KeyF13, Code: CodeType{0x1b, 0x5b, 0x25, 0x7e}},
	{Key: KeyF14, Code: CodeType{0x1b, 0x5b, 0x26, 0x7e}},
	// {Key: "F15", Code: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},  // Conflicts with CPR response
	{Key: KeyF15, Code: CodeType{0x1b, 0x5b, 0x28, 0x7e}},
	{Key: KeyF16, Code: CodeType{0x1b, 0x5b, 0x29, 0x7e}},
	{Key: KeyF17, Code: CodeType{0x1b, 0x5b, 0x31, 0x7e}},
	{Key: KeyF18, Code: CodeType{0x1b, 0x5b, 0x32, 0x7e}},
	{Key: KeyF19, Code: CodeType{0x1b, 0x5b, 0x33, 0x7e}},
	{Key: KeyF20, Code: CodeType{0x1b, 0x5b, 0x34, 0x7e}},

	{Key: KeyF13, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x50}},
	{Key: KeyF14, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x51}},
	{Key: KeyF16, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x52}},
	{Key: KeyF17, Code: CodeType{0x1b, 0x5b, 0x15, 0x3b, 0x32, 0x7e}},
	{Key: KeyF18, Code: CodeType{0x1b, 0x5b, 0x17, 0x3b, 0x32, 0x7e}},
	{Key: KeyF19, Code: CodeType{0x1b, 0x5b, 0x18, 0x3b, 0x32, 0x7e}},
	{Key: KeyF20, Code: CodeType{0x1b, 0x5b, 0x19, 0x3b, 0x32, 0x7e}},
	{Key: KeyF21, Code: CodeType{0x1b, 0x5b, 0x20, 0x3b, 0x32, 0x7e}},
	{Key: KeyF22, Code: CodeType{0x1b, 0x5b, 0x21, 0x3b, 0x32, 0x7e}},
	{Key: KeyF23, Code: CodeType{0x1b, 0x5b, 0x23, 0x3b, 0x32, 0x7e}},
	{Key: KeyF24, Code: CodeType{0x1b, 0x5b, 0x24, 0x3b, 0x32, 0x7e}},

	{Key: KeyControlUp, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x41}},
	{Key: KeyControlDown, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x42}},
	{Key: KeyControlRight, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x43}},
	{Key: KeyControlLeft, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x35, 0x44}},

	{Key: KeyShiftUp, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x41}},
	{Key: KeyShiftDown, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x42}},
	{Key: KeyShiftRight, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x43}},
	{Key: KeyShiftLeft, Code: CodeType{0x1b, 0x5b, 0x31, 0x3b, 0x32, 0x44}},

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	{Key: KeyUp, Code: CodeType{0x1b, 0x4f, 0x41}},
	{Key: KeyDown, Code: CodeType{0x1b, 0x4f, 0x42}},
	{Key: KeyRight, Code: CodeType{0x1b, 0x4f, 0x43}},
	{Key: KeyLeft, Code: CodeType{0x1b, 0x4f, 0x44}},

	{Key: KeyControlUp, Code: CodeType{0x1b, 0x5b, 0x35, 0x41}},
	{Key: KeyControlDown, Code: CodeType{0x1b, 0x5b, 0x35, 0x42}},
	{Key: KeyControlRight, Code: CodeType{0x1b, 0x5b, 0x35, 0x43}},
	{Key: KeyControlLeft, Code: CodeType{0x1b, 0x5b, 0x35, 0x44}},

	{Key: KeyControlRight, Code: CodeType{0x1b, 0x5b, 0x4f, 0x63}}, // rxvt
	{Key: KeyControlLeft, Code: CodeType{0x1b, 0x5b, 0x4f, 0x64}},  // rxvt

	{Key: KeyIgnore, Code: CodeType{0x1b, 0x5b, 0x45}}, // Xterm
	{Key: KeyIgnore, Code: CodeType{0x1b, 0x5b, 0x46}}, // Linux console
}

func GetKey(code CodeType) KeyType {
	for i, j := 0, len(Codes); i < j; i++ {
		if bytes.Equal(code, Codes[i].Code) {
			return Codes[i].Key
		}
	}
	return ""
}

func GetCode(key KeyType) CodeType {
	switch key {
	case KeyEnter:
		switch runtime.GOOS {
		case "darwin":
			return CodeControlM
		default:
			return CodeEnter
		}
	}
	return nil
}
