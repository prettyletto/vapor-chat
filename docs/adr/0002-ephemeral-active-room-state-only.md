# Ephemeral active-room state only

Phase 1 keeps only the minimum active room, participant, session, and file state required to run live rooms, and does not persist durable users, room history, or message history. We chose this because ephemerality is a core product property, and durable storage would change the product semantics and complicate lifecycle, privacy, and cleanup behavior too early.
