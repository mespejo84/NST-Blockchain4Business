# SCRIPT TO INSTALL MOST ENVIRONMENT TOOLS FOR USING COMPOSER
# AUTOR: MOISES ESPEJO MENA
# LAST MODIFIED: MAY 24, 2018

echo "INSTALL GIT"
sudo apt-get update
sudo apt-get install -y git
 
echo "INSTALL DOCKER"
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce

echo "FIX FOR DOCKER"
sudo usermod -a -G docker $USER
 
echo "INSTALL NODEJS"
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
echo "FIX FOR INSTALL GLOBAL PACKAGES AS NORMAL USER"
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
ex +'$s@$@\rexport PATH='${HOME_PATH}'/.npm-global/bin:$PATH' -cwq ~/.profile
source ~/.profile
 
echo "INSTALL DOCKER COMPOSE"
sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo curl -L https://raw.githubusercontent.com/docker/compose/1.18.0/contrib/completion/bash/docker-compose -o /etc/bash_completion.d/docker-compose
 
echo "INSTALL ESSENTIAL CLI TOOLS"
npm install -g composer-cli
echo "INSTALL COMPOSER REST SERVER"
npm install -g composer-rest-server
echo "INSTALL UTILITY FOR GENERATING APPLICATION"
npm install -g generator-hyperledger-composer
echo "INSTALL YEOMAN"
npm install -g yo
echo "INSTALL PLAYGROUND"
npm install -g composer-playground
