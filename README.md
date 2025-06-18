
````markdown
# Basic Web Server Project

This is a simple Go web server project that allows users to register via a form, view the homepage, and see a list of registered users.

---

## Features

- **User Registration Form**: Accessible at `/form` — allows users to register by providing their name, username, email, and password.
- **Homepage**: Accessible at `/home` — the main landing page of the application.
- **Registered Users List**: Accessible at `/registered-users` — displays a list of all users who have registered.

---

## Getting Started

### Prerequisites

- Go 1.18+ installed on your machine
- MongoDB instance running (local or cloud)
- `.env` file configured with your MongoDB URI (e.g. `MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/dbname`)

---

### Installation & Running

1. Clone the repository:

   ```bash
   git clone https://github.com/Deba00407/basic-web-server.git
   cd basic-web-server
````

2. Install dependencies and tidy the module:

   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory with the following content:

   ```
   MONGODB_URI=your_mongodb_connection_string
   ```

4. Run the server:

   ```bash
   go run main.go
   ```

   The server will start on port `5001`.

---

## Accessing the Webpages

* **Homepage**: [http://localhost:5001/home](http://localhost:5001/home)
  The main landing page.

* **User Registration Form**: [http://localhost:5001/form](http://localhost:5001/form)
  Fill in your details and submit to register a new user.

* **Registered Users List**: [http://localhost:5001/registered-users](http://localhost:5001/registered-users)
  View all users who have registered so far.

---

## Project Structure

```
basic-web-server/
├── controllers/        # Handlers for routes and business logic
│   └── ...             
├── database/           # MongoDB connection and CRUD functions
│   └── ...             
├── schemamodels/       # Data structures (e.g., User, FormData)
│   └── ...             
├── templates/          # HTML templates for rendering pages
│   └── ...             
├── .env                # Environment variables (MongoDB URI etc.)
├── go.mod              # Go module file
└── main.go             # Entry point, routing, and server startup
```

---

## Notes

* Ensure MongoDB is running and accessible via the URI provided in `.env`.
* You can extend this project by adding features like authentication, session management, password reset, etc.
---
