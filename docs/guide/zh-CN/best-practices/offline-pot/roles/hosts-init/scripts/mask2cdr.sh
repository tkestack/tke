#!/bin/sh

if [ -z $1  ]; then
  echo "please input a netmask,example: 255.255.255.0!" && exit 1;
else
  mask=$1
fi

mask2cdr ()
{
   # Assumes there's no "255." after a non-255 byte in the mask
   local x=${1##*255.}
   set -- 0^^^128^192^224^240^248^252^254^ $(( (${#1} - ${#x})*2 )) ${x%%.*}
   x=${1%%$3*}
   echo $(( $2 + (${#x}/4) ))
}

main(){
  mask2cdr $mask
}
main
