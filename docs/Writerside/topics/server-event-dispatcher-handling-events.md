# Handling events

As with event dispatching, you can process any event because the EventType is just a string.

```Go
dispatcher.RegisterEventHandler(“users.OnUserLogin”, func(ctx context.Context, data interface{}) {
    // process and work with data here, call other methods in your plugin, etc.
    // replace the interface{} with some type
})
```