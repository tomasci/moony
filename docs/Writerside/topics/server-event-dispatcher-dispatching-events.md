# Dispatching events

By default, `EventType` is set to `string`.
This means that you can create and handle any events you want.

Let's say you have created a Users plugin with `login` and `create` methods.

Now you can notify all the other plugins on the server that something has happened. Here's how you can do that:

```Go
dispatcher.Dispatch(“users.OnUserLogin”, ctx, {uid, username})
```

> This is not the actual code, I will replace it when the first plugins are ready.

It is important to add an event namespace like `users.` so that other plugins can distinguish events from different plugins. Of course, the namespace must be unique.

> Be careful with the data you share through events.
> Do not pass sensitive data such as passwords, keys, hashes, etc. {style="warning"}

> But also don't limit the amount of data too much, share enough data so that other plugins can work with it. 
> For example, if your plugin allows to move objects, share information about the object and coordinates (old and new). {style="note"}

> I hope to solve the namespace issue later and make it automatic. {style="warning"}