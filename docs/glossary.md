# Glossary

This document defines the shared language of the project.

Use these terms consistently in documentation, TypeSpec, API models, frontend, backend, and tests.

## Core Terms

### Host

The calendar owner.

The host defines available time and views upcoming bookings.

For MVP, there is only one host and no authentication.

### Guest

The person who books a meeting.

The guest opens the booking page, selects an available slot, enters their name and email, and confirms the booking.

### Availability

A time range published by the host as available for meetings.

Example:

```text
2026-06-15 10:00–12:00
```

Availability is not a booking.
It is the source from which available slots are derived.

### Slot

A 30-minute time interval that can be booked by a guest.

Example:

```text
2026-06-15 10:00–10:30
```

A slot can be:

* available
* booked
* unavailable

### Booking

A confirmed reservation of a slot by a guest.

A booking contains:

* selected slot
* guest name
* guest email
* creation time

### Meeting

The scheduled call represented by a booking.

For MVP, `Booking` is the main domain object.
Use `Meeting` only when describing the real-world event.

## Preferred Terms

Use these terms:

| Use          | Do not use                              |
| ------------ | --------------------------------------- |
| Host         | Owner, User, Organizer                  |
| Guest        | Customer, Visitor, Invitee              |
| Slot         | Event, Appointment, TimeCell            |
| Booking      | Reservation, Appointment, CalendarEvent |
| Availability | Schedule, WorkingHours, Calendar        |

## Domain Relationships

```text
Host publishes Availability
Availability produces Slots
Guest selects Slot
Guest creates Booking
Booking makes Slot unavailable
Host views Bookings
```

## Naming Rules

* Use `Host` for the calendar owner.
* Use `Guest` for the person booking the meeting.
* Use `Slot` only for a 30-minute bookable time interval.
* Use `Availability` only for the host's available time range.
* Use `Booking` for a confirmed reservation.
* Do not introduce new domain terms without updating this file.

## TBD

* Exact time zone rule
* Date/time display format
* Whether availability is preconfigured or editable through UI in MVP
