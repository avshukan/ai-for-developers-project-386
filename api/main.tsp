import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{
  title: "Call Booking API",
})
@server("http://localhost:8080", "Local development server")
namespace CallBooking {
  @doc("A bookable meeting type created by the host.")
  model EventType {
    @doc("System-generated event type id.")
    id: string;

    @minLength(1)
    title: string;

    @minLength(1)
    description: string;

    @doc("Fixed to 30 minutes in the MVP.")
    @minValue(30)
    @maxValue(30)
    durationMinutes: int32;
  }

  @doc("Request body for creating an event type.")
  model CreateEventTypeRequest {
    @minLength(1)
    title: string;

    @minLength(1)
    description: string;

    @doc("Must be 30 in the MVP.")
    @minValue(30)
    @maxValue(30)
    durationMinutes: int32;
  }

  @doc("A time range published by the host.")
  model Availability {
    @doc("System-generated availability id.")
    id: string;

    @doc("Availability start time in UTC.")
    startAt: utcDateTime;

    @doc("Availability end time in UTC.")
    endAt: utcDateTime;

    @doc("Creation time in UTC.")
    createdAt: utcDateTime;
  }

  @doc("Request body for creating an availability range.")
  model CreateAvailabilityRequest {
    @doc("Availability start time in UTC.")
    startAt: utcDateTime;

    @doc("Availability end time in UTC.")
    endAt: utcDateTime;
  }

  enum SlotStatus {
    available: "available",
    booked: "booked",
  }

  @doc("A 30-minute slot derived from host availability.")
  model Slot {
    @doc("Slot start time in UTC.")
    startAt: utcDateTime;

    @doc("Slot end time in UTC.")
    endAt: utcDateTime;

    status: SlotStatus;
  }

  @doc("A confirmed booking created by a guest.")
  model Booking {
    @doc("System-generated booking id.")
    id: string;

    eventTypeId: string;

    @doc("Booking start time in UTC.")
    startAt: utcDateTime;

    @doc("Booking end time in UTC.")
    endAt: utcDateTime;

    @minLength(1)
    guestName: string;

    @format("email")
    guestEmail: string;

    @doc("Creation time in UTC.")
    createdAt: utcDateTime;
  }

  @doc("Request body for creating a booking.")
  model CreateBookingRequest {
    eventTypeId: string;

    @doc("Selected slot start time in UTC. The backend calculates endAt as startAt plus 30 minutes.")
    startAt: utcDateTime;

    @minLength(1)
    guestName: string;

    @format("email")
    guestEmail: string;
  }

  model EventTypeCreatedResponse {
    @statusCode statusCode: 201;
    @body eventType: EventType;
  }

  model AvailabilityCreatedResponse {
    @statusCode statusCode: 201;
    @body availability: Availability;
  }

  model BookingCreatedResponse {
    @statusCode statusCode: 201;
    @body booking: Booking;
  }

  @doc("Returned when request validation fails.")
  model ValidationErrorResponse {
    @statusCode statusCode: 400;

    code: "validation_error";

    message: string;

    details?: string[];
  }

  @doc("Returned when a requested resource does not exist.")
  model NotFoundErrorResponse {
    @statusCode statusCode: 404;

    code: "not_found";

    message: string;
  }

  @doc("Returned when the selected slot is already booked or overlaps another booking.")
  model BookingConflictResponse {
    @statusCode statusCode: 409;

    code: "booking_conflict";

    message: string;
  }

  @summary("List event types")
  @doc("Returns all event types visible to guests.")
  @route("/event-types")
  @get
  op listEventTypes(): EventType[];

  @summary("Create event type")
  @doc("Creates a new event type. Host pages are public in the MVP.")
  @route("/host/event-types")
  @post
  op createEventType(
    @body request: CreateEventTypeRequest,
  ): EventTypeCreatedResponse | ValidationErrorResponse;

  @summary("Create availability")
  @doc("Creates a host availability range. The range is later split into 30-minute slots.")
  @route("/host/availability")
  @post
  op createAvailability(
    @body request: CreateAvailabilityRequest,
  ): AvailabilityCreatedResponse | ValidationErrorResponse;

  @summary("List available slots")
  @doc("Returns available 30-minute slots for the selected event type within the next 14 days.")
  @route("/event-types/{eventTypeId}/slots")
  @get
  op listSlots(
    @path eventTypeId: string,
  ): Slot[] | NotFoundErrorResponse;

  @summary("Create booking")
  @doc("Creates a booking for a selected slot. Double booking must be prevented across all event types.")
  @route("/bookings")
  @post
  op createBooking(
    @body request: CreateBookingRequest,
  ): BookingCreatedResponse | ValidationErrorResponse | NotFoundErrorResponse | BookingConflictResponse;

  @summary("List host bookings")
  @doc("Returns upcoming bookings across all event types. Host pages are public in the MVP.")
  @route("/host/bookings")
  @get
  op listHostBookings(): Booking[];
}