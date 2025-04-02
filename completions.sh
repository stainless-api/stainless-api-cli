#!/usr/bin/env bash

# This also works for zsh: https://zsh.sourceforge.io/Doc/Release/Completion-System.html#Completion-System
_main()
{
    COMPREPLY=()

    local subcommands="projects.config.commits.create projects.config.branches.create projects.config.branches.merge builds.create builds.retrieve targets.artifacts.retrieve"

    if [[ "$COMP_CWORD" -eq 1 ]]
    then
      local cur="${COMP_WORDS[COMP_CWORD]}"
      COMPREPLY=( $(compgen -W "$subcommands" -- "$cur") )
      return
    fi

    local subcommand="${COMP_WORDS[1]}"
    local flags
    case "$subcommand" in
      projects.config.commits.create)
        flags="--project-name --branch --commit-message --allow-empty --openapi-spec --stainless-config"
        ;;
      projects.config.branches.create)
        flags="--project-name --branch --branch-from"
        ;;
      projects.config.branches.merge)
        flags="--project-name --from --into"
        ;;
      builds.create)
        flags="--branch --config-commit --project --targets --+target"
        ;;
      builds.retrieve)
        flags="--build-id"
        ;;
      targets.artifacts.retrieve)
        flags="--build-id --target-name"
        ;;
      *)
        # Unknown subcommand
        return
        ;;
    esac

    local cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$COMP_CWORD" -eq 2 || $cur == -* ]] ; then
        COMPREPLY=( $(compgen -W "$flags" -- $cur) )
        return 0
    fi

    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    case "$subcommand" in
      builds.create)
        case "$prev" in
          --targets)
            COMPREPLY=( $(compgen -W "node typescript python go java kotlin ruby terraform cli" -- $cur) )
            ;;
          --+target)
            COMPREPLY=( $(compgen -W "node typescript python go java kotlin ruby terraform cli" -- $cur) )
            ;;
        esac
        ;;
      targets.artifacts.retrieve)
        case "$prev" in
          --target-name)
            COMPREPLY=( $(compgen -W "node typescript python go java kotlin ruby terraform cli" -- $cur) )
            ;;
        esac
        ;;
    esac
}
complete -F _main stainless-api-cli