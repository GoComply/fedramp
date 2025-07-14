# Production Deployment Guide for FedRAMP R5 Balance & 20x

## Prerequisites

### System Requirements
- **OS**: Linux (Ubuntu 20.04+ or RHEL 8+)
- **CPU**: 4+ cores recommended
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 50GB+ for database and logs
- **Network**: HTTPS-enabled load balancer

### Software Requirements
- Docker 20.10+
- Docker Compose 2.0+
- PostgreSQL 13+ (for production database)
- Redis 6+ (for caching)
- Nginx or similar reverse proxy

## Pre-Production Checklist

### 🔒 Security
- [ ] Enable TLS/SSL certificates
- [ ] Configure firewall rules
- [ ] Set up WAF (Web Application Firewall)
- [ ] Enable database encryption at rest
- [ ] Configure secure environment variables
- [ ] Set up audit logging
- [ ] Configure CORS properly
- [ ] Enable rate limiting

### 🔧 Configuration
- [ ] Update database connection strings
- [ ] Configure Redis connection
- [ ] Set production log levels
- [ ] Configure alert endpoints
- [ ] Set up backup schedules
- [ ] Configure monitoring endpoints

### 📊 Monitoring
- [ ] Set up Prometheus metrics
- [ ] Configure Grafana dashboards
- [ ] Set up ELK stack for logs
- [ ] Configure health check endpoints
- [ ] Set up uptime monitoring

## Production Configuration

### 1. Environment Variables

Create `.env.production`:

```bash
# Database Configuration
DB_HOST=your-postgres-host
DB_PORT=5432
DB_NAME=fedramp_prod
DB_USER=fedramp_user
DB_PASSWORD=your-secure-password
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENABLE_HTTPS=true
TLS_CERT_FILE=/etc/ssl/certs/fedramp.crt
TLS_KEY_FILE=/etc/ssl/private/fedramp.key

# Authentication
JWT_SECRET=your-very-long-random-secret
SESSION_SECRET=another-very-long-random-secret
OAUTH_CLIENT_ID=your-oauth-client-id
OAUTH_CLIENT_SECRET=your-oauth-client-secret

# Monitoring
ENABLE_METRICS=true
METRICS_PORT=9090
ENABLE_TRACING=true
JAEGER_ENDPOINT=http://jaeger:14268/api/traces

# Alerts
ALERT_WEBHOOK_URL=https://your-webhook-endpoint
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
PAGERDUTY_API_KEY=your-pagerduty-key
SMTP_HOST=smtp.your-provider.com
SMTP_PORT=587
SMTP_USER=alerts@yourdomain.com
SMTP_PASSWORD=your-smtp-password

# Feature Flags
ENABLE_DASHBOARD=true
ENABLE_API_DOCS=false
ENABLE_DEBUG=false
```

### 2. Docker Compose Production

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  fedramp-api:
    image: fedramp-server:latest
    container_name: fedramp-api
    restart: always
    ports:
      - "8080:8080"
    environment:
      - ENV=production
    env_file:
      - .env.production
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    depends_on:
      - postgres
      - redis
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - fedramp-network

  postgres:
    image: postgres:15-alpine
    container_name: fedramp-db
    restart: always
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./init-db.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - fedramp-network

  redis:
    image: redis:7-alpine
    container_name: fedramp-cache
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
    networks:
      - fedramp-network

  nginx:
    image: nginx:alpine
    container_name: fedramp-proxy
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
      - ./web/dashboard:/usr/share/nginx/html
    depends_on:
      - fedramp-api
    networks:
      - fedramp-network

  prometheus:
    image: prom/prometheus:latest
    container_name: fedramp-metrics
    restart: always
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    networks:
      - fedramp-network

  grafana:
    image: grafana/grafana:latest
    container_name: fedramp-grafana
    restart: always
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
    networks:
      - fedramp-network

volumes:
  postgres-data:
  redis-data:
  prometheus-data:
  grafana-data:

networks:
  fedramp-network:
    driver: bridge
