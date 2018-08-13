#启动命令所在目录
receive='/root/silkworm'

runall(){
    run
}

stopall(){
    stop
}

stop(){
    ps -ef | grep silkworm_linux_amd64_test | awk '{print $2}' | xargs sudo kill -9
}

run(){
    cd $receive
    nohup ./silkworm_linux_amd64_test 1>/dev/null 2>/dev/null &
}


case $1 in
    start)
        if [ "$2" == "silkworm" ]
        then
            run
            echo "run silkworm success!"
        fi
        if [ "$2" == "all" ]
        then
            runall
            echo "run all success!"
        fi
        ;;
    stop)
        if [ "$2" == "silkworm" ]
        then
            kill
            echo "stop silkworm success!"
        fi
        if [ "$2" == "all" ]
        then
            stopall
            echo "stop all success!"
        fi
        ;;
    restart)
        if [ "$2" == "silkworm" ]
        then
            stop
            sleep 2
            run
            echo "restart silkworm success!"
        fi
        if [ "$2" == "all" ]
        then
            stopall
            sleep 2
            runall
            echo "restart all success!"
        fi
        ;;
    *)
        echo "Usage: {start|stop|restart silkworm|all}"
        ;;
esac
exit 0