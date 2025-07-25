#!/usr/bin/env bash

# simple bash autocomplete for krypton-cli
# cli will need to support list friendly commands
# so that this can be further automated.

_krypton_cli_completions()
{
  local cur
  # current completion index
  cur=${COMP_WORDS[COMP_CWORD]}

  case "$COMP_CWORD" in
    1)
    COMPREPLY=($(compgen -W 'auth dcm dsts es fs fds iot ss util' -- "$cur"));;
    2)
    # switch on the last completed word (module)
    case "${COMP_WORDS[$((COMP_CWORD-1))]}" in
      auth)
      COMPREPLY=($(compgen -W 'app_token auth_code device_code' -- "$cur"));;
      dcm)
      COMPREPLY=($(compgen -W 'get_device list_devices' -- "$cur"));;
      dsts)
      COMPREPLY=($(compgen -W 'keys' -- "$cur"));;
      iot)
      COMPREPLY=($(compgen -W 'device_message login subscribe' -- "$cur"));;
      es)
      COMPREPLY=($(compgen -W 'enroll enroll_and_wait create_enroll_token delete_enroll_token get_certificate get_device_token get_enroll_token get_status renew_enroll unenroll' -- "$cur"));;
      fs)
      COMPREPLY=($(compgen -W 'create_file get_download_url get_file_details get_upload_url' -- "$cur"));;
      fds)
      COMPREPLY=($(compgen -W 'create_file create_common_file get_download_url get_file_details get_upload_url list_files' -- "$cur"));;
      ss)
      COMPREPLY=($(compgen -W 'create_task get_task' -- "$cur"));;
      util)
      COMPREPLY=($(compgen -W 'retry_url upload_file wait_for_server' -- "$cur"));;
    esac;;
  esac
}
complete -F _krypton_cli_completions krypton-cli
