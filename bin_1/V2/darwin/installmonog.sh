brew update
brew install mongodb
sudo mkdir -p /data/db
sudo chown -R `id -un` /data/db
mongod