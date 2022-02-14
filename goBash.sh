if [ $1 != rm ]; then
    docker build -f $1 -t $2 .
else
    ID=$(docker images | grep -w $2 | awk '{print $3}')
    echo 'are you sure on delete image' $ID?
    select reply in Y N
    do
        if [ $reply == Y ]; then
            echo "deleting"
            docker image $1 -f $ID
            break
        else
            echo "aborting"
            break
        fi
    done
    exit
fi