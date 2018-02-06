#!/bin/bash
pageNum=`gs -q -dNODISPLAY -c "($1) (r) file runpdfbegin pdfpagecount = quit"`
gs -dNumRenderingThreads=4 -dNOPAUSE -sDEVICE=pngalpha -dFirstPage=$3 -dLastPage=$3 -sOutputFile=$2$3.png -r200 -q $1 -c quit