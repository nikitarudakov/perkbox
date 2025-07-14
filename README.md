# Perkbox User Management Service

A REST API for managing users within businesses at Perkbox. Built with Go (Gin) and PostgreSQL, the service enforces role-based access control (RBAC) and provides endpoints for admins and regular users.

---

## 🧠 Design Notes

### Simulated RBAC via Headers

To simulate authentication, the API uses custom headers:

- `X-User-Id`
- `X-User-Role` (`admin` or `user`)
- `X-User-Business`

In a real-world system, these values would be extracted from a JWT or session. 
This header-based approach simplifies testing and avoids implementing a full auth layer for the assignment.

---

## ⚙️ Setup

### 1. Clone & run

```bash
git clone https://github.com/nikitarudakov/perkbox-user-service.git
cd perkbox-user-service
docker-compose up --build
```

### 2. Environment 

⚠️ .env is included for simplicity. In production, sensitive values like DB credentials would not be committed.


### 3. Testing

```bash
go test ./...
```

### 4. Structure

- cmd/api/ – app entry
- infra/ - DB setup
- internal/handlers – HTTP logic
- internal/repo – DB access (GORM)
- internal/domain – domain types