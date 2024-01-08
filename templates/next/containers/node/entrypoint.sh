#!/bin/sh
docker-entrypoint.sh
npm i
npm run build
npm run dev
