pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "lms-backend:latest"
        ENV_FILE_HOST = "/opt/lms/.env"    // path .env di server host
        ENV_FILE_CONTAINER = "/app/.env"   // path .env di dalam container
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/DannyASS/learning-management-system-be',
                    credentialsId: 'github-lms'
            }
        }

        stage('Build Docker') {
            steps {
                sh '''
                    docker build -t ${DOCKER_IMAGE} .
                '''
            }
        }

        stage('Deploy Docker') {
            steps {
                sh '''
                    docker build -t ${DOCKER_IMAGE} .

                    docker run -d \
                        --name lms-backend-new \
                        -v ${ENV_FILE_HOST}:${ENV_FILE_CONTAINER} \
                        -p 8083:8080 \
                        ${DOCKER_IMAGE}

                    # tunggu backend ready
                    sleep 10

                    docker rm -f lms-backend || true
                    docker rename lms-backend-new lms-backend
                '''
            }
        }
    }
}
