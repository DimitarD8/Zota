# Zota
Key-Value Store which randomly breaks

---

## Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/DimitarD8/Zota.git
cd zota
```

### 2. Set up environment variables
Fill in the `.env` file with your PostgreSQL connection details:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=zota
```

### 3. Run the application
You can start the project either from the terminal:
```bash
go run main.go
```

Or directly from your IDE.

---

## Project Structure
```
.
├── db/             # Database configuration and connection
├── store/          # Store implementation (Put, Get, Delete, Dump)
├── queries/        # SQL queries as constants
├── migrations/     # Schema initialization (previously initializer)
├── main.go         # Entry point of the application
└── go.mod
```

---

##  How it works
- PostgreSQL is used as the database for storing key-value pairs.
- The `Store` struct provides an API with methods:
    - `Put(key, value)` → insert or update a record
    - `Get(key)` → retrieve a value by key
    - `Delete(key)` → remove a record
    - `Dump()` → return the entire dataset
- The `isUnlucky()` function randomly fails in ~30% of cases, simulating unstable systems or network failures.

---

## Testing
The project uses pgxmock for unit testing.

- Tests are written only for the core functionality of the Store.
- Tests cover both happy paths (successful operations) and error scenarios (simulated failures).
- Run all tests with:
```bash
go test ./...
```