```

### 3. Nginx Configuration

Create `nginx.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    upstream fedramp_api {
        server fedramp-api:8080;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

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
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;

        # Security headers
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        # API endpoints
        location /api/ {
            limit_req zone=api_limit burst=20 nodelay;
            proxy_pass http://fedramp_api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Dashboard
        location / {
            root /usr/share/nginx/html;
            try_files $uri $uri/ /index.html;
        }

        # Health check
        location /health {
            access_log off;
            proxy_pass http://fedramp_api/api/v1/health;
        }
    }
}
```

## Deployment Steps

### 1. Build Production Image

```bash
# Build the production Docker image
docker build -t fedramp-server:latest -f Dockerfile.prod .

# Tag for your registry
docker tag fedramp-server:latest your-registry.com/fedramp-server:latest

# Push to registry
docker push your-registry.com/fedramp-server:latest
```

### 2. Database Setup

```bash
# Create production database
psql -h your-postgres-host -U postgres << EOF
CREATE DATABASE fedramp_prod;
CREATE USER fedramp_user WITH ENCRYPTED PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE fedramp_prod TO fedramp_user;
EOF

# Run migrations
docker run --rm \
  -e DB_HOST=your-postgres-host \
  -e DB_USER=fedramp_user \
  -e DB_PASSWORD=your-secure-password \
  -e DB_NAME=fedramp_prod \
  fedramp-server:latest \
  /app/fedramp-server migrate
```

### 3. Deploy with Docker Compose

```bash
# Deploy the stack
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 4. Verify Deployment

```bash
# Check health endpoint
curl https://your-domain.com/api/v1/health

# Test authentication
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your-password"}'

# Check metrics
curl http://your-domain.com:9090/metrics
```

## Production Monitoring

### 1. Set Up Alerts

Configure alerts for:
- API response time > 1s
- Error rate > 1%
- Database connection failures
- Memory usage > 80%
- CPU usage > 80%
- Disk usage > 80%

### 2. Log Aggregation

```bash
# View API logs
docker logs -f fedramp-api

# View all logs
docker-compose -f docker-compose.prod.yml logs -f

# Export logs
docker logs fedramp-api > api-logs-$(date +%Y%m%d).log
```

### 3. Backup Strategy

```bash
# Database backup
docker exec fedramp-db pg_dump -U fedramp_user fedramp_prod > backup-$(date +%Y%m%d).sql

# Automated daily backups
0 2 * * * docker exec fedramp-db pg_dump -U fedramp_user fedramp_prod | gzip > /backups/fedramp-$(date +\%Y\%m\%d).sql.gz
```

## Security Hardening

### 1. API Security

```bash
# Enable rate limiting (already in nginx.conf)
# Enable CORS restrictions
# Implement API key authentication
# Set up OAuth2/SAML for production
```

### 2. Database Security

```sql
-- Restrict database access
REVOKE ALL ON DATABASE fedramp_prod FROM PUBLIC;

-- Create read-only user for reporting
CREATE USER fedramp_readonly WITH ENCRYPTED PASSWORD 'readonly-password';
GRANT CONNECT ON DATABASE fedramp_prod TO fedramp_readonly;
GRANT USAGE ON SCHEMA public TO fedramp_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO fedramp_readonly;
```

### 3. Network Security

```bash
# Firewall rules (example for UFW)
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 9090/tcp  # Prometheus (restrict to monitoring network)
ufw enable
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check DB_HOST and credentials
   - Verify PostgreSQL is running
   - Check firewall rules

2. **High Memory Usage**
   - Adjust container memory limits
   - Check for memory leaks
   - Review query optimization

3. **Slow API Response**
   - Check database indexes
   - Review N+1 queries
   - Enable query caching

### Health Checks

```bash
# Full system health check
curl https://your-domain.com/api/v1/health/full

# Database health
docker exec fedramp-api /app/fedramp-server health db

# Redis health
docker exec fedramp-api /app/fedramp-server health redis
```

## Rollback Procedure

```bash
# Stop current deployment
docker-compose -f docker-compose.prod.yml down

# Restore database
psql -h your-postgres-host -U fedramp_user fedramp_prod < backup-20250110.sql

# Deploy previous version
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

## Support

For production support:
- Check logs first: `docker logs fedramp-api`
- Review metrics: Grafana dashboards
- Contact: devops@your-organization.com

---

**Note**: This guide assumes you're deploying to a cloud provider or on-premises infrastructure. Adjust configurations based on your specific environment and security requirements. 