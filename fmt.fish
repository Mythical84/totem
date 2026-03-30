#!/bin/fish
for x in */
	echo $x
	cd $x
	go fmt
	cd ..
end
