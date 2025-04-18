set -l subcommands \
  openapi.retrieve \
  projects.update \
  projects.branches.create \
  projects.branches.retrieve \
  projects.snippets.create_request \
  builds.create \
  builds.retrieve \
  builds.list \
  build_target_outputs.list \
  webhooks.postman.create_notification
complete -c stainless-api-cli --no-files \
  -n "not __fish_seen_subcommand_from $subcommands" \
  -a "$subcommands"

complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from openapi.retrieve" \
  -a ""
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
  -n "__fish_seen_subcommand_from projects.snippets.create_request" \
  -a "--project-name --language --request.method --request.parameters.in --request.parameters.name --request.+parameter --request.path --version"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.create" \
  -a "--project --revision --allow-empty --branch --commit-message --targets --+target"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.retrieve" \
  -a "--build-id"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from builds.list" \
  -a "--project --branch --cursor --limit"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from build_target_outputs.list" \
  -a "--build-id --target --type --output"
complete -c stainless-api-cli --no-files \
  -n "__fish_seen_subcommand_from webhooks.postman.create_notification" \
  -a "--collection-id"

 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from projects.snippets.create_request" \
   -l language \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from projects.snippets.create_request" \
   -l request.parameters.in \
   -ra "path query header cookie"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from projects.snippets.create_request" \
   -l version \
   -ra "next latest_released"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from builds.create" \
   -l targets \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from builds.create" \
   -l +target \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from build_target_outputs.list" \
   -l target \
   -ra "node typescript python go java kotlin ruby terraform cli"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from build_target_outputs.list" \
   -l type \
   -ra "source"
 complete -c stainless-api-cli --no-files \
   -n "__fish_seen_subcommand_from build_target_outputs.list" \
   -l output \
   -ra "url git"