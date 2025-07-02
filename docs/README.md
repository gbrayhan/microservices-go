# Documentation Index

Welcome to the comprehensive documentation for the **Microservices Go** application. This documentation provides detailed information about the architecture, API, deployment, and development guidelines.

## 📚 Documentation Structure

### 🏗️ Architecture & Design

- **[Clean Architecture Guide](README_CLEAN_ARCHITECTURE.md)** - Detailed implementation of Clean Architecture principles with diagrams and examples
- **[API Documentation](API_DOCUMENTATION.md)** - Complete API reference with endpoints, examples, and authentication flows

### 🔍 Features & Endpoints

- **[Search Endpoints](SEARCH_ENDPOINTS.md)** - Advanced search and pagination capabilities for Users and Medicines

### 🚀 Deployment & Operations

- **[Deployment Guide](DEPLOYMENT_GUIDE.md)** - Comprehensive deployment strategies for Docker, Kubernetes, and CI/CD

## 🎯 Quick Start

### 1. **Understanding the Architecture**
Start with the [Clean Architecture Guide](README_CLEAN_ARCHITECTURE.md) to understand the project structure and design principles.

### 2. **API Reference**
Review the [API Documentation](API_DOCUMENTATION.md) to understand all available endpoints and their usage.

### 3. **Search Features**
Explore the [Search Endpoints](SEARCH_ENDPOINTS.md) to learn about advanced filtering and pagination capabilities.

### 4. **Deployment**
Follow the [Deployment Guide](DEPLOYMENT_GUIDE.md) to deploy the application in your preferred environment.

## 📊 Documentation Overview

### Architecture Diagrams

The documentation includes comprehensive diagrams:

- **Clean Architecture Layers** - Shows the dependency flow and layer structure
- **Authentication Flow** - JWT token lifecycle and state management
- **Data Flow** - Request/response patterns and error handling
- **Deployment Architecture** - Docker and Kubernetes deployment strategies

### Code Examples

Each documentation section includes:

- **Go Code Examples** - Implementation patterns and best practices
- **API Examples** - Request/response formats with curl commands
- **Configuration Examples** - Environment variables and deployment configs
- **Testing Examples** - Unit and integration test patterns

## 🔧 Development Workflow

### 1. **Local Development**
```bash
# Clone and setup
git clone https://github.com/gbrayhan/microservices-go
cd microservices-go
cp .env.example .env

# Start services
docker-compose up --build -d

# Run tests
./coverage.sh
./scripts/run-integration-test.bash
```

### 2. **API Testing**
```bash
# Health check
curl http://localhost:8080/v1/health

# Test authentication
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### 3. **Code Quality**
```bash
# Linting
golangci-lint run ./...

# Security scan
trivy fs .
gosec ./...
```

## 📈 Key Features Documented

### ✅ Clean Architecture Implementation
- **Domain Layer** - Business entities and rules
- **Application Layer** - Use cases and business logic
- **Infrastructure Layer** - External implementations
- **Dependency Injection** - Centralized dependency management

### ✅ Authentication & Security
- **JWT Authentication** - Access and refresh tokens
- **Password Security** - bcrypt hashing with salt
- **Input Validation** - Request sanitization and validation
- **Error Handling** - Centralized error management

### ✅ Advanced Search & Pagination
- **Multi-field Search** - LIKE and exact match filters
- **Date Range Filtering** - Time-based queries
- **Sorting** - Multi-field sorting with direction control
- **Pagination** - Efficient result pagination

### ✅ Testing Strategy
- **Unit Tests** - Use case and controller testing
- **Integration Tests** - API endpoint testing
- **Acceptance Tests** - Cucumber-based BDD tests
- **Coverage Analysis** - Comprehensive test coverage reporting

### ✅ Deployment Options
- **Docker** - Containerized deployment
- **Kubernetes** - Orchestrated deployment
- **CI/CD** - Automated testing and deployment
- **Monitoring** - Health checks and metrics

## 🎯 Best Practices

### Code Quality
- **Clean Architecture** - Follow dependency inversion principles
- **Test Coverage** - Maintain ≥80% test coverage
- **Error Handling** - Use centralized error management
- **Logging** - Structured logging with correlation IDs

### Security
- **Authentication** - JWT with short-lived access tokens
- **Authorization** - Role-based access control
- **Input Validation** - Validate and sanitize all inputs
- **Security Headers** - Implement security headers

### Performance
- **Database Optimization** - Proper indexing and query optimization
- **Connection Pooling** - Efficient database connection management
- **Caching** - Implement caching strategies
- **Monitoring** - Performance metrics and health checks

## 🔄 Documentation Maintenance

### Contributing to Documentation
1. **Update Diagrams** - Use Mermaid for architecture diagrams
2. **Code Examples** - Keep examples up-to-date with code changes
3. **API Documentation** - Update when endpoints change
4. **Deployment Guides** - Update for new deployment options

### Documentation Standards
- **English Only** - All documentation in English
- **Clear Structure** - Use consistent headings and formatting
- **Code Examples** - Include working code examples
- **Diagrams** - Use Mermaid for visual documentation

## 📞 Support & Resources

### Getting Help
- **GitHub Issues** - Report bugs and request features
- **GitHub Discussions** - Ask questions and share ideas
- **Documentation** - This comprehensive documentation
- **API Documentation** - Complete REST API reference

### Additional Resources
- **Go Documentation** - [golang.org](https://golang.org/doc/)
- **Gin Framework** - [gin-gonic.com](https://gin-gonic.com/)
- **GORM** - [gorm.io](https://gorm.io/)
- **Clean Architecture** - [blog.cleancoder.com](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📝 Documentation Changelog

### v2.0.0 (Latest)
- ✅ **Complete English Translation** - All documentation translated to English
- ✅ **Comprehensive Diagrams** - Added Mermaid diagrams for all flows
- ✅ **API Documentation** - Complete API reference with examples
- ✅ **Deployment Guide** - Docker, Kubernetes, and CI/CD strategies
- ✅ **Search Documentation** - Advanced search and pagination features
- ✅ **Testing Documentation** - Unit, integration, and acceptance testing

### v1.0.0
- ✅ **Basic Architecture** - Initial Clean Architecture documentation
- ✅ **API Endpoints** - Basic endpoint documentation
- ✅ **Deployment** - Simple Docker deployment guide

---

**Last Updated**: 2024  
**Documentation Version**: 2.0.0  
**Status**: Complete and Up-to-Date 