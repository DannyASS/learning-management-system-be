pipeline {
    agent any

    environment {
        IMAGE_NAME = "lms-backend:latest"
        APP_DIR = "/opt/lms-backend"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                sh '''
                  docker build -t ${IMAGE_NAME} .
                '''
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                  cd ${APP_DIR}
                  docker compose down
                  docker compose up -d
                '''
            }
        }
    }

    post {
        success {
            echo "✅ Deploy sukses"
        }
        failure {
            echo "❌ Deploy gagal"
        }
    }
}
