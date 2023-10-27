# TODO

## MQTT
- [x] Create a way to make a new roasting session
- [x] Filter already existing session
- [x] How to exit the session
- [x] in Postgres: how to unset the session (just delete the row with the matching id)

# TODO

## SQL INSERT
- [x] For every new session, generate the variable for roast_sessions_id. To obtain this, first insert roaster_id and roast_date to roast_sessions.
- [x] Use roast_sessions_id from the above query and referenced it as session_id when updating session_measurements from mqtt callbacks.

## SQL TABLE REVISION
- [x] In roast session, there should be a user_id column

# TODO
- [ ] CONSIDERATION: should MQTT cb add timestamp in epoch automatically or became part or the mqtt message?
