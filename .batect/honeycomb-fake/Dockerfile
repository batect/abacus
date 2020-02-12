FROM node:13.8.0

RUN mkdir -p /app
WORKDIR /app

COPY package.json .
COPY yarn.lock .

RUN yarn install

COPY * .

CMD ["yarn", "start"]
