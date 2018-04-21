FROM node:7.9.0-alpine

WORKDIR /usr/src/app

COPY ./package.json /usr/src/app
RUN npm install
COPY ./app.js /usr/src/app

USER node
EXPOSE 3000

ENTRYPOINT [ "npm", "start" ]
