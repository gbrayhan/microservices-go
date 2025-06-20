# Deployment Guide

## 🚀 Deployment Overview

This guide covers deployment strategies for the Microservices Go application across different environments and platforms.

## 📋 Prerequisites

### System Requirements

- **Go**: 1.24.2+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **PostgreSQL**: 13+
- **Memory**: 2GB+ RAM
- **Storage**: 10GB+ free space

### Environment Variables

```bash
# Server Configuration
SERVER_PORT=8080
GO_ENV=production

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure_password
DB_NAME=microservices_go

# JWT Configuration
JWT_ACCESS_SECRET=your_very_secure_access_secret_key
JWT_REFRESH_SECRET=your_very_secure_refresh_secret_key
JWT_ACCESS_TIME_MINUTE=60
JWT_REFRESH_TIME_HOUR=24
```

## 🐳 Docker Deployment

### Development Environment

```bash
# Clone repository
git clone https://github.com/gbrayhan/microservices-go
cd microservices-go

# Copy environment file
cp .env.example .env

# Start services
docker-compose up --build -d

# Check services
docker-compose ps

# View logs
docker-compose logs -f app
```

### Production Environment

#### 1. Build Production Image

```dockerfile
# Dockerfile.prod
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

#### 2. Docker Compose for Production

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.prod
    ports:
      - "8080:8080"
    environment:
      - GO_ENV=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - JWT_ACCESS_SECRET=${JWT_ACCESS_SECRET}
      - JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
    depends_on:
      - postgres
    restart: unless-stopped
    networks:
      - app-network

  postgres:
    image: postgres:13-alpine
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

#### 3. Nginx Configuration

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

#### 4. Production Deployment

```bash
# Set environment variables
export DB_USER=postgres
export DB_PASSWORD=secure_password
export DB_NAME=microservices_go
export JWT_ACCESS_SECRET=your_very_secure_access_secret_key
export JWT_REFRESH_SECRET=your_very_secure_refresh_secret_key

# Deploy with production compose
docker-compose -f docker-compose.prod.yml up --build -d

# Check deployment
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f app
```

## ☸️ Kubernetes Deployment

### Namespace and ConfigMap

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: microservices-go
```

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: microservices-go
data:
  SERVER_PORT: "8080"
  GO_ENV: "production"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "microservices_go"
```

### Secret

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: microservices-go
type: Opaque
data:
  DB_USER: cG9zdGdyZXM=  # postgres
  DB_PASSWORD: c2VjdXJlX3Bhc3N3b3Jk  # secure_password
  JWT_ACCESS_SECRET: eW91cl92ZXJ5X3NlY3VyZV9hY2Nlc3Nfc2VjcmV0X2tleQ==
  JWT_REFRESH_SECRET: eW91cl92ZXJ5X3NlY3VyZV9yZWZyZXNoX3NlY3JldF9rZXk=
```

### PostgreSQL Deployment

```yaml
# postgres-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: microservices-go
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:13-alpine
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DB_NAME
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_USER
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: microservices-go
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: microservices-go
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

### Application Deployment

```yaml
# app-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: microservices-go
  namespace: microservices-go
spec:
  replicas: 3
  selector:
    matchLabels:
      app: microservices-go
  template:
    metadata:
      labels:
        app: microservices-go
    spec:
      containers:
      - name: app
        image: microservices-go:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_PORT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: SERVER_PORT
        - name: GO_ENV
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: GO_ENV
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DB_HOST
        - name: DB_PORT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DB_PORT
        - name: DB_NAME
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DB_NAME
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_USER
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
        - name: JWT_ACCESS_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: JWT_ACCESS_SECRET
        - name: JWT_REFRESH_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: JWT_REFRESH_SECRET
        livenessProbe:
          httpGet:
            path: /v1/health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: microservices-go-service
  namespace: microservices-go
spec:
  selector:
    app: microservices-go
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Ingress Configuration

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: microservices-go-ingress
  namespace: microservices-go
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - your-domain.com
    secretName: tls-secret
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: microservices-go-service
            port:
              number: 80
```

### Kubernetes Deployment Commands

```bash
# Create namespace
kubectl apply -f namespace.yaml

# Apply configurations
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml

# Deploy PostgreSQL
kubectl apply -f postgres-deployment.yaml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n microservices-go

# Deploy application
kubectl apply -f app-deployment.yaml

