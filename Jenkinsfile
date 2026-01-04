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
                    echo "deployment stage!"
                    # Coba docker compose (v2), jika gagal coba docker-compose (v1)
                    if command -v docker-compose &> /dev/null; then
                        echo "compose 1!"
                        docker-compose down
                        docker-compose up -d --build
                    elif docker compose version &> /dev/null; then
                        echo "compose 2!"
                        docker compose down
                        docker compose up -d --build
                    else
                        echo "ERROR: docker-compose not found!"
                        exit 1
                    fi
                """
            }
        }
    }
}
