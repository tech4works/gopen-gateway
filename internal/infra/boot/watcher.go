package boot

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/fsnotify/fsnotify"
)

func NewWatcher(env string, callback func()) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if helper.IsNotNil(err) {
		return nil, err
	}
	defer func() {
		if helper.IsNotNil(err) {
			watcher.Close()
		}
	}()

	go watchEvents(watcher, callback)

	for _, path := range []string{getEnvUri(env), getJsonUri(env)} {
		err = watcher.Add(path)
		if helper.IsNotNil(err) {
			return nil, err
		}
	}

	return watcher, nil
}

func watchEvents(watcher *fsnotify.Watcher, callback func()) {
	for {
		select {
		case _, ok := <-watcher.Events:
			executeEvent(ok, callback)
		}
	}
}

func executeEvent(ok bool, callback func()) {
	if !ok {
		return
	}
	callback()
}
