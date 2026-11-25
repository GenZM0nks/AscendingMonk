FROM mcr.microsoft.com/playwright:v1.65.1-noble

WORKDIR /playwright

COPY package*.json ./
COPY playwright.config.ts ./

RUN npm ci

COPY --from=tests / /playwright/test/e2e/

CMD ["npx", "playwright", "test"]