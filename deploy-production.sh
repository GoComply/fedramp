#!/bin/bash
# Production Deployment Script for FedRAMP R5 Balance & 20x

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REGISTRY="${DOCKER_REGISTRY:-your-registry.com}"
IMAGE_NAME="fedramp-server"
VERSION="${VERSION:-latest}"
NAMESPACE="${K8S_NAMESPACE:-fedramp-prod}"

echo -e "${GREEN}FedRAMP R5 Balance & 20x Production Deployment${NC}"
echo "================================================"

# Function to check prerequisites
check_prerequisites() {
    echo -e "\n${YELLOW}Checking prerequisites...${NC}"
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Docker is not installed${NC}"
        exit 1
    fi
    
    # Check kubectl (if using Kubernetes)
    if [[ "${DEPLOY_TARGET:-docker}" == "kubernetes" ]] && ! command -v kubectl &> /dev/null; then
        echo -e "${RED}kubectl is not installed${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Prerequisites satisfied${NC}"
}

# Function to build and push Docker image
build_and_push() {
    echo -e "\n${YELLOW}Building Docker image...${NC}"
    
    # Build production image
    docker build -f Dockerfile.prod -t ${IMAGE_NAME}:${VERSION} .
    
    # Tag for registry
    docker tag ${IMAGE_NAME}:${VERSION} ${REGISTRY}/${IMAGE_NAME}:${VERSION}
    docker tag ${IMAGE_NAME}:${VERSION} ${REGISTRY}/${IMAGE_NAME}:latest
    
    echo -e "\n${YELLOW}Pushing to registry...${NC}"
    docker push ${REGISTRY}/${IMAGE_NAME}:${VERSION}
    docker push ${REGISTRY}/${IMAGE_NAME}:latest
    
    echo -e "${GREEN}✓ Image built and pushed successfully${NC}"
}

# Function to deploy with Docker Compose
deploy_docker_compose() {
    echo -e "\n${YELLOW}Deploying with Docker Compose...${NC}"
    
    # Check for environment file
    if [[ ! -f ".env.production" ]]; then
        echo -e "${RED}Missing .env.production file${NC}"
        echo "Please create .env.production from .env.production.example"
        exit 1
    fi
    
    # Deploy stack
    docker-compose -f docker-compose.prod.yml up -d
    
    # Wait for health check
    echo -e "\n${YELLOW}Waiting for services to be healthy...${NC}"
    sleep 10
    
    # Check health
    if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ API is healthy${NC}"
    else
        echo -e "${RED}✗ API health check failed${NC}"
        docker-compose -f docker-compose.prod.yml logs fedramp-api
        exit 1
    fi
}

# Function to deploy to Kubernetes
deploy_kubernetes() {
    echo -e "\n${YELLOW}Deploying to Kubernetes...${NC}"
    
    # Create namespace if it doesn't exist
    kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply configurations
    kubectl apply -f k8s/production/ -n ${NAMESPACE}
    
    # Wait for rollout
    echo -e "\n${YELLOW}Waiting for deployment rollout...${NC}"
    kubectl rollout status deployment/fedramp-api -n ${NAMESPACE}
    
    # Check pod status
    kubectl get pods -n ${NAMESPACE} -l app=fedramp-api
    
    echo -e "${GREEN}✓ Kubernetes deployment successful${NC}"
}

# Function to run database migrations
run_migrations() {
    echo -e "\n${YELLOW}Running database migrations...${NC}"
    
    if [[ "${DEPLOY_TARGET:-docker}" == "kubernetes" ]]; then
        # Run migration job in Kubernetes
        kubectl run fedramp-migrate-${VERSION} \
            --image=${REGISTRY}/${IMAGE_NAME}:${VERSION} \
            --rm -i --tty \
            --restart=Never \
            -n ${NAMESPACE} \
            -- /app/fedramp-server migrate
    else
        # Run migration in Docker
        docker run --rm \
            --env-file .env.production \
            ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
            /app/fedramp-server migrate
    fi
    
    echo -e "${GREEN}✓ Migrations completed${NC}"
}

# Function to run smoke tests
run_smoke_tests() {
    echo -e "\n${YELLOW}Running smoke tests...${NC}"
    
    # Determine endpoint
    if [[ "${DEPLOY_TARGET:-docker}" == "kubernetes" ]]; then
        ENDPOINT=$(kubectl get service fedramp-api -n ${NAMESPACE} -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    else
        ENDPOINT="localhost:8080"
    fi
    
    # Test health endpoint
    echo -n "Testing health endpoint... "
    if curl -f http://${ENDPOINT}/api/v1/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        exit 1
    fi
    
    # Test KSI endpoint
    echo -n "Testing KSI validation endpoint... "
    if curl -f -X POST http://${ENDPOINT}/api/v1/ksi/validate \
        -H "Content-Type: application/json" \
        -d '{"csoId": "test-001"}' > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
    fi
    
    echo -e "${GREEN}✓ Smoke tests passed${NC}"
}

# Main deployment flow
main() {
    check_prerequisites
    
    # Parse command line arguments
    DEPLOY_TARGET="${1:-docker}"
    
    case ${DEPLOY_TARGET} in
        "docker")
            build_and_push
            deploy_docker_compose
            ;;
        "kubernetes"|"k8s")
            build_and_push
            deploy_kubernetes
            ;;
        "build-only")
            build_and_push
            ;;
        *)
            echo -e "${RED}Unknown deployment target: ${DEPLOY_TARGET}${NC}"
            echo "Usage: $0 [docker|kubernetes|build-only]"
            exit 1
            ;;
    esac
    
    # Run migrations if not build-only
    if [[ "${DEPLOY_TARGET}" != "build-only" ]]; then
        run_migrations
        run_smoke_tests
    fi
    
    echo -e "\n${GREEN}🚀 Deployment completed successfully!${NC}"
    echo -e "Access the API at: http://localhost:8080/api/v1/health"
    echo -e "Access the dashboard at: http://localhost:8080/dashboard/"
}

# Run main function
main "$@" 