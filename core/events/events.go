package events

/*
# events package

events package contains all events used in application.

in theory. all interaction should use events package.

the events are dispatched using eventbus package.

Here are some major events

- cmd: call cmd
- reply: call reply
- update: information updating event. usually issued by internal controller and broadcast to all channel


naming convention

- cmd: 'cmd.event.id.'
- reply: 'reply.same.same.cmd.id'
- update: 'update.event.id'

*/
