pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "lms-backend:latest"
        COMPOSE_DIR = "/opt/lms"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/DannyASS/learning-management-system-be',
                    credentialsId: 'github-lms'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE} ."
            }
        }

        stage('Deploy Backend') {
            steps {
                sh """
                    echo "=== Starting Deployment ==="
                    echo "Working directory: \$(pwd)"
                    echo "Docker Compose version: \$(docker-compose --version)"
                    
                    # Navigate to compose directory
                    cd ${COMPOSE_DIR}
                    echo "Now in: \$(pwd)"
                    
                    # Stop and remove containers
                    echo "Stopping existing containers..."
                    docker-compose down
                    
                    # Start new containers with build
                    echo "Starting new containers..."
                    docker-compose up -d --build
                    
                    # Check status
                    echo "Checking container status..."
                    sleep 3
                    docker-compose ps
                    
                    echo "=== Deployment Complete ==="
                """
            }
        }
    }
    
    post {
        success {
            echo "Pipeline succeeded! 🎉"
        }
        failure {
            echo "Pipeline failed! 😢"
        }
    }
}