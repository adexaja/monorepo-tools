#!/bin/sh
set -eu

if [ "$#" -ne 3 ]; then
  echo "usage: $0 DESTINATION PROJECT_NAME GO_MODULE_PATH" >&2
  exit 2
fi

destination=$1
project_name=$2
module_path=$3
template_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

case "$module_path" in
  /*|*' '*|*'	'*)
    echo "GO_MODULE_PATH must be a Go import path, for example github.com/you/my-project" >&2
    exit 2
    ;;
esac

case "$destination" in
  /*) target=$destination ;;
  *) target=$(pwd)/$destination ;;
esac

if [ -e "$target" ]; then
  echo "destination already exists: $target" >&2
  exit 1
fi

mkdir -p "$target"
rsync -a --exclude node_modules --exclude '.moon/cache' --exclude '.moon/toolchain' "$template_dir/files/" "$target/"

find "$target" -type f ! -name '.gitkeep' ! -name 'bun.lock' ! -name 'bun.lockb' -exec env LC_ALL=C sed -i.bak \
  -e "s|__PROJECT_NAME__|$project_name|g" \
  -e "s|__PROJECT_SLUG__|$(printf '%s' "$project_name" | tr '[:upper:] ' '[:lower:]-')|g" \
  -e "s|__GO_MODULE__|$module_path|g" {} \;
find "$target" -type f -name '*.bak' -delete

chmod +x "$target/scripts/check.sh"
echo "created $target"
