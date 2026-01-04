pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "lms-backend:latest"
        ENV_FILE_HOST = "/opt/lms/.env"
        ENV_FILE_CONTAINER = "/app/.env"
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
                    cd ${COMPOSE_DIR}
                    docker-compose down
                    docker-compose up -d
                """
            }
        }
    }
}
