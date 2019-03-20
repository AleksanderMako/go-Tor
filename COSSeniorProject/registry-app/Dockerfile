FROM node:latest
RUN mkdir /app
WORKDIR /app

COPY package.json /app
RUN npm install --save

COPY . /app
EXPOSE 4500