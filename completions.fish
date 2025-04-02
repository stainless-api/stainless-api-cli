set -l subcommands projects.config.commits.create projects.config.branches.create projects.config.branches.merge builds.create builds.retrieve targets.artifacts.retrieve
complete -c stainless-api-cli --no-files \
  -n "not __fish_seen_subcommand_from $subcommands" \
  -a "$subcommands"

complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.config.commits.create" \
  -a "--project-name --branch --commit-message --allow-empty --openapi-spec --stainless-config"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.config.branches.create" \
  -a "--project-name --branch --branch-from"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.config.branches.merge" \
  -a "--project-name --from --into"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.create" \
  -a "--branch --config-commit --project --targets --+target"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.retrieve" \
  -a "--build-id"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from targets.artifacts.retrieve" \
  -a "--build-id --target-name"

 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from builds.create" \
   -l targets \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from builds.create" \
   -l +target \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from targets.artifacts.retrieve" \
   -l target-name \
   -ra "node typescript python go java kotlin ruby terraform cli"