pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "lms-backend:latest"
        COMPOSE_DIR = "/opt/lms"      // Folder tempat docker-compose.yml & .env
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
                sh '''
                    docker build -t ${DOCKER_IMAGE} .
                '''
            }
        }

        stage('Deploy Backend') {
            steps {
                // Jalankan docker compose di host
                sh """
                    cd ${COMPOSE_DIR}
                    docker compose down
                    docker compose up -d
                """
            }
        }
    }
}
