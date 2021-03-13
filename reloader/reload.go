package reloader

import (
	"github.com/fsnotify/fsnotify"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/watcher"
)

type Reload struct {
	w watcher.Watcher
}

func New(c config.Configuration, filepath string) Reload {
	reloadActionMap := watcher.ActionMap{fsnotify.Create: c.Load, fsnotify.Write: c.Load}
	return Reload{
		w: watcher.New(filepath, reloadActionMap),
	}
}

func (r Reload) Run() {
	go r.w.Watch()
}
