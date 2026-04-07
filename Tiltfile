# Install Node deps on change of package.json
local_resource(
  'pnpm',
  cmd='pnpm i', # Not using the make target to edit the lockfile
  deps=['package.json'],
)

# Rebuild frontend if source files change
local_resource(
  'frontend',
  cmd='make frontend',
  deps=['frontend'],
  resource_deps=['pnpm'],
)

# Update go.sum file and ensure modules are available
local_resource(
  'go-sum',
  cmd='go mod tidy',
  deps=['go.mod'],
)

# Rebuild and run Go webserver on code changes
local_resource(
  'server',
  cmd='make build',
  deps=[
    'go.mod',
    'main.go',
    'pkg',
  ],
  ignore=[
    'accounting',
    'src',
  ],
  serve_cmd='envrun -- ./accounting --listen=:59813',
  readiness_probe=probe(
    http_get=http_get_action(59813, path='/api/healthz'),
    initial_delay_secs=1,
  ),
  resource_deps=[
    'frontend',
    'go-sum',
  ],
)
