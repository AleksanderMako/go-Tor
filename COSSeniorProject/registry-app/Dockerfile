FROM node:9
RUN mkdir /app
WORKDIR /app

COPY package.json /app
RUN npm install --save

COPY . /app
EXPOSE 4500