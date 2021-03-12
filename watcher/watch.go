package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/matankila/fenrir/logger"
)

type ActionOnEvent func(filePath string) error
type ActionMap map[fsnotify.Op]ActionOnEvent

type Watcher struct {
	FilePath  string
	actionMap ActionMap
}

func New(filePath string, actionMap ActionMap) Watcher {
	return Watcher{
		FilePath:  filePath,
		actionMap: actionMap,
	}
}

func (w Watcher) Watch() {
	l := logger.GetLogger(logger.Watcher)
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Write:
					if err := (w.actionMap[fsnotify.Write])(w.FilePath); err != nil {
						l.Error(err.Error())
					} else {
						l.Info("file write detected...")
					}
				case fsnotify.Create:
					if err := (w.actionMap[fsnotify.Create])(w.FilePath); err != nil {
						l.Error(err.Error())
					} else {
						l.Info("file creation detected...")
					}
				}
			// watch for errors
			case err := <-watcher.Errors:
				panic(err)
			}
		}
	}()

	if err := watcher.Add(w.FilePath); err != nil {
		panic(err)
	}
	<-done
}
