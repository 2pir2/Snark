# Snark
In order to use circom, it is recommended to use a linux system (I tried windows and some libraries do not work in windows).
## Installation
1. Install linux system. It is fine to stay in root. Otherwise create a linux account and make sure to grant the user sudo access.
2. Install Nodejs
    - sudo apt install nodejs npm -y
4. Install rust with the following command. (Current location /home/username)
    - curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
    - source $HOME/.cargo/env
    - those two commands allow you to install rust and use 'cargo' terminal command
5. Create a new directory where you can put the circom files. cd to that directory (Current location /home/username/Snark)
6. Install circom with following commands
    - git clone https://github.com/iden3/circom.git
    - cd circom
    - cargo build --release
    - echo 'export PATH="$HOME/circom/target/release:$PATH"' >> ~/.bashrc
    - source ~/.bashrc
7. Install Snarkjs (Current location /home/username)
    - npm install -g snarkjs

