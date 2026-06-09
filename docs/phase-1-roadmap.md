# Vapor Chat Phase 1 Roadmap

## Product Definition

Phase 1 proves that ephemeral room-based chat works reliably across terminal and web clients, with real-time text messaging as the core behavior. File sharing and room expiration are included, but they are secondary to the room event model.

## Locked Decisions

- No durable users, rooms, messages, or chat history.
- The server keeps only the minimum active state required to run a live room.
- Late joiners do not receive prior messages.
- Rooms are small-group rooms capped at 8 participants.
- The room TTL starts at creation time.
- A room ends when it becomes empty, when its TTL expires, or when the creator ends it early.
- The creator has exactly one special power: end the room early for everyone.
- The protocol is room-event based. Text is one event type, not a separate system.
- Participants identify themselves with temporary per-room display names.
- Display names must be unique within a room.
- Display name uniqueness is case-insensitive.
- Display names are trimmed before validation and uniqueness checks.
- Display names must be 3 to 32 characters long after trimming.
- Display names are restricted to alphanumeric characters in Phase 1.
- A message counts as sent when the server accepts it into the room event stream.
- File sharing is part of the room event model.
- The room event stream exposes high-level file events, not raw transfer chunks.
- File recipients choose to fetch a shared file; files are not auto-pushed.
- The join code is the room access secret.
- Join codes must be hard to guess and only need to be unique among active rooms.
- Phase 1 join codes are randomly generated 10-character uppercase alphanumeric values.
- A join code may be reused later for a different room after the original room is gone.
- Room creation is request/response first; the live room stream comes after creation.
- Creating a room automatically joins the creator as the first participant.
- The server issues the creator a session token during room creation.

## Phase 1 CreateRoom Contract

### Request

- Creator display name
- TTL preset

Allowed TTL presets:

- 15 minutes
- 30 minutes
- 1 hour
- 2 hours

### Response

- Join code
- TTL preset
- Expiration timestamp
- Creator session token

The browser join URL is built by the client or web layer from the join code rather than returned by the core room creation operation.

### Validation

- Display name must be present after trimming.
- Display name must be 3 to 32 characters long after trimming.
- Display name must be alphanumeric.
- TTL must be one of the allowed presets.

### Success Invariants

- An active room exists.
- The room expiration is derived from creation time.
- The creator is already the first participant in the room.
- The room has exactly one participant immediately after creation.
- The creator has a valid server-issued session token scoped to that room.
- The returned join code identifies that active room for its lifetime.
- If an active join code collision occurs during creation, the service retries with a new code up to 5 times before failing.

## Recommended Initial Package Direction

Top-level boundaries:

- `internal/server`: server-side room lifecycle and session behavior
- `internal/protocol`: shared room event and request/response shapes
- client packages later, built against the shared protocol

First server slice target:

- `internal/server/room/types.go`
- `internal/server/room/store.go`
- `internal/server/room/service.go`

Responsibility notes:

- `types.go`: minimum room, participant, session, join code, TTL, and expiration concepts for room creation
- `store.go`: minimal contract for storing active room state atomically
- `service.go`: `CreateRoom` validation, code/token generation, expiration calculation, and atomic room creation

## Implementation Slices

### Slice 1: Create Room Foundation

Goal:
The creator can create a room with a display name and TTL preset and receive the join code, expiration timestamp, and session token.

Done when:
- `CreateRoom` accepts valid input and rejects invalid input.
- The server creates one active room with one participant.
- The store enforces active join code uniqueness.
- The creator session token is issued by the server.

Tasks:
- Define `RoomCode`, `TTLPreset`, and `SessionToken` as named types.
- Define the minimum room, participant, and session shapes.
- Define the active-room store contract.
- Implement `CreateRoom` as one atomic application operation.
- Add tests for valid create, empty-after-trim name, short name, long name, non-alphanumeric name, invalid TTL, and join code collision retry behavior.

### Slice 2: Join Room

Goal:
Another participant can join an active room using the join code and a temporary display name.

Done when:
- Join succeeds for a valid active room.
- Duplicate display names are rejected.
- Full rooms are rejected.
- The joining participant receives a session token.

Tasks:
- Define the `JoinRoom` request and response.
- Add name uniqueness checks at join time.
- Enforce the 8-participant cap.
- Add tests for success, duplicate names, expired rooms, and full rooms.

### Slice 3: Live Room Connection

Goal:
An authenticated session can attach to the room's live event stream.

Done when:
- A client can open a live connection using its session token.
- Invalid or stale session tokens are rejected.
- Connected participants receive room lifecycle events.

Tasks:
- Define how live room authentication uses the session token.
- Emit room join and leave events through the shared protocol.
- Add tests for valid attach, invalid token, and stale token behavior.

### Slice 4: Real-Time Text Messaging

Goal:
Participants can send plain text messages into the room event stream.

Done when:
- The server accepts or rejects a text message explicitly.
- Accepted messages are broadcast as room events.
- No per-recipient delivery tracking is introduced.

Tasks:
- Define the text message event shape.
- Validate message input.
- Add tests for acceptance, rejection, and broadcast behavior.

### Slice 5: Leave and Empty-Room Cleanup

Goal:
Participant departure is handled correctly and empty rooms are deleted immediately.

Done when:
- Explicit leave and disconnect both remove the session.
- Leave events are broadcast.
- Empty rooms are destroyed immediately.

Tasks:
- Define leave behavior.
- Invalidate the departed session.
- Add tests for last-participant cleanup.

### Slice 6: TTL Expiration

Goal:
Rooms end automatically when their TTL runs out.

Done when:
- Expired rooms stop accepting joins and events.
- Connected participants receive a room-expired event.
- Active room and session state are destroyed.

Tasks:
- Define expiration enforcement.
- Add tests for expiry, stale session rejection, and cleanup.

### Slice 7: Creator Ends Room Early

Goal:
The creator can end the room before TTL expiry.

Done when:
- Only the creator can trigger early end.
- All connected participants receive a room-ended event.
- Room and session state are invalidated.

Tasks:
- Define creator authorization for early end.
- Add tests for creator success and non-creator rejection.

### Slice 8: Browser Client on the Same Contract

Goal:
The browser uses the same room creation, join, and event semantics as the terminal client.

Done when:
- Browser and terminal clients can join the same room.
- The browser consumes the shared protocol without browser-only room semantics.

Tasks:
- Validate that protocol fields are client-agnostic.
- Add one end-to-end scenario covering terminal and browser participation in the same room.

### Slice 9: File Share Metadata Flow

Goal:
Participants can announce a file share through room events.

Done when:
- A file share event includes metadata and lifecycle state.
- Other participants see the file share and may choose to fetch it.
- No raw chunk events appear in the shared room event stream.

Tasks:
- Define file share event types and states.
- Add tests for announce, visibility, and explicit fetch initiation.

### Slice 10: File Transfer Bytes and Cleanup

Goal:
File bytes can be transferred within the life of the room and disappear with the room.

Done when:
- File transfer is authorized by room/session state.
- Failed transfers surface clear file lifecycle events.
- Files are deleted when the room ends or expires.

Tasks:
- Define transfer authorization and lifetime rules.
- Add tests for successful fetch, failed fetch, and cleanup on room end.

## Immediate Next Step

Implement Slice 1 only.

Do not start with join, live streaming, browser behavior, or file transfer until `CreateRoom` is stable, tested, and its invariants are explicit in code.
