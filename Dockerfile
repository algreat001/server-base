# this not working - sample for b2b application

FROM node as b2b-frontend-build
ADD ./frontend /frontend
WORKDIR /frontend
RUN yarn install && yarn build

FROM golang as b2b-backend
ADD ./backend /go/src/bavsoft/b2b/
COPY --from=b2b-frontend-build /frontend/dist /go/src/bavsoft/b2b/web/static/
WORKDIR /go/src/bavsoft/b2b/
RUN export CGO_ENABLED=0 && make build
ENTRYPOINT ./apiserver
EXPOSE 8080