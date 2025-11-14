# Meeting Room Booking API

A comprehensive RESTful API for managing meeting room bookings with advanced search, filtering, and availability checking features.

## Features

### Core Functionality

- **User Management**: Admin can register users with different roles
- **Authentication**: JWT-based authentication system
- **Room Management**: Create, read, update, and delete meeting rooms
- **Booking System**: Book rooms, view schedules, and cancel bookings

### Advanced Features (NEW)

- **Smart Room Search**: Filter rooms by capacity, floor, amenities, and availability
- **Availability Checking**: Check if a room is available for specific time ranges
- **Detailed Schedule View**: View bookings with user and room details
- **Conflict Detection**: Get detailed information about conflicting bookings
- **Time Slot Suggestions**: Receive suggestions for available time slots

## Tech Stack

- **Language**: Go 1.24+
- **Database**: SQLite
- **Authentication**: JWT (JSON Web Tokens)
- **Router**: Gorilla Mux
- **API Documentation**: OpenAPI 3.0

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Docker (optional)

### Local Development

```bash
# Clone the repository
git clone https://github.com/amangirdhar210/meeting-room.git
cd meeting-room

# Install dependencies
go mod download

# Run the application
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Using Docker

```bash
# Build the image
docker build -t meeting-room-app .

# Run the container
docker run -p 8080:8080 meeting-room-app

# Run with persistent data
docker run -p 8080:8080 -v $(pwd)/data:/root meeting-room-app
```

## API Endpoints

### Authentication

- `POST /api/login` - Login and receive JWT token

### Users (Admin Only)

- `POST /api/register` - Register a new user
- `GET /api/users` - Get all users
- `DELETE /api/users/{id}` - Delete a user

### Rooms

- `POST /api/rooms` - Add a new room (admin only)
- `GET /api/rooms` - Get all rooms
- `GET /api/rooms/search` - **NEW** Search rooms with filters
- `POST /api/rooms/check-availability` - **NEW** Check room availability
- `GET /api/rooms/{id}` - Get room details
- `DELETE /api/rooms/{id}` - Delete a room (admin only)
- `GET /api/rooms/{id}/schedule` - Get room schedule with detailed booking information

### Bookings

- `POST /api/bookings` - Create a new booking
- `GET /api/bookings` - Get all bookings (admin only)
- `DELETE /api/bookings/{id}` - Cancel a booking (admin only)

## Frontend-Friendly Features

### 1. Room Search with Filters

**Endpoint**: `GET /api/rooms/search`

**Query Parameters**:

- `minCapacity` - Minimum room capacity (integer)
- `maxCapacity` - Maximum room capacity (integer)
- `floor` - Specific floor number (integer)
- `amenities` - Required amenities (string)
- `startTime` - Check availability from (RFC3339)
- `endTime` - Check availability until (RFC3339)

**Example Request**:

```bash
GET /api/rooms/search?minCapacity=5&floor=1&amenities=Projector&startTime=2025-11-14T09:00:00Z&endTime=2025-11-14T10:00:00Z
```

**Response**:

```json
[
  {
    "id": 1,
    "name": "Conference Room A",
    "roomNumber": 101,
    "capacity": 10,
    "floor": 1,
    "amenities": ["Projector", "Whiteboard"],
    "status": "Available",
    "location": "Building A, Floor 1"
  }
]
```

### 2. Availability Checking

**Endpoint**: `POST /api/rooms/check-availability`

**Request Body**:

```json
{
  "roomId": 1,
  "startTime": "2025-11-14T09:00:00Z",
  "endTime": "2025-11-14T10:00:00Z"
}
```

**Response**:

```json
{
  "available": false,
  "roomId": 1,
  "roomName": "Conference Room A",
  "requestedStart": "2025-11-14T09:00:00Z",
  "requestedEnd": "2025-11-14T10:00:00Z",
  "conflictingSlots": [
    {
      "bookingId": 5,
      "startTime": "2025-11-14T09:30:00Z",
      "endTime": "2025-11-14T10:30:00Z",
      "purpose": "Board meeting"
    }
  ],
  "suggestedSlots": []
}
```

### 3. Detailed Room Schedule

**Endpoint**: `GET /api/rooms/{id}/schedule`

**Response**:

```json
[
  {
    "id": 1,
    "user_id": 2,
    "userName": "John Doe",
    "userEmail": "john.doe@example.com",
    "room_id": 1,
    "roomName": "Conference Room A",
    "roomNumber": 101,
    "start_time": "2025-11-14T09:00:00Z",
    "end_time": "2025-11-14T10:00:00Z",
    "duration": 60,
    "purpose": "Team standup meeting",
    "status": "confirmed"
  }
]
```

## Authentication

All endpoints except `/health` and `/api/login` require JWT authentication.

1. Login to get your token:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

2. Use the token in subsequent requests:

```bash
curl -X GET http://localhost:8080/api/rooms \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Default Admin Credentials

