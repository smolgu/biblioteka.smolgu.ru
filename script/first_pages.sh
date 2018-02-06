#!/bin/bash
gs -sDEVICE=pdfwrite -dNOPAUSE -dBATCH -dSAFER \
       -dFirstPage=1 -dLastPage=2 \
       -sOutputFile=$1 $2