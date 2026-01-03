pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "lms-backend:latest"
        CONTAINER_NAME = "lms-backend"
        ENV_FILE = "/opt/lms/.env"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/username/lms-backend.git'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "docker build -t $DOCKER_IMAGE ."
            }
        }

        stage('Stop Old Container') {
            steps {
                sh "docker stop $CONTAINER_NAME || true"
                sh "docker rm $CONTAINER_NAME || true"
            }
        }

        stage('Run New Container') {
            steps {
                sh """
                docker run -d \\
                  --name $CONTAINER_NAME \\
                  --env-file $ENV_FILE \\
                  -p 8080:8080 \\
                  $DOCKER_IMAGE
                """
            }
        }
    }

    post {
        success {
            echo "✅ Backend deployed successfully!"
        }
        failure {
            echo "❌ Deployment failed!"
        }
    }
}
