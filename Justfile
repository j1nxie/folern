default:
    @just --choose

fmt:
    #!/bin/sh
    just fmt-backend fmt-web

fmt-backend:
    #!/bin/sh
    cd backend && go fmt ./

fmt-web:
    #!/bin/sh
    cd web && pnpm lint

web-install:
    #!/bin/sh
    cd web && pnpm i

web-dev:
    #!/bin/sh
    cd web && pnpm dev
