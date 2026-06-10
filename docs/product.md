# Product

## Overview

This project is a small meeting scheduling web application inspired by Cal.com.

The goal is to build a complete MVP using a Design First and API-first workflow.

The application allows a host to publish available time slots and allows a guest to book a 30-minute meeting.

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

### Guest

The guest books a meeting.

The guest wants to:

* see available slots
* choose a convenient time
* enter contact details
* confirm the booking

## MVP Scope

The MVP includes:

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
* complex time zone management

## Main Guest Scenario

1. Guest opens the booking page.
2. Guest sees available 30-minute slots.
3. Guest selects one available slot.
4. Guest enters name and email.
5. Guest confirms the booking.
6. System creates the booking.
7. The selected slot becomes unavailable.
8. Guest sees booking confirmation.

## Main Host Scenario

1. Host opens the bookings page.
2. Host sees upcoming bookings.
3. Each booking shows the slot time, guest name, and guest email.

## Business Rules

* Slot duration is 30 minutes.
* A slot can have only one booking.
* A booked slot must not be shown as available.
* Guest name is required.
* Guest email is required.
* Guest email must be valid enough for MVP validation.
* Booking must fail if the selected slot is already booked.
* Past slots must not be bookable.
* Host bookings page shows upcoming bookings only.

## Pages

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

* guest can see available slots
* guest can book a slot
* booked slot disappears from available slots
* host can see the booking
* double booking is prevented
* application can be run through Docker
* core scenarios are covered by tests

## TBD

* Database choice
* Exact time zone rule
* Initial availability setup
* Final deployment target
* Testing stack
