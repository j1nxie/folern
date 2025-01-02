default:
    @just --choose

everything:
    parallel --lb ::: \
        'FORCE_COLOR=1 cd backend && caddy run -c Caddyfile' \
        'docker compose -f docker-compose.dev.yml up -d' \
        'FORCE_COLOR=1 just backend' \
        'FORCE_COLOR=1 just web-dev'

fmt:
    #!/bin/sh
    just fmt-backend fmt-web

fmt-backend:
    #!/bin/sh
    cd backend && go fmt ./

backend:
    #!/bin/sh
    cd backend && wgo run .

fmt-web:
    #!/bin/sh
    cd web && pnpm lint

web-install:
    #!/bin/sh
    cd web && pnpm i

web-dev:
    #!/bin/sh
    cd web && pnpm dev
