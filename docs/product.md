# Product

## Overview

This project is a small meeting scheduling web application inspired by Cal.com.

The goal is to build a complete MVP using a Design First and API-first workflow.

The application allows a host to create event types, publish available time, and allows a guest to book a free slot.

## Product Goal

Build a small but complete application that demonstrates:

* clear product scope
* documented domain model
* TypeSpec API contract
* independent frontend and backend implementation
* tests for core user scenarios
* Docker-based delivery

## Users

### Host

The host owns the calendar.

For MVP:

* there is one predefined host
* there is no authentication
* host pages are public
* this is an intentional learning simplification, not a production security model

The host wants to:

* create event types
* publish available time
* see upcoming bookings across all event types

### Guest

The guest books a meeting without an account.

The guest wants to:

* see available event types
* choose an event type
* see free slots for the next 14 days
* choose a convenient slot
* enter contact details
* confirm the booking

## MVP Scope

The MVP includes:

* creating an event type
* viewing available event types
* creating a host availability range
* deriving slots from availability and event type duration
* showing free slots for the next 14 days
* booking an available slot
* collecting guest name
* collecting guest email
* preventing double booking across all event types
* viewing upcoming bookings as host

## Out of Scope

The MVP does not include:

* authentication
* user accounts
* multiple hosts
* teams
* external calendar integrations
* Google Calendar
* Outlook Calendar
* Zoom or video call integrations
* payments
* email notifications
* reminders
* cancellation
* rescheduling
* recurring events
* availability editing
* availability deletion
* complex time zone management

## Main Host Event Type Scenario

1. Host opens the event types page.
2. Host enters event type title.
3. Host enters event type description.
4. Host enters event type duration.
5. Host submits the event type.
6. System creates the event type.
7. Guest can see the event type on the public booking page.

## Main Host Availability Scenario

1. Host opens the availability page.
2. Host enters an availability start time.
3. Host enters an availability end time.
4. Host submits the availability range.
5. System saves the availability range.
6. System derives bookable slots from availability and event type duration.
7. Guest can see derived free slots on the booking page.

## Main Guest Scenario

1. Guest opens the booking page.
2. Guest sees available event types.
3. Guest selects an event type.
4. Guest sees free slots for the next 14 days.
5. Guest selects one free slot.
6. Guest enters name and email.
7. Guest confirms the booking.
8. System creates the booking.
9. The selected time becomes unavailable for all event types.
10. Guest sees booking confirmation.

## Main Host Bookings Scenario

1. Host opens the bookings page.
2. Host sees upcoming bookings across all event types.
3. Each booking shows event type, slot time, guest name, and guest email.

## Business Rules

### Event Type Rules

* Event type has id, title, description, and duration.
* Event type is created by the host.
* Guest can view available event types.
* Guest must choose an event type before selecting a slot.
* Event type duration is used to derive bookable slots.
* Event type duration must be positive.
* The default expected duration is 30 minutes.

### Availability Rules

* Availability has a start time and an end time.
* Availability is created by the host.
* Availability end time must be after start time.
* Availability is used to derive bookable slots.
* Recurring availability is out of scope.
* Editing and deleting availability are out of scope for MVP.

### Slot Rules

* Slots are derived from availability and event type duration.
* Guests see free slots for the next 14 days.
* A booked time must not be shown as available.
* Past slots must not be bookable.

### Booking Rules

* Guest name is required.
* Guest email is required.
* Guest email must be valid enough for MVP validation.
* Booking must reference an event type.
* Booking must reference a free slot.
* Booking must fail if the selected slot is already booked.
* Booking must fail if the selected slot is in the past.
* Double booking must be prevented on the backend.
* Two bookings cannot exist for the same time, even for different event types.

### Time Rules

* API date-time values use ISO 8601.
* API date-time values are represented in UTC.
* The server is the source of truth for current time.
* A past slot is a slot whose start time is earlier than the current server time.
* An upcoming booking is a booking whose slot start time is equal to or later than the current server time.
* UI date/time display format is `TBD`.

## Pages

### Host Event Types Page

Used by the host.

Shows:

* event type title field
* event type description field
* event type duration field
* submit button
* validation errors
* list of created event types

Authentication is not required for MVP.

### Host Availability Page

Used by the host.

Shows:

* availability start time field
* availability end time field
* submit button
* validation errors
* success state

Authentication is not required for MVP.

### Booking Page

Used by the guest.

Shows:

* available event types
* selected event type
* free slots for the next 14 days
* selected slot
* guest name field
* guest email field
* booking confirmation state
* validation errors

### Host Bookings Page

Used by the host.

Shows:

* list of upcoming bookings
* event type
* slot date and time
* guest name
* guest email

Authentication is not required for MVP.

## Success Criteria

The MVP is complete when:

* host can create an event type
* host can create an availability range
* guest can see available event types
* guest can choose an event type
* guest can see free slots for the next 14 days
* guest can book a slot
* booked time disappears from available slots
* host can see the booking
* double booking across all event types is prevented
* past slots are not bookable
* application can be run through Docker
* core scenarios are covered by tests

## TBD

* Database choice
* Exact UI date/time display format
* Final deployment target
* Testing stack
* Docker topology
* Whether event type duration can be changed after creation
