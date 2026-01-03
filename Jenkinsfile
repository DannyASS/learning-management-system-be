pipeline {
    agent any

    environment {
        // Path ke file .env di server
        ENV_FILE = "/opt/lms/.env"
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
                # Build Docker image
                docker build -t lms-backend:latest .
                '''
            }
        }

        stage('Deploy Docker') {
            steps {
                sh '''
                docker stop lms-backend || true
                docker rm lms-backend || true
                docker run -d --name lms-backend --env-file /opt/lms/.env -p 8082:8080 lms-backend:latest
                '''
            }
        }
    }

    post {
        success {
            echo "Deployment berhasil! Backend jalan di port 8082"
        }
        failure {
            echo "Pipeline gagal. Cek log untuk error"
        }
    }
}
