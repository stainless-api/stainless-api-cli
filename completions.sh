#!/usr/bin/env bash

# This also works for zsh: https://zsh.sourceforge.io/Doc/Release/Completion-System.html#Completion-System
_main()
{
    COMPREPLY=()

    local subcommands="projects.config.create_branch projects.config.create_commit projects.config.merge builds.retrieve builds.target.retrieve builds.target.artifacts.retrieve_source"

    if [[ "$COMP_CWORD" -eq 1 ]]
    then
      local cur="${COMP_WORDS[COMP_CWORD]}"
      COMPREPLY=( $(compgen -W "$subcommands" -- "$cur") )
      return
    fi

    local subcommand="${COMP_WORDS[1]}"
    local flags
    case "$subcommand" in
      projects.config.create_branch)
        flags="--project-name --branch --branch-from"
        ;;
      projects.config.create_commit)
        flags="--project-name --branch --commit-message --openapi-spec --stainless-config --allow-empty"
        ;;
      projects.config.merge)
        flags="--project-name --from --into"
        ;;
      builds.retrieve)
        flags="--build-id"
        ;;
      builds.target.retrieve)
        flags="--build-id --target-name"
        ;;
      builds.target.artifacts.retrieve_source)
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
      builds.target.retrieve)
        case "$prev" in
          --target-name)
            COMPREPLY=( $(compgen -W "node typescript python go java kotlin ruby terraform cli" -- $cur) )
            ;;
        esac
        ;;
      builds.target.artifacts.retrieve_source)
        case "$prev" in
          --target-name)
            COMPREPLY=( $(compgen -W "node typescript python go java kotlin ruby terraform cli" -- $cur) )
            ;;
        esac
        ;;
    esac
}
complete -F _main stainless-api-cli