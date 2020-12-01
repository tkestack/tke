#! /usr/bin/env bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2019 Tencent. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use
# this file except in compliance with the License. You may obtain a copy of the
# License at
#
# https://opensource.org/licenses/Apache-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

umask 0022
unset IFS
unset OFS
unset LD_PRELOAD
unset LD_LIBRARY_PATH

export PATH='/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'

die()
{
    echo '[FAIL] Operation failed.' >&2
    exit 1
}

cd `dirname "${0}"`     || exit 1
cwd=`pwd`               || exit 1
file=`basename "${0}"`  || exit 1

me="${cwd}/${file}"
tmp="${me}.tmp"

rm -rf "${tmp}"                                                 || die
mkdir "${tmp}" && cd "${tmp}"                                   || die

tailNum=`sed -n '/^#real installing packages append below/{=;q;}' ${me}`
tailNum=$((tailNum +1))
tail -n +${tailNum} "${me}" >package.tgz    || die
tar -zxf package.tgz                        || die

./install.sh $@                             || die
cd "${cwd}"                                 || die
exit 0

#real installing packages append below
