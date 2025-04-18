#!/usr/bin/env bash

# This also works for zsh: https://zsh.sourceforge.io/Doc/Release/Completion-System.html#Completion-System
_main()
{
    COMPREPLY=()

    local subcommands="\
      projects.retrieve \
      projects.update \
      projects.branches.create \
      projects.branches.retrieve \
      builds.create \
      builds.retrieve \
      builds.list \
      build_target_outputs.list"

    if [[ "$COMP_CWORD" -eq 1 ]]
    then
      local cur="${COMP_WORDS[COMP_CWORD]}"
      COMPREPLY=( $(compgen -W "$subcommands" -- "$cur") )
      return
    fi

    local subcommand="${COMP_WORDS[1]}"
    local flags
    case "$subcommand" in
      projects.retrieve)
        flags="--project-name"
        ;;
      projects.update)
        flags="--project-name --display-name"
        ;;
      projects.branches.create)
        flags="--project --branch --branch-from"
        ;;
      projects.branches.retrieve)
        flags="--project --branch"
        ;;
      builds.create)
        flags="--project --revision --allow-empty --branch --commit-message --targets --+target"
        ;;
      builds.retrieve)
        flags="--build-id"
        ;;
      builds.list)
        flags="--project --branch --cursor --limit"
        ;;
      build_target_outputs.list)
        flags="--build-id --target --type --output"
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
      build_target_outputs.list)
        case "$prev" in
          --target)
            COMPREPLY=( $(compgen -W "node typescript python go java kotlin ruby terraform cli" -- $cur) )
            ;;
          --type)
            COMPREPLY=( $(compgen -W "source" -- $cur) )
            ;;
          --output)
            COMPREPLY=( $(compgen -W "url git" -- $cur) )
            ;;
        esac
        ;;
    esac
}
complete -F _main stainless-api-cli