set -l subcommands projects.config.create_branch projects.config.create_commit projects.config.merge builds.retrieve builds.target.retrieve builds.target.artifacts.retrieve_source
complete -c stainless-v0-cli --no-files \
  -n "not __fish_seen_subcommand_from $subcommands" \
  -a "$subcommands"

complete -c stainless-v0-cli --no-files \
  -n "__fish_seen_subcommand_from projects.config.create_branch" \
  -a "--project-name --branch --branch-from"
complete -c stainless-v0-cli --no-files \
  -n "__fish_seen_subcommand_from projects.config.create_commit" \
  -a "--project-name --branch --commit-message --openapi-spec --stainless-config --allow-empty"
complete -c stainless-v0-cli --no-files \
  -n "__fish_seen_subcommand_from projects.config.merge" \
  -a "--project-name --from --into"
complete -c stainless-v0-cli --no-files \
  -n "__fish_seen_subcommand_from builds.retrieve" \
  -a "--build-id"
complete -c stainless-v0-cli --no-files \
  -n "__fish_seen_subcommand_from builds.target.retrieve" \
  -a "--build-id --target-name"
complete -c stainless-v0-cli --no-files \
  -n "__fish_seen_subcommand_from builds.target.artifacts.retrieve_source" \
  -a "--build-id --target-name"

 complete -c stainless-v0-cli --no-files \
   -n "__fish_seen_subcommand_from builds.target.retrieve" \
   -l target-name \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-v0-cli --no-files \
   -n "__fish_seen_subcommand_from builds.target.artifacts.retrieve_source" \
   -l target-name \
   -ra "node typescript python go java kotlin ruby terraform cli"