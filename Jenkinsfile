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
                    echo "=== Cleaning up old containers ==="
                    
                    # Hentikan dan hapus container lama
                    docker-compose down --remove-orphans || true
                    
                    # Pastikan container dengan nama tersebut tidak ada
                    docker stop lms-backend 2>/dev/null || true
                    docker rm lms-backend 2>/dev/null || true
                    
                    echo "=== Starting new container ==="
                    docker-compose up -d --build --force-recreate
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