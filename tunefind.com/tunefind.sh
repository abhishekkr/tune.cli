source "${TUNE_SH_DIR}/helpers.sh"

export TUNEFIND_BASE_URI="https://www.tunefind.com"
export YOUTUBE_SEARCH_BASE_URI="https://www.youtube.com/results?search_query="

###############################################################################

tunefine-cache(){
  local TUNEFIND_CACHE="${1}"
  local TUNEFIND_FILENAME="${2}"
  local TUNEFIND_CONTENT="${@:3}"

  local TUNEFIND_FILEPATH="${TUNEFIND_CACHE}/${TUNEFIND_FILENAME}"

  [[ ! -d "${TUNEFIND_CACHE}" ]] && mkdir -p "${TUNEFIND_CACHE}"

  echo "${TUNEFIND_CONTENT}" > "${TUNEFIND_FILEPATH}"
  echo "${TUNEFIND_FILEPATH}"
}

tunefind-url(){
  local _MEDIA_TITLE="${1}"
  local _MEDIA_TYPE="${2}"

  local _TUNEFIND_URL=""

  _TITLE=$(echo  ${_TITLE} |  tr '[:upper:]' '[:lower:]' | tr ' ' '-')

  case ${_MEDIA_TYPE} in
    "tv")
      _TUNEFIND_URL="${TUNEFIND_BASE_URI}/show/${_MEDIA_TITLE}"
      ;;
    "muv")
      _TUNEFIND_URL="${TUNEFIND_BASE_URI}/movies/${_MEDIA_TITLE}"
      ;;
    "search")
      _TUNEFIND_URL="${TUNEFIND_BASE_URI}/search/site?q=${_MEDIA_TITLE}"
      ;;
    *)
      echo "[ERROR] invalid media type" && return 1
      ;;  
  esac

  echo "${_TUNEFIND_URL}"
}

tunefind-tv-ispresent(){
  local _TV_SHOW="${1}"
  local _TV_SHOW_URL=$(tunefind-url "${_TV_SHOW}" "tv")

  local _HTTPCODE_TV_SHOW=$(curl -sIkL "${_TV_SHOW_URL}" | head -1 | cut -d' ' -f2)
  [[ "${_HTTPCODE_TV_SHOW:0:1}" -eq 2 ]] && return 0
  return 1
}

tunefind-muv-ispresent(){
  local _MUV="$1"
  local _MUV_SHOW=$(tunefind-url "${_TV_SHOW}" "muv")
  local _HTTPCODE_MUV_SHOW=$(curl -sIkL "${_MUV_SHOW}" | head -1 | cut -d' ' -f2)
  [[ "${_HTTPCODE_MUV_SHOW:0:1}" -eq 2 ]] && return 0
  return 1
}

tunefind-search(){
  local _SEARCH_FOR="$1"
  local _TUNEFIND_SEARCH_URL=$(tunefind-url "${_SEARCH_FOR}" "search")

  #[[]]
  #tunefind-cache "search--"$(echo "${_SEARCH_FOR}" | sed 's/[-_\ ]//g') 
  local _TUNEFIND_SEARCH_HTML=$( cache-my-url "${_TUNEFIND_SEARCH_URL}" )
  local _TUNEFIND_SEARCH_HREFS=$(echo "${_TUNEFIND_SEARCH_HTML}" | xmllint --html --xpath "(//div[@class='row tf-search-results']//a//@href)" - 2>/dev/null )

  for hrefs in $(echo ${_TUNEFIND_SEARCH_HREFS}); do
    local rel_href=$(echo "${hrefs}" | awk -F'=' '{print $2}' | xargs 2>/dev/null)
    local full_href="${TUNEFIND_BASE_URI}${rel_href}"
    echo "${full_href}"
  done
}

