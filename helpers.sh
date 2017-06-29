
export CURL_HEADER=" -H \"User-Agent: Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0\""

[[ -z "${TUNE_SH_CACHE}"  ]] && export TUNE_SH_CACHE="/tmp/.tune.sh"

###############################################################################

cache-my-url(){
  local _URL_TO_CACHE="$1"

  local _URL_CACHE_FILE=$(echo "${_URL_TO_CACHE}" | sed 's/[:\/?=]/-/g')
  _URL_CACHE_FILE="${TUNE_SH_CACHE}/${_URL_CACHE_FILE}"

  mkdir -p "${TUNE_SH_CACHE}"
  local _URL_HTML=$(cat "${_URL_CACHE_FILE}" 2>/dev/null)

  if [[ -z "${_URL_HTML}" ]]; then
    echo "not cached, fetching this url"
    _URL_HTML=$(curl -skL ${CURL_HEADER} "${_URL_TO_CACHE}" | tee "${_URL_CACHE_FILE}")
  else
    echo "this url is cached at ${_URL_CACHE_FILE}; delete it if you think new results aren't coming"
  fi

  echo ${_URL_HTML}

}

pick-a-option-interaction(){
  local options="$@"
  [[ -z "${options}" ]] && echo "[error] no options found" && return 1

  local options_idx=0
  for _option in $(echo $options); do
    ((options_idx=options_idx+1))
    echo "[${options_idx}] ${_option}"
  done
  echo -n "which index: "
}

pick-a-option-input(){
  read option_idx
  [[ -z "${option_idx}" ]] && option_idx=0
  echo ${option_idx}
}

pick-a-option(){
  local selection_idx="${1}"
  local options="${@:2}"
  [[ -z "${options}" ]] && echo "[error] no options found" && return 1
  
  local options_idx=0
  for _option in $(echo $options); do
    ((options_idx=options_idx+1))
    [[ "$options_idx" != "$selection_idx" ]] && continue
    echo "$_option" && return 0
  done
  echo "[error] failed in picking option"
  return 1
}
