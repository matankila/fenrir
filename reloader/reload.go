package reloader

import (
	"github.com/fsnotify/fsnotify"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/watcher"
)

type Reloader struct {
	w watcher.Watcher
}

func New(c config.Configuration, filepath string) Reloader {
	reloadActionMap := watcher.ActionMap{fsnotify.Create: c.Load, fsnotify.Write: c.Load}
	return Reloader{
		w: watcher.New(filepath, reloadActionMap),
	}
}

func (r Reloader) Run() {
	go r.w.Watch()
}
