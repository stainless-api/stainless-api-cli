set -l subcommands \
  projects.retrieve \
  projects.update \
  projects.branches.create \
  projects.branches.retrieve \
  builds.create \
  builds.retrieve \
  builds.list \
  build_target_outputs.retrieve
complete -c stainless-api-cli --no-files \
  -n "not __fish_seen_subcommand_from $subcommands" \
  -a "$subcommands"

complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.retrieve" \
  -a "--project-name"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.update" \
  -a "--project-name --display-name"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.branches.create" \
  -a "--project --branch --branch-from"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from projects.branches.retrieve" \
  -a "--project --branch"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.create" \
  -a "--project --revision --allow-empty --branch --commit-message --targets --+target"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.retrieve" \
  -a "--build-id"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.list" \
  -a "--project --branch --config-commit --cursor --limit"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from build_target_outputs.retrieve" \
  -a "--build-id --target --type --output"

 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from builds.create" \
   -l targets \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from builds.create" \
   -l +target \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from build_target_outputs.retrieve" \
   -l target \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from build_target_outputs.retrieve" \
   -l type \
   -ra "source"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from build_target_outputs.retrieve" \
   -l output \
   -ra "url git"