tunefind-search-songs(){
  local _SEARCH_FOR="$1"
  local _SEARCH_RESULTS=$(tunefind-search "${_SEARCH_FOR}")

  local _SED_ABLE_TUNEFIND_URL=$(echo "${TUNEFIND_BASE_URI}/" | sed -e 's/\//\\\//g' | sed -e 's/\:/\\\:/g')
  local _SEARCH_RESULTS_TITLES=$(echo "${_SEARCH_RESULTS}" | sed -e "s/${_SED_ABLE_TUNEFIND_URL}//g")

  pick-a-option-interaction "${_SEARCH_RESULTS}" "${_SEARCH_RESULTS_TITLES}"
  local _SEARCH_PICK_INDEX=$(pick-a-option-input)
  local _SEARCH_PICK=$(pick-a-option "${_SEARCH_PICK_INDEX}" "${_SEARCH_RESULTS}")

  if [[ $(echo "${_SEARCH_PICK}" | grep -c -E '^\[error\] ') -ne 0 ]]; then
    echo "************* you picked wrong option ************"
    tunefind-search-songs "${_SEARCH_FOR}"
    return 0
  fi
  
  echo "you picked: ${_SEARCH_PICK}"
  case "${_SEARCH_PICK}" in
    ${TUNEFIND_BASE_URI}/movie/*)
      tunefind-muv-seasons "${_SEARCH_PICK}"
      ;;
    ${TUNEFIND_BASE_URI}/show/*)
      tunefind-tv-seasons "${_SEARCH_PICK}"
      ;;
    ${TUNEFIND_BASE_URI}/artist/*)
      echo "WIP"
      ;;
    *)
      echo "[error] unhandled pick"
      ;;
  esac
}

tunefind-tv-seasons(){
  local _TV_SHOW_URL="${1}"

  local _TUNEFIND_TV_HTML=$( cache-my-url "${_TV_SHOW_URL}" )

  local _TUNEFIND_TV_HREFS=$(echo "${_TUNEFIND_TV_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//ul[@aria-labelledby='season-dropdown']//a[@role='menuitem']//@href)" - 2>/dev/null )
  _TUNEFIND_TV_HREFS=$(xmllint-hrefs-value-list "${_TUNEFIND_TV_HREFS}" )

  local _TUNEFIND_TV_HREFS_TITLES=$(echo "${_TUNEFIND_TV_HREFS}" | sed 's/\/show\///g')

  pick-a-option-interaction "${_TUNEFIND_TV_HREFS}" "${_TUNEFIND_TV_HREFS_TITLES}"
  local _SEASON_PICK_INDEX=$(pick-a-option-input)
  local _SEASON_PICK=$(pick-a-option "${_SEASON_PICK_INDEX}" "${_TUNEFIND_TV_HREFS}")

    local rel_href=$(echo "${_SEASON_PICK}")
    local full_href="${TUNEFIND_BASE_URI}${rel_href}"
    eval "tunefind-tv-seasons-episodes \"${full_href}\""

  return
# for hrefs in $(echo ${_TUNEFIND_TV_HREFS}); do
#   local rel_href=$(echo "${hrefs}" | awk -F'=' '{print $2}' | xargs 2>/dev/null)
#   local full_href="${TUNEFIND_BASE_URI}${rel_href}"
#   eval "tunefind-tv-seasons-episodes \"${full_href}\""
# done
}

tunefind-tv-seasons-episodes(){
  local _TV_SHOW_SEASONS_URL="${1}"

  echo "${_TV_SHOW_SEASONS_URL}"

  local _TUNEFIND_TV_SEASONS_HTML=$( cache-my-url "${_TV_SHOW_SEASONS_URL}" )
  local _TUNEFIND_TV_SEASONS_HREFS=$(echo "${_TUNEFIND_TV_SEASONS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//li[@class='MainList__item___fZ13_']//h3[@class='EpisodeListItem__title___32XUR']//a//@href)" - 2>/dev/null )
  _TUNEFIND_TV_SEASONS_HREFS=$(xmllint-hrefs-value-list "${_TUNEFIND_TV_SEASONS_HREFS}" )

  _TUNEFIND_TV_SEASONS_HREFS_TITLES=$(echo "${_TUNEFIND_TV_SEASONS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//li[@class='MainList__item___fZ13_']//h3[@class='EpisodeListItem__title___32XUR']//a//text())" - 2>/dev/null )
  _TUNEFIND_TV_SEASONS_HREFS_TITLES=$(echo "${_TUNEFIND_TV_SEASONS_HREFS_TITLES}" | sed 's/&Acirc;&middot;//g' | sed 's/ /-/g' | sed -E 's/(S[0-9]*--E[0-9]*)/ \1/g' | sed 's/^\s*//' )
 
  pick-a-option-interaction "${_TUNEFIND_TV_SEASONS_HREFS}" "${_TUNEFIND_TV_SEASONS_HREFS_TITLES}"
  local _EPISODE_PICK_INDEX=$(pick-a-option-input)
  local _EPISODE_PICK=$(pick-a-option "${_EPISODE_PICK_INDEX}" "${_TUNEFIND_TV_SEASONS_HREFS}")

  local rel_href=$(echo "${_EPISODE_PICK}")
  local full_href="${TUNEFIND_BASE_URI}${rel_href}#songs"

  echo "-------------------------------------"
  eval "tunefind-tv-seasons-episodes-songs \"${full_href}\""

  return

  for hrefs in $(echo ${_TUNEFIND_TV_SEASONS_HREFS}); do
    local rel_href=$(echo "${hrefs}" | awk -F'=' '{print $2}' | xargs 2>/dev/null)
    local full_href="${TUNEFIND_BASE_URI}${rel_href}#songs"

    echo "-------------------------------------"
    eval "tunefind-tv-seasons-episodes-songs \"${full_href}\""
  done
}


