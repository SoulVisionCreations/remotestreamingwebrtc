FROM nvidia/opengl:1.2-glvnd-runtime-ubuntu20.04
# update
ENV TZ=Asia/Kolkata \
    DEBIAN_FRONTEND=noninteractive

RUN apt -y update && apt -y upgrade && apt-get  update
# install go
RUN apt install -y unzip libnss3 python3-pip wget
RUN wget https://go.dev/dl/go1.19.2.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz
ENV PATH=${PATH}:/usr/local/go/bin
RUN go version
#RUN echo "PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
#RUN apt install -y golang
#RUN /bin/bash -c go version
# install chromedriver
RUN cd /tmp/
#RUN wget https://chromedriver.storage.googleapis.com/83.0.4103.39/chromedriver_linux64.zip
RUN wget https://chromedriver.storage.googleapis.com/106.0.5249.61/chromedriver_linux64.zip
RUN unzip chromedriver_linux64.zip
RUN mv chromedriver /usr/bin/chromedriver
# RUN chromedriver --version
#RUN apt install -y xauth
#RUN xauth add ip-172-31-19-126/unix:5  MIT-MAGIC-COOKIE-1  829d979ee71793875cc4f54783599841
#RUN xauth list
# install google chrome
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb 
RUN  apt install -y ./google-chrome-stable_current_amd64.deb
RUN google-chrome-stable --version
# install ffmpeg
RUN apt-get install -y ffmpeg
RUN apt install -y libnvidia-encode-515
#CMD ["google-chrome-stable",  "--no-sandbox", "https://threejs.org/examples/#webgl_animation_keyframes",  "--window-position=0,0", "--window-size=600,400", "-kiosk", " --use-gl=desktop"]
WORKDIR app
COPY . .
RUN go mod init avataar/webrtc
RUN go mod tidy
# ENTRYPOINT echo "PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc 
RUN apt-get install -y xdotool
ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=all
#ENV CUDA_VISIBLE_DEVICES=0,1
CMD ["go", "run", "matchmaker.go"]
