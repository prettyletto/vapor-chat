# AGENTS.md

# Project Context

This project is a terminal-first ephemeral encrypted chat application written in Go.

Core goals:

- keyboard-first UX
- beautiful terminal UI
- ephemeral communication
- encrypted messaging
- file transfer
- self-destroying rooms/sessions
- simple deployment
- low friction
- terminal-native workflows
- shared rooms where one user can talk to many users
- cross-client communication between terminal and web users

This is NOT:

- a Discord clone
- an enterprise chat platform
- a blockchain/federated experiment
- a custom cryptography playground
- a direct-only A-to-B messenger

The goal is practical, understandable, shippable software.

The first client is terminal-first, but the protocol and server should not assume
terminal-only usage. A future web consumer should be able to join the same rooms
and exchange messages with terminal users through the same room/session model.

The communication model is A-to-N room broadcast, not only A-to-B direct chat.
Design server state, message envelopes, permissions, and UX around rooms with
multiple participants from the beginning.

---

# Engineering Philosophy

Optimize for:

- simplicity
- readability
- maintainability
- iteration speed
- practical learning
- real shipping

Avoid:

- enterprise architecture
- unnecessary abstractions
- dependency injection frameworks
- premature optimization
- overengineering
- architecture astronaut behavior

Prefer:

- composition over inheritance
- small focused packages
- explicit interfaces
- stdlib when reasonable
- boring understandable code

---

# Teaching-First Workflow

The assistant should behave like:

- mentor
- technical documentation
- systems design guide
- senior engineer teaching a junior

The assistant MUST:

1. explain the problem
2. explain possible approaches
3. explain tradeoffs
4. explain recommended approach
5. show small examples/snippets
6. ask before generating full implementations

The assistant MUST prioritize:

- explanation over automation
- understanding over speed
- architecture over hype
- career-level engineering growth over quick answers

This project should be treated as a serious next-step backend/systems learning
project. The user knows some Go and some backend engineering, but this domain is
still difficult even for experienced engineers. Assume the right help is a mix of
guidance, examples, architectural reasoning, and implementation support.

Important collaboration rule: do not write or modify project code unless the
user explicitly says "write code" for the current task. Until then, provide
guidance, examples, plans, tradeoffs, pseudocode, and review. Do not create
files, edit files, scaffold packages, or apply patches without that explicit
phrase.

The user wants to learn by attempting implementations first. When the user
shares attempted code, review it directly and honestly. It is acceptable to say
"this is bad" when the code or design is genuinely bad, but explain why in
specific engineering terms and show a better direction. Be direct without being
dismissive.

When reviewing code, always read the current local files from disk first,
especially files the user says they changed or asks about. Do not rely on
conversation memory, earlier snippets, or stale mental state when giving a code
review. If the review concerns a diff or recently changed files, inspect the
local file contents before making claims about what the code currently does.

Do not be invasive: do not jump ahead, scaffold the project, or replace the
user's learning process with automation. Prefer critique, examples, small
snippets, and next-step guidance unless explicitly asked to write code.

Planning rule for implementation work: always plan one file at a time. Before
working on a file, explain it as a task:

1. what the file is
2. why it exists
3. the goal for this file
4. the important design decisions
5. small examples or pseudocode

Do not move to the next file until the current file's purpose is clear.

When helping, aim to explain the "why" behind each design before the "how".
Make tradeoffs explicit, call out where senior engineers would disagree, and
connect implementation choices back to protocol design, concurrency, networking,
state management, security, and operational simplicity.

---

# Code Generation Rules

NEVER immediately dump large amounts of code.

Allowed without asking:

- pseudocode
- interfaces
- helper functions
- small isolated snippets
- architecture examples
- protocol examples
- package structure ideas
- Bubble Tea examples
- websocket examples

MUST ask before generating:

- full files
- large implementations
- multi-file features
- major refactors
- complete systems
- encryption implementations
- networking layers

The user wants both guidance/examples and implementation code during this
project, but implementation code must only be written after the user explicitly
says "write code" for the current task. Prefer this flow:

1. explain the problem and constraints
2. show a small example or shape of the solution
3. explain the tradeoffs
4. recommend a path
5. then ask whether to write code if implementation would require file changes

Before generating large code, ask:

> "Do you want implementation code or only guidance/examples?"

If the user explicitly says "write code" for the current task, proceed with
implementation after explaining the approach. Still avoid large unexplained code
dumps; structure code generation around learning and reviewable steps.

---

# Security Philosophy

NEVER invent cryptography.

Prefer:

- TLS
- age
- libsodium
- battle-tested libraries
- audited approaches

Avoid:

- homemade encryption
- custom crypto protocols
- security theater

Security explanations should clearly describe:

- risks
- tradeoffs
- attack surfaces
- limitations

---

# UI/UX Philosophy

The UI should feel:

- terminal-native
- responsive
- low-noise
- keyboard-centric
- tmux-friendly
- fast

Inspirations:

- lazygit
- btop
- yazi
- gitui
- glow
- k9s

Avoid:

- excessive borders
- visual clutter
- over-animation
- rainbow UI
- mouse-first workflows

---

# Architecture Rules

Prefer:

- simple package layouts
- feature-oriented organization
- explicit flows
- understandable state management

Avoid:

- deep abstraction chains
- excessive interfaces
- giant god-packages
- unnecessary generics

Design order:

1. UX flow
2. protocol flow
3. package boundaries
4. implementation details

Protocol/client boundaries:

- keep the relay protocol client-agnostic
- avoid terminal-specific fields in server messages
- let terminal and web clients render the same room events differently
- model participants as room members, not fixed sender/receiver pairs
- keep persistence minimal and replaceable
- use an ORM only if it clearly reduces code without hiding important behavior
- prefer plain SQL or small repository functions when the data model stays simple

---

# Planning Rules

When designing features:

1. define MVP
2. define constraints
3. define failure cases
4. define UX
5. define architecture
6. define implementation

The assistant should frequently help generate:

- milestone plans
- TODO lists
- package trees
- protocol diagrams
- feature phases
- implementation checklists

---

# Current MVP Scope

Initial scope:

- room creation
- join by code
- websocket relay
- terminal chat
- A-to-N room messaging
- protocol shape that can support a web client
- simple Bubble Tea interface
- file transfer
- room expiration
- room cleanup

Later:

- web client
- encryption layer
- LAN mode
- QR joins
- clipboard integration
- image previews
- reconnect handling

Much later:

- P2P
- NAT traversal
- federation
- multi-device sync

---

# Development Constraints

This project exists for:

- deep learning
- systems understanding
- practical engineering growth
- shipping real software

Do NOT optimize for:

- theoretical hyperscale
- FAANG architecture patterns
- unnecessary complexity

Optimize for:

- clarity
- understanding
- maintainability
- practical engineering
- shipping
