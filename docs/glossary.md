# Glossary

This document defines the shared language of the project.

Use these terms consistently in documentation, TypeSpec, API models, frontend, backend, and tests.

## Core Terms

### Host

The calendar owner role.

The host manages event types, publishes available time, and views upcoming bookings.

For MVP:

* there is only one predefined host
* there is no authentication
* host pages are public
* `Host` does not have to be a persistent user entity

### Guest

The person who books a meeting.

The guest opens the booking page, selects an event type, selects an available slot, enters their name and email, and confirms the booking.

For MVP:

* guest identity is represented only by name and email in a booking
* `Guest` does not have to be a persistent user entity

### Event Type

A type of meeting that can be booked.

Example:

```text
Intro Call
30-minute introduction call
Duration: 30 minutes
```

An event type contains:

* id
* title
* description
* duration

For MVP:

* event types are created by the host
* event type duration is used to produce bookable slots
* in MVP, event type duration is fixed at 30 minutes
* guests can view available event types
* guests choose one event type before selecting a slot

### Availability

A time range published by the host as available for meetings.

Example:

```text
2026-06-15T10:00:00Z — 2026-06-15T12:00:00Z
```

Availability is not a booking.
Availability is the source from which available slots are derived.

For MVP:

* availability is created by the host
* availability is split into bookable slots
* recurring availability is out of scope
* editing and deleting availability are out of scope unless added later through documentation

### Slot

A 30-minute time interval derived from availability.

Example:

```text
2026-06-15T10:00:00Z — 2026-06-15T10:30:00Z
```

A slot can be:

* `available` — can be booked
* `booked` — already has a booking

Do not use a separate `unavailable` slot state in the MVP unless the domain model is updated.

### Booking

A confirmed booking of a slot by a guest.

A booking contains:

* event type
* selected slot start time
* selected slot end time
* guest name
* guest email
* creation time

A booking makes the selected time unavailable for future bookings.

Important MVP rule:

* two bookings cannot exist for the same time, even for different event types

### Meeting

The real-world call represented by a booking.

For MVP:

* `Booking` is the main domain object
* `Meeting` must not become a separate API model unless the domain model is updated

## Preferred Terms

Use these terms:

| Use          | Do not use                              |
| ------------ | --------------------------------------- |
| Host         | Owner, User, Organizer                  |
| Guest        | Customer, Visitor, Invitee              |
| Event Type   | Event, Meeting Type, Appointment Type   |
| Availability | Schedule, WorkingHours, Calendar        |
| Slot         | Event, Appointment, TimeCell            |
| Booking      | Reservation, Appointment, CalendarEvent |
| Meeting      | Event, CallEntity                       |

## Domain Relationships

```text
Host creates Event Type
Host publishes Availability
Guest views Event Types
Guest selects Event Type
Availability produces 30-minute Slots
Guest selects Slot
Guest creates Booking
Booking makes selected time unavailable
Host views Bookings
```

## Naming Rules

* Use `Host` for the calendar owner role.
* Use `Guest` for the person booking the meeting.
* Use `Event Type` for the bookable meeting type.
* Use `Availability` only for a host-published available time range.
* Use `Slot` only for a 30-minute bookable time interval.
* Use `Booking` for a confirmed booking.
* Use `Meeting` only when describing the real-world call.
* Do not introduce new domain terms without updating this file.

## Time Terms

### Past Slot

A slot whose start time is earlier than the current server time.

Past slots must not be bookable.

### Upcoming Booking

A booking whose slot start time is equal to or later than the current server time.

The host bookings page shows upcoming bookings only.

### Booking Horizon

The time window in which guests can see available slots.

For MVP:

* guests see available slots for the next 14 days

## TBD

* Whether configurable event type duration is needed after the first MVP
* Exact date/time display format in UI
* Whether availability editing is needed after the first MVP
