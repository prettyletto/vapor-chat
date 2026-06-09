# Shared room-event protocol across clients

Vapor Chat will define a shared room-event protocol from day one and use it across the server and terminal client, with the future web client expected to use the same contract. We chose this because the product is room-centric rather than client-centric, and separate per-client protocol shapes would make the server terminal-shaped now and expensive to correct later.
