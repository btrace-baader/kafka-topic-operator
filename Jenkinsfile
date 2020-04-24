def Name = 'btrace/kafka-topic-operator'
pipeline {
    agent {
        kubernetes {
            label "kafka-topic-operator"
            defaultContainer "jnlp"
            yaml """
apiVersion: v1
kind: Pod
metadata:
  namespace: jenkins
  labels:
    component: cicd
spec:
  # Use service account that can deploy to all namespaces
  serviceAccountName: jenkins
  nodeSelector:
    jenkins: bigworker
  tolerations:
    - key: "jenkins"
      operator: "Equal"
      value: "bigworker"
      effect: "NoSchedule"
  containers:
    - name: golang
      image: baader/golang-kubebuilder:1.12.5_2.3.0-alpine
      ImagePullPolicy: Always
      command:
        - cat
      tty: true
      resources:
        requests:
          memory: "128Mi"
          cpu: "256m"
    - name: udpate-manager
      image: xcnt/kubernetes-update-manager:stable
      tty: true
      command:
        - cat
      imagePullPolicy: Always
      resources:
        requests:
          memory: "128Mi"
          cpu: "100m"
    - name: kaniko
      image: gcr.io/kaniko-project/executor:debug-v0.16.0
      command:
        - /busybox/cat
      tty: true
      resources:
        requests:
          memory: "512Mi"
          cpu: "500m"
      volumeMounts:
        - name: docker-hub-credentials
          mountPath: /kaniko/.docker/
  volumes:
    - name: docker-hub-credentials
      secret:
        secretName: docker-hub-credentials
"""
        }
    }
     options {
        timeout(time: 25)
    }
    stages {
        stage('Setup') {
            steps {
                container('golang') {
                    initBuild()
                }
            }
        }

        stage ('Run Tests') {
            steps {
                container('golang') {
                    sh '''
                        make test
                    '''
                }
            }
        }

        stage('Build and push Docker Image') {
            steps {
                container(name: 'kaniko', shell: '/busybox/sh') {
                    sh 'cat /kaniko/.docker/config.json'
                    buildPublicRepoWithKaniko imageRepo: "${Name}"
                }
            }
        }

        stage('Deployment') {
             when {
                anyOf {
                    branch 'master'
                    branch 'develop'
                }
             }
            steps {
                container('udpate-manager') {
                     notifyUpdatePublicRepo imageRepo: "${Name}"
                }
            }
        }
    }
}
