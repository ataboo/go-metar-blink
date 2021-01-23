## Setup Docker Build Env

1. use apt to install docker
- Add repo
- `sudo apt install -y docker-ce docker-ce-cli containerd.io qemu-user-static`
2. make sure to add user to docker group
- `sudo usermod -aG docker $USER`
3. Log out and in then run hello world
- `docker run hello-world`
4. `docker buildx ls` should include 'linux/arm/v6'
5. Clone `github.com:rpi-ws281x/rpi-ws281x-go.git`
6. Run `docker buildx build --platform linux/arm/v6  --tag ws2811-builder --file docker/app-builder/Dockerfile .` to create a build container
7. `APP="go-metar-blink"` 
8. `docker run --rm -v "$PWD":/usr/src/$APP --platform linux/arm/v6 -w /usr/src/$APP ws2811-builder:latest go build -o "$APP-armv6" -v`
