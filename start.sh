while true
do
	clear
	go build -o api -buildvcs=false -ldflags="-s -w" && ./api :42069
	sleep 3
	
	clear
	echo "Restarting Oomph API in 1 second..."
	sleep 1
	clear
done
