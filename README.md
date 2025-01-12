目前只是开发了后台部分，后续还会继续跟进优化！

# Simple Bank Project

## Overview
Simple Bank is a backend API project designed to simulate basic banking functionalities. It is implemented using Go programming language and follows RESTful principles. The project includes features such as account management, transactions, and token-based authentication.

## Features
- **Account Management**:
  - Create, retrieve, and list bank accounts.
  - Handle account details securely.

- **Transactions**:
  - Perform money transfers between accounts.
  - Ensure data integrity and consistency with database transactions.

- **Authentication**:
  - Secure user sessions with token-based mechanisms (e.g., PASETO).
  - Validate and manage user authentication efficiently.

## Project Structure
```
.github/workflows/    # CI/CD pipeline configurations.
api/                  # RESTful API handlers.
db/                   # Database migrations and queries.
token/                # Token creation and verification logic.
util/                 # Utility functions and helper modules.
Dockerfile            # Docker image configuration.
Makefile              # Build and automation scripts.
docker-compose.yaml   # Docker Compose for multi-container setup.
app.env               # Environment configuration file.
go.mod, go.sum        # Go module dependencies.
main.go               # Entry point of the application.
sqlc.yaml             # SQLC configuration for code generation.
```

## Prerequisites
- Go (1.20 or higher).
- Docker and Docker Compose.
- PostgreSQL (configured in `app.env`).

## Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/Whuichenggong/Simple_bank.git
cd Simple_bank
```

### 2. Set Up Environment Variables
Configure the `app.env` file with your environment settings. Example:
```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://user:password@localhost:5432/simple_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=your_secret_key
ACCESS_TOKEN_DURATION=15m
```

### 3. Build and Run the Application
#### Using Docker Compose
```bash
docker-compose up --build
```

#### Without Docker
Install dependencies and run the application manually:
```bash
make run
```

## Testing
- Run unit tests to ensure functionality:
```bash
make test
```

## API Documentation
- Use tools like Postman or Swagger to explore and test API endpoints.
- API examples include:
  - **POST /accounts**: Create a new account.
  - **POST /transfers**: Perform a money transfer.

## Future Improvements
- Add more comprehensive test coverage.
- Integrate logging and monitoring tools.
- Implement rate limiting for API security.
- Develop a frontend for user interaction.

## Contributing
Contributions are welcome! Please fork the repository, create a feature branch, and submit a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

