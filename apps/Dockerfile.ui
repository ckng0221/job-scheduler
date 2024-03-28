FROM node:alpine As builder

WORKDIR /app

COPY apps/ui/package*.json ./

RUN npm install

COPY apps/ui/ .

RUN npm run build

FROM node:alpine as production

WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
USER nextjs
COPY --from=builder /app/next.config.mjs .
COPY --from=builder /app/package.json .
COPY --from=builder --chown=nextjs:nodejs  /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=builder --chown=nextjs:nodejs /app/public ./public


EXPOSE 3000

CMD [ "node", "server.js" ]


# docker build -f ./apps/Dockerfile.ui . -t job-scheduler-ui
# docker run -p 3000:3000 jobs-scheduler-ui
