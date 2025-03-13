<img align="right" width="150" src="./docs/tamago.png" alt="donburi" title="donburi" />
<h1>Tamago</h1>

<img src="https://github.com/thefishhat/tamago/actions/workflows/test.yaml/badge.svg" />[![Go Reference](https://pkg.go.dev/badge/github.com/thefishhat/tamago.svg)](https://pkg.go.dev/github.com/thefishhat/tamago)

Tamago is an editor that can be attached to
[Donburi](https://github.com/yottahmd/donburi)'s Entity
Component System
([ECS](https://en.wikipedia.org/wiki/Entity_component_system))
library for Ebitengine. It assists development by providing
a mechanism to view and edit internals of the ECS World
during runtime.

<img src="./docs/demo.gif" />

It comes with a [CLI](./cli) for users to interface with the
ECS World. The underlying [Go HTTP Client](./client)
implementation can be used to develop other customized
implementations.

## Contents

- [Contents](#contents)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Configuration](#configuration)
- [Architecture](#architecture)
- [To Do list](#to-do-list)

## Getting Started

### Installation

To add the **Editor** and **Client** to your project run:

```
go get github.com/thefishhat/tamago
```

To install the **CLI**:

```
go install github.com/thefishhat/tamago/cli@latest
```

To run the **CLI**, you can then use
`go run github.com/thefishhat/tamago/cli@latest`, or simply
`cli`.

### Usage

In order to boot up tamago, simply attach the **Editor**
from `github.com/thefishhat/tamago/editor` to the **donburi
ECS** instance.

```go
ecs := ecs.NewECS(donburi.NewWorld())
editor.Attach(ecs)
```

When running the project you should see a log, similar to
the following:

```
[server] 2010/11/12 01:23:45 Starting editor server on <SERVER_URL>
```

After the server has been started, you can
[run the CLI](#installation) to:

- navigate through entities
- inspect entity components
- explore and edit **exported** component fields

An example project can be found under
[./examples/platformer](./examples/platformer). It is
[donburi's platformer example](https://github.com/yottahmd/donburi/examples/platformer)
adapted to use the `tamago` editor.

### Configuration

You can create a `.env` file to configure the following
config:

- `SERVER_URL` - the URL (including port) where the server
  should start up. The CLI also uses the same variable to
  construct HTTP requests.

The environment variables can also be set manually if
preceded by the prefix `TAMAGO_`, e.g. `TAMAGO_SERVER_URL`.

## Architecture

The **editor** is split into 3 main components:

1. **Inspector**: periodically iterates over the donburi ECS
   world and populates the **Store**.
2. **Store**: in-memory cache of donburi internals of the
   game world.
3. **Server**: control layer that accepts HTTP traffic to
   operate on the game world, using the Store as the data
   layer.

The **CLI** is the main presentation layer - its purpose is
to query the server and format the data nicely for the
consumer.

<img src="./docs/architecture.png" />

## To Do list

- (CLI) Handle entries that are nulled between
  introspections
- (CLI) Option to clear fields (defaulting them - `""` for
  strings, `nil` for ptr, etc.)
- (CLI) Short polling / real-time communication with the
  Server
- (CLI) Loading indicator on I/O operations such as HTTP
  requests
