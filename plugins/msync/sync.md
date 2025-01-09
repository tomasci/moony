sync plugin supposed to track all connected clients 
and check if they're still alive

to do that this plugin will use Redis
which is not added to the server yet

and also, ping plugin must update sync with events:
because all users will be authorized, ping plugin will receive requests
with token; and sync plugin will subscribe to ping plugin events 
and update specific user data

ping data will be saved to redis for 10 seconds
if not updated - removed - client disconnected 

