# effective-mobile-trainee-assignment
# Example environment variables

```env
# Docker
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=songs_db

# Service
SERVICENAME=song_library_service
ENV=local

HTTP_PORT=8080
HTTP_TIMEOUT=5s

# Storage
STORAGE_PATH=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable

# External API
EXTERNALAPI=http://172.17.0.1:8081
```

# Update handler
Please use numerical values instead of additionalProp, as shown in the example.

![alt text](image-1.png)