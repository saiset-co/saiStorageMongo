#!/usr/bin/env bash

#export GOPATH="/home/kobs/MDocs/GOlang/go"
#$(pwd)
echo 'GOPATH:' $GOPATH

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

ret=""
function join_by { local IFS="$1"; shift; ret="$*"; }

package=$1
package_list=($(echo "$package" | tr '/' ' '))
exe_file=${package_list[*]:(-1)}
unset 'package_list[${#package_list[@]}-1]'
join_by / ${package_list[*]}
package_path=$ret
unset 'package_list[0]'
join_by / ${package_list[*]}
package_name=$ret
exe_file=($(echo "$exe_file" | tr '.' '\n'))
exe_name=${exe_file[0]}

platforms=(
    "linux/amd64"
    "linux/386"
    "windows/amd64"
    "windows/386"
    "darwin/amd64"
    "darwin/386"
)

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    echo "building for" $GOOS $GOARCH
    output_path=bin/$GOOS/$GOARCH/$package_name
    output_name=$output_path/$exe_name
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package

    for dependency in $(source $package_path/dependencies.sh)
    do
        cp ./$dependency ./$output_path/$dependency
    done

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
