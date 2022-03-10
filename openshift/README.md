# Openshift
This folder includes all the customizations needed to run Kubeflow on Openshift. To simplify the process of updating Kubeflow to new releases we added all our changes and customizations in this folder. 

## Process of Updating Kubeflow
To keep updating Kubeflow version smoothly, we decided to include all the changes in this folder and simply copy this folder from one release to another. 
After the Kubeflow community releases a stable Tag, we pull this new Tag in our fork and create a new branch from it for development of Openshift stack. The name convention is as follows: <tagname-openshift> for example for Kubeflow Tag v1.4.0 a new branch called v1.4.0-openshift will be created. The first step in development will be to copy all the work from the previous release to this new branch. The following is a list of steps:
1. Pull down the new Kubeflow Tag for a stable release for example Kubeflow Tag v1.4.0
2. Create a new branch from that Tag following the naming convention <tagname-openshift> for example v1.4.0-openshift
3. Create a new branch from the newly created openshift branch 
4. Copy the folder called ```openshift``` from the previous release openshift branch for example from v1.3.0-openshift
5. Make changes on fixes then submit the PR


