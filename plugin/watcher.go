package plugin

type PluginWatcher interface {
	Watch() (Plugin, error)
}

type pluginWatcher struct {
	root string
}

func NewPluginWatcher(root string) PluginWatcher {
	return &pluginWatcher{
		root: root,
	}
}

func (w *pluginWatcher) Watch() (plugin Plugin, err error) {
	// watch directory events(CREATE, DELETE, UPDATE)
	// check namespace conflict

	return
}
