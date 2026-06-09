# Vapor Chat

Vapor Chat is an ephemeral room-based chat system for real-time communication across terminal and web clients. The system exists to support short-lived conversations and file sharing without durable user or message history.

## Language

**Room**:
A temporary shared conversation space where participants exchange live events.
_Avoid_: Channel, thread, lobby

**Participant**:
A person currently present in a room under a temporary display name.
_Avoid_: User, account, member

**Creator**:
The participant who creates a room and receives the single extra authority to end it early.
_Avoid_: Admin, owner, host

**Join Code**:
The room's user-facing access secret used to join an active room.
_Avoid_: Room ID, invite link, pin

**Session**:
The server-recognized presence of one participant inside one active room.
_Avoid_: Account session, login

**Session Token**:
An opaque server-issued credential that proves a participant's session when opening or using a live room connection.
_Avoid_: API key, password, cookie

**Room Event**:
A fact that happened in a room and can be observed by connected participants, such as a join, leave, message, file share, or room end.
_Avoid_: Packet, frame, command

**TTL Preset**:
A fixed room lifetime option chosen at room creation.
_Avoid_: Custom duration, timeout value

**Expiration**:
The moment an active room ceases to exist because its lifetime has run out.
_Avoid_: Archive, retention

**Active Room**:
A room that currently exists and can still accept participants or events.
_Avoid_: Saved room, persistent room

**File Share**:
A room event that makes a file available to participants for explicit fetch during the life of the room.
_Avoid_: Attachment history, upload archive
