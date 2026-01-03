pipeline {
    agent any

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
                sh 'docker build -t lms-backend:latest .'
            }
        }

        stage('Deploy Docker') {
            steps {
                sh '''
                docker stop lms-backend || true
                docker rm lms-backend || true
                docker run -d --name lms-backend -p 8080:8080 lms-backend:latest
                '''
            }
        }
    }
}
