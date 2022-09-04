package gcli

import (
	"context"

	"github.com/gookit/goutil/maputil"
)

/*************************************************************
 * simple events manage
 *************************************************************/

// HookFunc definition.
//
// Returns:
//   - True go on handle. default is True
//   - False stop continue handle.
type HookFunc func(ctx *HookCtx) (stop bool)

// HookCtx struct
type HookCtx struct {
	context.Context
	maputil.Data
	App *App
	Cmd *Command

	stop bool // stop continue handle.
	err  error
	name string
}

func newHookCtx(name string, c *Command, data map[string]any) *HookCtx {
	if data == nil {
		data = make(maputil.Data)
	}

	return &HookCtx{
		name: name,
		Cmd:  c,
		Data: data,
		// with empty context
		Context: context.Background(),
	}
}

// Err of event
func (hc *HookCtx) Err() error {
	return hc.err
}

// Name of event
func (hc *HookCtx) Name() string {
	return hc.name
}

// Stopped value
func (hc *HookCtx) Stopped() bool {
	return hc.stop
}

// SetStop value
func (hc *HookCtx) SetStop(stop bool) bool {
	hc.stop = stop
	return hc.stop
}

// WithErr value
func (hc *HookCtx) WithErr(err error) *HookCtx {
	hc.err = err
	return hc
}

// WithData to ctx
func (hc *HookCtx) WithData(data map[string]any) *HookCtx {
	if data != nil {
		hc.Data = data
	}
	return hc
}

// WithApp to ctx
func (hc *HookCtx) WithApp(a *App) *HookCtx {
	hc.App = a
	return hc
}

// Hooks struct. hookManager
type Hooks struct {
	// Hooks can set some hooks func on running.
	hooks map[string]HookFunc
}

// On register event hook by name
func (h *Hooks) On(name string, handler HookFunc) {
	if handler != nil {
		if h.hooks == nil {
			h.hooks = make(map[string]HookFunc)
		}
		h.hooks[name] = handler
	}
}

// AddHook register on not exists hook.
func (h *Hooks) AddHook(name string, handler HookFunc) {
	if _, ok := h.hooks[name]; !ok {
		h.On(name, handler)
	}
}

// Fire event by name, allow with event data
func (h *Hooks) Fire(event string, ctx *HookCtx) (stop bool) {
	if handler, ok := h.hooks[event]; ok {
		return handler(ctx)
	}
	return false
}

// HasHook register
func (h *Hooks) HasHook(event string) bool {
	_, ok := h.hooks[event]
	return ok
}

// ResetHooks clear hooks
func (h *Hooks) ResetHooks() {
	h.hooks = nil
}
