package gocoa

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
// #include "window.h"
import "C"
import "unsafe"

// WindowEvent - different window delegate Events
type WindowEvent int

const (
	didResize        WindowEvent = 0
	didMove          WindowEvent = 1
	didMiniaturize   WindowEvent = 2
	didDeminiaturize WindowEvent = 3
)

// EventHandler - handler functions that accepts the updated window as parameter
type EventHandler func(wnd *Window)

// Window is just that.
type Window struct {
	title     string
	x         int
	y         int
	w         int
	h         int
	callbacks map[WindowEvent]EventHandler
	winPtr    unsafe.Pointer
}

var windows []*Window

// NewWindow constructs and returns a new window.
func NewWindow(title string, x int, y int, w int, h int) *Window {
	windowID := len(windows)
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	wnd := &Window{
		title:  title,
		x:      x,
		y:      y,
		w:      w,
		h:      h,
		callbacks: make(map[WindowEvent]EventHandler),
		winPtr: C.Window_New(C.int(windowID), C.int(x), C.int(y), C.int(w), C.int(h), cTitle)}
	windows = append(windows, wnd)
	return wnd
}

// MakeKeyAndOrderFront moves the window to the front of the screen list, within its
// level and it shows the window.
func (wnd *Window) MakeKeyAndOrderFront() {
	C.Window_MakeKeyAndOrderFront(wnd.winPtr)
}

// AddButton adds a Button to the window.
func (wnd *Window) AddButton(btn *Button) {
	C.Window_AddButton(wnd.winPtr, btn.buttonPtr)
}

// AddProgressIndicator adds a ProgressIndicator to the window.
func (wnd *Window) AddProgressIndicator(indicator *ProgressIndicator) {
	C.Window_AddProgressIndicator(wnd.winPtr, indicator.progressIndicatorPtr)
}

// Update - forces the whole window to repaint
func (wnd *Window) Update() {
	C.Window_Update(wnd.winPtr)
}

func (wnd *Window) OnDidResize(fn EventHandler) {
	wnd.callbacks[didResize] = fn
}

func (wnd *Window) OnDidMiniaturize(fn EventHandler) {
	wnd.callbacks[didMiniaturize] = fn
}

func (wnd *Window) OnDidDeminiaturize(fn EventHandler) {
	wnd.callbacks[didDeminiaturize] = fn
}

func (wnd *Window) OnDidMove(fn EventHandler) {
	wnd.callbacks[didMove] = fn
}

//export onWindowEvent
func onWindowEvent(id C.int, eventID C.int, x C.int, y C.int, w C.int, h C.int) {
	windowID := int(id)
	event := WindowEvent(eventID)
	if windowID < len(windows) && windows[windowID].callbacks[event] != nil {
		wnd := windows[windowID]
		windows[windowID].callbacks[event]( &Window{
			title:  wnd.title,
			x: int(x),
			y: int(y),
			w: int(w),
			h: int(h),
			winPtr: wnd.winPtr})
	}
}