# Apply ingress
kubectl apply -f ingress.yaml

# Check deployment
kubectl get pods -n microservices-go
kubectl get services -n microservices-go
kubectl get ingress -n microservices-go
```

## 🚀 CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24.2'
    
    - name: Run tests
      run: |
        go test -v ./...
        ./coverage.sh
    
    - name: Run security scan
      run: |
        trivy fs .
        gosec ./...

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Login to Docker Hub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: |
          your-registry/microservices-go:latest
          your-registry/microservices-go:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v3
    
    - name: Configure kubectl
      uses: azure/k8s-set-context@v3
      with:
        kubeconfig: ${{ secrets.KUBE_CONFIG }}
    
    - name: Deploy to Kubernetes
      run: |
        kubectl set image deployment/microservices-go app=your-registry/microservices-go:${{ github.sha }} -n microservices-go
        kubectl rollout status deployment/microservices-go -n microservices-go
```

### GitLab CI/CD Pipeline

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"

test:
  stage: test
  image: golang:1.24.2
  script:
    - go test -v ./...
    - ./coverage.sh
    - trivy fs .
    - gosec ./...

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - docker tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA $CI_REGISTRY_IMAGE:latest
    - docker push $CI_REGISTRY_IMAGE:latest
  only:
    - main

deploy:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/microservices-go app=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA -n microservices-go
    - kubectl rollout status deployment/microservices-go -n microservices-go
  only:
    - main
```

## 🔧 Monitoring and Observability

### Health Check Endpoint

```go
// Add to your routes
router.GET("/v1/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now().UTC(),
        "version": "1.0.0",
        "uptime": time.Since(startTime).String(),
    })
})
```

### Prometheus Metrics

```go
// Add Prometheus metrics
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
}

// Add metrics endpoint
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "Microservices Go Dashboard",
    "panels": [
      {
        "title": "HTTP Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

## 🔒 Security Considerations

### SSL/TLS Configuration

```bash
# Generate SSL certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ssl/key.pem -out ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=your-domain.com"
```

### Security Headers

```go
// Add security middleware
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
```

### Rate Limiting

```go
// Add rate limiting
import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(rate.Every(time.Second), 100)

func RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 📊 Performance Optimization

### Database Optimization

```sql
-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_medicines_name ON medicines(name);
CREATE INDEX idx_medicines_laboratory ON medicines(laboratory);

-- Optimize queries
EXPLAIN ANALYZE SELECT * FROM users WHERE email ILIKE '%john%';
```

### Application Optimization

```go
// Connection pooling
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func initDB() *gorm.DB {
    dsn := "host=localhost user=postgres password=password dbname=microservices_go port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        PrepareStmt: true,
    })
    
    sqlDB, err := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db
}
```

## 🔄 Backup and Recovery

### Database Backup

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="microservices_go"

# Create backup
pg_dump -h localhost -U postgres -d $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 7 days of backups
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
```

### Application Backup

```bash
#!/bin/bash
# app-backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/app-backups"

# Backup application data
tar -czf $BACKUP_DIR/app_backup_$DATE.tar.gz /app/data

# Backup configuration
cp /app/.env $BACKUP_DIR/env_backup_$DATE

# Keep only last 30 days of backups
find $BACKUP_DIR -name "app_backup_*.tar.gz" -mtime +30 -delete
find $BACKUP_DIR -name "env_backup_*" -mtime +30 -delete
```

## 🚨 Troubleshooting

### Common Issues

1. **Database Connection Issues**
   ```bash
   # Check database connectivity
   docker exec -it postgres psql -U postgres -d microservices_go -c "SELECT 1;"
   ```

2. **Application Startup Issues**
   ```bash
   # Check application logs
   docker-compose logs -f app
   kubectl logs -f deployment/microservices-go -n microservices-go
   ```

3. **Memory Issues**
   ```bash
   # Check memory usage
   docker stats
   kubectl top pods -n microservices-go
   ```

### Performance Monitoring

```bash
# Monitor application performance
curl -X GET http://localhost:8080/v1/health
curl -X GET http://localhost:8080/metrics

# Database performance
docker exec -it postgres psql -U postgres -d microservices_go -c "SELECT * FROM pg_stat_activity;"
```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)

---

This deployment guide provides comprehensive instructions for deploying the Microservices Go application across different environments and platforms. 