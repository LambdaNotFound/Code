package apidesign

import "time"

// === Core Data Models ===
type Room struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Location string `json:"location,omitempty"`
    Capacity int    `json:"capacity,omitempty"`
}

type Meeting struct {
    ID          string    `json:"id"`
    RoomID      string    `json:"room_id"`
    OrganizerID string    `json:"organizer_id"`
    Title       string    `json:"title"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Attendees   []string  `json:"attendees"`
    Status      string    `json:"status"`
}

// === Request/Response Contracts ===
type CreateMeetingRequest struct {
    RoomID      string    `json:"room_id" validate:"required"`
    OrganizerID string    `json:"organizer_id" validate:"required"`
    Title       string    `json:"title" validate:"required"`
    StartTime   time.Time `json:"start_time" validate:"required"`
    EndTime     time.Time `json:"end_time" validate:"required"`
    Attendees   []string  `json:"attendees"`
}

type CreateMeetingResponse struct {
    MeetingID string `json:"meeting_id"`
    Status    string `json:"status"`
}

type ListRoomsResponse struct {
    Rooms []Room `json:"rooms"`
}

type GetRoomAvailabilityResponse struct {
    AvailableRooms []Room `json:"available_rooms"`
}

type GetRoomMeetingsResponse struct {
    Meetings []Meeting `json:"meetings"`
}

type CancelMeetingResponse struct {
    Status string `json:"status"`
}

type MeetingService interface {
    CreateMeeting(req CreateMeetingRequest) (CreateMeetingResponse, error)
    CancelMeeting(id string) (CancelMeetingResponse, error)
    GetRoomAvailability(start, end string) (GetRoomAvailabilityResponse, error)
    ListRooms() (ListRoomsResponse, error)
    GetRoomMeetings(roomID string, date string) (GetRoomMeetingsResponse, error)
}

type Handler struct {
    Svc MeetingService
}
