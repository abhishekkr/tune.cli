
export CURL_HEADER=" -H \"User-Agent: Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0\""

[[ -z "${TUNE_SH_CACHE}"  ]] && export TUNE_SH_CACHE="/tmp/.tune.sh"

###############################################################################

xmllint-hrefs-value-list(){
  local _HREFS_LIST="$1"
  local _HREFS_VALUE_LIST=""

  for _href in $(echo $_HREFS_LIST); do
    local _href_value=$(echo $_href | awk -F'=' '{print $2}' | xargs 2>/dev/null )
    if [[ -z "${_HREFS_VALUE_LIST}" ]]; then
      _HREFS_VALUE_LIST="${_href_value}"
    else
      _HREFS_VALUE_LIST="${_HREFS_VALUE_LIST} ${_href_value}"
    fi
  done
  echo "${_HREFS_VALUE_LIST}"
}

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

url-params-sanitizer(){
  local _URL_PARAMS="$@"

  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '&' '%26')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '(' '%28')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr ')' '%29')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '[' '%5B')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr ']' '%5D')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '{' '%7B')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '}' '%7D')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr ',' '%2C')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr "\'" '%27')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '#' '%23')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '$' '%24')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '%' '%25')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr ':' '%3A')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr ';' '%3B')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '?' '%3F')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '/' '%2F')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '@' '%40')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '\\' '%5C')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '|' '%7C')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '[:blank:]' '+')
  _URL_PARAMS=$(echo "${_URL_PARAMS}" | tr '[:space:]' '+')

  echo "${_URL_PARAMS}"
}

pick-a-option-interaction(){
  local options="$1"
  local options_title="$2"
  [[ -z "${options}" ]] && echo "[error] no options found" && return 1

  local options_idx=0
  for _option in $(echo $options); do
    ((options_idx=options_idx+1))
    local option_title=$(pick-a-option "${options_idx}" "${options_title}")
    [[ -z "${options_title}" ]] && option_title="${_option}"
    [[ -z "${option_title}" ]] && option_title="${_option}"
    echo "[${options_idx}] ${option_title}"
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