- **Email**: admin@example.com
- **Password**: admin123

⚠️ **Important**: Change the default password in production!

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   ├── router.go              # Route definitions
│   │   └── server.go
│   ├── domain/
│   │   ├── booking.go             # Domain models
│   │   ├── room.go
│   │   ├── user.go
│   │   └── interfaces.go          # Service and repository interfaces
│   ├── http/
│   │   ├── dto/                   # Data Transfer Objects
│   │   │   ├── booking_dto.go
│   │   │   ├── room_dto.go
│   │   │   └── user_dto.go
│   │   ├── handlers/              # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── booking_handler.go
│   │   │   ├── room_handler.go
│   │   │   └── user_handler.go
│   │   └── middleware/            # Authentication & logging
│   ├── repositories/
│   │   └── mysql/                 # Database implementations
│   │       ├── booking_repo.go
│   │       ├── room_repo.go
│   │       └── user_repo.go
│   ├── service/                   # Business logic
│   │   ├── auth_service.go
│   │   ├── booking_service.go
│   │   ├── room_service.go
│   │   └── user_service.go
│   └── pkg/                       # Shared utilities
│       ├── jwt/
│       ├── logger/
│       └── utils/
├── dockerfile                     # Docker configuration
├── go.mod                         # Go dependencies
└── openapi.yaml                   # API documentation
```

## Best Practices Implemented

### 1. Clean Architecture

- Clear separation of concerns (domain, service, repository, handler)
- Dependency injection
- Interface-based design

### 2. Security

- JWT authentication
- Password hashing with bcrypt
- Role-based access control (RBAC)

### 3. Error Handling

- Consistent error responses
- Proper HTTP status codes
- Descriptive error messages

### 4. Code Quality

- Latest Go version (1.24+)
- Context-based timeouts
- Proper resource cleanup (defer)
- Idiomatic Go code

### 5. Database

- Prepared statements (SQL injection protection)
- Transaction support
- Foreign key constraints
- Indexed queries

## Frontend Integration Guide

### Recommended Frontend Flow

1. **Login Screen**

   - Call `/api/login` to get JWT token
   - Store token in localStorage or secure cookie

2. **Room Browse/Search Page**

   - Use `/api/rooms/search` with filters
   - Display room cards with capacity, amenities, location
   - Add "View Schedule" button for each room

3. **Room Schedule View**

   - Call `/api/rooms/{id}/schedule` to show all bookings
   - Display calendar/timeline view with:
     - Booked slots (with user name and purpose)
     - Available slots
     - Duration information

4. **Availability Check Before Booking**

   - Before showing booking form, call `/api/rooms/check-availability`
   - If conflicts exist, show them to the user
   - Display suggested alternative time slots

5. **Create Booking**
   - Call `/api/bookings` with selected time slot
   - Handle 409 Conflict gracefully
   - Show success message with booking details

### Example Frontend Workflow

```javascript
// 1. Search for available rooms
const searchRooms = async (minCapacity, floor, startTime, endTime) => {
  const response = await fetch(
    `/api/rooms/search?minCapacity=${minCapacity}&floor=${floor}&startTime=${startTime}&endTime=${endTime}`,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return await response.json();
};

// 2. Check specific room availability
const checkAvailability = async (roomId, startTime, endTime) => {
  const response = await fetch("/api/rooms/check-availability", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ roomId, startTime, endTime }),
  });
  return await response.json();
};

// 3. Get room schedule
const getRoomSchedule = async (roomId) => {
  const response = await fetch(`/api/rooms/${roomId}/schedule`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return await response.json();
};

// 4. Create booking
const createBooking = async (roomId, startTime, endTime, purpose) => {
  const response = await fetch("/api/bookings", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      room_id: roomId,
      start_time: startTime,
      end_time: endTime,
      purpose,
    }),
  });
  return await response.json();
};
```

## Testing with Postman

Import the `openapi.yaml` file into Postman to get a complete collection with:

- All endpoints pre-configured
- Example request bodies
- Environment variables setup
- Authentication flows

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on GitHub.
