# Design Patterns

Overview of the 23 classic Gang-of-Four design patterns. Each implemented pattern lives in its own self-contained `<pattern>_test.go` file (implementation + `Test_` cases together, mirroring the style used in [`utils/`](../../utils)), demonstrating the pattern's structure with a small, testable Go example rather than an abstract UML diagram.

## Creational

Concerned with object *creation* — decoupling a client from the concrete types it needs to instantiate.

| Pattern | Idea | Status |
|---|---|---|
| Singleton | Ensure a type has exactly one instance, with global access to it (`sync.Once` in Go) | not yet added |
| Factory Method | Defer object creation to a method/subtype so the caller doesn't need to know the concrete type | not yet added |
| Abstract Factory | A factory of factories: produces families of related objects (e.g. a `UIFactory` that makes matching `Button`+`Checkbox` for a theme) | not yet added |
| Builder | Construct a complex object step-by-step via chained calls, separating construction from representation | not yet added |
| Prototype | Create new objects by cloning an existing instance rather than instantiating from scratch | not yet added |

## Structural

Concerned with how objects/types are *composed* into larger structures.

| Pattern | Idea | Status |
|---|---|---|
| Facade | A single, simplified entry point in front of a complex subsystem | [facade_test.go](facade_test.go#L81) |
| Adapter | Convert one interface into another so incompatible types can work together | [adapter_test.go](adapter_test.go#L48) |
| Decorator | Wrap an object in another implementing the same interface, adding behavior; stacks freely without a subclass explosion | [decorator_test.go](decorator_test.go#L34) |
| Proxy | A stand-in that controls access to another object (lazy loading, access control, caching, logging) | not yet added |
| Composite | Treat individual objects and groups of objects uniformly through a shared interface (e.g. file/directory trees) | [composite_test.go](composite_test.go#L40) |
| Bridge | Decouple an abstraction from its implementation so both can vary independently (e.g. shape × renderer) | not yet added |
| Flyweight | Share common/immutable state across many objects to cut memory use (e.g. glyphs in a text editor) | not yet added |

## Behavioral

Concerned with how objects *communicate* and distribute responsibility.

| Pattern | Idea | Status |
|---|---|---|
| Mediator | Centralize communication between peer objects behind one mediator, decoupling them from each other | [mediator_test.go](mediator_test.go#L53) |
| Strategy | Encapsulate interchangeable algorithms behind a common interface, swappable at runtime | not yet added |
| Observer | Subjects notify a list of registered observers on state change (pub/sub, event listeners) | not yet added |
| State | An object changes behavior when its internal state changes, by delegating to state-specific objects | not yet added |
| Command | Encapsulate a request as an object, enabling queuing, undo/redo, logging | not yet added |
| Chain of Responsibility | Pass a request along a chain of handlers until one handles it (middleware pipelines) | not yet added |
| Template Method | Define an algorithm's skeleton in a base type, letting subtypes override specific steps | not yet added |
| Iterator | Sequential access to a collection's elements without exposing its internals (Go's `range` largely covers this natively) | not yet added |
| Visitor | Add operations to a type hierarchy without modifying the types themselves | not yet added |
| Memento | Capture/restore an object's internal state for undo | not yet added |
| Interpreter | Model a grammar as an object tree | not yet added |

## Other files in this directory

`affirm_card_game.go`, `affirm_filter_match.go`, `brex_stervo_game.go` (and their `_test.go` pairs) are mock-interview-style OO design problems from specific companies — not named GoF patterns, but exercises in applying SOLID principles (single responsibility, open/closed, encapsulation) to a from-scratch design.
