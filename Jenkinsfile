node {
     stage('Clone repository') {
         /* Let's make sure we have the repository cloned to our workspace */
         checkout scm
     }

    String dockerImage = "cjburchell/yasls"
    String branchName = env.BRANCH_NAME
    if(branchName == "master"){
        branchName = 'latest'
    }

    stage('Build image') {
        docker.build("${dockerImage}").tag("${branchName}")
    }

    stage ('Push image') {
        docker.withRegistry('https://390282485276.dkr.ecr.us-east-1.amazonaws.com', 'ecr:us-east-1:redpoint-ecr-credentials') {
            docker.image("${dockerImage}").push("${branchName}")
        }
    }
}