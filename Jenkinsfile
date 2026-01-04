pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "lms-backend:latest"
        DEPLOY_DIR = "/opt/lms"  // Directory deploy di server
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

        stage('Prepare Deploy Directory') {
            steps {
                sh """
                    echo "=== Preparing Deploy Directory ==="
                    
                    # Copy semua file yang diperlukan ke deploy directory
                    echo "Copying files to ${DEPLOY_DIR}..."
                    
                    # Copy docker-compose.yml
                    cp docker-compose.yml ${DEPLOY_DIR}/
                    
                    # Copy .env jika ada di repo, atau gunakan yang sudah ada di server
                    if [ -f ".env.example" ]; then
                        echo "Copying .env.example to ${DEPLOY_DIR}/.env"
                        cp .env.example ${DEPLOY_DIR}/.env
                    fi
                    
                    # Copy file lain yang dibutuhkan
                    if [ -d "config" ]; then
                        cp -r config ${DEPLOY_DIR}/
                    fi
                    
                    echo "Files in ${DEPLOY_DIR}:"
                    ls -la ${DEPLOY_DIR}/
                """
            }
        }

        stage('Deploy Backend') {
            steps {
                sh """
                    echo "=== Starting Deployment ==="
                    
                    # Navigate ke deploy directory
                    cd ${DEPLOY_DIR}
                    echo "Current directory: \$(pwd)"
                    
                    # Validasi docker-compose.yml
                    echo "Validating docker-compose.yml..."
                    docker-compose config
                    
                    # Deploy
                    echo "Stopping existing containers..."
                    docker-compose down
                    
                    echo "Starting new containers..."
                    docker-compose up -d --build
                    
                    # Check status
                    echo "Checking container status..."
                    sleep 5
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