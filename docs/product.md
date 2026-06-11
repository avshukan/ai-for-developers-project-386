# Product

## Overview

This project is a small meeting scheduling web application inspired by Cal.com.

The goal is to build a complete MVP using a Design First and API-first workflow.

The application allows a host to publish available time ranges and allows a guest to book a 30-minute meeting slot.

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

The host wants to:

* publish available time for meetings
* see upcoming bookings

For MVP, the host is not authenticated. Host pages are public.
This is an intentional learning simplification and not a production security model.

### Guest

The guest books a meeting.

The guest wants to:

* see available slots
* choose a convenient time
* enter contact details
* confirm the booking

## MVP Scope

The MVP includes:

* creating a host availability range
* deriving 30-minute slots from availability
* viewing available 30-minute slots
* booking an available slot
* collecting guest name
* collecting guest email
* preventing double booking
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

## Main Host Availability Scenario

1. Host opens the availability page.
2. Host enters an availability start time.
3. Host enters an availability end time.
4. Host submits the availability range.
5. System saves the availability range.
6. System derives 30-minute slots from the availability range.
7. Guest can see derived available slots on the booking page.

## Main Guest Scenario

1. Guest opens the booking page.
2. Guest sees available 30-minute slots.
3. Guest selects one available slot.
4. Guest enters name and email.
5. Guest confirms the booking.
6. System creates the booking.
7. The selected slot becomes unavailable.
8. Guest sees booking confirmation.

## Main Host Bookings Scenario

1. Host opens the bookings page.
2. Host sees upcoming bookings.
3. Each booking shows the slot time, guest name, and guest email.

## Business Rules

### Availability Rules

* Availability has a start time and an end time.
* Availability end time must be after start time.
* Availability is split into 30-minute slots.
* Only full 30-minute slots are bookable.
* Recurring availability is out of scope.
* Editing and deleting availability are out of scope for MVP.

### Slot Rules

* Slot duration is 30 minutes.
* A slot is derived from availability.
* A slot can have only one booking.
* A booked slot must not be shown as available.
* Past slots must not be bookable.

### Booking Rules

* Guest name is required.
* Guest email is required.
* Guest email must be valid enough for MVP validation.
* Booking must fail if the selected slot is already booked.
* Booking must fail if the selected slot is in the past.
* Double booking must be prevented on the backend.

### Time Rules

* API date-time values use ISO 8601.
* API date-time values are represented in UTC.
* The server is the source of truth for current time.
* A past slot is a slot whose start time is earlier than the current server time.
* An upcoming booking is a booking whose slot start time is equal to or later than the current server time.
* UI date/time display format is `TBD`.

## Pages

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

* available slots
* selected slot
* guest name field
* guest email field
* booking confirmation state
* validation errors

### Host Bookings Page

Used by the host.

Shows:

* list of upcoming bookings
* slot date and time
* guest name
* guest email

Authentication is not required for MVP.

## Success Criteria

The MVP is complete when:

* host can create an availability range
* guest can see available slots derived from availability
* guest can book a slot
* booked slot disappears from available slots
* host can see the booking
* double booking is prevented
* past slots are not bookable
* application can be run through Docker
* core scenarios are covered by tests

## TBD

* Database choice
* Exact UI date/time display format
* Final deployment target
* Testing stack
* Docker topology
