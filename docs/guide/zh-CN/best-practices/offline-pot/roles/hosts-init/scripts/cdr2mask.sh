#!/bin/sh

if [ -z $1  ]; then
  echo "please input a cdr,example: 24!" && exit 1;
else
  cdr=$1
fi


cdr2mask ()
{
   # Number of args to shift, 255..255, first non-255 byte, zeroes
   set -- $(( 5 - ($1 / 8) )) 255 255 255 255 $(( (255 << (8 - ($1 % 8))) & 255 )) 0 0 0
   [ $1 -gt 1 ] && shift $1 || shift
   echo ${1-0}.${2-0}.${3-0}.${4-0}
}

main(){
  cdr2mask $cdr
}
main