tunefind-tv-seasons-episodes-songs(){
  local _TV_SHOW_SEASONS_SONGS_URL="${1}"

  local _TUNEFIND_TV_SEASONS_SONGS_HTML=$( cache-my-url "${_TV_SHOW_SEASONS_SONGS_URL}" )
  local _TUNEFIND_TV_SEASONS_SONGS_HREFS=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a//@href)" - 2>/dev/null )
  _TUNEFIND_TV_SEASONS_SONGS_HREFS=$( xmllint-hrefs-value-list "${_TUNEFIND_TV_SEASONS_SONGS_HREFS}" )

  local _TUNEFIND_TV_SEASONS_SONGS_HREFS_TITLE=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a//@title)" - 2>/dev/null )
  _TUNEFIND_TV_SEASONS_SONGS_HREFS_TITLE=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HREFS_TITLE}" | sed 's/title=//g' | tr '[:blank:]' '-' | tr '[:space:]' '-' | sed  's/"-"/" "/g' | sed 's/^-"/"/' | sed 's/"-$/"/')

  pick-a-option-interaction "${_TUNEFIND_TV_SEASONS_SONGS_HREFS}" "${_TUNEFIND_TV_SEASONS_SONGS_HREFS_TITLE}"
  local _SONG_PICK_INDEX=$(pick-a-option-input)
  local _SONG_PICK=$(pick-a-option "${_SONG_PICK_INDEX}" "${_TUNEFIND_TV_SEASONS_SONGS_HREFS}")

  echo "-------------------------------------"
  tunefind-tv-seasons-episodes-songs-detail "${_SONG_PICK}" "${_TUNEFIND_TV_SEASONS_SONGS_HTML}"
  return

  for hrefs in $(echo ${_TUNEFIND_TV_SEASONS_SONGS_HREFS}); do
    local full_href=$(echo "${hrefs}" | awk -F'=' '{print $2}' | xargs 2>/dev/null)

    local _TUNEFIND_TV_SEASONS_SONGS_TITLE=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${full_href}']//text())" - 2>/dev/null )

    local _TUNEFIND_TV_SEASONS_SONGS_ARTIST_NAME=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${full_href}']//..//..//..//a[@class='Tunefind__Artist']//text())" - 2>/dev/null)

    local _TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${full_href}']//..//..//..//a[@class='Tunefind__Artist']//@href)" - 2>/dev/null | awk -F'=' '{print $2}' | xargs 2>/dev/null)
    _TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL="${TUNEFIND_BASE_URI}${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL}"

    local _TUNEFIND_TV_SEASONS_SONGS_YOUTUBE=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${full_href}']//..//..//..//..//a[@class='StoreLinks__youtube___2MHoI']//@href)" - 2>/dev/null | cut -d'=' -f 2- | xargs 2>/dev/null)

    echo "${_TUNEFIND_TV_SEASONS_SONGS_TITLE}__(${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_NAME})" | sed 's/ /-/g'
    return

    full_href="${TUNEFIND_BASE_URI}${full_href}"
    echo "[*] ${_TUNEFIND_TV_SEASONS_SONGS_TITLE} - ${full_href}"
    echo " |-- ${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_NAME} ( ${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL} )"
    echo " |-- ${_TUNEFIND_TV_SEASONS_SONGS_YOUTUBE}"
    echo " '~~~ "
  done
}

tunefind-tv-seasons-episodes-songs-detail(){
  local _SONG_PICK="$1"
  local _TUNEFIND_TV_SEASONS_SONGS_HTML="${@:2}"

    local _TUNEFIND_TV_SEASONS_SONGS_TITLE=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${_SONG_PICK}']//text())" - 2>/dev/null )

    local _TUNEFIND_TV_SEASONS_SONGS_ARTIST_NAME=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${_SONG_PICK}']//..//..//..//a[@class='Tunefind__Artist']//text())" - 2>/dev/null)

    local _TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL=$(echo "${_TUNEFIND_TV_SEASONS_SONGS_HTML}" | xmllint --html --xpath "(//div[@class='Tunefind__Content']//div[@class='SongRow__center___1I0Cg']//h4[@class='SongTitle__heading___3kxXK']//a[@href='${_SONG_PICK}']//..//..//..//a[@class='Tunefind__Artist']//@href)" - 2>/dev/null | awk -F'=' '{print $2}' | xargs 2>/dev/null)
    _TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL="${TUNEFIND_BASE_URI}${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL}"

    local _TUNEFIND_TV_SEASONS_SONGS_YOUTUBE="${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_NAME} ${_TUNEFIND_TV_SEASONS_SONGS_TITLE}"
    _TUNEFIND_TV_SEASONS_SONGS_YOUTUBE=$(url-params-sanitizer "${_TUNEFIND_TV_SEASONS_SONGS_YOUTUBE}")
    _TUNEFIND_TV_SEASONS_SONGS_YOUTUBE="${YOUTUBE_SEARCH_BASE_URI}${_TUNEFIND_TV_SEASONS_SONGS_YOUTUBE}"

    echo "[*] ${_TUNEFIND_TV_SEASONS_SONGS_TITLE}"
    echo " |-- ${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_NAME} ( ${_TUNEFIND_TV_SEASONS_SONGS_ARTIST_URL} )"
    echo " |-- ${_TUNEFIND_TV_SEASONS_SONGS_YOUTUBE}"
    echo " '~~~ "
}


###############################################################################

